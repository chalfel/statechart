package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chalfel/statechart/internal"
	"github.com/spf13/cobra"
)

var projectDir string
var outputDir string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a project directory for state machine files and generate Mermaid.js diagrams",
	Run: func(cmd *cobra.Command, args []string) {
		if projectDir == "" {
			log.Fatal("Error: --project flag is required")
		}

		// Scan for files
		files := []string{}
		err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
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
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				log.Fatalf("Failed to create output directory: %v", err)
			}
		}

		// Generate Mermaid diagrams for each file
		for _, file := range files {
			diagram, err := internal.GenerateMermaidFromFile(file)
			if err != nil {
				log.Printf("Error generating Mermaid chart for %s: %v\n", file, err)
				continue
			}

			outputFile := filepath.Join(outputDir, filepath.Base(file)+".mmd")
			err = os.WriteFile(outputFile, []byte(diagram), 0644)
			if err != nil {
				log.Printf("Error writing Mermaid chart for %s: %v\n", file, err)
				continue
			}

			fmt.Printf("Generated Mermaid chart for %s -> %s\n", file, outputFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&projectDir, "project", "p", ".", "Root directory of the project to scan")
	scanCmd.Flags().StringVarP(&outputDir, "output", "o", "./charts", "Directory to save generated Mermaid charts")
}
