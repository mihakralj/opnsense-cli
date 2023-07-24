package internal

import (
	"fmt"
	"strings"
	"github.com/beevik/etree"
	"github.com/clbanning/mxj"
)

func recursiveFocus(el *etree.Element, parts []string, full bool) *etree.Element {
	if el.Tag != parts[0] {
		return nil
	}

	focused := etree.NewElement(el.Tag)
	for _, attr := range el.Attr {
		focused.CreateAttr(attr.Key, attr.Value)
	}
	focused.SetText(el.Text())

	if len(parts) == 1 {
		if full {
			for _, child := range el.ChildElements() {
				focused.AddChild(child.Copy())
			}
		} else {
			for _, child := range el.ChildElements() {
				shallow := etree.NewElement(child.Tag)
				for _, attr := range child.Attr {
					shallow.CreateAttr(attr.Key, attr.Value)
				}
				shallow.SetText(child.Text())
				focused.AddChild(shallow)
			}
		}
	} else {
		for _, child := range el.ChildElements() {
			focusedChild := recursiveFocus(child, parts[1:], full)
			if focusedChild != nil {
				focused.AddChild(focusedChild)
			}
		}
	}

	return focused
}

func FocusTree(el *etree.Element, path string, full bool) *etree.Element {
	// Remove leading slash if it exists
	cleanedPath := strings.TrimPrefix(path, "/")
	parts := strings.Split(cleanedPath, "/")
	return recursiveFocus(el, parts, full)
}

func EtreeToYaml(el *etree.Element, indent int) string {
	indentation := strings.Repeat("  ", indent)
	var result strings.Builder

	if len(el.ChildElements()) > 0 {
		// For an element with child elements, we don't want to append the content yet
		result.WriteString(fmt.Sprintf("%s%s:", indentation, el.Tag))
		result.WriteString("\n")
	} else {
		content := strings.ReplaceAll(el.Text(), "\n", "")
		if content == "" {
			content = "~"
		}
		result.WriteString(fmt.Sprintf("%s%s: %s", indentation, el.Tag, content))

		// Append the attributes as comments at the end of the line
		if len(el.Attr) > 0 {
			result.WriteString("  # ")
			for i, attr := range el.Attr {
				if i != 0 {
					result.WriteString(", ")
				}
				result.WriteString(fmt.Sprintf("%s=%s", attr.Key, attr.Value))
			}
		}

		result.WriteString("\n") // append newline here for leaf nodes
	}

	for _, child := range el.ChildElements() {
		result.WriteString(EtreeToYaml(child, indent+1))
	}

	return result.String()
}

func EtreeToJSON(el *etree.Element) (string, error) {
	doc := etree.NewDocument()
	doc.SetRoot(el.Copy())
	str, err := doc.WriteToString()
	if err != nil {
		return "", err
	}

	mv, err := mxj.NewMapXml([]byte(str)) // parse xml to map
	if err != nil {
		return "", err
	}

	jsonStr, err := mv.JsonIndent("", "  ") // convert map to json
	if err != nil {
		return "", err
	}

	return string(jsonStr), nil
}
