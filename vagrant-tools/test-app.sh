#!/bin/bash
set -x
set -e

cd $(dirname $0)/dev-vm/

# Vagrant up
vagrant up

# Set infrastructure

ansible-playbook -i ../vagrant.py -i ../vagrant-groups.list ../../ansible/lib/setApp.yml
vagrant ssh -c "cd /opt/light-cicd/ ; source venv/bin/activate; python tests/test-app.py"
