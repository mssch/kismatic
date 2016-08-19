## Manually run Ansible playbooks
### Prerequisites
* `ansible-playbook`
* `cfssl`

On a Mac run `brew install ansible cfssl`
* manually modify `/etc/hosts` of the control machine public IPv4 and hosts for all the nodes in the cluster
* in `../tls/kubernetes-csr.json` replace with the internal and external IPv4 for all nodes in the cluster
* `./install.sh "@runtime_vars.yaml"`

To just run the playbook
`ansible-playbook -i inventory.ini -s kubernetes.yaml --extra-vars "@runtime_vars.yaml"`

### TODO
Manage certs correctly
