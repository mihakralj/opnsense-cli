package internal

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
)

func PrintDocument(doc *etree.Document, path string) {
	var output string
	switch {
	case xmlFlag:
		output = ConfigToXML(doc, path)
	case jsonFlag:
		output = ConfigToJSON(doc, path)
	case yamlFlag:
		output = ConfigToYAML(doc, path)
	default:
		output = ConfigToTTY(doc, path)
	}
	fmt.Println(output)
}

func LoadXMLFile(filename string, host string) *etree.Document {
	if !strings.HasSuffix(filename, ".xml") {
		Log(1, "filename %s does not end with .xml", filename)
	}
	doc := etree.NewDocument()
	bash := fmt.Sprintf(`test -f "%s" && cat "%s" || echo "missing"`, filename, filename)
	content := ExecuteCmd(bash, host)
	if strings.TrimSpace(content) != "missing" {
		err := doc.ReadFromString(content)
		if err != nil {
			Log(1, "%s is not an XML file", filename)
		}
		return doc
	}
	return nil
}

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
