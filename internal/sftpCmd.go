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
package internal

import (
	"os"

	"github.com/pkg/sftp"
)

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
