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
	"regexp"
	"strings"

	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// loadCmd represents the load command
var restoreCmd = &cobra.Command{
	Use:     "restore [<backup.xml>]",
	Aliases: []string{"load"},
	Short:   `Restore active configuration from a backup XML file`,
	Long: `The 'restore' command restores the active configuration of the OPNsense firewall system using a backup file from the '/conf/backup' directory. When no filename is provided, the system defaults to using the most recent backup. The command also has alias 'load'.`,
	Example: `  opnsense restore              Restore from the most recent backup in '/conf/backup'
  opnsense load config-123.xml  Restore from the specified backup file in '/conf/backup'`,

	Run: func(cmd *cobra.Command, args []string) {

		var filename string
		if len(args) < 1 {
			bash := "ls -t /conf/backup/*.xml | head -n 1"
			filename = strings.TrimSpace(internal.ExecuteCmd(bash, host))
			if filename == "" {
				internal.Log(1, "No backup files found in /conf/backup.")
				return
			}
		} else {
			filename = args[0]

		}
		filename = strings.TrimPrefix(filename, "/conf/backup/")
		filename = strings.TrimPrefix(filename, "conf/backup/")
		if !strings.HasSuffix(filename, ".xml") {
			filename += ".xml"
		}
		validFilenamePattern := "^[a-zA-Z0-9_.-]+$"
		match, err := regexp.MatchString(validFilenamePattern, filename)
		if err != nil || !match {
			internal.Log(1, "%s is not a valid filename.", filename)
			return
		}
		filename = "/conf/backup/"+filename
		internal.Checkos()

		configdoc := internal.LoadXMLFile(configfile, host, false)
		saveddoc := internal.LoadXMLFile(filename, host, false)

		depthset := cmd.LocalFlags().Lookup("depth")
		if depthset != nil && !depthset.Changed {
		internal.FullDepth()
		}

		deltadoc := internal.DiffXML(configdoc, saveddoc, false)

		internal.PrintDocument(deltadoc, "opnsense")

		internal.Log(2, "Stage %s into %s",filename, stagingfile)
		internal.SaveXMLFile(stagingfile, saveddoc, host, true)
		fmt.Printf("The file %s has been staged into %s.\n", filename, stagingfile)
	},
}

func init() {
	restoreCmd.Flags().IntVarP(&depth, "depth", "d", 1, "Specifies number of depth levels of returned tree (default: 1)")
	rootCmd.AddCommand(restoreCmd)

}
