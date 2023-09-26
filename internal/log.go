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
package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Log(verbosity int, format string, args ...interface{}) {
	levels := []string{"",
		c["red"] + "Error:\t " + c["nil"],
		c["yel"] + "Warning: " + c["nil"],
		c["grn"] + "Info:\t " + c["nil"],
		c["blu"] + "Note:\t " + c["nil"],
		c["wht"] + "Debug:\t " + c["nil"]}

	formatted := fmt.Sprintf(format, args...)
	if len(formatted) > 2000 {
		formatted = formatted[:1000] + "\n...\n" + formatted[len(formatted)-200:]
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
			fmt.Fprintln(os.Stderr, "action canceled")
			os.Exit(1)
		}
	}
	if verbosity == 1 {
		os.Exit(1)
	}
}
