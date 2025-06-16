/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"mfile/utils"
	"net"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// accessCmd represents the access command
var accessCmd = &cobra.Command{
	Use:   "access",
	Short: "To check access to a particular server",
	Run: func(cmd *cobra.Command, args []string) {

		if server == "" {
			fmt.Println("Server is required!")
		}

		if username == "" && password == "" {
			username = viper.GetString("username")
			password = viper.GetString("password")
		} else if username == "" && password != "" {
			fmt.Println("You need to enter username!")
		} else if username != "" && password == "" {
			fmt.Println("You need to enter password!")
		}

		ips, err := net.LookupIP(server)
		if err != nil {
			fmt.Printf("Error encountered:%v\n", err)
		}
		if len(ips) == 0 {
			fmt.Println("Server could not be resolved. Check again!")
		} else {
			for _, ip := range ips {
				ipSSH := fmt.Sprintf("%s:22", ip.String())
				conn, session, err := utils.CreateConnection(ipSSH, username, password)
				if session == nil && err != nil {
					fmt.Printf("Error encountered: %v\n", err)
				}
				if session != nil && err == nil {
					fmt.Printf("%s has access to the server %s", username, server)
					utils.Close(session, conn)
					break
				}
			}
		}
	},
}

// var sudoCheck = &cobra.Command{
// 	Use:   "Check Sudo",
// 	Short: "Check sudo access to a server",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		if server == "" {
// 			fmt.Println("Server is required!")
// 		}

// 		if username == "" && password == "" {
// 			username = viper.GetString("username")
// 			password = viper.GetString("password")
// 		} else if username == "" && password != "" {
// 			fmt.Println("You need to enter username!")
// 		} else if username != "" && password == "" {
// 			fmt.Println("You need to enter password!")
// 		}

// 		ips, err := net.LookupIP(server)
// 		if err != nil {
// 			fmt.Printf("Error encountered:%v\n", err)
// 		}
// 		if len(ips) == 0 {
// 			fmt.Println("Server could not be resolved. Check again!")
// 		} else {
// 			for _, ip := range ips {
// 				ipSSH := fmt.Sprintf("%s:22", ip.String())
// 				conn, session, err := utils.CreateConnection(ipSSH, username, password)
// 				if session == nil && err != nil {
// 					fmt.Printf("Error encountered: %v\n", err)
// 				}
// 				if session != nil && err == nil {
// 					sudo := utils.SudoAccess(session, password)
// 					if sudo {
// 					} else {
// 						fmt.Printf("%s does not have sudo access to the server %s", username, server)
// 					}
// 					utils.Close(session, conn)
// 					break
// 				}
// 			}
// 		}
// 	},
// }

func init() {
	rootCmd.AddCommand(accessCmd)
	accessCmd.Flags().StringVar(&username, "username", "u", "User inputs username")
	accessCmd.Flags().StringVar(&server, "server", "s", "User inputs username")
	accessCmd.Flags().StringVar(&password, "password", "p", "User inputs username")

	// rootCmd.AddCommand(sudoCheck)
	// sudoCheck.Flags().StringVar(&username, "username", "u", "User inputs username")
	// sudoCheck.Flags().StringVar(&server, "server", "s", "User inputs username")
	// sudoCheck.Flags().StringVar(&password, "password", "p", "User inputs username")
}
