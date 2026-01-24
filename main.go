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
	addonManager, err := addon.NewManager(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize addon manager: %v\n", err)
		os.Exit(1)
	}

	var rootCmd = &cobra.Command{
		Use:   "gmod-addon-manager",
		Short: "A TUI for managing Garry's Mod addons",
		Long:  "A terminal-based application for downloading, installing, and managing Garry's Mod addons",
	}

	rootCmd.AddCommand(initGetCmd(addonManager))
	rootCmd.AddCommand(initEnableCmd(addonManager))
	rootCmd.AddCommand(initDisableCmd(addonManager))
	rootCmd.AddCommand(initListCmd(addonManager))
	rootCmd.AddCommand(initInfoCmd(addonManager))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initGetCmd(manager *addon.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "get [addon-id]",
		Short: "Download and install an addon from Steam Workshop",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := manager.GetAddon(args[0])
			if err != nil {
				fmt.Printf("Error getting addon: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Successfully downloaded and installed addon %s\n", args[0])
		},
	}
}

func initEnableCmd(manager *addon.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "enable [addon-id]",
		Short: "Enable an installed addon",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := manager.EnableAddon(args[0])
			if err != nil {
				fmt.Printf("Error enabling addon: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Addon %s enabled successfully\n", args[0])
		},
	}
}

func initDisableCmd(manager *addon.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "disable [addon-id]",
		Short: "Disable an installed addon",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := manager.DisableAddon(args[0])
			if err != nil {
				fmt.Printf("Error disabling addon: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Addon %s disabled successfully\n", args[0])
		},
	}
}

func formatAddonInfo(addon addon.Addon) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ID: %s\n", addon.ID))
	if addon.Title != "" {
		sb.WriteString(fmt.Sprintf("Title: %s\n", addon.Title))
	}
	if addon.Author != "" {
		sb.WriteString(fmt.Sprintf("Author: %s\n", addon.Author))
	}
	if len(addon.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(addon.Tags, ", ")))
	}
	sb.WriteString(fmt.Sprintf("Enabled: %t\n", addon.Enabled))
	return sb.String()
}

func initListCmd(manager *addon.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all installed addons",
		Run: func(cmd *cobra.Command, args []string) {
			addons, err := manager.GetAddonsInfo()
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
				fmt.Print(formatAddonInfo(addon))
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

			fmt.Println("Addon Information:")
			fmt.Println("==================")
			fmt.Print(formatAddonInfo(*addonInfo))
			fmt.Printf("Installed: %t\n", addonInfo.Installed)
		},
	}
}
