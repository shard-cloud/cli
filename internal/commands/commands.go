package commands

import (
	"context"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/shard-cloud/cli/internal/cli"
	"github.com/shard-cloud/cli/internal/config"
)

type Commands struct {
	cli *cli.ShardCloud
}

func NewCommands(cli *cli.ShardCloud) *Commands {
	return &Commands{
		cli: cli,
	}
}

func (c *Commands) RegisterCommands() {
	c.cli.Cmd.AddCommand(c.RegisterLoginCommand())
	c.cli.Cmd.AddCommand(c.RegisterCreateAppCommand())
	c.cli.Cmd.AddCommand(c.RegisterCommitCommand())
	c.cli.Cmd.AddCommand(c.RegisterMeCommand())
	c.cli.Cmd.AddCommand(c.RegisterStartCommand())
	c.cli.Cmd.AddCommand(c.RegisterStopCommand())
	c.cli.Cmd.AddCommand(c.RegisterRestartCommand())
	c.cli.Cmd.AddCommand(c.RegisterLogoutCommand())
	c.cli.Cmd.AddCommand(c.RegisterDeleteCommand())
	c.cli.Cmd.AddCommand(c.RegisterStatusCommand())
	c.cli.Cmd.AddCommand(c.RegisterLogsCommand())
	c.cli.Cmd.AddCommand(c.RegisterBackupCommand())
	c.cli.Cmd.AddCommand(c.RegisterRestoreCommand())
}

func (c *Commands) getAppId(args []string) string {
	if len(args) > 0 {
		return args[0]
	}

	cwd, err := os.Getwd()
	if err == nil {
		config, _ := config.GetAppConfig(cwd)

		if config != nil && config.AppID != "" {
			return config.AppID
		}
	}

	apps, err := c.cli.GetAppsQuiet(context.Background())
	if err != nil {
		return ""
	}

	items := make([]string, len(apps))
	for i, app := range apps {
		items[i] = app.App.Name
	}
	prompt := promptui.Select{
		Label: "Select App",
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return ""
	}
	var appID string
	for _, app := range apps {
		if app.App.Name == result {
			appID = app.App.ID.String()
			break
		}
	}

	return appID
}
