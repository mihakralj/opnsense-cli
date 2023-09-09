package internal

import (
	"encoding/json"
	"strings"

	"github.com/beevik/etree"
	"github.com/clbanning/mxj"
	"gopkg.in/yaml.v3"
)

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
