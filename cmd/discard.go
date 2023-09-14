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

		internal.Checkos()

		configdoc := internal.LoadXMLFile(configfile, host)
		stagingdoc := internal.LoadXMLFile(stagingfile, host)
		if stagingdoc.Root() == nil {
			stagingdoc = configdoc
		}
		path := "opensense"

		if len(args) < 1 {
			stagingdoc = configdoc
		} else {
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

			configElement := configdoc.FindElement(path)
			stagingElement := stagingdoc.FindElement(path)
			stagingParent := stagingElement.Parent()
			stagingParent.RemoveChild(stagingElement)
			stagingParent.AddChild(configElement.Copy())

		}
		internal.SaveXMLFile(stagingfile, stagingdoc, host, false)
		fmt.Printf("Discarded staged changes in %s\n", path)

	},
}

func init() {
	rootCmd.AddCommand(discardCmd)
}
