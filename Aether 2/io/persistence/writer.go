// Persistence > Writer
// This file collects all of the functions that write to the database. UI uses this for insertions, as well as the Fetcher.

package persistence

import (
	"aether-core/io/api"
	"fmt"
	// _ "github.com/mattn/go-sqlite3"
	"aether-core/backend/metrics"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"errors"
	"github.com/fatih/color"
	// "github.com/jmoiron/sqlx/types"
	"strconv"
	"strings"
	"time"
)

// Node is a non-communicating entity that holds the LastCheckin timestamps of each of the entities provided in the remote node. There is no way to send this data over to somebody, this is entirely local. There is also no batch processing because there is no situation in which you would need to insert multiple nodes at the same time (since you won't be connecting to multiple nodes simultaneously)

func InsertNode(n DbNode) error {
	err := insertNode(n)
	if err != nil {
		if strings.Contains(err.Error(), "Database was locked") {
			logging.Log(1, err)
			if strings.Contains(err.Error(), "Database was locked") {
				logging.Log(1, "This transaction was not committed because database was locked. We'll wait 10 seconds and retry the transaction.")
				time.Sleep(10 * time.Second)
				logging.Log(1, "Retrying the previously failed InsertNode transaction.")
				err2 := insertNode(n)
				if err2 != nil {
					logging.LogCrash(fmt.Sprintf("The second attempt to commit this data to the database failed. The first attempt had failed because the database was locked. The second attempt failed with the error: %s This database is corrupted. Quitting.", err2))
				} else { // If the reattempted transaction succeeds
					logging.Log(1, "The retry attempt of the failed transaction succeeded.")
				}
			}
		}
	}
	return nil
}
func insertNode(n DbNode) error {
	// fmt.Println("Node to be inserted:")
	// fmt.Printf("%#v\n", n)
	// TODO: Consider whether this needs a enforceNoEmptyIdentityFields or enforceNoEmptyRequiredFields
	if api.Fingerprint(globals.BackendConfig.GetNodeId()) == n.Fingerprint {
		return errors.New(fmt.Sprintf("The node ID that was attempted to be inserted is the SAME AS the local node's ID. This could be an attempted attack. Node ID of the remote: %s", n.Fingerprint))
	}
	tx, err := globals.DbInstance.Beginx()
	if err != nil {
		return err
	}
	_, err2 := tx.NamedExec(nodeInsert, n)
	if err2 != nil {
		return err2
	}
	err3 := tx.Commit()
	if err3 != nil {
		tx.Rollback()
		logging.Log(1, fmt.Sprintf("InsertNode encountered an error when trying to commit to the database. Error is: %s", err3))
		if strings.Contains(err3.Error(), "database is locked") {
			logging.Log(1, fmt.Sprintf("This database seems to be locked. We'll sleep 10 seconds to give it the time it needs to recover. This mostly happens when the app has crashed and there is a hot journal - and SQLite is in the process of repairing the database. THE DATA IN THIS TRANSACTION WAS NOT COMMITTED. PLEASE RETRY."))
			return errors.New("Database was locked. THE DATA IN THIS TRANSACTION WAS NOT COMMITTED. PLEASE RETRY.")
		}
		return err3
	}

	nodeAsMap := make(map[string]string)
	nodeAsMap["Fingerprint"] = string(n.Fingerprint)
	nodeAsMap["BoardsLastCheckin"] = strconv.Itoa(int(n.BoardsLastCheckin))
	nodeAsMap["ThreadsLastCheckin"] = strconv.Itoa(int(n.ThreadsLastCheckin))
	nodeAsMap["PostsLastCheckin"] = strconv.Itoa(int(n.PostsLastCheckin))
	nodeAsMap["VotesLastCheckin"] = strconv.Itoa(int(n.VotesLastCheckin))
	nodeAsMap["KeysLastCheckin"] = strconv.Itoa(int(n.KeysLastCheckin))
	nodeAsMap["TruststatesLastCheckin"] = strconv.Itoa(int(n.TruststatesLastCheckin))
	nodeAsMap["AddressesLastCheckin"] = strconv.Itoa(int(n.AddressesLastCheckin))
	metrics.CollateMetrics("NodeInsertionsSinceLastMetricsDbg", nodeAsMap)
	client, conn := metrics.StartConnection()
	defer conn.Close()
	metrics.SendMetrics(client)
	return nil
}

// InsertOrUpdateAddresses is the multi-entry of the core function InsertOrUpdateAddress. This is the only public API, and it should be used exclusively, because this is where we have the connection retry logic that we need.
func InsertOrUpdateAddresses(a *[]api.Address) {
	addresses := *a
	logging.Log(2, fmt.Sprintf("We got an insert or update address request for these addresses: %#v", addresses))
	for i, _ := range addresses {
		err := insertOrUpdateAddress(addresses[i])
		if err != nil {
			logging.Log(1, err)
			if strings.Contains(err.Error(), "Database was locked") {
				logging.Log(1, "This transaction was not committed because database was locked. We'll wait 10 seconds and retry the transaction.")
				time.Sleep(10 * time.Second)
				logging.Log(1, "Retrying the previously failed InsertOrUpdateAddresses transaction.")
				err2 := insertOrUpdateAddress(addresses[i])
				if err2 != nil {
					fmt.Println(err2)
					logging.LogCrash(fmt.Sprintf("The second attempt to commit this data to the database failed. The first attempt had failed because the database was locked. The second attempt failed with the error: %s This database is corrupted. Quitting.", err2))
				} else { // If the reattempted transaction succeeds
					logging.Log(1, "The retry attempt of the failed transaction succeeded.")
				}
			}
		}
	}
}

// insertOrUpdateAddress is the ONLY way to update an address in the database. Be very careful with this, careless use of this function can result in entry of untrusted data from the remotes into the local database. The only legitimate use of this is to put in the details of nodes that this local machine has personally connected to.
func insertOrUpdateAddress(a api.Address) error {
	addressPackAsInterface, err := APItoDB(a)
	if err != nil {
		return errors.New(fmt.Sprint(
			"Error raised from APItoDB function used in Batch insert. Error: ", err))
	}
	addressPack := addressPackAsInterface.(AddressPack)
	err2 := enforceNoEmptyIdentityFields(addressPack)
	if err2 != nil {
		// If this unit does have empty identity fields, we pass on adding it to the database.
		return err2
	}
	err7 := enforceNoEmptyRequiredFields(addressPack)
	if err7 != nil {
		// If this unit does have empty identity fields, we pass on adding it to the database.
		return err7
	}
	dbAddress := DbAddress{}
	dbSubprotocols := []DbSubprotocol{}
	dbJunctionItems := []DbAddressSubprotocol{} // Junction table.
	dbAddress = addressPack.Address             // We only have one address.
	for _, dbSubprot := range addressPack.Subprotocols {
		dbSubprotocols = append(dbSubprotocols, dbSubprot)
		jItem := generateAdrSprotJunctionItem(addressPack.Address, dbSubprot)
		dbJunctionItems = append(dbJunctionItems, jItem)
	}
	tx, err3 := globals.DbInstance.Beginx()
	if err3 != nil {
		logging.LogCrash(err3)
	}
	_, err4 := tx.NamedExec(getSQLCommand("dbAddressUpdate"), dbAddress)
	if err4 != nil {
		logging.LogCrash(err4)
	}
	if len(dbSubprotocols) > 0 {
		for _, dbSubprotocol := range dbSubprotocols {
			_, err5 := tx.NamedExec(getSQLCommand("dbSubprotocol"), dbSubprotocol)
			if err5 != nil {
				logging.LogCrash(err5)
			}
		}
	}
	if len(dbJunctionItems) > 0 {
		for _, dbJunctionItem := range dbJunctionItems {
			_, err5 := tx.NamedExec(getSQLCommand("dbAddressSubprotocol"), dbJunctionItem)
			if err5 != nil {
				logging.LogCrash(err5)
			}
		}
	}
	err6 := tx.Commit()
	if err6 != nil {
		tx.Rollback()
		logging.Log(1, fmt.Sprintf("InsertOrUpdateAddress encountered an error when trying to commit to the database. Error is: %s", err6))
		if strings.Contains(err6.Error(), "database is locked") {
			logging.Log(1, fmt.Sprintf("This database seems to be locked. We'll sleep 10 seconds to give it the time it needs to recover. This mostly happens when the app has crashed and there is a hot journal - and SQLite is in the process of repairing the database. THE DATA IN THIS TRANSACTION WAS NOT COMMITTED. PLEASE RETRY."))
			return errors.New("Database was locked. THE DATA IN THIS TRANSACTION WAS NOT COMMITTED. PLEASE RETRY.")
		}
		return err6
	}
	return nil
}

func generateAdrSprotJunctionItem(addr DbAddress, sprot DbSubprotocol) DbAddressSubprotocol {
	var adrSprot DbAddressSubprotocol
	adrSprot.AddressLocation = addr.Location
	adrSprot.AddressSublocation = addr.Sublocation
	adrSprot.AddressPort = addr.Port
	adrSprot.SubprotocolFingerprint = sprot.Fingerprint
	return adrSprot
}

// This is where we capture DB errors like 'DB is locked' and take action, such as retrying.
func BatchInsert(apiObjects []interface{}) error {
	err := batchInsert(&apiObjects)
	if err != nil {
		if strings.Contains(err.Error(), "Database was locked") {
			logging.Log(1, "This transaction was not committed because database was locked. We'll wait 10 seconds and retry the transaction.")
			time.Sleep(10 * time.Second)
			logging.Log(1, "Retrying the previously failed BatchInsert transaction.")
			err2 := batchInsert(&apiObjects)
			if err2 != nil {
				logging.LogCrash(fmt.Sprintf("The second attempt to commit this data to the database failed. The first attempt had failed because the database was locked. The second attempt failed with the error: %s This database is corrupted. Quitting.", err2))
			} else { // If the reattempted transaction succeeds
				logging.Log(1, "The retry attempt of the failed transaction succeeded.")
			}
		}
	}
	return err
}

type batchBucket struct {
	DbBoards      []DbBoard
	DbThreads     []DbThread
	DbPosts       []DbPost
	DbVotes       []DbVote
	DbKeys        []DbKey
	DbTruststates []DbTruststate
	DbAddresses   []DbAddress
	// Sub objects
	// // Parent: Board
	DbBoardOwners         []DbBoardOwner
	DbBoardOwnerDeletions []DbBoardOwner
	// // Parent: Address
	// dbSubprotocols := []DbSubprotocol{}
	// WHY? Because this is untrusted address entry, and subprotocol info coming from the environment is not committed, alongside many other parts of the address data.
	// // Parent: Key
	DbCurrencyAddresses        []DbCurrencyAddress
	DbCurrencyAddressDeletions []DbCurrencyAddress
}

// TODO: Mind that any errors happening within the transaction, if they need to bail from the transaction, they need to close it! otherwise you get database is locked.
// TODO: Should this take a pointer instead? It's dealing with some big amounts of data.
// BatchInsert insert a set of objects in a batch as a transaction.
func batchInsert(apiObjectsPtr *[]interface{}) error {
	apiObjects := *apiObjectsPtr
	logging.Log(1, "Batch insert starting.")
	defer logging.Log(1, "Batch insert is complete.")
	numberOfObjectsCommitted := len(apiObjects)
	logging.Log(2, fmt.Sprintf("%v objects are being committed.", numberOfObjectsCommitted))
	start := time.Now()
	bb := batchBucket{}
	// For each API object, convert to DB object and add to transaction.
	for _, apiObject := range apiObjects {
		// apiObject: API type, dbObj: DB type.
		dbo, err := APItoDB(apiObject) // does not hit DB
		if err != nil {
			return errors.New(fmt.Sprint(
				"Error raised from APItoDB function used in Batch insert. Error: ", err))
		}
		err2 := enforceNoEmptyIdentityFields(dbo) // does not hit DB
		if err2 != nil {
			// If this unit does have empty identity fields, we pass on adding it to the database.
			logging.Log(2, err2)
			continue
		}
		err3 := enforceNoEmptyRequiredFields(dbo) // does not hit DB
		if err3 != nil {
			// If this unit does have empty identity fields, we pass on adding it to the database.
			logging.Log(2, err3)
			continue
		}
		switch dbObject := dbo.(type) {
		case BoardPack:
			if packShouldBeCommitted(dbObject) { // HITS DB
				bb.DbBoards = append(bb.DbBoards, dbObject.Board)
				// Get the list of board owners before the transaction.
				boardBoardOwnersBeforeTx, err := getBoardOwnersBeforeTx(dbObject.Board.Fingerprint) // HITS DB
				if err != nil {
					logging.LogCrash(err)
				}
				// Get the changelist.
				changelist := generateBoardOwnerChangelist( // Does not hit DB
					boardBoardOwnersBeforeTx, dbObject.BoardOwners)
				for boardOwner, keepBoardOwner := range changelist {
					if keepBoardOwner == true {
						// We keep the owner's existence. This can either be a creation or an update, SQL deals with that.
						bb.DbBoardOwners = append(bb.DbBoardOwners, boardOwner)
					} else {
						// The owner is deleted. So we remove it from the database.
						bb.DbBoardOwnerDeletions = append(bb.DbBoardOwnerDeletions, boardOwner)
					}
				}
			}
		case DbThread:
			bb.DbThreads = append(bb.DbThreads, dbObject)
		case DbPost:
			bb.DbPosts = append(bb.DbPosts, dbObject)
		case DbVote:
			bb.DbVotes = append(bb.DbVotes, dbObject)
		case AddressPack:
			// In case of address, we strip out everything except the primary keys. This is because we cannot trust the data that is coming from the network. We just add the primary key set, and the local node will take care of directly connecting to these nodes and getting the details.

			// The other types of address inputs are not affected by this because they use InsertOrUpdateAddress, not this batch insert. If you're batch inserting addresses, it's by definition third party data.

			// This also means that we will actually be not using the Subprotocols data, as that would be untrusted data.

			dbObject.Address.LocationType = 0 // IPv4 or 6
			dbObject.Address.Type = 0         // 2 = live, 255 = static
			dbObject.Address.LastOnline = 0   // We cannot trust someone else's last online timestamp
			dbObject.Address.ProtocolVersionMajor = 0
			dbObject.Address.ProtocolVersionMinor = 0
			dbObject.Address.ClientVersionMajor = 0
			dbObject.Address.ClientVersionMinor = 0
			dbObject.Address.ClientVersionPatch = 0
			dbObject.Address.ClientName = ""
			bb.DbAddresses = append(bb.DbAddresses, dbObject.Address)
		case KeyPack:
			if packShouldBeCommitted(dbObject) {
				bb.DbKeys = append(bb.DbKeys, dbObject.Key)
				currencyAddressesBeforeTx, err := getCurrencyAddressesBeforeTx(dbObject.Key.Fingerprint)
				if err != nil {
					logging.LogCrash(err)
				}
				// Get the changelist.
				changelist := generateCurrencyAddressChangelist(
					currencyAddressesBeforeTx, dbObject.CurrencyAddresses)
				for currencyAddress, keepcurrencyAddress := range changelist {
					if keepcurrencyAddress == true {
						// We keep the owner's existence. This can either be a creation or an update, SQL deals with that.
						bb.DbCurrencyAddresses = append(bb.DbCurrencyAddresses, currencyAddress)
					} else {
						// The owner is deleted. So we remove it from the database.
						bb.DbCurrencyAddressDeletions = append(bb.DbCurrencyAddressDeletions, currencyAddress)
					}
				}
			}
		case DbTruststate:
			bb.DbTruststates = append(bb.DbTruststates, dbObject)
		default:
			return errors.New(
				fmt.Sprintf(
					"This object type is something batch insert does not understand. Your object: %#v\n", dbObject))
		}
	}
	err := insert(&bb)
	if err != nil {
		return err
	}
	elapsed := time.Since(start)
	clr := color.New(color.FgCyan)
	logging.Log(1, clr.Sprintf("It took %v to insert %v objects. %s", elapsed.Round(time.Millisecond), numberOfObjectsCommitted, generateInsertLog(&bb)))
	committedToDb := len(bb.DbBoards) + len(bb.DbThreads) + len(bb.DbPosts) + len(bb.DbVotes) + len(bb.DbKeys) + +len(bb.DbTruststates) + len(bb.DbAddresses)
	if (committedToDb != numberOfObjectsCommitted) && numberOfObjectsCommitted == 1 {
		clr2 := color.New(color.FgRed)
		logging.Log(1, clr2.Sprintf("There is a discrepancy between the number of entities in the inbound package, and those that end up being committed. Inbound entities count: %d, Committed to DB: %d", numberOfObjectsCommitted, committedToDb))
		logging.Log(1, clr.Sprintf("Inbound entities: %#v", apiObjects))
	}
	if len(apiObjects) > 0 {
		metrics.CollateMetrics("ArrivedEntitiesSinceLastMetricsDbg", apiObjects)
		client, conn := metrics.StartConnection()
		defer conn.Close()
		metrics.SendMetrics(client)
	}
	return nil
}

func generateInsertLog(bb *batchBucket) string {
	str := "Type:"
	if len(bb.DbBoards) > 0 {
		str = str + fmt.Sprintf(" %d Boards", len(bb.DbBoards))
	}
	if len(bb.DbThreads) > 0 {
		str = str + fmt.Sprintf(" %d Threads", len(bb.DbThreads))
	}
	if len(bb.DbPosts) > 0 {
		str = str + fmt.Sprintf(" %d Posts", len(bb.DbPosts))
	}
	if len(bb.DbVotes) > 0 {
		str = str + fmt.Sprintf(" %d Votes", len(bb.DbVotes))
	}
	if len(bb.DbKeys) > 0 {
		str = str + fmt.Sprintf(" %d Keys", len(bb.DbKeys))
	}
	if len(bb.DbTruststates) > 0 {
		str = str + fmt.Sprintf(" %d Truststates", len(bb.DbTruststates))
	}
	if len(bb.DbAddresses) > 0 {
		str = str + fmt.Sprintf(" %d Untrusted Addresses", len(bb.DbAddresses))
	}
	// Disabled because they don't count as individual entities.
	// if len(bb.DbBoardOwners) > 0 {
	// 	str = str + fmt.Sprintf(" %d Board Owners", len(bb.DbBoardOwners))
	// }
	// if len(bb.DbCurrencyAddresses) > 0 {
	// 	str = str + fmt.Sprintf(" %d Currency Addresses", len(bb.DbCurrencyAddresses))
	// }
	if len(bb.DbBoards) == 0 &&
		len(bb.DbThreads) == 0 &&
		len(bb.DbPosts) == 0 &&
		len(bb.DbVotes) == 0 &&
		len(bb.DbKeys) == 0 &&
		len(bb.DbTruststates) == 0 &&
		len(bb.DbAddresses) == 0 {
		str = str + " Nothing."
	} else {
		str = str + ". Nothing else."
	}
	return str
}

func insert(batchBucket *batchBucket) error {
	// We have our final list of entries. Add these objects to DB and let DB deal with what is a new addition and what is an update.
	// (Hot code path.) Start transaction.
	bb := *batchBucket
	tx, err := globals.DbInstance.Beginx()
	if err != nil {
		logging.LogCrash(err)
	}
	if len(bb.DbBoards) > 0 {
		for _, dbBoard := range bb.DbBoards {
			_, err := tx.NamedExec(getSQLCommand("dbBoard"), dbBoard)
			if err != nil {
				logging.LogCrash(err)
			}
		}
	}
	if len(bb.DbThreads) > 0 {
		for _, dbThread := range bb.DbThreads {
			_, err := tx.NamedExec(getSQLCommand("dbThread"), dbThread)
			if err != nil {
				logging.LogCrash(err)
			}
		}
	}
	if len(bb.DbPosts) > 0 {
		for _, dbPost := range bb.DbPosts {
			_, err := tx.NamedExec(getSQLCommand("dbPost"), dbPost)
			if err != nil {
				logging.LogCrash(err)
			}
		}
	}
	if len(bb.DbVotes) > 0 {
		for _, dbVote := range bb.DbVotes {
			_, err := tx.NamedExec(getSQLCommand("dbVote"), dbVote)
			if err != nil {
				logging.LogCrash(err)
			}
		}
	}
	if len(bb.DbKeys) > 0 {
		for _, dbKey := range bb.DbKeys {
			_, err := tx.NamedExec(getSQLCommand("dbKey"), dbKey)
			if err != nil {
				logging.LogCrash(err)
			}
		}
	}
	if len(bb.DbTruststates) > 0 {
		for _, dbTruststate := range bb.DbTruststates {
			_, err := tx.NamedExec(getSQLCommand("dbTruststate"), dbTruststate)
			if err != nil {
				logging.LogCrash(err)
			}
		}
	}
	if len(bb.DbAddresses) > 0 {
		for _, dbAddress := range bb.DbAddresses {
			_, err := tx.NamedExec(getSQLCommand("dbAddress"), dbAddress)
			if err != nil {
				logging.LogCrash(err)
			}
		}
	}
	if len(bb.DbBoardOwners) > 0 {
		for _, dbBoardOwner := range bb.DbBoardOwners {
			_, err := tx.NamedExec(getSQLCommand("dbBoardOwner"), dbBoardOwner)
			if err != nil {
				logging.LogCrash(err)
			}
		}
	}
	if len(bb.DbBoardOwnerDeletions) > 0 {
		for _, dbBoardOwner := range bb.DbBoardOwnerDeletions {
			_, err := tx.NamedExec(getSQLCommand("dbBoardOwnerDeletion"), dbBoardOwner)
			if err != nil {
				logging.LogCrash(err)
			}
		}
	}
	if len(bb.DbBoardOwnerDeletions) > 0 {
		for _, dbBoardOwner := range bb.DbBoardOwnerDeletions {
			_, err := tx.NamedExec(getSQLCommand("dbBoardOwnerDeletion"), dbBoardOwner)
			if err != nil {
				logging.LogCrash(err)
			}
		}
	}
	if len(bb.DbCurrencyAddresses) > 0 {
		for _, dbCurrencyAddress := range bb.DbCurrencyAddresses {
			_, err := tx.NamedExec(getSQLCommand("dbCurrencyAddress"), dbCurrencyAddress)
			if err != nil {
				logging.LogCrash(err)
			}
		}
	}
	if len(bb.DbCurrencyAddressDeletions) > 0 {
		for _, dbCurrencyAddress := range bb.DbCurrencyAddressDeletions {
			_, err := tx.NamedExec(getSQLCommand("dbCurrencyAddressDeletions"), dbCurrencyAddress)
			if err != nil {
				logging.LogCrash(err)
			}
		}
	}
	err2 := tx.Commit()
	if err2 != nil {
		tx.Rollback()
		logging.Log(1, fmt.Sprintf("BatchInsert encountered an error when trying to commit to the database. Error is: %s", err2))
		if strings.Contains(err.Error(), "database is locked") {
			logging.Log(1, fmt.Sprintf("This database seems to be locked. We'll sleep 10 seconds to give it the time it needs to recover. This mostly happens when the app has crashed and there is a hot journal - and SQLite is in the process of repairing the database. THE DATA IN THIS TRANSACTION WAS NOT COMMITTED. PLEASE RETRY."))
			return errors.New("Database was locked. THE DATA IN THIS TRANSACTION WAS NOT COMMITTED. PLEASE RETRY.")
		}
		return err2
	}
	fmt.Println("batchInsert successfully completed.")
	return nil
}

func getSQLCommand(dbType string) string {
	var sqlstr string
	if dbType == "dbBoard" {
		sqlstr = boardInsert
	} else if dbType == "dbThread" {
		if globals.BackendConfig.GetDbEngine() == "mysql" {
			sqlstr = threadInsertMySQL
		} else if globals.BackendConfig.GetDbEngine() == "sqlite" {
			sqlstr = threadInsertSQLite
		} else {
			logging.LogCrash(fmt.Sprintf("Db Engine type not recognised."))
		}
	} else if dbType == "dbPost" {
		if globals.BackendConfig.GetDbEngine() == "mysql" {
			sqlstr = postInsertMySQL
		} else if globals.BackendConfig.GetDbEngine() == "sqlite" {
			sqlstr = postInsertSQLite
		} else {
			logging.LogCrash(fmt.Sprintf("Db Engine type not recognised."))
		}
	} else if dbType == "dbVote" {
		sqlstr = voteInsert
	} else if dbType == "dbKey" {
		sqlstr = keyInsert
	} else if dbType == "dbTruststate" {
		sqlstr = truststateInsert
	} else if dbType == "dbAddress" { // untrusted address
		if globals.BackendConfig.GetDbEngine() == "mysql" {
			sqlstr = addressInsertMySQL
		} else if globals.BackendConfig.GetDbEngine() == "sqlite" {
			sqlstr = addressInsertSQLite
		} else {
			logging.LogCrash(fmt.Sprintf("Db Engine type not recognised."))
		}
	} else if dbType == "dbAddressUpdate" { // trusted address
		sqlstr = addressUpdateInsert
	} else if dbType == "dbBoardOwner" {
		sqlstr = boardOwnerInsert
	} else if dbType == "dbBoardOwnerDeletion" {
		sqlstr = boardOwnerDelete
	} else if dbType == "dbSubprotocol" {
		sqlstr = subprotocolInsert
	} else if dbType == "dbCurrencyAddress" {
		sqlstr = currencyAddressInsert
	} else if dbType == "dbCurrencyAddressDeletion" {
		sqlstr = currencyAddressDelete
	} else if dbType == "dbAddressSubprotocol" {
		if globals.BackendConfig.GetDbEngine() == "mysql" {
			sqlstr = addressSubprotocolInsertMySQL
		} else if globals.BackendConfig.GetDbEngine() == "sqlite" {
			sqlstr = addressSubprotocolInsertSQLite
		} else {
			logging.LogCrash(fmt.Sprintf("Db Engine type not recognised."))
		}
	}
	return sqlstr
}

// packShouldBeCommitted hits DB
func packShouldBeCommitted(pack interface{}) bool {
	switch pack := pack.(type) {
	case BoardPack:
		// First, pull the item if it exists in the database.
		resp, err := ReadBoards(
			[]api.Fingerprint{pack.Board.Fingerprint}, 0, 0)
		if err != nil {
			// logging.LogCrash(err)
			return false
		}
		// fmt.Printf("%#v\n", resp)
		// If the response is empty, this is a new board. Insert.
		if len(resp) == 0 {
			return true
		} else if len(resp) > 0 {
			// If the response is not empty, then do the regular update check.
			extantBoard := resp[0]

			if pack.Board.LastUpdate > extantBoard.LastUpdate &&
				pack.Board.LastUpdate > extantBoard.Creation {

				return true
			} else {

				return false
			}
		}
	case KeyPack:
		resp, err := ReadKeys([]api.Fingerprint{pack.Key.Fingerprint}, 0, 0)
		if err != nil {
			return false
		}
		// If the response is empty, this is a new key. Insert.
		if len(resp) == 0 {
			return true
		} else if len(resp) > 0 {
			// If the response is not empty, then do the regular update check.
			extantKey := resp[0]
			if pack.Key.LastUpdate > extantKey.LastUpdate &&
				pack.Key.LastUpdate > extantKey.Creation {
				return true
			} else {
				return false
			}
		}

	}
	return false
}

// getBoardOwnersBeforeTx hits DB.

// getBoardOwnersBeforeTx is an internal function of the Writer. This gets the pre-transaction state of the board owners of the board that is being inserted.
func getBoardOwnersBeforeTx(boardFingerprint api.Fingerprint) ([]DbBoardOwner, error) {
	var boardBoardOwnersBeforeTx []DbBoardOwner
	// Fetch all Board BoardOwners of this board that is already in database.
	rowsOfBoardOwnersBeforeTx, err := globals.DbInstance.Queryx("SELECT * from BoardOwners WHERE BoardFingerprint = ?", boardFingerprint)
	if err != nil {
		return boardBoardOwnersBeforeTx, err
	}
	// Do the struct scan into the row.
	for rowsOfBoardOwnersBeforeTx.Next() {
		var bo DbBoardOwner
		err := rowsOfBoardOwnersBeforeTx.StructScan(&bo)
		if err != nil {
			return boardBoardOwnersBeforeTx, err
		}
		boardBoardOwnersBeforeTx = append(boardBoardOwnersBeforeTx, bo)
	}
	rowsOfBoardOwnersBeforeTx.Close()
	return boardBoardOwnersBeforeTx, nil
}

// getCurrencyAddressesBeforeTx is an internal function of the Writer. This gets the pre-transaction state of the currency addresses of the key that is being inserted.
func getCurrencyAddressesBeforeTx(keyFingerprint api.Fingerprint) ([]DbCurrencyAddress, error) {
	var currencyAddressesBeforeTx []DbCurrencyAddress
	// Fetch all currency addresses of this key that is already in database.
	rowsOfCurrencyAddressesBeforeTx, err := globals.DbInstance.Queryx("SELECT * from CurrencyAddresses WHERE KeyFingerprint = ?", keyFingerprint)
	if err != nil {
		return currencyAddressesBeforeTx, err
	}
	// Do the struct scan into the row.
	for rowsOfCurrencyAddressesBeforeTx.Next() {
		var ca DbCurrencyAddress
		err := rowsOfCurrencyAddressesBeforeTx.StructScan(&ca)
		if err != nil {
			return currencyAddressesBeforeTx, err
		}
		currencyAddressesBeforeTx = append(currencyAddressesBeforeTx, ca)
	}
	return currencyAddressesBeforeTx, nil
}

// generateBoardOwnerDeletionList creates the list which shows which board owners will have to be deleted.
func generateBoardOwnerChangelist(
	currentBoardOwners []DbBoardOwner,
	candidateBoardOwners []DbBoardOwner) map[DbBoardOwner]bool {
	set := make(map[DbBoardOwner]bool)
	// Add both current board owners and candidate owners into the set.
	for _, currentBoardOwner := range currentBoardOwners {
		// If not in the candidate list, will be removed.
		set[currentBoardOwner] = false
	}
	for _, candidateBoardOwner := range candidateBoardOwners {
		// Everything in the candidate list will be added by default, hence true.
		set[candidateBoardOwner] = true
	}
	return set
}

// generateCurrencyAddressChangelist creates the list which shows which currency addresses will have to be deleted.
func generateCurrencyAddressChangelist(
	currentCurrencyAddresses []DbCurrencyAddress,
	candidateCurrencyAddresses []DbCurrencyAddress) map[DbCurrencyAddress]bool {
	set := make(map[DbCurrencyAddress]bool)
	// Add both current and candidate currency addresses into the set.
	for _, currentCurrencyAddress := range currentCurrencyAddresses {
		// Make so that if not in the candidate list, will be removed.
		set[currentCurrencyAddress] = false
	}
	for _, candidateCurrencyAddress := range candidateCurrencyAddresses {
		// Everything in the candidate list will be added by default, hence true.
		set[candidateCurrencyAddress] = true
	}
	return set
}

// enforceNoEmptyIdentityFields enforces that nothing will enter the database without having proper identity columns. For most objects this is Fingerprint(s), for some, like address, it's a combination of multiple fields.
func enforceNoEmptyIdentityFields(object interface{}) error {
	switch obj := object.(type) {
	case BoardPack:
		if obj.Board.Fingerprint == "" {
			return errors.New(
				fmt.Sprintf(
					"This board has an empty primary key. BoardPack: %#v\n", obj))
		}
		for _, bo := range obj.BoardOwners {
			if bo.BoardFingerprint == "" || bo.KeyFingerprint == "" {
				return errors.New(
					fmt.Sprintf(
						"This board owner has one or more empty primary key(s). BoardPack: %#v\n", obj))
			}
		}
	case DbThread:
		if obj.Fingerprint == "" {
			return errors.New(
				fmt.Sprintf(
					"This thread has an empty primary key. Thread: %#v\n", obj))
		}
	case DbPost:
		if obj.Fingerprint == "" {
			return errors.New(
				fmt.Sprintf(
					"This post has an empty primary key. Post: %#v\n", obj))
		}
	case DbVote:
		if obj.Fingerprint == "" {
			return errors.New(
				fmt.Sprintf(
					"This vote has an empty primary key. Vote: %#v\n", obj))
		}

	case AddressPack:
		if obj.Address.Location == "" || obj.Address.Port == 0 {
			return errors.New(
				fmt.Sprintf(
					"This address has one or more empty primary key(s). Address: %#v\n", obj))
		}
	case KeyPack:
		if obj.Key.Fingerprint == "" {
			return errors.New(
				fmt.Sprintf(
					"This key has an empty primary key. KeyPack: %#v\n", obj))
		}
		for _, ca := range obj.CurrencyAddresses {
			if ca.KeyFingerprint == "" || ca.Address == "" {
				return errors.New(
					fmt.Sprintf(
						"This currency address has one or more empty primary key(s). KeyPack: %#v\n", obj))
			}

		}
	case DbTruststate:
		if obj.Fingerprint == "" {
			return errors.New(
				fmt.Sprintf(
					"This trust state has an empty primary key. Truststate: %#v\n", obj))
		}
	}
	return nil
}

// enforceNoEmptyRequiredFields enforces that nothing will enter the database without having proper required columns. What columns are required depends on the type of the entity. See documentation for details.
func enforceNoEmptyRequiredFields(object interface{}) error {
	// TODO: This needs to be able to also defend against unicode replacement char or unicode rune error characters, as well as fields that are somehow only composed of spaces. There have been occurrences in the past where people tried to get past this by editing their own local database. The local machine assumes zero trust, everything that is coming in needs to be fully checked for sanity.
	switch obj := object.(type) {
	case BoardPack:
		if obj.Board.Name == "" || obj.Board.Creation == 0 || obj.Board.ProofOfWork == "" {
			return errors.New(
				fmt.Sprintf(
					"This board has some required fields empty (One or more of: Name, Creation, PoW). BoardPack: %#v\n", obj))
		}
		for _, bo := range obj.BoardOwners {
			if bo.Level == 0 {
				return errors.New(
					fmt.Sprintf(
						"This boardowner has some required fields empty (One or more of: Level). BoardPack: %#v\n", obj))
			}
		}
	case DbThread:
		if obj.Board == "" || obj.Name == "" || obj.Creation == 0 || obj.ProofOfWork == "" {
			return errors.New(
				fmt.Sprintf(
					"This thread has some required fields empty (One or more of: Board, Name, Creation, PoW). Thread: %#v\n", obj))
		}
	case DbPost:
		if obj.Board == "" || obj.Thread == "" || obj.Parent == "" || string(obj.Body) == "" || obj.Creation == 0 || obj.ProofOfWork == "" {
			return errors.New(
				fmt.Sprintf(
					"This post has some required fields empty (One or more of: Board, Thread, Parent, Body, Creation, PoW). Post: %#v\n", obj))
		}
	case DbVote:
		if obj.Board == "" || obj.Thread == "" || obj.Target == "" || obj.Owner == "" || obj.Type == 0 || obj.Creation == 0 || obj.Signature == "" || obj.ProofOfWork == "" {
			return errors.New(
				fmt.Sprintf(
					"This vote has some required fields empty (One or more of: Board, Thread, Target, Owner, Type, Creation, Signature, PoW). Vote: %#v\n", obj))
		}
	case AddressPack: // fix
		// if obj.Address.LocationType == 0 || obj.Address.LastOnline == 0 || obj.Address.ProtocolVersionMajor == 0 || obj.Address.ClientVersionMajor == 0 || obj.Address.ClientName == "" || len(obj.Subprotocols) < 1 {
		// 	return errors.New(
		// 		fmt.Sprintf(
		// 			"This address has some required fields empty (One or more of: LocationType, LastOnline, ProtocolVersionMajor, Subprotocols, ClientVersionMajor, ClientName). Address: %#v\n", obj))
		// }
		// for _, subprot := range obj.Subprotocols {
		// 	if subprot.Fingerprint == "" || subprot.Name == "" || subprot.VersionMajor == 0 || subprot.SupportedEntities == "" {
		// 		return errors.New(
		// 			fmt.Sprintf(
		// 				"This address' subprotocol has some required fields empty (One or more of: Fingerprint, Name, VersionMajor, SupportedEntities). Address: %#v\n Subprotocol: %#v\n", obj, subprot))
		// 	}
		// }

		/*
			Why only these? Address is special. When address traverses over the network, it is mostly emptied out, because the information it contains is untrustable - a remote might be maliciously replacing those fields to get the network to do its bidding.
			The only address entry that is trustable is gained first-person, that is, this node connects to the node on the address, and that direct connection can update this address entity with real, first-party data.
		*/

		if obj.Address.Location == "" || obj.Address.Port == 0 {
			return errors.New(
				fmt.Sprintf(
					"This address has some required fields empty (One or more of: Location, Port. Address: %#v", obj))
		}

	case KeyPack:
		if obj.Key.Type == "" || obj.Key.PublicKey == "" || obj.Key.Creation == 0 || obj.Key.ProofOfWork == "" || obj.Key.Signature == "" {
			return errors.New(
				fmt.Sprintf(
					"This key has some required fields empty (One or more of: Type, PublicKey, Creation, PoW, Signature). KeyPack: %#v\n", obj))
		}
		for _, ca := range obj.CurrencyAddresses {
			if ca.CurrencyCode == "" {
				return errors.New(
					fmt.Sprintf(
						"This currency address has some required fields empty (One or more of: CurrencyCode). KeyPack: %#v\n", obj))
			}
		}
	case DbTruststate:
		if obj.Target == "" || obj.Owner == "" || obj.Type == 0 || obj.Creation == 0 || obj.ProofOfWork == "" || obj.Signature == "" {
			return errors.New(
				fmt.Sprintf(
					"This trust state has some required fields empty (One or more of: Target, Owner, Type, Creation, PoW, Signature). Truststate: %#v\n", obj))
		}
	}
	return nil
}
