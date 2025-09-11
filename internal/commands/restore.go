package commands

import (
	"context"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func (c *Commands) RegisterRestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore",
		Short: "restore the current app",
		RunE:  c.runRestoreCommand,
	}
	return cmd
}

func (c *Commands) runRestoreCommand(cmd *cobra.Command, args []string) error {
	appID := c.getAppId(args)
	backups, err := c.cli.GetBackups(context.Background(), appID)
	if err != nil {
		return err
	}

	items := make([]string, len(backups.Backups))
	for i, backup := range backups.Backups {
		items[i] = backup.CreatedAt.Format("2006-01-02 15:04:05") + " " + backup.ID.String()
	}
	prompt := promptui.Select{
		Label: "Select Backup",
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return err
	}

	var backupID string
	for _, backup := range backups.Backups {
		if backup.CreatedAt.Format("2006-01-02 15:04:05")+" "+backup.ID.String() == result {
			backupID = backup.ID.String()
			break
		}
	}

	if backupID == "" {
		return fmt.Errorf("backup not found")
	}

	err = c.cli.RestoreBackup(context.Background(), appID, backupID)
	if err != nil {
		return err
	}

	fmt.Printf("Backup restored with ID: %s", backupID)

	return nil
}
