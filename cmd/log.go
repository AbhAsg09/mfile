package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("log called")
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}

func tryResolveServerNames(base string) []string {
	var valid []string

	if _, err := net.LookupIP(base); err == nil {
		valid = append(valid, base)
		return valid
	}

	suffixes := []string{"a", "b", "c", "d"}
	for _, suffix := range suffixes {
		host := fmt.Sprintf("%s%s", base, suffix)
		if _, err := net.LookupIP(host); err == nil {
			valid = append(valid, host)
		}
	}
	return valid
}
