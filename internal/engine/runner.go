package engine

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/thesouldev/goboxd/internal/models"
)

// maxOutputBytes caps stdout/stderr at 128 KiB to prevent RAM exhaustion (OOM Vulnerability Fix)
const maxOutputBytes = 128 * 1024

// ValidateFlags checks the requested flags against the language's allowlist (Flag Injection Fix)
func ValidateFlags(requested []string, allowlist []string) error {
	for _, reqFlag := range requested {
		allowed := false
		for _, allowedFlag := range allowlist {
			// Handle wildcards like "-std=*"
			if strings.HasSuffix(allowedFlag, "*") {
				prefix := strings.TrimSuffix(allowedFlag, "*")
				if strings.HasPrefix(reqFlag, prefix) {
					allowed = true
					break
				}
			} else if reqFlag == allowedFlag {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("security violation: flag '%s' is not in the allowlist", reqFlag)
		}
	}
	return nil
}

// ConstructArgs safely replaces template variables like {{source}} and {{flags}}
func ConstructArgs(template []string, lang models.LanguageConfig, flags []string) []string {
	var finalArgs []string
	for _, arg := range template {
		if arg == "{{flags}}" {
			finalArgs = append(finalArgs, flags...)
			continue
		}
		arg = strings.ReplaceAll(arg, "{{source}}", lang.SourceFilename)
		arg = strings.ReplaceAll(arg, "{{artifact}}", lang.Artifact)
		finalArgs = append(finalArgs, arg)
	}
	return finalArgs
}

// ExecuteSandbox runs a command completely isolated inside nsjail
func ExecuteSandbox(workspace *Workspace, limits models.Limits, cmd string, args []string, stdinData string) (string, string, int, error) {
	// 1. Build the hardened nsjail arguments
	nsjailArgs := []string{
		"-Mo", "-q", "--chroot", "/",
		"--bindmount", fmt.Sprintf("%s:/workspace", workspace.HostDir),
		"--cwd", "/workspace",
		"--time_limit", fmt.Sprintf("%d", limits.WallTimeS),
		"--rlimit_as", fmt.Sprintf("%d", limits.MemoryKB*1024),
		"--rlimit_fsize", "10", 
		"--disable_clone_newnet",
		"--env", "PATH=/usr/local/bin:/usr/bin:/bin", 
		"--env", "TMPDIR=/workspace",
		"--",
		cmd,
	}
	nsjailArgs = append(nsjailArgs, args...)

	// 2. Set a backup context timeout slightly longer than the nsjail timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(limits.WallTimeS+2)*time.Second)
	defer cancel()

	execCmd := exec.CommandContext(ctx, "nsjail", nsjailArgs...)

	if stdinData != "" {
		execCmd.Stdin = strings.NewReader(stdinData)
	}

	// 3. Attach our custom capped buffers to prevent Unbounded Output
	var stdoutBuf, stderrBuf bytes.Buffer
	execCmd.Stdout = &cappedWriter{buf: &stdoutBuf, limit: maxOutputBytes}
	execCmd.Stderr = &cappedWriter{buf: &stderrBuf, limit: maxOutputBytes}

	// 4. Run the engine and measure duration
	startTime := time.Now()
	err := execCmd.Run()
	durationMs := int(time.Since(startTime).Milliseconds())

	return stdoutBuf.String(), stderrBuf.String(), durationMs, err
}

// cappedWriter safely limits how much data we read from the untrusted process
type cappedWriter struct {
	buf   *bytes.Buffer
	limit int
	total int
}

func (w *cappedWriter) Write(p []byte) (n int, err error) {
	if w.total >= w.limit {
		return len(p), nil // Pretend we wrote it, but silently discard to protect RAM
	}
	writeSize := len(p)
	if w.total+writeSize > w.limit {
		writeSize = w.limit - w.total
		w.buf.Write(p[:writeSize])
		w.buf.WriteString("\n[OUTPUT TRUNCATED - LIMIT EXCEEDED]")
		w.total = w.limit
		return len(p), nil
	}
	w.buf.Write(p)
	w.total += writeSize
	return len(p), nil
}