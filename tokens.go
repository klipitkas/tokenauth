package tokenauth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

// MaxTokenLength is the maximum token length.
const MaxTokenLength = 1000

// NewToken generates a new token of length l.
func NewToken(l int) (string, error) {
	if l <= 0 {
		return "", errors.New("a non-negative, non-zero integer is required for the length of the token")
	}
	if l > MaxTokenLength {
		return "", fmt.Errorf("the maximum token length is limited to %d characters", MaxTokenLength)
	}
	token, err := randomString(l)
	if err != nil {
		return "", fmt.Errorf("generate random string: %v", err)
	}
	return token, nil
}

func randomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", fmt.Errorf("generate uniform random value: %v", err)
		}
		ret[i] = letters[num.Int64()]
	}
	return string(ret), nil
}
