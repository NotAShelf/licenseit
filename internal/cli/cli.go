package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type CLIOptions struct {
	Author      string
	ConfigFile  string
	LicenseFile string
	LicenseDir  string
	LicenseName string
	ShowHelp    bool
	ShowPreview bool
}

func ParseCLIArgs() *CLIOptions {
	opts := &CLIOptions{}
	var rootCmd = &cobra.Command{
		Use:   filepath.Base(os.Args[0]) + " <template-name> [options]",
		Short: "Generate a license file based on a template and configuration.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				opts.LicenseName = args[0]
			}
		},
	}

	rootCmd.Flags().StringVarP(&opts.Author, "author", "a", "", "Author's name for the license")
	rootCmd.Flags().StringVarP(&opts.ConfigFile, "config", "c", "", "Path to the configuration file")
	rootCmd.Flags().StringVarP(&opts.LicenseFile, "file", "f", "", "Name of the generated license file")
	rootCmd.Flags().StringVarP(&opts.LicenseDir, "dir", "d", ".", "Directory to save generated licenses")

	var previewCmd = &cobra.Command{
		Use:   "preview",
		Short: "Show available license templates",
		Run: func(cmd *cobra.Command, args []string) {
			opts.ShowPreview = true
		},
	}

	rootCmd.AddCommand(previewCmd)

	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		opts.ShowHelp = true
	})

	// Parse the CLI
	_ = rootCmd.Execute()

	// If no args and no subcommand, show help
	if opts.LicenseName == "" && !opts.ShowPreview && !opts.ShowHelp {
		opts.ShowHelp = true
	}

	return opts
}

func PromptForAuthor() string {
	fmt.Print("Author not provided. Please enter the author's name: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
	return strings.TrimSpace(input)
}

func PrintHelp() {
	programName := filepath.Base(os.Args[0])
	fmt.Printf("Usage: %s <template-name> [options]\n", programName)
	fmt.Println("Generate a license file based on a template and configuration.")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  preview              Show available license templates")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  <template-name>      The base name of the license template to use (e.g., 'MIT')")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -a, --author <name>  The author of the license")
	fmt.Println("  -c, --config <file>  Path to the configuration file")
	fmt.Println("  -f, --file <name>    Name of the generated license file")
	fmt.Println("  -d, --dir <dir>      Directory to save generated licenses")
	fmt.Println("  -h, --help           Show this help message and exit.")
}
