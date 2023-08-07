package internal

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/clbanning/mxj"
)

func recursiveFocus(el *etree.Element, parts []string, depth int) *etree.Element {
	if el.Tag != parts[0] {
		return nil
	}

	focused := etree.NewElement(el.Tag)
	for _, attr := range el.Attr {
		focused.CreateAttr(attr.Key, attr.Value)
	}
	focused.SetText(el.Text())

	if len(parts) > 1 {
		for _, child := range el.ChildElements() {
			focusedChild := recursiveFocus(child, parts[1:], depth)
			if focusedChild != nil {
				focused.AddChild(focusedChild)
			}
		}
	} else if len(parts) == 1 {
		addChildren(focused, el, depth)
	}

	return focused
}

func addChildren(focused, el *etree.Element, depth int) {
	if depth == 0 {
		return
	}

	for _, child := range el.ChildElements() {
		shallow := etree.NewElement(child.Tag)
		for _, attr := range child.Attr {
			shallow.CreateAttr(attr.Key, attr.Value)
		}
		shallow.SetText(child.Text())
		focused.AddChild(shallow)
		if depth == -1 {
			addChildren(shallow, child, depth)
		} else {
			addChildren(shallow, child, depth-1)
		}
	}
}

func EtreeToTTY(el *etree.Element, level int, indent int) string {
	indentation := strings.Repeat("    ", indent)
	var result strings.Builder
	// comments
	commentstr := ""
	for _, token := range el.Child {
		if comment, ok := token.(*etree.Comment); ok {
			commentstr = fmt.Sprintf(c["red"]+" %s"+c["nil"], strings.TrimSpace(comment.Data))
			//commentstr = fmt.Sprintf("# %s", strings.TrimSpace(comment.Data))
		}
	}
	// attributes
	attributestr := ""
	for _, attr := range el.Attr {
		attributestr = fmt.Sprintf(c["ita"]+c["blu"]+" (%s=\"%s\")"+c["nil"], attr.Key, attr.Value)
	}

	if len(el.ChildElements()) > 0 {
		result.WriteString(fmt.Sprintf("%s%s: {", indentation, el.Tag))
		result.WriteString(attributestr)

		if level > 0 {
			// Print comments
			result.WriteString(" "+commentstr+"\n")

			for _, child := range el.ChildElements() {
				result.WriteString(EtreeToTTY(child, level-1, indent+1))
			}
		} else {
			if attributestr != "" {
				result.WriteString(" ")
			}
			//result.WriteString("\033[32m\u2026\033[0m")
			result.WriteString(c["grn"]+c["ell"]+c["nil"])

			//result.WriteString("...")
			// this part is shit
			indentation = ""
		}

		result.WriteString(fmt.Sprintf("%s}", indentation))
		if level == 0 {
			result.WriteString(fmt.Sprintf(" %s", commentstr))
		}
		result.WriteString("\n")

	} else {
		content := strings.ReplaceAll(strings.TrimSpace(el.Text()), "\n", "")
		result.WriteString(fmt.Sprintf("%s%s:%s "+c["grn"]+"\"%s\""+c["nil"], indentation, el.Tag, attributestr, content))

		//result.WriteString(fmt.Sprintf("%s%s:%s \033[32m\"%s\"\033[0m", indentation, el.Tag, attributestr, content))
		//result.WriteString(fmt.Sprintf("%s%s:%s \"%s\"", indentation, el.Tag, attributestr, content))
		// Print comments
		result.WriteString(" "+commentstr+"\n")
	}
	return result.String()
}

func ConfigToTTY(doc *etree.Document, path string, level int) string {
	// check

	parts := strings.Split(path, "/")
	focused := etree.NewElement(parts[0])
	depth := len(parts)

		if depth > 1 {
			parts = parts[:depth-1]
			current := focused
			for i := 1; i < len(parts); i++ {

				current = current.CreateElement(parts[i])
			}
			if doc.FindElement(path) != nil {
				current.AddChild(doc.FindElement(path).Copy())
			}
		} else {
			focused = doc.Root()
		}

	return EtreeToTTY(focused, depth+level-1, 0)
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

func ConfigToJSON(doc *etree.Document, path string) string {
	focused := etree.NewElement("opnsense")
	parts := strings.Split(path, "/")
	depth := len(parts)

	if depth > 1 {
		parts = parts[:depth-1]
		current := focused
		for i := 1; i < len(parts); i++ {
			current = current.CreateElement(parts[i])
		}
		current.AddChild(doc.FindElement(path).Copy())
	} else {
		focused = doc.Root()
	}
	res, _ := EtreeToJSON(focused)
	return res
}
