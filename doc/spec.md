# opnsense-cli features and ideas

## Scope and intent

There is a gap between using OPNsense GUI that offers fail-safe but limited configuration capabilities and using FreeBSD terminal that offers direct access to all functionality of OPNsense but at great risk of messing things up.

opnsense-cli utility should bridge this gap:
- allow quick view of basic firweal vitals
- provide command-line access to local or remote OPNsense firewall
- allow staged and controlled changes to conf/config.xml with rollback option
- allow controlled execution of OPNsense configctl commands by calling the same pre-configured commands as GUI

## Commands

**info** - retrieves system information from the firewall

**show** - display hierarchical segments of config.xml

**action** - display (and execute) commands that are available on OPNsense

### Planned

- set
- discard
- bakcup
- commit
- revert/restore/rollback

## Flags

- **--target (-t)** - sets the target OPNsense in the form of user@hostname[:port]
- **--force (-f)** - removes checks and prompts before config.xml or configctl are touched
- **--depth (-d)** -  specifies number of branch levels to show
- **--verbose (-v)** - set verbosity (1=error, 2=warning, 3=info, 4=note, 5=debug)
- **--no-color (n)** - remove ansi colors from printout
- **--xml** - display results in XML format
- **--json** - display results in JSON format
- **--yaml** - display results in YAML format