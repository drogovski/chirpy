package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
	}

	err = validateChirp(params.Body)
	if err != nil {
		respondWithError(w, 400, err.Error(), err)
		return
	}

	respondWithJSON(w, 200, returnVals{
		Valid: true,
	})
}

func validateChirp(chirp string) error {
	const maxChirpLength = 140

	chirpLength := len(chirp)
	if chirpLength > maxChirpLength {
		return fmt.Errorf("the chirp is to long: %d", chirpLength)
	}
	return nil
}
