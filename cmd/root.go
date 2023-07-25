package cmd

import (
	"fmt"
	"os"

	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

var (
	verbose int
	host string
	configfile string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&host, "target", "t", "", "Target host (-t user@hostname[:port])")
	rootCmd.PersistentFlags().StringVarP(&configfile, "config", "c", "/conf/config.xml", "path to target config.xml")
	rootCmd.PersistentFlags().IntVarP(&verbose, "verbose", "v", 1, "Set verbosity level (1-5)")

	cobra.OnInitialize(func() {
		internal.SetFlags(verbose, host, configfile)
		//other initializations
	})

}

var rootCmd = &cobra.Command{
	Use:   "opnsense",
	Short: "opnsense - command line ",
	Long: `opnsense is a super fancy CLI (kidding)

One can use opnsense to inspect opnsense configuration straight from the terminal`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
