package engine

import (
	"fmt"
	"strings"
)

// ValidateFilename ensures the filename is a single path component without any directory traversal characters.
func ValidateFilename(name string) error {
	if name == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return fmt.Errorf("filename cannot contain slashes")
	}
	if strings.Contains(name, "..") {
		return fmt.Errorf("filename cannot contain directory traversal (..)")
	}
	if strings.HasPrefix(name, ".") {
		return fmt.Errorf("filename cannot start with a dot")
	}
	if len(name) > 64 {
		return fmt.Errorf("filename exceeds maximum length of 64 characters")
	}
	return nil
}