package main

import (
	"net/http"
	"time"

	"github.com/drogovski/chirpy/internal/auth"
	"github.com/drogovski/chirpy/internal/database"
)

func (ac *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	q := database.New(ac.db)
	refreshToken, err := q.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "provided refresh token is wrong", err)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "provided refresh token is no longer correct", err)
		return
	}

	newJWT, err := auth.MakeJWT(refreshToken.UserID, ac.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: newJWT,
	})
}
