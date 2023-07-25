package cmd

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := "opnsense"
		if len(args) >= 1 {
			trimmedArg := strings.Trim(args[0], "/")
			if trimmedArg != "" {
				path = trimmedArg
			}
		}
		internal.Checkos()
		bash := "cat /conf/config.xml"
		output, err := internal.ExecuteCmd(bash, host)
		if err != nil {
			panic(err)
		}

		doc := etree.NewDocument()
		if err := doc.ReadFromString(output); err != nil {
			panic(err)
		}

		focused := internal.FocusTree(doc.Root(), path, 1)

		fmt.Println(internal.EtreeToYaml(focused, 0))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
