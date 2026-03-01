package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var apiMethod string
var apiData string

var apiCmd = &cobra.Command{
	Use:   "api <path> [-d data]",
	Short: "Make an authenticated API request",
	Long: `Make an authenticated request to the Sunshine API.

The path is relative to /api/v1/. Method defaults to GET, or POST if -d is given.

Examples:
  sunshine api /scooters
  sunshine api /scooters/3
  sunshine api /scooters/3/lock -X POST
  sunshine api /scooters/3/blinkers -X POST -d '{"state":"left"}'
  sunshine api /scooters/3/trips?limit=5`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		method := strings.ToUpper(apiMethod)
		if method == "" {
			if apiData != "" {
				method = "POST"
			} else {
				method = "GET"
			}
		}

		var body interface{}
		if apiData != "" {
			var parsed interface{}
			if err := json.Unmarshal([]byte(apiData), &parsed); err != nil {
				fmt.Fprintf(os.Stderr, "Error: invalid JSON data: %v\n", err)
				os.Exit(1)
			}
			body = parsed
		}

		respBody, statusCode, err := apiClient.DoRaw(method, path, body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Pretty-print JSON if possible
		var parsed interface{}
		if json.Unmarshal(respBody, &parsed) == nil {
			formatted, _ := json.MarshalIndent(parsed, "", "  ")
			fmt.Println(string(formatted))
		} else {
			fmt.Println(string(respBody))
		}

		if statusCode >= 400 {
			os.Exit(1)
		}
	},
}

func init() {
	apiCmd.Flags().StringVarP(&apiMethod, "method", "X", "", "HTTP method (default: GET, or POST with -d)")
	apiCmd.Flags().StringVarP(&apiData, "data", "d", "", "JSON request body")
	rootCmd.AddCommand(apiCmd)
}
