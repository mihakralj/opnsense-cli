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
	"bytes"
	"os/exec"
	"strings"
)

func ExecuteCmd(command, host string) string {
	if host == "" {
		Log(5, "local shell: %s", command)
		out, err := exec.Command("sh", "-c", command).Output()
		if err != nil {
			Log(1, "failed to execute command: %s %s", command, err.Error())
		}
		Log(5, "received from local shell: %s", out)

		return string(out)
	}

	sshClient, err := getSSHClient(host)
	if err != nil {
		Log(1, "failed to initiate ssh client. %s", err.Error())
	}

	if sshClient.Session == nil {
		session, err := sshClient.Client.NewSession()
		if err != nil {
			Log(1, "failed to initiate ssh session. %s", err.Error())
		}
		var stdout bytes.Buffer
		session.Stdout = &stdout
		sshClient.Session = session
	}

	var stdoutBuf bytes.Buffer
	sshClient.Session.Stdout = &stdoutBuf
	Log(5, "ssh: %s", command)
	err = sshClient.Session.Run(command)
	if err != nil {
		Log(1, "failed to execute ssh command: %s", err.Error())
	}
	Log(5, "received from ssh: %s", strings.TrimRight(stdoutBuf.String(), "\n"))
	return strings.TrimRight(stdoutBuf.String(), "\n")
}
