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
	"math/rand"
	"strings"
	"time"

	"github.com/beevik/etree"
)

func SaveXMLFile(filename string, doc *etree.Document, host string, forced bool) {
	configout := ConfigToXML(doc, "opnsense")
	bash := `test -f "` + filename + `" && echo "exists" || echo "missing"`
	fileexists := ExecuteCmd(bash, host)

	if strings.TrimSpace(fileexists) == "exists" {
		if !forced {
			Log(2, "%s already exists and will be overwritten.", filename)
		}
		// delete the file
		bash = "sudo rm " + filename
		ExecuteCmd(bash, host)
	}
	sftpCmd(configout, filename, host)

	// check that file was written
	bash = `test -f "` + filename + `" && echo "exists" || echo "missing"`
	fileexists = ExecuteCmd(bash, host)

	if strings.TrimSpace(fileexists) == "exists" {
		Log(4, "%s has been succesfully saved.\n", filename)
	} else {
		Log(1, "error writing file %s", filename)
	}
}

func GenerateBackupFilename() string {
	timestamp := time.Now().Unix()
	randomNumber := rand.Intn(10000)
	filename := fmt.Sprintf("config-%d.%04d.xml", timestamp, randomNumber)
	return filename
}
