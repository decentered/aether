// API > Create / Verify
// This file provides the creation and verification commands that we use for each entity.

package api

import (
	// "fmt"
	// "aether-core/services/fingerprinting"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"aether-core/services/signaturing"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/ed25519"
	// "github.com/davecgh/go-spew/spew"
)

// // Create ProofOfWork
func (b *Board) CreatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if b.GetVersion() == 1 {
		return createBoardPoW_V1(b, keyPair, difficulty)
	} else {
		return errors.New(fmt.Sprintf("PoW creation of this version of this entity is not supported in this version of the app. Entity: %#v", b))
	}
}

func (t *Thread) CreatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if t.GetVersion() == 1 {
		return createThreadPoW_V1(t, keyPair, difficulty)
	} else {
		return errors.New(fmt.Sprintf("PoW creation of this version of this entity is not supported in this version of the app. Entity: %#v", t))
	}
}

func (p *Post) CreatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if p.GetVersion() == 1 {
		return createPostPoW_V1(p, keyPair, difficulty)
	} else {
		return errors.New(fmt.Sprintf("PoW creation of this version of this entity is not supported in this version of the app. Entity: %#v", p))
	}
}

func (v *Vote) CreatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if v.GetVersion() == 1 {
		return createVotePoW_V1(v, keyPair, difficulty)
	} else {
		return errors.New(fmt.Sprintf("PoW creation of this version of this entity is not supported in this version of the app. Entity: %#v", v))
	}
}

func (k *Key) CreatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if k.GetVersion() == 1 {
		return createKeyPoW_V1(k, keyPair, difficulty)
	} else {
		return errors.New(fmt.Sprintf("PoW creation of this version of this entity is not supported in this version of the app. Entity: %#v", k))
	}
}

func (ts *Truststate) CreatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if ts.GetVersion() == 1 {
		return createTruststatePoW_V1(ts, keyPair, difficulty)
	} else {
		return errors.New(fmt.Sprintf("PoW creation of this version of this entity is not supported in this version of the app. Entity: %#v", ts))
	}
}

// Create UpdateProofOfWork

func (b *Board) CreateUpdatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if b.GetVersion() == 1 {
		return createBoardUpdatePoW_V1(b, keyPair, difficulty)
	} else {
		return errors.New(fmt.Sprintf("PoW update creation of this version of this entity is not supported in this version of the app. Entity: %#v", b))
	}
}

func (t *Thread) CreateUpdatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if t.GetVersion() == 1 {
		return createThreadUpdatePoW_V1(t, keyPair, difficulty)
	} else {
		return errors.New(fmt.Sprintf("PoW update creation of this version of this entity is not supported in this version of the app. Entity: %#v", t))
	}
}

func (p *Post) CreateUpdatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if p.GetVersion() == 1 {
		return createPostUpdatePoW_V1(p, keyPair, difficulty)
	} else {
		return errors.New(fmt.Sprintf("PoW update creation of this version of this entity is not supported in this version of the app. Entity: %#v", p))
	}
}

func (v *Vote) CreateUpdatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if v.GetVersion() == 1 {
		return createVoteUpdatePoW_V1(v, keyPair, difficulty)
	} else {
		return errors.New(fmt.Sprintf("PoW update creation of this version of this entity is not supported in this version of the app. Entity: %#v", v))
	}
}

func (k *Key) CreateUpdatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if k.GetVersion() == 1 {
		return createKeyUpdatePoW_V1(k, keyPair, difficulty)
	} else {
		return errors.New(fmt.Sprintf("PoW update creation of this version of this entity is not supported in this version of the app. Entity: %#v", k))
	}
}

func (ts *Truststate) CreateUpdatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if ts.GetVersion() == 1 {
		return createTruststateUpdatePoW_V1(ts, keyPair, difficulty)
	} else {
		return errors.New(fmt.Sprintf("PoW update creation of this version of this entity is not supported in this version of the app. Entity: %#v", ts))
	}
}

// Verify ProofOfWork

func (b *Board) VerifyPoW(pubKey string) (bool, error) {
	if !globals.BackendTransientConfig.ProofOfWorkCheckEnabled {
		return true, nil
	}
	if b.GetVersion() == 1 {
		return verifyBoardPoW_V1(b, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("PoW verification of this version of this entity is not supported in this version of the app. Entity: %#v", b))
		return false, nil
	}
}

func (t *Thread) VerifyPoW(pubKey string) (bool, error) {
	if !globals.BackendTransientConfig.ProofOfWorkCheckEnabled {
		return true, nil
	}
	if t.GetVersion() == 1 {
		return verifyThreadPoW_V1(t, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("PoW verification of this version of this entity is not supported in this version of the app. Entity: %#v", t))
		return false, nil
	}
}

func (p *Post) VerifyPoW(pubKey string) (bool, error) {
	if !globals.BackendTransientConfig.ProofOfWorkCheckEnabled {
		return true, nil
	}
	if p.GetVersion() == 1 {
		return verifyPostPoW_V1(p, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("PoW verification of this version of this entity is not supported in this version of the app. Entity: %#v", p))
		return false, nil
	}
}

func (v *Vote) VerifyPoW(pubKey string) (bool, error) {
	if !globals.BackendTransientConfig.ProofOfWorkCheckEnabled {
		return true, nil
	}
	if v.GetVersion() == 1 {
		return verifyVotePoW_V1(v, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("PoW verification of this version of this entity is not supported in this version of the app. Entity: %#v", v))
		return false, nil
	}
}

func (k *Key) VerifyPoW(pubKey string) (bool, error) {
	if !globals.BackendTransientConfig.ProofOfWorkCheckEnabled {
		return true, nil
	}
	if k.GetVersion() == 1 {
		return verifyKeyPoW_V1(k, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("PoW verification of this version of this entity is not supported in this version of the app. Entity: %#v", k))
		return false, nil
	}
}

func (ts *Truststate) VerifyPoW(pubKey string) (bool, error) {
	if !globals.BackendTransientConfig.ProofOfWorkCheckEnabled {
		return true, nil
	}
	if ts.GetVersion() == 1 {
		return verifyTruststatePoW_V1(ts, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("PoW verification of this version of this entity is not supported in this version of the app. Entity: %#v", ts))
		return false, nil
	}
}

// Create Fingerprint

func (b *Board) CreateFingerprint() error {
	if b.GetVersion() == 1 {
		createBoardFp_V1(b)
		return nil
	} else {
		return errors.New(fmt.Sprintf("Fingerprint creation of this version of this entity is not supported in this version of the app. Entity: %#v", b))
	}
}

func (t *Thread) CreateFingerprint() error {
	if t.GetVersion() == 1 {
		createThreadFp_V1(t)
		return nil
	} else {
		return errors.New(fmt.Sprintf("Fingerprint creation of this version of this entity is not supported in this version of the app. Entity: %#v", t))
	}
}

func (p *Post) CreateFingerprint() error {
	if p.GetVersion() == 1 {
		createPostFp_V1(p)
		return nil
	} else {
		return errors.New(fmt.Sprintf("Fingerprint creation of this version of this entity is not supported in this version of the app. Entity: %#v", p))
	}
}

func (v *Vote) CreateFingerprint() error {
	if v.GetVersion() == 1 {
		createVoteFp_V1(v)
		return nil
	} else {
		return errors.New(fmt.Sprintf("Fingerprint creation of this version of this entity is not supported in this version of the app. Entity: %#v", v))
	}
}

func (k *Key) CreateFingerprint() error {
	if k.GetVersion() == 1 {
		createKeyFp_V1(k)
		return nil
	} else {
		return errors.New(fmt.Sprintf("Fingerprint creation of this version of this entity is not supported in this version of the app. Entity: %#v", k))
	}
}

func (ts *Truststate) CreateFingerprint() error {
	if ts.GetVersion() == 1 {
		createTruststateFp_V1(ts)
		return nil
	} else {
		return errors.New(fmt.Sprintf("Fingerprint creation of this version of this entity is not supported in this version of the app. Entity: %#v", ts))
	}
}

// Verify Fingerprint
func (b *Board) VerifyFingerprint() bool {
	if !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	if b.GetVersion() == 1 {
		return verifyBoardFingerprint_V1(b)
	} else {
		logging.Log(1, fmt.Sprintf("Fingerprint verification of this version of this entity is not supported in this version of the app. Entity: %#v", b))
		return false
	}
}

func (t *Thread) VerifyFingerprint() bool {
	if !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	if t.GetVersion() == 1 {
		return verifyThreadFingerprint_V1(t)
	} else {
		logging.Log(1, fmt.Sprintf("Fingerprint verification of this version of this entity is not supported in this version of the app. Entity: %#v", t))
		return false
	}
}

func (p *Post) VerifyFingerprint() bool {
	if !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	if p.GetVersion() == 1 {
		return verifyPostFingerprint_V1(p)
	} else {
		logging.Log(1, fmt.Sprintf("Fingerprint verification of this version of this entity is not supported in this version of the app. Entity: %#v", p))
		return false
	}
}

func (v *Vote) VerifyFingerprint() bool {
	if !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	if v.GetVersion() == 1 {
		return verifyVoteFingerprint_V1(v)
	} else {
		logging.Log(1, fmt.Sprintf("Fingerprint verification of this version of this entity is not supported in this version of the app. Entity: %#v", v))
		return false
	}
}

func (k *Key) VerifyFingerprint() bool {
	if !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	if k.GetVersion() == 1 {
		return verifyKeyFingerprint_V1(k)
	} else {
		logging.Log(1, fmt.Sprintf("Fingerprint verification of this version of this entity is not supported in this version of the app. Entity: %#v", k))
		return false
	}
}

func (ts *Truststate) VerifyFingerprint() bool {
	if !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	if ts.GetVersion() == 1 {
		return verifyTruststateFingerprint_V1(ts)
	} else {
		logging.Log(1, fmt.Sprintf("Fingerprint verification of this version of this entity is not supported in this version of the app. Entity: %#v", ts))
		return false
	}
}

// Signature

func (b *Board) CreateSignature(keyPair *ed25519.PrivateKey) error {
	if b.GetVersion() == 1 {
		return createBoardSignature_V1(b, keyPair)
	} else {
		return errors.New(fmt.Sprintf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", b))
	}
}

func (t *Thread) CreateSignature(keyPair *ed25519.PrivateKey) error {
	if t.GetVersion() == 1 {
		return createThreadSignature_V1(t, keyPair)
	} else {
		return errors.New(fmt.Sprintf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", t))
	}
}

func (p *Post) CreateSignature(keyPair *ed25519.PrivateKey) error {
	if p.GetVersion() == 1 {
		return createPostSignature_V1(p, keyPair)
	} else {
		return errors.New(fmt.Sprintf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", p))
	}
}

func (v *Vote) CreateSignature(keyPair *ed25519.PrivateKey) error {
	if v.GetVersion() == 1 {
		return createVoteSignature_V1(v, keyPair)
	} else {
		return errors.New(fmt.Sprintf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", v))
	}
}

func (k *Key) CreateSignature(keyPair *ed25519.PrivateKey) error {
	if k.GetVersion() == 1 {
		return createKeySignature_V1(k, keyPair)
	} else {
		return errors.New(fmt.Sprintf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", k))
	}
}

func (ts *Truststate) CreateSignature(keyPair *ed25519.PrivateKey) error {
	if ts.GetVersion() == 1 {
		return createTruststateSignature_V1(ts, keyPair)
	} else {
		return errors.New(fmt.Sprintf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", ts))
	}
}

// Create UpdateSignature

func (b *Board) CreateUpdateSignature(keyPair *ed25519.PrivateKey) error {
	if b.GetVersion() == 1 {
		return createBoardUpdateSignature_V1(b, keyPair)
	} else {
		return errors.New(fmt.Sprintf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", b))
	}
}

func (t *Thread) CreateUpdateSignature(keyPair *ed25519.PrivateKey) error {
	if t.GetVersion() == 1 {
		return createThreadUpdateSignature_V1(t, keyPair)
	} else {
		return errors.New(fmt.Sprintf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", t))
	}
}

func (p *Post) CreateUpdateSignature(keyPair *ed25519.PrivateKey) error {
	if p.GetVersion() == 1 {
		return createPostUpdateSignature_V1(p, keyPair)
	} else {
		return errors.New(fmt.Sprintf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", p))
	}
}

func (v *Vote) CreateUpdateSignature(keyPair *ed25519.PrivateKey) error {
	if v.GetVersion() == 1 {
		return createVoteUpdateSignature_V1(v, keyPair)
	} else {
		return errors.New(fmt.Sprintf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", v))
	}
}

func (k *Key) CreateUpdateSignature(keyPair *ed25519.PrivateKey) error {
	if k.GetVersion() == 1 {
		return createKeyUpdateSignature_V1(k, keyPair)
	} else {
		return errors.New(fmt.Sprintf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", k))
	}
}

func (ts *Truststate) CreateUpdateSignature(keyPair *ed25519.PrivateKey) error {
	if ts.GetVersion() == 1 {
		return createTruststateUpdateSignature_V1(ts, keyPair)
	} else {
		return errors.New(fmt.Sprintf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", ts))
	}
}

// Verify Signature

func (b *Board) VerifySignature(pubKey string) (bool, error) {
	if !globals.BackendTransientConfig.SignatureCheckEnabled {
		// If signature check is disabled with a debug flag, then we unconditionally return true.
		return true, nil
	}
	if globals.BackendConfig.GetAllowUnsignedEntities() && len(b.Signature) == 0 {
		// If Allow Unsigned Entities is true, we allow for anonymous posts without signature, but if there is a signature present, we still want to do the signature check. Allow Unsigned Entities does not mean that we will allow invalid signatures.
		return true, nil
	}
	if b.GetVersion() == 1 {
		return verifyBoardSignature_V1(b, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("Signature verification of this version of this entity is not supported in this version of the app. Entity: %#v", b))
		return false, nil
	}
}

func (t *Thread) VerifySignature(pubKey string) (bool, error) {
	if !globals.BackendTransientConfig.SignatureCheckEnabled {
		// If signature check is disabled with a debug flag, then we unconditionally return true.
		return true, nil
	}
	if globals.BackendConfig.GetAllowUnsignedEntities() && len(t.Signature) == 0 {
		// If Allow Unsigned Entities is true, we allow for anonymous posts without signature, but if there is a signature present, we still want to do the signature check. Allow Unsigned Entities does not mean that we will allow invalid signatures.
		return true, nil
	}
	if t.GetVersion() == 1 {
		return verifyThreadSignature_V1(t, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("Signature verification of this version of this entity is not supported in this version of the app. Entity: %#v", t))
		return false, nil
	}
}

func (p *Post) VerifySignature(pubKey string) (bool, error) {
	if !globals.BackendTransientConfig.SignatureCheckEnabled {
		// If signature check is disabled with a debug flag, then we unconditionally return true.
		return true, nil
	}
	if globals.BackendConfig.GetAllowUnsignedEntities() && len(p.Signature) == 0 {
		// If Allow Unsigned Entities is true, we allow for anonymous posts without signature, but if there is a signature present, we still want to do the signature check. Allow Unsigned Entities does not mean that we will allow invalid signatures.
		return true, nil
	}
	if p.GetVersion() == 1 {
		return verifyPostSignature_V1(p, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("Signature verification of this version of this entity is not supported in this version of the app. Entity: %#v", p))
		return false, nil
	}
}

func (v *Vote) VerifySignature(pubKey string) (bool, error) {
	if !globals.BackendTransientConfig.SignatureCheckEnabled {
		// If signature check is disabled with a debug flag, then we unconditionally return true.
		return true, nil
	}
	if globals.BackendConfig.GetAllowUnsignedEntities() && len(v.Signature) == 0 {
		// If Allow Unsigned Entities is true, we allow for anonymous posts without signature, but if there is a signature present, we still want to do the signature check. Allow Unsigned Entities does not mean that we will allow invalid signatures.
		return true, nil
	}
	if v.GetVersion() == 1 {
		return verifyVoteSignature_V1(v, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("Signature verification of this version of this entity is not supported in this version of the app. Entity: %#v", v))
		return false, nil
	}
}

func (k *Key) VerifySignature(pubKey string) (bool, error) {
	if !globals.BackendTransientConfig.SignatureCheckEnabled {
		// If signature check is disabled with a debug flag, then we unconditionally return true.
		return true, nil
	}
	if globals.BackendConfig.GetAllowUnsignedEntities() && len(k.Signature) == 0 {
		// If Allow Unsigned Entities is true, we allow for anonymous posts without signature, but if there is a signature present, we still want to do the signature check. Allow Unsigned Entities does not mean that we will allow invalid signatures.
		return true, nil
	}
	if k.GetVersion() == 1 {
		return verifyKeySignature_V1(k, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("Signature verification of this version of this entity is not supported in this version of the app. Entity: %#v", k))
		return false, nil
	}
}

func (ts *Truststate) VerifySignature(pubKey string) (bool, error) {
	if !globals.BackendTransientConfig.SignatureCheckEnabled {
		// If signature check is disabled with a debug flag, then we unconditionally return true.
		return true, nil
	}
	if globals.BackendConfig.GetAllowUnsignedEntities() && len(ts.Signature) == 0 {
		// If Allow Unsigned Entities is true, we allow for anonymous posts without signature, but if there is a signature present, we still want to do the signature check. Allow Unsigned Entities does not mean that we will allow invalid signatures.
		return true, nil
	}
	if ts.GetVersion() == 1 {
		return verifyTruststateSignature_V1(ts, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("Signature verification of this version of this entity is not supported in this version of the app. Entity: %#v", ts))
		return false, nil
	}
}

// Api Response Signature Create / Verify

func (ar *ApiResponse) CreateSignature(keyPair *ed25519.PrivateKey) error {
	// Unlike other signatures, ApiResponse signature includes the key that it is signed by itself, because it does not have a separate fingerprint field. By including the key within the signature, we protect the key under the seal of the signature, as well.
	cpI := *ar
	// Remove signature just in case, if it's been accidentally set.
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	ar.Signature = Signature(signature)
	return nil
}

// VerifySignature verifies the signature of the page. Since the public key the page is verified by is within the page itself, it does not need the public key to be given from the outside.
func (ar *ApiResponse) VerifySignature() (bool, error) {
	// 1) Check if signature check is enabled.
	if !globals.BackendTransientConfig.PageSignatureCheckEnabled {
		return true, nil
	}
	// 2) Check if required fields are empty.
	if !(len(ar.NodePublicKey) > 0 && len(ar.Signature) > 0) {
		return false, errors.New(fmt.Sprintf(
			"Page signature check is enabled, but the page has some fields (Public Key or Signature) empty. Public Key: %s, Signature: %s", ar.NodePublicKey, ar.Signature))
	}
	// 3) Verify signature.
	cpI := *ar
	var signature string
	// Determine if we are checking for original or update signature
	// Save signature to be verified
	signature = string(cpI.Signature)
	// This happens *after* Signature, so should be empty here.
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Signature
	verifyResult := signaturing.Verify(string(res), signature, ar.NodePublicKey)
	// If the Signature is valid
	if verifyResult {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprintf(
			"This signature is invalid, but no reason given as to why. Signature: %s", signature))
	}
}

// Verification for the provable and for the response.
func Verify(e interface{}) error {
	switch entity := e.(type) {
	case Provable:
		encrypted := len(entity.GetEncrContent()) > 0
		if encrypted {
			return errors.New(fmt.Sprintf("This item appears to be encrypted. Please decrypt before requesting verification. EncrContent: %s, Entity: %#v", entity.GetEncrContent(), entity))
		}
		realmed := len(entity.GetRealmId()) > 0
		if realmed {
			return errors.New(fmt.Sprintf("This item appears to belong to a realm that is different than the mainnet. Non-mainnet realms are currently not supported, but might be in the future. RealmId: %s, Entity: %#v", entity.GetRealmId(), entity))
		}
		boundsOk, err := entity.CheckBounds()
		if err != nil {
			return err
		}
		if !boundsOk {
			return errors.New(fmt.Sprintf("Field boundaries of this entity is invalid. Entity: %#v", entity))
		}
		fpOk := entity.VerifyFingerprint()
		if !fpOk {
			return errors.New(fmt.Sprintf(
				"Fingerprint of this entity is invalid. Fingerprint: %s, Entity: %#v\n", entity.GetFingerprint(), entity))
		}
		// Bounds ok, Fp ok
		powOk, err2 := entity.VerifyPoW(entity.GetOwnerPublicKey())
		if err2 != nil {
			return err2
		}
		if !powOk {
			return errors.New(fmt.Sprintf(
				"ProofOfWork of this entity is invalid. ProofOfWork: %s, Entity: %#v\n", entity.GetProofOfWork(), entity))
		}
		// Bounds ok, Fp ok, PoW ok
		sigOk, err3 := entity.VerifySignature(entity.GetOwnerPublicKey())
		if err3 != nil {
			return err3
		}
		if !sigOk {
			return errors.New(fmt.Sprintf(
				"Signature of this entity is invalid. Signature: %s, Entity: %#v\n", entity.GetSignature(), entity))
		}
		// Bounds ok, Fp ok, PoW ok, Sig ok
		entity.SetVerified(true)
		return nil

	case *Address:
		boundsOk, err := entity.CheckBounds()
		if err != nil {
			return err
		}
		if !boundsOk {
			return errors.New(fmt.Sprintf("Field boundaries of this entity is invalid. Entity: %#v", entity))
		}
		// Bounds ok
		entity.SetVerified(true)
		return nil

	default:
		return errors.New(fmt.Sprintf("Verify could not recognise this entity type. Entity: %#v", entity))
	}

}
