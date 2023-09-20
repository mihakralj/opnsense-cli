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
	"regexp"
	"strings"

	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// configCmd represents the show command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Displays information about configuration stored in config.xml",
	Long: `The show command displays configuration elements in config.xml, including interfaces, routes, firewall rules, and other system settings.
	 Use this command to view the current system configuration and troubleshoot issues.`,
	Run: func(cmd *cobra.Command, args []string) {

		path := "opnsense"
		if len(args) >= 1 {
			trimmedArg := strings.Trim(args[0], "/")
			if matched, _ := regexp.MatchString(`\[0\]`, trimmedArg); matched {
				internal.Log(1, "XPath indexing of elements starts with 1, not 0")
			}
			if trimmedArg != "" {
				path = trimmedArg
			}
			parts := strings.Split(path, "/")
			if parts[0] != "opnsense" {
				path = "opnsense/" + path
			}
		}

		internal.Checkos()

		configdoc := internal.LoadXMLFile(configfile, host)
		stagingdoc := internal.LoadXMLFile(stagingfile, host)
		if stagingdoc.Root() == nil {
			stagingdoc = configdoc
		}

		deltadoc := internal.DiffXML(configdoc, stagingdoc, true)
		internal.PrintDocument(deltadoc, path)

	},
}

func init() {
	configCmd.Flags().IntVarP(&depth, "depth", "d", 1, "Specifies number of levels of returned tree (1-5)")
	showCmd.AddCommand(configCmd)
}
