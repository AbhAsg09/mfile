/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mfile",
	Short: "Test",
	Long:  `Test`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Set default config path and name
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("mfile") // ~/mfile.yaml
	}

	// Automatically read environment variables that match
	viper.AutomaticEnv()

	// If config file doesn't exist, create it
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Create empty config file
			home, _ := os.UserHomeDir()
			configPath := filepath.Join(home, "mfile.yaml")
			err := os.WriteFile(configPath, []byte{}, 0644)
			if err != nil {
				fmt.Printf("Failed to create config file: %v\n", err)
				os.Exit(1)
			}
			viper.SetConfigFile(configPath)
		} else {
			fmt.Printf("Failed to read config file: %v\n", err)
			os.Exit(1)
		}
	}
}
