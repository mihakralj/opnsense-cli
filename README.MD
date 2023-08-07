# OPNsense CLI
opnsense is a command-line utility for managing, configuring, and monitoring OPNsense firewall systems.
It facilitates non-GUI administration, both directly in the shell and remotely via an SSH tunnel.
All interactions with OPNsense utilize the same mechanisms as the Web GUI,
including staged modifications of config.xml and execution of available configd commands.

opnsense can run on Windows, macOs, Linux and OpenBSD - and can target opnsense system remotely.


Usage:

  `opnsense [flags] command`


Available Commands:
  action      A brief description of your command
  info        Retrieves system information
  show        Displays information about configuration stored in config.xml

Flags:
  -t, --target          Target host (-t user@hostname[:port])
  -d, --depth           Specifies number of levels of returned tree (1-5) (default 1)
  -n, --no-color        Turn off ANSI colored output
  -v, --verbose         Set verbosity level (1-5) (default 1)

Examples:

   Local execution in opnsense shell: `opnsense show ntpd`
   ```
   opnsense: {
    ntpd: {
        prefer: "time-a-g.nist.gov 0.north-america.pool.ntp.org clock.sjc.he.net time.windows.com time.apple.com 0.opnsense.pool.ntp.org"
        interface: "wan"
        iburst: "time.windows.com time.apple.com 0.opnsense.pool.ntp.org"
        kod: "1"
        limited: "1"
        nomodify: "1"
        nopeer: "1"
        notrap: "1"
    }
}
   ```

   remote execution (targeting 192.168.1.1), showing details about lan: `opnsense -t admin@192.168.1.1 show interfaces/lan`
   ```
   opnsense: {
    interfaces: {
        lan: {
            if: "igb1"
            descr: "LAN1"
            enable: "1"
            lock: "1"
            spoofmac: ""
            blockbogons: "1"
            media: "autoselect"
            mediaopt: ""
            ipaddr: "192.168.1.12"
            subnet: "24"
            ipaddrv6: "2001:470:a920:1::1"
            subnetv6: "64"
        }
    }
}
   ```

   show two levels of depth of system information: `opnsense -t admin@192.168.1.1 info/hardware -d2`
   ```
   system: {
    hardware: {
        model: "Intel(R) Core(TM) i7-7500U CPU @ 2.70GHz"
        machine_arch: "amd64"
        machine: "amd64"
        clockrate: "2900"
        ncpu: "2"
        smp.cpus: "2"
        realmem: "32.00GB"
        physmem: "31.87GB"
        usermem: "24.45GB"
        disks: "ada0"
    }
}
   ```

