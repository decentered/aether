package create_test

import (
	"aether-core/io/api"
	"aether-core/services/create"
	"aether-core/services/globals"
	// "aether-core/services/signaturing"
	"aether-core/services/verify"
	// "fmt"
	"os"
	"strings"
	"testing"
)

// Infrastructure, setup and teardown

func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	teardown()
	os.Exit(exitVal)
}

var UserKeyEntity api.Key

func setup() {
	globals.GenerateUserKeyPair()
	globals.SetBailoutTime()
	globals.SetMinPoWStrengths(16)
	UserKeyEntity, _ = create.CreateKey("", globals.MarshaledPubKey, "", *new([]api.CurrencyAddress), "")
}

func teardown() {
}

// Tests

// Tests for sub-entity creation

func TestCreateBoardOwner_Success(t *testing.T) {
	_, err :=
		create.CreateBoardOwner(
			"my key fingerprint",
			api.Timestamp(12345678),
			uint8(1))
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
}

func TestCreateCurrencyAddress_Success(t *testing.T) {
	_, err :=
		create.CreateCurrencyAddress(
			"XFK",
			"My currency address")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
}

// Tests for main entity creation.

func TestCreateBoard_Success(t *testing.T) {
	bo, _ :=
		create.CreateBoardOwner(UserKeyEntity.GetOwner(), api.Timestamp(12345678), uint8(1))
	entity, err :=
		create.CreateBoard(
			"My board name",
			bo.KeyFingerprint,
			[]api.BoardOwner{bo},
			"my description")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	result, err2 := verify.Verify(&entity, UserKeyEntity)
	if err2 != nil {
		t.Errorf("Object verification process failed. Err: '%s'", err2)
	}
	if result != true {
		t.Errorf("This object should be valid, but it is invalid. Entity: '%#v\n'", entity)
	}
}

func TestCreateThread_Success(t *testing.T) {
	entity, err :=
		create.CreateThread(
			"Thread parent (board) fingerprint",
			"thread name",
			"thread body",
			"thread link",
			UserKeyEntity.GetOwner())
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	result, err2 := verify.Verify(&entity, UserKeyEntity)
	if err2 != nil {
		t.Errorf("Object verification process failed. Err: '%s'", err2)
	}
	if result != true {
		t.Errorf("This object should be valid, but it is invalid. Entity: '%#v\n'", entity)
	}
}

func TestCreatePost_Success(t *testing.T) {
	entity, err :=
		create.CreatePost(
			"Post parent (board) fingerprint",
			"Post parent (thread) fingerprint",
			"Post parent (post or thread) fingerprint",
			"Post body",
			UserKeyEntity.GetOwner())
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	result, err2 := verify.Verify(&entity, UserKeyEntity)
	if err2 != nil {
		t.Errorf("Object verification process failed. Err: '%s'", err2)
	}
	if result != true {
		t.Errorf("This object should be valid, but it is invalid. Entity: '%#v\n'", entity)
	}
}

func TestCreateVote_Success(t *testing.T) {
	entity, err :=
		create.CreateVote(
			"board fp",
			"thread fp",
			"target fp",
			UserKeyEntity.GetOwner(),
			uint8(1))
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	result, err2 := verify.Verify(&entity, UserKeyEntity)
	if err2 != nil {
		t.Errorf("Object verification process failed. Err: '%s'", err2)
	}
	if result != true {
		t.Errorf("This object should be valid, but it is invalid. Entity: '%#v\n'", entity)
	}
}

func TestCreateAddress_Success(t *testing.T) {
	_, err :=
		create.CreateAddress(
			api.Location("127.0.0.1"),
			api.Location("test_subdir"),
			uint8(1),
			uint16(4732),
			uint8(1),
			api.Timestamp(12345678),
			uint8(1),
			uint16(1),
			[]api.Subprotocol{api.Subprotocol{"c0", 1, 0, []string{"board", "thread", "post", "vote", "key", "truststate"}}},
			uint8(1),
			uint16(1),
			uint16(0),
			"my client name",
		)
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
}

func TestCreateKey_Success(t *testing.T) {
	entity, err :=
		create.CreateKey(
			"key type",
			globals.MarshaledPubKey,
			"user name",
			*new([]api.CurrencyAddress),
			"key info")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	// fmt.Printf("%#v\n", entity)
	result, err2 := verify.Verify(&entity, entity)
	if err2 != nil {
		t.Errorf("Object verification process failed. Err: '%s'", err2)
	}
	if result != true {
		t.Errorf("This object should be valid, but it is invalid. Entity: '%#v\n'", entity)
	}
}

func TestCreateTruststate_Success(t *testing.T) {
	entity, err :=
		create.CreateTruststate(
			"target fp",
			UserKeyEntity.GetFingerprint(),
			uint8(1),
			[]api.Fingerprint{"domain1fp", "domain2fp"},
			api.Timestamp(12345678))
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	result, err2 := verify.Verify(&entity, UserKeyEntity)
	if err2 != nil {
		t.Errorf("Object verification process failed. Err: '%s'", err2)
	}
	if result != true {
		t.Errorf("This object should be valid, but it is invalid. Entity: '%#v\n'", entity)
	}
}

// Entity updates

func TestUpdateBoard_Success(t *testing.T) {
	bo, _ :=
		create.CreateBoardOwner(UserKeyEntity.GetFingerprint(), api.Timestamp(12345678), uint8(1))
	entity, err :=
		create.CreateBoard(
			"My board name",
			bo.KeyFingerprint,
			[]api.BoardOwner{bo},
			"my description")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	updatereq := create.BoardUpdateRequest{}
	updatereq.Entity = &entity
	updatereq.DescriptionUpdated = true
	updatereq.NewDescription = "I changed the board description!"
	create.UpdateBoard(updatereq)
	result, err2 := verify.Verify(&entity, UserKeyEntity)
	if err2 != nil {
		t.Errorf("Object verification process failed. Err: '%s'", err2)
	}
	if result != true {
		t.Errorf("This object should be valid, but it is invalid. Entity: '%#v\n'", entity)
	}
	// fmt.Printf("%#v\n", entity)
	// fmt.Printf("%#v\n", result)
}

func TestUpdateVote_Success(t *testing.T) {
	entity, err :=
		create.CreateVote(
			"board fp",
			"thread fp",
			"target fp",
			UserKeyEntity.GetFingerprint(),
			uint8(1))
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	updatereq := create.VoteUpdateRequest{}
	updatereq.Entity = &entity
	updatereq.TypeUpdated = true
	updatereq.NewType = 0
	create.UpdateVote(updatereq)
	result, err2 := verify.Verify(&entity, UserKeyEntity)
	if err2 != nil {
		t.Errorf("Object verification process failed. Err: '%s'", err2)
	}
	if result != true {
		t.Errorf("This object should be valid, but it is invalid. Entity: '%#v\n'", entity)
	}
	// fmt.Printf("%#v\n", entity)
	// fmt.Printf("%#v\n", result)
}

func TestUpdateKey_Success(t *testing.T) {
	entity, err :=
		create.CreateKey(
			"key type",
			globals.MarshaledPubKey,
			"user name",
			*new([]api.CurrencyAddress),
			"key info")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	updatereq := create.KeyUpdateRequest{}
	updatereq.Entity = &entity
	updatereq.InfoUpdated = true
	updatereq.NewInfo = "This is my new key info."
	create.UpdateKey(updatereq)
	result, err2 := verify.Verify(&entity, entity)
	if err2 != nil {
		t.Errorf("Object verification process failed. Err: '%s'", err2)
	}
	if result != true {
		t.Errorf("This object should be valid, but it is invalid. Entity: '%#v\n'", entity)
	}
}

func TestUpdateTruststate_Success(t *testing.T) {
	entity, err :=
		create.CreateTruststate(
			"target fp",
			UserKeyEntity.GetFingerprint(),
			uint8(1),
			[]api.Fingerprint{"domain1fp", "domain2fp"},
			api.Timestamp(12345678))
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	updatereq := create.TruststateUpdateRequest{}
	updatereq.Entity = &entity
	updatereq.TypeUpdated = true
	updatereq.NewType = 3
	create.UpdateTruststate(updatereq)
	result, err2 := verify.Verify(&entity, UserKeyEntity)
	if err2 != nil {
		t.Errorf("Object verification process failed. Err: '%s'", err2)
	}
	if result != true {
		t.Errorf("This object should be valid, but it is invalid. Entity: '%#v\n'", entity)
	}
	// fmt.Printf("%#v\n", entity)
	// fmt.Printf("%#v\n", result)
}

func TestUpdateVote_EditAfter_Fail(t *testing.T) {
	entity, err :=
		create.CreateVote(
			"board fp",
			"thread fp",
			"target fp",
			"owner fp",
			uint8(1))
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	updatereq := create.VoteUpdateRequest{}
	updatereq.Entity = &entity
	updatereq.TypeUpdated = true
	updatereq.NewType = 0
	create.UpdateVote(updatereq)
	entity.Type = 2
	errMessage := "This proof of work is invalid or malformed"
	result, err2 := verify.Verify(&entity, UserKeyEntity)
	if err2 == nil || result == true {
		t.Errorf("Expected an error to be raised from this test.")
	}
	if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
	// fmt.Printf("%#v\n", entity)
	// fmt.Printf("%#v\n", result)
}
