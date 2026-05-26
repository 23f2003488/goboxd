package main

import (
	"encoding/json"
	"log"
	"net/http"
	"goboxd/internal/models"
)

// maxRequestSize is 256 KiB as per the spec
const maxRequestSize = 256 * 1024

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	// Security Fix: Reject payloads that are larger than 256 Kib to prevent memory exhaustion
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)

	var req models.RunRequest
	decoder := json.NewDecoder(r.body)
	//Disallow unknown fields to strictly enforce API contract
	decoder.DisallowUnknownFields()      
	
	if err := decoder.Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"code": "bad_request",
				"message": "INVALID JSON or playload exceeds 256 KiB limit",
			},
		})
		return
	}
	// TODO: Pass the validated request to the SandBox Engine

	// Placeholder Response until the engine is built
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.RunResponse{Status: "accepted"})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthzHandler)
	mux.HandleFunc("POST /run", runHandler)

	log.Println("Starting goboxd server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}