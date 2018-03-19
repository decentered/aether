package persistence_test

import (
	"aether-core/io/api"
	"aether-core/io/persistence"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

// Infrastructure, setup and teardown

func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	teardown()
	os.Exit(exitVal)
}

func setup() {
	// Create the database.
	persistence.CreateDatabase()
	// Insert some basic data.
	createNodeData()
}

func teardown() {
	persistence.DeleteDatabase()
	// Mind that this isn't as optional as you think. There are some tests, especially those related to updates below that need the database to be clean. Because we automatically switch to update when something is there, it breaks the creation tests (they end up being updates.)
}

func ValidateTest(expected interface{}, actual interface{}, t *testing.T) {
	if actual != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
	}
}

func createNodeData() {
	var b api.Board
	var b2 api.Board
	var bo api.BoardOwner
	var bo2 api.BoardOwner
	var t api.Thread
	var p api.Post
	var v api.Vote
	var a api.Address
	var s api.Subprotocol
	var k api.Key
	var k2 api.Key
	var ca api.CurrencyAddress
	var ts api.Truststate
	// Insert the data for board. Only the bare minimum of data required.
	k.Fingerprint = "2389749283fasdf"
	k2.Fingerprint = "asdfasfdfa9023423"
	ca.Address = "1241245342513412341234123412341234"
	ca.CurrencyCode = "XXX"
	k.CurrencyAddresses = append(k.CurrencyAddresses, ca)
	k.Key = "public key"
	k.Creation = 1
	k.ProofOfWork = "pow"
	k.Signature = "sig"
	k.Type = "key type"

	b.Fingerprint = "my board fingerprint"
	b.Name = "alice"
	b.Creation = 1
	b.ProofOfWork = "pow"
	b.Owner = k.Fingerprint
	bo.KeyFingerprint = k.Fingerprint
	bo.Level = 255
	b.BoardOwners = append(b.BoardOwners, bo)

	b2.Fingerprint = "my board fingerprint_second"
	b2.Name = "alice"
	b2.Creation = 1
	b2.ProofOfWork = "pow"
	bo2.KeyFingerprint = k.Fingerprint
	bo2.Level = 1
	b2.BoardOwners = append(b2.BoardOwners, bo2)

	t.Fingerprint = "my thread fingerprint"
	t.Board = b.Fingerprint
	t.Name = "alice"
	t.Creation = 1
	t.ProofOfWork = "pow"

	p.Fingerprint = "my post fingerprint"
	p.Board = b.Fingerprint
	p.Thread = p.Fingerprint
	p.Parent = p.Fingerprint
	p.Body = "a"
	p.Creation = 1
	p.ProofOfWork = "pow"

	v.Fingerprint = "my vote fingerprint"
	v.Board = b.Fingerprint
	v.Thread = t.Fingerprint
	v.Target = t.Fingerprint
	v.Owner = k2.Fingerprint
	v.Type = 1
	v.Creation = 1
	v.Signature = "sig"
	v.ProofOfWork = "pow"

	s.Name = "c0"
	s.VersionMajor = 1
	s.VersionMinor = 0
	s.SupportedEntities = []string{"board", "thread", "post", "vote", "key", "truststate"}
	a.Location = "www.example.com"
	a.Sublocation = "example"
	a.Port = 8090
	a.LocationType = 1
	a.LastOnline = 1
	a.Protocol.VersionMajor = 1
	a.Protocol.Subprotocols = []api.Subprotocol{s}
	a.Client.VersionMajor = 1
	a.Client.ClientName = "client name"

	ts.Fingerprint = "my truststate fingerprint"
	ts.Target = k.Fingerprint
	ts.Owner = k2.Fingerprint
	ts.Type = 1
	ts.Creation = 1
	ts.Signature = "sig"
	ts.ProofOfWork = "pow"
	// Create the container and insert the items into
	var batch []interface{}
	batch = append(batch, b)
	batch = append(batch, b2)
	batch = append(batch, t)
	batch = append(batch, p)
	batch = append(batch, v)
	batch = append(batch, a)
	batch = append(batch, k)
	batch = append(batch, ts)
	err := persistence.BatchInsert(batch)
	if err != nil {
		log.Fatal(err)
	}

}

// Tests

// High level API (Read) tests

func TestRead_Success(t *testing.T) {
	fp := api.Fingerprint("my board fingerprint")
	resp, err := persistence.Read("boards", []api.Fingerprint{api.Fingerprint(fp)}, []string{}, 0, 0)
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Boards) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp.Boards[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Fingerprint: '%s'", resp.Boards[0].Fingerprint)
	}
}

func TestRead_SingleEmbed_BoardEmbedThread_Success(t *testing.T) {
	fp := api.Fingerprint("my board fingerprint")
	resp, err := persistence.Read("boards", []api.Fingerprint{api.Fingerprint(fp)}, []string{"threads"}, 0, 0)
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Boards) != 1 || len(resp.Threads) != 1 {
		t.Errorf("Test failed, the response has missing data. Response: %#v", resp)
	} else if resp.Boards[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Fingerprint: '%s'", resp.Boards[0].Fingerprint)
	}
}

func TestRead_MultiEmbed_Success(t *testing.T) {
	// This test checks for embedding order.
	// The embed provided below are Boards main, and thread and key embeds. The correct behaviour is that key embed runs last, and it includes keys referred to by BOTH the board and the thread.

	// 1 board, 1 thread, 2 keys
	var bOne api.Board
	var tOne api.Thread
	var kOne api.Key
	var kTwo api.Key

	kOne.Fingerprint = "first key fingerprint"
	kOne.Key = "public key"
	kOne.Creation = 1
	kOne.ProofOfWork = "pow"
	kOne.Signature = "sig"
	kOne.Type = "key type"

	kTwo.Fingerprint = "second key fingerprint"
	kTwo.Key = "public key"
	kTwo.Creation = 1
	kTwo.ProofOfWork = "pow"
	kTwo.Signature = "sig"
	kTwo.Type = "key type"

	bOne.Fingerprint = "my board fingerprint multi entity batch test"
	bOne.Name = "alice"
	bOne.Creation = 1
	bOne.ProofOfWork = "pow"
	bOne.Owner = kOne.Fingerprint

	tOne.Fingerprint = "my thread fingerprint multi entity batch test"
	tOne.Board = bOne.Fingerprint
	tOne.Name = "alice"
	tOne.Creation = 1
	tOne.ProofOfWork = "pow"
	tOne.Owner = kTwo.Fingerprint

	var batch []interface{}
	batch = append(batch, bOne)
	batch = append(batch, tOne)
	batch = append(batch, kOne)
	batch = append(batch, kTwo)
	err := persistence.BatchInsert(batch)

	fp := api.Fingerprint("my board fingerprint multi entity batch test")
	resp, err := persistence.Read("boards", []api.Fingerprint{api.Fingerprint(fp)}, []string{"threads", "keys"}, 0, 0)
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Boards) != 1 || len(resp.Threads) != 1 || len(resp.Keys) != 2 {
		t.Errorf("Test failed, the response has different data than expected. Response: %#v", resp)
	} else if resp.Boards[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Fingerprint: '%s'", resp.Boards[0].Fingerprint)
	}
}

func TestRead_PostEmbedVote_Success(t *testing.T) {

	var p api.Post
	var v api.Vote

	p.Fingerprint = "my post fingerprint99"
	p.Board = "board fingerprint"
	p.Thread = "thread fingerprint"
	p.Parent = "thread fingerprint"
	p.Body = "a"
	p.Creation = 1
	p.ProofOfWork = "pow"

	v.Fingerprint = "my vote fingerprint100"
	v.Board = "board fingerprint"
	v.Thread = "thread fingerprint"
	v.Target = p.Fingerprint
	v.Owner = "key fingerprint"
	v.Type = 1
	v.Creation = 1
	v.Signature = "sig"
	v.ProofOfWork = "pow"

	var batch []interface{}
	batch = append(batch, p)
	batch = append(batch, v)
	err := persistence.BatchInsert(batch)

	fp := api.Fingerprint("my post fingerprint99")
	resp, err := persistence.Read("posts", []api.Fingerprint{api.Fingerprint(fp)}, []string{"votes"}, 0, 0)
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Posts) != 1 || len(resp.Votes) != 1 {
		t.Errorf("Test failed, the response has different data than expected. Response: %#v", resp)
	}
}

func TestRead_ThreadEmbedPost_Success(t *testing.T) {

	var p api.Post
	var t2 api.Thread

	p.Fingerprint = "my post fingerprint991"
	p.Board = "board fingerprint"
	p.Thread = "my thread fingerprint99"
	p.Parent = "my thread fingerprint99"
	p.Body = "a"
	p.Creation = 1
	p.ProofOfWork = "pow"

	t2.Fingerprint = "my thread fingerprint99"
	t2.Board = "board fingerprint"
	t2.Name = "alice"
	t2.Creation = 1
	t2.ProofOfWork = "pow"

	var batch []interface{}
	batch = append(batch, p)
	batch = append(batch, t2)
	err := persistence.BatchInsert(batch)

	fp := api.Fingerprint("my thread fingerprint99")
	resp, err := persistence.Read("threads", []api.Fingerprint{api.Fingerprint(fp)}, []string{"posts"}, 0, 0)
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Threads) != 1 || len(resp.Posts) != 1 {
		t.Errorf("Test failed, the response has different data than expected. Response: %#v", resp)
	}
}

func TestRead_TruststateEmbedKey_Success(t *testing.T) {

	var ts4 api.Truststate
	var k4 api.Key
	var ca4 api.CurrencyAddress

	k4.Fingerprint = "2389749283fasdf"
	ca4.Address = "1241245342513412341234123412341234"
	ca4.CurrencyCode = "XXX"
	k4.CurrencyAddresses = append(k4.CurrencyAddresses, ca4)
	k4.Key = "public key"
	k4.Creation = 1
	k4.ProofOfWork = "pow"
	k4.Signature = "sig"
	k4.Type = "key type"

	ts4.Fingerprint = "my truststate fingerprint99"
	ts4.Target = "an awesome user's key"
	ts4.Owner = k4.Fingerprint
	ts4.Type = 1
	ts4.Creation = 1
	ts4.Signature = "sig"
	ts4.ProofOfWork = "pow"

	var batch []interface{}
	batch = append(batch, ts4)
	batch = append(batch, k4)
	err := persistence.BatchInsert(batch)

	fp := api.Fingerprint("my truststate fingerprint99")
	resp, err := persistence.Read("truststates", []api.Fingerprint{api.Fingerprint(fp)}, []string{"keys"}, 0, 0)
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Truststates) != 1 || len(resp.Keys) != 1 {
		t.Errorf("Test failed, the response has different data than expected. Response: %#v", resp)
	}
}

func TestReadBasedOnArrivalTimeRange(t *testing.T) {

	var b10 api.Board
	b10.Fingerprint = "my board fingerprint102"
	b10.Name = "alice"
	b10.Creation = 1
	b10.ProofOfWork = "pow"
	b10.Owner = "board owner"
	var batch []interface{}
	batch = append(batch, b10)
	persistence.BatchInsert(batch)

	time.Sleep(1000 * time.Millisecond) // Wait a bit so we have a decent range.
	now := api.Timestamp(time.Now().Unix())
	// fmt.Printf("%#v\n", now)
	resp, err := persistence.Read("boards", []api.Fingerprint{}, []string{}, 0, now)
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Boards) == 0 {
		t.Errorf("Test failed, the response is empty.")
	}
}

// // Reader tests (Medium level ReadX)

func TestReadBoard_Success(t *testing.T) {
	fp := api.Fingerprint("my board fingerprint")
	resp, err := persistence.ReadBoards(
		[]api.Fingerprint{api.Fingerprint(fp)}, 0, 0)
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Fingerprint: '%s'", resp[0].Fingerprint)
	}
}

func TestReadMultipleBoards_Success(t *testing.T) {
	fp := api.Fingerprint("my board fingerprint")
	fp2 := api.Fingerprint("my board fingerprint_second")
	resp, err := persistence.ReadBoards(
		[]api.Fingerprint{fp, fp2}, 0, 0)
	// fmt.Printf("%#v\n", resp)
	// fmt.Printf("%#v\n", len(resp))
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Fingerprint: '%s'", resp[0].Fingerprint)
	} else if len(resp) != 2 {
		t.Errorf("The response received doesn't have the right number of items. Fingerprint: '%#v'", resp)
	}
}

func TestReadBoard_Empty(t *testing.T) {
	resp, err := persistence.ReadBoards(
		[]api.Fingerprint{api.Fingerprint("fake board fingerprint")}, 0, 0)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) > 0 {
		t.Errorf("Test failed, the response is expected to be empty, but is not. Response: '%s'", resp)
	}
}

func TestReadDBBoardOwner_Success(t *testing.T) {
	fp := api.Fingerprint("my board fingerprint")
	keyFp := api.Fingerprint("2389749283fasdf")
	resp, err := persistence.ReadDBBoardOwners(
		fp, keyFp)
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp[0].KeyFingerprint != keyFp {
		t.Errorf("The response received isn't the expected one. Key Fingerprint: '%s'", resp[0].KeyFingerprint)
	}
}

func TestReadDBBoardOwner_PartialData_Success(t *testing.T) {
	fp := api.Fingerprint("my board fingerprint")
	keyFp := api.Fingerprint("2389749283fasdf")
	resp, err := persistence.ReadDBBoardOwners(
		fp, "")
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp[0].KeyFingerprint != keyFp {
		t.Errorf("The response received isn't the expected one. Key Fingerprint: '%s'", resp[0].KeyFingerprint)
	}
}
func TestReadDBBoardOwner_Empty(t *testing.T) {
	resp, err := persistence.ReadDBBoardOwners(
		"fake board owner fingerprint", "fake key fingerprint")
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) > 0 {
		t.Errorf("Test failed, the response is expected to be empty, but is not. Response: '%s'", resp)
	}
}

func TestReadThread_Success(t *testing.T) {
	fp := api.Fingerprint("my thread fingerprint")
	resp, err := persistence.ReadThreads([]api.Fingerprint{fp}, 0, 0)
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Fingerprint: '%s'", resp[0].Fingerprint)
	}
}

func TestReadThread_Empty(t *testing.T) {
	resp, err := persistence.ReadThreads([]api.Fingerprint{"fake thread fingerprint"}, 0, 0)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) > 0 {
		t.Errorf("Test failed, the response is expected to be empty, but is not. Response: '%s'", resp)
	}
}

func TestReadPost_Success(t *testing.T) {
	fp := api.Fingerprint("my post fingerprint")
	resp, err := persistence.ReadPosts([]api.Fingerprint{fp}, 0, 0)
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Fingerprint: '%s'", resp[0].Fingerprint)
	}
}

func TestReadPost_Empty(t *testing.T) {
	resp, err := persistence.ReadPosts([]api.Fingerprint{"fake post fingerprint"}, 0, 0)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) > 0 {
		t.Errorf("Test failed, the response is expected to be empty, but is not. Response: '%s'", resp)
	}
}

func TestReadVote_Success(t *testing.T) {
	fp := api.Fingerprint("my vote fingerprint")
	resp, err := persistence.ReadVotes([]api.Fingerprint{fp}, 0, 0)
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Fingerprint: '%s'", resp[0].Fingerprint)
	}
}

func TestReadVote_Empty(t *testing.T) {
	resp, err := persistence.ReadVotes([]api.Fingerprint{"fake vote fingerprint"}, 0, 0)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) > 0 {
		t.Errorf("Test failed, the response is expected to be empty, but is not. Response: '%s'", resp)
	}
}

func TestReadAddress_Success(t *testing.T) {
	loc := api.Location("www.example.com")
	subloc := api.Location("example")
	port := uint16(8090)
	resp, err := persistence.ReadAddresses(
		loc, subloc, port, 0, 0, 0, 0, 0, "")
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp[0].Location != loc {
		t.Errorf("The response received isn't the expected one. Location: '%s'", resp[0].Location)
	}
}

func TestFirstPartyInsertReadAddress_Success(t *testing.T) {
	// Insert test for new type address
	var a2 api.Address
	var s1 api.Subprotocol
	s1.Name = "dweb"
	s1.VersionMajor = 1
	s1.VersionMinor = 0
	s1.SupportedEntities = []string{"page"}
	var s2 api.Subprotocol
	s2.Name = "c0"
	s2.VersionMajor = 1
	s2.VersionMinor = 0
	s2.SupportedEntities = []string{"board", "thread", "post", "vote", "key", "truststate"}
	a2.Location = "www.example33.com"
	a2.Sublocation = "example33"
	a2.Port = 1111
	a2.LocationType = 1
	a2.LastOnline = 1
	a2.Protocol.VersionMajor = 1
	a2.Protocol.Subprotocols = []api.Subprotocol{s1, s2}
	a2.Client.VersionMajor = 1
	a2.Client.ClientName = "client name"
	addressSet := []api.Address{a2}
	persistence.InsertOrUpdateAddresses(&addressSet)

	loc := api.Location("www.example33.com")
	subloc := api.Location("example33")
	port := uint16(1111)

	resp, err := persistence.ReadAddresses(
		loc, subloc, port, 0, 0, 0, 0, 0, "")
	fmt.Printf("%#v\n", resp)
	if resp[0].Protocol.Subprotocols[0].Name != "c0" {
		t.Errorf(fmt.Sprintf("Test failed, the subprotocol information has not been committed. Response: %#v", resp))
	}
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp[0].Location != loc {
		t.Errorf("The response received isn't the expected one. Location: '%s'", resp[0].Location)
	}

}

func TestReadAddress_Empty(t *testing.T) {
	resp, err := persistence.ReadAddresses(
		"fake loc", "fake subloc", 9090, 0, 0, 0, 0, 0, "")
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) > 0 {
		t.Errorf("Test failed, the response is expected to be empty, but is not. Response: '%s'", resp)
	}
}

func TestReadKey_Success(t *testing.T) {
	fp := api.Fingerprint("2389749283fasdf")
	resp, err := persistence.ReadKeys([]api.Fingerprint{fp}, 0, 0)
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Fingerprint: '%s'", resp[0].Fingerprint)
	}
}

func TestReadKey_Empty(t *testing.T) {
	resp, err := persistence.ReadKeys([]api.Fingerprint{"fake key fingerprint"}, 0, 0)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) > 0 {
		t.Errorf("Test failed, the response is expected to be empty, but is not. Response: '%s'", resp)
	}
}

func TestReadDBCurrencyAddress_Success(t *testing.T) {
	keyFp := api.Fingerprint("2389749283fasdf")
	currAddr := "1241245342513412341234123412341234"
	resp, err := persistence.ReadDBCurrencyAddresses(
		keyFp, currAddr)
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp[0].KeyFingerprint != keyFp {
		t.Errorf("The response received isn't the expected one. Key fingerprint: '%s'", resp[0].KeyFingerprint)
	}
}

func TestReadDBCurrencyAddress_PartialData_Success(t *testing.T) {
	keyFp := api.Fingerprint("2389749283fasdf")
	currAddr := "" // empty means we don't know about it.
	resp, err := persistence.ReadDBCurrencyAddresses(
		keyFp, currAddr)
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp[0].KeyFingerprint != keyFp {
		t.Errorf("The response received isn't the expected one. Key fingerprint: '%s'", resp[0].KeyFingerprint)
	}
}

func TestReadDBCurrencyAddress_Empty(t *testing.T) {
	resp, err := persistence.ReadDBCurrencyAddresses(
		"fake key fingerprint", "fake addr")
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) > 0 {
		t.Errorf("Test failed, the response is expected to be empty, but is not. Response: '%s'", resp)
	}
}

func TestReadTruststate_Success(t *testing.T) {
	fp := api.Fingerprint("my truststate fingerprint")
	resp, err := persistence.ReadTruststates([]api.Fingerprint{fp}, 0, 0)
	// fmt.Printf("%#v\n", resp)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Fingerprint: '%s'", resp[0].Fingerprint)
	}
}

func TestReadTruststate_Empty(t *testing.T) {
	resp, err := persistence.ReadTruststates([]api.Fingerprint{"fake truststate fingerprint"}, 0, 0)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp) > 0 {
		t.Errorf("Test failed, the response is expected to be empty, but is not. Response: '%s'", resp)
	}
}

func TestDbToApi_Success(t *testing.T) {
	var ts persistence.DbTruststate
	ts.Fingerprint = "my awesome truststate fingerprint"
	ts.Target = "my target key"
	ts.Owner = "my owner's key fingerprint"
	ts.Domains = "my first domain fingerprint, my second domain fingerprint, my third domain fingerprint"
	apiObj, err := persistence.DBtoAPI(ts)
	obj := apiObj.(api.Truststate)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(obj.Domains) == 0 {
		t.Errorf("Test failed, the domains response is empty.")
	} else if obj.Fingerprint != ts.Fingerprint {
		t.Errorf("The response received isn't the expected one. Fingerprint: '%s'", obj.Fingerprint)
	}
}

func TestDbToApi_ItemLengthLongerThanAllowed(t *testing.T) {
	var ts persistence.DbTruststate
	ts.Fingerprint = "my awesome truststate fingerprint"
	ts.Target = "my target key"
	ts.Owner = "my owner's key fingerprint"
	ts.Domains = "my first domain fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint fingerprint, my second domain fingerprint, my third domain fingerprint"
	_, err := persistence.DBtoAPI(ts)
	errMessage := "This string is too long for this field."
	if err == nil {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestDbToApi_RepeatedItems(t *testing.T) {
	var ts persistence.DbTruststate
	ts.Fingerprint = "my awesome truststate fingerprint"
	ts.Target = "my target key"
	ts.Owner = "my owner's key fingerprint"
	ts.Domains = "alice,alice,alice"
	_, err := persistence.DBtoAPI(ts)
	errMessage := "This list includes items that are duplicates."
	if err == nil {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestDbToApi_TooManyItems(t *testing.T) {
	var ts persistence.DbTruststate
	ts.Fingerprint = "my awesome truststate fingerprint"
	ts.Target = "my target key"
	ts.Owner = "my owner's key fingerprint"
	ts.Domains = "a,a1,a2,a3,a4,a5,a6,a7,a8,a9,a10,a11,a12,a13,a14,a15,a16,a17,a18,a19,a20,a21,a22,a23,a24,a25,a26,a27,a28,a29,a30,a31,a32,a33,a34,a35,a36,a37,a38,a39,a40,a41,a42,a43,a44,a45,a46,a47,a48,a49,a50,a51,a52,a53,a54,a55,a56,a57,a58,a59,a60,a61,a62,a63,a64,a65,a66,a67,a68,a69,a70,a71,a72,a73,a74,a75,a76,a77,a78,a79,a80,a81,a82,a83,a84,a85,a86,a87,a88,a89,a90,a91,a92,a93,a94,a95,a96,a97,a98,a99,a100,a101,a102,a103,a104,a105"
	_, err := persistence.DBtoAPI(ts)
	errMessage := "The string provided has too many items"
	if err == nil {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestApiToDb_Success(t *testing.T) {
	var a api.Address
	a.Location = "www.example.com"
	a.Sublocation = "hello"
	a.Port = uint16(8090)
	var s api.Subprotocol
	s.Name = "c0"
	s.VersionMajor = 1
	s.VersionMinor = 0
	s.SupportedEntities = []string{"board", "thread", "post", "vote", "key", "truststate"}
	a.Protocol.Subprotocols = []api.Subprotocol{s}
	addressPack, err := persistence.APItoDB(a)
	obj := addressPack.(persistence.AddressPack)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(obj.Subprotocols) == 0 {
		t.Errorf("Test failed, the locations response is empty.")
	} else if obj.Address.Location != a.Location {
		t.Errorf("The response received isn't the expected one. Location: '%s'", obj.Address.Location)
	}
}

// TOFIX: make sure this detection works.
func TestApiToDb_RepeatedItems(t *testing.T) {
	var a api.Address
	a.Location = "www.example.com"
	a.Sublocation = "hello"
	a.Port = uint16(8090)
	var s1 api.Subprotocol
	s1.Name = "c0"
	s1.VersionMajor = 1
	s1.VersionMinor = 0
	s1.SupportedEntities = []string{"board", "board"}
	a.Protocol.Subprotocols = []api.Subprotocol{s1}
	_, err := persistence.APItoDB(a)
	errMessage := "This list includes items that are duplicates."
	if err == nil {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestApiToDb_TooManyItems(t *testing.T) {
	var a api.Address
	a.Location = "www.example.com"
	a.Sublocation = "hello"
	a.Port = uint16(8090)
	var s api.Subprotocol
	s.Name = "c0"
	s.VersionMajor = 1
	s.VersionMinor = 0
	s.SupportedEntities = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "111", "211", "311", "411", "511", "611", "711", "811", "911", "1011", "1111", "1211", "1311", "1411", "1511", "1611", "1711", "1811", "1911", "2011", "2111", "2211", "2311", "2411", "2511", "2611", "2711", "2811", "2911", "3011", "11111", "21111", "31111", "41111", "51111", "61111", "71111", "81111", "91111", "101111", "111111", "121111", "131111", "141111", "151111", "161111", "171111", "181111", "191111", "201111", "211111", "221111", "231111", "241111", "251111", "261111", "271111", "281111", "291111", "301111", "1111111", "2111111", "3111111", "4111111", "5111111", "6111111", "7111111", "8111111", "9111111", "10111111", "11111111", "12111111", "13111111", "14111111", "15111111", "16111111", "17111111", "18111111", "19111111", "20111111", "21111111", "22111111", "23111111", "24111111", "25111111", "26111111", "27111111", "28111111", "29111111", "30111111"}
	a.Protocol.Subprotocols = []api.Subprotocol{s}
	_, err := persistence.APItoDB(a)
	errMessage := "The string slice provided has too many items."
	if err == nil {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestApiToDb_ItemLengthLongerThanAllowed(t *testing.T) {
	var a api.Address
	a.Location = "www.example.com"
	a.Sublocation = "hello"
	a.Port = uint16(8090)
	var s api.Subprotocol
	s.Name = "c0"
	s.VersionMajor = 1
	s.VersionMinor = 0
	s.SupportedEntities = []string{"boaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaard"}
	a.Protocol.Subprotocols = []api.Subprotocol{s}
	_, err := persistence.APItoDB(a)
	errMessage := "This string is too long for this field."
	if err == nil {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

// // Writer tests

func TestSingleInsert_Success(t *testing.T) {
	fp := api.Fingerprint("my awesome vote fingerprint")
	var vote api.Vote
	vote.Fingerprint = fp
	vote.Board = "board fingerprint"
	vote.Thread = "thread fingerprint"
	vote.Target = "target fingerprint"
	vote.Owner = "owner fingerprint"
	vote.Type = 1
	vote.Creation = 1
	vote.Signature = "sig"
	vote.ProofOfWork = "pow"
	err := persistence.BatchInsert([]interface{}{vote})
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	}
	resp, err2 := persistence.ReadVotes([]api.Fingerprint{fp}, 0, 0)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	} else if len(resp) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Fingerprint: '%s'", resp[0].Fingerprint)
	}
}

func TestSingleInsert_Duplicate(t *testing.T) {
	fp := api.Fingerprint("my awesome vote fingerprint2")
	var vote api.Vote
	vote.Fingerprint = fp
	vote.Board = "board fingerprint"
	vote.Thread = "thread fingerprint"
	vote.Target = "target fingerprint"
	vote.Owner = "owner fingerprint"
	vote.Type = 1
	vote.Creation = 1
	vote.Signature = "sig"
	vote.ProofOfWork = "pow"
	var vote2 api.Vote
	vote2.Fingerprint = fp
	vote2.Board = "board fingerprint"
	vote2.Thread = "thread fingerprint"
	vote2.Target = "target fingerprint"
	vote2.Owner = "owner fingerprint"
	vote2.Type = 1
	vote2.Creation = 1
	vote2.Signature = "sig"
	vote2.ProofOfWork = "pow"
	err := persistence.BatchInsert([]interface{}{vote, vote2})
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	}
	resp, err2 := persistence.ReadVotes([]api.Fingerprint{fp}, 0, 0)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	} else if len(resp) > 1 {
		t.Errorf("Test failed, the response has more than one vote.")
	} else if resp[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Fingerprint: '%s'", resp[0].Fingerprint)
	}
}

func TestBatchInsert_Success(t *testing.T) {
	fp := api.Fingerprint("my awesome vote fingerprint3")
	fp2 := api.Fingerprint("my awesome vote fingerprint4")

	var vote api.Vote
	vote.Fingerprint = fp
	vote.Board = "board fingerprint"
	vote.Thread = "thread fingerprint"
	vote.Target = "target fingerprint"
	vote.Owner = "owner fingerprint"
	vote.Type = 1
	vote.Creation = 1
	vote.Signature = "sig"
	vote.ProofOfWork = "pow"
	var vote2 api.Vote
	vote2.Fingerprint = fp2
	vote2.Board = "board fingerprint"
	vote2.Thread = "thread fingerprint"
	vote2.Target = "target fingerprint"
	vote2.Owner = "owner fingerprint"
	vote2.Type = 1
	vote2.Creation = 1
	vote2.Signature = "sig"
	vote2.ProofOfWork = "pow"

	err := persistence.BatchInsert([]interface{}{vote, vote2})
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	}
	// Check for first
	resp, err2 := persistence.ReadVotes([]api.Fingerprint{fp}, 0, 0)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	} else if len(resp) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Fingerprint: '%s'", resp[0].Fingerprint)
	}
	// Check for second
	resp2, err3 := persistence.ReadVotes([]api.Fingerprint{fp2}, 0, 0)
	if err3 != nil {
		t.Errorf("Test failed, err: '%s'", err3)
	} else if len(resp2) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp2[0].Fingerprint != fp2 {
		t.Errorf("The response received isn't the expected one. Fingerprint: '%s'", resp2[0].Fingerprint)
	}
}

func TestInsert_MultipleTypes_Success(t *testing.T) {
	tfp := api.Fingerprint("my awesome truststate fp")
	addressLoc := api.Location("www.example2.com")
	addressSubloc := api.Location("subloc")
	addressPort := uint16(9080)
	var ts api.Truststate
	var a api.Address
	ts.Fingerprint = tfp
	ts.Target = "ts target"
	ts.Owner = "ts owner"
	ts.Type = 1
	ts.Creation = 1
	ts.Signature = "sig"
	ts.ProofOfWork = "pow"

	var s api.Subprotocol
	s.Name = "c0"
	s.VersionMajor = 1
	s.VersionMinor = 0
	s.SupportedEntities = []string{"board", "thread", "post", "vote", "key", "truststate"}

	a.Location = addressLoc
	a.Sublocation = addressSubloc
	a.Port = addressPort
	a.LocationType = 1
	a.LastOnline = 1
	a.Protocol.VersionMajor = 1
	a.Protocol.Subprotocols = []api.Subprotocol{s, s}
	a.Client.VersionMajor = 1
	a.Client.ClientName = "client name"

	err := persistence.BatchInsert([]interface{}{ts, a})
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	}
	// Check for first
	resp, err2 := persistence.ReadAddresses(addressLoc, addressSubloc, addressPort, 0, 0, 0, 0, 0, "")
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	} else if len(resp) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp[0].Location != addressLoc {
		t.Errorf("The response received isn't the expected one. Address: '%s'", resp[0])
	}
	// Check for second
	resp2, err3 := persistence.ReadTruststates([]api.Fingerprint{tfp}, 0, 0)
	if err3 != nil {
		t.Errorf("Test failed, err: '%s'", err3)
	} else if len(resp2) == 0 {
		t.Errorf("Test failed, the response is empty.")
	} else if resp2[0].Fingerprint != tfp {
		t.Errorf("The response received isn't the expected one. Fingerprint: '%s'", resp2[0].Fingerprint)
	}
}

func TestInsert_ItemsWithUpdates_Board_SimpleField_Success(t *testing.T) {
	// Insert a board.
	var b api.Board
	fp := api.Fingerprint("my cool board fingerprint")
	b.Fingerprint = fp
	b.Creation = 2
	b.Name = "alice"
	b.ProofOfWork = "pow"
	// var bo api.BoardOwner
	// bo.KeyFingerprint = "key fingerprint"
	// bo.Level = 1
	// b.BoardOwners = append(b.BoardOwners, bo)

	err := persistence.BatchInsert([]interface{}{b})
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	}
	resp, err2 := persistence.ReadBoards(
		[]api.Fingerprint{api.Fingerprint(fp)}, 0, 0)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	}
	if resp[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Board: '%s'", resp[0])
	}

	// // Checking for field changes based on last update.

	// Change a board, but make last update earlier than creation. This means it won't enter the database.
	b.Description = "hello there!"
	b.LastUpdate = 1
	err9 := persistence.BatchInsert([]interface{}{b})
	if err9 != nil {
		t.Errorf("Test failed, err: '%s'", err9)
	}
	resp3, err10 := persistence.ReadBoards(
		[]api.Fingerprint{api.Fingerprint(fp)}, 0, 0)
	if err10 != nil {
		t.Errorf("Test failed, err: '%s'", err10)
	}
	if len(resp3[0].Description) != 0 {
		t.Errorf("The description shouldn't have gotten in because it's from an update that is earlier than creation. Board: '%#v\n'", resp3[0])
		t.Fatal()
	}

	// Change a board, but make last update same as creation. This means it won't enter the database.
	b.Description = "hello there!2"
	b.LastUpdate = 2
	err11 := persistence.BatchInsert([]interface{}{b})
	if err11 != nil {
		t.Errorf("Test failed, err: '%s'", err11)
	}
	resp4, err12 := persistence.ReadBoards(
		[]api.Fingerprint{api.Fingerprint(fp)}, 0, 0)
	if err12 != nil {
		t.Errorf("Test failed, err: '%s'", err12)
	}
	if len(resp4[0].Description) != 0 {
		t.Errorf("The description shouldn't have gotten in because it's from an update that has the same timestamp as creation. Board: '%#v\n'", resp4[0])
		t.Fatal()
	}

	// Change a board, but make last update after creation. This means it should enter the database.
	b.Description = "hello there!3"
	b.LastUpdate = 3
	err13 := persistence.BatchInsert([]interface{}{b})
	if err13 != nil {
		t.Errorf("Test failed, err: '%s'", err13)
	}
	resp5, err14 := persistence.ReadBoards(
		[]api.Fingerprint{api.Fingerprint(fp)}, 0, 0)
	if err14 != nil {
		t.Errorf("Test failed, err: '%s'", err14)
	}
	if len(resp5[0].Description) == 0 {
		t.Errorf("The description should have gotten in because it's from an update that has a later timestamp than creation. Board: '%#v\n'", resp5[0])
		t.Fatal()
	}
}

func TestInsert_ItemsWithUpdates_Board_SubObject_Success(t *testing.T) {
	// Insert a board.
	var b api.Board
	b.Fingerprint = "my board fingerprint"
	b.Name = "alice"
	b.Creation = 1
	b.ProofOfWork = "pow"
	var bo1 api.BoardOwner
	bo1.KeyFingerprint = "key fingerprint1"
	bo1.Level = 1
	var bo2 api.BoardOwner
	bo2.KeyFingerprint = "key fingerprint2"
	bo2.Level = 1
	var bo3 api.BoardOwner
	bo3.KeyFingerprint = "key fingerprint3"
	bo3.Level = 1
	fp := api.Fingerprint("my cool board fingerprint subobject test")
	b.Fingerprint = fp
	b.BoardOwners = []api.BoardOwner{bo1, bo2}
	b.Creation = 2
	err := persistence.BatchInsert([]interface{}{b})
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	}
	resp, err2 := persistence.ReadBoards(
		[]api.Fingerprint{api.Fingerprint(fp)}, 0, 0)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	}
	if resp[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Board: '%s'", resp[0])
	}

	// // Checking sub entity changes based on last update.

	// Change a board, but make last update earlier than creation. This means it won't enter the database.
	b.BoardOwners = []api.BoardOwner{bo1, bo2, bo3}
	b.LastUpdate = 1
	err3 := persistence.BatchInsert([]interface{}{b})
	if err3 != nil {
		t.Errorf("Test failed, err: '%s'", err3)
	}
	resp2, err4 := persistence.ReadBoards(
		[]api.Fingerprint{api.Fingerprint(fp)}, 0, 0)
	if err4 != nil {
		t.Errorf("Test failed, err: '%s'", err4)
	}
	if len(resp2[0].BoardOwners) > 2 {
		t.Errorf("The board owner shouldn't have gotten in because it's from an update that is earlier than creation. Current Board: '%#v\n', Attempted board: '%#v\n'", resp2[0], b)
		t.Fatal()
	}
	// Now change the last update to be the same as creation. This also should not go in.
	b.LastUpdate = 2
	err5 := persistence.BatchInsert([]interface{}{b})
	if err5 != nil {
		t.Errorf("Test failed, err: '%s'", err5)
	}
	resp3, err6 := persistence.ReadBoards(
		[]api.Fingerprint{api.Fingerprint(fp)}, 0, 0)
	if err6 != nil {
		t.Errorf("Test failed, err: '%s'", err6)
	}
	if len(resp3[0].BoardOwners) > 2 {
		t.Errorf("The board owner shouldn't have gotten in because it's from an update that is the same date as creation. Board: '%#v\n'", resp3[0])
		t.Fatal()
	}
	// Now change the last update to be after creation. This should go in.
	b.LastUpdate = 3
	err7 := persistence.BatchInsert([]interface{}{b})
	if err7 != nil {
		t.Errorf("Test failed, err: '%s'", err7)
	}
	resp4, err8 := persistence.ReadBoards(
		[]api.Fingerprint{api.Fingerprint(fp)}, 0, 0)
	if err8 != nil {
		t.Errorf("Test failed, err: '%s'", err8)
	}
	if len(resp4[0].BoardOwners) < 3 {
		fmt.Println(resp4[0].BoardOwners)
		t.Errorf("The board owner should have gotten in (but did not) because it's from an update that is later than creation. Current Board: '%#v\n', Attempted board: '%#v\n'", resp4[0], b)
		t.Fatal()
	}
}

func TestInsert_ItemsWithUpdates_Key_SimpleField_Success(t *testing.T) {
	// Insert a key.
	var k api.Key
	fp := api.Fingerprint("my cool key fingerprint5")
	k.Fingerprint = fp
	k.Creation = 2
	k.Key = "public key"
	k.ProofOfWork = "pow"
	k.Signature = "sig"
	k.Type = "key type"
	err := persistence.BatchInsert([]interface{}{k})
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	}
	resp, err2 := persistence.ReadKeys([]api.Fingerprint{fp}, 0, 0)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	}
	if resp[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Board: '%s'", resp[0])
	}

	// // Checking for field changes based on last update.

	// Change a board, but make last update earlier than creation. This means it won't enter the database.
	k.Name = "hola!"
	k.LastUpdate = 1
	err9 := persistence.BatchInsert([]interface{}{k})
	if err9 != nil {
		t.Errorf("Test failed, err: '%s'", err9)
	}
	resp3, err10 := persistence.ReadKeys([]api.Fingerprint{fp}, 0, 0)
	if err10 != nil {
		t.Errorf("Test failed, err: '%s'", err10)
	}
	if len(resp3[0].Name) != 0 {
		t.Errorf("The name shouldn't have gotten in because it's from an update that is earlier than creation. Key: '%#v\n'", resp3[0])
		t.Fatal()
	}

	// Change a board, but make last update same as creation. This means it won't enter the database.
	k.Name = "hola!2"
	k.LastUpdate = 2
	err11 := persistence.BatchInsert([]interface{}{k})
	if err11 != nil {
		t.Errorf("Test failed, err: '%s'", err11)
	}
	resp4, err12 := persistence.ReadKeys([]api.Fingerprint{fp}, 0, 0)
	if err12 != nil {
		t.Errorf("Test failed, err: '%s'", err12)
	}
	if len(resp4[0].Name) != 0 {
		t.Errorf("The name shouldn't have gotten in because it's from an update that has the same timestamp as creation. Key: '%#v\n'", resp4[0])
		t.Fatal()
	}

	// Change a board, but make last update after creation. This means it should enter the database.
	k.Name = "hola!3"
	k.LastUpdate = 3
	err13 := persistence.BatchInsert([]interface{}{k})
	if err13 != nil {
		t.Errorf("Test failed, err: '%s'", err13)
	}
	resp5, err14 := persistence.ReadKeys([]api.Fingerprint{fp}, 0, 0)
	if err14 != nil {
		t.Errorf("Test failed, err: '%s'", err14)
	}
	if len(resp5[0].Name) == 0 {
		t.Errorf("The name should have gotten in because it's from an update that has a later timestamp than creation. Key: '%#v\n'", resp5[0])
		t.Fatal()
	}
}

func TestInsert_ItemsWithUpdates_Key_SubObject_Success(t *testing.T) {
	// Insert a board.
	var k api.Key
	var ca1 api.CurrencyAddress
	var ca2 api.CurrencyAddress
	var ca3 api.CurrencyAddress
	ca1.Address = "my awesome address"
	ca1.CurrencyCode = "XXX"
	ca2.Address = "my awesome address2"
	ca2.CurrencyCode = "XXX"
	ca3.Address = "my awesome address3"
	ca3.CurrencyCode = "XXX"
	fp := api.Fingerprint("my cool key fingerprint subobject test")
	k.Fingerprint = fp
	k.CurrencyAddresses = []api.CurrencyAddress{ca1, ca2}
	k.Creation = 2
	k.Key = "public key"
	k.ProofOfWork = "pow"
	k.Signature = "sig"
	k.Type = "key type"

	err := persistence.BatchInsert([]interface{}{k})
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	}
	resp, err2 := persistence.ReadKeys([]api.Fingerprint{fp}, 0, 0)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	}
	if resp[0].Fingerprint != fp {
		t.Errorf("The response received isn't the expected one. Board: '%s'", resp[0])
	}

	// // Checking sub entity changes based on last update.

	// Change a board, but make last update earlier than creation. This means it won't enter the database.
	k.CurrencyAddresses = []api.CurrencyAddress{ca1, ca2, ca3}
	k.LastUpdate = 1
	err3 := persistence.BatchInsert([]interface{}{k})
	if err3 != nil {
		t.Errorf("Test failed, err: '%s'", err3)
	}
	resp2, err4 := persistence.ReadKeys([]api.Fingerprint{fp}, 0, 0)
	if err4 != nil {
		t.Errorf("Test failed, err: '%s'", err4)
	}
	if len(resp2[0].CurrencyAddresses) > 2 {
		t.Errorf("The currency address shouldn't have gotten in because it's from an update that is earlier than creation. Current key: '%#v\n', Attempted key: '%#v\n'", resp2[0], k)
		t.Fatal()
	}
	// Now change the last update to be the same as creation. This also should not go in.
	k.LastUpdate = 2
	err5 := persistence.BatchInsert([]interface{}{k})
	if err5 != nil {
		t.Errorf("Test failed, err: '%s'", err5)
	}
	resp3, err6 := persistence.ReadKeys([]api.Fingerprint{fp}, 0, 0)
	if err6 != nil {
		t.Errorf("Test failed, err: '%s'", err6)
	}
	if len(resp3[0].CurrencyAddresses) > 2 {
		t.Errorf("The currency address shouldn't have gotten in because it's from an update that is the same date as creation. Key: '%#v\n'", resp3[0])
		t.Fatal()
	}
	// Now change the last update to be after creation. This should go in.
	k.LastUpdate = 3
	err7 := persistence.BatchInsert([]interface{}{k})
	if err7 != nil {
		t.Errorf("Test failed, err: '%s'", err7)
	}
	resp4, err8 := persistence.ReadKeys([]api.Fingerprint{fp}, 0, 0)
	if err8 != nil {
		t.Errorf("Test failed, err: '%s'", err8)
	}
	if len(resp4[0].CurrencyAddresses) < 3 {
		t.Errorf("The currency address should have gotten in (but did not) because it's from an update that is later than creation. Current Key: '%#v\n', Attempted key: '%#v\n'", resp4[0], k)
		t.Fatal()
	}
}

func TestInsert_NonsensicalItem_Success(t *testing.T) {
	// Try to insert a DB item via Batch insert (which takes API items)
	var thr persistence.DbThread
	thr.Fingerprint = "my cool thread fingerprint"
	err := persistence.BatchInsert([]interface{}{thr})
	errMessage := "APItoDB only takes API (not DB) objects."
	if err == nil {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestInsert_AddRemoveBoardOwner(t *testing.T) {
	var b api.Board
	fp := api.Fingerprint("hello this is a board")
	b.Fingerprint = fp
	b.Name = "alice"
	b.Creation = 1
	b.ProofOfWork = "pow"
	var bo1 api.BoardOwner
	var bo2 api.BoardOwner
	var bo3 api.BoardOwner
	var bo4 api.BoardOwner
	bo1.KeyFingerprint = "hello"
	bo1.Level = 1
	bo2.KeyFingerprint = "hello2"
	bo2.Level = 1
	bo3.KeyFingerprint = "hello3"
	bo3.Level = 1
	bo4.KeyFingerprint = "hello4"
	bo4.Level = 1
	b.BoardOwners = []api.BoardOwner{bo1, bo2, bo3}
	err := persistence.BatchInsert([]interface{}{b})
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	}
	b.BoardOwners = []api.BoardOwner{bo1, bo2, bo4}
	b.LastUpdate = 1 // Remember, we need to do this otherwise it won't go in.
	err2 := persistence.BatchInsert([]interface{}{b})
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	}
	resp, err3 := persistence.ReadBoards(
		[]api.Fingerprint{api.Fingerprint(fp)}, 0, 0)
	// fmt.Printf("%#v\n", resp[0].BoardOwners)
	if err3 != nil {
		t.Errorf("Test failed, err: '%s'", err3)
	} else if len(resp[0].BoardOwners) > 3 {
		t.Errorf("This should have returned 3 board owners. Error: '%#v\n' Current Board: '%#v\n', Board owners: '%#v\n'", err3, resp[0], resp[0].BoardOwners)
	}
}

func TestInsert_AddRemoveAndEditBoardOwner(t *testing.T) {
	var b api.Board
	fp := api.Fingerprint("hello this is a board3")
	b.Fingerprint = fp
	b.Name = "alice"
	b.Creation = 1
	b.ProofOfWork = "pow"
	var bo1 api.BoardOwner
	var bo2 api.BoardOwner
	var bo3 api.BoardOwner
	var bo4 api.BoardOwner
	bo1.KeyFingerprint = "hello"
	bo1.Level = 1
	bo2.KeyFingerprint = "hello2"
	bo2.Level = 1
	bo3.KeyFingerprint = "hello3"
	bo3.Level = 1
	bo4.KeyFingerprint = "hello4"
	bo4.Level = 1
	b.BoardOwners = []api.BoardOwner{bo1, bo2, bo3}
	err := persistence.BatchInsert([]interface{}{b})
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	}
	// Change bo2. Let's see if the update persists.
	bo2.Level = 2
	b.BoardOwners = []api.BoardOwner{bo1, bo2, bo4}
	b.LastUpdate = 2 // Remember, we need to do this otherwise it won't go in.
	err2 := persistence.BatchInsert([]interface{}{b})
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	}
	resp, err3 := persistence.ReadBoards(
		[]api.Fingerprint{api.Fingerprint(fp)}, 0, 0)
	// fmt.Printf("%#v\n", resp[0].BoardOwners)
	if err3 != nil {
		t.Errorf("Test failed, err: '%s'", err3)
	} else if len(resp[0].BoardOwners) > 3 {
		t.Errorf("This should have returned 3 board owners. Error: '%#v\n' Current Board: '%#v\n', Board owners: '%#v\n'", err3, resp[0], resp[0].BoardOwners)
	} else {
		for _, bo := range resp[0].BoardOwners {
			if bo.KeyFingerprint == "hello2" {
				// This is the board we've changed.
				if bo.Level != 2 {
					t.Errorf("We've changed the level of this board, but it did not persist. Error: '%#v\n' Current Board: '%#v\n', Board owner: '%#v\n'", err3, resp[0], bo)
				}
			}
		}
	}
}

func TestInsert_DuplicateBoardOwner(t *testing.T) {
	var b api.Board
	fp := api.Fingerprint("hello this is a boardie")
	b.Fingerprint = fp
	b.Name = "alice"
	b.Creation = 1
	b.ProofOfWork = "pow"
	var bo1 api.BoardOwner
	var bo2 api.BoardOwner
	var bo3 api.BoardOwner
	bo1.KeyFingerprint = "hello5"
	bo1.Level = 1
	bo2.KeyFingerprint = "hello6"
	bo2.Level = 1
	bo3.KeyFingerprint = "hello6"
	bo3.Level = 1
	b.BoardOwners = []api.BoardOwner{bo1, bo2, bo3}
	err := persistence.BatchInsert([]interface{}{b})
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	}
	resp, err3 := persistence.ReadBoards(
		[]api.Fingerprint{api.Fingerprint(fp)}, 0, 0)
	if err3 != nil {
		t.Errorf("Test failed, err: '%s'", err3)
	} else if len(resp[0].BoardOwners) > 2 {
		t.Errorf("This should have returned 2 board owners. Error: '%#v\n' Current Board: '%#v\n', Board owners: '%#v\n'", err3, resp[0], resp[0].BoardOwners)
	}
}

func TestInsert_AddRemoveCurrencyAddress(t *testing.T) {
	var k api.Key
	fp := api.Fingerprint("hello this is a key")
	k.Fingerprint = fp
	var ca1 api.CurrencyAddress
	var ca2 api.CurrencyAddress
	var ca3 api.CurrencyAddress
	var ca4 api.CurrencyAddress
	ca1.Address = "hello"
	ca1.CurrencyCode = "XXX"
	ca2.Address = "hello2"
	ca2.CurrencyCode = "XXX"
	ca3.Address = "hello3"
	ca3.CurrencyCode = "XXX"
	ca4.Address = "hello4"
	ca4.CurrencyCode = "XXX"
	k.CurrencyAddresses = []api.CurrencyAddress{ca1, ca2, ca3}
	k.Key = "public key"
	k.Creation = 1
	k.ProofOfWork = "pow"
	k.Signature = "sig"
	k.Type = "key type"
	err := persistence.BatchInsert([]interface{}{k})
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	}
	k.CurrencyAddresses = []api.CurrencyAddress{ca1, ca2, ca4}
	k.LastUpdate = 1 // Remember, we need to do this otherwise it won't go in.
	err2 := persistence.BatchInsert([]interface{}{k})
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	}
	resp, err3 := persistence.ReadKeys([]api.Fingerprint{fp}, 0, 0)
	// fmt.Printf("%#v\n", resp[0].CurrencyAddresses)
	if err3 != nil {
		t.Errorf("Test failed, err: '%s'", err3)
	} else if len(resp[0].CurrencyAddresses) > 3 {
		t.Errorf("This should have returned 3 currency addresses. Error: '%#v\n' Current Key: '%#v\n', Currency addresses: '%#v\n'", err3, resp[0], resp[0].CurrencyAddresses)
	}
}

func TestInsert_DuplicateCurrencyAddress(t *testing.T) {
	var k api.Key
	fp := api.Fingerprint("hello this is a key")
	k.Fingerprint = fp
	var ca1 api.CurrencyAddress
	var ca2 api.CurrencyAddress
	var ca3 api.CurrencyAddress
	var ca4 api.CurrencyAddress
	ca1.Address = "hello"
	ca1.CurrencyCode = "XXX"
	ca2.Address = "hello2"
	ca2.CurrencyCode = "XXX"
	ca3.Address = "hello2"
	ca3.CurrencyCode = "XXX"
	ca4.Address = "hello4"
	ca4.CurrencyCode = "XXX"
	k.CurrencyAddresses = []api.CurrencyAddress{ca1, ca2, ca3, ca4}
	k.Key = "public key"
	k.Creation = 1
	k.ProofOfWork = "pow"
	k.Signature = "sig"
	k.Type = "key type"
	err := persistence.BatchInsert([]interface{}{k})
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	}
	resp, err3 := persistence.ReadKeys([]api.Fingerprint{fp}, 0, 0)
	// fmt.Printf("%#v\n", resp[0].CurrencyAddresses)
	if err3 != nil {
		t.Errorf("Test failed, err: '%s'", err3)
	} else if len(resp[0].CurrencyAddresses) > 3 {
		t.Errorf("This should have returned 3 currency addresses. Error: '%#v\n' Current Key: '%#v\n', Currency addresses: '%#v\n'", err3, resp[0], resp[0].CurrencyAddresses)
	}
}
