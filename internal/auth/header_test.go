package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"

	validHeaders := http.Header{}
	validHeaders.Add("Authorization", "Bearer "+bearerToken)

	tests := []struct {
		name          string
		headers       http.Header
		expectedToken string
		wantErr       bool
	}{
		{
			name:          "Found",
			headers:       validHeaders,
			expectedToken: bearerToken,
			wantErr:       false,
		},
		{
			name:          "Not Found",
			headers:       http.Header{},
			expectedToken: "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() err = %v; wantErr %v", err, tt.wantErr)
				return
			}

			if token != tt.expectedToken {
				t.Errorf("GetBearerToken() = %v; expected %v", token, tt.expectedToken)
			}
		})
	}
}
