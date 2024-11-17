package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "statechart",
	Short: "StateChart generates Mermaid.js diagrams from Go state machines",
	Long: `StateChart is a tool to scan Go projects for *_state_machine.go files
and generate Mermaid.js state diagrams for visualizing state transitions.`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
