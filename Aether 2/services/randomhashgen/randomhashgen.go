// Services > RandomHashGen
// This module provides a random hash generation function. At this stage, this is not worth its own package, but it makes sense to do so here for the purpose of avoiding import cycles.
package randomhashgen

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
)

func GenerateRandomHash() (string, error) {
	const LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	saltBytes := make([]byte, 16)
	for i := range saltBytes {
		randNum, err := rand.Int(rand.Reader, big.NewInt(int64(len(LETTERS))))
		if err != nil {
			return "", errors.New(fmt.Sprint(
				"Random number generator generated an error. err: ", err))
		}
		saltBytes[i] = LETTERS[int(randNum.Int64())]
	}
	calculator := sha256.New()
	calculator.Write(saltBytes)
	resultHex := fmt.Sprintf("%x", calculator.Sum(nil))
	return resultHex, nil
}
