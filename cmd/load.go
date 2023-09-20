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
	"regexp"
	"strings"

	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:     "load [<backup.xml>]",
	Aliases: []string{"restore", "revert"},
	Short:   "Restores the firewall's active configuration from a backup file.",
	Long: `The 'load' command enables you to restore the active configuration of the OPNsense firewall system from a specified backup file located in the /conf/backup directory. If no filename is provided, the most recent backup will be loaded into config.xml.
Command has aliases 'restore' and 'revert' .

Examples:
  opnsense load                  - Restore from the most recent backup in /conf/backup.
  opnsense load config-123.xml   - Restore from a backup file /conf/backup/config-123.xml.
  opnsense load --force          - Restore from the most recent backup without interactive confirmation.
`,
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
		configdoc := internal.LoadXMLFile(filename, host)
		internal.Log(2, "Load %s into /conf/staging.xml.",filename)
		internal.SaveXMLFile(stagingfile, configdoc, host, true)
		fmt.Printf("The file %s has been loaded into /conf/staging.xml.\n", filename)
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)

}
