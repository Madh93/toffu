package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:     "toffu",
	Version: "0.2.0",
	Short:   "Toffu - Woffu presence from your terminal",
	Long:    "Toffu is a simple CLI to manage presence in Woffu from your terminal",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global Flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file for toffu")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable debug mode")
	rootCmd.PersistentFlags().VisitAll(bindCustomFlag)

	// Local Flags
	rootCmd.Flags().VisitAll(bindCustomFlag)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Config name
		viper.SetConfigName("toffu.json")
		viper.SetConfigType("json")

		// Default config file
		defaultConfigPath := filepath.Join(os.Getenv("HOME"), ".toffu")
		defaultConfigFile := "toffu.json"
		err := ensureConfigFileExists(defaultConfigPath, defaultConfigFile)
		if err != nil {
			fmt.Println("error creating default config file,", err)
			return
		}

		// Config location
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".toffu"))
	}

	// Environment variables
	viper.SetEnvPrefix("TOFFU")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("debug") {
			log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
			log.Println("Using config file:", viper.ConfigFileUsed())
		}
	} else {
		if viper.GetBool("debug") {
			log.Println("error reading config file,", err)
		}
	}
}

func bindCustomFlag(flag *pflag.Flag) {
	if flag.Name == "config" {
		return
	}
	name := strings.ReplaceAll(flag.Name, "-", "_")
	viper.BindPFlag(name, flag)
}

func ensureConfigFileExists(configPath, configFile string) error {
	err := os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		return err
	}

	filePath := filepath.Join(configPath, configFile)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		_, err = file.Write([]byte("{}"))
		if err != nil {
			return err
		}
		defer file.Close()
	}

	return nil
}
