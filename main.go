package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Author string `json:"author"`
}

//go:embed templates/*
var licenseTemplates embed.FS

func readConfig(configFile string) (string, error) {
	// If -config flag is specified, use that path
	if configFile != "" {
		content, err := os.ReadFile(configFile)
		if err != nil {
			return "", fmt.Errorf("could not read config file '%s': %w", configFile, err)
		}

		var config Config

		err = json.Unmarshal(content, &config)
		if err != nil {
			return "", fmt.Errorf("could not parse config file '%s': %w", configFile, err)
		}

		return config.Author, nil
	}

	// Otherwise, check $XDG_CONFIG_HOME
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		configFile = filepath.Join(xdgConfigHome, "licenseit", "config.json")
		content, err := os.ReadFile(configFile)
		if err == nil {
			var config Config
			err = json.Unmarshal(content, &config)
			if err == nil {
				return config.Author, nil
			}
		}
	}

	return "", nil
}

func findTemplate(templateBaseName string) (string, error) {
	extensions := []string{".txt", ".md"}

	for _, ext := range extensions {
		templateName := templateBaseName + ext
		_, err := licenseTemplates.ReadFile("templates/" + templateName)
		if err == nil {
			return templateName, nil
		}
	}

	_, err := licenseTemplates.ReadFile("templates/" + templateBaseName)
	if err == nil {
		return templateBaseName, nil
	}

	return "", fmt.Errorf("could not find a template for '%s' with supported extensions", templateBaseName)
}

func loadTemplate(templateName string) (string, error) {
	// Read the template from the embedded file system
	content, err := licenseTemplates.ReadFile("templates/" + templateName)
	if err != nil {
		return "", fmt.Errorf("could not read template file '%s': %w", templateName, err)
	}

	return string(content), nil
}

func generateLicense(template string, author string, date string) string {
	template = strings.ReplaceAll(template, "{author}", author)
	template = strings.ReplaceAll(template, "{date}", date)
	return template
}

func saveLicense(licenseContent string, licenseName string, licenseDir string) error {
	err := os.MkdirAll(licenseDir, 0755)
	if err != nil {
		return fmt.Errorf("could not create directory '%s': %w", licenseDir, err)
	}

	filePath := filepath.Join(licenseDir, licenseName)

	// If file exists, prompt for overwrite
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("File '%s' already exists. Overwrite? (y/N): ", filePath)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("could not read user input: %w", err)
		}
		response = strings.TrimSpace(response)
		if strings.ToLower(response) != "y" {
			return fmt.Errorf("operation aborted by user; file exists")
		}
	}

	err = os.WriteFile(filePath, []byte(licenseContent), 0644)
	if err != nil {
		return fmt.Errorf("could not write license to file '%s': %w", filePath, err)
	}

	return nil
}

func listTemplates() error {
	entries, err := licenseTemplates.ReadDir("templates")
	if err != nil {
		return fmt.Errorf("could not read templates directory: %w", err)
	}

	fmt.Println("Available license templates:")
	for _, entry := range entries {
		fmt.Println(" -", entry.Name())
	}

	return nil
}

func printHelp() {
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

func main() {
	// Show help if no arguments provided
	if len(os.Args) < 2 || os.Args[1] == "-help" || os.Args[1] == "--help" {
		printHelp()
		os.Exit(0)
	}

	if os.Args[1] == "preview" {
		err := listTemplates()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	licenseBaseName := os.Args[1]
	newArgs := []string{os.Args[0]}
	newArgs = append(newArgs, os.Args[2:]...)
	os.Args = newArgs

	authorFlag := flag.String("author", "", "Author's name for the license")
	configFileFlag := flag.String("config", "", "Path to the configuration file")
	licenseFileFlag := flag.String("file", "", "Name of the generated license file")
	licenseDirFlag := flag.String("dir", ".", "Directory to save generated licenses")
	flag.Parse()

	author := *authorFlag
	if author == "" {
		authorFromConfig, err := readConfig(*configFileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}
		if authorFromConfig != "" {
			author = authorFromConfig
		}
	}

	if author == "" {
		fmt.Print("Author not provided. Please enter the author's name: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
		author = strings.TrimSpace(input)
	}

	if author == "" {
		fmt.Fprintf(os.Stderr, "Error: Author is required.\n")
		printHelp()
		os.Exit(1)
	}

	templateName, err := findTemplate(licenseBaseName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading template: %v\n", err)
		os.Exit(1)
	}

	template, err := loadTemplate(templateName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading template: %v\n", err)
		os.Exit(1)
	}

	currentYear := time.Now().Year()
	licenseContent := generateLicense(template, author, fmt.Sprintf("%d", currentYear))

	// Determine the generated license file name from --file flag or fall back to template name
	licenseFileName := *licenseFileFlag
	if licenseFileName == "" {
		licenseFileName = licenseBaseName + filepath.Ext(templateName)
	}

	err = saveLicense(licenseContent, licenseFileName, *licenseDirFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving license: %v\n", err)
		os.Exit(1)
	}

	// Compute absolute path for the output file
	absPath, err := filepath.Abs(filepath.Join(*licenseDirFlag, licenseFileName))
	if err != nil {
		absPath = filepath.Join(*licenseDirFlag, licenseFileName)
	}
	fmt.Printf("License '%s' created at '%s' for author: %s\n", licenseFileName, absPath, author)
}
