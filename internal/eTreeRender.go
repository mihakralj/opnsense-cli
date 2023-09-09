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
    focused.CreateComment("XPath: " + path)

    // Get the space of the found element
    space := foundElement.Space

    depth := len(parts)
    if depth > 1 {
        parts = parts[:depth-1]
        current := focused
        for i := 1; i < len(parts); i++ {
            newElem := current.CreateElement(parts[i])
            // Set the space to "att" if the found element has space "att"
            if space == "att" {
                newElem.Space = "att"
            }
            current = newElem
        }
        current.AddChild(foundElement.Copy())
    } else {
        focused = doc.Root()
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
		linePrefix = c["nil"]+c["blu"] + "!   "
	case "new":
		linePrefix = c["grn"] + "+   " + c["grn"]
	case "chg":
		linePrefix = c["grn"] + "~   "
	case "del":
		linePrefix = c["bred"] + "-   " + c["red"]
	default:
		linePrefix = c["cyn"] + "    "
	}

	var attributestr, chgstr string
	for _, attr := range el.Attr {
		if attr.Space == "chg" {
			attributestr += c["ita"] + fmt.Sprintf(" (%s=\"%s\")", attr.Key, attr.Value)
		} else if attr.Key == "old" {
			chgstr = c["red"] + attr.Value + c["nil"] + c["arw"]
		}
	}
	match, _ := regexp.MatchString(`\.\d+$`, el.Tag)
	if match {
		lastIndex := strings.LastIndex(el.Tag, ".")
		el.Tag = el.Tag[:lastIndex] + "[" + el.Tag[lastIndex+1:] + "]"
	}
	if len(el.ChildElements()) > 0 {
		result.WriteString(linePrefix + indentation + el.Tag + ":"+c["cyn"]+" {" + c["nil"])

		if attributestr != "" {
			result.WriteString(linePrefix + indentation + "    " + attributestr + c["nil"] + "\n")
		}

		if level > 0 {
			result.WriteString("\n")
			for _, child := range el.ChildElements() {
				result.WriteString(EtreeToTTY(child, level-1, indent+1))
			}

			result.WriteString(c["cyn"] + "    " + indentation + "}" + c["nil"] + "\n")
		} else {
			result.WriteString(c["nil"] + c["bwht"] + c["ell"] + c["cyn"] + "}\n")
		}

	} else {
		elText := el.Text()
		switch el.Space {
		case "chg":
			elText = c["nil"] + c["bred"] + strings.Replace(elText, "|||", c["nil"]+c["arw"]+c["bgrn"], 1)
		case "del":
			elText = c["nil"] + c["bred"] + strings.TrimSpace(elText)
		case "new":
			elText = c["nil"] + c["grn"] + strings.TrimSpace(elText)
		default:
			elText = c["nil"] + c["bwht"] + strings.TrimSpace(elText)
		}
		content := chgstr + elText + c["nil"]
		result.WriteString(linePrefix + indentation + el.Tag + ": " + content + c["nil"] + "\n")

		if attributestr != "" {
			result.WriteString(linePrefix + indentation + "    " + attributestr + c["nil"] + "\n")
		}
	}
	return result.String()
}
