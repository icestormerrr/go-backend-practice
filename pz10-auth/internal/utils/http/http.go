package http_utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("[ERROR] Failed to encode JSON response: %v", err)
		return
	}

	log.Printf("[RESPONSE] %d %T %+v", http.StatusOK, v, v)
}

func WriteError(w http.ResponseWriter, code int, msg string, details any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := map[string]any{
		"error":   msg,
		"details": details,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[ERROR] Failed to encode error response: %v", err)
		return
	}

	log.Printf("[ERROR RESPONSE] %d %s | details: %+v", code, msg, details)
}
