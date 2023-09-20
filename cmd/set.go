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

var deleteFlag bool = false

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set <xpath> <value>",
	Short: "Sets a value for a specific node in the staging.xml file.",
	Long: `The 'set' command allows you to assign a new value to a specific node within the staging.xml file, effectively modifying the configuration in a controlled manner.

Before the changes can take effect, you need to use the 'commit' command to move the staging.xml file to the active config.xml. If at any point you decide to discard the changes made, you can use the 'discard' command.

The XPath parameter allows for precise targeting of the nodes in the XML structure, helping you to navigate to the exact setting or property that you wish to update.

Examples:
  opnsense set interfaces/wan/if igb0   - sets the 'interfaces/wan/if' node with the value 'igb0'
  opnsense set system/hostname myrouter - assigns 'myrouter' as the new hostname in the staging.xml file.

Make sure to validate your XPath expressions to avoid any unintended changes.`,
	Run: func(cmd *cobra.Command, args []string) {

		internal.Checkos()

		configdoc := internal.LoadXMLFile(configfile, host)
		stagingdoc := internal.LoadXMLFile(stagingfile, host)
		if stagingdoc.Root() == nil {
			stagingdoc = configdoc
		}

		if len(args) == 0 {
			internal.Log(1, "XPath not provided")
			return
		}

		path := strings.Trim(args[0], "/")
		if !strings.HasPrefix(path, "opnsense/") {
			path = "opnsense/" + path
		}
		if matched, _ := regexp.MatchString(`\[0\]`, path); matched {
			internal.Log(1, "XPath indexing of elements starts with 1, not 0")
			return
		}

		var attribute, value string

		if len(args) == 2 {
			if isAttribute(args[1]) {
				attribute = args[1]
			} else {
				value = strings.Trim(args[1], " ")
			}

		}
		if len(args) == 3 {
			if isAttribute(args[1]) {
				attribute = args[1]
				if !isAttribute((args[2])) {
					value = strings.Trim(args[2], " ")
				} else {
					internal.Log(1, "Too many attributes provided")
				}
			} else {
				value = strings.Trim(args[1], " ")
				if isAttribute(args[2]) {
					attribute = args[2]
				} else {
					internal.Log(1, "Too many values provided")
				}
			}
		}

		element := stagingdoc.FindElement(path)

		if element != nil {
			element.SetText(value)
		} else {
			element := stagingdoc.Root()
			parts := strings.Split(path, "/")
			for i, part := range parts {
				if i == 0 && part == "opnsense" {
					continue
				}
				if element.SelectElement(part) == nil && !deleteFlag {
					element.CreateElement(part)
				}
				element = element.SelectElement(part)
			}
			if deleteFlag {
				element.Parent().RemoveChild(element)
			} else {
				element.SetText(value)
			}
		}
		deltadoc := internal.DiffXML(configdoc, stagingdoc, true)
		internal.PrintDocument(deltadoc, path)

		internal.SaveXMLFile(stagingfile, stagingdoc, host, false)

		fmt.Printf("Path: %s\t Attribute: %s\t Value: %s\n", path, attribute, value)

	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().BoolVarP(&deleteFlag, "delete", "d", false, "Delete a node")

}

func isAttribute(s string) bool {
	re := regexp.MustCompile(`^\([^=]+=[^=]+\)$`)
	return re.MatchString(s)
}
