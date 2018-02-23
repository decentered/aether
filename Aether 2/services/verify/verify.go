// Verify
// This package verifies all entities. It uses the other lower level services in the services directory and it provides a single endpoint, comprehensive solution to deal with the verification.

package verify

import (
	"aether-core/io/api"
	"aether-core/io/persistence"
	"aether-core/services/logging"
	"errors"
	"fmt"
)

// findKey finds the public key of any given key fingerprint by looking first into the response that is received, and then into the database.
func findKey(
	fpToLookFor api.Fingerprint,
	resp api.Response) (api.Key, error) {
	// findKey looks for the appropriate key first in the given API response, then in the database. When found, it verifies the key and returns. If not found, the appropriate error is returned.
	var result api.Key
	var unverifiedKey api.Key

	if fpToLookFor == "" {
		// Before anything else, if the fpToLookFor is empty, directly return blank. It's an anonymous entity.
		return unverifiedKey, nil
	}
	// If this is not an anonymou entity, first search in the incoming response.
	var foundInResp bool
	for _, key := range resp.Keys {
		if key.Fingerprint == fpToLookFor {
			unverifiedKey = key
			foundInResp = true
			break
		}
	}
	// Second, if the key does not exist in the response, search in the local DB.
	if !foundInResp {
		// Not found in the response received from the remote
		dbResp, err := persistence.ReadKeys([]api.Fingerprint{fpToLookFor}, 0, 0)
		if err != nil {
			return result, err
		}
		if len(dbResp) > 0 {
			// We have a DB response with at least one item.
			if len(dbResp) < 2 {
				// We have a DB response with exactly one item.
				unverifiedKey = dbResp[0]
			} else {
				// We have a response with *multiple* items. Fail. This should never happen, it indicates DB corruption.
				return result, errors.New(fmt.Sprintf(
					"WARNING: A database search for a public key found multiple results. This should never happen - this indicates DB corruption. Key with multiple results: %s\n", fpToLookFor))
			}
		} else {
			// We have a DB response with zero items. Fail.
			return result, errors.New(fmt.Sprintf(
				"The key requested does not exist in the database or in the incoming response. Key: %s\n", fpToLookFor))
		}
	}
	// We have the unverified key. Now run verify over it.
	isVerified, err2 := Verify(&unverifiedKey, unverifiedKey)
	if err2 != nil {
		return result, err2
	}
	if isVerified {
		result = unverifiedKey
	} else {
		// If found in the response:
		if foundInResp {
			return result, errors.New(fmt.Sprintf(
				"The key requested is present, however it failed the verification. Key: %#v\n", unverifiedKey))
		} else {
			// If it came from the database
			return result, errors.New(fmt.Sprintf(
				"WARNING: A database search for a public key found a key entity that did not validate. This should never happen - this indicates DB corruption. Key with multiple results: %s\n", fpToLookFor))
		}
	}
	return result, nil
}

// verifyProvable verifies any given api.Provable. It automatically handles key finding.
func verifyProvable(resp api.Response, entity api.Provable) (bool, error) {
	// Find the key that is needed to validate the item.
	owner := entity.GetOwner()
	// If this is not an anonymous entity, look for the key.
	key, err := findKey(owner, resp)
	if err != nil {
		return false, errors.New(fmt.Sprintf(
			"An error occurred when the key for this entity was being searched for. Entity: %#v, Error: %s\n", entity, err))
	}
	// Then verify it with that key.
	isVerified, err2 := Verify(entity, key)
	if err2 != nil {
		// We have an error in validation process.
		return false, errors.New(fmt.Sprintf(
			"There occurred an error in the validation process. Entity: %#v, Error: %s\n", entity, err2))
	} else if !isVerified {
		// We do not have an error in validation but the item could not be validated..
		return false, errors.New(fmt.Sprintf(
			"An entity in the response could not be validated. Entity: %#v, Coming from: %#v\n", entity, resp.Addresses))
	}
	if isVerified {
		// We do not have an error, and the validation was successful.
		return true, nil
	}
	return false, nil
}

// VerifyResponse filters through an API Response and removes the broken / unverifiable items.
func VerifyResponse(resp api.Response) api.Response {
	var cleanedResp api.Response
	for _, entity := range resp.Boards {
		isVerified, err := verifyProvable(resp, &entity)
		if isVerified {
			cleanedResp.Boards = append(cleanedResp.Boards, entity)
		} else {
			logging.Log(1, fmt.Sprintf("Verification failed for this entity. Entity: %#v, Error: %s", entity, err))
		}
	}

	for _, entity := range resp.Threads {
		isVerified, err := verifyProvable(resp, &entity)
		if isVerified {
			cleanedResp.Threads = append(cleanedResp.Threads, entity)
		} else {
			logging.Log(1, fmt.Sprintf("Verification failed for this entity. Entity: %#v, Error: %s", entity, err))
		}
	}

	for _, entity := range resp.Posts {
		isVerified, err := verifyProvable(resp, &entity)
		if isVerified {
			cleanedResp.Posts = append(cleanedResp.Posts, entity)
		} else {
			logging.Log(1, fmt.Sprintf("Verification failed for this entity. Entity: %#v, Error: %s", entity, err))
		}
	}

	for _, entity := range resp.Votes {
		isVerified, err := verifyProvable(resp, &entity)
		if isVerified {
			cleanedResp.Votes = append(cleanedResp.Votes, entity)
		} else {
			logging.Log(1, fmt.Sprintf("Verification failed for this entity. Entity: %#v, Error: %s", entity, err))
		}
	}

	// Auto pass the addresses - those cannot be verified.
	cleanedResp.Addresses = resp.Addresses

	for _, entity := range resp.Keys {
		isVerified, err := verifyProvable(resp, &entity)
		if isVerified {
			cleanedResp.Keys = append(cleanedResp.Keys, entity)
		} else {
			logging.Log(1, fmt.Sprintf("Verification failed for this entity. Entity: %#v, Error: %s", entity, err))
		}
	}

	for _, entity := range resp.Truststates {
		isVerified, err := verifyProvable(resp, &entity)
		if isVerified {
			cleanedResp.Truststates = append(cleanedResp.Truststates, entity)
		} else {
			logging.Log(1, fmt.Sprintf("Verification failed for this entity. Entity: %#v, Error: %s", entity, err))
		}
	}
	return cleanedResp
}

func Verify(entity api.Provable, keyEntity api.Key) (bool, error) {
	pubKey := keyEntity.Key
	fpOk := entity.VerifyFingerprint()
	if !fpOk {
		return false, errors.New(fmt.Sprintf(
			"Fingerprint of this entity is invalid. Fingerprint: %s, Entity: %#v\n", entity.GetFingerprint(), entity))
	} else {
		// Fp ok
		powOk, err2 := entity.VerifyPoW(pubKey)
		if err2 != nil {
			return false, err2
		}
		if !powOk {
			return false, errors.New(fmt.Sprintf(
				"ProofOfWork of this entity is invalid. ProofOfWork: %s, Entity: %#v\n", entity.GetProofOfWork(), entity))
		} else {
			// Fp ok, PoW ok
			// Check that the fingerprint of the key matches owner fingerprint in the object.
			if entity.GetOwner() == keyEntity.GetFingerprint() {
				sigOk, err3 := entity.VerifySignature(pubKey)
				if err3 != nil {
					return false, err3
				}
				if !sigOk {
					return false, errors.New(fmt.Sprintf(
						"Signature of this entity is invalid. Signature: %s, Entity: %#v\n", entity.GetSignature(), entity))
				} else {
					// Fp ok, PoW ok, Sig ok
					return true, nil
				}
			} else {
				// Entity owner isn't the signature given to the method
				return false, errors.New(fmt.Sprintf(
					"A wrong key is provided for this signature. Entity Signature: %s, Provided Signature Fingerprint: %#v\n", entity.GetSignature(), keyEntity.Fingerprint))
			}
		}
	}
}
