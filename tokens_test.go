package tokenauth

import (
	"testing"

	"github.com/gofiber/fiber/v2/utils"
)

// go test -run Test_NewToken
func Test_NewToken(t *testing.T) {
	tests := []struct {
		length  int
		wantErr bool
	}{
		{
			length:  -1,
			wantErr: true,
		},
		{
			length:  0,
			wantErr: true,
		},
		{
			length:  16,
			wantErr: false,
		}, {
			length:  100000,
			wantErr: true,
		},
	}

	for _, tt := range tests {
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
	}
}
