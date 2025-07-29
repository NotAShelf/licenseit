package template

import (
	"embed"
	"fmt"
)

// Searches for a template file with supported extensions in the provided embedded filesystem.
func FindTemplate(fs embed.FS, templateBaseName string) (string, error) {
	extensions := []string{".txt", ".md"}

	for _, ext := range extensions {
		templateName := templateBaseName + ext
		_, err := fs.ReadFile("templates/" + templateName)
		if err == nil {
			return templateName, nil
		}
	}

	_, err := fs.ReadFile("templates/" + templateBaseName)
	if err == nil {
		return templateBaseName, nil
	}

	return "", fmt.Errorf("could not find a template for '%s' with supported extensions", templateBaseName)
}

// Reads the template content from the provided embedded filesystem.
func LoadTemplate(fs embed.FS, templateName string) (string, error) {
	content, err := fs.ReadFile("templates/" + templateName)
	if err != nil {
		return "", fmt.Errorf("could not read template file '%s': %w", templateName, err)
	}
	return string(content), nil
}
