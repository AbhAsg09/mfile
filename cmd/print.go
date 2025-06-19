package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "print",
	Short: "Print \"Hello World!\"",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World!")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
