package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validSecret := "mySecret"
	validToken, err := MakeJWT(userID, validSecret, time.Hour)
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name           string
		tokenString    string
		tokenSecret    string
		expectedUserID uuid.UUID
		wantErr        bool
	}{
		{
			"Valid JWT",
			validToken,
			validSecret,
			userID,
			false,
		},
		{
			"Invalid JWT",
			"is.not.valid",
			validSecret,
			uuid.Nil,
			true,
		},
		{
			"Wrong Token Secret",
			validToken,
			"wrongSecret",
			uuid.Nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if resultID != tt.expectedUserID {
				t.Errorf("ValidateJWT() = %s; expected %s", resultID, tt.expectedUserID)
			}
		})
	}
}
