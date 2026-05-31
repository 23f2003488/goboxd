# Architecture

This document outlines how a request flows through the goboxd system. The codebase is separated into three internal packages to maintain a clean separation of concerns.

## 1. Config Layer (`internal/config`)
On startup, the system parses `languages.yml` and loads it into a global `GlobalRegistry` map. We use the standard library `//go:embed` directive to compile the YAML directly into the binary, preventing file-path resolution issues inside the Docker container. 

## 2. API Layer (`cmd/goboxd`)
The API strictly enforces the JSON contract using `json.Decoder.DisallowUnknownFields()`. Request sizes are capped at 256 KiB at the HTTP socket level using `http.MaxBytesReader` to prevent basic DoS attacks before the JSON parser even runs. 

## 3. Engine Layer (`internal/engine`)
The engine executes the core sandboxing logic in three steps:

* **Workspace Generation:** Uses `os.MkdirTemp` to generate a collision-proof directory on the host. `defer workspace.Cleanup()` guarantees deletion on all exit paths, preventing stale directories.
* **Security Validation:** Incoming compiler flags are strictly checked against a per-language prefix allowlist. Incoming filenames are scanned for path traversal characters (`../`, `/`).
* **Execution:** We wrap the `nsjail` process via `os/exec`. We inject `PATH` and `TMPDIR` environment variables so standard toolchains function within the read-only chroot. Output is piped through a custom `cappedWriter` struct that forcibly truncates `stdout`/`stderr` at 128 KiB to prevent runaway processes from exhausting host memory. Output whitespace is trimmed to conform to standard competitive programming evaluation rules.