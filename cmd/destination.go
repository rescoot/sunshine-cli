package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rescoot/sunshine-cli/internal/output"
	"github.com/spf13/cobra"
)

var destinationCmd = &cobra.Command{
	Use:   "destination",
	Short: "Manage navigation destination",
}

var destinationShowCmd = &cobra.Command{
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

var destAddress string

var destinationSetCmd = &cobra.Command{
	Use:   "set [scooter-id] <lat> <lng>",
	Short: "Set navigation destination",
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		var id int
		var latStr, lngStr string
		var err error

		if len(args) == 3 {
			id, err = resolveScooterID(args[:1])
			latStr = args[1]
			lngStr = args[2]
		} else {
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

		resp, err := apiClient.SetDestination(id, lat, lng, destAddress)
		output.PrintCommandResponse(resp, err)
	},
}

var destinationClearCmd = &cobra.Command{
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
	destinationSetCmd.Flags().StringVar(&destAddress, "address", "", "Address label for the destination")
	destinationCmd.AddCommand(destinationShowCmd)
	destinationCmd.AddCommand(destinationSetCmd)
	destinationCmd.AddCommand(destinationClearCmd)
	rootCmd.AddCommand(destinationCmd)
}
