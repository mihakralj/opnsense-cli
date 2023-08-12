package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/beevik/etree"
	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

var actionCmd = &cobra.Command{
	Use:   "action",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.`,
	Run: func(cmd *cobra.Command, args []string) {

		path := "actions"
		trimmedArg := ""
		if len(args) >= 1 {
			trimmedArg = args[0]
			if len(args) > 1 {
				trimmedArg = trimmedArg+"/"+args[1]
			}
			if len(args) > 2 {
				trimmedArg += "."+args[2]
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
		config, err := internal.ExecuteCmd(bash, host)
		if err != nil {
			panic(err)
		}
		configdoc := etree.NewDocument()
		configdoc.ReadFromString(config)
		node := configdoc.FindElement(path + "/command")

		if verbose > 2 || node == nil{
			configtty := internal.ConfigToTTY(configdoc, path)
			fmt.Println(configtty)
		}

		if node != nil {
			path = strings.Replace(path, "actions/", "", 1)
			command := "configctl " + regexp.MustCompile(`[/\.]`).ReplaceAllString(path, " ")
			internal.Log(2,"sending command: %s ",command)
			ret, err := internal.ExecuteCmd(command, host)
			if err != nil {
				panic(err)
			}
			fmt.Println(ret)
		}

	},
}

func init() {
	rootCmd.AddCommand(actionCmd)
	// Here you will define your flags and configuration settings.
}
