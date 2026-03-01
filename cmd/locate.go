package cmd

import (
	"fmt"
	"os"

	"github.com/rescoot/sunshine-cli/internal/output"
	"github.com/spf13/cobra"
)

var locateCmd = &cobra.Command{
	Use:   "locate",
	Short: "Locate the scooter (honk + flash)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := resolveScooterID(nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		resp, err := apiClient.Locate(id)
		output.PrintCommandResponse(resp, err)
	},
}

func init() {
	rootCmd.AddCommand(locateCmd)
}
