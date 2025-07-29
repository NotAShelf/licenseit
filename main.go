package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"notashelf.dev/licenseit/internal/cli"
	"notashelf.dev/licenseit/internal/config"
	"notashelf.dev/licenseit/internal/license"
	"notashelf.dev/licenseit/internal/template"
)

//go:embed templates/*
var licenseTemplates embed.FS

func main() {
	opts := cli.ParseCLIArgs()

	if opts.ShowHelp {
		cli.PrintHelp()
		os.Exit(0)
	}

	if opts.ShowPreview {
		entries, err := licenseTemplates.ReadDir("templates")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not read templates directory: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Available license templates:")
		for _, entry := range entries {
			fmt.Println(" -", entry.Name())
		}
		os.Exit(0)
	}

	author := opts.Author
	if author == "" {
		authorFromConfig, err := config.ReadConfig(opts.ConfigFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}
		if authorFromConfig != "" {
			author = authorFromConfig
		}
	}

	if author == "" {
		author = cli.PromptForAuthor()
	}

	if author == "" {
		fmt.Fprintf(os.Stderr, "Error: Author is required.\n")
		cli.PrintHelp()
		os.Exit(1)
	}

	templateName, err := template.FindTemplate(licenseTemplates, opts.LicenseName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading template: %v\n", err)
		os.Exit(1)
	}

	templateContent, err := template.LoadTemplate(licenseTemplates, templateName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading template: %v\n", err)
		os.Exit(1)
	}

	currentYear := time.Now().Year()
	licenseContent := license.GenerateLicense(templateContent, author, fmt.Sprintf("%d", currentYear))

	licenseFileName := opts.LicenseFile
	if licenseFileName == "" {
		licenseFileName = opts.LicenseName + filepath.Ext(templateName)
	}

	err = license.SaveLicense(licenseContent, licenseFileName, opts.LicenseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving license: %v\n", err)
		os.Exit(1)
	}

	absPath, err := filepath.Abs(filepath.Join(opts.LicenseDir, licenseFileName))
	if err != nil {
		absPath = filepath.Join(opts.LicenseDir, licenseFileName)
	}
	fmt.Printf("License '%s' created at '%s' for author: %s\n", licenseFileName, absPath, author)
}
