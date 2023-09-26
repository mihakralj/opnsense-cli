/*
Copyright © 2023 Miha miha.kralj@outlook.com

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
	"os"

	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

var (
	Version     string = "0.13.0"
	verbose     int
	force       bool
	host        string
	configfile  string
	stagingfile string
	nocolor     bool
	depth       int = 1
	xmlFlag     bool
	yamlFlag    bool
	jsonFlag    bool
	ver_flag    bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&host, "target", "t", "", "Specify target host (user@hostname[:port])")
	rootCmd.PersistentFlags().IntVarP(&verbose, "verbose", "v", 1, "Set verbosity level (range: 1-5, default: 1)")
	rootCmd.PersistentFlags().BoolVarP(&nocolor, "no-color", "n", false, "Disable ANSI color output")
	rootCmd.PersistentFlags().BoolVarP(&xmlFlag, "xml", "x", false, "Output results in XML format")
	rootCmd.PersistentFlags().BoolVarP(&jsonFlag, "json", "j", false, "Output results in JSON format")
	rootCmd.PersistentFlags().BoolVarP(&yamlFlag, "yaml", "y", false, "Output results in YAML format")
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Bypass checks and prompts (force action)")
	rootCmd.Flags().BoolVarP(&ver_flag, "version", "V", false, "Display the version of opnsense")
	//rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	cobra.OnInitialize(func() {
		configfile = "/conf/config.xml"
		stagingfile = "/conf/staging.xml"
		internal.SetFlags(verbose, force, host, configfile, nocolor, depth, xmlFlag, yamlFlag, jsonFlag)
		//other initializations
	})
}

var rootCmd = &cobra.Command{
	Use:   "opnsense [command]",
	Short: "CLI tool for managing and monitoring OPNsense firewall systems",
	Long: `The 'opnsense' command-line utility provides non-GUI administration of OPNsense firewall systems. It can be run locally on the firewall or remotely via an SSH tunnel.

To streamline remote operations, add your private key to the SSH agent using 'ssh-add' and the matching public key to the admin account on OPNsense.`,
	Example: `  opnsense help [COMMAND]                Display help for specific commands
  opnsense show interfaces/wan           Show details for the interfaces/wan node in config.xml
  opnsense -t admin@192.168.1.1 sysinfo  Retrieve system information from a remote firewall`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if ver_flag {
			fmt.Println("opnsense-CLI version", Version)
			os.Exit(0)
		}

		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.DisableNoDescFlag = true
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing CLI '%s'", err)
		os.Exit(1)
	}
}
