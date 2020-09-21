package tokenauth

import (
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2/utils"
)

// go test -run Test_NewToken
func Test_NewToken(t *testing.T) {
	tests := []struct {
		length   int
		wantErr  bool
		alphabet string
	}{
		{
			length:   -1,
			alphabet: Alphabet,
			wantErr:  true,
		},
		{
			length:   0,
			alphabet: Alphabet,
			wantErr:  true,
		},
		{
			length:   16,
			alphabet: Alphabet,
			wantErr:  false,
		}, {
			length:   16,
			alphabet: "abc",
			wantErr:  false,
		}, {
			length:   100000,
			alphabet: "abcdef",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		if len(tt.alphabet) > 0 {
			Alphabet = tt.alphabet
		}
		token, err := NewToken(tt.length)
		if !tt.wantErr {
			utils.AssertEqual(t, nil, err)
		}
		if tt.wantErr {
			utils.AssertEqual(t, "", token)
		}
		if tt.length > 0 && tt.length < MaxTokenLength {
			utils.AssertEqual(t, len(token), tt.length)
		}
		if len(tt.alphabet) > 0 {
			for _, r := range token {
				if !strings.ContainsRune(tt.alphabet, r) {
					t.Fatalf("string %s does not contain rune: %c", tt.alphabet, r)
				}
			}
		}
	}
}
