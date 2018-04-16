package api_test

import (
	"aether-core/io/api"
	"aether-core/services/configstore"
	"aether-core/services/create"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"aether-core/services/signaturing"
	"aether-core/services/verify"
	"crypto/elliptic"
	"encoding/hex"
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

var MarshaledPubKey string

func startconfigs() {
	becfg, err := configstore.EstablishBackendConfig()
	if err != nil {
		logging.LogCrash(err)
	}
	becfg.Cycle()
	globals.BackendConfig = becfg

	fecfg, err := configstore.EstablishFrontendConfig()
	if err != nil {
		logging.LogCrash(err)
	}
	fecfg.Cycle()
	globals.FrontendConfig = fecfg
}

func setup() {
	startconfigs()
	globals.BackendConfig.SetMinimumPoWStrengths(16)
	MarshaledPubKey = hex.EncodeToString(elliptic.Marshal(elliptic.P521(), globals.FrontendConfig.GetUserKeyPair().PublicKey.X, globals.FrontendConfig.GetUserKeyPair().PublicKey.Y))
}

func teardown() {
}

// Tests

func TestVerify_Success(t *testing.T) {
	keyEntity, err3 := create.CreateKey(
		"", MarshaledPubKey, "", "")
	if err3 != nil {
		t.Errorf("Object creation failed. Err: '%s'", err3)
	}
	thr, err :=
		create.CreateThread(
			"my board fingerprint",
			"my thread name",
			"my thread body",
			"my thread link",
			keyEntity.Fingerprint)
	if err != nil {
		t.Errorf("Object creation failed. Err: '%#v\n'", err)
	}
	result, err2 := verify.Verify(&thr, keyEntity)
	if err2 != nil {
		t.Errorf("Object verification process failed. Error: '%#v\n'", err2)
	}
	if result != true {
		t.Errorf("This object should be valid, but it is invalid. Object: '%#v\n'", thr)
	}
}

func TestVerify_BrokenFingerprint_Fail(t *testing.T) {
	keyEntity, err3 := create.CreateKey(
		"", MarshaledPubKey, "", "")
	if err3 != nil {
		t.Errorf("Object creation failed. Err: '%s'", err3)
	}
	thr, err :=
		create.CreateThread(
			"my board fingerprint",
			"my thread name",
			"my thread body",
			"my thread link",
			"my owner fingerprint")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	thr.Board = "I'm changing the board of this thread to fail the test."
	result, err2 := verify.Verify(&thr, keyEntity)
	errMessage := "Fingerprint of this entity is invalid"
	if err2 == nil || result == true {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
}

func TestVerify_BrokenPoW1_Fail(t *testing.T) {
	keyEntity, err3 := create.CreateKey(
		"", MarshaledPubKey, "", "")
	if err3 != nil {
		t.Errorf("Object creation failed. Err: '%s'", err3)
	}
	thr, err :=
		create.CreateThread(
			"my board fingerprint",
			"my thread name",
			"my thread body",
			"my thread link",
			"my owner fingerprint")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	thr.ProofOfWork = "I'm changing the board pow to fail the test."
	// Re-shrink-wrap
	thr.CreateFingerprint()
	errMessage := "PoW had more or less fields than expected"
	result, err2 := verify.Verify(&thr, keyEntity)
	if err2 == nil || result == true {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
	// fmt.Printf("%#v\n", thr)
	// fmt.Printf("%#v\n", result)
}

func TestVerify_BrokenPoW2_Fail(t *testing.T) {
	// Changing a mutable element, but not actually running update.
	keyEntity, err3 := create.CreateKey(
		"", MarshaledPubKey, "", "")
	if err3 != nil {
		t.Errorf("Object creation failed. Err: '%s'", err3)
	}
	entity, err :=
		create.CreateBoard(
			"board name",
			"my board owner fingerprint",
			*new([]api.BoardOwner),
			"my board description")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	entity.Description = "Changing an element that is not protected by fingerprint"
	result, err2 := verify.Verify(&entity, keyEntity)
	errMessage := "This proof of work is invalid or malformed"
	if err2 == nil || result == true {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
}

func TestVerify_BrokenSignature_Fail(t *testing.T) {
	keyEntity, err3 := create.CreateKey(
		"", MarshaledPubKey, "", "")
	if err3 != nil {
		t.Errorf("Object creation failed. Err: '%s'", err3)
	}
	privKey, err := signaturing.CreateKeyPair()
	if err != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err)
	}
	thr, err :=
		create.CreateThread(
			"my board fingerprint",
			"my thread name",
			"my thread body",
			"my thread link",
			"my owner fingerprint")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	// Re-shrink-wrap
	thr.CreateSignature(privKey) // Signing it with a new key
	thr.CreatePoW(globals.FrontendConfig.GetUserKeyPair(), 20)
	thr.CreateFingerprint()
	errMessage := "A wrong key is provided for this signature"
	result, err2 := verify.Verify(&thr, keyEntity)
	if err2 == nil || result == true {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
	// fmt.Printf("%#v\n", thr)
	// fmt.Printf("%#v\n", result)
}

// Test update verify success / fail

func TestVerify_UpdatedItemSuccess(t *testing.T) {
	keyEntity, err3 := create.CreateKey(
		"", MarshaledPubKey, "", "")
	if err3 != nil {
		t.Errorf("Object creation failed. Err: '%s'", err3)
	}
	board, err :=
		create.CreateBoard(
			"my board name",
			keyEntity.Fingerprint,
			[]api.BoardOwner{},
			"my board description")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%#v\n'", err)
	}
	updatereq := create.BoardUpdateRequest{}
	updatereq.Entity = &board
	updatereq.DescriptionUpdated = true
	updatereq.NewDescription = "I changed the board description!"
	create.UpdateBoard(updatereq)
	result, err2 := verify.Verify(&board, keyEntity)
	if err2 != nil {
		t.Errorf("Object verification process failed. Error: '%#v\n'", err2)
	}
	if result != true {
		t.Errorf("This object should be valid, but it is invalid. Object: '%#v\n'", board)
	}
}

func TestVerify_UpdatedItemFailure_Pow(t *testing.T) {
	// Failed to call the update request.
	keyEntity, err3 := create.CreateKey(
		"", MarshaledPubKey, "", "")
	if err3 != nil {
		t.Errorf("Object creation failed. Err: '%s'", err3)
	}
	board, err :=
		create.CreateBoard(
			"my board name",
			keyEntity.Fingerprint,
			[]api.BoardOwner{},
			"my board description")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%#v\n'", err)
	}
	// Since this is a mutable field, the fingerprint does not protect against it. The field to fail first will be pow.
	board.Description = "new description"
	result, err2 := verify.Verify(&board, keyEntity)
	errMessage := "This proof of work is invalid or malformed"
	if err2 == nil || result == true {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
}

func TestVerify_UpdatedItemFailure_Fingerprint(t *testing.T) {
	// Failed to call the update request.
	keyEntity, err3 := create.CreateKey(
		"", MarshaledPubKey, "", "")
	if err3 != nil {
		t.Errorf("Object creation failed. Err: '%s'", err3)
	}
	board, err :=
		create.CreateBoard(
			"my board name",
			keyEntity.Fingerprint,
			[]api.BoardOwner{},
			"my board description")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%#v\n'", err)
	}
	// Since this is an immutable field, it will trip up the fingerprint first.
	board.Name = "New name"
	result, err2 := verify.Verify(&board, keyEntity)
	errMessage := "Fingerprint of this entity is invalid"
	if err2 == nil || result == true {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
}

func TestVerify_UpdatedItemFailure_Signature(t *testing.T) {
	// Failed to call the update request.
	keyEntity, err3 := create.CreateKey(
		"", MarshaledPubKey, "", "")
	if err3 != nil {
		t.Errorf("Object creation failed. Err: '%s'", err3)
	}
	board, err :=
		create.CreateBoard(
			"my board name",
			keyEntity.Fingerprint,
			[]api.BoardOwner{},
			"my board description")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%#v\n'", err)
	}
	// Since this is an immutable field, it will trip up the fingerprint first.
	keyEntity.Fingerprint = "changed key entity fingerprint"
	result, err2 := verify.Verify(&board, keyEntity)
	errMessage := "A wrong key is provided for this signature"
	if err2 == nil || result == true {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
}
