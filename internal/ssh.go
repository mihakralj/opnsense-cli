package internal

import (
	"fmt"
	"net"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type SSHClient struct {
	Session *ssh.Session
	Client  *ssh.Client
}

var (
	SshClient *SSHClient
	SSHTarget string
	config    *ssh.ClientConfig
)

func getSSHClient(target string) (*SSHClient, error) {
	var user, host, port string

	userhost, port, err := net.SplitHostPort(target)
	if err != nil {
		userhost = target
	}
	if port == "" {
		port = "22"
	}
	split := strings.SplitN(userhost, "@", 2)
	if len(split) == 2 {
		user = split[0]
		host = split[1]
	} else {
		user = "admin"
		host = userhost
	}

	var connection *ssh.Client

	if config == nil {
		config = &ssh.ClientConfig{
			User:            user,
			Auth:            []ssh.AuthMethod{},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		//try to get sshAgent
		sshAgent, err := GetSSHAgent()
		if err == nil {
			config.Auth = append(config.Auth, sshAgent)
			if len(config.Auth) > 0 {
				connection, err = ssh.Dial("tcp", host+":"+port, config)
				if err == nil {
					return &SSHClient{Client: connection}, nil
				}
			}
		}
		fmt.Println("No authorized SSH keys found in local ssh agent, reverting to password-based access.\nTo enable seamless access, use the 'ssh-add' to add the authorized key for user", user)
		fmt.Printf("Enter password for %s@%s: ", user, host)
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			Log(5, "failed to read password: %v", err)
		}
		password := string(bytePassword)
		config.Auth = []ssh.AuthMethod{ssh.Password(password)}
		connection, err = ssh.Dial("tcp", host+":"+port, config)
		if err != nil {
			fmt.Println()
			Log(1, "%v", err)
		} else {
			fmt.Println()
			return &SSHClient{Client: connection}, nil
		}
	}
	connection, err = ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		Log(1, "%v", err)
	}
	SshClient = &SSHClient{Client: connection}

	return SshClient, nil
}
