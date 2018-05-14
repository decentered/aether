// Backend > Event Horizon
// This library handles the impermanence of all things.

/*
  The birds have vanished down the sky.
  Now the last cloud drains away.

  We sit together, the mountain and me,
  until only the mountain remains.

  — Li Bai 李白
*/

package eventhorizon

import (
	"aether-core/services/globals"
	"aether-core/services/logging"
	"fmt"
	"os"
	"time"
)

type Timestamp int64

const (
	Day = time.Duration(24) * time.Hour
)

func delete(ts Timestamp, entityType string) {
	tableName := ""
	switch entityType {
	case "boards":
		tableName = "Boards"
	case "threads":
		tableName = "Threads"
	case "posts":
		tableName = "Posts"
	case "votes":
		tableName = "Votes"
	case "keys":
		tableName = "Keys"
	case "truststates":
		tableName = "Truststates"
	case "addresses":
		tableName = "Addresses"
	default:
		return
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE LastReferenced < ?", tableName)
	tx, err := globals.DbInstance.Beginx()
	tx.Exec(query, ts)
	tx.Commit()
}

func cnvToCutoff(days int) Timestamp {
	return Timestamp(time.Now().Add(-(time.Duration(days) * time.Hour * time.Duration(24))).Unix())
}

func max(ts1 Timestamp, ts2 Timestamp) Timestamp {
	if ts1 > ts2 {
		return ts1
	}
	return ts2
}

func deleteUpToLocalMemory() {
	lmD := globals.BackendConfig.GetLocalMemoryDays()
	lmCutoff := cnvToCutoff(lmD)
	vmD := globals.BackendConfig.GetVotesMemoryDays()
	vmCutoff := cnvToCutoff(vmD)
	delete(lmCutoff, "boards")
	delete(lmCutoff, "threads")
	delete(lmCutoff, "posts")
	delete(lmCutoff, "keys")
	delete(lmCutoff, "truststates")
	delete(lmCutoff, "addresses")
	// These are the special ones
	delete(vmCutoff, "votes")
}

func deleteUpToEH(eventhorizon Timestamp) {
	delete(eventhorizon, "boards")
	delete(eventhorizon, "threads")
	delete(eventhorizon, "posts")
	delete(eventhorizon, "keys")
	delete(eventhorizon, "truststates")
	// Addresses is limited to 1000 items and it has its own cycling logic. No need to delete based on event horizon, it will likely yield not many items. The LM cutoff deletion (deleteUpToLocalMemory) does that for us.
	// delete(eventhorizon, "addresses")
}

// improve this based on page couting for sqlite and whatever is needed for mysql. TODO
func getDbSize() int {
	switch globals.BackendConfig.GetDbEngine() {
	case "mysql":
		query := `
      SELECT
      SUM(
        ROUND(
          ((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024 ), 2)
        )
      AS "SIZE IN MB"
      FROM INFORMATION_SCHEMA.TABLES
      WHERE TABLE_SCHEMA = "AetherDB"
    `
		var size int
		err := db.Get(&size, query)
		if err != nil {
			logging.LogCrash("The attempt to read the MySQL database size failed.")
		}
		return size
	case "sqlite":
		dbLoc := fmt.Sprintf("%s/AetherDB.db", globals.BackendConfig.GetUserDirectory())
		fi, _ := os.Stat(dbLoc)
		// get the size
		size := fi.Size() / 1048576 // Assuming 1Mb = 1048576 bytes (binary, not decimal)
		return int(size)
	default:
		logging.LogCrash(fmt.Sprintf("This database type is not supported: %s", globals.BackendConfig.GetDbEngine()))
	}
	return -1 // Should never happen
}

func PruneDB() {
	lmD := globals.BackendConfig.GetLocalMemoryDays()
	lmCutoff := cnvToCutoff(lmD)
	tempeh := globals.BackendConfig.GetEventHorizonTimestamp()
	if tempeh <= lmCutoff {
		deleteUpToLocalMemory()
	}
	if getDbSize() <= globals.BackendConfig.GetMaxDbSizeMb() {
		tempeh = lmCutoff
	}
	itercount := 0
	for getDbSize() > globals.BackendConfig.GetMaxDbSizeMb() {
		/*
		   Logic:
		   - When we're scanning through things to delete, we back off event horizon one day, and start scanning from one day behind. That way, if the network pressure relieves, event horizon will gradually extend itself back. But if it stays high, it'll try to delete from one day behind, then iterate day by day, so we only lose one day cycle delete, which is not a big loss.
		   - Gotcha: When there is no network pressure (when the allotted size of the local database can contain the entire local memory days), event horizon is the same as local memory. In that case, starting from one day behind would be Local memory +1 days, in default, that would make it 181 days. Not that big of a deal, but we don't want that, so that's why the max() exists. It basically starts it from day 179 (assuming LM is = 180) regardless of what EH says.
		*/
		tempeh = tempeh + max(
			Timestamp((time.Duration(itercount)*Day)-Day),
			Timestamp(lmCutoff)+Timestamp(Day),
		)
		itercount++
	}
	globals.BackendConfig.SetEventHorizonTimestamp(tempeh)
}
