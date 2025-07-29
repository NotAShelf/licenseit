package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Author string `json:"author"`
}

func ReadConfig(configFile string) (string, error) {
	if configFile != "" {
		content, err := os.ReadFile(configFile)
		if err != nil {
			return "", fmt.Errorf("could not read config file '%s': %w", configFile, err)
		}

		var config Config
		err = json.Unmarshal(content, &config)
		if err != nil {
			return "", fmt.Errorf("could not parse config file '%s': %w", configFile, err)
		}

		return config.Author, nil
	}

	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		configFile = filepath.Join(xdgConfigHome, "licenseit", "config.json")
		content, err := os.ReadFile(configFile)
		if err == nil {
			var config Config
			err = json.Unmarshal(content, &config)
			if err == nil {
				return config.Author, nil
			}
		}
	}

	return "", nil
}
