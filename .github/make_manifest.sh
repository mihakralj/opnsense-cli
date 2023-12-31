#!/bin/sh

mkdir -p ./dist/pkg/usr/local/bin
export GOARCH=amd64
export GOOS=freebsd
export GOCACHE=~/Github/opnsense-cli/dist/pkg/cache
go mod vendor
go build -trimpath -ldflags "-w -s" -mod=vendor -o ./dist/pkg/usr/local/bin/opnsense opnsense.go
chmod +x ./dist/pkg/usr/local/bin/opnsense
SHA=$(shasum -a 256 ./dist/pkg/usr/local/bin/opnsense | awk '{ print $1 }')

VERSION=$1
FLATSIZE=$(du -b -s ./dist/pkg/usr/local/bin | cut -f1)
MANIFEST="./dist/pkg/+MANIFEST"
echo -e "{\n\"name\": \"opnsense-cli\"," > $MANIFEST
echo -e "\"version\": \"${VERSION}\"," >> $MANIFEST
echo -e "\"origin\": \"net-mgmt/opnsense-cli\"," >> $MANIFEST
echo -e "\"comment\": \"CLI to manage and monitor OPNsense firewall configuration, check status, change settings, and execute commands.\"," >> $MANIFEST
echo -e "\"desc\": \"opnsense is a command-line utility for managing, configuring, and monitoring OPNsense firewall systems. It facilitates non-GUI administration, both directly in the shell and remotely via an SSH tunnel.\"," >> $MANIFEST
echo -e "\"maintainer\": \"miha.kralj@outlook.com\"," >> $MANIFEST
echo -e "\"www\": \"https://github.com/mihakralj/opnsense-cli\"," >> $MANIFEST 
echo -e "\"abi\": \"FreeBSD:*:amd64\"," >> $MANIFEST
echo -e "\"prefix\": \"/usr/local\"," >> $MANIFEST
echo -e "\"flatsize\": ${FLATSIZE}," >> $MANIFEST
echo -e "\"licenselogic\": \"single\"," >> $MANIFEST
echo -e "\"licenses\": [\"APACHE20\"]," >> $MANIFEST
echo -e "\"files\": {" >> $MANIFEST
echo -e "\"/usr/local/bin/opnsense\": \"SHA256:$SHA\"" >> $MANIFEST
echo -e "}" >> $MANIFEST
echo -e "}" >> $MANIFEST

MANIFEST="./dist/pkg/+COMPACT_MANIFEST"
echo -e "{\n\"name\": \"opnsense-cli\"," > $MANIFEST
echo -e "\"version\": \"${VERSION}\"," >> $MANIFEST
echo -e "\"origin\": \"net-mgmt/opnsense-cli\"," >> $MANIFEST
echo -e "\"comment\": \"CLI to manage and monitor OPNsense firewall configuration, check status, change settings, and execute commands.\"," >> $MANIFEST
echo -e "\"www\": \"https://github.com/mihakralj/opnsense-cli\"," >> $MANIFEST 
echo -e "\"abi\": \"FreeBSD:*:amd64\"," >> $MANIFEST
echo -e "}" >> $MANIFEST

cd ./dist/pkg
tar -cJf ../opnsense-cli-${VERSION}.txz -s'|^\./|/|' -P +MANIFEST +COMPACT_MANIFEST ./usr/local/bin/opnsense
