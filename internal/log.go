package internal

import (
	"fmt"
	"os"
)

var (
	verbose int
	host string
	configfile string
)

func SetFlags(v int, h string, c string) {
	if v < 1 || v > 5 {
		Log(1,"invalid verbosity level %d. It should be between 1 and 5", v)
	}
	verbose = v
	host = h
	configfile = c
	Log(5,"flags:\tverbose=%d, host=%s, config=%s",verbose,host,configfile)
}

func Log(verbosity int, format string, args ...interface{}) {
	levels := []string{"", "Error:\t", "Warn:\t", "Info:\t", "Note:\t", "Debug:\t"}
    message := levels[verbosity]+fmt.Sprintf(format, args...)

    if verbose >= verbosity || verbosity==1{
        fmt.Println(message)
    }
	if verbosity == 1 {
		os.Exit(1)
	}


}
