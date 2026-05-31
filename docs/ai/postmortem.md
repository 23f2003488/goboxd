# Phase 1 Postmortem

**What we thought would be hard that turned out to be easy:**
We expected adding new languages to be highly complex, requiring deep Go routing logic. By utilizing a central `GlobalRegistry` map populated by a YAML file, the Go code acts entirely agnostically. Adding Node.js or Rust in Stage 2 will simply require a YAML block and a Dockerfile `apt-get` addition, with zero Go changes.

**What we thought would be easy that turned out to be hard:**
We vastly underestimated the difficulty of satisfying `nsjail`'s environment stripping. We successfully installed `build-essential` via Docker, but `g++` consistently crashed with `cannot find 'ld'`. It took hours to realize `nsjail` strips the `PATH` variable, leaving the compiler blind. Explicitly injecting `--env PATH=/usr/local/bin:/usr/bin:/bin` into the sandbox arguments solved an incredibly frustrating invisible bug.

**Where AI gave confident-sounding advice that was wrong:**
The AI repeatedly insisted that our missing compiler error was due to Docker caching mechanisms and suggested `--no-cache` builds. While logically sound for typical Docker issues, it masked the true root cause (the `nsjail` environment isolation). It cost us several slow rebuild cycles before we inspected the `nsjail` constraints directly.

**If we started Phase 1 again tomorrow, what would we change?**
We would write the `cappedWriter` (to prevent OOM) and the `defer workspace.Cleanup()` (to prevent disk leaks) before writing a single line of API routing. Retrofitting security onto a working API is harder than building the core execution primitives securely from day one.