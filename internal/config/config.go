package config

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed prompts.yaml
var promptsYAML []byte

type Prompts struct {
	SystemPrompt string `yaml:"system_prompt"`
}

var promptsConfig *Prompts

func LoadPrompts() (*Prompts, error) {
	if promptsConfig != nil {
		// Return cached config
		return promptsConfig, nil
	}

	var p Prompts
	if err := yaml.Unmarshal(promptsYAML, &p); err != nil {
		return nil, fmt.Errorf("failed to unmarshal prompts config: %w", err)
	}
	promptsConfig = &p
	return promptsConfig, nil
}
