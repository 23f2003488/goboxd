package config

import (
	_ "embed"
	"fmt"
	"github.com/thesouldev/goboxd/internal/models"
	"gopkg.in/yaml.v3"
)

//go:embed registry.yaml
var languagesYAML []byte

// GlobalRegistry holds our loaded languages in a map for O(1) instant lookups
var GlobalRegistry map[string]models.LanguageConfig

// LoadLanguages parses the embedded YAML file and populates the GlobalRegistry map
func LoadLanguages() error {
	var registry models.LanguageRegistry
	
	// We now unmarshal directly from the embedded byte array instead of reading from disk
	if err := yaml.Unmarshal(languagesYAML, &registry); err != nil {
		return fmt.Errorf("failed to parse embedded YAML: %v", err)
	}

	// Initialize the map
	GlobalRegistry = make(map[string]models.LanguageConfig)

	// Populate the map
	for _, lang := range registry.Languages {
		GlobalRegistry[lang.ID] = lang
	}

	return nil
}