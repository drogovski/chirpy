package main

import (
	"net/http"

	"github.com/drogovski/chirpy/internal/database"
	"github.com/google/uuid"
)

func (ac *apiConfig) handlerChirpsGetOne(w http.ResponseWriter, r *http.Request) {
	chirpIdString := r.PathValue("id")

	chirpId, err := uuid.Parse(chirpIdString)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	q := database.New(ac.db)
	chirp, err := q.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp with this id.", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
