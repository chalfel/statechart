package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/chalfel/statechart/internal"
	"github.com/spf13/cobra"
)

var file string
var output string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a Mermaid.js diagram for a single state machine file",
	Run: func(cmd *cobra.Command, args []string) {
		if file == "" {
			log.Fatal("Error: --file flag is required")
		}

		diagram, err := internal.GenerateMermaidFromFile(file)
		if err != nil {
			log.Fatalf("Error generating Mermaid diagram: %v", err)
		}

		if output != "" {
			err = os.WriteFile(output, []byte(diagram), 0644)
			if err != nil {
				log.Fatalf("Error writing Mermaid diagram to file: %v", err)
			}
			fmt.Printf("Mermaid diagram saved to %s\n", output)
		} else {
			fmt.Println("Generated Mermaid diagram:")
			fmt.Println(diagram)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&file, "file", "f", "", "Path to the Go state machine file")
	generateCmd.Flags().StringVarP(&output, "output", "o", "", "Path to save the generated Mermaid diagram")
}
