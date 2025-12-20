package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type healthCheckResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func GrtHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	formattedTime := now.Format(time.RFC3339)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(healthCheckResponse{
		Status: "ok",
		Time:   formattedTime,
	})
}
