package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chalfel/statechart/internal"
)

func main() {
	// Define CLI flags
	projectDir := flag.String("project", ".", "Path to the root directory of the project to scan for _state_machine.go files")
	outputDir := flag.String("output", "./charts", "Directory to save generated Mermaid charts")
	flag.Parse()

	// Scan the entire project for _state_machine.go files
	files := []string{}
	err := filepath.Walk(*projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".go" && filepath.Base(path[len(path)-16:]) == "_state_machine.go" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error scanning project: %v", err)
	}

	// Ensure output directory exists
	if _, err := os.Stat(*outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(*outputDir, 0755); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}
	}

	// Process each file and generate Mermaid charts
	for _, file := range files {
		diagram, err := internal.GenerateMermaidFromFile(file)
		if err != nil {
			log.Printf("Error generating Mermaid chart for %s: %v\n", file, err)
			continue
		}

		outputFile := filepath.Join(*outputDir, filepath.Base(file)+".mmd")
		err = os.WriteFile(outputFile, []byte(diagram), 0644)
		if err != nil {
			log.Printf("Error writing Mermaid chart for %s: %v\n", file, err)
			continue
		}

		fmt.Printf("Generated Mermaid chart for %s -> %s\n", file, outputFile)
	}
}
