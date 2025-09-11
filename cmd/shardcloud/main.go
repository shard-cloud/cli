package main

import (
	"context"
	"fmt"

	"github.com/shard-cloud/cli/internal/cli"
	"github.com/shard-cloud/cli/internal/commands"
	"github.com/shard-cloud/cli/internal/config"
	"github.com/spf13/cobra"
)

var version = "1.0.0"

func main() {
	configManager := config.NewConfigManager()
	config := configManager.LoadConfigFile()
	cmd := &cobra.Command{
		Use:   "shard COMMAND",
		Short: "A CLI for Shard Cloud",
		//	SilenceErrors:     true,
		SilenceUsage:      true,
		TraverseChildren:  true,
		ValidArgsFunction: cobra.NoFileCompletions,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   false,
			HiddenDefaultCmd:    true,
			DisableDescriptions: true,
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			isHelpCommand := cmd.Name() == "help" || cmd.Name() == cobra.ShellCompRequestCmd || cmd.Name() == cobra.ShellCompNoDescRequestCmd || cmd.Name() == "shard"
			if !config.IsAuthenticated() && cmd.Annotations["skipAuth"] != "true" && !isHelpCommand {
				return fmt.Errorf("not authenticated")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return fmt.Errorf("%s is not a command. See 'shardcloud --help'", args[0])
		},

		Version: version,
	}

	cmd.SetVersionTemplate("Shard Cloud CLI version {{.Version}}\n")
	cmd.Flags().BoolP("version", "v", false, "Print CLI version")

	cmd.Flags().BoolP("debug", "d", false, "Debug Mode")
	cmd.Flags().MarkHidden("debug")
	shardCloud := cli.NewShardCloud(config.Token, configManager, cmd)
	commands := commands.NewCommands(shardCloud)
	commands.RegisterCommands()
	shardCloud.Cmd.ExecuteContext(context.Background())

	latestVersion, err := cli.GetLatestVersion()
	if err != nil {
		return
	}
	if latestVersion != version {
		fmt.Printf("\nA new version of the Shard Cloud CLI %s is available. Please update to the latest version.\n", latestVersion)
	}
}
