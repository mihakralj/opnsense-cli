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
}

var (
	SshClient *SSHClient
	SSHTarget string
    config *ssh.ClientConfig
)

func getSSHClient(target string) (*SSHClient, error) {
	var user, host, port string

	userhost, port, err := net.SplitHostPort(target)
	if err != nil {
		userhost = target
	}
	if port=="" {
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


	if config == nil {
        config = &ssh.ClientConfig{
            User:            user,
            HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        }

        sshAgent, err := GetSSHAgent()
        if err == nil {
            config.Auth = []ssh.AuthMethod{sshAgent}
        }

        if len(config.Auth) == 0 {
            fmt.Println("No suitable SSH identities found in ssh-agent.\nFor enhanced security add SSH key to the ssh agent")
            fmt.Printf("Enter password for %s@%s: \n", user, host)
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
            fmt.Println()
            if err != nil {
                return nil, fmt.Errorf("failed to read password: %v", err)
            }
            password := string(bytePassword)
            config.Auth = []ssh.AuthMethod{ssh.Password(password)}
        }
    }


	connection, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		Log(1, "%v",err)
	}

	session, err := connection.NewSession()
	if err != nil {
		Log(1, "%v",err)
	}

	SshClient = &SSHClient{Session: session}
	return SshClient, nil
}
