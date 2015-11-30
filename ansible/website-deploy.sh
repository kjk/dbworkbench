#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

cd ansible/website-deploy
ansible-playbook -i inventory website-deploy-systemctl.yml
