/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/beevik/etree"
	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Lists available backup configurations or a specific backup",
	Long: `The 'backup' command allows you to view the available backup configurations in the OPNsense system or retrieve details of a specific backup. It's an essential tool for managing and understanding the saved configurations within your firewall system.

Example usage:
- show backup: Lists all available backup configurations.
- show backup <config>: Displays details of a specific backup configuration identified by <config>.`,
	Run: func(cmd *cobra.Command, args []string) {

		backupdir := "/conf/backup/"
		path := "backups"
		internal.Checkos()
		backupdoc := etree.NewDocument()
		bash := fmt.Sprintf(`count=$(find %s -type f | wc -l | sed -e 's/^[[:space:]]*//') && echo "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<backups count=\"$count\">" &&
		find /conf/backup -type f -exec sh -c 'echo $(stat -f "%%m" "$1") $(basename "$1") $(stat -f "%%z" "$1") $(md5sum "$1")' sh {} \; | sort -nr -k1 |
		awk '{
			date = strftime("%%Y-%%m-%%dT%%H:%%M:%%S", $1)
			delta = systime() - $1;
			days = int(delta / 86400);
			hours = int((delta %% 86400) / 3600);
			minutes = int((delta %% 3600) / 60);
			seconds = int(delta %% 60);
			age = days "d " hours "h " minutes "m " seconds "s";
			print "  <" $2 " age=\"" age "\"><date>" date "</date><size>" $3 "</size><md5>" $4 "</md5></" $2 ">";}
		END { print "</backups>"; }'`, backupdir)

		backups := internal.ExecuteCmd(bash, host)

		err := backupdoc.ReadFromString(backups)
		if err != nil {
			internal.Log(1, "did not receive XML")
		}

		backupout := ""
		if xmlFlag {
			backupout = internal.ConfigToXML(backupdoc, path)
		} else if jsonFlag {
			backupout = internal.ConfigToJSON(backupdoc, path)
		} else if yamlFlag {
			backupout = internal.ConfigToJSON(backupdoc, path)
		} else {
			backupout = internal.ConfigToTTY(backupdoc, path)
		}

		fmt.Println(backupout)

	},
}

func init() {
	backupCmd.Flags().IntVarP(&depth, "depth", "d", 1, "Specifies number of levels of returned tree (1-5)")
	rootCmd.AddCommand(backupCmd)
}
