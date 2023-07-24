package internal

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type SSHClient struct {
    Session *ssh.Session
}

var (
    SshClient *SSHClient
    SSHTarget string
)

func SetSSHTarget(user, host, port string) {
    if host != "" {
        SSHTarget = fmt.Sprintf("%s@%s:%s", user, host, port)
    }
}

func getSSHClient(user, addr string) (*SSHClient, error) {
    config := &ssh.ClientConfig{
        User:            user,
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }

    sshAgent, err := GetSSHAgent()
    if err == nil {
        config.Auth = []ssh.AuthMethod{sshAgent}
    }

    if len(config.Auth) == 0 {
        fmt.Println("No identities found in the SSH agent. Falling back to password authentication.")
        fmt.Printf("Enter Password for %s@%s: ", user, addr)
        bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
        fmt.Println()
        if err != nil {
            return nil, fmt.Errorf("failed to read password: %v", err)
        }
        password := string(bytePassword)
        config.Auth = []ssh.AuthMethod{ssh.Password(password)}
    }

    if !strings.Contains(addr, ":") {
        addr = addr + ":22"
    }

    connection, err := ssh.Dial("tcp", addr, config)
    if err != nil {
        fmt.Println("Failed to dial SSH server:", err)
        return nil, fmt.Errorf("failed to dial: %v", err)
    }

    session, err := connection.NewSession()
    if err != nil {
        return nil, fmt.Errorf("failed to create session: %v", err)
    }

    SshClient = &SSHClient{Session: session}
    return SshClient, nil
}

func ExecuteCmd(command, host string) (string, error) {
	if host == "" {
		out, err := exec.Command("sh", "-c", command).Output()
		if err != nil {
			return "", fmt.Errorf("failed to execute command locally: %v", err)
		}
		return string(out), nil
	}

	hostParts := strings.Split(host, "@")
	user := hostParts[0]
	addr := hostParts[1]

	sshClient, err := getSSHClient(user, addr)
	if err != nil {
		return "", err
	}

	session := sshClient.Session
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	//defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(command)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %v", err)
	}

	return stdoutBuf.String(), nil
}
