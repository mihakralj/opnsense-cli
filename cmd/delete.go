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
	"regexp"
	"strconv"
	"strings"

	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [filename|command]",
	Short: "Deletes a backup configuration from /conf/backup",
	Long: `The 'delete' command allows you to delete a specific backup configuration from the OPNsense system. It is limited to /conf/backup directory only.

Examples:
  opnsense delete filename.xml  - Delete a specific file in /conf/backup.
  opnsense delete age 10        - Delete all files older than 10 days
  opnsense delete trim 10       - Delete the oldest 10 backup files
  opnsense delete keep 10       - Delete all backup files except the most recent 10
`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			// display long help
			cmd.Help()
			return
		}
		bash := ""
		filename := args[0]

		switch filename {
		case "age", "keep", "trim":
			if len(args) < 2 {
				internal.Log(1, "missing required value for %s",filename)
				return
			}

			value := args[1]
			_, err := strconv.Atoi(value) //value needs to be a number
			if err != nil {
				internal.Log(1, "%s is not a valid number",value)
			}

			if filename == "age" {
				bash = "find /conf/backup -type f -mtime +"+value
			}
			if filename == "keep" {
				bash = "find /conf/backup -type f -print0 | xargs -0 ls -lt | tail -n +"+value+" | awk '{print $NF}'"
			}
			if filename == "trim" {
				bash = "find /conf/backup -type f -print0 | xargs -0 ls -lt | tail -n "+value+" | awk '{print $NF}'"
			}

			ret := internal.ExecuteCmd(bash+" | wc -l", host)

			cnt := strings.TrimSpace(ret)
			if cnt=="0" {
				fmt.Println("no files meeting criteria")
				return
			}
			internal.Log(2, "deleting %s files from /conf/backup",cnt)
			ret = internal.ExecuteCmd(bash+" | sudo xargs rm", host)
			if ret == "" {
				fmt.Printf("%s files have been deleted.\n", cnt)
			}
		default:
			validFilenamePattern := "^[a-zA-Z0-9_.-]+$"
			match, err := regexp.MatchString(validFilenamePattern, filename)
			if err != nil || !match {
				internal.Log(1, "%s is not a valid file in /conf/backup.", filename)
			}
			internal.Checkos()
			bash = `if [ -e "/conf/backup/` + filename + `" ]; then echo "ok"; else echo "missing"; fi`
			fileexists := internal.ExecuteCmd(bash, host)

			if strings.TrimSpace(fileexists) == "missing" {
				internal.Log(1, "file %s not found", filename)
			}

			internal.Log(2, "deleting file %s", filename)

			bash = "sudo chmod a+w /conf/backup/" + filename + " && sudo rm -f /conf/backup/" + filename
			result := internal.ExecuteCmd(bash, host)

			if result == "" {
				fmt.Printf("%s has been deleted.\n", filename)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
