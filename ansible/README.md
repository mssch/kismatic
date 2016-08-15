## Manually run Ansible playbooks
### Prerequisites
* `ansible-playbook`
* `cfssl`

On a mac run `brew install ansible cfssl`
* manually modify `inventory.ini` with the IPs/hosts
* manually modify `../tls/kubernetes-csr.json` with the IPs/hosts
* manually modify `group_vars\etcd_k8s.yaml` and `group_vars\etcd_networking.yaml` with the IPs/hosts
* `./install.sh`

### TODO
Manage certs correctly
