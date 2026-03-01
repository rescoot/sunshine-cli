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
	Short: "Control the alarm system",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// No subcommand: show current alarm state from scooter detail
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
		if output.JSONOutput {
			output.PrintJSON(map[string]interface{}{
				"alarm_state":            scooter.AlarmState,
				"alarm_state_humanized":  scooter.AlarmStateHumanized,
				"alarm_triggered":        scooter.AlarmTriggered,
				"alarm_state_updated_at": scooter.AlarmStateUpdatedAt,
			})
		} else {
			fmt.Printf("Alarm: %s\n", scooter.AlarmStateHumanized)
		}
	},
}

var alarmTriggerCmd = &cobra.Command{
	Use:   "trigger [scooter-id]",
	Short: "Trigger the alarm for a duration",
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

var alarmArmCmd = &cobra.Command{
	Use:   "arm [scooter-id]",
	Short: "Arm the alarm system",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := resolveScooterID(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		resp, err := apiClient.AlarmArm(id)
		output.PrintCommandResponse(resp, err)
	},
}

var alarmDisarmCmd = &cobra.Command{
	Use:   "disarm [scooter-id]",
	Short: "Disarm the alarm system",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := resolveScooterID(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		resp, err := apiClient.AlarmDisarm(id)
		output.PrintCommandResponse(resp, err)
	},
}

var alarmStopCmd = &cobra.Command{
	Use:   "stop [scooter-id]",
	Short: "Stop an active alarm",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := resolveScooterID(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		resp, err := apiClient.AlarmStop(id)
		output.PrintCommandResponse(resp, err)
	},
}

func init() {
	alarmTriggerCmd.Flags().StringVar(&alarmDuration, "duration", "5s", "Alarm duration (e.g. 5s, 10s, 30s)")
	alarmCmd.AddCommand(alarmTriggerCmd)
	alarmCmd.AddCommand(alarmArmCmd)
	alarmCmd.AddCommand(alarmDisarmCmd)
	alarmCmd.AddCommand(alarmStopCmd)
	rootCmd.AddCommand(alarmCmd)
}
