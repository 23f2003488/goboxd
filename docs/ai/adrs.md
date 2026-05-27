## Safe WorkSpace Generation

**Context:**
The Python + Flask implementation is generating temporary directories by picking a random number between 30,000 to 60,000 and creating them using shell commands(`os.system("mkdir...)`). For phase 3 where it will be tested under the sustained load, this guarantees UID collisions and recursive retries. It is even blindly trusting the user's `source_filename` which allows path traversal.

**Options considered:**
1. Replicate the same Python logic with a larger random number range.
2. Using Go's native `os.MkdirTemp` and explicitly validating filenames.

**Decision:**
We will go with option 2. We built a `workspace` manager that uses `os.MkdirTemp` to guarantee atomicity and prevent collision in directory creation. We coupled this with a strict `ValidateFilename` function which rejects slashes and `..` characters. Finally, we used Go's `defer` keyword on the workspace cleanup function to guarantee directory deletion even if the HTTP handler panics.

**Rationale:**
This approach closes three known vulnerabilities (UID collisions, path traversal, and stale directories) using standard library features without relying on external dependencies or unsafe shell execution.