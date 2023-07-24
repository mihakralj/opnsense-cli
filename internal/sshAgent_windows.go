// +build windows
package internal

import (
    "errors"
    "github.com/Microsoft/go-winio"
    "golang.org/x/crypto/ssh"
    "golang.org/x/crypto/ssh/agent"
)

func GetSSHAgent() (ssh.AuthMethod, error) {
    var agentClient agent.Agent
    conn, err := winio.DialPipe(`\\.\pipe\openssh-ssh-agent`, nil)
    if err != nil {
        return nil, err
    }
    agentClient = agent.NewClient(conn)

    signers, err := agentClient.Signers()
    if err != nil {
        return nil, err
    }
    if len(signers) == 0 {
        return nil, errors.New("SSH agent has no keys")
    }

    return ssh.PublicKeysCallback(agentClient.Signers), nil
}
