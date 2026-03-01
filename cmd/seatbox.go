package cmd

import (
	"fmt"
	"os"

	"github.com/rescoot/sunshine-cli/internal/output"
	"github.com/spf13/cobra"
)

var seatboxCmd = &cobra.Command{
	Use:   "seatbox",
	Short: "Open the seatbox",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := resolveScooterID(nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		resp, err := apiClient.OpenSeatbox(id)
		output.PrintCommandResponse(resp, err)
	},
}

func init() {
	rootCmd.AddCommand(seatboxCmd)
}
