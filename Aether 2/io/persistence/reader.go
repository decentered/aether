// Persistence > Reader
// This file collects all functions that read from a database. The Server uses this API, as well as the UI.

package persistence

import (
	"aether-core/io/api"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

// ReadNode provides the ability to seek a specific node.
func ReadNode(fingerprint api.Fingerprint) (DbNode, error) {
	var n DbNode
	if len(fingerprint) > 0 {
		query, args, err := sqlx.In("SELECT * FROM Nodes WHERE Fingerprint IN (?);", fingerprint)
		if err != nil {
			return n, err
		}
		rows, err := globals.DbInstance.Queryx(query, args...)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return n, err
		}
		for rows.Next() {
			err := rows.StructScan(&n)
			if err != nil {
				return n, err
			}
		}
		rows.Close()
	}
	if len(n.Fingerprint) == 0 {
		return n, errors.New(fmt.Sprintf("The node you have asked for could not be found. You asked for: %s", fingerprint))
	}
	return n, nil
}

// enforceReadValidity enforces that, in a ReadX function (medium level API below), either a time range or a list of fingerprints are asked, and not both.
func enforceReadValidity(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) error {
	valid := false
	/*
		IF:
			A: BeginTS blank AND EndTS blank AND fingerprint filter extant: GOOD. Fingerprint search.

			B: BeginTS OR EndTS is blank (or both are not), and fingerprint filter blank: GOOD. One-way bounded  or two-way bounded time range search.

			C: All fields are blank: GOOD. Time range search where BeginTS is last cache generation timestamp or network head.

			D: Anything else: BAD.
	*/
	if (beginTimestamp == 0 && endTimestamp == 0 && len(fingerprints) > 0) ||
		((beginTimestamp != 0 || endTimestamp != 0) && len(fingerprints) == 0) ||
		(beginTimestamp == 0 && endTimestamp == 0 && len(fingerprints) == 0) {
		valid = true
	}
	if !valid {
		return errors.New(fmt.Sprintf("You can either search for a time range, or for fingerprint(s). You can't do both or neither at the same time - you have to do one. Asked fingerprints: %#v, BeginTimestamp: %s, EndTimestamp: %s", fingerprints, strconv.Itoa(int(beginTimestamp)), strconv.Itoa(int(endTimestamp))))
	}
	return nil
}

// sanitiseTimeRange validates and cleans the time range used in ReadX functions (medium level API below)
func sanitiseTimeRange(
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp, now api.Timestamp) (api.Timestamp, api.Timestamp, error) {
	// If there is no end timestamp, the end is now.
	if endTimestamp == 0 || endTimestamp > now {
		endTimestamp = now
	}
	// If the begin is newer than the end, flip. We haven't started to enforce limits yet, so the change here will be entirely coming from the remote.
	if beginTimestamp > endTimestamp {
		return beginTimestamp, endTimestamp, errors.New(fmt.Sprintf("Your BeginTimestamp is larger than your EndTimestamp. BeginTimestamp: %s, EndTimestamp: %s", strconv.Itoa(int(beginTimestamp)), strconv.Itoa(int(endTimestamp))))
	}
	// Internal processing starts.

	// If beginTimestamp is older than our last cache, start from the end of the last cache.
	if beginTimestamp < api.Timestamp(globals.BackendConfig.GetLastCacheGenerationTimestamp()) {
		beginTimestamp = api.Timestamp(globals.BackendConfig.GetLastCacheGenerationTimestamp())
		endTimestamp = now // Because, in thecase of begin 3 and end 5, begin going to 145000000 will make begin much bigger than end. Prevent that by moving the end also.
	}
	if beginTimestamp == 0 {
		// If there are no caches, lastCache will be 0 and this will return everything in the database. To prevent this, we limit the results to the duration of the network head.
		nhd := globals.BackendConfig.GetNetworkHeadDays()
		delta := time.Duration(nhd) * time.Hour * 24
		beginTimestamp = api.Timestamp(time.Now().Add(-delta).Unix())
	}
	//If beginTimestamp is in the future, return error.
	if beginTimestamp > now {
		return beginTimestamp, endTimestamp, errors.New(fmt.Sprintf("Your beginTimestamp is in the future. BeginTimestamp: %s, Now: %s", strconv.Itoa(int(beginTimestamp)), strconv.Itoa(int(now))))
	}
	// End of internal processing
	// After we do these things, if we end up with a begin timestamp that is newer than the end, the end timestamp will be 'now'. This can happen in the case where both the start and end timestamps are within the cached period.
	if beginTimestamp > endTimestamp {
		endTimestamp = now
	}
	return beginTimestamp, endTimestamp, nil
}

// Read is the high level API for DB reads. It provides filtering support. It can return multiple types if requested by the embeds.
func Read(
	entityType string, // boards, threads, posts, votes, addresses, keys, truststates
	fingerprints []api.Fingerprint,
	embeds []string,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) (api.Response, error) {

	var result api.Response
	now := api.Timestamp(time.Now().Unix())
	// Fingerprints search and start/end timestamp search are mutually exclusive. Make sure that is enforced.
	err := enforceReadValidity(fingerprints, beginTimestamp, endTimestamp)
	if err != nil {
		return result, err
	}
	var sanitisedBeginTimestamp api.Timestamp
	var sanitisedEndTimestamp api.Timestamp
	var err2 error
	if len(fingerprints) == 0 {
		// If this is a time range search:
		sanitisedBeginTimestamp, sanitisedEndTimestamp, err2 = sanitiseTimeRange(beginTimestamp, endTimestamp, now)
		if err2 != nil {
			return result, err2
		}
	} // If not a time range search, the ranges are 0 and 0, and we have fingerprints.

	// This thing below is for embeds. This is the container within which we fill the fingerprints for the item requested (board fps, etc.) as []api.Provable. It's used to do []api.Board to []api.Provable transition essentially.
	var provableArr []api.Provable
	// Now we switch based on the entity type.
	switch entityType {
	case "boards":
		entities, err := ReadBoards(fingerprints, sanitisedBeginTimestamp, sanitisedEndTimestamp)
		if err != nil {
			return result, err
		}
		result.Boards = entities
		// Convert the result to []api.Provable

		for i, _ := range entities {
			provableArr = append(provableArr, &entities[i])
		}

	case "threads":
		entities, err := ReadThreads(fingerprints, sanitisedBeginTimestamp, sanitisedEndTimestamp)
		if err != nil {
			return result, err
		}
		result.Threads = entities
		// Convert the result to []api.Provable
		for i, _ := range entities {
			provableArr = append(provableArr, &entities[i])
		}
	case "posts":
		entities, err := ReadPosts(fingerprints, sanitisedBeginTimestamp, sanitisedEndTimestamp)
		if err != nil {
			return result, err
		}
		result.Posts = entities
		// Convert the result to []api.Provable
		for i, _ := range entities {
			provableArr = append(provableArr, &entities[i])
		}

	case "votes":
		entities, err := ReadVotes(fingerprints, sanitisedBeginTimestamp, sanitisedEndTimestamp)
		if err != nil {
			return result, err
		}
		result.Votes = entities
		// Convert the result to []api.Provable
		for i, _ := range entities {
			provableArr = append(provableArr, &entities[i])
		}
	case "addresses":
		return result, errors.New(fmt.Sprint("You tried to supply an address into the high level Read API. This API only provides reads for entities that fulfil the api.Provable interface. Please use ReadAddress directly."))
	case "keys":
		entities, err := ReadKeys(fingerprints, sanitisedBeginTimestamp, sanitisedEndTimestamp)
		if err != nil {
			return result, err
		}
		result.Keys = entities
		// Convert the result to []api.Provable
		for i, _ := range entities {
			provableArr = append(provableArr, &entities[i])
		}
	case "truststates":
		entities, err := ReadTruststates(fingerprints, sanitisedBeginTimestamp, sanitisedEndTimestamp)
		if err != nil {
			return result, err
		}
		result.Truststates = entities
		// Convert the result to []api.Provable
		for i, _ := range entities {
			provableArr = append(provableArr, &entities[i])
		}
	}
	// We deal with filling the embedded fields. Embed handler has all the code for the different types of embeds.
	embedErr := handleEmbeds(provableArr, &result, embeds)
	if embedErr != nil {
		return result, embedErr
	}
	return result, nil
}

// Embeds APIs. These provide the supporting embeds for the primary calls at the medium level. This is consumed by the high level API above.

// Embed helpers used by the Read high level API above.

func existsInEmbed(asked string, embeds []string) bool {
	if len(embeds) > 0 {
		for i, _ := range embeds {
			if embeds[i] == asked {
				return true
			}
		}
	}
	return false
}

func handleEmbeds(entities []api.Provable, result *api.Response, embeds []string) error {
	// This holds the results of the first embed so we can add it to the main results before sending it into the keys. See the comment below for context.
	var firstEmbedCache []api.Provable
	if existsInEmbed("threads", embeds) {
		thr, err := ReadThreadEmbed(entities)
		if err != nil {
			return err
		}
		result.Threads = thr
		for i, _ := range thr {
			firstEmbedCache = append(firstEmbedCache, &thr[i])
		}

	}

	if existsInEmbed("posts", embeds) {
		posts, err := ReadPostEmbed(entities)
		if err != nil {
			return err
		}
		result.Posts = posts
		for i, _ := range posts {
			firstEmbedCache = append(firstEmbedCache, &posts[i])
		}
	}
	if existsInEmbed("votes", embeds) {
		votes, err := ReadVoteEmbed(entities)
		if err != nil {
			return err
		}
		result.Votes = votes
		for i, _ := range votes {
			firstEmbedCache = append(firstEmbedCache, &votes[i])
		}
	}
	// Keys being at the end is significant. This is always the last processed embed, because, for example, if you do Board with Threads embed, both boards and threads will have keys to refer to, and without the thread keys the data will be incomplete. So the keys need to read both the original data, and the data coming from any other embeds when providing the keys.

	//TODO: Embedding needs to be able to go five levels deep. Boards, (Threads, Posts, Votes, Keys, Truststates). This is useful when a node needs a copy of a board, and the whole board and anything it links can be dissected from another node and sent over.

	// But for the time being, let's treat the key as a special case, as if the key are not fully provided, the embeds with more than one layer don't work at all. The embedded objects will not be able to be validated otherwise. The five layer embed thing could be constructed from a series of queries but the absence of keys for the first embed is a serious problem.
	if existsInEmbed("keys", embeds) {
		keys, err := ReadKeyEmbed(entities, firstEmbedCache) // <- firstEmbedCache
		if err != nil {
			return err
		}
		result.Keys = keys
	}
	return nil
}

// Core read embeds.

// ReadThreadEmbed gets the threads linked from the entities provided.
// Only available for: Boards
func ReadThreadEmbed(entities []api.Provable) ([]api.Thread, error) {
	var arr []api.Thread
	var dbArr []DbThread
	var entityFingerprints []api.Fingerprint
	if len(entities) == 0 {
		logging.Log(1, fmt.Sprintf("The entities list given to the thread embed is empty."))
		return arr, nil
	}
	switch entity := entities[0].(type) {
	// entity: typed API object.
	case *api.Board:
		// Only defined for boards. No other entity has thread embeds.
		entity = entity // Stop complaining
		for i, _ := range entities {
			entityFingerprints = append(entityFingerprints, entities[i].GetFingerprint())
		}
		query, args, err := sqlx.In("SELECT DISTINCT * FROM Threads WHERE Board IN (?);", entityFingerprints)
		if err != nil {
			return arr, err
		}
		rows, err := globals.DbInstance.Queryx(query, args...)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbThread
			err := rows.StructScan(&entity)
			if err != nil {
				return []api.Thread{}, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close() // Close it ASAP, do not call any new DB queries while still in this.
		for _, entity := range dbArr {
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Thread))
		}
	}
	return arr, nil
}

// ReadPostEmbed gets the posts linked from the existing entities provided.
// Only available for: Threads
func ReadPostEmbed(entities []api.Provable) ([]api.Post, error) {
	var arr []api.Post
	var dbArr []DbPost
	var entityFingerprints []api.Fingerprint
	if len(entities) == 0 {
		logging.Log(1, fmt.Sprintf("The entities list given to the post embed is empty."))
		return arr, nil
	}
	switch entity := entities[0].(type) {
	// entity: typed API object.
	case *api.Thread:
		// Only defined for threads. No other entity has post embeds.
		entity = entity // Stop complaining
		for i, _ := range entities {
			entityFingerprints = append(entityFingerprints, entities[i].GetFingerprint())
		}
		query, args, err := sqlx.In("SELECT DISTINCT * FROM Posts WHERE Thread IN (?);", entityFingerprints)
		if err != nil {
			return arr, err
		}
		rows, err := globals.DbInstance.Queryx(query, args...)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbPost
			err := rows.StructScan(&entity)
			if err != nil {
				return []api.Post{}, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close() // Close it ASAP, do not call any new DB queries while still in this.
		for _, entity := range dbArr {
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Post))
		}
	}
	return arr, nil
}

// ReadVoteEmbed gets the votes linking to the entities provided.
// Only available for: Posts
func ReadVoteEmbed(entities []api.Provable) ([]api.Vote, error) {
	var arr []api.Vote
	var dbArr []DbVote
	var entityFingerprints []api.Fingerprint
	if len(entities) == 0 {
		logging.Log(1, fmt.Sprintf("The entities list given to the vote embed is empty."))
		return arr, nil
	}
	switch entity := entities[0].(type) {
	// entity: typed API object.
	case *api.Post:
		// Only defined for posts. No other entity has vote embeds.
		entity = entity // Stop complaining
		for i, _ := range entities {
			entityFingerprints = append(entityFingerprints, entities[i].GetFingerprint())
		}
		query, args, err := sqlx.In("SELECT DISTINCT * FROM Votes WHERE Target IN (?);", entityFingerprints)
		if err != nil {
			return arr, err
		}
		rows, err := globals.DbInstance.Queryx(query, args...)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbVote
			err := rows.StructScan(&entity)
			if err != nil {
				return []api.Vote{}, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close() // Close it ASAP, do not call any new DB queries while still in this.
		for _, entity := range dbArr {
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Vote))
		}
	}
	return arr, nil
}

// ReadKeyEmbed gets the keys linked from the existing entities provided.
// Only available for: Boards, Threads, Posts, Truststates
func ReadKeyEmbed(entities []api.Provable, firstEmbedCache []api.Provable) ([]api.Key, error) {
	var arr []api.Key
	var dbArr []DbKey
	var entityOwners []api.Fingerprint
	if len(entities) == 0 {
		logging.Log(1, fmt.Sprintf("The entities list given to the key embed is empty."))
		return arr, nil
	}
	entities = append(entities, firstEmbedCache...)
	for i, _ := range entities {
		switch entity := entities[i].(type) {
		// entity: typed API object.
		case *api.Board:
			entityOwners = append(entityOwners, entity.GetOwner())
			for j, _ := range entity.BoardOwners {
				entityOwners = append(entityOwners, entity.BoardOwners[j].KeyFingerprint)
			}
		case *api.Thread, *api.Post, *api.Truststate:
			entityOwners = append(entityOwners, entity.GetOwner())
		}
	}
	// The thing below is the same as read keys.
	query, args, err := sqlx.In("SELECT DISTINCT * FROM PublicKeys WHERE Fingerprint IN (?);", entityOwners)
	if err != nil {
		return arr, err
	}
	rows, err := globals.DbInstance.Queryx(query, args...)
	defer rows.Close() // In case of premature exit.
	if err != nil {
		return arr, err
	}
	for rows.Next() {
		var entity DbKey
		err := rows.StructScan(&entity)
		if err != nil {
			return []api.Key{}, err
		}
		dbArr = append(dbArr, entity)
	}
	rows.Close() // Close it ASAP, do not call any new DB queries while still in this.
	for _, entity := range dbArr {
		apiEntity, err := DBtoAPI(entity)
		if err != nil {
			// Log the problem and go to the next iteration without saving this one.
			logging.Log(1, err)
			continue
		}
		arr = append(arr, apiEntity.(api.Key))
	}
	return arr, nil
}

// Medium Level API. You should not use these directly. Use the high level API (above) because that one has the embed support and returns a proper api.Response object.

// ReadBoards reads threads from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadBoards(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]api.Board, error) {
	var arr []api.Board
	dbArr, err := ReadDbBoards(fingerprints, beginTimestamp, endTimestamp)
	if err != nil {
		return arr, err
	}
	for _, entity := range dbArr {
		apiEntity, err := DBtoAPI(entity)
		if err != nil {
			// Log the problem and go to the next iteration without saving this one.
			logging.Log(1, err)
			continue
		}
		arr = append(arr, apiEntity.(api.Board))
	}
	return arr, nil
}

// ReadDbBoards returns the search result in DB form. This layer has a few fields like LocalArrival and LastReferenced not exposed to the API layer that allows for internal decision making.
func ReadDbBoards(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]DbBoard, error) {

	var dbArr []DbBoard
	if len(fingerprints) > 0 { // Fingerprints array search.
		query, args, err := sqlx.In("SELECT * FROM Boards WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return dbArr, err
		}
		rows, err := globals.DbInstance.Queryx(query, args...)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return dbArr, err
		}
		for rows.Next() {
			var entity DbBoard
			err := rows.StructScan(&entity)
			if err != nil {
				return dbArr, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close() // Close it ASAP, do not call any new DB queries while still in this.
	} else { // Time range search
		// This should result in:
		// - Entities that has landed to local after the beginning and before the end
		rows, err := globals.DbInstance.Queryx("SELECT DISTINCT * from Boards WHERE (LocalArrival > ? AND LocalArrival < ? ) ORDER BY LocalArrival DESC", beginTimestamp, endTimestamp)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return dbArr, err
		}
		for rows.Next() {
			var entity DbBoard
			err := rows.StructScan(&entity)
			if err != nil {
				return dbArr, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close()
	}
	return dbArr, nil
}

// ReadThreads reads threads from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadThreads(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]api.Thread, error) {
	var arr []api.Thread
	dbArr, err := ReadDbThreads(fingerprints, beginTimestamp, endTimestamp)
	if err != nil {
		return arr, err
	}
	for _, entity := range dbArr {
		apiEntity, err := DBtoAPI(entity)
		if err != nil {
			// Log the problem and go to the next iteration without saving this one.
			logging.Log(1, err)
			continue
		}
		arr = append(arr, apiEntity.(api.Thread))
	}
	return arr, nil
}

func ReadDbThreads(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]DbThread, error) {
	var dbArr []DbThread
	if len(fingerprints) > 0 { // Fingerprints array search.
		query, args, err := sqlx.In("SELECT * FROM Threads WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return dbArr, err
		}
		rows, err := globals.DbInstance.Queryx(query, args...)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return dbArr, err
		}
		for rows.Next() {
			var entity DbThread
			err := rows.StructScan(&entity)
			if err != nil {
				return dbArr, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close()
	} else { // Time range search
		// This should result in:
		// - Entities that has landed to local after the beginning and before the end
		rows, err := globals.DbInstance.Queryx("SELECT DISTINCT * from Threads WHERE (LocalArrival > ? AND LocalArrival < ?) ORDER BY LocalArrival DESC", beginTimestamp, endTimestamp)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return dbArr, err
		}
		for rows.Next() {
			var entity DbThread
			err := rows.StructScan(&entity)
			if err != nil {
				return dbArr, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close()
	}
	return dbArr, nil
}

// ReadPosts reads posts from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadPosts(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]api.Post, error) {
	var arr []api.Post
	dbArr, err := ReadDbPosts(fingerprints, beginTimestamp, endTimestamp)
	if err != nil {
		return arr, err
	}
	for _, entity := range dbArr {
		apiEntity, err := DBtoAPI(entity)
		if err != nil {
			// Log the problem and go to the next iteration without saving this one.
			logging.Log(1, err)
			continue
		}
		arr = append(arr, apiEntity.(api.Post))
	}
	return arr, nil
}

func ReadDbPosts(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]DbPost, error) {
	var dbArr []DbPost
	if len(fingerprints) > 0 { // Fingerprints array search.
		query, args, err := sqlx.In("SELECT * FROM Posts WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return dbArr, err
		}
		rows, err := globals.DbInstance.Queryx(query, args...)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return dbArr, err
		}
		for rows.Next() {
			var entity DbPost
			err := rows.StructScan(&entity)
			if err != nil {
				return dbArr, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close()

	} else { // Time range search
		// This should result in:
		// - Entities that has landed to local after the beginning and before the end
		rows, err := globals.DbInstance.Queryx("SELECT DISTINCT * from Posts WHERE (LocalArrival > ? AND LocalArrival < ?) ORDER BY LocalArrival DESC", beginTimestamp, endTimestamp)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return dbArr, err
		}
		for rows.Next() {
			var entity DbPost
			err := rows.StructScan(&entity)
			if err != nil {
				return dbArr, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close()
	}
	return dbArr, nil
}

// ReadVotes reads votes from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadVotes(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]api.Vote, error) {
	var arr []api.Vote
	dbArr, err := ReadDbVotes(fingerprints, beginTimestamp, endTimestamp)
	if err != nil {
		return arr, err
	}
	for _, entity := range dbArr {
		apiEntity, err := DBtoAPI(entity)
		if err != nil {
			// Log the problem and go to the next iteration without saving this one.
			logging.Log(1, err)
			continue
		}
		arr = append(arr, apiEntity.(api.Vote))
	}
	return arr, nil
}

func ReadDbVotes(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]DbVote, error) {
	var dbArr []DbVote
	if len(fingerprints) > 0 { // Fingerprints array search.
		query, args, err := sqlx.In("SELECT * FROM Votes WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return dbArr, err
		}
		rows, err := globals.DbInstance.Queryx(query, args...)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return dbArr, err
		}
		for rows.Next() {
			var entity DbVote
			err := rows.StructScan(&entity)
			if err != nil {
				return dbArr, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close()
	} else { // Time range search
		// This should result in:
		// - Entities that has landed to local after the beginning and before the end
		rows, err := globals.DbInstance.Queryx("SELECT DISTINCT * from Votes WHERE (LocalArrival > ? AND LocalArrival < ?) ORDER BY LocalArrival DESC", beginTimestamp, endTimestamp)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return dbArr, err
		}
		for rows.Next() {
			var entity DbVote
			err := rows.StructScan(&entity)
			if err != nil {
				return dbArr, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close()
	}
	return dbArr, nil
}

func readDbAddressesBasicSearch(Location api.Location, Sublocation api.Location, Port uint16) (*[]DbAddress, error) {
	var dbArr []DbAddress
	if len(Location) > 0 && Port > 0 { // Regular address search.
		rows, err := globals.DbInstance.Queryx("SELECT * from Addresses WHERE Location = ? AND Sublocation = ? AND Port = ?", Location, Sublocation, Port)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return &dbArr, err
		}
		for rows.Next() {
			var entity DbAddress
			err := rows.StructScan(&entity)
			if err != nil {
				return &dbArr, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close()
	}
	return &dbArr, nil
}

func readDbAddressesFirstXResultsSearch(maxResults int, offset int, addrType uint8) (*[]DbAddress, error) {
	var dbArr []DbAddress
	if maxResults > 0 {
		// First X results search.
		var query string
		var err error
		// You have to provide a addrtype, if you search for 0, that will find the nodes you haven't connected yet.
		query = "SELECT * from Addresses WHERE AddressType = ? ORDER BY LocalArrival DESC LIMIT ? OFFSET ?"
		rows, err := globals.DbInstance.Queryx(query, addrType, maxResults, offset)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return &dbArr, err
		}
		for rows.Next() {
			var entity DbAddress
			err := rows.StructScan(&entity)
			if err != nil {
				return &dbArr, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close()
	}
	return &dbArr, nil
}

func readDbAddressesTimeRangeSearch(
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	offset int,
	timeRangeSearchType string) (*[]DbAddress, error) {
	var dbArr []DbAddress
	// Time range search
	// This should result in:
	// - Entities that has landed to local after the beginning and before the end
	// If the end timestamp is 0, it's assumed that endTs is right now.
	var endTs api.Timestamp
	if endTimestamp == 0 {
		endTs = api.Timestamp(time.Now().Unix())
	} else {
		endTs = endTimestamp
	}
	var rangeSearchColumn string
	// Options: all, connected.
	if timeRangeSearchType == "" || timeRangeSearchType == "all" {
		// Default case. Default is LocalArrival.
		rangeSearchColumn = "LocalArrival"
	} else if timeRangeSearchType == "connected" {
		rangeSearchColumn = "LastOnline"
	} else {
		return &dbArr, errors.New(fmt.Sprintf("You have provided an invalid time range search type. You provided: %s", timeRangeSearchType))
	}
	query := fmt.Sprintf("SELECT DISTINCT * from Addresses WHERE (%s > ? AND %s < ?) ORDER BY %s DESC", rangeSearchColumn, rangeSearchColumn, rangeSearchColumn)
	rows, err := globals.DbInstance.Queryx(query, beginTimestamp, endTs)
	defer rows.Close() // In case of premature exit.
	if err != nil {
		return &dbArr, err
	}
	for rows.Next() {
		var entity DbAddress
		err := rows.StructScan(&entity)
		if err != nil {
			return &dbArr, err
		}
		dbArr = append(dbArr, entity)
	}
	rows.Close()
	return &dbArr, nil
}

// ReadAddresses reads addresses from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadAddresses(
	Location api.Location,
	Sublocation api.Location,
	Port uint16,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	maxResults int, offset int, addrType uint8,
	timeRangeSearchType string) ([]api.Address, error) {
	var arr []api.Address
	dbArr, err := ReadDbAddresses(Location, Sublocation, Port, beginTimestamp, endTimestamp, maxResults, offset, addrType, timeRangeSearchType)
	if err != nil {
		return arr, err
	}
	for _, entity := range dbArr {
		apiEntity, err := DBtoAPI(entity)
		if err != nil {
			// Log the problem and go to the next iteration without saving this one.
			logging.Log(1, err)
			continue
		}
		arr = append(arr, apiEntity.(api.Address))
	}
	return arr, nil
}

func ReadDbAddresses(
	Location api.Location,
	Sublocation api.Location,
	Port uint16,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	maxResults int, offset int, addrType uint8,
	timeRangeSearchType string) ([]DbAddress, error) {
	/*
		There are three ways you can use this.
		1) Provide Location, sublocation, port, and nothing else = regular search
		2) Provide MaxResults, Maybe provide offset, addrType, and nothing else = first X results search
		3) Provide Begin, End timestamp, nothing else = time range search.
		None of these can be combined. Provide only one set - not a combination.

		timeRangeSearchType:
		Options: all, connected. It only affects time range search (#3.)
		Allows the caller to specify whether you want this search to be done based on localArrival (all addresses in the database is processed) or lastOnline (only the addresses the computer has personally connected to returns).
		If nothing is given (i.e. ""), it defaults to "all".
	*/
	var dbArr []DbAddress
	if len(Location) > 0 && Port > 0 && maxResults == 0 { // Regular address search.
		logging.Log(1, "This is an address search. Type: Basic.")
		dbArrPointer, err := readDbAddressesBasicSearch(Location, Sublocation, Port)
		dbArr = *dbArrPointer
		if err != nil {
			return dbArr, err
		}
	} else if maxResults > 0 && len(Location) == 0 && Port == 0 {
		// First X results search.
		logging.Log(2, "This is an address search. Type: First X results.")
		dbArrPointer, err := readDbAddressesFirstXResultsSearch(maxResults, offset, addrType)
		dbArr = *dbArrPointer
		if err != nil {
			return dbArr, err
		}
	} else if maxResults == 0 && len(Location) == 0 && Port == 0 { // Time range search
		// This should result in:
		// - Entities that has landed to local after the beginning and before the end
		// If the end timestamp is 0, it's assumed that endTs is right now.
		logging.Log(2, "This is an address search. Type: Time Range.")
		arrPointer, err := readDbAddressesTimeRangeSearch(beginTimestamp, endTimestamp, offset, timeRangeSearchType)
		dbArr = *arrPointer
		if err != nil {
			return dbArr, err
		}
	} else {
		// Invalid configuration coming from the address. Return error.
		return dbArr, errors.New("You have requested data from ReadAddresses in an invalid configuration. It can provide a) search by IP and Port, b) Return last X updated addresses, c) Return the addresses that were updated in a given time range. These cannot be combined. In other words,if you want to use any of the options, you need to zero out the inputs required for the other two. You have provided inputs for more than one option. ")
	}
	// if arrPointer != nil {
	// 	arr = *arrPointer
	// } else {
	// 	arr = []api.Address{}
	// }
	return dbArr, nil
}

// ReadKeys reads keys from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadKeys(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]api.Key, error) {
	var arr []api.Key
	dbArr, err := ReadDbKeys(fingerprints, beginTimestamp, endTimestamp)
	if err != nil {
		return arr, err
	}
	for _, entity := range dbArr {
		apiEntity, err := DBtoAPI(entity)
		if err != nil {
			// Log the problem and go to the next iteration without saving this one.
			logging.Log(1, err)
			continue
		}
		arr = append(arr, apiEntity.(api.Key))
	}
	return arr, nil
}

func ReadDbKeys(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]DbKey, error) {
	var dbArr []DbKey
	if len(fingerprints) > 0 { // Fingerprints array search.
		query, args, err := sqlx.In("SELECT * FROM PublicKeys WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return dbArr, err
		}
		rows, err := globals.DbInstance.Queryx(query, args...)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return dbArr, err
		}
		for rows.Next() {
			var entity DbKey
			err := rows.StructScan(&entity)
			if err != nil {
				return dbArr, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close()
	} else { // Time range search
		// This should result in:
		// - Entities that has landed to local after the beginning and before the end
		rows, err := globals.DbInstance.Queryx("SELECT DISTINCT * from PublicKeys WHERE (LocalArrival > ? AND LocalArrival < ?) ORDER BY LocalArrival DESC", beginTimestamp, endTimestamp)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return dbArr, err
		}
		for rows.Next() {
			var entity DbKey
			err := rows.StructScan(&entity)
			if err != nil {
				return dbArr, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close()
	}
	return dbArr, nil
}

// ReadTrustStates reads trust states from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadTruststates(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]api.Truststate, error) {
	var arr []api.Truststate
	dbArr, err := ReadDbTruststates(fingerprints, beginTimestamp, endTimestamp)
	if err != nil {
		return arr, err
	}
	for _, entity := range dbArr {
		apiEntity, err := DBtoAPI(entity)
		if err != nil {
			// Log the problem and go to the next iteration without saving this one.
			logging.Log(1, err)
			continue
		}
		arr = append(arr, apiEntity.(api.Truststate))
	}
	return arr, nil
}

func ReadDbTruststates(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]DbTruststate, error) {
	var dbArr []DbTruststate
	if len(fingerprints) > 0 { // Fingerprints array search.
		query, args, err := sqlx.In("SELECT * FROM Truststates WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return dbArr, err
		}
		rows, err := globals.DbInstance.Queryx(query, args...)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return dbArr, err
		}
		for rows.Next() {
			var entity DbTruststate
			err := rows.StructScan(&entity)
			if err != nil {
				return dbArr, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close()
	} else { // Time range search
		// This should result in:
		// - Entities that has landed to local after the beginning and before the end
		rows, err := globals.DbInstance.Queryx("SELECT DISTINCT * from Truststates WHERE (LocalArrival > ? AND LocalArrival < ?) ORDER BY LocalArrival DESC", beginTimestamp, endTimestamp)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			return dbArr, err
		}
		for rows.Next() {
			var entity DbTruststate
			err := rows.StructScan(&entity)
			if err != nil {
				return dbArr, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close()
	}
	return dbArr, nil
}

// The Reader functions that return DB instances, rather than API ones.

// ReadDBBoardOwners reads board owners from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.

// This is left as a single-select, not multiple, because it already supports returning multiple entities, and there is no demand for these to be fetched in bulk.
func ReadDBBoardOwners(BoardFingerprint api.Fingerprint,
	KeyFingerprint api.Fingerprint) ([]DbBoardOwner, error) {
	var arr []DbBoardOwner
	// If this query is without a key fingerprint (we want all addresses with that board fingerprint), change the query as such.
	if KeyFingerprint == "" {
		rows, err := globals.DbInstance.Queryx("SELECT * from BoardOwners WHERE BoardFingerprint = ?", BoardFingerprint)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			logging.LogCrash(err)
		}
		for rows.Next() {
			var boardOwner DbBoardOwner
			err := rows.StructScan(&boardOwner)
			if err != nil {
				logging.LogCrash(err)
			}
			arr = append(arr, boardOwner)
		}
		rows.Close()
	} else {
		rows, err := globals.DbInstance.Queryx("SELECT * from BoardOwners WHERE BoardFingerprint = ? AND KeyFingerprint = ?", BoardFingerprint, KeyFingerprint)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			logging.LogCrash(err)
		}
		for rows.Next() {
			var boardOwner DbBoardOwner
			err := rows.StructScan(&boardOwner)
			if err != nil {
				logging.LogCrash(err)
			}
			arr = append(arr, boardOwner)
		}
		rows.Close()
	}
	return arr, nil
}

// ReadDBSubprotocols reads the subprotocols of a given address from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.

func ReadDBSubprotocols(Location api.Location, Sublocation api.Location, Port uint16) ([]DbSubprotocol, error) {
	var fpArr []api.Fingerprint
	rows, err := globals.DbInstance.Queryx("SELECT * from AddressesSubprotocols WHERE AddressLocation = ? AND AddressSublocation = ? AND AddressPort = ?", Location, Sublocation, Port)
	defer rows.Close() // In case of premature exit.
	if err != nil {
		logging.LogCrash(err)
	}
	// Get the Subprotocol fingerprints from the junction table.
	for rows.Next() {
		var dbAddressSubprot DbAddressSubprotocol
		err := rows.StructScan(&dbAddressSubprot)
		if err != nil {
			logging.LogCrash(err)
		}
		fpArr = append(fpArr, dbAddressSubprot.SubprotocolFingerprint)
	}
	rows.Close()
	// For each fingerprint, get the matching subprotocol.
	var subprotArr []DbSubprotocol
	for _, val := range fpArr {
		rows, err := globals.DbInstance.Queryx("SELECT * from Subprotocols WHERE Fingerprint = ?", val)
		defer rows.Close() // In case of premature exit.
		if err != nil {
			logging.LogCrash(err)
		}

		for rows.Next() {
			var subprot DbSubprotocol
			err := rows.StructScan(&subprot)
			if err != nil {
				logging.LogCrash(err)
			}
			subprotArr = append(subprotArr, subprot)
		}
		rows.Close()
	}
	return subprotArr, nil
}
