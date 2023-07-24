package cmd

import (
	"fmt"
	"os"
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
		if len(args) < 1 {
			fmt.Println("Please provide a path as a command-line argument.")
			os.Exit(1)
		}
		path := args[0]

		output, err := internal.ExecuteCmd("cat /conf/config.xml", internal.SSHTarget)
		if err != nil {
			panic(err)
		}

		doc := etree.NewDocument()
		if err := doc.ReadFromString(output); err != nil {
			panic(err)
		}

		focused := internal.FocusTree(doc.Root(), path, false)

		fmt.Println(internal.EtreeToYaml(focused, 0))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
