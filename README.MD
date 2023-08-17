# OPNsense CLI

The **opnsense** CLI is a  command-line utility designed to manage, configure, and monitor OPNsense firewall systems. It provides administrators and power users an alternative way to interact with the firewall system without relying on the graphical user interface (GUI). The CLI facilitates both local and remote management via an SSH tunnel, making it versatile for a range of deployment scenarios. Whether it's executing direct configurations, orchestrating complex automation tasks, or troubleshooting on-the-fly, opnsense CLI streamlines these operations through FreeBSD shell commands, employing the same commands used in the Web GUI. Designed with transparency and verifiability in mind, the CLI's shell commands are traceable and each change to the firewall is seeking interactive confirmation (unless supressed with --force flag). Opnsense CLI supports various operating systems including Windows, macOS, Linux, and OpenBSD, allowing creation of administration scripts on either locally on OPNsense system or remotely from admin workstation.

The opnsense CLI can be installed on Windows, macOS, Linux, and OpenBSD, and can target an OPNsense system either locally or remotely.

## Why Use opnsense CLI?

The OPNsense GUI offers fail-safe but limited configuration capabilities, while the FreeBSD terminal provides direct access to all functionalities of OPNsense but with a greater risk of errors. The `opnsense-cli` utility bridges this gap by:

- Allowing a quick view of basic firewall vitals.
- Providing command-line access to a local or remote OPNsense firewall.
- Enabling staged and controlled changes to `conf/config.xml` with rollback options.
- Allowing controlled execution of OPNsense `configctl` commands, calling the same pre-configured commands as the GUI.

## Usage

`opnsense [flags] command [parameters]`

### (Planned) Commands

- **`show system <xpath>`**: Retrieves system information from the firewall.
- **`show config <xpath>`**: Displays hierarchical segments of `config.xml`.
- **`show backup <config>`**: Lists available backup configs or a specific backup.
- **`run <service> <command>`**: Executes commands on OPNsense.
- **`set <xpath> value <value>`**: Sets a value of a specific node in the staging config.
- **`commit`**: Moves staging config to active config.
- **`compare [<config>] [<config>]`**: Compares staging config and active config or any two configs.
- **`discard [<xpath>]`**: Discards a value in the staging config or all changes in the staging config.
- **`save [<config>]`**: Creates a new backup config.
- **`load <config>`**: Moves a specific backup config to active config.
- **`delete <config>`**: Deletes a specific config file.

### Flags

- **`--target (-t)`**: Sets the target OPNsense in the form of `user@hostname[:port]`.
- **`--force (-f)`**: Removes checks and prompts before `config.xml` or `configctl` are touched.
- **`--depth (-d)`**: Specifies the number of branch levels to show.
- **`--verbose (-v)`**: Sets verbosity (1=error, 2=warning, 3=info, 4=note, 5=debug).
- **`--no-color (-n)`**: Removes ANSI colors from the printout.
- **`--xml`**: Displays results in XML format.
- **`--json`**: Displays results in JSON format.
- **`--yaml`**: Displays results in YAML format.

## Installation

BSD12/13 install of opnsense.txz package (into _usr/local/bin_):
`sudo pkg add https://github.com/mihakralj/opnsense-cli/releases/download/beta/opnsense.txz`

Signed macos pkg (into _/usr/local/bin_):
`wget https://github.com/mihakralj/opnsense-cli/releases/download/beta/opnsense-macos.pkg`
`sudo installer -pkg opnsense-macos.pkg -target /`

Direct download of signed macos opnsense binary:
`wget https://github.com/mihakralj/opnsense-cli/releases/download/beta/opnsense-macos -O opnsense`
`sudo chmod +x opnsense-macos`

Ubuntu/Debian opnsense.deb install (into _usr/local/bin_):
`wget https://github.com/mihakralj/opnsense-cli/releases/download/beta/opnsense.deb`
`dpkg -i opnsense.deb`

Windows opnsense.exe:
` Invoke-WebRequest -Uri https://github.com/mihakralj/opnsense-cli/releases/download/beta/opnsense.exe -OutFile opnsense.exe`