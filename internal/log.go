package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var c = map[string]string{
	"tag": "\033[0m",
	"txt": "\033[36m",
	"atr": "\033[35m",
	"new": "\033[32m",

	"red":  "\033[31m",
	"grn":  "\033[32m",
	"yel":  "\033[33m",
	"blu":  "\033[34m",
	"mgn":  "\033[35m",
	"cyn":  "\033[36m",
	"wht":  "\033[37m",
	"gry":  "\033[90m",
	"dred": "\033[2m\033[31m",
	"dgrn": "\033[2m\033[32m",
	"dyel": "\033[2m\033[33m",
	"dblu": "\033[2m\033[34m",
	"dmgn": "\033[2m\033[35m",
	"dcyn": "\033[2m\033[36m",
	"dwht": "\033[2m\033[37m",
	"dgry": "\033[2m\033[90m",

	"bred": "\033[91m",
	"bgrn": "\033[92m",
	"byel": "\033[93m",
	"bblu": "\033[94m",
	"bmgn": "\033[95m",
	"bcyn": "\033[96m",
	"bwht": "\033[97m",

	"bgr": "\033[41m",
	"bgg": "\033[42m",
	"bgy": "\033[43m",
	"bgb": "\033[44m",
	"bgm": "\033[45m",
	"bgc": "\033[46m",
	"bgw": "\033[47m",
	"ita": "\033[3m", // italics
	"bld": "\033[1m", // bold
	"stk": "\033[9m", // strikethroough
	"und": "\033[4m",
	"rev": "\033[7m", // reverse colors

	"ell": "\u2026",
	"arw": " \u2192 ",
	"nil": "\033[0m",
}

func Log(verbosity int, format string, args ...interface{}) {
	levels := []string{"",
		c["red"] + "Error:\t " + c["nil"],
		c["yel"] + "Warning: " + c["nil"],
		c["grn"] + "Info:\t " + c["nil"],
		c["blu"] + "Note:\t " + c["nil"],
		c["wht"] + "Debug:\t " + c["nil"]}

		formatted := fmt.Sprintf(format, args...)
		if len(formatted) > 210 {
			formatted = formatted[:100] + "\n...\n" + formatted[len(formatted)-100:]
		}
		message := levels[verbosity] + formatted

	if (verbose >= verbosity || verbosity == 1) && verbosity != 2 {
		fmt.Fprintln(os.Stderr, message)
	}
	if verbosity == 2 && !force {
		fmt.Print(message + "\nAre you sure? (Y/N): ")
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
