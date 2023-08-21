package internal

import (
	"bytes"
	"os/exec"
)

func ExecuteCmd(command, host string) (string, error) {
	Log(4, "sh -c %s", command)
	if host == "" {
		Log(5, "no target provided; executing command locally.")
		out, err := exec.Command("sh", "-c", command).Output()
		if err != nil {
			Log(3, "failed to execute sh command. %s", err.Error())
			return "", err
		}
		Log(5, "received results from executed command.")
		return string(out), nil
	}

	sshClient, err := getSSHClient(host)
	if err != nil {
		return "", err
	}

	if sshClient.Session == nil {
		session, err := sshClient.Client.NewSession()
		if err != nil {
			Log(1, "%s", err)
			return "", err
		}
		sshClient.Session = session
	}

	var stdoutBuf bytes.Buffer
	sshClient.Session.Stdout = &stdoutBuf
	err = sshClient.Session.Run(command)
	if err != nil {
		Log(3, "failed to execute sh command. %s", err.Error())
		return "", err
	}
	Log(5, "received results from executed command.")
	return stdoutBuf.String(), nil
}
