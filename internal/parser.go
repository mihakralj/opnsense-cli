package internal

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/clbanning/mxj"
	"gopkg.in/yaml.v3"
)

/*
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
*/

func FocusEtree(doc *etree.Document, path string) *etree.Element {
	if doc.FindElement(path) == nil {
		Log(1, "Xpath element \"%s\" does not exist",path)
	}
	parts := strings.Split(path, "/")
	focused := etree.NewElement(parts[0])
	//focused.CreateAttr("xmlns", "https://opnsense.org/namespace")
	focused.CreateComment("XPath: " + path)
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
	return focused
}

//////////

func EtreeToTTY(el *etree.Element, level int, indent int) string {
	indentation := strings.Repeat("    ", indent)
	var result strings.Builder
	// comments
	commentstr := ""
	for _, token := range el.Child {
		if comment, ok := token.(*etree.Comment); ok {
			commentstr = fmt.Sprintf(c["red"]+" %s"+c["nil"], strings.TrimSpace(comment.Data))
		}
	}
	// attributes
	attributestr := ""
	for _, attr := range el.Attr {
		attributestr += fmt.Sprintf(c["ita"]+c["blu"]+" (%s=\"%s\")"+c["nil"], attr.Key, attr.Value)
	}
	if len(el.ChildElements()) > 0 {
		result.WriteString(fmt.Sprintf("%s%s: {", indentation, el.Tag))
		result.WriteString(attributestr)
		if level > 0 {
			result.WriteString(fmt.Sprintf(" %s\n", commentstr))

			for _, child := range el.ChildElements() {
				result.WriteString(EtreeToTTY(child, level-1, indent+1))
			}
		} else {
			if attributestr != "" {
				result.WriteString(" ")
			}
			result.WriteString(c["grn"] + c["ell"] + c["nil"])
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
		result.WriteString(fmt.Sprintf(" %s\n", commentstr))
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

//////////

func ConfigToTTY(doc *etree.Document, path string) string {
	focused := FocusEtree(doc, path)
	//calculate depth of path
	return EtreeToTTY(focused, depth+len(strings.Split(path, "/"))-1, 0)
}

func ConfigToXML(doc *etree.Document, path string) string {
	focused := FocusEtree(doc, path)
	newDoc := etree.NewDocument()
	newDoc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	newDoc.SetRoot(focused)
	newDoc.Indent(2)
	xmlStr, err := newDoc.WriteToString()
	if err != nil {
		return ""
	}
	return xmlStr
}

func ConfigToJSON(doc *etree.Document, path string) string {
	focused := FocusEtree(doc, path)
	res, _ := EtreeToJSON(focused)
	return res
}

func ConfigToYAML(doc *etree.Document, path string) string {
	focused := FocusEtree(doc, path)
	jsonStr, _ := EtreeToJSON(focused)
	var jsonObj interface{}
	json.Unmarshal([]byte(jsonStr), &jsonObj)

	yamlBytes, _ := yaml.Marshal(jsonObj)
	return string(yamlBytes)
}
