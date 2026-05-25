## 24-05-2026 Docker WSL2 Backend Exhausting C: Drive Space

**What we were trying to do:**
I was trying to run the `tools` container for the very first time to manually test out the `nsjail` isolation which requires pulling out large Debian Linux images and compiling `nsjail` from source.

**What went wrong:**
Now, my machine's C: drive has very limited space left, and the massive Docker build cache immediately filled the drive causing my system to lock up. I attempted standard ways to move WSL data to E: drive using `wsl --export` and `wsl --import` but they failed; modern Docker desktop was forcefully bypassing these settings on startup and recreated its default data folders back on the C: drive.

**How we resolved it:**
Then I stopped fighting with the Docker's internal settings and used the operating system to trick it. I shut down the WSL, deleted the default `$env:LOCALAPPDATA\Docker\wsl` folder and used windows PowerShell to create a hard directory junction (`mklink /J`) pointing that excat C: drive path to a new clean folder on E: drive.

**What we learned:**
Docker Desktop on windows is very strict about enforcing its default installation paths. Using OS level hard symlinks is a reliable way of controlling storage routing than fighting with WSL configuration registeries.