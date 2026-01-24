package main

import (
	"fmt"
	"os"
	"strings"

	"gmod-addon-manager/addon"
	"gmod-addon-manager/config"

	"github.com/spf13/cobra"
)

func main() {
	cfg := config.NewDefaultConfig()
	addonManager := addon.NewManager(cfg)

	var rootCmd = &cobra.Command{
		Use:   "gmod-addon-manager",
		Short: "A TUI for managing Garry's Mod addons",
		Long:  "A terminal-based application for downloading, installing, and managing Garry's Mod addons",
	}

	rootCmd.AddCommand(initDownloadCmd(addonManager))
	rootCmd.AddCommand(initListCmd(addonManager))
	rootCmd.AddCommand(initInfoCmd(addonManager))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initDownloadCmd(manager *addon.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "download [addon-id]",
		Short: "Download and install an addon from Steam Workshop",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := manager.DownloadAddon(args[0])
			if err != nil {
				fmt.Printf("Error downloading addon: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Successfully downloaded and installed addon %s\n", args[0])
		},
	}
}
// AI! also printt Enabled status
func initListCmd(manager *addon.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all installed addons",
		Run: func(cmd *cobra.Command, args []string) {
			addons, err := manager.ListAddons()
			if err != nil {
				fmt.Printf("Error listing addons: %v\n", err)
				os.Exit(1)
			}

			if len(addons) == 0 {
				fmt.Println("No addons installed")
				return
			}

			fmt.Println("Installed Addons:")
			fmt.Println("=================")
			for _, addon := range addons {
				fmt.Printf("ID: %s\n", addon.ID)
				if addon.Title != "" {
					fmt.Printf("Title: %s\n", addon.Title)
				}
				if addon.Author != "" {
					fmt.Printf("Author: %s\n", addon.Author)
				}
				if len(addon.Tags) > 0 {
					fmt.Printf("Tags: %s\n", strings.Join(addon.Tags, ", "))
				}
				fmt.Println("------------------")
			}
		},
	}
}

func initInfoCmd(manager *addon.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "info [addon-id]",
		Short: "Show information about an addon",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			addonInfo, err := manager.GetAddonInfo(args[0])
			if err != nil {
				fmt.Printf("Error getting addon info: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Addon Information:\n")
			fmt.Printf("==================\n")
			fmt.Printf("ID: %s\n", addonInfo.ID)
			if addonInfo.Title != "" {
				fmt.Printf("Title: %s\n", addonInfo.Title)
			}
			if addonInfo.Author != "" {
				fmt.Printf("Author: %s\n", addonInfo.Author)
			}
			if addonInfo.Description != "" {
				fmt.Printf("Description: %s\n", addonInfo.Description)
			}
			if len(addonInfo.Tags) > 0 {
				fmt.Printf("Tags: %s\n", strings.Join(addonInfo.Tags, ", "))
			}
			fmt.Printf("Installed: %t\n", addonInfo.Installed)
			fmt.Printf("Enabled: %t\n", addonInfo.Enabled)
		},
	}
}
