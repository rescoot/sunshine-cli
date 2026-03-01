package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rescoot/sunshine-cli/internal/output"
	"github.com/spf13/cobra"
)

var navigateCmd = &cobra.Command{
	Use:   "navigate [scooter-id] <lat> <lng> [title]",
	Short: "Set navigation destination",
	Long:  "Set a navigation destination. Pass lat/lng as positional args, optionally with a title.\nUse 'navigate show' to view or 'navigate clear' to remove the current destination.",
	Args:  cobra.RangeArgs(2, 4),
	Run: func(cmd *cobra.Command, args []string) {
		var id int
		var latStr, lngStr, title string
		var err error

		// Try parsing first arg as scooter ID
		if len(args) >= 3 {
			if tryID, parseErr := strconv.Atoi(args[0]); parseErr == nil && tryID > 0 {
				id, err = resolveScooterID(args[:1])
				latStr = args[1]
				lngStr = args[2]
				if len(args) == 4 {
					title = args[3]
				}
			} else {
				// First arg is lat, not an ID
				id, err = resolveScooterID(nil)
				latStr = args[0]
				lngStr = args[1]
				title = args[2]
			}
		} else {
			// Exactly 2 args: lat lng
			id, err = resolveScooterID(nil)
			latStr = args[0]
			lngStr = args[1]
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid latitude: %s\n", latStr)
			os.Exit(1)
		}
		lng, err := strconv.ParseFloat(lngStr, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid longitude: %s\n", lngStr)
			os.Exit(1)
		}

		resp, err := apiClient.SetDestination(id, lat, lng, title)
		output.PrintCommandResponse(resp, err)
	},
}

var navigateShowCmd = &cobra.Command{
	Use:   "show [scooter-id]",
	Short: "Show current navigation destination",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := resolveScooterID(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		dest, err := apiClient.GetDestination(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		output.PrintDestination(dest)
	},
}

var navigateClearCmd = &cobra.Command{
	Use:   "clear [scooter-id]",
	Short: "Clear navigation destination",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := resolveScooterID(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := apiClient.ClearDestination(id); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("OK")
	},
}

func init() {
	navigateCmd.AddCommand(navigateShowCmd)
	navigateCmd.AddCommand(navigateClearCmd)
	rootCmd.AddCommand(navigateCmd)
}
