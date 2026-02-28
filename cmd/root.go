package cmd

import (
	"fmt"
	"os"

	"github.com/rescoot/sunshine-cli/internal/api"
	"github.com/rescoot/sunshine-cli/internal/config"
	"github.com/rescoot/sunshine-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	cfg        *config.Config
	apiClient  *api.Client
	jsonOutput bool
	serverURL  string
	scooterID  int
)

var rootCmd = &cobra.Command{
	Use:   "sunshine",
	Short: "CLI for controlling Rescoot scooters",
	Long:  "sunshine is a command-line interface for the Rescoot Sunshine API.\nControl your electric scooter from the terminal.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		output.JSONOutput = jsonOutput

		var err error
		cfg, err = config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		if serverURL != "" {
			cfg.Server = serverURL
		}

		apiClient = api.NewClient(cfg)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", "", "Server URL (overrides config)")
	rootCmd.PersistentFlags().IntVar(&scooterID, "scooter", 0, "Scooter ID (overrides default)")
}

// resolveScooterID returns the scooter ID from args, flag, or config default.
func resolveScooterID(args []string) (int, error) {
	if len(args) > 0 {
		var id int
		if _, err := fmt.Sscanf(args[0], "%d", &id); err == nil && id > 0 {
			return id, nil
		}
	}

	if scooterID > 0 {
		return scooterID, nil
	}

	if cfg != nil && cfg.DefaultScooter > 0 {
		return cfg.DefaultScooter, nil
	}

	return 0, fmt.Errorf("scooter ID required — pass as argument, use --scooter flag, or set default_scooter in config")
}
