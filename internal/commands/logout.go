package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (c *Commands) RegisterLogoutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from your Shard Cloud account",
		RunE:  c.runLogoutCommand,
	}
	return cmd
}

func (c *Commands) runLogoutCommand(cmd *cobra.Command, args []string) error {
	config := c.cli.Config.LoadConfigFile()
	config.Token = ""
	c.cli.Config.SaveConfigFile(config)

	fmt.Printf("Logged out")
	return nil
}
