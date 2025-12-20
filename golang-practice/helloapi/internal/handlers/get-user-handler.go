package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type userResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(userResponse{
		ID:   uuid.NewString(),
		Name: "Gopher",
	})
}
