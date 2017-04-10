#!/bin/bash
set -u -e -o pipefail

cd ansible/website-deploy
ansible-playbook -i inventory website-deploy.yml
