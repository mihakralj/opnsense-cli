package internal

import (
	"fmt"
	"strings"

)

func Checkos() (string, error) {
	//check that the target is OPNsense
	osstr, err := ExecuteCmd("echo $(uname; opnsense-version -N)", host)
	if err != nil {
		Log(1,"OPNsense not detected")
		//return "", fmt.Errorf("failed to execute command: %v", err)
	}
	osstr = strings.TrimSpace(osstr)
	if osstr != "FreeBSD OPNsense" {
		return "", fmt.Errorf("the target system is not FreeBSD OPNsense")
	}
	Log(4,"OPNsense detected")
	return osstr, nil
}
