package main

import "net/http"

func (ac *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	ac.fileserverHits.Store(0)
	w.Write([]byte("Hits reset to 0"))
}
