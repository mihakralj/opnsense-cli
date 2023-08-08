/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// infoCmd represents the status command
var infoCmd = &cobra.Command{
	Use:   "info [node]",
	Short: "Retrieves system information",
	Long: `
Info command provides a comprehensive overview of the OPNsense's current state.
The output is divided into multiple branches, each offering details about different aspects of the system:

- hardware: Presents details about the system's hardware, including the CPU, memory and recognized disks.
- system: Contains information about OPNSense, OS, its release version, boot time, and uptime.
- storage: Displays the configurations of available storage devices and any associated zpools.
- network: Lists all network interfaces available on the system.

you can use xpath to dive deeper into the result tree:
opnsense info hardware
opnsense info storage/disk0
opnsense info network/igb0/mtu
`,
	Run: func(cmd *cobra.Command, args []string) {

		path := "system"
		if len(args) >= 1 {
			trimmedArg := strings.Trim(args[0], "/")
			if trimmedArg != "" {
				path = trimmedArg
			}
			parts := strings.Split(path, "/")
			if parts[0] != "system" {
				path = "system/" + path
			}
		}
		internal.Checkos()
		bash := `echo -e "<system>\n<hardware>" && sysctl hw.model hw.machine_arch hw.machine hw.clockrate hw.ncpu kern.smp.cpus hw.realmem hw.physmem hw.usermem kern.disks | awk -F: '{ gsub(/^hw\./, "", $1); gsub(/^kern\./, "", $1); content = substr($2, 2); if ($1 ~ /mem$/) { content = sprintf("%.2fGB", content/1073741824) }; printf "<%s>%s</%s>\n", $1, content, $1 }' && echo -e "</hardware>"
		echo "<os>"&&opnsense-version -O | sed -n -e '/{/d' -e '/}/d' -e 's/^[[:space:]]*"product_\([^"]*\)":[[:space:]]*"\([^"]*\)".*/<\1>\2<\/\1>/p'&&sysctl kern.ostype kern.osrelease kern.version kern.hostname kern.hostid kern.hostuuid | awk 'NR>1 && !/^kern.version/ && !NF {next} {print}' | awk -F: '{ gsub(/^kern\./, "", $1); printf "<%s>%s</%s>\n", $1, substr($2, 2), $1 }'&&epochtime=$(sysctl kern.boottime | awk '{print $5}' | tr -d ','); now=$(date "+%s"); diff=$((now - epochtime)); days=$((diff / 86400)); hours=$(( (diff % 86400) / 3600)); minutes=$(( (diff % 3600) / 60)); seconds=$((diff % 60)); boottime=$(date -j -r "$epochtime" "+%Y-%m-%d %H:%M:%S"); printf "<boottime>%s</boottime>\n<boottime_epoch>%s</boottime_epoch>\n<uptime>%dd %dh %dm %ds</uptime>\n<uptime_seconds>%s</uptime_seconds>\n" "$boottime" "$epochtime" "$days" "$hours" "$minutes" "$seconds" "$diff"&&echo "</os>"&&echo -e "<storage>"&&mount | awk '{print $1}' | grep '^/dev/' | sort | uniq | xargs df -h | awk 'NR>1 {print}' | awk -v OFS='\t' -v disk_num=0 '{split($1, arr, "/"); if (arr[3] == "gpt") arr[3] = "rootfs"; print "<disk" disk_num ">\n\t<name>" arr[3] "</name>\n\t<type>gpt</type>\n\t<size>" $2 "</size>\n\t<used>" $3 "</used>\n\t<free>" $4 "</free>\n\t<capacity>" $5 "</capacity>\n\t</disk" disk_num ">"; disk_num++}'&&zpool list | awk 'NR>1 {print}' | awk -v OFS='\t' -v pool_num=0 '{print "<zpool" pool_num ">\n\t<name>" $1 "</name>\n\t<type>zfs</type>\n\t<size>" $2 "</size>\n\t<used>" $3 "</used>\n\t<free>" $4 "</free>\n\t<capacity>" $8 "</capacity>\n</zpool" pool_num ">";pool_num++}'&&echo -e"</storage>\n<network>" && ifconfig -a | sed -E 's/metric ([0-9]+)/\n metric: \1/;s/mtu ([0-9]+)/\n mtu: \1/' | sed -E 's/=/: /g; s/<([^>]*)>/ (\1)/g; s/nd6 options/nd6_options/g; s/^([a-zA-Z0-9]+) /\1: /; s/^([[:space:]]+)([a-zA-Z0-9]+)([ \t])/\1\2:\3/' | awk 'BEGIN {ORS=""} /^[a-zA-Z0-9]+: / { if (iface) print "</" iface ">"; iface=$1; sub(/:$/, "", iface); print "\n<" iface ">"; next } { sub(/:$/, "", $1); key=$1; $1=""; gsub(/^ /, ""); printf "\n\t<%s>%s</%s>", key, $0, key } END { if (iface) print "\n</" iface ">"; }' && echo -e '\n</network>'
		echo -e "</system>"`
		config, err := internal.ExecuteCmd(bash, host)
		if err != nil {
			panic(err)
		}
		configdoc := etree.NewDocument()
		configdoc.ReadFromString(config)
		configtty := internal.ConfigToTTY(configdoc, path, depth)
		fmt.Println(configtty)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
