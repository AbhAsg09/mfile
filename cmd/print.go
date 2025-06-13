package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print \"Hello World!\"",
	Long:  `Print \"Hello World!\`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World!")
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}
