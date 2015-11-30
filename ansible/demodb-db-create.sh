#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

cd ansible/demodb-db-create
ansible-playbook demodb-db-create.yml
