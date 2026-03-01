package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rescoot/sunshine-cli/internal/output"
	"github.com/spf13/cobra"
)

var tripListLimit int
var tripListOffset int

var tripsCmd = &cobra.Command{
	Use:   "trips",
	Short: "Manage trips",
}

var tripsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List trips for a scooter",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := resolveScooterID(nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		trips, err := apiClient.ListTrips(id, tripListLimit, tripListOffset)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		output.PrintTripList(trips)
	},
}

var tripsShowCmd = &cobra.Command{
	Use:   "show <trip-id>",
	Short: "Show trip details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tripID, err := strconv.Atoi(args[0])
		if err != nil || tripID <= 0 {
			fmt.Fprintf(os.Stderr, "Error: invalid trip ID %q\n", args[0])
			os.Exit(1)
		}

		scooterID, err := resolveScooterID(nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		trip, err := apiClient.GetTrip(scooterID, tripID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		output.PrintTripDetail(trip)
	},
}

func init() {
	tripsListCmd.Flags().IntVar(&tripListLimit, "limit", 20, "Maximum number of trips to return (0 for all)")
	tripsListCmd.Flags().IntVar(&tripListOffset, "offset", 0, "Number of trips to skip")
	tripsCmd.AddCommand(tripsListCmd)
	tripsCmd.AddCommand(tripsShowCmd)
	rootCmd.AddCommand(tripsCmd)
}
