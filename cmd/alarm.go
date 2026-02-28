package cmd

import (
	"fmt"
	"os"

	"github.com/rescoot/sunshine-cli/internal/output"
	"github.com/spf13/cobra"
)

var alarmDuration string

var alarmCmd = &cobra.Command{
	Use:   "alarm [scooter-id]",
	Short: "Trigger the alarm",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := resolveScooterID(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		resp, err := apiClient.Alarm(id, alarmDuration)
		output.PrintCommandResponse(resp, err)
	},
}

func init() {
	alarmCmd.Flags().StringVar(&alarmDuration, "duration", "5s", "Alarm duration (e.g. 5s, 10s, 30s)")
	rootCmd.AddCommand(alarmCmd)
}
