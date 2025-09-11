package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/shard-cloud/cli/internal/config"
	"github.com/spf13/cobra"
)

func (c *Commands) RegisterLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "login",
		Short:       "Login to your Shard Cloud account",
		Annotations: map[string]string{"skipAuth": "true"},
		RunE:        c.runLoginCommand,
	}

	cmd.Flags().String("token", "", "")
	return cmd
}

func (c *Commands) runLoginCommand(cmd *cobra.Command, args []string) error {
	var token string
	tkn, err := cmd.Flags().GetString("token")

	if tkn == "" {
		prompt := promptui.Prompt{
			Label: "Enter the your shard cloud token",
			Validate: func(input string) error {
				if len(input) < 40 {
					return errors.New("token is invalid")
				}
				return nil
			},
		}

		token, err = prompt.Run()
		if err != nil {
			return err
		}
	}

	if tkn != "" {
		if err != nil {
			return fmt.Errorf("invalid token")
		}

		token = tkn
	}
	c.cli.SetToken(token)

	result, err := c.cli.Me(context.Background())
	if err != nil {
		return fmt.Errorf("invalid token")

	}

	fmt.Println("Logged in as", result.User.Name)
	if err := c.cli.Config.SaveConfigFile(&config.ConfigFile{
		Token: token,
	}); err != nil {
		return fmt.Errorf("failed to save config file: %w", err)
	}

	return nil
}
