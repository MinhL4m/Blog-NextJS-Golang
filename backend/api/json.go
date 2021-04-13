package api

import (
	"encoding/json"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}
