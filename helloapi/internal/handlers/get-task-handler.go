package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type taskResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(taskResponse{
		ID:   uuid.NewString(),
		Name: "Deploy service to vps",
		Done: true,
	})
}
