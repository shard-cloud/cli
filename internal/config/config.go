package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

const (
	CONFIG_FILE_NAME = "shardcloud.json"
)

type ConfigFile struct {
	Token string `json:"token"`
}

func (c *ConfigFile) IsAuthenticated() bool {
	return c.Token != ""
}

type ConfigManager struct{}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{}
}

func (c *ConfigManager) SaveConfigFile(config *ConfigFile) error {
	configDir, err := c.ConfigDir()
	if err != nil {
		return err
	}
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return err
	}
	configFile := filepath.Join(configDir, CONFIG_FILE_NAME)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		_, err = os.Create(configFile)
		if err != nil {
			return err
		}
	}
	content, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, content, 0644)
}

func (c *ConfigManager) LoadConfigFile() *ConfigFile {
	configDir, err := c.ConfigDir()
	configFile := ConfigFile{
		Token: "",
	}
	if err != nil {
		return &configFile
	}
	configFilePath := filepath.Join(configDir, CONFIG_FILE_NAME)
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return &configFile
	}
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return &configFile
	}
	err = json.Unmarshal(content, &configFile)
	if err != nil {
		return &configFile
	}
	return &configFile

}

func (c *ConfigManager) ConfigDir() (string, error) {
	if appData := os.Getenv("AppData"); runtime.GOOS == "windows" && appData != "" {
		return filepath.Join(appData, "Shard Cloud CLI"), nil
	} else if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "shardcloud"), nil
	} else {
		dir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(dir, ".config", "shardcloud"), nil
	}
}
