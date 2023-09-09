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

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set <xpath> <value>",
	Short: "Sets a value for a specific node in the staging.xml file.",
	Long: `The 'set' command allows you to assign a new value to a specific node within the staging.xml file, effectively modifying the configuration in a controlled manner.

Before the changes can take effect, you need to use the 'commit' command to move the staging.xml file to the active config.xml. If at any point you decide to discard the changes made, you can use the 'discard' command.

The XPath parameter allows for precise targeting of the nodes in the XML structure, helping you to navigate to the exact setting or property that you wish to update.

Examples:
  opnsense set interfaces/wan/if igb0   - sets the 'interfaces/wan/if' node with the value 'igb0'
  opnsense set system/hostname myrouter - assigns 'myrouter' as the new hostname in the staging.xml file.

Make sure to validate your XPath expressions to avoid any unintended changes.`,
	Run: func(cmd *cobra.Command, args []string) {

		cmd.Help()

		fmt.Println("\n\033[34mSet command is not implemented yet\033[0m")

		// capture args
		//read and parse staging doc
		// modify the branch according to args
		// write to staging.xml
		// check that staging.xml was written

		/*
			internal.Checkos()

			//read and parse staging doc
			stagingdoc := etree.NewDocument()
			bash := fmt.Sprintf(`if [ -f %s ]; then cat %s; else cat %s; fi`, stagingfile, stagingfile, configfile)
			staging := internal.ExecuteCmd(bash, host)
			err := stagingdoc.ReadFromString(staging)
			if err != nil {
				internal.Log(1, "%s is not an XML", stagingfile)
			}

			// modify the branch according to args

			// write/append to staging.xml
			internal.ExecuteCmd(fmt.Sprintf(`sudo rm -fv %s`, stagingfile), host)
			stagingout := internal.ConfigToXML(stagingdoc, "opnsense")
			chunkSize := 200000
			totalLength := len(stagingout)
			for i := 0; i < totalLength; i += chunkSize {
				end := i + chunkSize
				if end > totalLength {
					end = totalLength
				}
				chunk := stagingout[i:end]
				bash = fmt.Sprintf(`echo -n '%s' | sudo tee -a %s`, chunk, stagingfile)
				internal.ExecuteCmd(bash, host)
			}

			// check that staging.xml was written
			bash = `if [ -f "` + stagingfile + `" ]; then echo "exists"; fi`
			fileexists := internal.ExecuteCmd(bash, host)
			if fileexists == "exists" {
				fmt.Printf("%s has been succesfully saved. ", stagingfile)
			} else {
				internal.Log(1, "error writing file %s", stagingfile)
			}
		*/
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
