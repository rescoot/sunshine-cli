package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rescoot/sunshine-cli/internal/output"
	"github.com/spf13/cobra"
)

var navigateCmd = &cobra.Command{
	Use:   "navigate <lat> <lng> [title]",
	Short: "Set navigation destination",
	Long:  "Set a navigation destination by providing lat/lng coordinates and an optional title.\nUse 'navigate show' to view or 'navigate clear' to remove the current destination.",
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := resolveScooterID(nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		lat, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid latitude: %s\n", args[0])
			os.Exit(1)
		}
		lng, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid longitude: %s\n", args[1])
			os.Exit(1)
		}

		title := ""
		if len(args) == 3 {
			title = args[2]
		}

		resp, err := apiClient.SetDestination(id, lat, lng, title)
		output.PrintCommandResponse(resp, err)
	},
}

var navigateShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current navigation destination",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := resolveScooterID(nil)
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
	Use:   "clear",
	Short: "Clear navigation destination",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := resolveScooterID(nil)
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
