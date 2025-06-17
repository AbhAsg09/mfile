/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"mfile/utils"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Will be used to retrieve file from the server",
	Run: func(cmd *cobra.Command, args []string) {
		if server != "" {
			fmt.Println("Server name is required.")
			return
		}
		if fileName == "" {
			fmt.Println("File name is required")
			return
		}
		username = viper.GetString("username")
		password = viper.GetString("password")
		if filePath == "" {
			filePath = viper.GetString("defaultPath")
		}

		conn, sftpClient, err := utils.CreateSFTP(server, username, password)
		if err != nil || sftpClient == nil {
			fmt.Printf("Error connecting to server: %s\n", server)
			return
		}

		file := fmt.Sprintf("%s/%s", filePath, fileName)

		srcFile, err := sftpClient.Open(file)
		if err != nil {
			panic("Failed to open remote file: " + err.Error())
		}
		defer srcFile.Close()

		dstFile, err := os.Create(fileName)
		if err != nil {
			panic("Failed to create local file: " + err.Error())
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			panic("Failed to copy: " + err.Error())
		}

		fmt.Println("File downloaded successfully!")

		if sftpClient != nil || conn != nil {
			utils.CloseSFTP(sftpClient, conn)
		}

	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVar(&server, "server", "s", "Server name that stores the file")
	getCmd.Flags().StringVar(&fileName, "file", "f", "Filename to be retrieved")
	getCmd.Flags().StringVar(&filePath, "path", "p", "File path given")
}
