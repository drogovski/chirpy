package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, 400, err.Error(), err)
		return
	}

	respondWithJSON(w, 200, returnVals{
		CleanedBody: cleaned,
	})
}

func validateChirp(chirp string) (string, error) {
	const maxChirpLength = 140

	chirpLength := len(chirp)
	if chirpLength > maxChirpLength {
		return "", fmt.Errorf("the chirp is to long: %d", chirpLength)
	}

	return validateBadWords(chirp), nil
}

func validateBadWords(chirp string) string {
	forbiddenWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	words := strings.Fields(chirp)

	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if forbiddenWords[lowerWord] {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
