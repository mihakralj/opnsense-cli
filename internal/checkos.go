package internal

import (
	"strings"
)

func Checkos() (string, error) {
	//check that the target is OPNsense
	osstr := ExecuteCmd("echo $(uname; opnsense-version -N)", host)
	osstr = strings.TrimSpace(osstr)
	if osstr != "FreeBSD OPNsense" {
		Log(1, "%s is not OPNsense system", osstr)
	}
	Log(4, "OPNsense detected")
	return osstr, nil
}
