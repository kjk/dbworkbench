#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

# ansible_ssh_private_key_file=$HOME/.ssh/id_rsa_apptranslator
#ansible dbworkbench -m ping
cd ansible/demodb-server-setup
ansible-playbook demodb-server-setup.yml
