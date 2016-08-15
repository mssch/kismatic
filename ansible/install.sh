#!/bin/bash
cd ../tls/
./tls-bootstrap.sh
cd ../ansible
ansible-playbook -i inventory.ini -s kubernetes.yaml
