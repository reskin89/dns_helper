/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/reskin89/dns_helper/dyndns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var dynCfg string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "d2",
	Short: "d2 (dyndns) is a small command to configure dynamic dns",
	Long: `d2 (dyndns) is a dyndns binary, intended to be run by a scheduler, to attain the 
	public ip of the running machine (even if behind NAT, the public NAT IP will be used)
	and update a provided route53 zone and dns A record if the IP and A record are different. 
	Currently only support IPv4`,
	Run: func(cmd *cobra.Command, args []string) {
		err := RunUpdater(args)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&dynCfg, "dynconfig", "", "Config File for DynDns (yml)")

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.d2.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".d2" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".d2")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func RunUpdater(args []string) error {
	var conf Configuration

	if dynCfg != nil {
		conf, err := dns_helper.NewConfigurationFromFile(dynCfg)
		if err != nil {
			log.Println(err)
			conf, err := dns_helper.NewConfigurationFromEnvironment()
			if err != nil {
				log.Println("Unable to load config from given file or environment...exiting....")
				log.Fatal(err)
			}
		}
	} else {
		conf, err := dns_helper.NewConfigurationFromEnvironment()
			if err != nil {
				log.Println("Unable to load config from environment...exiting....")
				log.Fatal(err)
			}
	}

	conf 
	log.Println(args)

	return nil
}
