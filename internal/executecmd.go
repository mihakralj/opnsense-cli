package internal

import (
	"bytes"
	"os/exec"
	"strings"
)

func ExecuteCmd(command, host string) (string) {
	Log(4, "sh -c %q", command)
	if host == "" {
		Log(5, "no target provided; executing command locally.")
		out, err := exec.Command("sh", "-c", command).Output()
		if err != nil {
			Log(1, "failed to execute shell command. %s", err.Error())
		}
		Log(5, "received results from local shell: %s", out)

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
	err = sshClient.Session.Run(command)
	if err != nil {
		Log(1, "failed to execute ssh command: %s", err.Error())
	}
	Log(5, "received results from ssh host: %s", strings.TrimRight(stdoutBuf.String(), "\n"))
	return strings.TrimRight(stdoutBuf.String(), "\n")
}
