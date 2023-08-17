/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/beevik/etree"
	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// grab config.xml
		configdoc := etree.NewDocument()
		bash := fmt.Sprintf("cat %s", configfile)
		config, err := internal.ExecuteCmd(bash, host)
		if err != nil {
			internal.Log(1, "execution error: %s", err.Error())
		}
		err = configdoc.ReadFromString(config)
		if err != nil {
			internal.Log(1, "%s is not an XML", configfile)
		}
		configout := internal.ConfigToXML(configdoc, "")

		// generate backup name
		backupdir := "/conf/backup/"
		backupfile := generateBackupFilename()

		// generate hash

		// check for existence of backup folder

		// pour backup file into the backup folder
		bash = fmt.Sprintf("mkdir -p %s&&cat >%s%s <<EOF\n%s\nEOF", backupdir, backupdir, backupfile, configout)
		ret, err := internal.ExecuteCmd(bash, host)
		if err != nil {
			internal.Log(1, "execution error: %s", err.Error())
		}
		fmt.Println(ret)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func generateBackupFilename() string {
	timestamp := time.Now().Unix()
	randomNumber := rand.Intn(10000)
	filename := fmt.Sprintf("config-%d.%d.xml", timestamp, randomNumber)
	return filename
}
