package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "mfile",
	Short: "Test",
	Long:  `Test`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
}

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
		viper.SetConfigFile(cfgFile)
	} else {
		cwd, err := os.Getwd()
		cobra.CheckErr(err)

		configPath := filepath.Join(cwd, "mfile.yaml")
		viper.SetConfigFile(configPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("mfile")
		viper.AddConfigPath(cwd)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			configFile := viper.ConfigFileUsed()
			if configFile == "" {
				configFile = "mfile.yaml"
			}
			err := os.WriteFile(configFile, []byte{}, 0644)
			if err != nil {
				fmt.Printf("Failed to create config file: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Failed to read config file: %v\n", err)
			os.Exit(1)
		}
	}
}
