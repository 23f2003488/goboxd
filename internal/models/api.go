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
	WallTimeS    int `json:"wall_time_s,omitempty" yaml:"wall_time_s"`
	MemoryKB     int `json:"memory_kb,omitempty" yaml:"memory_kb"`
	MaxProcesses int `json:"max_processes,omitempty" yaml:"max_processes"`
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

// LanguageConfig maps the structure of a single language entry inside the YAML file
type LanguageConfig struct {
	ID             string       `yaml:"id"`
	Name           string       `yaml:"name"`
	SourceFilename string       `yaml:"source_filename"`
	Artifact       string       `yaml:"artifact,omitempty"`
	Build          *LangCmdSpec `yaml:"build,omitempty"`
	Run            LangCmdSpec  `yaml:"run"`
}

type LangCmdSpec struct {
	Cmd           string   `yaml:"cmd"`
	Args          []string `yaml:"args,omitempty"`
	Limits        Limits   `yaml:"limits"`
	FlagAllowlist []string `yaml:"flag_allowlist,omitempty"`
}

// Registry wrapper to capture the root array from the YAML file
type LanguageRegistry struct {
	Languages []LanguageConfig `yaml:"languages"`
}