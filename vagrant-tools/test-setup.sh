#!/bin/bash
set -x
set -e
set -u


tempDir="$(mktemp -d)"
buildDir="$(mktemp -d)"
installDir="opt/web-infrared"
prefixPackage="web-infrared"
venvPackage="venv.tar"

# Create package
oldPwd="$(pwd)"
cp -ra $(dirname $0)/../* "$tempDir"
cd "$tempDir"
cp -ra * ${buildDir}

DEBIAN/venvCreation.sh -b ${buildDir} \
                       -i /${installDir} \
                       -o ${venvPackage}
cp -va "${tempDir}/${venvPackage}" ${buildDir}/

DEBIAN/packetize.sh -b ${buildDir} \
                    -i ${installDir} \
                    -p ${venvPackage} \
                    -o . \
                    -f ${prefixPackage}

# Vagrant up
cd "$oldPwd"
cd $(dirname $0)/test-vm/
vagrant up

# Setup all
ansible-playbook -i ../vagrant.py -i ../vagrant-groups.list -e src_folder="$tempDir" ../../ansible/setup.yml 2>&1 

# Test
vagrant ssh -c "cd /opt/web-infrared/ ; source venv/bin/activate; python tests/test-infra.py"

# Clean
rm -rf "$tempDir"
rm -rf "$buildDir"
