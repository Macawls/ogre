package style

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// LoadTailwindConfig reads a Tailwind config from disk. The file's extension
// selects the format: .css → v4-style @theme block; .json → the
// TailwindConfig struct literal.
func LoadTailwindConfig(path string) (*TailwindConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lower := strings.ToLower(path)
	switch {
	case strings.HasSuffix(lower, ".css"):
		return ParseTailwindThemeCSS(string(data))
	case strings.HasSuffix(lower, ".json"):
		var cfg TailwindConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("parse json config: %w", err)
		}
		return &cfg, nil
	default:
		return nil, fmt.Errorf("unrecognised extension for tailwind config: %s (want .css or .json)", path)
	}
}
