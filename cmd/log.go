package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "To get server logs",
	Run: func(cmd *cobra.Command, args []string) {
		if server == "" {
			fmt.Println("Server name is required.")
			return
		}

		if fileName == "" {
			fmt.Println("File name is required")
			return
		}

		fmt.Println("Starting download...")

		sname := tryResolveServerNames(server)
		if server == "" {
			fmt.Println("Server name either incorrect or Server Down! Please recheck and try later.")
			return
		}

		username = viper.GetString("username")
		password = viper.GetString("password")
		filePath = viper.GetString("logs")
		getLogs(sname, server, username, password)

	},
}

func getLogs(sname []string, server string, username string, password string) {
	for _, i := range sname {
		fmt.Println(i)
	}
}

func init() {
	getCmd.Flags().StringVarP(&server, "server", "s", "", "Server name that stores the file")
	getCmd.Flags().StringVarP(&fileName, "file", "f", "", "Filename to be retrieved")
	getCmd.Flags().StringVarP(&search, "search", "S", "", "String to be searched in file")
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
