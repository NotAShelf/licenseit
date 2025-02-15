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
	content, err := os.ReadFile(configFile)
	if err != nil {
		return "", err
	}
	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return "", err
	}
	return config.Author, nil
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

func main() {
	if len(os.Args) < 2 || strings.HasPrefix(os.Args[1], "-") {
		fmt.Println("Error: License name (template) must be specified as the first argument.")
		os.Exit(1)
	}

	licenseBaseName := os.Args[1]
	newArgs := []string{os.Args[0]}
	newArgs = append(newArgs, os.Args[2:]...)
	os.Args = newArgs

	authorFlag := flag.String("author", "", "Author's name for the license")
	configFileFlag := flag.String("config", "config.json", "Path to the configuration file (optional)")
	licenseDirFlag := flag.String("dir", ".", "Directory to save generated licenses")
	flag.Parse()

	currentDate := time.Now().Format("2006-01-02")

	// Get the name of the author from the -author flag if provided. Try the config file
	// otherwise. If both are missing, exit with an error.
	author := *authorFlag
	if author == "" {
		authorFromConfig, err := readConfig(*configFileFlag)
		if err != nil || authorFromConfig == "" {
			fmt.Println("Error: Author must be specified via flag or config file.")
			flag.Usage()
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

	licenseContent := generateLicense(template, author, currentDate)

	err = saveLicense(licenseContent, licenseBaseName+filepath.Ext(templateName), *licenseDirFlag)
	if err != nil {
		fmt.Println("Error saving license:", err)
		os.Exit(1)
	}

	fmt.Printf("License '%s' created and saved for author: %s\n", licenseBaseName, author)
}
