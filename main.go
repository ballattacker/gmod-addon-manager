package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gmod-addon-manager",
		Short: "A TUI for managing Garry's Mod addons",
		Long:  "A terminal-based application for downloading, installing, and managing Garry's Mod addons",
	}

	rootCmd.AddCommand(initDownloadCmd())
	rootCmd.AddCommand(initListCmd())
	rootCmd.AddCommand(initInfoCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initDownloadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "download [addon-id]",
		Short: "Download and install an addon from Steam Workshop",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Would download addon with ID: %s\n", args[0])
		},
	}
}

func initListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all installed addons",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Would list all installed addons")
		},
	}
}

func initInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info [addon-id]",
		Short: "Show information about an addon",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Would show info for addon with ID: %s\n", args[0])
		},
	}
}
