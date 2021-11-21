#!/bin/bash
set -x
set -e

cd $(dirname $0)/dev-vm/

# Vagrant up
vagrant up

# Set infrastructure

# for cleaning the venv
# rm venv/*

ansible-playbook -i ../vagrant.py -i ../vagrant-groups.list ../../ansible/lib/setApp.yml
vagrant ssh -c "cd /opt/photobook/ ; source venv/bin/activate; python app.py"
