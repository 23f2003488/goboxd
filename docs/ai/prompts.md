## 2026-05-24 · Resolving nsjail clone() permission errors

**Prompt:**
When running nsjail in the tools container, I got "clone() failed: Operation not permitted". How do I fix this in my Docker environment?

**Response summary:**
The AI explained that Docker's default security profile blocks the CAP_SYS_ADMIN and clone privileges needed by nsjail. Suggested adding `privileged: true` to the tools service in docker-compose.yml.

**What we used / didn't use:**
Used the `privileged: true` flag in docker-compose.yml for the tools container. This successfully allowed nsjail to execute standalone commands and unblocked the environment.