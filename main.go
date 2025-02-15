package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type License struct {
	Author       string
	Date         string
	TemplateName string
	Content      string
}

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
	// Read the template from the embedded files system.
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
	err = os.WriteFile(filePath, []byte(licenseContent), 0644)
	if err != nil {
		return fmt.Errorf("could not write license to file '%s': %w", filePath, err)
	}

	return nil
}

func printHelp() {
	programName := filepath.Base(os.Args[0])
	fmt.Printf("Usage: %s <template-name> [options]\n", programName)
	fmt.Println("Generate a license file based on a template and configuration.")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  <template-name>  The base name of the license template to use (e.g., 'MIT')")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -author <name>        The author of the license (if not specified, it will be read from the config file)")
	fmt.Println("  -config <file>        Path to the configuration file (optional)")
	fmt.Println("  -file <filename>      Name of the generated license file")
	fmt.Println("  -dir <directory>      Directory to save generated licenses")
	fmt.Println("  -help                 Show this help message and exit.")
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "-help" || os.Args[1] == "--help" {
		printHelp()
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

	// Get the name of the author from the -author flag if provided. Try the config file
	// otherwise. If both are missing, exit with an error.
	author := *authorFlag
	if author == "" {
		authorFromConfig, err := readConfig(*configFileFlag)
		if err != nil {
			fmt.Println("Error:", err)
			printHelp()
			os.Exit(1)
		}

		// If config didn't provide an author, exit with an error
		if authorFromConfig == "" {
			fmt.Println("Error: Author must be specified via flag or config file.")
			printHelp()
			os.Exit(1)
		}
		author = authorFromConfig
	}

	templateName, err := findTemplate(licenseBaseName)
	if err != nil {
		fmt.Println("Error loading template:", err)
		os.Exit(1)
	}

	template, err := loadTemplate(templateName)
	if err != nil {
		fmt.Println("Error loading template:", err)
		os.Exit(1)
	}

	currentYear := time.Now().Year()
	licenseContent := generateLicense(template, author, fmt.Sprintf("%d", currentYear))

	// Determine the generated license file name from --file flag or fall back to template name
	licenseFileName := *licenseFileFlag
	if licenseFileName == "" {
		licenseFileName = licenseBaseName + filepath.Ext(templateName)
	}

	// Save the generated license to the specified directory
	err = saveLicense(licenseContent, licenseFileName, *licenseDirFlag)
	if err != nil {
		fmt.Println("Error saving license:", err)
		os.Exit(1)
	}

	// Confirm the license file creation
	fmt.Printf("License '%s' created and saved for author: %s\n", licenseFileName, author)
}
