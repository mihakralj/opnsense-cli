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
	Short: `Discard changes made to the 'staging.xml' file`,
	Long:  `The 'discard' command reverses staged changes in the 'staging.xml' file. You can target specific nodes using an XPath expression. If no XPath is provided, all staged changes are discarded, effectively reverting 'staging.xml' to match the active 'config.xml'.`,
	Example: `  opnsense discard interfaces/wan/if   Discard changes to the 'if' node under the 'wan' interface
  opnsense discard                     Discard all staged changes in 'staging.xml'

To review staged changes, use 'show' or 'compare' command with no arguments.
Use the 'discard' command cautiously to avoid losing uncommitted changes.`,

	Run: func(cmd *cobra.Command, args []string) {

		internal.Checkos()

		configdoc := internal.LoadXMLFile(configfile, host)
		if configdoc == nil {
			internal.Log(1, "failed to get data from %s", configfile)
		}
		stagingdoc := internal.LoadXMLFile(stagingfile, host)
		if stagingdoc == nil {
			stagingdoc = configdoc
		}
		path := "opnsense"

		if len(args) < 1 {
			internal.Log(2, "Discarding all staged configuration changes.")
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

			if stagingElement != nil {
				if configElement != nil {
					stagingParent := stagingElement.Parent()
					stagingParent.RemoveChild(stagingElement)
					stagingParent.AddChild(configElement.Copy())
				} else {
					stagingParent := stagingElement.Parent()
					stagingParent.RemoveChild(stagingElement)
				}
			} else {
				stagingdoc.Root().AddChild(configElement.Copy())
			}
		}
		internal.SaveXMLFile(stagingfile, stagingdoc, host, true)
		fmt.Printf("Discarded all staged changes in %s\n", path)

	},
}

func init() {
	rootCmd.AddCommand(discardCmd)
}
