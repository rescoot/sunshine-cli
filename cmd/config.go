package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rescoot/sunshine-cli/internal/config"
	"github.com/rescoot/sunshine-cli/internal/output"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Run: func(cmd *cobra.Command, args []string) {
		showConfig()
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		showConfig()
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long:  "Set a configuration value. Supported keys: server, default_scooter, client_id, output",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key, value := args[0], args[1]

		switch key {
		case "server":
			cfg.Server = value
		case "default_scooter":
			id, err := strconv.Atoi(value)
			if err != nil || id < 0 {
				fmt.Fprintf(os.Stderr, "Error: default_scooter must be a non-negative integer\n")
				os.Exit(1)
			}
			cfg.DefaultScooter = id
		case "client_id":
			cfg.ClientID = value
		case "output":
			if value != "text" && value != "json" {
				fmt.Fprintf(os.Stderr, "Error: output must be 'text' or 'json'\n")
				os.Exit(1)
			}
			cfg.Output = value
		default:
			fmt.Fprintf(os.Stderr, "Error: unknown key %q (supported: server, default_scooter, client_id, output)\n", key)
			os.Exit(1)
		}

		if err := cfg.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s = %s\n", key, value)
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print the config file path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Path())
	},
}

func showConfig() {
	if output.JSONOutput {
		output.PrintJSON(cfg)
		return
	}

	kv := [][]string{
		{"Server", cfg.Server},
		{"Client ID", cfg.ClientID},
		{"Output", cfg.Output},
	}
	if cfg.DefaultScooter > 0 {
		kv = append(kv, []string{"Default Scooter", strconv.Itoa(cfg.DefaultScooter)})
	}
	kv = append(kv, []string{"Config File", config.Path()})

	maxKey := 0
	for _, p := range kv {
		if len(p[0]) > maxKey {
			maxKey = len(p[0])
		}
	}
	for _, p := range kv {
		fmt.Printf("%-*s  %s\n", maxKey+1, p[0]+":", p[1])
	}
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configPathCmd)
	rootCmd.AddCommand(configCmd)
}
