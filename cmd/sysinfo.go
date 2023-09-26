/*
Copyright © 2023 Miha miha.kralj@outlook.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

// systemCmd represents the status command
var sysinfoCmd = &cobra.Command{
	Use:   "sysinfo [node]",
	Short: "Retrieve comprehensive system information",
	Long: `The 'sysinfo' command provides an extensive overview of your OPNsense firewall system, including hardware, operating system, storage, and network configurations. The output is organized into multiple branches, each containing details on various aspects of the system:`,
	Example: `  opnsense sysinfo hardware          Display hardware details
  opnsense sysinfo storage/disk0     Information about the first disk
  opnsense sysinfo network/igb0/mtu  Show the MTU for the igb0 network interface`,

	Run: func(cmd *cobra.Command, args []string) {
		if changed := cmd.Flags().Changed("depth"); !changed {
			depth = 2
			internal.SetFlags(verbose, force, host, configfile, nocolor, depth, xmlFlag, yamlFlag, jsonFlag)
		}

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
		echo "<os>" && opnsense-version -O | sed -n -e '/{/d' -e '/}/d' -e 's/^[[:space:]]*"product_\([^"]*\)":[[:space:]]*"\([^"]*\)".*/<\1>\2<\/\1>/p' && sysctl kern.ostype kern.osrelease kern.version kern.hostname kern.hostid kern.hostuuid | sed 's/^kern.version:/kern.osversion:/' | awk 'NR>1 && !/^kern.opnversion/ && !NF {next} {print}' | awk -F: '{ gsub(/^kern\./, "", $1); printf "<%s>%s</%s>\n", $1, substr($2, 2), $1 }' && epochtime=$(sysctl kern.boottime | awk '{print $5}' | tr -d ','); now=$(date "+%s"); diff=$((now - epochtime)); days=$((diff / 86400)); hours=$(( (diff % 86400) / 3600)); minutes=$(( (diff % 3600) / 60)); seconds=$((diff % 60)); boottime=$(date -j -r "$epochtime" "+%Y-%m-%d %H:%M:%S"); printf "<boottime>%s</boottime>\n<boottime_epoch>%s</boottime_epoch>\n<uptime>%dd %dh %dm %ds</uptime>\n<uptime_seconds>%s</uptime_seconds>\n" "$boottime" "$epochtime" "$days" "$hours" "$minutes" "$seconds" "$diff" && echo "</os>"
		echo -e "<storage>" && mount | awk '{print $1}' | grep '^/dev/' | sort | uniq | xargs df -h | awk 'NR>1 {print}' | awk -v OFS='\t' -v disk_num=0 '{split($1, arr, "/"); if (arr[3] == "gpt") arr[3] = "rootfs"; print "<disk" disk_num ">\n\t<name>" arr[3] "</name>\n\t<type>gpt</type>\n\t<size>" $2 "</size>\n\t<used>" $3 "</used>\n\t<free>" $4 "</free>\n\t<capacity>" $5 "</capacity>\n\t</disk" disk_num ">"; disk_num++}' && zpool list | awk 'NR>1 {print}' | awk -v OFS='\t' -v pool_num=0 '{print "<zpool" pool_num ">\n\t<name>" $1 "</name>\n\t<type>zfs</type>\n\t<size>" $2 "</size>\n\t<used>" $3 "</used>\n\t<free>" $4 "</free>\n\t<capacity>" $8 "</capacity>\n</zpool" pool_num ">";pool_num++}' && echo -e"</storage>\n<network>" && ifconfig -a | sed -E 's/metric ([0-9]+)/\n metric: \1/;s/mtu ([0-9]+)/\n mtu: \1/' | sed -E 's/=/: /g; s/<([^>]*)>/ (\1)/g; s/nd6 options/nd6_options/g; s/^([a-zA-Z0-9]+) /\1: /; s/^([[:space:]]+)([a-zA-Z0-9]+)([ \t])/\1\2:\3/' | awk 'BEGIN {ORS=""} /^[a-zA-Z0-9]+: / { if (iface) print "</" iface ">"; iface=$1; sub(/:$/, "", iface); print "\n<" iface ">"; next } { sub(/:$/, "", $1); key=$1; $1=""; gsub(/^ /, ""); printf "\n\t<%s>%s</%s>", key, $0, key } END { if (iface) print "\n</" iface ">"; }' && echo -e '\n</network>\n</system>'
		echo -e "</system>"`
		config := internal.ExecuteCmd(bash, host)

		configdoc := etree.NewDocument()
		configdoc.ReadFromString(config)

		configout := ""
		if xmlFlag {
			configout = internal.ConfigToXML(configdoc, path)
		} else if jsonFlag {
			configout = internal.ConfigToJSON(configdoc, path)
		} else if yamlFlag {
			configout = internal.ConfigToJSON(configdoc, path)
		} else {
			configout = internal.ConfigToTTY(configdoc, path)
		}
		fmt.Println(configout)
	},
}

func init() {
	sysinfoCmd.Flags().IntVarP(&depth, "depth", "d", 1, "Specifies number of levels of returned tree (1-5)")
	rootCmd.AddCommand(sysinfoCmd)
}
