package engine

import "testing"

func TestValidateFlags(t *testing.T) {
	allowlist := []string{"-O0", "-O1", "-O2", "-O3", "-Wall", "-std=*"}

	tests := []struct {
		name      string
		requested []string
		wantErr   bool
	}{
		{"Empty requested flags", []string{}, false},
		{"Exact match allowed flags", []string{"-O2", "-Wall"}, false},
		{"Wildcard matched flag", []string{"-std=c++17"}, false},
		{"Mixed valid exact and wildcard", []string{"-O3", "-std=c++20"}, false},
		{"Single malicious flag", []string{"-fplugin=malicious.so"}, true},
		{"Mixed valid and malicious flag", []string{"-O2", "-Wl,@exploit"}, true},
		{"Sneaky prefix flag", []string{"-Walll"}, true}, // Extra 'l' shouldn't pass
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFlags(tt.requested, allowlist)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFlags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}