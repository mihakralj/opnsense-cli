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

func LoadXMLFile(filename string, host string, canBeMissing bool) *etree.Document {
	if !strings.HasSuffix(filename, ".xml") {
		Log(1, "filename %s does not end with .xml", filename)
	}
	doc := etree.NewDocument()
	bash := fmt.Sprintf(`test -f "%s" && cat "%s" || echo "missing"`, filename, filename)
	content := ExecuteCmd(bash, host)
	if strings.TrimSpace(content) == "missing" {
		if canBeMissing {
			return nil
		} else {
			Log(1, "failed to get data from %s", filename)
		}
	}
	err := doc.ReadFromString(content)
	if err != nil {
		Log(1, "%s is not an XML file", filename)
	}
	return doc

}
