package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func (c *Commands) RegisterMeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "me",
		Short: "Get your Shard Cloud account information",
		RunE:  c.runMeCommand,
	}

	return cmd
}

func (c *Commands) runMeCommand(cmd *cobra.Command, args []string) error {
	result, err := c.cli.Me(context.Background())
	if err != nil {
		return fmt.Errorf("invalid token")
	}
	fmt.Println("Logged in as", result.User.Name)
	return nil
}
