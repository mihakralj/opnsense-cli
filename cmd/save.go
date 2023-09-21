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

	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save [filename]",
	Short: "Create a new backup XML configuration in '/conf/backup' directory",
	Long: `The 'save' command generates a new backup of the existing configuration, storing it in the '/conf/backup' directory. You can specify a filename for the backup, or if no filename is provided, the system will generate a default name based on the current epoch time.`,
	Example: `  opnsense save                 Save current config as '/conf/backup/config-<epoch_time>.xml'
  opnsense save filename.xml    Save current config as '/conf/backup/filename.xml'`,

	Run: func(cmd *cobra.Command, args []string) {
		filename := ""
		if len(args) < 1 {
			filename = generateBackupFilename()
		} else {
			filename = args[0]
			filename = strings.TrimPrefix(filename, "/conf/backup/")
			filename = strings.TrimPrefix(filename, "conf/backup/")
			if !strings.HasSuffix(filename, ".xml") {
				filename += ".xml"
			}
			validFilenamePattern := "^[a-zA-Z0-9_.-]+$"
			match, err := regexp.MatchString(validFilenamePattern, filename)
			if err != nil || !match {
				internal.Log(1, "%s is not a valid filename to save in /conf/backup.", filename)
			}
		}

		filename = "/conf/backup/"+filename
		internal.Checkos()
		configdoc := internal.LoadXMLFile(configfile, host)
		if configdoc == nil {
			internal.Log(1,"failed to get data from %s",configfile)
		}
		internal.SaveXMLFile(filename, configdoc, host, false)
		fmt.Printf("%s saved to %s\n", configfile, filename)
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
