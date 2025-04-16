package main

import (
	"encoding/json"
	"net/http"

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
		respondWithError(w, http.StatusInternalServerError, "Couldn't decoded request parameters", err)
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
	q := database.New(ac.db)
	err := q.UpgradeToChirpyRed(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find user to upgrade", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
