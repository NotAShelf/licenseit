package license

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Replaces placeholders in the template with author and date.
func GenerateLicense(template, author, date string) string {
	template = strings.ReplaceAll(template, "{author}", author)
	template = strings.ReplaceAll(template, "{date}", date)
	return template
}

// Writes the license content to a file. If the target file already exists, then
// we show a confirmation prompt.
func SaveLicense(licenseContent, licenseName, licenseDir string) error {
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
