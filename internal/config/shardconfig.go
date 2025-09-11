package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type ShardCloudConfig struct {
	AppID       string
	DisplayName string
	Memory      string
	Version     string
}

const (
	SHARD_CLOUD_CONFIG_NAME = ".shardcloud"
)

func AddLine(basePath string, line string) error {
	file, err := os.OpenFile(filepath.Join(basePath, SHARD_CLOUD_CONFIG_NAME), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(line)
	return err
}

func GetAppConfig(basePath string) (*ShardCloudConfig, error) {
	file, err := os.OpenFile(filepath.Join(basePath, SHARD_CLOUD_CONFIG_NAME), os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	configMap := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			configMap[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(configMap) == 0 {
		return nil, nil
	}

	config := &ShardCloudConfig{
		AppID:       configMap["APPID"],
		DisplayName: configMap["DISPLAY_NAME"],
		Memory:      configMap["MEMORY"],
		Version:     configMap["VERSION"],
	}

	return config, nil
}
