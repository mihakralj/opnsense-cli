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
	"strings"
)

// Checkos checks that the target is an OPNsense system
func Checkos() (string, error) {
	//check that the target is OPNsense
	osstr := ExecuteCmd("echo `uname` `opnsense-version -N`", host)
	osstr = strings.TrimSpace(osstr)
	if osstr != "FreeBSD OPNsense" {
		Log(1, "%s is not OPNsense system", osstr)
	}
	Log(4, "OPNsense system detected")
	return osstr, nil
}
