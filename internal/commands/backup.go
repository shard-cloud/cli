package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func (c *Commands) RegisterBackupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup the current app",
		RunE:  c.runBackupCommand,
	}
	return cmd
}

func (c *Commands) runBackupCommand(cmd *cobra.Command, args []string) error {
	appID := c.getAppId(args)
	backup, err := c.cli.CreateBackup(context.Background(), appID)
	if err != nil {
		return err
	}

	fmt.Printf("Backup created with ID: %s", backup.Backup.ID)

	return nil
}
