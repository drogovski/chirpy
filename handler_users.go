package main

import (
	"encoding/json"
	"net/http"

	"github.com/drogovski/chirpy/internal/database"
)

func (ac *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	q := database.New(ac.db)
	user, err := q.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 400, "Couldn't create new user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
