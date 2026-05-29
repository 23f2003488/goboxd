package engine

import (
	"strings"
	"testing"
)

func TestValidateFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{"Valid standard filename", "solution.cpp", false},
		{"Valid short filename", "a.c", false},
		{"Empty filename rejected", "", true},
		{"Forward slash path traversal", "dir/file.c", true},
		{"Backward slash path traversal", "..\\etc\\passwd", true},
		{"Dot-dot traversal", "../secret.txt", true},
		{"Starts with a dot (hidden file)", ".bashrc", true},
		{"Exceeds length limit", strings.Repeat("a", 65) + ".py", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFilename(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFilename() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}