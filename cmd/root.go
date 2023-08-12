package cmd

import (
	"fmt"
	"os"

	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

var (
	Version    string
	verbose    int
	force      bool
	host       string
	configfile string
	nocolor    bool
	depth      int
	xmlFlag	   bool
	yamlFlag   bool
	jsonFlag   bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&host, "target", "t", "", "Target host (-t user@hostname[:port])")
	rootCmd.PersistentFlags().IntVarP(&verbose, "verbose", "v", 1, "Set verbosity level (1-5)")
	rootCmd.PersistentFlags().BoolVarP(&nocolor, "no-color", "n", false, "Turn off ANSI colored output")
	rootCmd.PersistentFlags().IntVarP(&depth, "depth", "d", 1, "Specifies number of levels of returned tree (1-5)")
	rootCmd.PersistentFlags().BoolVar(&xmlFlag, "xml", false, "Output in XML format")
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&yamlFlag, "yaml", false, "Output in YAML format")
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Accept or bypass checks and prompts")
	//rootCmd.PersistentFlags().StringVarP(&configfile, "config", "c", "/conf/config.xml", "path to target config.xml")


	cobra.OnInitialize(func() {
		configfile = "/conf/config.xml"
		internal.SetFlags(verbose, force, host, configfile, nocolor, depth, xmlFlag, yamlFlag, jsonFlag)
		//other initializations
	})

}

var rootCmd = &cobra.Command{
	Use:   "opnsense",
	Short: "opnsense is a CLI to manage and monitor OPNsense firewall configuration, check status, change settings, and execute commands.",
	Long: `
Description:
  opnsense is a command-line utility for managing, configuring, and monitoring OPNsense firewall systems.
  It facilitates non-GUI administration, both directly in the shell and remotely via an SSH tunnel.
  All interactions with OPNsense utilize the same mechanisms as the Web GUI,
  including staged modifications of config.xml and execution of available configd commands.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Println(cmd.Long)
		}
	},
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.DisableNoDescFlag = true
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing CLI '%s'", err)
		os.Exit(1)
	}
}
