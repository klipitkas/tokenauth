package tokenauth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

// MaxTokenLength is the maximum token length.
const MaxTokenLength = 1000

// Alphabet for generating the tokens.
var Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// NewToken generates a new token of length l.
func NewToken(l int) (string, error) {
	if l <= 0 {
		return "", errors.New("a non-negative, non-zero integer is required for the length of the token")
	}
	if l > MaxTokenLength {
		return "", fmt.Errorf("the maximum token length is limited to %d characters", MaxTokenLength)
	}
	token, err := randomString(l, Alphabet)
	if err != nil {
		return "", fmt.Errorf("generate random string: %v", err)
	}
	return token, nil
}

// Securely generate a random string.
func randomString(n int, alphabet string) (string, error) {
	if len(alphabet) == 0 {
		return "", errors.New("alphabet is empty")
	}
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", fmt.Errorf("generate uniform random value: %v", err)
		}
		ret[i] = alphabet[num.Int64()]
	}
	return string(ret), nil
}
