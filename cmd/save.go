/*
Copyright © 2023 MihaK mihak09@gmail.com

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
package cmd

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save [filename]",
	Short: "Creates a new backup of the current firewall configuration in /conf/backup directory.",
	Long: `The 'save' command parses existing configuration and generates a copy in /conf/backup directory.

Examples:
  opnsense save               - saves current config as /conf/backup/config-<epoch_time>.xml
  opnsense save filename      - saves current config as /conf/backup/filename.xml
  opnsense save filename.xml  - saves current config as /conf/backup/filename.xml.
	  `,
	Run: func(cmd *cobra.Command, args []string) {
		filename := ""
		if len(args) < 1 {
			filename = generateBackupFilename()
		} else {
			filename = args[0]
			if !strings.HasSuffix(filename, ".xml") {
				filename += ".xml"
			}
			validFilenamePattern := "^[a-zA-Z0-9_.-]+$"
			match, err := regexp.MatchString(validFilenamePattern, filename)
			if err != nil || !match {
				internal.Log(1, "%s is not a valid file in /conf/backup.", filename)
			}
		}

		internal.Checkos()
		// check if filename already exists
		bash := `if [ -f "/conf/backup/` + filename + `" ]; then echo "exists"; fi`
		fileexists, err := internal.ExecuteCmd(bash, host)
		if err != nil {
			internal.Log(1, "execution error: %s", err.Error())
		}
		if strings.TrimSpace(fileexists) == "exists" {
			internal.Log(2, "%s already exists and will be overwritten.", filename)
			// delete the file
			bash = "sudo rm /conf/backup/" + filename
			_, err := internal.ExecuteCmd(bash, host)
			if err != nil {
				internal.Log(1, "execution error: %s", err.Error())
			}
		}
		// read and parse the config.xml file
		configdoc := etree.NewDocument()
		bash = "cat /conf/config.xml"
		config, err := internal.ExecuteCmd(bash, host)
		if err != nil {
			internal.Log(1, "execution error: %s", err.Error())
		}
		err = configdoc.ReadFromString(config)
		if err != nil {
			internal.Log(1, "could not parse /conf/config.xml")
		}

		configout := internal.ConfigToXML(configdoc, "opnsense")

		// chunking the long config.xml to upload in pieces
		chunkSize := 200000
		totalLength := len(configout)
		for i := 0; i < totalLength; i += chunkSize {
			end := i + chunkSize
			if end > totalLength {
				end = totalLength
			}
			chunk := configout[i:end]
			//escapedChunk := strconv.Quote(chunk) // This will escape necessary characters
			bash = fmt.Sprintf(`echo -n '%s' | sudo tee -a /conf/backup/%s`, chunk, filename)
			_, err = internal.ExecuteCmd(bash, host)
			if err != nil {
				internal.Log(1, "ssh execution error: %s", err.Error())
			}
		}

		// check that file was made
		bash = `if [ -f "/conf/backup/` + filename + `" ]; then echo "exists"; fi`
		fileexists, err = internal.ExecuteCmd(bash, host)
		if err != nil {
			internal.Log(1, "execution error: %s", err.Error())
		}

		if fileexists == "exists" {
			fmt.Printf("%s has been succesfully saved to /conf/backup.", filename)
		} else {
			internal.Log(1, "error writing file %s", filename)
		}

	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
}

func generateBackupFilename() string {
	timestamp := time.Now().Unix()
	randomNumber := rand.Intn(10000)
	filename := fmt.Sprintf("config-%d.%04d.xml", timestamp, randomNumber)
	return filename
}
