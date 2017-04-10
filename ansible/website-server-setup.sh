#!/bin/bash
set -u -e -o pipefail

# ansible_ssh_private_key_file=$HOME/.ssh/id_rsa_apptranslator
#ansible dbheroapp -m ping
cd ansible/website-server-setup
ansible-playbook  -i inventory website-server-setup.yml
