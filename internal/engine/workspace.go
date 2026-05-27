package engine

import (
	"fmt"
	"os"
	"path/filepath"
)

// Workspace holds the paths for a single isolated request
type Workspace struct {
	HostDir    string // The folder on the server
	SourceFile string // The path to the user's code
}

// SetupEnvironment safely creates a collision-proof temporary directory
// and writes the user's source code into it.
func SetupEnvironment(sourceCode, filename string) (*Workspace, error) {
	// 1. Validate the filename to prevent Path Traversal
	if err := ValidateFilename(filename); err != nil {
		return nil, fmt.Errorf("security violation: %v", err)
	}

	// 2. Fix UID Collisions: os.MkdirTemp safely generates a unique folder
	// We use the default temp directory (usually /tmp in Linux)
	hostDir, err := os.MkdirTemp("", "goboxd-run-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create secure workspace: %v", err)
	}

	// 3. Construct the safe path and write the code
	sourcePath := filepath.Join(hostDir, filename)
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0644); err != nil {
		// If writing fails, we MUST clean up the directory we just created
		os.RemoveAll(hostDir) 
		return nil, fmt.Errorf("failed to write source code: %v", err)
	}

	return &Workspace{
		HostDir:    hostDir,
		SourceFile: sourcePath,
	}, nil
}

// Cleanup permanently deletes the workspace. 
// We will call this using `defer` in main.go to prevent Stale Directories.
func (w *Workspace) Cleanup() error {
	if w.HostDir != "" {
		return os.RemoveAll(w.HostDir)
	}
	return nil
}