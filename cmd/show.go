/*
Copyright © 2023 MihaK mihak09@gmail.com

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
	"strings"

	"github.com/beevik/etree"
	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Displays information about configuration stored in config.xml",
	Long: `The show command displays configuration elements in config.xml, including interfaces, routes, firewall rules, and other system settings.
	 Use this command to view the current system configuration and troubleshoot issues.`,
	Run: func(cmd *cobra.Command, args []string) {

		path := "opnsense"
		if len(args) >= 1 {
			trimmedArg := strings.Trim(args[0], "/")
			if trimmedArg != "" {
				path = trimmedArg
			}
			parts := strings.Split(path, "/")
			if parts[0] != "opnsense" {
				path = "opnsense/"+path
			}
		}
		internal.Checkos()
		bash := "cat /conf/config.xml"
		config, err := internal.ExecuteCmd(bash, host)
		if err != nil {
			panic(err)
		}
		config_doc := etree.NewDocument()
		config_doc.ReadFromString(config)
		configtty := internal.ConfigToTTY(config_doc, path, depth)
		fmt.Println(configtty)

		//read config, path, depth
		// display ConfigToTTY

	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
