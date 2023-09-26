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
	"encoding/json"
	"strings"

	"github.com/beevik/etree"
	"gopkg.in/yaml.v3"
)

func ConfigToTTY(doc *etree.Document, path string) string {

	path = strings.TrimPrefix(path, "/")
	focused := FocusEtree(doc, path)
	d := depth + len(strings.Split(path, "/")) - 1
	if len(doc.FindElements(path)) > 1 {
		d -= 1
	}
	return EtreeToTTY(focused, d, 0)
}

func ConfigToXML(doc *etree.Document, path string) string {
	focused := FocusEtree(doc, path)
	newDoc := etree.NewDocument()
	newDoc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	//newDoc.CreateComment("XPath: " + path)
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
