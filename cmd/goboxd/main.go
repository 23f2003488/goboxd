package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// HealthResponse defines the JSON structure for the /healthz endpoint
type HealthResponse struct {
	Status string `json:"status"`
}

func main() {
	// Define our router
	mux := http.NewServeMux()

	// Register the /healthz endpoint
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := HealthResponse{Status: "ok"}
		json.NewEncoder(w).Encode(response)
	})

	// Start the server on port 8080
	log.Println("Starting goboxd server on :8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}