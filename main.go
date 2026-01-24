package main

import (
	"fmt"
	"os"
	"strings"

	"gmod-addon-manager/addon"
	"gmod-addon-manager/config"
	"gmod-addon-manager/tui"

	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	addonManager, err := addon.NewManager(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize addon manager: %v\n", err)
		os.Exit(1)
	}

	// Check if we should run in TUI mode (no arguments)
	if len(os.Args) == 1 {
		// Disable verbose output for TUI mode
		addonManager.SetVerbose(false)
		runTUI(addonManager)
		return
	}

	// Otherwise run in CLI mode (verbose output is already enabled by default)
	runCLI(addonManager, cfg)
}

func runTUI(manager *addon.Manager) {
	p := tea.NewProgram(tui.NewModel(manager), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

func runCLI(manager *addon.Manager, cfg *config.Config) {
	var rootCmd = &cobra.Command{
		Use:   "gmod-addon-manager",
		Short: "A TUI for managing Garry's Mod addons",
		Long:  "A terminal-based application for downloading, installing, and managing Garry's Mod addons",
	}

	rootCmd.AddCommand(initGetCmd(manager))
	rootCmd.AddCommand(initEnableCmd(manager))
	rootCmd.AddCommand(initDisableCmd(manager))
	rootCmd.AddCommand(initRemoveCmd(manager))
	rootCmd.AddCommand(initListCmd(manager))
	rootCmd.AddCommand(initInfoCmd(manager))
	rootCmd.AddCommand(initConfigCmd(cfg))

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

func initRemoveCmd(manager *addon.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "remove [addon-id]",
		Short: "Remove an addon (removes files and disables it)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := manager.RemoveAddon(args[0])
			if err != nil {
				fmt.Printf("Error removing addon: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Addon %s removed successfully\n", args[0])
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

func initConfigCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Run: func(cmd *cobra.Command, args []string) {
			// Show current config
			fmt.Println("Current Configuration:")
			fmt.Println("======================")
			fmt.Printf("GMod Directory: %s\n", cfg.GModDir)
			fmt.Printf("Download Directory: %s\n", cfg.DownloadDir)
			fmt.Printf("Addon Directory: %s\n", cfg.AddonDir)
			fmt.Printf("Output Directory: %s\n", cfg.OutDir)
			fmt.Printf("SteamCMD Path: %s\n", cfg.SteamCmdPath)
			fmt.Printf("GMAD Path: %s\n", cfg.GMADPath)
			fmt.Printf("Steam API Key: %s\n", cfg.SteamAPIKey)

			// Show config file location
			configPath, err := config.GetConfigPath()
			if err != nil {
				fmt.Printf("\nWarning: Could not determine config file location: %v\n", err)
			} else {
				fmt.Printf("\nConfig file location: %s\n", configPath)
				fmt.Println("You can edit this file to customize your configuration.")
			}
		},
	}
}
