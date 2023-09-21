# OPNsense Command Line Interface (CLI)

*OPNsense CLI* is a command-line utility for FreeBSD, Linux, MacOS and Windows that empowers administrators and power users to manage, configure, and monitor OPNsense firewall systems. The CLI provides an alternative method to browser-based GUI to interact with the firewall system.

## Features and Benefits
- *Versatility*: Can operate both locally and remotely (via SSH),and is suitable for various deployment scenarios.
- *Transparency and Control*: All opnsense-cli Commands are bash scripts (and not API calls), with interactive confirmation for changes (bypassable with the --force flag).
- *Cross-Platform Support*: Works with macOS, Windows, Linux, and OpenBSD.
- *Streamlined Operations*: Facilitates repeatable configurations, troubleshooting and complex automations.

## Why Use opnsense CLI?

Administrators  typically only two options how to manage the firewall: either through browser-based GUI, or using FreeBSD commands in bash. The opnsense CLI utility provides the middle option by using command line to quickly perform common tasks:

- A quick view of firewall settings such as configuration, backups, system vitals.
- Controlled changes of `conf/config.xml`, including staging and rollback options.
- Controlled execution of any OPNsense command available through `configctl`.

## Usage

`opnsense [flags] command [parameters]`

### Commands

- **`sysctl [<xpath>]`**: Retrieves system information from the firewall.
- **`show [<xpath>]`**: Displays xpath segment of config.xml.
- **`backups [<backup.xml>]`**: Lists available backup configs or displays a specific backup.
- **`run <service> <command>`**: Executes commands on OPNsense.
- **`set <xpath> value <value>`**: Sets a value of a specific node in staging.xml file.
- **`commit`**: Moves staging.xml to active config.xml.
- **`compare [<staging.xml>] [<config.xml>]`**: Compares two config files.
- **`discard [<xpath>]`**: Discards a value (or all changes) in the staging.xml.
- **`save [<backup.xml>]`**: Creates a new backup.xml.
- **`load [<backup.xml>]`**: Restores config.xml from a specific backup.xml. (alias: `restore`)
- **`delete <backup.xml>`**: Deletes a specific backup.xml.
- **`delete age [days]`**: Deletes all backups older than specified days
- **`delete keep [count]`**: Keeps specified number of backups and deletes the rest
- **`delete trim [count]`**: Deletes number of the oldest backups
- **`set <xpath> [value] [(attribute)]`**: Adds a new branch, value and/or attribute
- **`set <xpath> [value] [(attribute)] -d`**: Deletes branch, value and/or attribute

### Flags

- **`--target (-t)`**: Sets the target OPNsense in the form of `user@hostname[:port]`.
- **`--no-color (-n)`**: Disable ANSI color output
- **`--force (-f)`**: Removes checks and prompts before `config.xml` or `configctl` are touched.
- **`--verbose (-v)`**: Sets verbosity (1=error, 2=warning, 3=info, 4=note, 5=debug).
- **`--no-color (-n)`**: Removes ANSI colors from the printout.
- **`--xml (-x)`**: Displays results in XML format.
- **`--json (-j)`**: Displays results in JSON format.
- **`--yaml (-y)`**: Displays results in YAML format.

