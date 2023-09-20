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
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show [segment]",
	Short: "Display information related to OPNsense system",
	Long: `The 'show' command retrieves various details about the OPNsense system.

Examples:
  show config <xpath>  - Hierarchical segments of config.xml
  show system <xpath>  - System information about OPNsense firewall
  show backup -d2      - Details about files in /conf/backup`,
	Run: func(cmd *cobra.Command, args []string) {
		//cmd.Help()
		configCmd.Run(cmd, args)
	},
}

func init() {
	showCmd.Flags().IntVarP(&depth, "depth", "d", 1, "Specifies number of levels of returned tree (1-5)")
	rootCmd.AddCommand(showCmd)
}
