package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func (c *Commands) RegisterStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the current app",
		RunE:  c.runStartCommand,
	}
	return cmd
}

func (c *Commands) runStartCommand(cmd *cobra.Command, args []string) error {
	appID := c.getAppId(args)
	err := c.cli.UpdateAppStatus(context.Background(), appID, "run")
	if err != nil {
		return err
	}

	fmt.Printf("App started")

	return nil
}
