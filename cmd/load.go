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

	"github.com/spf13/cobra"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load [<backup.xml>]",
	Aliases: []string{"restore","revert"},
	Short: "Restores the firewall's active configuration from the most recent or from a specific backup file.",
	Long: `The 'load' command enables you to restore the active configuration of the OPNsense firewall system from a specified backup file located in the /conf/backup directory. If no filename is provided, the most recent backup will be loaded into config.xml.
Command has aliases 'restore' and 'revert' .

Examples:
  opnsense load                  - Restore from the most recent backup in /conf/backup.
  opnsense load config-123.xml   - Restore from a backup file /conf/backup/config-123.xml.
  opnsense load --force          - Restore from the most recent backup without interactive confirmation.
`,
	Run: func(cmd *cobra.Command, args []string) {

		// if no parameter, find the latest backup file
		// copy backup file to config.xml
		// configctl config reload

		cmd.Help()

		fmt.Println("\n\033[34mLoad command is not implemented yet\033[0m")
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)

}
