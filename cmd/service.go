package cmd

import (
	"fmt"
	"strings"

	//"github.com/beevik/etree"
	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := "services"
		if len(args) >= 1 {
			trimmedArg := strings.Trim(args[0], "/")
			if trimmedArg != "" {
				path = trimmedArg
			}
		}

		internal.Checkos()
		bash := `for file in /usr/local/opnsense/service/conf/actions.d/actions_*.conf; do service=$(basename $file | awk -F'[_.]' '{print $2}'); echo "${service}:"; awk -F'[][]' '/^\[/{print "  " $2 ":"; next} {if($1 ~ /:/ && $1 !~ /: /){gsub(/:/,": ",$1)}; if(NF>0) print "    " $1}' $file; done`
		output, err := internal.ExecuteCmd(bash, host)
		if err != nil {
			panic(err)
		}
/*
		doc := etree.NewDocument()
		root := doc.CreateElement("opnsense")
		//merge multilines
		lines := strings.Split(output, "\n")
		var previousLine string
		var trimmedOutput []string

		for i, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if (strings.Contains(trimmedLine, ":") || strings.HasPrefix(trimmedLine, "command:")) && i != 0 {
				// Add the completed previous line to trimmedOutput
				trimmedOutput = append(trimmedOutput, previousLine)
				previousLine = line
			} else {
				// Attach this line to the previous line, but only if it's not empty
				if trimmedLine != "" {
					if strings.TrimSpace(previousLine) != "" {
						previousLine += "; " + trimmedLine
					} else {
						previousLine = trimmedLine
					}
				}
			}
		}
		// Add the last line to trimmedOutput
		if previousLine != "" {
			trimmedOutput = append(trimmedOutput, previousLine)
		}


		var serviceElem, commandElem *etree.Element

		for _, line := range trimmedOutput {
			trimmedLine := strings.TrimSpace(line)

			// Determine the level of the line based on the leading spaces
			switch strings.Count(line, "  ") {
			case 0: // Service level
				serviceElem = root.CreateElement(trimmedLine)
				commandElem = nil
			case 1: // Command level
				if serviceElem != nil {
					commandElem = serviceElem.CreateElement(trimmedLine)
				}
			case 2: // Parameter level
				if commandElem != nil {
					commandElem.CreateElement(trimmedLine)
				}
			}
		}

		focused := internal.FocusTree(doc.Root(), path, -1)
		fmt.Println(internal.EtreeToYaml(focused, 0))
		// Join the lines into a single string

		trimmedOutputStr := strings.Join(trimmedOutput, "\n")
		*/
		fmt.Println(path)
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
