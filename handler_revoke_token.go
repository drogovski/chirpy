package main

import (
	"net/http"

	"github.com/drogovski/chirpy/internal/auth"
	"github.com/drogovski/chirpy/internal/database"
)

func (ac *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	q := database.New(ac.db)
	err = q.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
