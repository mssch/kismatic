## Manually run Ansible playbooks
### Prerequisites
* `ansible`
* `cfssl`
* `jq`

On a Mac run `brew install ansible cfssl jq`

On packet create 3 etcd servers, 2 master servers and some workers (only CentOS7 has been tested)  
Get list for inventory  
`curl -H 'X-Auth-Token: CuUyhwHviRwHxJk3qruZH8ui9V3EN2Rr' https://api.packet.net/projects/001400bb-f3e4-46c1-bbf4-7bbf9ab849ad/devices\?per_page\=100 | jq -r '.devices|=sort_by(.hostname)|.devices[]|.hostname+" ansible_host="+.ip_addresses[0].address+" internal_ipv4="+.ip_addresses[2].address'`

Get list for servers  
`curl -H 'X-Auth-Token: CuUyhwHviRwHxJk3qruZH8ui9V3EN2Rr' https://api.packet.net/projects/001400bb-f3e4-46c1-bbf4-7bbf9ab849ad/devices\?per_page\=100 | jq -r '.devices|=sort_by(.hostname)|.devices[]|.hostname+","+.ip_addresses[0].address+","+.ip_addresses[2].address'`

* manually modify `inventory.ini` with the hostnames, `ansible_host`=public IPv4 and `internal_ipv4` of the machines you created
* manually modify tls/servers (create the file if doesn't exist) with the output of the 1st command above for your nodes
* run `./install.sh`

To just run the playbook
`ansible-playbook -i inventory.ini -s kubernetes.yaml --extra-vars "@runtime_vars.yaml"`

### TODO
Manage certs correctly
