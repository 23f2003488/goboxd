package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/thesouldev/goboxd/internal/config"
	"github.com/thesouldev/goboxd/internal/engine"
	"github.com/thesouldev/goboxd/internal/models"
)

const maxRequestSize = 256 * 1024

func init() {
	if err := config.LoadLanguages(); err != nil {
		log.Fatalf("Fatal: Could not load language registry: %v", err)
	}
	log.Printf("Successfully loaded %d languages into registry", len(config.GlobalRegistry))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// sendError is a quick helper to format API errors consistently
func sendError(w http.ResponseWriter, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{"code": code, "message": message},
	})
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)

	var req models.RunRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		sendError(w, "bad_request", "Invalid JSON or payload exceeds size limit")
		return
	}

	langConfig, exists := config.GlobalRegistry[req.Language]
	if !exists {
		sendError(w, "invalid_language", "The requested language id is not supported")
		return
	}

	filename := req.SourceFilename
	if filename == "" {
		filename = langConfig.SourceFilename
	}

	workspace, err := engine.SetupEnvironment(req.Source, filename)
	if err != nil {
		sendError(w, "bad_request", err.Error())
		return
	}
	defer func() { _ = workspace.Cleanup() }()

	// --- PHASE 1: BUILD ---
	var buildResult *models.StageResult
	if langConfig.Build != nil {
		buildLimits := langConfig.Build.Limits
		var reqFlags []string
		
		if req.Build != nil {
			reqFlags = req.Build.Flags
			if req.Build.Limits != nil {
				if req.Build.Limits.WallTimeS > 0 { buildLimits.WallTimeS = req.Build.Limits.WallTimeS }
				if req.Build.Limits.MemoryKB > 0 { buildLimits.MemoryKB = req.Build.Limits.MemoryKB }
			}
		}

		// Security: Validate compiler flags against allowlist
		if err := engine.ValidateFlags(reqFlags, langConfig.Build.FlagAllowlist); err != nil {
			sendError(w, "disallowed_flag", err.Error())
			return
		}

		buildArgs := engine.ConstructArgs(langConfig.Build.Args, langConfig, reqFlags)
		stdout, stderr, duration, err := engine.ExecuteSandbox(workspace, buildLimits, langConfig.Build.Cmd, buildArgs, "")

		buildResult = &models.StageResult{
			Stdout:     stdout,
			Stderr:     stderr,
			DurationMs: duration,
			Status:     "ok",
		}

		if err != nil {
			buildResult.Status = "failed"
			var skippedTests []models.TestResult
			for range req.Tests {
				skippedTests = append(skippedTests, models.TestResult{Status: "not_executed"})
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(models.RunResponse{
				Status: "build_failed",
				Build:  buildResult,
				Tests:  skippedTests,
			})
			return
		}
	}

	// --- PHASE 2: RUN ---
	runLimits := langConfig.Run.Limits
	var runFlags []string
	if req.Run != nil {
		runFlags = req.Run.Flags
		if req.Run.Limits != nil {
			if req.Run.Limits.WallTimeS > 0 { runLimits.WallTimeS = req.Run.Limits.WallTimeS }
			if req.Run.Limits.MemoryKB > 0 { runLimits.MemoryKB = req.Run.Limits.MemoryKB }
		}
	}

	runArgs := engine.ConstructArgs(langConfig.Run.Args, langConfig, runFlags)
	var testResults []models.TestResult
	topLevelStatus := "accepted"

	for _, test := range req.Tests {
		runCmd := strings.ReplaceAll(langConfig.Run.Cmd, "{{artifact}}", langConfig.Artifact)
		runCmd = strings.ReplaceAll(runCmd, "{{source}}", filename)

		stdout, stderr, duration, err := engine.ExecuteSandbox(workspace, runLimits, runCmd, runArgs, test.Stdin)

		status := "accepted"
		if err != nil {
			// Basic status mapping for execution errors
			if duration >= runLimits.WallTimeS*1000 {
				status = "time_exceeded"
			} else {
				status = "runtime_error"
			}
		} else {
			// CP Standard: Strip trailing newlines from both before comparing
			cleanStdout := strings.TrimRight(stdout, "\r\n")
			cleanExpected := strings.TrimRight(test.ExpectedStdout, "\r\n")

			if cleanStdout != cleanExpected {
				// If they still don't match, check if it's just a space/tab issue
				if strings.TrimSpace(cleanStdout) == strings.TrimSpace(cleanExpected) {
					status = "output_whitespace_mismatch"
				} else {
					status = "wrong_output"
				}
			}
		}

		testResults = append(testResults, models.TestResult{
			Status:     status,
			Stdout:     stdout,
			Stderr:     stderr,
			DurationMs: duration,
		})

		// Track the first failure for the top-level status
		if status != "accepted" && topLevelStatus == "accepted" {
			topLevelStatus = status
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(models.RunResponse{
		Status: topLevelStatus,
		Build:  buildResult,
		Tests:  testResults,
	})
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