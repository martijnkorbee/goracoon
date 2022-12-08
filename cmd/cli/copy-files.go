package main

import (
	"embed"
	"os"

	"github.com/fatih/color"
)

//go:embed templates
var templateFS embed.FS

func copyFileFromTemplate(templatePath string, targetFile string) error {
	// check if file exist
	if fileExists(targetFile) {
		color.Yellow("Warning: skipped copying: " + targetFile + " already exists")
		return nil
	}

	data, err := templateFS.ReadFile(templatePath)
	if err != nil {
		exitGracefully(err)
	}

	err = os.WriteFile(targetFile, []byte(data), 0644)
	if err != nil {
		exitGracefully(err)
	}

	return nil
}

func fileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}

	return true
}
