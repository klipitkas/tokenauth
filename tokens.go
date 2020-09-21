package tokenauth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

// MinTokenLength is the minimum token length.
const MinTokenLength = 16

// MaxTokenLength is the maximum token length.
const MaxTokenLength = 512

// TokenConfig is the struct that contains the config
// for token generation.
type TokenConfig struct {
	// Alphabet contains the characters that can be used as parts of the token.
	Alphabet string
	// Length is the length of the token to generate.
	Length int
	// Generator is the function that can be used to generate tokens.
	Generator func(int, string) (string, error)
}

// TokenConfigDefault is the default configuration for
// generating tokens.
var TokenConfigDefault = TokenConfig{
	// https://github.com/ai/nanoid/blob/master/url-alphabet/index.js
	Alphabet: "0123456789ABCDEFGHNRVfgctiUvz_KqYTJkLxpZXIjQW",
	Length:   MinTokenLength,
	Generator: func(n int, ab string) (string, error) {
		if len(ab) == 0 {
			return "", errors.New("alphabet is empty")
		}
		m := make([]byte, n)
		for i := 0; i < n; i++ {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(ab))))
			if err != nil {
				return "", fmt.Errorf("generate uniform random value: %v", err)
			}
			m[i] = ab[num.Int64()]
		}
		return string(m), nil
	},
}

// NewToken generates a new token of length l.
func NewToken(cfg TokenConfig) (string, error) {
	if cfg.Alphabet == "" {
		cfg.Alphabet = TokenConfigDefault.Alphabet
	}
	if cfg.Length < MinTokenLength {
		cfg.Length = MinTokenLength
	}
	if cfg.Length > MaxTokenLength {
		cfg.Length = MaxTokenLength
	}
	if cfg.Generator == nil {
		cfg.Generator = TokenConfigDefault.Generator
	}
	return cfg.Generator(cfg.Length, cfg.Alphabet)
}
