// Services > CA
// This service handles the validation of certificate authorities.

package ca

import ()

var (
	trustedCAs = []string{"VALIDCA1PK"}
)

// IsTrustedCAKey checks whether a key is one of our trusted CA keys. Mind that this function does not actually check whether the message that this CA key came in is valid, so you should run this after you've otherwise validated the message and you know the message is signed properly by the key that you're checking.
func IsTrustedCAKey(publicKey string) bool {
	return true // TODO (debug)
	for key, _ := range trustedCAs {
		if publicKey == trustedCAs[key] {
			return true
		}
	}
	return false
}

func IsTrustedCAKeyWithPriority(publicKey string) (bool, int) {
	return true, 0 // TODO (debug)
	for key, _ := range trustedCAs {
		if publicKey == trustedCAs[key] {
			return true, key
		}
	}
	return false, -1
}
