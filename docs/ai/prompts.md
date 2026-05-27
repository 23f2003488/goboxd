## 25-05-2026 · Resolving nsjail clone() permission errors

**Prompt:**
When running nsjail in the tools container, I got "clone() failed: Operation not permitted". How do I fix this in my Docker environment?

**Response summary:**
The AI explained that Docker's default security profile blocks the CAP_SYS_ADMIN and clone privileges needed by nsjail. Suggested adding `privileged: true` to the tools service in docker-compose.yml.

**What we used / didn't use:**
Used the `privileged: true` flag in docker-compose.yml for the tools container. This successfully allowed nsjail to execute standalone commands and unblocked the environment.


## 26-05-2026 · Implementing request size limits for POST /run

**Prompt:**
How do I cleanly restrict the incoming JSON payload size in Go's native net/http server to protect the host against Out-Of-Memory (OOM) attacks from oversized requests?

**Response Summary:**
Suggested using `http.MaxBytesReader(w, r.Body, limit)` inside the HTTP handler function. This native wrapper limits the number of bytes read from the request body, returning a non-EOF error if the threshold is breached and terminating further memory allocation.

**What we used / didn't use:**
Used `http.MaxBytesReader` configured to 256 KiB (`256 * 1024` bytes) inside our `runHandler` to strictly enforce the server-imposed max size limit and block large adversarial inputs before they hit our decoder.


## 27-05-2026 · Securing workspace generation

**Prompt:**
How do I safely generate isolated temporary directories in Go for untrusted code execution, avoiding path traversal and UID collisions?

**Response summary:**
Suggested using Go's `os.MkdirTemp` to atomically generate collision-proof folders, writing a strict filename validator using `strings.Contains` to block `../` or slashes, and using `defer workspace.Cleanup()` to guarantee the folder is deleted when the request finishes.

**What we used / didn't use:**
Used the entire pattern. Created a `Workspace` struct to manage the lifecycle and enforce cleanup via `defer`, permanently closing the stale directory vulnerability.


## 27-05-2026 · Sandbox Execution & Output Bounding

**Prompt:**
How do I safely wrap `nsjail` execution in Go to prevent compiler flag injection and protect the host server from unbounded memory exhaustion (OOM) caused by runaway child output?

**Response summary:**
The AI suggested creating a custom `io.Writer` implementation (`cappedWriter`) that automatically truncates output and discards bytes after a 128 KiB threshold to prevent RAM exhaustion. For flag injection, it suggested a strict prefix-aware validation function against the YAML allowlist.

**What we used / didn't use:**
Implemented the custom `cappedWriter` and attached it directly to `exec.CommandContext.Stdout` and `Stderr`. Implemented the `ValidateFlags` function to strictly filter incoming flags against the language registry, explicitly closing two of the reference vulnerabilities.