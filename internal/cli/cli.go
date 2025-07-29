package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	if len(os.Args) < 2 || os.Args[1] == "-help" || os.Args[1] == "--help" {
		opts.ShowHelp = true
		return opts
	}

	if os.Args[1] == "preview" {
		opts.ShowPreview = true
		return opts
	}

	opts.LicenseName = os.Args[1]
	newArgs := []string{os.Args[0]}
	newArgs = append(newArgs, os.Args[2:]...)
	os.Args = newArgs

	authorFlag := flag.String("author", "", "Author's name for the license")
	configFileFlag := flag.String("config", "", "Path to the configuration file")
	licenseFileFlag := flag.String("file", "", "Name of the generated license file")
	licenseDirFlag := flag.String("dir", ".", "Directory to save generated licenses")
	flag.Parse()

	opts.Author = *authorFlag
	opts.ConfigFile = *configFileFlag
	opts.LicenseFile = *licenseFileFlag
	opts.LicenseDir = *licenseDirFlag

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
	fmt.Println("  -author <name>       The author of the license")
	fmt.Println("  -config <file>       Path to the configuration file")
	fmt.Println("  -file <filename>     Name of the generated license file")
	fmt.Println("  -dir <directory>     Directory to save generated licenses")
	fmt.Println("  -help                Show this help message and exit.")
}
