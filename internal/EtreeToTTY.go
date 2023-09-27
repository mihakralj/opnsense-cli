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
package internal

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
)

// EtreeToTTY returns a string representation of the etree.Element in a TTY-friendly format.
func EtreeToTTY(el *etree.Element, level int, indent int) string {
	// Enumerate list elements
	EnumerateListElements(el)

	// Set the indentation string for hierarchy
	indentation := strings.Repeat("    ", indent)

	// Set the line prefix based on the element's space
	var result strings.Builder

	// lead stores the leading spaces before root tag
	lead := c["nil"] + "  "
	linePrefix := ""
	switch el.Space {
	case "att":
		linePrefix = c["chg"] + "!" + lead + c["chg"]
	case "add":
		linePrefix = c["add"] + "+" + lead
		indentation = indentation + c["grn"]
	case "chg":
		linePrefix = c["chg"] + "~" + lead + c["chg"]
	case "del":
		linePrefix = c["del"] + "-" + lead
		indentation = indentation + c["del"]
	default:
		linePrefix = c["tag"] + " " + lead + c["tag"]
	}

	// Build the attribute string
	var attributestr, chgstr string
	for _, attr := range el.Attr {
		switch {
		case attr.Space == "del":
			attributestr += " " + c["ita"] + c["del"] + fmt.Sprintf("(%s=\"%s\")"+c["nil"], attr.Key, attr.Value)
			if el.Space == "" {
				linePrefix = c["del"] + "-" + lead + c["nil"]
			}
		case attr.Space == "add":
			attributestr += " " + c["ita"] + c["add"] + fmt.Sprintf("(%s=\"%s\")"+c["nil"], attr.Key, attr.Value)
			if el.Space == "" {
				linePrefix = c["add"] + "+" + lead + c["nil"]
			}
		case attr.Space == "chg":
			attributestr += c["tag"] + " (" + c["ita"] + c["chg"] + fmt.Sprintf("%s"+c["tag"]+"=\""+c["del"]+"%s"+c["tag"]+"\")"+c["nil"], attr.Key, strings.Replace(attr.Value, "|||", c["nil"]+c["tag"]+"\""+c["arw"]+"\""+c["grn"], 1))
			if el.Space == "" {
				linePrefix = c["chg"] + "~" + lead + c["nil"]
			}
		default:
			attributestr += c["tag"] + " (" + c["ita"] + c["atr"] + fmt.Sprintf("%s"+c["tag"]+"=\""+c["atr"]+"%s"+c["tag"]+"\")"+c["nil"], attr.Key, attr.Value)
		}
	}

	// Replace ".n" with "[n]" in the tag name
	el.Tag = ReverseEnumeratePath(el.Tag)
	/*
		match, _ := regexp.MatchString(`\.\d+$`, el.Tag)
		if match {
			lastIndex := strings.LastIndex(el.Tag, ".")
			el.Tag = el.Tag[:lastIndex] + "[" + el.Tag[lastIndex+1:] + "]"
		}
	*/

	// Build the content string
	if len(el.ChildElements()) > 0 {
		// If the element has child elements, build a block of nested elements
		result.WriteString(linePrefix + indentation + el.Tag + ":" + c["atr"] + attributestr + c["tag"] + " {" + c["nil"])

		if level > 0 {
			result.WriteString("\n")
			for _, child := range el.ChildElements() {
				result.WriteString(EtreeToTTY(child, level-1, indent+1))
			}
			result.WriteString(lead + " " + indentation + c["tag"] + "}" + c["nil"] + "\n")
		} else {
			result.WriteString(c["nil"] + c["txt"] + c["ell"] + c["tag"] + "}\n")
		}

	} else {
		// If the element has no child elements, build a single-line representation
		elText := el.Text()
		switch el.Space {
		case "chg":
			elText = c["nil"] + c["del"] + strings.Replace(elText, "|||", c["nil"]+c["arw"]+c["grn"], 1)
		case "del":
			elText = c["nil"] + c["del"] + strings.TrimSpace(elText)
		case "add":
			elText = c["nil"] + c["grn"] + strings.TrimSpace(elText)
		default:
			elText = c["nil"] + c["txt"] + strings.TrimSpace(elText)
		}
		content := chgstr + elText + c["nil"]
		if el.Parent().GetPath() == "/" && len(el.ChildElements()) == 0 {
			result.WriteString(linePrefix + indentation + el.Tag + ": {\n" + linePrefix + "}")

		} else {
			result.WriteString(linePrefix + indentation + el.Tag + ":" + c["atr"] + attributestr + c["nil"] + " " + content + c["nil"] + "\n")
		}
	}

	return result.String()
}
