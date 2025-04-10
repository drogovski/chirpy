package main

import (
	"net/http"

	"github.com/drogovski/chirpy/internal/database"
)

func (ac *apiConfig) handlerChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	q := database.New(ac.db)
	chirps, err := q.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps from db.", err)
		return
	}

	chirpsToJson := []Chirp{}

	for _, chirp := range chirps {
		chirpsToJson = append(chirpsToJson, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirpsToJson)
}
