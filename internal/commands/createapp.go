package commands

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/briandowns/spinner"
	"github.com/shard-cloud/cli/internal/config"
	"github.com/shard-cloud/cli/pkg/zip"
	"github.com/spf13/cobra"
)

func (c *Commands) RegisterCreateAppCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new app",
		RunE:  c.runCreateAppCommand,
	}
	return cmd
}

func (c *Commands) runCreateAppCommand(cmd *cobra.Command, args []string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	file, _ := cmd.Flags().GetString("file")
	var zipData []byte
	if file == "" {

		zipData, err = zip.ZipFolderToBytes(currentDir, []string{})
		if err != nil {
			return err
		}

	} else {
		zipData, err = os.ReadFile(path.Join(currentDir, file))
		if err != nil {
			return err
		}
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = "Creating app..."
	s.Start()
	app, err := c.cli.CreateApp(context.Background(), zipData)
	if err != nil {
		return err
	}
	s.Stop()

	fmt.Println("App created with ID:", app.ID)
	config.AddLine(currentDir, fmt.Sprintf("\nAPPID=%s", app.ID.String()))
	return nil
}
