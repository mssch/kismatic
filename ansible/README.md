## Manually run Ansible playbooks
### Prerequisites
* `ansible`
* `cfssl`

On a Mac run `brew install ansible cfssl`

On packet create 3 etcd servers, 2 master servers and some workers (only CentOS7 has been tested)
* manually modify `inventory.ini` with the hostnames, `ansible_host`=public IPv4 and `internal_ipv4` of the machines you created
* in `../tls/kubernetes-csr.json` replace with the internal and hostnames, external and internal IPv4 for all the nodes in the cluster
* run `./install.sh`

To just run the playbook
`ansible-playbook -i inventory.ini -s kubernetes.yaml --extra-vars "@runtime_vars.yaml"`

### TODO
Manage certs correctly
