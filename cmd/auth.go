package cmd

import (
	"fmt"
	"os"

	"github.com/rescoot/sunshine-cli/internal/auth"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Sunshine via OAuth2",
	Run: func(cmd *cobra.Command, args []string) {
		tokens, err := auth.Login(cfg.Server, cfg.ClientID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Login failed: %v\n", err)
			os.Exit(1)
		}

		if err := auth.SaveTokens(tokens); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving tokens: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Authenticated successfully.")
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear stored authentication tokens",
	Run: func(cmd *cobra.Command, args []string) {
		if err := auth.ClearTokens(); err != nil {
			fmt.Fprintf(os.Stderr, "Error clearing tokens: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Logged out.")
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Run: func(cmd *cobra.Command, args []string) {
		tokens, err := auth.LoadTokens()
		if err != nil || tokens == nil {
			fmt.Println("Not authenticated. Run 'sunshine auth login' to authenticate.")
			return
		}

		fmt.Printf("Server:  %s\n", cfg.Server)
		fmt.Printf("Scopes:  %s\n", tokens.Scopes)
		if tokens.IsExpired() {
			fmt.Println("Status:  expired (will refresh on next request)")
		} else if tokens.ExpiresAt.IsZero() {
			fmt.Println("Status:  authenticated (no expiry)")
		} else {
			fmt.Printf("Status:  authenticated (expires %s)\n", tokens.ExpiresAt.Local().Format("2006-01-02 15:04"))
		}
	},
}

func init() {
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authStatusCmd)
	rootCmd.AddCommand(authCmd)
}
