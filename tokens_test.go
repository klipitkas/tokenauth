package tokenauth

import (
	"strings"
	"testing"
)

// go test -run Test_NewToken
func Test_NewToken(t *testing.T) {
	tests := []struct {
		desc            string
		conf            TokenConfig
		wantError       bool
		wantTokenLength int
		wantCharacters  string
	}{
		{
			desc:            "Token length is -1.",
			conf:            TokenConfig{Length: -1},
			wantError:       false,
			wantTokenLength: MinTokenLength,
			wantCharacters:  TokenConfigDefault.Alphabet,
		},
		{
			desc:            "Token length is zero.",
			conf:            TokenConfig{Length: 0},
			wantError:       false,
			wantTokenLength: MinTokenLength,
			wantCharacters:  TokenConfigDefault.Alphabet,
		},
		{
			desc:            "Token length is 20 and alphabet is default.",
			conf:            TokenConfig{Length: 20},
			wantError:       false,
			wantTokenLength: 20,
			wantCharacters:  TokenConfigDefault.Alphabet,
		}, {
			desc:            "Token length is 24 and alphabet is \"ABCabc012\"",
			conf:            TokenConfig{Length: 24, Alphabet: "ABCabc012"},
			wantError:       false,
			wantTokenLength: 24,
			wantCharacters:  "ABCabc012",
		}, {
			desc:            "Token length is 10000 and alphabet is default.",
			conf:            TokenConfig{Length: 10000},
			wantError:       false,
			wantTokenLength: MaxTokenLength,
			wantCharacters:  TokenConfigDefault.Alphabet,
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/kyoh86/scopelint/issues/4#issuecomment-471661062
		t.Run(tt.desc, func(t *testing.T) {
			token, err := NewToken(tt.conf)
			if tt.wantError && err == nil {
				t.Errorf("Wanted error for test: %s but got error = nil", tt.desc)
			}
			if len(token) != tt.wantTokenLength {
				t.Errorf("Want token length %d but got token length %d", tt.wantTokenLength, len(token))
			}
			for _, c := range token {
				if !strings.ContainsRune(tt.wantCharacters, c) {
					t.Errorf("Want character %c in alphabet %q, but it was not found", c, tt.wantCharacters)
				}
			}
		})
	}
}
