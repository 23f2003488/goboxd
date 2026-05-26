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