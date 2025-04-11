package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.hash, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckJWTTokenGeneration(t *testing.T) {
	userID := uuid.New()
	secret1 := "secret1"
	secret2 := "secret2"
	durationHour := 1 * time.Hour
	durationHourAgo := -1 * time.Hour
	correctToken, _ := MakeJWT(userID, secret1, durationHour)
	expiredToken, _ := MakeJWT(userID, secret1, durationHourAgo)

	tests := []struct {
		name             string
		userID           uuid.UUID
		token            string
		validationSecret string
		wantErr          bool
	}{
		{
			name:             "Correct secret and time not expired",
			userID:           userID,
			token:            correctToken,
			validationSecret: secret1,
			wantErr:          false,
		},
		{
			name:             "Wrong secret and time not expired",
			userID:           userID,
			token:            correctToken,
			validationSecret: secret2,
			wantErr:          true,
		},
		{
			name:             "Correct secret and time expired",
			userID:           userID,
			token:            expiredToken,
			validationSecret: secret1,
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := ValidateJWT(tt.token, tt.validationSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
