package auth

import (
	"errors"
	"net/http"
	"testing"

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

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validSecret := "secret"
	invalidSecret := "someWrongSecret"
	validToken, _ := MakeJWT(userID, validSecret)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: validSecret,
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "some.invalid.token",
			tokenSecret: validSecret,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: invalidSecret,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	header1 := http.Header{}
	header1.Add("Authorization", "Bearer somelongtoken")
	header2 := http.Header{}
	header3 := http.Header{}
	header3.Add("Authorization", "Some token")
	header4 := http.Header{}
	header4.Add("Authorization", "Bearer")

	tests := []struct {
		name    string
		header  http.Header
		want    string
		wantErr error
	}{
		{
			name:    "Proper header format",
			header:  header1,
			want:    "somelongtoken",
			wantErr: nil,
		},
		{
			name:    "No authorization header",
			header:  header2,
			want:    "",
			wantErr: errors.New("the request doesn't have authorization header"),
		},
		{
			name:    "Authorization header with wrong format",
			header:  header3,
			want:    "",
			wantErr: errors.New("provided token has wrong format"),
		},
		{
			name:    "Authorization header with correct prefix but without token",
			header:  header4,
			want:    "",
			wantErr: errors.New("there was no token value: Bearer <token>"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBearerToken(tt.header)

			if tt.wantErr != nil && err == nil {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && err != nil {
				t.Errorf("GetBearerToken() unexpected error: %v", err)
			}

			if tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}

			if got != tt.want {
				t.Errorf("expected result %s, got %s", tt.want, got)
			}
		})
	}
}
