package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsDir string

var docsCmd = &cobra.Command{
	Use:    "docs",
	Short:  "Generate documentation",
	Hidden: true,
}

var docsManCmd = &cobra.Command{
	Use:   "man",
	Short: "Generate man pages",
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.MkdirAll(docsDir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating dir: %v\n", err)
			os.Exit(1)
		}
		header := &doc.GenManHeader{
			Title:   "SUNSHINE",
			Section: "1",
			Source:  "sunshine " + Version,
		}
		if err := doc.GenManTree(rootCmd, header, docsDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating man pages: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Man pages written to %s/\n", docsDir)
	},
}

var docsMarkdownCmd = &cobra.Command{
	Use:   "markdown",
	Short: "Generate markdown documentation",
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.MkdirAll(docsDir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating dir: %v\n", err)
			os.Exit(1)
		}
		if err := doc.GenMarkdownTree(rootCmd, docsDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating docs: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Markdown docs written to %s/\n", docsDir)
	},
}

func init() {
	docsCmd.PersistentFlags().StringVar(&docsDir, "dir", "doc", "Output directory")
	docsCmd.AddCommand(docsManCmd)
	docsCmd.AddCommand(docsMarkdownCmd)
	rootCmd.AddCommand(docsCmd)
}
