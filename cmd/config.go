/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set config values using flags",
	Run: func(cmd *cobra.Command, args []string) {
		if err := viper.ReadInConfig(); err != nil {
			// Handle missing config gracefully
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Try to write a new one
				err = viper.SafeWriteConfig()
				if err != nil {
					fmt.Printf("Error creating config: %v\n", err)
					return
				}
			} else {
				fmt.Printf("Error reading config: %v\n", err)
				return
			}
		}
		fmt.Println("Do you wanna update username?(Y/N)")
		fmt.Scanln(&ans)
		if ans == "Y" || ans == "y" {
			fmt.Println("Enter username:")
			fmt.Scanln(&username)
			viper.Set("username", username)
			if err := viper.WriteConfig(); err != nil {
				fmt.Printf("Error updating username:%v\n", err)
				return
			}
			fmt.Printf("Username updated to: %s\n", username)
		} else if ans != "N" && ans != "n" {
			fmt.Println("Invalid input! Exiting...")
			return
		}

		fmt.Println("Do you wanna update password?(Y/N)")
		fmt.Scan(&ans)
		if ans == "Y" || ans == "y" {
			password, err := promptPassword()
			if err != nil {
				fmt.Printf("Error reading password: %v\n", err)
				return
			}
			viper.Set("password", password)
			if err := viper.WriteConfig(); err != nil {
				fmt.Printf("Error updating password: %v", err)
				return
			}
			fmt.Println("Password updated successfully.")
		} else if ans != "N" && ans != "n" {
			fmt.Println("Invalid input! Exiting...")
			return
		}

	},
}

func promptPassword() (string, error) {
	fmt.Print("Enter password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Move to the next line after input
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View and Edit the mfile configuration",
	Run: func(cmd *cobra.Command, args []string) {

		if addPath != "" && dropPath != "" {
			fmt.Println("Error: Cannot use --addPath and --dropPath together!")
			return
		}

		//To add a new path in the config file
		if addPath != "" {
			if err := viper.ReadInConfig(); err != nil {
				// Handle missing config gracefully
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					// Try to write a new one
					err = viper.SafeWriteConfig()
					if err != nil {
						fmt.Printf("Error creating config: %v\n", err)
						return
					}
				} else {
					fmt.Printf("Error reading config: %v\n", err)
					return
				}
			}

			// Now safe to set values
			parts := strings.SplitN(addPath, "=", 2)
			if len(parts) != 2 {
				fmt.Println("Invalid format for --addPath. Use key=value format.")
				return
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			viper.Set(key, value)

			// Overwrite with updated values
			if err := viper.WriteConfig(); err != nil {
				fmt.Printf("Error writing config: %v\n", err)
				return
			}
			fmt.Printf("Added %s=%s to the configuration\n", key, value)
		} else {
			allSettings := viper.AllSettings()
			if len(allSettings) == 0 {
				fmt.Println("No configuration settings found.")
				return
			}
			fmt.Println("Current mfile configuration settings:")
			for key, value := range allSettings {
				fmt.Printf("%s: %v\n", key, value)
			}

		}

		//To drop an existing path in the config file
		if dropPath != "" {
			settings := viper.AllSettings()

			delete(settings, dropPath)
			configFile := viper.ConfigFileUsed()

			if configFile == "" {
				fmt.Println("No config file found!")
				return
			}

			yamlData, err := yaml.Marshal(settings)
			if err != nil {
				fmt.Printf("Error marshaling config: %v\n", err)
				return
			}

			err = os.WriteFile(configFile, yamlData, 0644)
			if err != nil {
				fmt.Printf("Error writing config file: %v\n", err)
				return
			}

			fmt.Printf("Removed config key: %s\n", dropPath)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.Flags().StringVar(&addPath, "addPath", "", "Add a path in config file")
	configCmd.Flags().StringVar(&dropPath, "dropPath", "", "Drop a path from config file")
}
