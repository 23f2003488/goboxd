package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/thesouldev/goboxd/internal/config"
	"github.com/thesouldev/goboxd/internal/models"
)

// maxRequestSize is 256 KiB as per the spec
const maxRequestSize = 256 * 1024

// init() runs automatically before main() starts
func init() {
	// Load the embedded YAML file at startup.
	if err := config.LoadLanguages(); err != nil {
		log.Fatalf("Fatal: Could not load language registry: %v", err)
	}
	log.Printf("Successfully loaded %d languages into registry", len(config.GlobalRegistry))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	// Security Fix: Reject payloads that are larger than 256 Kib to prevent memory exhaustion
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)

	var req models.RunRequest
	decoder := json.NewDecoder(r.Body)
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

	// API contract validation - Language ID check
	_, exists := config.GlobalRegistry[req.Language]
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"code":    "invalid_language",
				"message": "The requested language id is not supported",
			},
		})
		return
	}

	// TODO: Pass req and langConfig to the Sandbox Engine
	
	// Temporary success response
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