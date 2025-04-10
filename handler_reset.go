package main

import (
	"fmt"
	"net/http"

	"github.com/drogovski/chirpy/internal/database"
)

func (ac *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	type respVals struct {
		Message string `json:"message"`
	}

	if ac.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Couldn't reset the user table",
			fmt.Errorf("you cannot reset user table on not dev enviroment"))
		return
	}

	q := database.New(ac.db)
	err := q.Reset(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't reset the user table", err)
		return
	}

	respondWithJSON(w, http.StatusOK, respVals{
		Message: "The table was reset successfully",
	})
}
