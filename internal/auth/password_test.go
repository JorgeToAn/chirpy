package auth

import "testing"

func TestCheckPasswordHash(t *testing.T) {
	validPassword := "password123"
	validHash, err := HashPassword(validPassword)
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name          string
		password      string
		hash          string
		expectedMatch bool
		wantErr       bool
	}{
		{
			name:          "Password Match",
			password:      validPassword,
			hash:          validHash,
			expectedMatch: true,
			wantErr:       false,
		},
		{
			name:          "Password No Match",
			password:      "other-passw0rd",
			hash:          validHash,
			expectedMatch: false,
			wantErr:       false,
		},
		{
			name:          "Invalid Hash",
			password:      validPassword,
			hash:          "not-valid",
			expectedMatch: false,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() err = %v; wantErr = %v", err, tt.wantErr)
				return
			}

			if match != tt.expectedMatch {
				t.Errorf("CheckPasswordHash() = %v; expected = %v", match, tt.expectedMatch)
			}
		})
	}
}
