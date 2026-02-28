package cmd

import (
	"fmt"
	"os"

	"github.com/rescoot/sunshine-cli/internal/output"
	"github.com/spf13/cobra"
)

var blinkersCmd = &cobra.Command{
	Use:   "blinkers [scooter-id] <left|right|both|off>",
	Short: "Control blinker lights",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var id int
		var state string
		var err error

		if len(args) == 2 {
			id, err = resolveScooterID(args[:1])
			state = args[1]
		} else {
			// Single arg — is it a state or an ID?
			if args[0] == "left" || args[0] == "right" || args[0] == "both" || args[0] == "off" {
				state = args[0]
				id, err = resolveScooterID(nil)
			} else {
				fmt.Fprintf(os.Stderr, "Usage: sunshine blinkers [scooter-id] <left|right|both|off>\n")
				os.Exit(1)
				return
			}
		}

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
