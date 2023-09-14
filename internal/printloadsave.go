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
	bash := fmt.Sprintf(`if [ -f %s ]; then cat %s; else echo "missing"; fi`, filename, filename)
	content := ExecuteCmd(bash, host)
	if content != "missing" {
		err := doc.ReadFromString(content)
		if err != nil {
			Log(1, "%s is not an XML file", filename)
		}
	}
	return doc
}

func SaveXMLFile(filename string, doc *etree.Document, host string, forced bool) {
	configout := ConfigToXML(doc, "opnsense")

	bash := `if [ -f "` + filename + `" ]; then echo "exists"; fi`
	fileexists := ExecuteCmd(bash, host)

	if strings.TrimSpace(fileexists) == "exists" {
		if !forced {
			Log(2, "%s already exists and will be overwritten.", filename)
		}
		// delete the file
		bash = "sudo rm " + filename
		ExecuteCmd(bash, host)
	}
	// chunking the long config.xml to upload in pieces
	chunkSize := 200000
	totalLength := len(configout)
	for i := 0; i < totalLength; i += chunkSize {
		end := i + chunkSize
		if end > totalLength {
			end = totalLength
		}
		chunk := configout[i:end]
		bash := fmt.Sprintf(`echo -n '%s' | sudo tee -a %s`, chunk, filename)
		ExecuteCmd(bash, host)
	}

	// check that file was written
	bash = `if [ -f "` + filename + `" ]; then echo "exists"; fi`
	fileexists = ExecuteCmd(bash, host)

	if fileexists == "exists" {
		Log(4, "%s has been succesfully saved.\n", filename)
	} else {
		Log(1, "error writing file %s", filename)
	}
}
