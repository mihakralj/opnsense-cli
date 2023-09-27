# opnsense-cli features

### Scope and intent

There is a gap between using OPNsense web GUI that offers fail-safe (but limited) configuration capabilities and using FreeBSD command terminal that offers direct access to all functionality of FreeBSD and OPNsense but exposes a great risk of messing things up for anyone that is not well versed in shell commands.

__opnsense-cli__ utility bridges this gap by providing command-line access to local or remote OPNsense firewall. For remote access, it requires `ssh` service to be enabled, as it uses ssh to communicate with the firewall. Every action of __opnsense-cli__ is translated to a shell command that is then executed on OPNsense.

### Features and Benefits
- **Versatility**: Can operate both locally and remotely (via ssh),and is suitable for various deployment scenarios.
- **Transparency and Control**: All opnsense-cli Commands are translated to shell scripts (not API calls), with interactive confirmation for critical changes (bypassable with the --force flag).
- **Cross-Platform Support**: Works on macOS, Windows, Linux, and OpenBSD.
- **Streamlined Operations**: Facilitates repeatable configurations, troubleshooting and complex automations.

### Mechanics

__opnsense-cli__ is focusing on `config.xml` manipulation of OPNsense. All configuration settings are stored in `config.xml` file and OPNSense web GUI actions primarily change data in config XML elements. To protect the integrity of configuration, __opnsense-cli__ is not changing `config.xml` directly - all changes are staged in a separate `staging.xml` file. Configuration elements can be added, removed, modified, discarded and imported - all changes will impact only `staging.xml` until 'commit' command is issued. That's when __opnsense-cli__ will create a backup of `config.xml` and replace it with content from `staging.xml`.

__opnsense-cli__ is also providing commands to manage backup copies in `/conf/backup` directory of OPNsense. It can show all available backups, display details of a specific backup file (including XML diffs between backup file and config.xml), save, restore, load and delete backup files. It can trim number of backup files based on age and desired count of files in the directory.

__opnsense-cli__ also offers (very basic) system management commands. `sysinfo` will display core information about OPNsense instance, `run` command will list and execute all commands that are available through __configctl__ process on OPNsense.

### using ssh identity with __opnsense-cli__

When connecting remotely to OPNsense using ssh, __opnsense-cli__ will try to use private key stored in `ssh-agent` to authenticate. Only when no identities are present or match the public key on OPNsense server, the fallback to *password* will be initiated. As __opnsense-cli__ stores no data locally, the password request will pop-up every time when __opnsense-cli__ initiates the ssh call. Very annoying.

To use ssh identity, both server and client need to be configured with the access key. OPNsense server requires the public key in the format `ssh-rsa AAAAB3NC7we...wIfhtcSt==` and is assigned to a specific user (under System/Access/Users) in the field 'Authorized keys'.

Client needs to support `ssh-agent` and accepts the private key in the format:
```
-----BEGIN RSA PRIVATE KEY-----
[BASE64 ENCODED DATA]
-----END RSA PRIVATE KEY-----
```
The command to add the private key to `ssh-agent`:
```
eval "$(ssh-agent -s)"  # Start the ssh-agent in the background
ssh-add id_rsa          # Add your SSH private key to the ssh-agent
```

