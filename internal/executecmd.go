package internal

import (
	"bytes"
	"fmt"
	"os/exec"
)

func ExecuteCmd(command, host string) (string, error) {

	if host == "" {
		out, err := exec.Command("sh", "-c", command).Output()
		if err != nil {
			return "", fmt.Errorf("failed to execute command locally: %v", err)
		}
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

	Log(4, "%s", command) //info level

	err = session.Run(command)

	if err != nil {
		return "", fmt.Errorf("failed to run command: %v", err)
	}

	return stdoutBuf.String(), nil
}
