#!/bin/bash
set -u -e -o pipefail

# ansible_ssh_private_key_file=$HOME/.ssh/id_rsa_apptranslator
#ansible dbheroapp -m ping
cd ansible/demodb-server-setup
ansible-playbook demodb-server-setup.yml
