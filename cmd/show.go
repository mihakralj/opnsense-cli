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
	"strings"

	"github.com/beevik/etree"
	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Displays information about configuration stored in config.xml",
	Long: `The show command displays configuration elements in config.xml, including interfaces, routes, firewall rules, and other system settings.
	 Use this command to view the current system configuration and troubleshoot issues.`,
	Run: func(cmd *cobra.Command, args []string) {

		path := "opnsense"
		if len(args) >= 1 {
			trimmedArg := strings.Trim(args[0], "/")
			if matched, _ := regexp.MatchString(`\[0\]`, trimmedArg); matched {
				internal.Log(1, "XPath indexing of elements starts with 1, not 0")
			}
			if trimmedArg != "" {
				path = trimmedArg
			}
			parts := strings.Split(path, "/")
			if parts[0] != "opnsense" {
				path = "opnsense/" + path
			}
		}
		bash := ""
		//internal.Checkos()
		configdoc := etree.NewDocument()
		bash = fmt.Sprintf("cat %s", configfile)
		config, err := internal.ExecuteCmd(bash, host)
		if err != nil {
			internal.Log(1, "execution error: %s", err.Error())
		}
		err = configdoc.ReadFromString(config)
		if err != nil {
			internal.Log(1, "%s is not an XML", configfile)
		}

		stagingdoc := etree.NewDocument()
		bash = fmt.Sprintf("if [ -f %s ]; then cat %s; else cat %s; fi", stagingfile, stagingfile, configfile)
		staging, err := internal.ExecuteCmd(bash, host)
		if err != nil {
			internal.Log(1, "execution error: %s", err.Error())
		}
		err = stagingdoc.ReadFromString(staging)
		if err != nil {
			internal.Log(1, "%s is not an XML", stagingfile)
		}

		if false {
			fmt.Println(stagingdoc)
		}

		configout := ""
		if xmlFlag {
			configout = internal.ConfigToXML(configdoc, path)
		} else if jsonFlag {
			configout = internal.ConfigToJSON(configdoc, path)
		} else if yamlFlag {
			configout = internal.ConfigToJSON(configdoc, path)
		} else {
			configout = internal.ConfigToTTY(configdoc, path)
		}

		fmt.Println(configout)

	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
