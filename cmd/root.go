package cmd

import (
	"fmt"
	"os"

	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

var (
	Version     string = "0.11.0"
	verbose     int
	force       bool
	host        string
	configfile  string
	stagingfile string
	nocolor     bool
	depth       int = 1
	xmlFlag     bool
	yamlFlag    bool
	jsonFlag    bool
	ver_flag    bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&host, "target", "t", "", "Target host (-t user@hostname[:port])")
	rootCmd.PersistentFlags().IntVarP(&verbose, "verbose", "v", 1, "Set verbosity level (1-5)")
	rootCmd.PersistentFlags().BoolVarP(&nocolor, "no-color", "n", false, "Turn off ANSI colored output")
	//rootCmd.PersistentFlags().IntVarP(&depth, "depth", "d", 1, "Specifies number of levels of returned tree (1-5)")
	rootCmd.PersistentFlags().BoolVar(&xmlFlag, "xml", false, "Output in XML format")
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&yamlFlag, "yaml", false, "Output in YAML format")
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Accept or bypass checks and prompts")
	rootCmd.Flags().BoolVar(&ver_flag, "version", false, "display version of opnsense")

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	cobra.OnInitialize(func() {
		configfile = "/conf/config.xml"
		stagingfile = "/conf/staging.xml"
		internal.SetFlags(verbose, force, host, configfile, nocolor, depth, xmlFlag, yamlFlag, jsonFlag)
		//other initializations
	})
}

var rootCmd = &cobra.Command{
	Use:   "opnsense [command]",
	Short: "CLI to manage and monitor OPNsense firewall systems.",
	Long: `Command Line utility to interact with OPNsense firewall.

opnsense CLI is a command-line utility for managing, configuring, and monitoring OPNsense firewall systems.
It facilitates non-GUI administration, both locally on the firewall and remotely via an SSH tunnel.
To avoid entering passwords for each remote call, use 'ssh-add' to add private key to your ssh-agent.`,

	Example: `  opnsense show interfaces/wan       - Show the inerfaces/wan of config.xml in json format
  opnsense sysinfo                   - Show system information on remote OPNsense
  opnsense backup                    - Show backup files and their age
  opnsense run firmware reboot -f    - Reboot OPNsense, force (no confirmation)
  opnsense commit                    - Commit staged changes`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if ver_flag {
			fmt.Println("opnsense-CLI version", Version)
			os.Exit(0)
		}

		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
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
