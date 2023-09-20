package internal

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/beevik/etree"
)

func FocusEtree(doc *etree.Document, path string) *etree.Element {
	foundElement := doc.FindElement(path)
	if foundElement == nil {
		Log(1, "Xpath element \"%s\" does not exist", path)
		return nil
	}

	parts := strings.Split(path, "/")
	focused := etree.NewElement(parts[0])

	// Get the space of the found element
	space := foundElement.Space

	depth := len(parts)
	if depth > 1 {
		parts = parts[:depth-1]
		current := focused
		for i := 1; i < len(parts); i++ {
			newElem := current.CreateElement(parts[i])
			// Find the element in the document and copy its attributes
			element := doc.FindElement(strings.Join(parts[:i+1], "/"))
			space = element.Space
			if space != "" {
				newElem.Space = space
			}
			if element != nil {
				for _, attr := range element.Attr {
					newElem.CreateAttr(attr.Key, attr.Value)
				}
			}
			current = newElem
		}
		current.AddChild(foundElement.Copy())
	} else {
		focused = doc.Root()
	}
	if space != "" {
		focused.Space = "att"
	}

	return focused
}

//////////

func EtreeToTTY(el *etree.Element, level int, indent int) string {
	EnumerateListElements(el)
	indentation := strings.Repeat("    ", indent)
	var result strings.Builder
	linePrefix := ""

	switch el.Space {
	case "att":
		linePrefix = c["chg"] + "!   "
	case "new":
		linePrefix = c["grn"] + "+   "
	case "chg":
		linePrefix = c["chg"] + "~   "
	case "del":
		linePrefix = c["red"] + "-   "
	default:
		linePrefix = c["tag"] + "    "
	}

	var attributestr, chgstr string
	for _, attr := range el.Attr {
		switch {
		case attr.Space == "new":
			attributestr += c["ita"] + c["new"] + fmt.Sprintf(" (%s=\"%s\")", attr.Key, attr.Value)
		case attr.Space == "chg":
			attributestr += c["ita"] + c["atr"] + fmt.Sprintf(" (%s=\""+c["red"]+"%s"+c["atr"]+"\")", attr.Key, strings.Replace(attr.Value, "|||", c["atr"]+"\""+c["arw"]+"\""+c["grn"], 1))
		case attr.Space == "del":
			attributestr += c["ita"] + c["red"] + fmt.Sprintf(" (%s=\"%s\")", attr.Key, attr.Value)
		default:
			attributestr += c["ita"] + c["atr"] + fmt.Sprintf(" (%s=\"%s\")", attr.Key, attr.Value)
		}
	}
	match, _ := regexp.MatchString(`\.\d+$`, el.Tag)
	if match {
		lastIndex := strings.LastIndex(el.Tag, ".")
		el.Tag = el.Tag[:lastIndex] + "[" + el.Tag[lastIndex+1:] + "]"
	}
	if len(el.ChildElements()) > 0 {
		result.WriteString(linePrefix + indentation + el.Tag + ":" + c["atr"] + attributestr + c["tag"] + " {" + c["nil"])

		if level > 0 {
			result.WriteString("\n")
			for _, child := range el.ChildElements() {
				result.WriteString(EtreeToTTY(child, level-1, indent+1))
			}

			result.WriteString("    " + indentation + c["tag"] + "}" + c["nil"] + "\n")
		} else {
			result.WriteString(c["nil"] + c["txt"] + c["ell"] + c["tag"] + "}\n")
		}

	} else {
		elText := el.Text()
		switch el.Space {
		case "chg":
			elText = c["nil"] + c["red"] + strings.Replace(elText, "|||", c["nil"]+c["arw"]+c["grn"], 1)
		case "del":
			elText = c["nil"] + c["red"] + strings.TrimSpace(elText)
		case "new":
			elText = c["nil"] + c["grn"] + strings.TrimSpace(elText)
		default:
			elText = c["nil"] + c["txt"] + strings.TrimSpace(elText)
		}
		content := chgstr + elText + c["nil"]
		result.WriteString(linePrefix + indentation + el.Tag + ":" + c["atr"] + attributestr + " " + content + c["nil"] + "\n")
	}
	return result.String()
}
