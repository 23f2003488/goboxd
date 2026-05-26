package models

// RunRequest represents the incoming JSON Payload for POST /run
type RunRequest struct {
	Language         string       `json:"language"`
	Source           string       `json:"source"`
	SourceFilename   string       `json:"source_filename,omitempty"`
	ArtifactFilename string       `json:"artifact_filename,omitempty"`
	Build            *CommandSpec `json:"build,omitempty"`
	Run              *CommandSpec `json:"run,omitempty"`
	Tests            []TestCase   `json:"tests"`
}

type CommandSpec struct {
	Limits *Limits  `json:"limits,omitempty"`
	Flags  []string `json:"flags,omitempty"`
}

type Limits struct {
	WallTimeS    int `json:"wall_time_s,omitempty"`
	MemoryKB     int `json:"memory_kb,omitempty"`
	MaxProcesses int `json:"max_processes,omitempty"`
}

type TestCase struct {
	Stdin          string `json:"stdin"`
	ExpectedStdout string `json:"expected_stdout"`
}

// RunResponse represents the outgoing JSON Payload
type RunResponse struct {
	Status string       `json:"status"`
	Build  *StageResult `json:"build,omitempty"`
	Tests  []TestResult `json:"tests,omitempty"`
}

type StageResult struct {
	Status     string `json:"status"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	DurationMs int    `json:"duration_ms"`
}

type TestResult struct {
	Status       string `json:"status"`
	Stdout       string `json:"stdout"`
	Stderr       string `json:"stderr"`
	DurationMs   int    `json:"duration_ms"`
	MemoryPeakKB int    `json:"memory_peak_kb,omitempty"`
}