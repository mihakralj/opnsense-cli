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

	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: `Commit changes from the 'staging.xml' to the active 'config.xml'`,
	Long: `The 'commit' command finalizes the staged changes made to the 'staging.xml' file, making them the active configuration for the OPNsense firewall system. This operation is the last step in a sequence that typically involves the 'set' and optionally 'discard' commands. The 'commit' action creates a backup of the active 'config.xml', moves 'staging.xml' to 'config.xml', and reloads the 'configd' service.
	`,

	Example: `  opnsense commit          Commit the changes in 'staging.xml' to become the active 'config.xml'
  opnsense commit --force  Commit the changes without requiring interactive confirmation.`,
	Run: func(cmd *cobra.Command, args []string) {

		// check if staging.xml exists
		internal.Checkos()
		bash := `test -f "` + stagingfile + `" && echo "exists" || echo "missing"`
		fileexists := internal.ExecuteCmd(bash, host)
		if strings.TrimSpace(fileexists) != "exists" {
			fmt.Println("no staging.xml detected - nothing to commit.")
			return
		}
		bash = `cmp -s "` + configfile + `" "` + stagingfile + `" && echo "same" || echo "diff"`
		filesame := internal.ExecuteCmd(bash, host)
		if strings.TrimSpace(filesame) != "diff" {
			fmt.Println("staging.xml and config.xml are the same - nothing to commit.")
		}

		configdoc := internal.LoadXMLFile(configfile, host, false)
		stagingdoc := internal.LoadXMLFile(stagingfile, host, false)

		deltadoc := internal.DiffXML(configdoc, stagingdoc, false)
		fmt.Println("\nChanges to be commited:")
		internal.PrintDocument(deltadoc, "opnsense")

		internal.Log(2, "commiting %s to %s and reloading all services", stagingfile, configfile)

		// copy config.xml to /conf/backup dir
		backupname := internal.GenerateBackupFilename()
		bash = `sudo cp -f ` + configfile + ` /conf/backup/` + backupname + ` && sudo mv -f /conf/staging.xml ` + configfile
		internal.ExecuteCmd(bash, host)

		include := "php -r \"require_once('/usr/local/etc/inc/config.inc'); require_once('/usr/local/etc/inc/interfaces.inc'); require_once('/usr/local/etc/inc/filter.inc'); require_once('/usr/local/etc/inc/auth.inc'); require_once('/usr/local/etc/inc/rrd.inc'); require_once('/usr/local/etc/inc/util.inc'); require_once('/usr/local/etc/inc/system.inc'); require_once('/usr/local/etc/inc/interfaces.inc'); "

		var result string

		result = internal.ExecuteCmd(include+"system_firmware_configure(true); \"", host)
		fmt.Println(result)
		result = internal.ExecuteCmd(include+"system_trust_configure(true); \"", host)
		fmt.Println(result)
		result = internal.ExecuteCmd(include+"system_login_configure(true); \"", host)
		fmt.Println(result)
		result = internal.ExecuteCmd(include+"system_cron_configure(true); \"", host)
		fmt.Println(result)
		result = internal.ExecuteCmd(include+"system_timezone_configure(true); \"", host)
		fmt.Println(result)
		result = internal.ExecuteCmd(include+"system_hostname_configure(true); \"", host)
		fmt.Println(result)
		result = internal.ExecuteCmd(include+"system_resolver_configure(true); \"", host)
		fmt.Println(result)
		result = internal.ExecuteCmd(include+"interfaces_configure(true); \"", host)
		fmt.Println(result)
		result = internal.ExecuteCmd(include+"system_routing_configure(true); \"", host)
		fmt.Println(result)
		result = internal.ExecuteCmd(include+"rrd_configure(true); \"", host)
		fmt.Println(result)
		result = internal.ExecuteCmd(include+"filter_configure(true); \"", host)
		fmt.Println(result)
		result = internal.ExecuteCmd(include+"plugins_configure('vpn', true); \"", host)
		fmt.Println(result)
		result = internal.ExecuteCmd(include+"plugins_configure('local', true); \"", host)
		fmt.Println(result)

	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
