package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
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

var c = map[string]string{
	"red": "\033[31m",
	"grn": "\033[32m",
	"yel": "\033[33m",
	"blu": "\033[34m",
	"mgn": "\033[35m",
	"cyn": "\033[36m",
	"wht": "\033[37m",
	"bgr": "\033[41m",
	"bgg": "\033[42m",
	"bgy": "\033[43m",
	"bgb": "\033[44m",
	"bgm": "\033[45m",
	"bgc": "\033[46m",
	"bgw": "\033[47m",
	"ita": "\033[3m",
	//"ell": "\u2026",
	"ell": "...",
	"nil": "\033[0m",
}

func SetFlags(v int, f bool, h string, config string, nc bool, dpt int, x bool, y bool, j bool) {
	if v < 1 || v > 5 {
		Log(1, "invalid verbosity level %d. It should be between 1 and 5", v)
	}
	verbose = v
	force = f
	host = h
	configfile = config
	nocolor = nc
	depth = dpt
	xmlFlag = x
	yamlFlag = y
	jsonFlag = j
	Log(5, "flags:\tverbose=%d, host=%s, config=%s", verbose, host, configfile)
	if nc {
		for key := range c {
			delete(c, key)
		}
		c["ell"] = "..."
	}
}

func Log(verbosity int, format string, args ...interface{}) {
	levels := []string{"",
		c["red"] + "Error:\t " + c["nil"],
		c["yel"] + "Warning: " + c["nil"],
		c["grn"] + "Info:\t " + c["nil"],
		c["blu"] + "Note:\t " + c["nil"],
		c["wht"] + "Debug:\t " + c["nil"]}
	message := levels[verbosity] + fmt.Sprintf(format, args...)

	if (verbose >= verbosity || verbosity == 1) && verbosity != 2{
		fmt.Fprintln(os.Stderr, message)
	}
	if verbosity == 2 && !force {
		fmt.Print(message+"\nAre you sure? (Y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			Log(1, "error reading input")
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return
		} else {
			Log(1, "canceled action")
		}
	}
	if verbosity == 1 {
		os.Exit(1)
	}

}
