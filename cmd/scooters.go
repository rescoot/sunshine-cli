package cmd

import (
	"fmt"
	"os"

	"github.com/rescoot/sunshine-cli/internal/output"
	"github.com/spf13/cobra"
)

var scootersCmd = &cobra.Command{
	Use:   "scooters",
	Short: "Manage scooters",
}

var scootersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your scooters",
	Run: func(cmd *cobra.Command, args []string) {
		scooters, err := apiClient.ListScooters()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		output.PrintScooterList(scooters)
	},
}

var scootersShowCmd = &cobra.Command{
	Use:   "show [scooter-id]",
	Short: "Show scooter details",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := resolveScooterID(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		scooter, err := apiClient.GetScooter(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		output.PrintScooterDetail(scooter)
	},
}

func init() {
	scootersCmd.AddCommand(scootersListCmd)
	scootersCmd.AddCommand(scootersShowCmd)
	rootCmd.AddCommand(scootersCmd)
}
