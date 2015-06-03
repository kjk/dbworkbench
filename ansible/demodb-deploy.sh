#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

cd ansible
ansible-playbook demodb-deploy.yml
