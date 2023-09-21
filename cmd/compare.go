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
	"strings"

	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// compareCmd represents the compare command
var compareCmd = &cobra.Command{
	Use:   "compare [<original.xml>] [<modified.xml>]",
	Short: `Compare differences between two XML configuration files`,
	Long: `The 'compare' command identifies differences between two XML configuration files for the OPNsense firewall system. When only one filename is provided, it shows the differences between that file and the current 'config.xml'. When no filenames are provided, it compares the current 'config.xml' with 'staging.xml', akin to the 'show' command.`,
	Example: `  opnsense compare b1.xml b2.xml  Compare differences from 'b1.xml' to 'b2.xml'
  opnsense compare backup.xml     Compare differences from 'backup.xml' to 'config.xml'
  opnsense compare                Compare differences from 'config.xml' to 'staging.xml'`,

	Run: func(cmd *cobra.Command, args []string) {
		internal.SetFlags(verbose, force, host, configfile, nocolor, depth, xmlFlag, yamlFlag, jsonFlag)
		var oldconfig, newconfig, path string

		switch len(args) {
		case 3:
			oldconfig = "/conf/backup/" + args[0]
			newconfig = "/conf/backup/" + args[1]
			path = strings.Trim(args[2], "/")
		case 2:
			if strings.HasSuffix(args[1], ".xml") {
				oldconfig = "/conf/backup/" + args[0]
				newconfig = "/conf/backup/" + args[1]
				path = "opnsense"
			} else {
				oldconfig = "/conf/backup/" + args[0]
				newconfig = "/conf/config.xml"
				path = strings.Trim(args[1], "/")
			}
		case 1:
			if strings.HasSuffix(args[0], ".xml") {
				newconfig = "/conf/config.xml"
				oldconfig = "/conf/backup/" + args[0]
				path = "opnsense"
			} else {
				newconfig = "/conf/staging.xml"
				oldconfig = "/conf/config.xml"
				path = strings.Trim(args[0], "/")
			}
		default:
			oldconfig = "/conf/config.xml"
			newconfig = "/conf/staging.xml"
			path = "opnsense"
		}
		parts := strings.Split(path, "/")
		if parts[0] != "opnsense" {
			path = "opnsense/" + path
		}

		internal.Checkos()
		olddoc := internal.LoadXMLFile(oldconfig, host)
		newdoc := internal.LoadXMLFile(newconfig, host)
		if newdoc.Root() == nil {
			newdoc = olddoc
		}
		deltadoc := internal.DiffXML(olddoc, newdoc, true)
		internal.PrintDocument(deltadoc, path)

	},
}

func init() {
	compareCmd.Flags().IntVarP(&depth, "depth", "d", 1, "Specifies number of levels of returned tree (1-5)")
	rootCmd.AddCommand(compareCmd)
}
