package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func (c *Commands) RegisterStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the current app",
		RunE:  c.runStopCommand,
	}
	return cmd
}

func (c *Commands) runStopCommand(cmd *cobra.Command, args []string) error {
	appID := c.getAppId(args)
	err := c.cli.UpdateAppStatus(context.Background(), appID, "stop")
	if err != nil {
		return err
	}

	fmt.Printf("App stopped")

	return nil
}
