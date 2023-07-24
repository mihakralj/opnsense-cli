package cmd

import (
	"fmt"
	"os"
    "strings"
	"github.com/spf13/cobra"
    "github.com/mihakralj/opnsense/internal"
)

var (
	user      string
	host      string
	port      string
)

func init() {
	cobra.OnInitialize(func() {
        internal.SetSSHTarget(user, host, port)
        
        //check that the target is OPNsense
        osstr, _ := internal.ExecuteCmd("uname", internal.SSHTarget)
        osstr = strings.TrimSpace(osstr)
        if osstr != "FreeBSD" {
            fmt.Println("The target system is not FreeBSD")
            os.Exit(1)
        }
        opn, _ := internal.ExecuteCmd("opnsense-version -N", internal.SSHTarget)
        opn = strings.TrimSpace(opn)
        if opn != "OPNsense" {
            fmt.Println("The target system is not OPNsense")
            os.Exit(1)
        }
        fmt.Println(osstr, opn)

		//check that the target is OPNsense
	})
	rootCmd.PersistentFlags().StringVarP(&host, "target", "t", "", "Target hostname for SSH")
	rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "admin", "Username for SSH")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "22", "Port for SSH")

}

var rootCmd = &cobra.Command{
	Use:   "opnsense",
	Short: "opnsense - command line ",
	Long: `opnsense is a super fancy CLI (kidding)

One can use opnsense to inspect opnsense configuration straight from the terminal`,
	Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("hello root command")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
