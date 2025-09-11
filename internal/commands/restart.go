package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func (c *Commands) RegisterRestartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart the current app",
		RunE:  c.runRestartCommand,
	}
	return cmd
}

func (c *Commands) runRestartCommand(cmd *cobra.Command, args []string) error {
	appID := c.getAppId(args)
	err := c.cli.UpdateAppStatus(context.Background(), appID, "restart")
	if err != nil {
		return err
	}

	fmt.Printf("App restarted")

	return nil
}
