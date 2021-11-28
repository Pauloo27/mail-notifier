package storage

import (
	"encoding/json"
	"os"
	"strings"
)

var (
	configFolder = ""
)

func init() {
	// tries to read env... if it's not found, use ~/.config
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "." // use working dir as fallback
		}
		configHome = home + "/.config/"
	}
	// add missing trailing `/` when it's needed
	if !strings.HasSuffix(configHome, "/") {
		configHome += "/"
	}
	configFolder = configHome + "mail-notifier/"
}

func LoadConfig() (*Config, error) {
	buf, err := os.ReadFile(configFolder + "config.json")
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(buf, &config)
	return &config, err
}

// TODO: save config
