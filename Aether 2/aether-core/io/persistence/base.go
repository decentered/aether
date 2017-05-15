// Persistence > Base
// This file contains the SQL statements that are necessary to handle database action, as well as basic maintenance functions for the database itself.

package persistence

import (
	"github.com/jmoiron/sqlx"
	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
	// "fmt"
	// _ "github.com/lib/pq"
	// "os"
)

// Global Objects

// Creates the database connection to be used from this point on.
// var DbInstance = sqlx.MustConnect("sqlite3", "./test.db")
var DbInstance = sqlx.MustConnect("mysql", "root:@/aether_test")

// var DbInstance = sqlx.MustConnect("postgres", "user=burak password=12345 dbname=aether_test sslmode=disable")

// func SetMaxOpenConn() {
// 	DbInstance.SetMaxOpenConns(10000000)
// }

// DeleteDatabase removes the existing database in the default location.
func DeleteDatabase() {
	// os.Remove("./test.db")
	DbInstance.MustExec("DROP TABLE `aether_test`.`Addresses`, `aether_test`.`BoardOwners`, `aether_test`.`Boards`, `aether_test`.`CurrencyAddresses`, `aether_test`.`Posts`, `aether_test`.`PublicKeys`, `aether_test`.`Threads`, `aether_test`.`Truststates`, `aether_test`.`Votes`;")
}

// CreateDatabase creates a new database in the default location and places into it the database schema.
func CreateDatabase() {
	schema1 := `
    CREATE TABLE IF NOT EXISTS BoardOwners (
      BoardFingerprint VARCHAR(64) NOT NULL,
      KeyFingerprint VARCHAR(64) NOT NULL,
      Expiry BIGINT NOT NULL,
      Level SMALLINT NOT NULL,
      PRIMARY KEY(BoardFingerprint, KeyFingerprint)
    );
    `
	schema2 := `
    CREATE TABLE IF NOT EXISTS CurrencyAddresses (
      KeyFingerprint VARCHAR(64) NOT NULL,
      CurrencyCode VARCHAR(5) NOT NULL,
      Address VARCHAR(512) NOT NULL, -- Changed from 1024
      PRIMARY KEY(KeyFingerprint, Address)
    );`
	schema3 := `
    CREATE TABLE IF NOT EXISTS Boards (
      Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
      Name VARCHAR(255) NOT NULL,
      Owner VARCHAR(64) NOT NULL,
      -- BoardOwners field will have to be constructed on the fly.
      Description TEXT NOT NULL,  -- Converted from varchar(65535) to text, because it doesn't fit into a MYSQL table. Enforce max 65535 chars on the application layer.
      Creation BIGINT NOT NULL,
      ProofOfWork VARCHAR(1024) NOT NULL,
      Signature VARCHAR(512) NOT NULL,
      LastUpdate BIGINT NOT NULL,
      UpdateProofOfWork VARCHAR(1024) NOT NULL,
      UpdateSignature VARCHAR(512) NOT NULL,
      LocalArrival BIGINT NOT NULL
    );`
	schema4 := `
    CREATE TABLE IF NOT EXISTS Threads (
      Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
      Board VARCHAR(64) NOT NULL,
      Name VARCHAR(255) NOT NULL,
      Body TEXT NOT NULL,
      Link VARCHAR(5000) NOT NULL,
      Owner VARCHAR(64) NOT NULL,
      Creation BIGINT NOT NULL,
      ProofOfWork VARCHAR(1024) NOT NULL,
      Signature VARCHAR(512) NOT NULL,
      LocalArrival BIGINT NOT NULL,
      INDEX (Board)
    );`
	schema5 := `
    CREATE TABLE IF NOT EXISTS Posts (
      Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
      Board VARCHAR(64) NOT NULL,
      Thread VARCHAR(64) NOT NULL,
      Parent VARCHAR(64) NOT NULL,
      Body TEXT NOT NULL,
      Owner VARCHAR(64) NOT NULL,
      Creation BIGINT NOT NULL,
      ProofOfWork VARCHAR(1024) NOT NULL,
      Signature VARCHAR(512) NOT NULL,
      LocalArrival BIGINT NOT NULL,
      INDEX (Thread)
    );`
	schema6 := `
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
      INDEX (Target)
    );`
	schema7 := `
    CREATE TABLE IF NOT EXISTS Addresses (
      Location VARCHAR(256) NOT NULL, -- From 2500
      Sublocation VARCHAR(256) NOT NULL, -- From 2500
      Port INTEGER NOT NULL,
      IPType SMALLINT NOT NULL,
      AddressType SMALLINT NOT NULL,
      LastOnline BIGINT NOT NULL,
      ProtocolVersionMajor SMALLINT NOT NULL,
      ProtocolVersionMinor INTEGER NOT NULL,
      ProtocolExtensions VARCHAR(5000) NOT NULL,
      ClientVersionMajor SMALLINT NOT NULL,
      ClientVersionMinor INTEGER NOT NULL,
      ClientVersionPatch INTEGER NOT NULL,
      ClientName VARCHAR(255) NOT NULL,
      LocalArrival BIGINT NOT NULL,
      PRIMARY KEY(Location, Sublocation, Port)
    );`
	schema8 := `
    CREATE TABLE IF NOT EXISTS PublicKeys (
      Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
      Type VARCHAR(64) NOT NULL,
      PublicKey TEXT NOT NULL,
      Name VARCHAR(64) NOT NULL,
      -- CurrencyAddresses will have to be constructed on the fly.
      Info VARCHAR(1024) NOT NULL,
      Creation BIGINT NOT NULL,
      ProofOfWork VARCHAR(1024) NOT NULL,
      Signature VARCHAR(512) NOT NULL,
      LastUpdate BIGINT NOT NULL,
      UpdateProofOfWork VARCHAR(1024) NOT NULL,
      UpdateSignature VARCHAR(512) NOT NULL,
      LocalArrival BIGINT NOT NULL
    );`
	schema9 := `
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
      LocalArrival BIGINT NOT NULL
    );
  `
	schema10 := `
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
	var creationSchemas []string
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

	for _, schema := range creationSchemas {
		// fmt.Println(schema)
		DbInstance.MustExec(schema)
	}
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
    Creation, ProofOfWork, Signature,
    LastUpdate, UpdateProofOfWork, UpdateSignature
  ) VALUES (
    :Fingerprint, :Name, :Owner, :Description, :LocalArrival,
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

// Deletion for BoardOwner. This triggers when a person is no longer a moderator, etc.
var boardOwnerDelete = `DELETE FROM BoardOwners WHERE BoardFingerprint = :BoardFingerprint AND KeyFingerprint = :KeyFingerprint`

// Immutable
var threadInsert = `INSERT IGNORE INTO Threads
(
  Fingerprint, Board, Name, Body, Link, Owner, LocalArrival,
  Creation, ProofOfWork, Signature
) VALUES (
  :Fingerprint, :Board, :Name, :Body, :Link, :Owner, :LocalArrival,
  :Creation, :ProofOfWork, :Signature
)`

// Immutable
var postInsert = `INSERT IGNORE INTO Posts
(
  Fingerprint, Board, Thread, Parent, Body, Owner, LocalArrival,
  Creation, ProofOfWork, Signature
) VALUES (
  :Fingerprint, :Board, :Thread, :Parent, :Body, :Owner, :LocalArrival,
  :Creation, :ProofOfWork, :Signature
)`

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
          :LocalArrival AS LocalArrival
          ) AS Candidate
  LEFT JOIN Votes ON Candidate.Fingerprint = Votes.Fingerprint
  WHERE (Candidate.LastUpdate > Votes.LastUpdate AND Candidate.LastUpdate > Votes.Creation)
  OR Votes.Fingerprint IS NULL`

// Address insert is immutable. This is used for when a node receives data from an address from a node that is not at the aforementioned address. In other words, an address object coming from a third party node not at that address cannot change an existing address saved in the database.
var addressInsert = `INSERT IGNORE INTO Addresses
(
  Location, Sublocation, Port, IPType, AddressType, LastOnline,
  ProtocolVersionMajor, ProtocolVersionMinor, ProtocolExtensions,
  ClientVersionMajor, ClientVersionMinor, ClientVersionPatch, ClientName,
  LocalArrival
) VALUES (
  :Location, :Sublocation, :Port,:IPType, :AddressType, :LastOnline,
  :ProtocolVersionMajor, :ProtocolVersionMinor, :ProtocolExtensions,
  :ClientVersionMajor, :ClientVersionMinor, :ClientVersionPatch, :ClientName,
  :LocalArrival
)`

// Address update insert is mutable. This is used when the node connects to the address itself. Example: When a node connects to 256.253.231.123:8080, it will update the entry for that address with the data coming from the remote node. This is the only way to mutate an address object.
var addressUpdateInsert = `REPLACE INTO Addresses
(
  Location, Sublocation, Port, IPType, AddressType, LastOnline,
  ProtocolVersionMajor, ProtocolVersionMinor, ProtocolExtensions,
  ClientVersionMajor, ClientVersionMinor, ClientVersionPatch, ClientName,
  LocalArrival
) VALUES (
  :Location, :Sublocation, :Port,:IPType, :AddressType, :LastOnline,
  :ProtocolVersionMajor, :ProtocolVersionMinor, :ProtocolExtensions,
  :ClientVersionMajor, :ClientVersionMinor, :ClientVersionPatch, :ClientName,
  :LocalArrival
)`

// Key insert does insert or replace without checking because we're handling the logic that decides whether we should update or not in the database layer.
var keyInsert = `REPLACE INTO PublicKeys
  (
    Fingerprint, Type, PublicKey, Name, Info, LocalArrival,
    Creation, ProofOfWork, Signature,
    LastUpdate, UpdateProofOfWork, UpdateSignature
  ) VALUES (
    :Fingerprint, :Type, :PublicKey, :Name, :Info, :LocalArrival,
    :Creation, :ProofOfWork, :Signature,
    :LastUpdate, :UpdateProofOfWork, :UpdateSignature
  )`

// CurrencyAddresses are mutable, but the condition of mutation is handled in the application layer. The only place the REPLACE could trigger is change of CurrencyCode. The Address and KeyFingerprint are identity columns, so anything with different data on those will be committed as a new item.
var currencyAddressInsert = `REPLACE INTO CurrencyAddresses
(
  KeyFingerprint, CurrencyCode, Address
) VALUES (
  :KeyFingerprint, :CurrencyCode, :Address
)`

// This is used when an user removes the currency address on his own key.
var currencyAddressDelete = `DELETE FROM CurrencyAddresses WHERE KeyFingerprint = :KeyFingerprint AND Address = :Address`

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
          :LocalArrival AS LocalArrival
          ) AS Candidate
  LEFT JOIN Truststates ON Candidate.Fingerprint = Truststates.Fingerprint
  WHERE (Candidate.LastUpdate > Truststates.LastUpdate AND Candidate.LastUpdate > Truststates.Creation)
  OR Truststates.Fingerprint IS NULL`
