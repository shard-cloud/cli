package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func (c *Commands) RegisterStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Status the current app",
		RunE:  c.runStatusCommand,
	}
	return cmd
}

func (c *Commands) runStatusCommand(cmd *cobra.Command, args []string) error {
	// Handle Ctrl+C gracefully
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nExiting...")
		cancel()
	}()

	// Clear screen initially
	clearScreen()

	ticker := time.NewTicker(8 * time.Second)
	defer ticker.Stop()

	// Display stats immediately, then every 8 seconds
	if err := c.displayAppStats(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := c.displayAppStats(ctx); err != nil {
				fmt.Printf("Error fetching stats: %v\n", err)
				continue
			}
		}
	}
}

func (c *Commands) displayAppStats(ctx context.Context) error {
	apps, err := c.cli.GetApps(ctx)
	if err != nil {
		return err
	}

	// Move cursor to top of screen (preserve clear screen effect)
	fmt.Print("\033[H")

	// Print header with timestamp
	fmt.Printf("SHARD CLOUD APPS STATUS - %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	if len(apps) == 0 {
		fmt.Println("No apps found.")
		return nil
	}

	// Print table header with fixed-width formatting
	fmt.Printf("%-15s %-10s %-22s %-18s %-16s\n",
		"NAME", "STATUS", "RAM USAGE", "CPU USAGE", "CREATED")
	fmt.Println(strings.Repeat("-", 82))

	for _, appResp := range apps {
		app := appResp.App

		// Format RAM usage
		ramUsage := formatMemory(appResp.RealTimeRam / 1024)
		ramLimit := formatMemory(app.Ram)
		ramPercent := float64(appResp.RealTimeRam/1024) / float64(app.Ram) * 100
		ramDisplay := fmt.Sprintf("%s / %s (%.1f%%)", ramUsage, ramLimit, ramPercent)

		// Format CPU usage
		cpuUsage := fmt.Sprintf("%.1f", float64(appResp.RealTimeVcpu)/1000)
		cpuLimit := fmt.Sprintf("%.1f", float64(app.Vcpu)/1000)
		cpuPercent := float64(appResp.RealTimeVcpu) / float64(app.Vcpu) * 100
		cpuDisplay := fmt.Sprintf("%s / %s (%.1f%%)", cpuUsage, cpuLimit, cpuPercent)

		// Format creation time
		createdAt := app.CreatedAt.Format("2006-01-02 15:04")

		// Truncate long names
		name := app.Name
		if len(name) > 13 {
			name = name[:10] + "..."
		}

		// Print with fixed-width columns
		fmt.Printf("%-15s %-10s %-22s %-18s %-16s\n",
			name,
			getStatusDisplay(appResp.Status),
			ramDisplay,
			cpuDisplay,
			createdAt,
		)
	}
	fmt.Println("\nPress Ctrl+C to exit...")
	return nil
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func formatMemory(mb int) string {
	if mb < 1024 {
		return fmt.Sprintf("%dMB", mb)
	}
	return fmt.Sprintf("%.1fGB", float64(mb)/1024.0)
}

func getStatusDisplay(status string) string {
	switch strings.ToLower(status) {
	case "running":
		return "🟢 Running"
	case "stopped":
		return "🔴 Stopped"
	case "pending":
		return "🟡 Pending"
	case "error":
		return "❌ Error"
	default:
		return status
	}
}
