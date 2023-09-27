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
	"strings"

	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// compareCmd represents the compare command
var exportCmd = &cobra.Command{
	Use:   "export [<original.xml>] [<modified.xml>]",
	Short: `Export differences between two XML configuration files`,
	Long:  `The 'export' command generates a patch XML between two configuration files for the OPNsense firewall system. When only one filename is provided, it exports differences between that file and the current 'config.xml'. When no filenames are provided, it exports the patch from current 'config.xml' to 'staging.xml'`,
	Example: `  opnsense export b1.xml b2.xml  Exports XML patch from 'b1.xml' to 'b2.xml'
  opnsense export backup.xml     Exports XML patch from 'backup.xml' to 'config.xml'
  opnsense export                Exports XML patch from 'config.xml' to 'staging.xml'`,

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
		olddoc := internal.LoadXMLFile(oldconfig, host, false)
		newdoc := internal.LoadXMLFile(newconfig, host, true)
		if newdoc == nil {
			newdoc = olddoc
		}

//TODO: add support for xml lists

		deltadoc := internal.DiffXML(olddoc, newdoc, false)
		internal.RemoveChgSpace(deltadoc.Root())
		output := internal.ConfigToXML(deltadoc, path)
		fmt.Print(output)
	},
}

func init() {

	rootCmd.AddCommand(exportCmd)
}
