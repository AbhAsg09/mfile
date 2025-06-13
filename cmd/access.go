/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// accessCmd represents the access command
var accessCmd = &cobra.Command{
	Use:   "access",
	Short: "To check access to a particular server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("access called")
	},
}

func init() {
	rootCmd.AddCommand(accessCmd)
}
