package commands

import (
	"context"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func (c *Commands) RegisterDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete the current app",
		RunE:  c.runDeleteCommand,
	}
	return cmd
}

func (c *Commands) runDeleteCommand(cmd *cobra.Command, args []string) error {
	appID := c.getAppId(args)

	app, err := c.cli.GetApp(context.Background(), appID)
	if err != nil {
		return err
	}
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Are you sure you want to delete app '%s' (%s)", app.App.Name, appID),
		IsConfirm: true,
	}

	_, err = prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			fmt.Println("Operation cancelled.")
			return nil
		}
		return fmt.Errorf("failed to get confirmation: %w", err)
	}
	err = c.cli.DeleteApp(context.Background(), appID)
	if err != nil {
		return err
	}

	fmt.Printf("App deleted")

	return nil
}
