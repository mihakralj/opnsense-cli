/*
Copyright © 2023 Miha miha.kralj@outlook.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
//go:build windows
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

func createAgentClient() (agent.Agent, error) {
	conn, err := winio.DialPipe(`\\.\pipe\openssh-ssh-agent`, nil)
	if err != nil {
		return nil, err
	}
	return agent.NewClient(conn), nil
}
