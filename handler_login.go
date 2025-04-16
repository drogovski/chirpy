package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/drogovski/chirpy/internal/auth"
	"github.com/drogovski/chirpy/internal/database"
	"github.com/google/uuid"
)

func (ac *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	q := database.New(ac.db)
	user, err := q.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, ac.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate access JWT", err)
		return
	}

	refreshToken, err := prepareRefreshToken(user.ID, r.Context(), ac.db)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

func prepareRefreshToken(userID uuid.UUID, context context.Context, db *sql.DB) (string, error) {
	const refreshTokenDuration = 60 * 24 * time.Hour

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		return "", err
	}

	q := database.New(db)
	savedToken, err := q.CreateRefreshToken(context, database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userID,
		ExpiresAt: time.Now().Add(refreshTokenDuration),
	})
	if err != nil {
		return "", err
	}

	return savedToken.Token, nil
}
