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

func (c *Commands) RegisterCommitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Commit the current app",
		RunE:  c.runCommitCommand,
	}
	return cmd
}

func (c *Commands) runCommitCommand(cmd *cobra.Command, args []string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	config, err := config.GetAppConfig(currentDir)
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
	s.Suffix = "Committing app..."
	s.Start()
	_, err = c.cli.UpdateApp(context.Background(), config.AppID, zipData)
	if err != nil {
		return err
	}
	s.Stop()

	fmt.Printf("App %s committed", config.DisplayName)

	return nil
}
