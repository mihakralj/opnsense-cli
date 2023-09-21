package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/beevik/etree"
	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [service] [command] [parameters]",
	Short: "Execute registered commands on OPNsense firewall",
	Long: `The 'run' command allows you to list and execute specific commands registered with the 'configctl' utility on the OPNsense firewall system. You can run various service-specific commands and other operational tasks.`,
	Example: `  opnsense run                       List all available configd services
  opnsense run dns                   List available commands for DNS service
  opnsense run dhcpd list leases     Show DHCP leases
  opnsense run interface flush arp   Flush ARP table
  opnsense run firmware reboot       Initiate a system reboot`,


	Run: func(cmd *cobra.Command, args []string) {

		path := "actions"

		// need better parsing of tokens, so the command can be recognized better and parameters passed
		trimmedArg := ""
		if len(args) >= 1 {
			trimmedArg = args[0]
			if len(args) > 1 {
				trimmedArg = trimmedArg + "/" + args[1]
			}
			if len(args) > 2 {
				trimmedArg += "." + args[2]
			}
			//trimmedArg = strings.Trim(args[0], "/")
			if trimmedArg != "" {
				path = trimmedArg
			}
			parts := strings.Split(path, "/")
			if parts[0] != "actions" {
				path = "actions/" + path
			}
		}
		internal.Checkos()
		bash := `echo "<actions>" && for file in /usr/local/opnsense/service/conf/actions.d/actions_*.conf; do service_name=$(basename "$file" | sed 's/actions_\(.*\).conf/\1/'); echo "  <${service_name}>"; awk 'function escape_xml(str) { gsub(/&/, "&amp;", str); gsub(/</, "&lt;", str); gsub(/>/, "&gt;", str); return str; } BEGIN {FS=":"; action = "";} /\[.*\]/ { if (action != "") {print "    </" action ">"} action = substr($0, 2, length($0) - 2); print "    <" action ">";} !/\[.*\]/ && NF > 1 { gsub(/^[ \t]+|[ \t]+$/, "", $2); value = escape_xml($2); print "      <" $1 ">" value "</" $1 ">";} END { if (action != "") {print "    </" action ">"} }' "$file"; echo "  </${service_name}>"; done && echo "</actions>"`
		config := internal.ExecuteCmd(bash, host)

		configdoc := etree.NewDocument()
		configdoc.ReadFromString(config)
		node := configdoc.FindElement(path + "/command")

		if !force {
			configtty := internal.ConfigToTTY(configdoc, path)
			fmt.Println(configtty)
		}

		if node != nil {
			path = strings.Replace(path, "actions/", "", 1)
			command := "configctl " + regexp.MustCompile(`[/\.]`).ReplaceAllString(path, " ")
			internal.Log(2, "sending command: %s ", command)
			ret := internal.ExecuteCmd(command, host)

			var js json.RawMessage
			if err := json.Unmarshal([]byte(ret), &js); err == nil {
				var obj interface{}
				if err := json.Unmarshal(js, &obj); err != nil {
					fmt.Println("Error unmarshaling JSON:", err)
					return
				}
				prettyJSON, err := json.MarshalIndent(obj, "", "  ")
				if err != nil {
					fmt.Println("Error formatting JSON:", err)
					return
				}
				fmt.Println(string(prettyJSON))
			} else {
				fmt.Println(ret)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
