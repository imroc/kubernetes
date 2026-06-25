package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// getSchemasDirectory resolves the output directory for schema files.
// Priority: flag arg > OUTPUT_DIR env > ./schemas
func getSchemasDirectory(flagValue string) string {
	outputDir := flagValue
	if outputDir == "" {
		outputDir = os.Getenv("OUTPUT_DIR")
	}
	if outputDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		outputDir = filepath.Join(cwd, "schemas")
	}
	if strings.Contains(outputDir, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		outputDir = strings.Replace(outputDir, "~", homeDir, 1)
	}
	fileInfo, err := os.Stat(outputDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	} else {
		if !fileInfo.IsDir() {
			log.Fatalf("%s is not a directory", outputDir)
		}
	}
	outputDir, err = filepath.EvalSymlinks(outputDir)
	if err != nil {
		log.Fatal(err)
	}
	return outputDir
}
