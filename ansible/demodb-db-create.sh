#!/bin/bash
set -u -e -o pipefail

cd ansible/demodb-db-create
ansible-playbook demodb-db-create.yml
