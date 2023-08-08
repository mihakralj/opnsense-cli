//go:build !windows
// +build !windows

package internal

import (
	"errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
)

func GetSSHAgent() (ssh.AuthMethod, error) {
	var agentClient agent.Agent
	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, err
	}
	agentClient = agent.NewClient(sshAgent)

	signers, err := agentClient.Signers()
	if err != nil {
		return nil, err
	}
	if len(signers) == 0 {
		return nil, errors.New("SSH agent has no keys")
	}

	return ssh.PublicKeysCallback(agentClient.Signers), nil
}

func createAgentClient() (agent.Agent, error) {
	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, err
	}
	return agent.NewClient(sshAgent), nil
}
