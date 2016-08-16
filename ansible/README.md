## Manually run Ansible playbooks
### Prerequisites
* `ansible-playbook`
* `cfssl`

On a Mac run `brew install ansible cfssl`
* manually modify `inventory.ini` with the IPs/hosts
* manually modify `../tls/kubernetes-csr.json` with the IPs/hosts
* manually modify `runtime_vars.yaml` with the IPs/hosts
* `./install.sh "@runtime_vars.yaml"`

To just run the playbook
`ansible-playbook -i inventory.ini -s kubernetes.yaml --extra-vars "@runtime_vars.yaml"`

### TODO
Manage certs correctly
