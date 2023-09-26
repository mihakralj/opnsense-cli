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

		configdoc := internal.LoadXMLFile(configfile, host, false)
		stagingdoc := internal.LoadXMLFile(stagingfile, host, true)
		if stagingdoc == nil {
			stagingdoc = configdoc
		}
		path := "opnsense"

		if len(args) < 1 {
			internal.Log(2, "Discarding all staged configuration changes.")
			stagingdoc = configdoc
		} else {

			if matched, _ := regexp.MatchString(`\[0\]`, args[0]); matched {
				internal.Log(1, "XPath indexing of elements starts with 1, not 0")
			}
			if args[0] != "" {
				path = args[0]
			}

			parts := strings.Split(path, "/")
			if parts[0] != "opnsense" {
				path = "opnsense/" + path
			}
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}

			stagingEl := stagingdoc.FindElement(path)
			configEl := configdoc.FindElement(path)

			if configEl == nil && stagingEl != nil {
				// Element is new in staging, remove it
				parent := stagingEl.Parent()
				parent.RemoveChild(stagingEl)

				// Remove the last part of the path
				lastSlash := strings.LastIndex(path, "/")
				if lastSlash != -1 {
					path = path[:lastSlash]
				}
			} else if configEl != nil && stagingEl != nil {
				// Element exists in both configdoc and stagingdoc, restore it
				parent := stagingEl.Parent()
				parent.RemoveChild(stagingEl)
				parent.AddChild(configEl.Copy())

				// Restore attributes
				configAttrs := configEl.Attr
				stagingEl = parent.FindElement(configEl.Tag)
				if stagingEl != nil {
					for _, attr := range configAttrs {
						stagingEl.CreateAttr(attr.Key, attr.Value)
					}
				}
			} else if configEl != nil && stagingEl == nil {
				// Element exists in configdoc but not in stagingdoc, add it to stagingdoc
				stagingdoc.Root().AddChild(configEl.Copy())

				// Copy attributes
				configAttrs := configEl.Attr
				stagingEl = stagingdoc.Root().FindElement(configEl.Tag)
				if stagingEl != nil {
					for _, attr := range configAttrs {
						stagingEl.CreateAttr(attr.Key, attr.Value)
					}
				}
			}
		}

		if len(args) < 1 {
			fmt.Printf("Discarded all staged configuration changes")
		} else {
			fmt.Printf("Discarded staged changes in node %s:\n\n", path)
			deltadoc := internal.DiffXML(configdoc, stagingdoc, true)
			internal.PrintDocument(deltadoc, path)
		}
		internal.SaveXMLFile(stagingfile, stagingdoc, host, true)
	},
}

func init() {
	rootCmd.AddCommand(discardCmd)
}
