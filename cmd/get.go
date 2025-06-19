/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"mfile/utils"
	"net"
	"os"
	"path"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Will be used to retrieve file from the server",
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

		server = tryResolveServerName(server)
		if server == "" {
			fmt.Println("Server name either incorrect or Server Down! Please recheck and try later.")
			return
		}

		username = viper.GetString("username")
		password = viper.GetString("password")
		if filePath == "" {
			filePath = viper.GetString("defaultPath")
		} else {
			filePath = viper.GetString(filePath)
			if filePath == "" {
				fmt.Println("File path not found")
				return
			}
		}

		success := SFTP(server, username, password, fileName, filePath)
		if success {
			fmt.Println("File downloaded successfully!")
		} else {
			fmt.Println("File could not be downloaded!")

		}
	},
}

func SFTP(server string, username string, password string, fileName string, filePath string) bool {
	conn, sftpClient, err := utils.CreateSFTP(server, username, password)
	if err != nil || sftpClient == nil {
		fmt.Printf("Error connecting to server: %s\n", server)
		return false
	}

	file := path.Join(filePath, fileName)

	srcFile, err := sftpClient.Open(file)
	fmt.Printf("Connected to server %s\n", server)
	fmt.Printf("Filepath: %s\n", file)
	if err != nil {
		fmt.Printf("Failed to find remote file: %v\n", err)
		utils.CloseSFTP(sftpClient, conn)
		return false
	}
	defer srcFile.Close()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to get your home directory")
		utils.CloseSFTP(sftpClient, conn)
		return false
	}

	downloadsDir := filepath.Join(homeDir, "Downloads", fileName)
	dstFile, err := os.Create(downloadsDir)
	if err != nil {
		fmt.Printf("Failed to create local file: %v\n", err)
		utils.CloseSFTP(sftpClient, conn)
		return false
	}
	defer dstFile.Close()

	fmt.Printf("Downloading %s from %s...\n", fileName, server)

	stat, _ := srcFile.Stat()
	// Optional: Add progress bar
	bar := progressbar.DefaultBytes(
		stat.Size(),
		fmt.Sprintf("Downloading from %s", server),
	)

	_, err = io.Copy(io.MultiWriter(dstFile, bar), srcFile)
	if err != nil {
		fmt.Printf("Failed to copy file: %v\n", err)
		return false
	}

	utils.CloseSFTP(sftpClient, conn)
	return true
}

func tryResolveServerName(base string) string {

	// Try original
	if _, err := net.LookupIP(base); err == nil {
		return base
	}

	// Try suffix variations
	suffix := "a"
	base = fmt.Sprintf("%s%s", base, suffix)
	if _, err := net.LookupIP(base); err == nil {
		return base
	}
	return ""
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&server, "server", "s", "", "Server name that stores the file")
	getCmd.Flags().StringVarP(&fileName, "file", "f", "", "Filename to be retrieved")
	getCmd.Flags().StringVarP(&filePath, "path", "p", "", "Remote file path (optional)")

}
