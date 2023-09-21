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
	"html"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

var deleteFlag bool = false

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set <xpath> <value> <(att=value)>",
	Short: "Sets a value and an attribute for a specific node in the staging.xml file.",
	Long: `The 'set' command allows you to assign a new value to a specific node within the staging.xml file, effectively modifying the configuration in a controlled manner.

Before the changes can take effect, you need to use the 'commit' command to move the staging.xml file to the active config.xml. If at any point you decide to discard the changes made, you can use the 'discard' command.

The XPath parameter allows for precise targeting of the nodes in the XML structure, helping you to navigate to the exact setting or property that you wish to update.

Examples:
  opnsense set interfaces/wan/if igb0         - sets the 'interfaces/wan/if' node with the value 'igb0'
  opnsense set system/hostname myrouter       - assigns 'myrouter' as the new hostname in the staging.xml file.
  opnsense set interfaces "(version=2.0)"     - assign an attribute 'version' with value '2.0' to node 'interfaces'
  opnsense set system/hostname -d             - remove the node and all its content, attributes and children.

Make sure to validate your XPath expressions to avoid any unintended changes.`,
	Run: func(cmd *cobra.Command, args []string) {

		internal.Checkos()

		configdoc := internal.LoadXMLFile(configfile, host)
		internal.EnumerateListElements(configdoc.Root())

		stagingdoc := internal.LoadXMLFile(stagingfile, host)
		internal.EnumerateListElements(stagingdoc.Root())

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
		path = strings.ReplaceAll(path, "[", ".")
		path = strings.ReplaceAll(path, "]", "")

		if matched, _ := regexp.MatchString(`\[0\]`, path); matched {
			internal.Log(1, "XPath indexing of elements starts with 1, not 0")
			return
		}

		var attribute, value string

		if len(args) == 2 {
			if isAttribute(args[1]) {
				attribute = escapeXML(args[1])

			} else {
				value = escapeXML(strings.Trim(args[1], " "))
			}

		}
		if len(args) == 3 {
			if isAttribute(args[1]) {
				attribute = escapeXML(args[1])
				if !isAttribute((args[2])) {
					value = escapeXML(strings.Trim(args[2], " "))
				} else {
					internal.Log(1, "Too many attributes provided")
				}
			} else {
				value = strings.Trim(args[1], " ")
				if isAttribute(args[2]) {
					attribute = escapeXML(args[2])
				} else {
					internal.Log(1, "Too many values provided")
				}
			}
		}

		element := stagingdoc.FindElement(path)

		if !deleteFlag {
			if element == nil {
				element = stagingdoc.Root()
				parts := strings.Split(path, "/")
				for i, part := range parts {
					part = fixXMLName(part)
					if i == 0 && part == "opnsense" {
						continue
					}
					if element.SelectElement(part) == nil {
						if element.SelectElement(part+".1") != nil {
							var maxIndex int
							for _, child := range element.ChildElements() {
								if strings.HasPrefix(child.Tag, part+".") {
									indexStr := strings.TrimPrefix(child.Tag, part+".")
									index, err := strconv.Atoi(indexStr)
									if err == nil && index > maxIndex {
										maxIndex = index
									}
								}
							}
							part = fmt.Sprintf("%s.%d", part, maxIndex+1)
						}
						element.CreateElement(part)
					}
					element = element.SelectElement(part)
					fmt.Println(part, element.GetPath())
				}
				path = element.GetPath()
			}
			if value != "" {
				children := element.ChildElements()
				if len(children) > 0 {
					internal.Log(1, "%s is an element container and cannot store content.", element.GetPath())
				}
				element.SetText(value)
				path = element.GetPath()
			}
			if attribute != "" {
				attribute = strings.Trim(attribute, "()") // remove parentheses
				parts := strings.Split(attribute, "=")
				if len(parts) == 2 {
					key := fixXMLName(parts[0])
					val := escapeXML(parts[1])
					element.CreateAttr(key, val)
				} else {
					internal.Log(1, "Invalid attribute format")
				}
			}
		} else {
			if value == "" && attribute == "" && element != nil {
				parent := element.Parent()
				if parent != nil {
					parent.RemoveChild(element)
					path = parent.GetPath()
				} else {
					internal.Log(1, "Cannot delete the root element")
				}
			}
			if value != "" {
				element.SetText("")
				path = element.GetPath()
			}
			if attribute != "" {
				attribute = strings.Trim(attribute, "()") // remove parentheses
				parts := strings.Split(attribute, "=")
				if len(parts) == 2 {
					key := fixXMLName(parts[0])
					element.RemoveAttr(key)
					fmt.Println("deleted attribute", key)

				} else {
					internal.Log(1, "Invalid attribute format")
				}
			}
		}
		deltadoc := internal.DiffXML(configdoc, stagingdoc, true)

		internal.ReverseEnumerateListElements(configdoc.Root())
		internal.ReverseEnumerateListElements(stagingdoc.Root())
		re := regexp.MustCompile(`\.(\d+)`)
		path = re.ReplaceAllString(path, "[$1]")
		internal.PrintDocument(deltadoc, path)

		internal.SaveXMLFile(stagingfile, stagingdoc, host, true)
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

func escapeXML(value string) string {
	value = strings.TrimSpace(value)
	escapedValue := html.EscapeString(value)
	return escapedValue
}

func fixXMLName(value string) string {
	// Trim the input string
	value = strings.TrimSpace(value)
	if value == "" {
		return "_"
	}

	// Ensure the first character is a valid start character
	for len(value) > 0 && !isXMLNameStartChar(rune(value[0])) {
		value = value[1:]
	}

	// If no valid start character was found, prepend an underscore
	if len(value) == 0 {
		value = "_"
	}

	// Ensure all other characters are valid name characters
	runes := []rune(value)
	for i, r := range runes {
		if !isXMLNameChar(r) {
			runes[i] = '_'
		}
	}

	return string(runes)
}

// Checks if a rune is a valid XML name start character
func isXMLNameStartChar(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

// Checks if a rune is a valid XML name character
func isXMLNameChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '-' || r == '_' || r == ':' || r == '[' || r == ']'
}
