// Persistence > Base
// This file contains the SQL statements that are necessary to handle database action, as well as basic maintenance functions for the database itself.

package persistence

import (
	"aether-core/services/globals"
	"aether-core/services/logging"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	// "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	// _ "github.com/lib/pq"
	"errors"
	"github.com/fatih/color"
	"math/rand"
	"os"
	"strings"
	"time"
)

// DeleteDatabase removes the existing database in the default location.
func DeleteDatabase() {
	if globals.BackendConfig.GetDbEngine() == "sqlite" {
		os.RemoveAll(globals.BackendConfig.GetUserDirectory())
	} else if globals.BackendConfig.GetDbEngine() == "mysql" {
		globals.DbInstance.MustExec("DROP DATABASE `AetherDB`;")
	}
}

// CreateDatabase creates a new database in the default location and places into it the database schema.

func CreateDatabase() {
	err := createDatabase()
	if err != nil {
		if strings.Contains(err.Error(), "Database was locked") {
			logging.Log(1, err)
			if strings.Contains(err.Error(), "Database was locked") {
				logging.Log(1, "This transaction was not committed because database was locked. We'll wait 10 seconds and retry the transaction.")
				time.Sleep(10 * time.Second)
				logging.Log(1, "Retrying the previously failed CreateDatabase transaction.")
				err2 := createDatabase()
				if err2 != nil {
					logging.LogCrash(fmt.Sprintf("The second attempt to commit this data to the database failed. The first attempt had failed because the database was locked. The second attempt failed with the error: %s This database is corrupted. Quitting.", err2))
				} else { // If the reattempted transaction succeeds
					logging.Log(1, "The retry attempt of the failed transaction succeeded.")
				}
			}
		}
	}
}
func createDatabase() error {
	var schemaPrep1 string
	var schemaPrep2 string
	var schemaPrep3 string
	var schema1 string
	var schema2 string
	var schema3 string
	var schema4 string
	var schema5 string
	var schema6 string
	var schema7 string
	var schema8 string
	var schema9 string
	var schema10 string
	var schema11 string
	var schema12 string
	var schema13 string
	var schema14 string
	var schema15 string
	var schema16 string

	if globals.BackendConfig.DbEngine == "mysql" {
		schemaPrep1 = `
      CREATE DATABASE IF NOT EXISTS AetherDB
      DEFAULT CHARACTER SET utf8
      DEFAULT COLLATE utf8_general_ci;
    `
		schemaPrep2 = `USE AetherDB;`
		schemaPrep3 = `SET sql_mode = 'STRICT_TRANS_TABLES';`
		schema1 = `
        CREATE TABLE IF NOT EXISTS BoardOwners (
          BoardFingerprint VARCHAR(64) NOT NULL,
          KeyFingerprint VARCHAR(64) NOT NULL,
          Expiry BIGINT NOT NULL,
          Level SMALLINT NOT NULL,
          PRIMARY KEY(BoardFingerprint, KeyFingerprint)
        );
        `
		schema2 = ``
		schema3 = `
        CREATE TABLE IF NOT EXISTS Boards (
          Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
          Name VARCHAR(255) NOT NULL,
          Owner VARCHAR(64) NOT NULL,
          -- BoardOwners field will have to be constructed on the fly.
          Description LONGBLOB NOT NULL,  -- Converted from varchar(65535) to text, because it doesn't fit into a MYSQL table. Enforce max 65535 chars on the application layer.
          Creation BIGINT NOT NULL,
          ProofOfWork VARCHAR(1024) NOT NULL,
          Signature VARCHAR(512) NOT NULL,
          LastUpdate BIGINT NOT NULL,
          UpdateProofOfWork VARCHAR(1024) NOT NULL,
          UpdateSignature VARCHAR(512) NOT NULL,
          LocalArrival BIGINT NOT NULL,
          Meta JSON NOT NULL,
          RealmId VARCHAR(64) NOT NULL,
          EncrConcent LONGBLOB NOT NULL
        );`
		schema4 = `
        CREATE TABLE IF NOT EXISTS Threads (
          Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
          Board VARCHAR(64) NOT NULL,
          Name VARCHAR(255) NOT NULL,
          Body LONGBLOB NOT NULL,
          Link VARCHAR(5000) NOT NULL,
          Owner VARCHAR(64) NOT NULL,
          Creation BIGINT NOT NULL,
          ProofOfWork VARCHAR(1024) NOT NULL,
          Signature VARCHAR(512) NOT NULL,
          LastUpdate BIGINT NOT NULL,
          UpdateProofOfWork VARCHAR(1024) NOT NULL,
          UpdateSignature VARCHAR(512) NOT NULL,
          LocalArrival BIGINT NOT NULL,
          Meta JSON NOT NULL,
          RealmId VARCHAR(64) NOT NULL,
          EncrConcent LONGBLOB NOT NULL,
          INDEX (Board)
        );`
		schema5 = `
        CREATE TABLE IF NOT EXISTS Posts (
          Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
          Board VARCHAR(64) NOT NULL,
          Thread VARCHAR(64) NOT NULL,
          Parent VARCHAR(64) NOT NULL,
          Body LONGBLOB NOT NULL,
          Owner VARCHAR(64) NOT NULL,
          Creation BIGINT NOT NULL,
          ProofOfWork VARCHAR(1024) NOT NULL,
          Signature VARCHAR(512) NOT NULL,
          LastUpdate BIGINT NOT NULL,
          UpdateProofOfWork VARCHAR(1024) NOT NULL,
          UpdateSignature VARCHAR(512) NOT NULL,
          LocalArrival BIGINT NOT NULL,
          Meta JSON NOT NULL,
          RealmId VARCHAR(64) NOT NULL,
          EncrConcent LONGBLOB NOT NULL,
          INDEX (Thread)
        );`
		schema6 = `
        CREATE TABLE IF NOT EXISTS Votes (
          Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
          Board VARCHAR(64) NOT NULL,
          Thread VARCHAR(64) NOT NULL,
          Target VARCHAR(64) NOT NULL,
          Owner VARCHAR(64) NOT NULL,
          Type SMALLINT NOT NULL,
          Creation BIGINT NOT NULL,
          ProofOfWork VARCHAR(1024) NOT NULL,
          Signature VARCHAR(512) NOT NULL,
          LastUpdate BIGINT NOT NULL,
          UpdateProofOfWork VARCHAR(1024) NOT NULL,
          UpdateSignature VARCHAR(512) NOT NULL,
          LocalArrival BIGINT NOT NULL,
          Meta JSON NOT NULL,
          RealmId VARCHAR(64) NOT NULL,
          EncrConcent LONGBLOB NOT NULL,
          INDEX (Target)
        );`
		schema7 = `
        CREATE TABLE IF NOT EXISTS Addresses (
          Location VARCHAR(256) NOT NULL, -- From 2500
          Sublocation VARCHAR(256) NOT NULL, -- From 2500
          Port INTEGER NOT NULL,
          IPType SMALLINT NOT NULL,
          AddressType SMALLINT NOT NULL,
          LastOnline BIGINT NOT NULL,
          ProtocolVersionMajor SMALLINT NOT NULL,
          ProtocolVersionMinor INTEGER NOT NULL,
          ClientVersionMajor SMALLINT NOT NULL,
          ClientVersionMinor INTEGER NOT NULL,
          ClientVersionPatch INTEGER NOT NULL,
          ClientName VARCHAR(255) NOT NULL,
          LocalArrival BIGINT NOT NULL,
          RealmId VARCHAR(64) NOT NULL,
          PRIMARY KEY(Location, Sublocation, Port)
        );`
		schema8 = `
        CREATE TABLE IF NOT EXISTS PublicKeys (
          Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
          Type VARCHAR(64) NOT NULL,
          PublicKey TEXT NOT NULL,
          Name VARCHAR(64) NOT NULL,
          Info LONGBLOB NOT NULL,
          Creation BIGINT NOT NULL,
          ProofOfWork VARCHAR(1024) NOT NULL,
          Signature VARCHAR(512) NOT NULL,
          LastUpdate BIGINT NOT NULL,
          UpdateProofOfWork VARCHAR(1024) NOT NULL,
          UpdateSignature VARCHAR(512) NOT NULL,
          LocalArrival BIGINT NOT NULL,
          Meta JSON NOT NULL,
          RealmId VARCHAR(64) NOT NULL,
          EncrConcent LONGBLOB NOT NULL
        );`
		schema9 = `
        CREATE TABLE IF NOT EXISTS Truststates (
          Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
          Target VARCHAR(64) NOT NULL,
          Owner VARCHAR(64) NOT NULL,
          Type SMALLINT NOT NULL,
          Domains VARCHAR(7000) NOT NULL,
          Expiry BIGINT NOT NULL,
          Creation BIGINT NOT NULL,
          ProofOfWork VARCHAR(1024) NOT NULL,
          Signature VARCHAR(512) NOT NULL,
          LastUpdate BIGINT NOT NULL,
          UpdateProofOfWork VARCHAR(1024) NOT NULL,
          UpdateSignature VARCHAR(512) NOT NULL,
          LocalArrival BIGINT NOT NULL,
          Meta JSON NOT NULL,
          RealmId VARCHAR(64) NOT NULL,
          EncrConcent LONGBLOB NOT NULL
        );
      `
		schema10 = `
          CREATE TABLE IF NOT EXISTS Nodes (
            Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
            BoardsLastCheckin BIGINT NOT NULL,
            ThreadsLastCheckin BIGINT NOT NULL,
            PostsLastCheckin BIGINT NOT NULL,
            VotesLastCheckin BIGINT NOT NULL,
            AddressesLastCheckin BIGINT NOT NULL,
            KeysLastCheckin BIGINT NOT NULL,
            TruststatesLastCheckin BIGINT NOT NULL
          );
        `
		schema11 = `
          CREATE TABLE IF NOT EXISTS Subprotocols (
            Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
            Name VARCHAR(64) NOT NULL,
            VersionMajor SMALLINT NOT NULL,
            VersionMinor INTEGER NOT NULL,
            SupportedEntities VARCHAR(5000) NOT NULL
          );`

		schema12 = `
          CREATE TABLE IF NOT EXISTS AddressesSubprotocols (
            AddressLocation VARCHAR(256) NOT NULL,
            AddressSublocation VARCHAR(256) NOT NULL,
            AddressPort INTEGER NOT NULL,
            SubprotocolFingerprint VARCHAR(64) NOT NULL,
            PRIMARY KEY(AddressLocation, AddressSublocation, AddressPort, SubprotocolFingerprint)
          );`
		schema16 = `
          CREATE TABLE IF NOT EXISTS Diagnostics (
            DbRoundtripTestField BIGINT PRIMARY KEY NOT NULL
          );`
	} else if globals.BackendConfig.DbEngine == "sqlite" {
		schemaPrep1 = ``
		schemaPrep2 = ``
		schemaPrep3 = ``
		schema1 = `
        CREATE TABLE IF NOT EXISTS "BoardOwners" (
          "BoardFingerprint" varchar(64) NOT NULL
        ,  "KeyFingerprint" varchar(64) NOT NULL
        ,  "Expiry" integer NOT NULL
        ,  "Level" integer NOT NULL
        ,  PRIMARY KEY ("BoardFingerprint","KeyFingerprint")
        );`
		schema2 = ``
		schema3 = `
        CREATE TABLE IF NOT EXISTS "Boards" (
          "Fingerprint" varchar(64) NOT NULL
        ,  "Name" varchar(255) NOT NULL
        ,  "Owner" varchar(64) NOT NULL
        ,  "Description" blob NOT NULL
        ,  "Creation" integer NOT NULL
        ,  "ProofOfWork" varchar(1024) NOT NULL
        ,  "Signature" varchar(512) NOT NULL
        ,  "LastUpdate" integer NOT NULL
        ,  "UpdateProofOfWork" varchar(1024) NOT NULL
        ,  "UpdateSignature" varchar(512) NOT NULL
        ,  "LocalArrival" integer NOT NULL
        ,  "Meta" text NOT NULL
        ,  "RealmId" varchar(64) NOT NULL
        ,  "EncrContent" blob NOT NULL
        ,  PRIMARY KEY ("Fingerprint")
        );`
		schema4 = `
        CREATE TABLE IF NOT EXISTS "Threads" (
          "Fingerprint" varchar(64) NOT NULL
        ,  "Board" varchar(64) NOT NULL
        ,  "Name" varchar(255) NOT NULL
        ,  "Body" blob NOT NULL
        ,  "Link" varchar(5000) NOT NULL
        ,  "Owner" varchar(64) NOT NULL
        ,  "Creation" integer NOT NULL
        ,  "ProofOfWork" varchar(1024) NOT NULL
        ,  "Signature" varchar(512) NOT NULL
        ,  "LastUpdate" integer NOT NULL
        ,  "UpdateProofOfWork" varchar(1024) NOT NULL
        ,  "UpdateSignature" varchar(512) NOT NULL
        ,  "LocalArrival" integer NOT NULL
        ,  "Meta" text NOT NULL
        ,  "RealmId" varchar(64) NOT NULL
        ,  "EncrContent" blob NOT NULL
        ,  PRIMARY KEY ("Fingerprint")
        );`
		schema5 = `
        CREATE TABLE IF NOT EXISTS "Posts" (
          "Fingerprint" varchar(64) NOT NULL
        ,  "Board" varchar(64) NOT NULL
        ,  "Thread" varchar(64) NOT NULL
        ,  "Parent" varchar(64) NOT NULL
        ,  "Body" blob NOT NULL
        ,  "Owner" varchar(64) NOT NULL
        ,  "Creation" integer NOT NULL
        ,  "ProofOfWork" varchar(1024) NOT NULL
        ,  "Signature" varchar(512) NOT NULL
        ,  "LastUpdate" integer NOT NULL
        ,  "UpdateProofOfWork" varchar(1024) NOT NULL
        ,  "UpdateSignature" varchar(512) NOT NULL
        ,  "LocalArrival" integer NOT NULL
        ,  "Meta" text NOT NULL
        ,  "RealmId" varchar(64) NOT NULL
        ,  "EncrContent" blob NOT NULL
        ,  PRIMARY KEY ("Fingerprint")
        );`
		schema6 = `
        CREATE TABLE IF NOT EXISTS "Votes" (
          "Fingerprint" varchar(64) NOT NULL
        ,  "Board" varchar(64) NOT NULL
        ,  "Thread" varchar(64) NOT NULL
        ,  "Target" varchar(64) NOT NULL
        ,  "Owner" varchar(64) NOT NULL
        ,  "Type" integer NOT NULL
        ,  "Creation" integer NOT NULL
        ,  "ProofOfWork" varchar(1024) NOT NULL
        ,  "Signature" varchar(512) NOT NULL
        ,  "LastUpdate" integer NOT NULL
        ,  "UpdateProofOfWork" varchar(1024) NOT NULL
        ,  "UpdateSignature" varchar(512) NOT NULL
        ,  "LocalArrival" integer NOT NULL
        ,  "Meta" text NOT NULL
        ,  "RealmId" varchar(64) NOT NULL
        ,  "EncrContent" blob NOT NULL
        ,  PRIMARY KEY ("Fingerprint")
        );`
		schema7 = `
        CREATE TABLE IF NOT EXISTS "Addresses" (
          "Location" varchar(256) NOT NULL
        ,  "Sublocation" varchar(256) NOT NULL
        ,  "Port" integer NOT NULL
        ,  "IPType" integer NOT NULL
        ,  "AddressType" integer NOT NULL
        ,  "LastOnline" integer NOT NULL
        ,  "ProtocolVersionMajor" integer NOT NULL
        ,  "ProtocolVersionMinor" integer NOT NULL
        ,  "ClientVersionMajor" integer NOT NULL
        ,  "ClientVersionMinor" integer NOT NULL
        ,  "ClientVersionPatch" integer NOT NULL
        ,  "ClientName" varchar(255) NOT NULL
        ,  "LocalArrival" integer NOT NULL
        ,  "RealmId" varchar(64) NOT NULL
        ,  PRIMARY KEY ("Location","Sublocation","Port")
        );`
		schema8 = `
        CREATE TABLE IF NOT EXISTS "PublicKeys" (
          "Fingerprint" varchar(64) NOT NULL
        ,  "Type" varchar(64) NOT NULL
        ,  "PublicKey" text NOT NULL
        ,  "Name" varchar(64) NOT NULL
        ,  "Info" blob NOT NULL
        ,  "Creation" integer NOT NULL
        ,  "ProofOfWork" varchar(1024) NOT NULL
        ,  "Signature" varchar(512) NOT NULL
        ,  "LastUpdate" integer NOT NULL
        ,  "UpdateProofOfWork" varchar(1024) NOT NULL
        ,  "UpdateSignature" varchar(512) NOT NULL
        ,  "LocalArrival" integer NOT NULL
        ,  "Meta" text NOT NULL
        ,  "RealmId" varchar(64) NOT NULL
        ,  "EncrContent" blob NOT NULL
        ,  PRIMARY KEY ("Fingerprint")
        );`
		schema9 = `
        CREATE TABLE IF NOT EXISTS "Truststates" (
          "Fingerprint" varchar(64) NOT NULL
        ,  "Target" varchar(64) NOT NULL
        ,  "Owner" varchar(64) NOT NULL
        ,  "Type" integer NOT NULL
        ,  "Domains" varchar(7000) NOT NULL
        ,  "Expiry" integer NOT NULL
        ,  "Creation" integer NOT NULL
        ,  "ProofOfWork" varchar(1024) NOT NULL
        ,  "Signature" varchar(512) NOT NULL
        ,  "LastUpdate" integer NOT NULL
        ,  "UpdateProofOfWork" varchar(1024) NOT NULL
        ,  "UpdateSignature" varchar(512) NOT NULL
        ,  "LocalArrival" integer NOT NULL
        ,  "Meta" text NOT NULL
        ,  "RealmId" varchar(64) NOT NULL
        ,  "EncrContent" blob NOT NULL
        ,  PRIMARY KEY ("Fingerprint")
        );`
		schema10 = `
          CREATE TABLE IF NOT EXISTS "Nodes" (
            "Fingerprint" varchar(64) NOT NULL
          ,  "BoardsLastCheckin" integer NOT NULL
          ,  "ThreadsLastCheckin" integer NOT NULL
          ,  "PostsLastCheckin" integer NOT NULL
          ,  "VotesLastCheckin" integer NOT NULL
          ,  "AddressesLastCheckin" integer NOT NULL
          ,  "KeysLastCheckin" integer NOT NULL
          ,  "TruststatesLastCheckin" integer NOT NULL
          ,  PRIMARY KEY ("Fingerprint")
          );`
		schema11 = `
          CREATE TABLE IF NOT EXISTS "Subprotocols" (
            "Fingerprint" varchar(64) NOT NULL
          ,  "Name" varchar(64) NOT NULL
          ,  "VersionMajor" integer NOT NULL
          ,  "VersionMinor" integer NOT NULL
          ,  "SupportedEntities" varchar(5000) NOT NULL
          ,  PRIMARY KEY ("Fingerprint")
          );`
		schema12 = `
          CREATE TABLE IF NOT EXISTS "AddressesSubprotocols" (
            "AddressLocation" varchar(256) NOT NULL
          ,  "AddressSublocation" varchar(256) NOT NULL
          ,  "AddressPort" integer NOT NULL
          ,  "SubprotocolFingerprint" varchar(64) NOT NULL
          ,  PRIMARY KEY ("AddressLocation","AddressSublocation","AddressPort","SubprotocolFingerprint")
          );`
		schema13 = `
          CREATE INDEX IF NOT EXISTS "idx_Posts_Thread" ON "Posts" ("Thread");
          `
		schema14 = `
          CREATE INDEX IF NOT EXISTS "idx_Threads_Board" ON "Threads" ("Board");
          `
		schema15 = `
          CREATE INDEX IF NOT EXISTS "idx_Votes_Target" ON "Votes" ("Target");
          `
		schema16 = `
            CREATE TABLE IF NOT EXISTS "Diagnostics" (
              "DbRoundtripTestField" integer NOT NULL
            ,  PRIMARY KEY ("DbRoundtripTestField")
            );`
	} else {
		logging.LogCrash(fmt.Sprintf("Storage engine you've inputted is not supported. Please change it from the backend user config into something that is supported. You've provided: %s", globals.BackendConfig.GetDbEngine()))
	}

	var creationSchemas []string
	if len(schemaPrep1) > 0 {
		// MySQL specific DB creation schema.
		creationSchemas = append(creationSchemas, schemaPrep1)
		creationSchemas = append(creationSchemas, schemaPrep2)
		creationSchemas = append(creationSchemas, schemaPrep3)
	}
	creationSchemas = append(creationSchemas, schema1)
	creationSchemas = append(creationSchemas, schema2)
	creationSchemas = append(creationSchemas, schema3)
	creationSchemas = append(creationSchemas, schema4)
	creationSchemas = append(creationSchemas, schema5)
	creationSchemas = append(creationSchemas, schema6)
	creationSchemas = append(creationSchemas, schema7)
	creationSchemas = append(creationSchemas, schema8)
	creationSchemas = append(creationSchemas, schema9)
	creationSchemas = append(creationSchemas, schema10)
	creationSchemas = append(creationSchemas, schema11)
	creationSchemas = append(creationSchemas, schema12)
	// Indexes in SQLite require separate statements
	if len(schema13) > 0 && len(schema14) > 0 && len(schema15) > 0 {
		creationSchemas = append(creationSchemas, schema13)
		creationSchemas = append(creationSchemas, schema14)
		creationSchemas = append(creationSchemas, schema15)
	}
	creationSchemas = append(creationSchemas, schema16)

	tx, err := globals.DbInstance.Beginx()
	if err != nil {
		logging.LogCrash(err)
	}
	for _, schema := range creationSchemas {
		// fmt.Println(schema)
		_, err2 := tx.Exec(schema)
		if err2 != nil {
			logging.LogCrash(err2)
		}
	}
	err3 := tx.Commit()
	if err3 != nil {
		tx.Rollback()
		logging.Log(1, fmt.Sprintf("CreateDatabase encountered an error when trying to commit to the database. Error is: %s", err3))
		if strings.Contains(err3.Error(), "database is locked") {
			logging.Log(1, fmt.Sprintf("This database seems to be locked. We'll sleep 10 seconds to give it the time it needs to recover. This mostly happens when the app has crashed and there is a hot journal - and SQLite is in the process of repairing the database. THE DATA IN THIS TRANSACTION WAS NOT COMMITTED. PLEASE RETRY."))
			return errors.New("Database was locked. THE DATA IN THIS TRANSACTION WAS NOT COMMITTED. PLEASE RETRY.")
		}
		return err3
	}
	return nil
}

func CheckDatabaseReady() {
	DiagInsert := `REPLACE INTO Diagnostics
  (
    DbRoundtripTestField
  ) VALUES (
    :DbRoundtripTestField
  )`
	DiagDelete := `DELETE FROM Diagnostics`
	rand.Seed(time.Now().UTC().UnixNano())
	// We're using time.now because we don't want DB to optimise out the write and not test the connection that way. Get a random number between 0- 65535 for entry test.
	ss := map[string]interface{}{"DbRoundtripTestField": rand.Intn(65535)}
	// First, remove everything.
	tx, err := globals.DbInstance.Beginx()
	if err != nil {
		logging.LogCrash(err)
	}
	_, err2 := tx.Exec(DiagDelete)
	if err2 != nil {
		logging.LogCrash(err2)
	}
	tx.Commit()
	// Second, insert a new item.
	tx2, err3 := globals.DbInstance.Beginx()
	if err3 != nil {
		logging.LogCrash(err3)
	}
	_, err4 := tx2.NamedExec(DiagInsert, ss)
	if err4 != nil {
		logging.LogCrash(err4)
	}
	err5 := tx2.Commit()
	if err5 != nil {
		logging.LogCrash(err4)
	}
	color.Cyan("Database is ready. Just verified by removing and inserting data successfully.")
}

// Insertion SQL code used by the writer.

// NodeInsert just inserts the Node details into the entry. This is mutable.
var nodeInsert = `REPLACE INTO Nodes
(
  Fingerprint, BoardsLastCheckin, ThreadsLastCheckin, PostsLastCheckin,
  VotesLastCheckin, AddressesLastCheckin, KeysLastCheckin,
  TruststatesLastCheckin
) VALUES (
  :Fingerprint, :BoardsLastCheckin, :ThreadsLastCheckin, :PostsLastCheckin,
  :VotesLastCheckin, :AddressesLastCheckin, :KeysLastCheckin,
  :TruststatesLastCheckin
)`

// Board insert does insert or replace without checking because we're handling the logic that decides whether we should update or not in the database layer.
var boardInsert = `REPLACE INTO Boards
  (
    Fingerprint, Name, Owner, Description, LocalArrival,
    Meta, RealmId, EncrContent,
    Creation, ProofOfWork, Signature,
    LastUpdate, UpdateProofOfWork, UpdateSignature
  ) VALUES (
    :Fingerprint, :Name, :Owner, :Description, :LocalArrival,
    :Meta, :RealmId, :EncrContent,
    :Creation, :ProofOfWork, :Signature,
    :LastUpdate, :UpdateProofOfWork, :UpdateSignature
  )`

// BoardOwners are mutable, but the condition of mutation is handled in the application layer. The only place the REPLACE could trigger is change of Expiry and level. The BoardFingerprint and KeyFingerprint are identity columns, so anything with different data on those will be committed as a new item.
var boardOwnerInsert = `REPLACE INTO BoardOwners
(
  BoardFingerprint, KeyFingerprint, Expiry, Level
) VALUES (
  :BoardFingerprint, :KeyFingerprint, :Expiry, :Level
)`

// Deletion for BoardOwner. This triggers when a person is no longer a moderator, etc. This is the only deletion here because nothing else really gets deleted.
var boardOwnerDelete = `DELETE FROM BoardOwners WHERE BoardFingerprint = :BoardFingerprint AND KeyFingerprint = :KeyFingerprint`

var threadInsert = `REPLACE INTO Threads
  SELECT Candidate.* FROM
  (SELECT :Fingerprint AS Fingerprint,
          :Board AS Board,
          :Name AS Name,
          :Body AS Body,
          :Link AS Link,
          :Owner AS Owner,
          :Creation AS Creation,
          :ProofOfWork AS ProofOfWork,
          :Signature AS Signature,
          :LastUpdate AS LastUpdate,
          :UpdateProofOfWork AS UpdateProofOfWork,
          :UpdateSignature AS UpdateSignature,
          :LocalArrival AS LocalArrival,
          :Meta AS Meta,
          :RealmId AS RealmId,
          :EncrContent AS EncrContent
          ) AS Candidate
    LEFT JOIN Threads ON Candidate.Fingerprint = Threads.Fingerprint
    WHERE (Candidate.LastUpdate > Threads.LastUpdate AND Candidate.LastUpdate > Threads.Creation)
    OR Threads.Fingerprint IS NULL`

var postInsert = `REPLACE INTO Posts
  SELECT Candidate.* FROM
  (SELECT :Fingerprint AS Fingerprint,
          :Board AS Board,
          :Thread AS Thread,
          :Parent AS Parent,
          :Body AS Body,
          :Owner AS Owner,
          :Creation AS Creation,
          :ProofOfWork AS ProofOfWork,
          :Signature AS Signature,
          :LastUpdate AS LastUpdate,
          :UpdateProofOfWork AS UpdateProofOfWork,
          :UpdateSignature AS UpdateSignature,
          :LocalArrival AS LocalArrival,
          :Meta AS Meta,
          :RealmId AS RealmId,
          :EncrContent AS EncrContent
          ) AS Candidate
    LEFT JOIN Posts ON Candidate.Fingerprint = Posts.Fingerprint
    WHERE (Candidate.LastUpdate > Posts.LastUpdate AND Candidate.LastUpdate > Posts.Creation)
    OR Posts.Fingerprint IS NULL`

var voteInsert = `REPLACE INTO Votes
  SELECT Candidate.* FROM
  (SELECT :Fingerprint AS Fingerprint,
          :Board AS Board,
          :Thread AS Thread,
          :Target AS Target,
          :Owner AS Owner,
          :Type AS Type,
          :Creation AS Creation,
          :ProofOfWork AS ProofOfWork,
          :Signature AS Signature,
          :LastUpdate AS LastUpdate,
          :UpdateProofOfWork AS UpdateProofOfWork,
          :UpdateSignature AS UpdateSignature,
          :LocalArrival AS LocalArrival,
          :Meta AS Meta,
          :RealmId AS RealmId,
          :EncrContent AS EncrContent
          ) AS Candidate
  LEFT JOIN Votes ON Candidate.Fingerprint = Votes.Fingerprint
  WHERE (Candidate.LastUpdate > Votes.LastUpdate AND Candidate.LastUpdate > Votes.Creation)
  OR Votes.Fingerprint IS NULL`

var keyInsert = `REPLACE INTO PublicKeys
    SELECT Candidate.* FROM
    (SELECT :Fingerprint AS Fingerprint,
            :Type AS Type,
            :PublicKey AS PublicKey,
            :Name AS Name,
            :Info AS Info,
            :Creation AS Creation,
            :ProofOfWork AS ProofOfWork,
            :Signature AS Signature,
            :LastUpdate AS LastUpdate,
            :UpdateProofOfWork AS UpdateProofOfWork,
            :UpdateSignature AS UpdateSignature,
            :LocalArrival AS LocalArrival,
            :Meta AS Meta,
            :RealmId AS RealmId,
            :EncrContent AS EncrContent
            ) AS Candidate
    LEFT JOIN PublicKeys ON Candidate.Fingerprint = PublicKeys.Fingerprint
    WHERE (Candidate.LastUpdate > PublicKeys.LastUpdate AND Candidate.LastUpdate > PublicKeys.Creation)
    OR PublicKeys.Fingerprint IS NULL`

var truststateInsert = `REPLACE INTO Truststates
  SELECT Candidate.* FROM
  (SELECT :Fingerprint AS Fingerprint,
          :Target AS Target,
          :Owner AS Owner,
          :Type AS Type,
          :Domains AS Domains,
          :Expiry AS Expiry,
          :Creation AS Creation,
          :ProofOfWork AS ProofOfWork,
          :Signature AS Signature,
          :LastUpdate AS LastUpdate,
          :UpdateProofOfWork AS UpdateProofOfWork,
          :UpdateSignature AS UpdateSignature,
          :LocalArrival AS LocalArrival,
          :Meta AS Meta,
          :RealmId AS RealmId,
          :EncrContent AS EncrContent
          ) AS Candidate
  LEFT JOIN Truststates ON Candidate.Fingerprint = Truststates.Fingerprint
  WHERE (Candidate.LastUpdate > Truststates.LastUpdate AND Candidate.LastUpdate > Truststates.Creation)
  OR Truststates.Fingerprint IS NULL`

// Address insert is immutable. This is used for when a node receives data from an address from a node that is not at the aforementioned address. In other words, an address object coming from a third party node not at that address cannot change an existing address saved in the database.
var addressInsertMySQL = `INSERT IGNORE INTO Addresses
  (
    Location, Sublocation, Port, IPType, AddressType, LastOnline,
    ProtocolVersionMajor, ProtocolVersionMinor, ClientVersionMajor,
    ClientVersionMinor, ClientVersionPatch, ClientName, RealmId,
    LocalArrival
  ) VALUES (
    :Location, :Sublocation, :Port,:IPType, :AddressType, :LastOnline,
    :ProtocolVersionMajor, :ProtocolVersionMinor, :ClientVersionMajor,
    :ClientVersionMinor, :ClientVersionPatch, :ClientName, :RealmId,
    :LocalArrival
  )`

var addressInsertSQLite = `INSERT OR IGNORE INTO Addresses
  (
    Location, Sublocation, Port, IPType, AddressType, LastOnline,
    ProtocolVersionMajor, ProtocolVersionMinor, ClientVersionMajor,
    ClientVersionMinor, ClientVersionPatch, ClientName, RealmId,
    LocalArrival
  ) VALUES (
    :Location, :Sublocation, :Port,:IPType, :AddressType, :LastOnline,
    :ProtocolVersionMajor, :ProtocolVersionMinor, :ClientVersionMajor,
    :ClientVersionMinor, :ClientVersionPatch, :ClientName, :RealmId,
    :LocalArrival
  )`

// Address update insert is mutable. This is used when the node connects to the address itself. Example: When a node connects to 256.253.231.123:8080, it will update the entry for that address with the data coming from the remote node. This is the only way to mutate an address object.
var addressUpdateInsert = `REPLACE INTO Addresses
  (
    Location, Sublocation, Port, IPType, AddressType, LastOnline,
    ProtocolVersionMajor, ProtocolVersionMinor, ClientVersionMajor,
    ClientVersionMinor, ClientVersionPatch, ClientName, RealmId,
    LocalArrival
  ) VALUES (
    :Location, :Sublocation, :Port,:IPType, :AddressType, :LastOnline,
    :ProtocolVersionMajor, :ProtocolVersionMinor, :ClientVersionMajor,
    :ClientVersionMinor, :ClientVersionPatch, :ClientName, :RealmId,
    :LocalArrival
  )`

// Subprotocol insert is the part of address insertion series. This makes it so that we have a list of all subprotocols flying around.
var subprotocolInsert = `REPLACE INTO Subprotocols
  (
    Fingerprint, Name, VersionMajor, VersionMinor, SupportedEntities
  ) VALUES (
    :Fingerprint, :Name, :VersionMajor, :VersionMinor, :SupportedEntities
  )`

// AddressSubprotocolInsert inserts into the many to many junction table so that we can keep track of the subprotocols an address uses.

var addressSubprotocolInsertMySQL = `INSERT IGNORE INTO AddressesSubprotocols
   (
     AddressLocation, AddressSublocation, AddressPort, SubprotocolFingerprint
   ) VALUES (
     :AddressLocation, :AddressSublocation, :AddressPort, :SubprotocolFingerprint
   )`

var addressSubprotocolInsertSQLite = `INSERT OR IGNORE INTO AddressesSubprotocols
    (
      AddressLocation, AddressSublocation, AddressPort, SubprotocolFingerprint
    ) VALUES (
      :AddressLocation, :AddressSublocation, :AddressPort, :SubprotocolFingerprint
    )`
