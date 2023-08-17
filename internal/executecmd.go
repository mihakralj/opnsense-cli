package internal

import (
	"bytes"
	"fmt"
	"os/exec"
)

func ExecuteCmd(command, host string) (string, error) {

	Log(4,"sh -c %s",command)
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

	session := sshClient.Session
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	err = session.Run(command)
	if err != nil {
		Log(3, "failed to execute sh command. %s", err.Error())
		return "", err
	}
	Log(5, "received results from executed command.")
	return stdoutBuf.String(), nil
}
