// Persistence > Writer
// This file collects all of the functions that write to the database. UI uses this for insertions, as well as the Fetcher.

package persistence

import (
	"aether-core/io/api"
	"fmt"
	// _ "github.com/mattn/go-sqlite3"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"errors"
	"time"
)

// Node is a non-communicating entity that holds the LastCheckin timestamps of each of the entities provided in the remote node. There is no way to send this data over to somebody, this is entirely local. There is also no batch processing because there is no situation in which you would need to insert multiple nodes at the same time (since you won't be connecting to multiple nodes simultaneously)
func InsertNode(n DbNode) error {
	// fmt.Println("Node to be inserted:")
	// fmt.Printf("%#v\n", n)
	// TODO: Consider whether this needs a enforceNoEmptyIdentityFields or enforceNoEmptyRequiredFields
	if api.Fingerprint(globals.NodeId) == n.Fingerprint {
		return errors.New(fmt.Sprintf("The node ID that was attempted to be inserted is the SAME AS the local node's ID. This could be an attempted attack. Node ID of the remote: %s", n.Fingerprint))
	}
	tx, err := DbInstance.Beginx()
	if err != nil {
		return err
	}
	_, err2 := tx.NamedExec(nodeInsert, n)
	if err2 != nil {
		return err2
	}
	err3 := tx.Commit()
	if err3 != nil {
		return err3
	}
	return nil
}

// InsertOrUpdateAddresses is the multi-entry of the core function InsertOrUpdateAddress.
func InsertOrUpdateAddresses(a *[]api.Address) {
	addresses := *a
	for i, _ := range addresses {
		err := InsertOrUpdateAddress(addresses[i])
		if err != nil {
			logging.Log(1, err)
		}
	}
}

// InsertOrUpdateAddress is the ONLY way to update an address in the database. Be very careful with this, careless use of this function can result in entry of untrusted data from the remotes into the local database. The only legitimate use of this is to put in the details of nodes that this local machine has personally connected to.
func InsertOrUpdateAddress(a api.Address) error {
	dbA, err := APItoDB(a)
	if err != nil {
		return errors.New(fmt.Sprint(
			"Error raised from APItoDB function used in Batch insert. Error: ", err))
	}
	err2 := enforceNoEmptyIdentityFields(dbA)
	if err2 != nil {
		// If this unit does have empty identity fields, we pass on adding it to the database.
		return err2
	}
	// err3 := enforceNoEmptyRequiredFields(dbA)
	// if err3 != nil {
	// 	// If this unit does have empty identity fields, we pass on adding it to the database.
	// 	logging.Log(err3)
	// }
	tx, err4 := DbInstance.Beginx()
	if err4 != nil {
		logging.LogCrash(err4)
	}
	_, err5 := tx.NamedExec(addressUpdateInsert, dbA)
	if err5 != nil {
		logging.LogCrash(err5)
	}
	err6 := tx.Commit()
	if err6 != nil {
		return err6
	}
	return nil
}

// TODO: Mind that any errors happening within the transaction, if they need to bail from the transaction, they need to close it! otherwise you get database is locked.
// TODO: Should this take a pointer instead? It's dealing with some big amounts of data.
// BatchInsert insert a set of objects in a batch as a transaction.
func BatchInsert(apiObjects []interface{}) error {
	logging.Log(2, "Batch insert starting.")
	defer logging.Log(2, "Batch insert is complete.")
	numberOfObjectsCommitted := len(apiObjects)
	logging.Log(2, fmt.Sprintf("%v objects are being committed.", numberOfObjectsCommitted))

	start := time.Now()
	// fmt.Printf("%#v\n", apiObjects)
	// Begin transaction.
	tx, err := DbInstance.Beginx()
	if err != nil {
		logging.LogCrash(err)
	}
	// For each API object, convert to DB object and add to transaction.
	for _, apiObject := range apiObjects {
		// apiObject: API type, dbObj: DB type.
		dbo, err := APItoDB(apiObject)
		if err != nil {
			return errors.New(fmt.Sprint(
				"Error raised from APItoDB function used in Batch insert. Error: ", err))
		}
		err2 := enforceNoEmptyIdentityFields(dbo)
		if err2 != nil {
			// If this unit does have empty identity fields, we pass on adding it to the database.
			logging.Log(1, err2)
			continue
		}
		err3 := enforceNoEmptyRequiredFields(dbo)
		if err3 != nil {
			// If this unit does have empty identity fields, we pass on adding it to the database.
			logging.Log(1, err3)
			continue
		}
		switch dbObject := dbo.(type) {
		// case BoardPack:
		// 	if packShouldBeCommitted(dbObject) {
		// 		_, err := tx.NamedExec(boardInsert, dbObject.Board)
		// 		if err != nil {
		// 			logging.LogCrash(err)
		// 		}
		// 	}

		case BoardPack:
			if packShouldBeCommitted(dbObject) {
				_, err := tx.NamedExec(boardInsert, dbObject.Board)
				if err != nil {
					logging.LogCrash(err)
				}
				// Get the list of board owners before the transaction.
				boardBoardOwnersBeforeTx, err := getBoardOwnersBeforeTx(dbObject.Board.Fingerprint)
				if err != nil {
					logging.LogCrash(err)
				}
				// Get the changelist.
				changelist := generateBoardOwnerChangelist(
					boardBoardOwnersBeforeTx, dbObject.BoardOwners)
				for boardOwner, keepBoardOwner := range changelist {
					if keepBoardOwner == true {
						// We keep the owner's existence. This can either be a creation or an update, SQL deals with that.
						_, err := tx.NamedExec(boardOwnerInsert, boardOwner)
						if err != nil {
							// fmt.Printf("%#v\n", err)
							logging.LogCrash(err)
						}
					} else {
						// The owner is deleted. So we remove it from the database.
						_, err := tx.NamedExec(boardOwnerDelete, boardOwner)
						if err != nil {
							// fmt.Printf("%#v\n", err)
							logging.LogCrash(err)
						}
					}
				}
			}
		case DbThread:
			_, err := tx.NamedExec(threadInsert, dbObject)
			if err != nil {
				logging.LogCrash(err)
			}
		case DbPost:
			_, err := tx.NamedExec(postInsert, dbObject)
			if err != nil {
				logging.LogCrash(err)
			}
		case DbVote:
			_, err := tx.NamedExec(voteInsert, dbObject)
			if err != nil {
				logging.LogCrash(err)
			}
		case DbAddress:
			// In case of address, we strip out everything except the primary keys. This is because we cannot trust the data that is coming from the network. We just add the primary key set, and the local node will take care of directly connecting to these nodes and getting the details.
			// The other types of address inputs are not affected by this because they use InsertOrUpdateAddress, not this batch insert. If you're batch inserting addresses, it's by definition third party data.
			dbObject.LocationType = 0 // IPv4 or 6
			dbObject.Type = 0         // 2 = live, 255 = static
			dbObject.LastOnline = 0   // We cannot trust someone else's last online timestamp
			dbObject.ProtocolVersionMajor = 0
			dbObject.ProtocolVersionMinor = 0
			dbObject.ProtocolExtensions = ""
			dbObject.ClientVersionMajor = 0
			dbObject.ClientVersionMinor = 0
			dbObject.ClientVersionPatch = 0
			dbObject.ClientName = ""
			_, err := tx.NamedExec(addressInsert, dbObject)
			if err != nil {
				logging.LogCrash(err)
			}
		case KeyPack:
			if packShouldBeCommitted(dbObject) {
				_, err := tx.NamedExec(keyInsert, dbObject.Key)
				if err != nil {
					logging.LogCrash(err)
				}
				// Get the list of currency addresses before the transaction.
				currencyAddressesBeforeTx, err := getCurrencyAddressesBeforeTx(dbObject.Key.Fingerprint)
				// Get the changelist.
				changelist := generateCurrencyAddressChangelist(
					currencyAddressesBeforeTx, dbObject.CurrencyAddresses)
				for currencyAddress, keepcurrencyAddress := range changelist {
					if keepcurrencyAddress == true {
						// We keep the owner's existence. This can either be a creation or an update, SQL deals with that.
						_, err := tx.NamedExec(currencyAddressInsert, currencyAddress)
						if err != nil {
							// fmt.Printf("%#v\n", err)
							logging.LogCrash(err)
						}
					} else {
						// The owner is deleted. So we remove it from the database.
						_, err := tx.NamedExec(currencyAddressDelete, currencyAddress)
						if err != nil {
							// fmt.Printf("%#v\n", err)
							logging.LogCrash(err)
						}
					}
				}
			}
		case DbTruststate:
			_, err := tx.NamedExec(truststateInsert, dbObject)
			if err != nil {
				logging.LogCrash(err)
			}
		default:
			return errors.New(
				fmt.Sprintf(
					"This object type is something batch insert does not understand. Your object: %#v\n", dbObject))
		}
		// TODO: Create a prepared statement for each of those that allows for insertion.
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	elapsed := time.Since(start)
	logging.Log(2, fmt.Sprintf("It took %v to insert %v objects.", elapsed, numberOfObjectsCommitted))
	return nil
}

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

// getBoardOwnersBeforeTx is an internal function of the Writer. This gets the pre-transaction state of the board owners of the board that is being inserted.
func getBoardOwnersBeforeTx(boardFingerprint api.Fingerprint) ([]DbBoardOwner, error) {
	var boardBoardOwnersBeforeTx []DbBoardOwner
	// Fetch all Board BoardOwners of this board that is already in database.
	rowsOfBoardOwnersBeforeTx, err := DbInstance.Queryx("SELECT * from BoardOwners WHERE BoardFingerprint = ?", boardFingerprint)
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
	rowsOfCurrencyAddressesBeforeTx, err := DbInstance.Queryx("SELECT * from CurrencyAddresses WHERE KeyFingerprint = ?", keyFingerprint)
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

	case DbAddress:
		if obj.Location == "" || obj.Port == 0 {
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
		if obj.Board == "" || obj.Thread == "" || obj.Parent == "" || obj.Body == "" || obj.Creation == 0 || obj.ProofOfWork == "" {
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

	case DbAddress:
		if obj.LocationType == 0 || obj.LastOnline == 0 || obj.ProtocolVersionMajor == 0 || obj.ProtocolExtensions == "" || obj.ClientVersionMajor == 0 || obj.ClientName == "" {
			return errors.New(
				fmt.Sprintf(
					"This address has some required fields empty (One or more of: LocationType, LastOnline, ProtocolVersionMajor, ProtocolExtensions, ClientVersionMajor, ClientName). Address: %#v\n", obj))
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
