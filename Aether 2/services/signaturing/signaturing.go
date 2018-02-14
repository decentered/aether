// Services > Fingerprinting
// This module handles the creation of signatures, signing entities, and checking for signatures.

package signaturing

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

/*
This is where we create keys, sign entities, and validate them. This does have a dependency on the config system, in that the user's private key is needs to be saved to the configuration, which is a subsystem of its own.
*/

func CreateKeyPair() (*ecdsa.PrivateKey, error) {
	// This emits a key pair entity. The private key emitted here also includes the public key, too.
	curve := elliptic.P521()
	privKey := new(ecdsa.PrivateKey) // Also includes public key in itself.
	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return new(ecdsa.PrivateKey), errors.New(fmt.Sprint(
			"Key pair generation failed. err: ", err))
	}
	// Note: accessing the public key is shown below.
	// var pubKey ecdsa.PublicKey
	// pubKey = privKey.PublicKey
	return privKey, nil
}

func Sign(input string, privKey *ecdsa.PrivateKey) (string, error) {
	// This creates an ECDSA signature for the input provided.
	// Mind that signatures are generated on hashes of items, not the item itself. So the first step in either signing or verifying is to hash the input provided.
	inputByte := []byte(input)
	hasher := sha256.New()
	hasher.Write(inputByte)
	hash := hasher.Sum(nil)
	r := big.NewInt(0)
	s := big.NewInt(0)
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash)
	if err != nil {
		return "", errors.New(fmt.Sprint(
			"Signing failed. err: ", err))
	}
	// Why "..."? Because append takes only one byte and s.Bytes is a []byte, not byte.
	signature := fmt.Sprint(
		hex.EncodeToString(r.Bytes()),
		"-", hex.EncodeToString(s.Bytes()))
	return signature, nil
}

func Verify(input string, signature string, pubKey string) bool {
	// This verifies the input provided by a given signature and public key.
	// Mind that signatures are generated on hashes of items, not the item itself. So the first step in either signing or verifying is to hash the input provided.
	// First of all, let's make sure that this is not an anonymous post. If anonymous, the signature and pubKey will be empty, and this should pass.
	if signature == "" && pubKey == "" {
		return true
	} else if signature == "" && pubKey != "" {
		return false
	} else if signature != "" && pubKey == "" {
		return false
	}
	inputByte := []byte(input)
	hasher := sha256.New()
	hasher.Write(inputByte)
	hash := hasher.Sum(nil)
	// Signature = hex(byte(r))-hex(byte(s)), - being dash.
	// Find dash to split numbers.
	dashCaret := strings.Index(signature, "-")
	if dashCaret == -1 {
		return false
	}
	r := big.NewInt(0)
	s := big.NewInt(0)
	rAsHex := signature[0:dashCaret]
	sAsHex := signature[dashCaret+1:]
	rAsByte, err := hex.DecodeString(rAsHex)
	sAsByte, err2 := hex.DecodeString(sAsHex)
	if err != nil || err2 != nil {
		return false
	}
	r.SetBytes(rAsByte)
	s.SetBytes(sAsByte)
	// Unmarshal the public key
	// x, y := elliptic.Unmarshal(elliptic.P521(), []byte(pubKey))
	pubKeyAsByte, err := hex.DecodeString(pubKey)
	if err != nil {
		return false
	}
	x, y := elliptic.Unmarshal(elliptic.P521(), pubKeyAsByte)
	// Create the public key struct
	if x == nil || y == nil {
		return false
	}
	publicKeyStruct := new(ecdsa.PublicKey)
	publicKeyStruct.Curve = elliptic.P521()
	publicKeyStruct.X = x
	publicKeyStruct.Y = y
	result := ecdsa.Verify(publicKeyStruct, hash, r, s)
	return result
}
