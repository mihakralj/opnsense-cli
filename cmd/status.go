package cmd

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := "interfaces"
		if len(args) >= 1 {
			trimmedArg := strings.Trim(args[0], "/")
			if trimmedArg != "" {
				path = trimmedArg
			}
		}

		internal.Checkos()
		bash := "ifconfig -a | sed -E 's/^([a-zA-Z0-9]*:)(.*)/\\1\\n       \\2/; s/=/\\: /g' | awk '{ if (NF > 0 && substr($1,length($1)) != \":\") {$1 = $1 \":\"; print \"        \" $0} else {print $0}}'"
		output, err := internal.ExecuteCmd(bash, host)
		if err != nil {
			panic(err)
		}

		doc := etree.NewDocument()
		root := doc.CreateElement("interfaces")
		lines := strings.Split(output, "\n")
		var currentElement *etree.Element
		for _, line := range lines {
			if strings.HasSuffix(line, ":") {
				// This is an interface name
				currentElement = root.CreateElement(strings.TrimSuffix(line, ":"))
			} else if currentElement != nil {
				// This is a property of the current interface
				parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
				if len(parts) == 2 {
					currentElement.CreateElement(parts[0]).SetText(parts[1])
				}
			}
		}
		focused := internal.FocusTree(doc.Root(), path, -1)

		fmt.Println(internal.EtreeToJSON(focused))

	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
