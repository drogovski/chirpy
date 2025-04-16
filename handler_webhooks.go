package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/drogovski/chirpy/internal/auth"
	"github.com/drogovski/chirpy/internal/database"
	"github.com/google/uuid"
)

type parameters struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (ac *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request parameters", err)
		return
	}

	switch params.Event {
	case "user.upgraded":
		ac.handleUserUpgrade(w, r, params.Data.UserID)
		return
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}

func (ac *apiConfig) handleUserUpgrade(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find api key", err)
		return
	}

	if apiKey != ac.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Provided api key is incorrect", err)
		return
	}

	q := database.New(ac.db)
	_, err = q.UpgradeToChirpyRed(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
