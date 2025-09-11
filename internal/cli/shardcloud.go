package cli

import (
	"encoding/json"
	"fmt"
	"time"

	httpGO "net/http"

	"github.com/shard-cloud/cli/internal/config"
	"github.com/shard-cloud/cli/internal/http"
	"github.com/spf13/cobra"
)

const (
	BaseURL = "https://shardcloud.app/api"
)

var No_auth_commands = []string{
	"login",
}

type ShardCloud struct {
	*http.ShardCloudClient
	Config *config.ConfigManager
	Cmd    *cobra.Command
}

func (c *ShardCloud) SetToken(token string) {
	c.HttpClient.SetToken(token)
}

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

func GetLatestVersion() (string, error) {
	url := "https://api.github.com/repos/shard-cloud/cli/releases/latest"

	client := &httpGO.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpGO.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if release.TagName == "" {
		return "", fmt.Errorf("no tag_name found in release")
	}

	return release.TagName, nil
}

func NewShardCloud(token string, config *config.ConfigManager, cmd *cobra.Command) *ShardCloud {
	httpClient := http.NewShardCloudClient(token, BaseURL)
	shardCloud := &ShardCloud{
		ShardCloudClient: httpClient,
		Config:           config,
		Cmd:              cmd,
	}

	return shardCloud
}
