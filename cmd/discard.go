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

// discardCmd represents the discard command
var discardCmd = &cobra.Command{
	Use:   "discard [<xpath>]",
	Short: "Discards changes in the staging.xml file.",
	Long: `The 'discard' command discard changes staged in the staging.xml file. You can discard changes for a specific node identified
by an XPath or discard all changes if no XPath is provided.

Using an XPath allows targeting specific nodes in staging.xml without affecting other modifications that are staged for commit.
If no XPath is specified, all changes in the staging.xml file will be discarded, removing any staged changes to the config.xml.

Examples:
  opnsense discard interfaces/wan/if   - discards changes made to the 'if' node under 'wan' interface
  opnsense discard system/hostname     - discards changes made to the 'hostname' node
  opnsense discard                     - discards all changes in the staging.xml file.

To review staged changes, use 'opnsense show config' or 'opnsense compare' (without arguments).
Always use discard command with caution to avoid losing uncommitted work.`,
	Run: func(cmd *cobra.Command, args []string) {

		// check if arg is providede
		// if not, delete the full staging.xml
		// else:
		// read config.xml
		// read staging.xml
		// replace the xpath node in staging.xml with node from config.xml
		// save staging.xml

		cmd.Help()

		fmt.Println("\n\033[34mDiscard command is not implemented yet\033[0m")
	},
}

func init() {
	rootCmd.AddCommand(discardCmd)
}
