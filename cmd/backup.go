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
package cmd

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: `List available backup configurations in '/conf/backup' directory`,
	Long:  `The 'backup' command provides functionalities for managing and viewing backup XML configurations within your OPNsense firewall system. You can list all backup configurations or get details about a specific one.`,
	Example: `  show backup           Lists all backup XML configurations.
  show backup <config>  Show details of a specific backup XML configuration`,
	Run: func(cmd *cobra.Command, args []string) {

		backupdir := "/conf/backup/"
		path := "backups"
		filename := ""

		if len(args) > 0 {
			filename = strings.TrimPrefix(args[0], "/")
			if !strings.HasSuffix(filename, ".xml") {
				filename = filename + ".xml"
			}
			path = path + "/" + filename
		}

		internal.Checkos()
		rootdoc := etree.NewDocument()

		bash := fmt.Sprintf(`echo -n '<?xml version="1.0" encoding="UTF-8"?>' && echo -n '<backups count="' && find %s -type f | wc -l | awk '{$1=$1};1' | tr -d '\n' && echo -n '">' | sed 's/##/"/g'`, backupdir)
		bash = bash + fmt.Sprintf(` && find %s -type f -exec sh -c 'echo $(stat -f "%%m" "$1") $(basename "$1") $(stat -f "%%z" "$1") $(md5sum "$1")' sh {} \; | sort -nr -k1`, backupdir)
		bash = bash + `| awk '{ date = strftime("%Y-%m-%dT%H:%M:%S", $1); delta = systime() - $1; days = int(delta / 86400); hours = int((delta % 86400) / 3600); minutes = int((delta % 3600) / 60); seconds = int(delta % 60); age = days "d " hours "h " minutes "m " seconds "s"; print "  <" $2 " age=\"" age "\"><date>" date "</date><size>" $3 "</size><md5>" $4 "</md5></" $2 ">"; } END { print "</backups>"; }'`

		backups := internal.ExecuteCmd(bash, host)
		err := rootdoc.ReadFromString(backups)
		if err != nil {
			internal.Log(1, "format is not XML")
		}
		if len(args) > 0 {
			internal.FullDepth()

			configdoc := internal.LoadXMLFile(configfile, host, false)
			backupdoc := internal.LoadXMLFile(backupdir+filename, host, false)

			deltadoc := internal.DiffXML(backupdoc, configdoc, false)


			// append all differences to the rootdoc
			diffEl := rootdoc.FindElement(path).CreateElement("diff")
			for _, child := range deltadoc.Root().ChildElements() {
				diffEl.AddChild(child.Copy())
			}

		}
		fmt.Println()
		internal.PrintDocument(rootdoc, path)

	},
}

func init() {
	backupCmd.Flags().IntVarP(&depth, "depth", "d", 1, "Specifies number of depth levels of returned tree (default: 1)")
	rootCmd.AddCommand(backupCmd)
}
