#!/bin/bash
cd ../tls/
./tls-bootstrap.sh
cd ../ansible
if [[ -z "$var" ]]; then
  ansible-playbook -i inventory.ini -s kubernetes.yaml
else
  ansible-playbook -i inventory.ini -s kubernetes.yaml --extra-vars $1
fi
