#!/bin/bash
cd tls/
curl -H 'X-Auth-Token: CuUyhwHviRwHxJk3qruZH8ui9V3EN2Rr' https://api.packet.net/projects/001400bb-f3e4-46c1-bbf4-7bbf9ab849ad/devices\?per_page\=100 | jq -r '.devices|=sort_by(.hostname)|.devices[]|.hostname+","+.ip_addresses[0].address+","+.ip_addresses[2].address' > servers
./tls-bootstrap.sh
cd ../
if [[ -z "$var" ]]; then
  ansible-playbook -i inventory.ini -s kubernetes.yaml
else
  ansible-playbook -i inventory.ini -s kubernetes.yaml --extra-vars $1
fi
