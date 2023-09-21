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
var backupsCmd = &cobra.Command{
	Use:   "backups",
	Short: `List available backup configurations in '/conf/backup' directory`,
	Long: `The 'backups' command provides functionalities for managing and viewing backup XML configurations within your OPNsense firewall system. You can list all backup configurations or get details about a specific one.`,
	Example:`  show backup           Lists all backup XML configurations.
  show backup <config>  Show details of a specific backup XML configuration`,
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
	backupsCmd.Flags().IntVarP(&depth, "depth", "d", 1, "Specify the depth of shown hierarchy (Default: 1)")
	rootCmd.AddCommand(backupsCmd)
}
