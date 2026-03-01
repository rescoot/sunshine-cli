package cmd

import (
	"fmt"
	"os"

	"github.com/rescoot/sunshine-cli/internal/output"
	"github.com/spf13/cobra"
)

var blinkersCmd = &cobra.Command{
	Use:   "blinkers <left|right|both|off>",
	Short: "Control blinker lights",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		state := args[0]
		if state != "left" && state != "right" && state != "both" && state != "off" {
			fmt.Fprintf(os.Stderr, "Error: state must be left, right, both, or off\n")
			os.Exit(1)
		}

		id, err := resolveScooterID(nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		resp, apiErr := apiClient.SetBlinkers(id, state)
		output.PrintCommandResponse(resp, apiErr)
	},
}

func init() {
	rootCmd.AddCommand(blinkersCmd)
}
