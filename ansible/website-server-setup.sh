#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

# ansible_ssh_private_key_file=$HOME/.ssh/id_rsa_apptranslator
#ansible dbworkbench -m ping
cd ansible/website-server-setup
ansible-playbook  -i inventory website-server-setup.yml
