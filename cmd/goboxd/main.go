package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/thesouldev/goboxd/internal/config"
	"github.com/thesouldev/goboxd/internal/models"
	"github.com/thesouldev/goboxd/internal/engine"
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

	// API Contract Validation - Language ID Check
	langConfig, exists := config.GlobalRegistry[req.Language]
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

	// NEW: Determine the filename (fallback to registry default if not provided)
	filename := req.SourceFilename
	if filename == "" {
		filename = langConfig.SourceFilename
	}

	// NEW: Setup the secure workspace
	workspace, err := engine.SetupEnvironment(req.Source, filename)
	if err != nil {
		// Log the actual error for our debugging, but return 400 Bad Request to the user
		log.Printf("Workspace error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"code":    "bad_request",
				"message": err.Error(),
			},
		})
		return
	}
	// FIX: Stale Directories. 'defer' guarantees Cleanup() runs when this function exits!
	defer workspace.Cleanup() 

	// Temporary success response showing we successfully created and deleted the folder
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