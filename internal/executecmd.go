package internal

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/sftp"
)

func ExecuteCmd(command, host string) string {
	if host == "" {
		Log(5, "local shell: %s", command)
		out, err := exec.Command("sh", "-c", command).Output()
		if err != nil {
			Log(1, "failed to execute command: %s %s",  command, err.Error())
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

func sftpCmd(data, filename, host string) {
	var sftpClient *sftp.Client
	var err error

	if host == "" {
		// If host is empty, save the data to a local file
		err = os.WriteFile(filename, []byte(data), 0644)
		if err != nil {
			Log(1, "Failed to save data to local file. %s", err.Error())
		}
		Log(4, "Successfully saved data to local file %s", filename)
		return
	}

	sshClient, err := getSSHClient(host)
	if err != nil {
		Log(1, "failed to initiate ssh client. %s", err.Error())
	}

	if sshClient.Client == nil {
		Log(1, "SSH client is nil. Cannot perform SFTP operation.")
	}

	// Create an SFTP client
	sftpClient, err = sftp.NewClient(sshClient.Client)
	if err != nil {
		Log(1, "Failed to initiate SFTP client. %s", err.Error())
	}
	defer sftpClient.Close()

	// Create remote file
	remoteFile, err := sftpClient.Create(filename)
	if err != nil {
		Log(1, "Failed to create remote file. %s", err.Error())
	}
	defer remoteFile.Close()

	// Write data to remote file
	_, err = remoteFile.Write([]byte(data))
	if err != nil {
		Log(1, "Failed to write to remote file. %s", err.Error())
	}

	Log(4, "Successfully transferred data to %s on host %s", filename, host)
}
