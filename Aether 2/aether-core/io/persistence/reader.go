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
	"time"
)

// These are utility methods that need to read from the database for miscelleaneous purposes.

func LocalNodeIsMature() (bool, error) {
	var nrOfRows int
	err := DbInstance.Get(&nrOfRows, "SELECT count(1) FROM Nodes;")
	if err != nil {
		return false, errors.New(fmt.Sprintf("LocalNodeIsMature failed to get the number of rows from the Nodes database. Error: %#v", err))
	}
	if nrOfRows >= 3 {
		logging.Log(2, "A maturity check was requested. Local node is mature.")
		return true, nil
	}
	logging.Log(2, "A maturity check was requested. Local node is NOT mature.")
	return false, nil
}

// ReadNode provides the ability to seek a specific node.
func ReadNode(fingerprint api.Fingerprint) (DbNode, error) {
	var n DbNode
	if len(fingerprint) > 0 {
		query, args, err := sqlx.In("SELECT * FROM Nodes WHERE Fingerprint IN (?);", fingerprint)
		if err != nil {
			return n, err
		}
		rows, err := DbInstance.Queryx(query, args...)
		if err != nil {
			return n, err
		}
		for rows.Next() {
			err = rows.StructScan(&n)
			if err != nil {
				return n, err
			}
		}
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
	if beginTimestamp == 0 && endTimestamp == 0 && len(fingerprints) > 0 {
		// This is a fingerprints search.
		valid = true
	} else if (beginTimestamp != 0 || endTimestamp != 0) &&
		len(fingerprints) == 0 {
		// If begin and end timestamps are not both zero, this is a time range search.
		valid = true
	}
	if !valid {
		return errors.New(fmt.Sprintf("You can either search for a time range, or for fingerprint(s). You can't do both or neither at the same time - you have to do one. Asked fingerprints: %#v, BeginTimestamp: %s, EndTimestamp: %s", fingerprints, beginTimestamp, endTimestamp))
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
		return beginTimestamp, endTimestamp, errors.New(fmt.Sprintf("Your BeginTimestamp is larger than your EndTimestamp. BeginTimestamp: %s, EndTimestamp: %s", beginTimestamp, endTimestamp))
	}
	// Internal processing starts.

	// If beginTimestamp is older than our last cache, start from the end of the last cache. If there are no caches, lastCache will be 0 and this will return everything in the database.
	if beginTimestamp < api.Timestamp(globals.LastCacheGenerationTimestamp) {
		beginTimestamp = api.Timestamp(globals.LastCacheGenerationTimestamp)
		endTimestamp = now // Because, in thecase of begin 3 and end 5, begin going to 145000000 will make begin much bigger than end. Prevent that by moving the end also.
	}
	//If beginTimestamp is in the future, return error.
	if beginTimestamp > now {
		return beginTimestamp, endTimestamp, errors.New(fmt.Sprintf("Your beginTimestamp is in the future. BeginTimestamp: %s, Now: %s", beginTimestamp, now))
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
		result.AvailableTypes = append(result.AvailableTypes, "Boards")
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
		result.AvailableTypes = append(result.AvailableTypes, "Threads")
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
		result.AvailableTypes = append(result.AvailableTypes, "Posts")
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
		result.AvailableTypes = append(result.AvailableTypes, "Votes")
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
		result.AvailableTypes = append(result.AvailableTypes, "Keys")
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
		result.AvailableTypes = append(result.AvailableTypes, "Truststates")
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
		result.AvailableTypes = append(result.AvailableTypes, "Threads")
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
		result.AvailableTypes = append(result.AvailableTypes, "Posts")
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
		result.AvailableTypes = append(result.AvailableTypes, "Votes")
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
		result.AvailableTypes = append(result.AvailableTypes, "Keys")
	}
	return nil
}

// Core read embeds.

// ReadThreadEmbed gets the threads linked from the entities provided.
// Only available for: Boards
func ReadThreadEmbed(entities []api.Provable) ([]api.Thread, error) {
	var arr []api.Thread
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
		rows, err := DbInstance.Queryx(query, args...)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbThread
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
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
		rows, err := DbInstance.Queryx(query, args...)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbPost
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
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
		rows, err := DbInstance.Queryx(query, args...)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbVote
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
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

	// switch entity := entities[0].(type) {
	// // entity: typed API object.
	// case *api.Board:
	// 	entity = entity // So that the compiler does not complain about entity not being used.
	// 	for i, _ := range entities {
	// 		entityOwners = append(entityOwners, entities[i].GetOwner())
	// 		b := entities[i].(*api.Board)
	// 		for j, _ := range b.BoardOwners {
	// 			entityOwners = append(entityOwners, b.BoardOwners[j].KeyFingerprint)
	// 		}
	// 	}
	// case *api.Thread, *api.Post, *api.Truststate:
	// 	for i, _ := range entities {
	// 		entityOwners = append(entityOwners, entities[i].GetOwner())
	// 	}
	// }
	// The thing below is the same as read keys.
	query, args, err := sqlx.In("SELECT DISTINCT * FROM PublicKeys WHERE Fingerprint IN (?);", entityOwners)
	if err != nil {
		return arr, err
	}
	rows, err := DbInstance.Queryx(query, args...)
	if err != nil {
		return arr, err
	}
	for rows.Next() {
		var entity DbKey
		err = rows.StructScan(&entity)
		if err != nil {
			return arr, err
		}
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
	if len(fingerprints) > 0 { // Fingerprints array search.
		query, args, err := sqlx.In("SELECT * FROM Boards WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return arr, err
		}
		rows, err := DbInstance.Queryx(query, args...)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbBoard
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Board))
		}
	} else { // Time range search
		// This should result in:
		// - Entities that has landed to local after the beginning and before the end
		rows, err := DbInstance.Queryx("SELECT DISTINCT * from Boards WHERE (LocalArrival > ? AND LocalArrival < ?) ", beginTimestamp, endTimestamp)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbBoard
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Board))
		}
	}
	return arr, nil
}

// ReadThreads reads threads from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadThreads(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]api.Thread, error) {
	var arr []api.Thread
	if len(fingerprints) > 0 { // Fingerprints array search.
		query, args, err := sqlx.In("SELECT * FROM Threads WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return arr, err
		}
		rows, err := DbInstance.Queryx(query, args...)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbThread
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Thread))
		}
	} else { // Time range search
		// This should result in:
		// - Entities that has landed to local after the beginning and before the end
		rows, err := DbInstance.Queryx("SELECT DISTINCT * from Threads WHERE (LocalArrival > ? AND LocalArrival < ?) ", beginTimestamp, endTimestamp)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbThread
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
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

// func ReadThreads(Fingerprint api.Fingerprint) (
// 	[]api.Thread, error) {
// 	var arr []api.Thread
// 	rows, err := DbInstance.Queryx("SELECT * from Threads WHERE Fingerprint = ?", Fingerprint)
// 	if err != nil {
// 		logging.LogCrash(err)
// 	}
// 	for rows.Next() {
// 		var thread DbThread
// 		err = rows.StructScan(&thread)
// 		if err != nil {
// 			logging.LogCrash(err)
// 		}
// 		apiThread, err := DBtoAPI(thread)
// 		if err != nil {
// 			// Log the problem and go to the next iteration without saving this one.
// 			logging.Log(1, err)
// 			continue
// 		}
// 		arr = append(arr, apiThread.(api.Thread))
// 	}
// 	return arr, nil
// }

// ReadPosts reads posts from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadPosts(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]api.Post, error) {
	var arr []api.Post
	if len(fingerprints) > 0 { // Fingerprints array search.
		query, args, err := sqlx.In("SELECT * FROM Posts WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return arr, err
		}
		rows, err := DbInstance.Queryx(query, args...)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbPost
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Post))
		}
	} else { // Time range search
		// This should result in:
		// - Entities that has landed to local after the beginning and before the end
		rows, err := DbInstance.Queryx("SELECT DISTINCT * from Posts WHERE (LocalArrival > ? AND LocalArrival < ?) ", beginTimestamp, endTimestamp)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbPost
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
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

// ReadVotes reads votes from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadVotes(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]api.Vote, error) {
	var arr []api.Vote
	if len(fingerprints) > 0 { // Fingerprints array search.
		query, args, err := sqlx.In("SELECT * FROM Votes WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return arr, err
		}
		rows, err := DbInstance.Queryx(query, args...)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbVote
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Vote))
		}
	} else { // Time range search
		// This should result in:
		// - Entities that has landed to local after the beginning and before the end
		rows, err := DbInstance.Queryx("SELECT DISTINCT * from Votes WHERE (LocalArrival > ? AND LocalArrival < ?) ", beginTimestamp, endTimestamp)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbVote
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
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

// func ReadVotes(Fingerprint api.Fingerprint) ([]api.Vote, error) {
// 	var arr []api.Vote
// 	rows, err := DbInstance.Queryx("SELECT * from Votes WHERE Fingerprint = ?", Fingerprint)
// 	if err != nil {
// 		logging.LogCrash(err)
// 	}
// 	for rows.Next() {
// 		var vote DbVote
// 		err = rows.StructScan(&vote)
// 		if err != nil {
// 			logging.LogCrash(err)
// 		}
// 		apiVote, err := DBtoAPI(vote)
// 		if err != nil {
// 			// Log the problem and go to the next iteration without saving this one.
// 			logging.Log(1, err)
// 			continue
// 		}
// 		arr = append(arr, apiVote.(api.Vote))
// 	}
// 	return arr, nil
// }

// ReadAddresses reads addresses from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadAddresses(
	Location api.Location,
	Sublocation api.Location,
	Port uint16,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	maxResults int, offset int, addrType uint8) ([]api.Address, error) {
	/*
		There are three ways you can use this.
		1) Provide Location, sublocation, port, and nothing else = regular search
		2) Provide MaxResults, Maybe provide offset, addrType, and nothing else = first X results search
		3) Provide Begin, End timestamp, nothing else = time range search.
		None of these can be combined. Provide only one set - not a combination.
	*/
	// TODO: Split this into 3 functions, probably.
	var arr []api.Address
	if len(Location) > 0 && Port > 0 && maxResults == 0 { // Regular address search.
		rows, err := DbInstance.Queryx("SELECT * from Addresses WHERE Location = ? AND Sublocation = ? AND Port = ?", Location, Sublocation, Port)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbAddress
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Address))
		}
	} else if maxResults > 0 && len(Location) == 0 && Port == 0 {
		// First X results search.
		var query string
		var rows *sqlx.Rows
		var err error
		// You have to provide a addrtype, if you search for 0, that will find the nodes you haven't connected yet.
		query = "SELECT * from Addresses WHERE AddressType = ? ORDER BY LocalArrival DESC LIMIT ? OFFSET ?"
		rows, err = DbInstance.Queryx(query, addrType, maxResults, offset)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbAddress
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Address))
		}
	} else if maxResults == 0 && len(Location) == 0 && Port == 0 { // Time range search
		// This should result in:
		// - Entities that has landed to local after the beginning and before the end
		// If the end timestamp is 0, it's assumed that endTs is right now.
		var endTs api.Timestamp
		if endTimestamp == 0 {
			endTs = api.Timestamp(time.Now().Unix())
		}
		rows, err := DbInstance.Queryx("SELECT DISTINCT * from Addresses WHERE (LocalArrival > ? AND LocalArrival < ?) ", beginTimestamp, endTs)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbAddress
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Address))
		}
	} else {
		// Invalid configuration coming from the address. Return error.
		return arr, errors.New("You have requested data from ReadAddresses in an invalid configuration. It can provide a) search by IP and Port, b) Return last X updated addresses, c) Return the addresses that were updated in a given time range. These cannot be combined. In other words,if you want to use any of the options, you need to zero out the inputs required for the other two. You have provided inputs for more than one option. ")
	}
	return arr, nil
}

// func ReadAddresses(Location api.Location,
// 	Sublocation api.Location, Port uint16) ([]api.Address, error) {
// 	var arr []api.Address
// 	rows, err := DbInstance.Queryx("SELECT * from Addresses WHERE Location = ? AND Sublocation = ? AND Port = ?", Location, Sublocation, Port)
// 	if err != nil {
// 		logging.LogCrash(err)
// 	}
// 	for rows.Next() {
// 		var address DbAddress
// 		err = rows.StructScan(&address)
// 		if err != nil {
// 			logging.LogCrash(err)
// 		}
// 		apiAddress, err := DBtoAPI(address)
// 		if err != nil {
// 			// Log the problem and go to the next iteration without saving this one.
// 			logging.Log(1, err)
// 			continue
// 		}
// 		arr = append(arr, apiAddress.(api.Address))
// 	}
// 	return arr, nil
// }

// ReadKeys reads keys from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadKeys(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]api.Key, error) {
	var arr []api.Key
	if len(fingerprints) > 0 { // Fingerprints array search.
		query, args, err := sqlx.In("SELECT * FROM PublicKeys WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return arr, err
		}
		rows, err := DbInstance.Queryx(query, args...)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbKey
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Key))
		}
	} else { // Time range search
		// This should result in:
		// - Entities that has landed to local after the beginning and before the end
		rows, err := DbInstance.Queryx("SELECT DISTINCT * from PublicKeys WHERE (LocalArrival > ? AND LocalArrival < ?) ", beginTimestamp, endTimestamp)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbKey
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Key))
		}
	}
	return arr, nil
}

// func ReadKeys(Fingerprint api.Fingerprint) ([]api.Key, error) {
// 	var arr []api.Key
// 	rows, err := DbInstance.Queryx("SELECT * from PublicKeys WHERE Fingerprint = ?", Fingerprint)
// 	if err != nil {
// 		logging.LogCrash(err)
// 	}
// 	for rows.Next() {
// 		var key DbKey
// 		err = rows.StructScan(&key)
// 		if err != nil {
// 			logging.LogCrash(err)
// 		}
// 		apiKey, err := DBtoAPI(key)
// 		if err != nil {
// 			// Log the problem and go to the next iteration without saving this one.
// 			logging.Log(1, err)
// 			continue
// 		}
// 		arr = append(arr, apiKey.(api.Key))
// 	}
// 	return arr, nil
// }

// ReadTrustStates reads trust states from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadTruststates(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp) ([]api.Truststate, error) {
	var arr []api.Truststate
	if len(fingerprints) > 0 { // Fingerprints array search.
		query, args, err := sqlx.In("SELECT * FROM Truststates WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return arr, err
		}
		rows, err := DbInstance.Queryx(query, args...)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbTruststate
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Truststate))
		}
	} else { // Time range search
		// This should result in:
		// - Entities that has landed to local after the beginning and before the end
		rows, err := DbInstance.Queryx("SELECT DISTINCT * from Truststates WHERE (LocalArrival > ? AND LocalArrival < ?) ", beginTimestamp, endTimestamp)
		if err != nil {
			return arr, err
		}
		for rows.Next() {
			var entity DbTruststate
			err = rows.StructScan(&entity)
			if err != nil {
				return arr, err
			}
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Truststate))
		}
	}
	return arr, nil
}

// func ReadTruststates(Fingerprint api.Fingerprint) (
// 	[]api.Truststate, error) {
// 	var arr []api.Truststate
// 	rows, err := DbInstance.Queryx("SELECT * from Truststates WHERE Fingerprint = ?", Fingerprint)
// 	if err != nil {
// 		logging.LogCrash(err)
// 	}
// 	for rows.Next() {
// 		var truststate DbTruststate
// 		err = rows.StructScan(&truststate)
// 		if err != nil {
// 			logging.LogCrash(err)
// 		}
// 		apiTruststate, err := DBtoAPI(truststate)
// 		if err != nil {
// 			// Log the problem and go to the next iteration without saving this one.
// 			logging.Log(1, err)
// 			continue
// 		}
// 		arr = append(arr, apiTruststate.(api.Truststate))
// 	}
// 	return arr, nil
// }

// The Reader functions that return DB instances, rather than API ones.

// ReadDBCurrencyAddresses reads currency addresses from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.

// This is left as a single-select, not multiple, because it already supports returning multiple entities, and there is no demand for these to be fetched in bulk.
func ReadDBCurrencyAddresses(KeyFingerprint api.Fingerprint,
	Address string) ([]DbCurrencyAddress, error) {
	var arr []DbCurrencyAddress
	// If this query is without address (we want all addresses with that key fingerprint), change the query as such.
	if Address == "" {
		rows, err := DbInstance.Queryx("SELECT * from CurrencyAddresses WHERE KeyFingerprint = ?", KeyFingerprint)
		if err != nil {
			logging.LogCrash(err)
		}
		for rows.Next() {
			var currencyAddress DbCurrencyAddress
			err = rows.StructScan(&currencyAddress)
			if err != nil {
				logging.LogCrash(err)
			}
			arr = append(arr, currencyAddress)
		}
	} else {
		rows, err := DbInstance.Queryx("SELECT * from CurrencyAddresses WHERE KeyFingerprint = ? AND Address = ?", KeyFingerprint, Address)
		if err != nil {
			logging.LogCrash(err)
		}
		for rows.Next() {
			var currencyAddress DbCurrencyAddress
			err = rows.StructScan(&currencyAddress)
			if err != nil {
				logging.LogCrash(err)
			}
			arr = append(arr, currencyAddress)
		}
	}
	return arr, nil
}

// ReadDBBoardOwners reads board owners from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.

// This is left as a single-select, not multiple, because it already supports returning multiple entities, and there is no demand for these to be fetched in bulk.
func ReadDBBoardOwners(BoardFingerprint api.Fingerprint,
	KeyFingerprint api.Fingerprint) ([]DbBoardOwner, error) {
	var arr []DbBoardOwner
	// If this query is without a key fingerprint (we want all addresses with that board fingerprint), change the query as such.
	if KeyFingerprint == "" {
		rows, err := DbInstance.Queryx("SELECT * from BoardOwners WHERE BoardFingerprint = ?", BoardFingerprint)
		if err != nil {
			logging.LogCrash(err)
		}
		for rows.Next() {
			var boardOwner DbBoardOwner
			err = rows.StructScan(&boardOwner)
			if err != nil {
				logging.LogCrash(err)
			}
			arr = append(arr, boardOwner)
		}
	} else {
		rows, err := DbInstance.Queryx("SELECT * from BoardOwners WHERE BoardFingerprint = ? AND KeyFingerprint = ?", BoardFingerprint, KeyFingerprint)
		if err != nil {
			logging.LogCrash(err)
		}
		for rows.Next() {
			var boardOwner DbBoardOwner
			err = rows.StructScan(&boardOwner)
			if err != nil {
				logging.LogCrash(err)
			}
			arr = append(arr, boardOwner)
		}
	}
	return arr, nil
}
