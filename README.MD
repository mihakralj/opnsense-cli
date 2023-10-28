# OPNsense Command Line Interface (CLI)

*OPNsense CLI* is a command-line utility for FreeBSD, Linux, MacOS and Windows that empowers administrators and power users to manage, configure, and monitor OPNsense firewall systems. The CLI provides an alternative method to browser-based GUI to interact with the firewall system.

[Why this thing exists?](/doc/scope.md)


## Usage

`opnsense [flags] command [parameters]`

### Commands

- **`show [<xpath>]`**: Displays config.xml or the Xpath segment in it
- **`compare [<staging.xml>] [<config.xml>]`**: Compares two config files
- **`set <xpath> [value] [(attribute)]`**: Adds a new branch, value and/or attribute
- **`set <xpath> [value] [(attribute)] -d`**: Deletes branch, value and/or attribute
- **`discard [<xpath>]`**: Discards a value (or all changes) in the 'staging.xml'
- **`commit`**: Moves staging.xml to active 'config.xml'
- **`export [<source.xml>] [<target.xml>]`**: Extracts a patch file
- **`import [patch.xml]`**: Reads provided XML patch and injects it into 'staging.xml'
- **`backup [<backup.xml>]`**: Lists available backup configs or displays a specific backup
- **`restore [<backup.xml>]`**: Restores config.xml from a specific backup.xml. (alias: `load`)
- **`save [<file.xml>]`**: Creates a new /conf/backup/file.xml
- **`delete <backup.xml>`**: Deletes a specific backup.xml.
- **`delete age [days]`**: Deletes all backups older than specified days
- **`delete keep [count]`**: Keeps specified number of backups and deletes the rest
- **`delete trim [count]`**: Deletes number of the oldest backups
- **`sysinfo [<xpath>]`**: Retrieves system information from the firewall
- **`run <service> <command>`**: Executes commands on OPNsense.

### Flags

- **`--target (-t)`**: Sets the target OPNsense in the form of `user@hostname[:port]`.
- **`--no-color (-n)`**: Disable ANSI color output
- **`--force (-f)`**: Removes checks and prompts before `config.xml` or `configctl` are touched.
- **`--verbose (-v)`**: Sets verbosity (1=error, 2=warning, 3=info, 4=note, 5=debug).
- **`--no-color (-n)`**: Removes ANSI colors from the printout.
- **`--xml (-x)`**: Displays results in XML format.
- **`--json (-j)`**: Displays results in JSON format.
- **`--yaml (-y)`**: Displays results in YAML format.

