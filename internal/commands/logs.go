package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func (c *Commands) RegisterLogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "See app logs",
		RunE:  c.runLogsCommand,
	}
	return cmd
}

func (c *Commands) runLogsCommand(cmd *cobra.Command, args []string) error {
	appID := c.getAppId(args)
	channel, err := c.cli.GetLogs(context.Background(), appID)
	if err != nil {
		return err
	}
	for log := range channel {
		split := strings.Split(log, " ")
		date := split[0]
		message := strings.Join(split[1:], " ")
		formatedDate, err := time.Parse("2006-01-02T15:04:05.000000000Z", date)
		if err != nil {
			return err
		}
		fmt.Println(formatedDate.Format("2006-01-02 15:04:05"), message)
	}

	return nil
}
