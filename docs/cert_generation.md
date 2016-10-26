# Certificate Generation

## What are certs used for?
* Internal communications between platform components
  * Kubectl -> API server
  * API Server -> Etcd
  * Etcd <-> Etcd
  * Calico -> Etcd
* Client authentication against the API server
 
## What certs get generated?
* Self-signed certificate to be used as CA
* One certificate for each node on the cluster
* One certificate for an admin user
* One certificate for the private Docker registry, to use that Docker registry each Docker engine that push/pulls from must also have the CA
 
## How are certs generated?
* Using cfssl (https://github.com/cloudflare/cfssl
  * Algorithm: RSA
  * Key Size: 2048
* Common Name:
  * CA certificate => cluster name
  * Node certificate => node’s machine name (could be hostname or FQDN)
  * User certificate => user name (this is a K8s requirement)
* Subject:
  * Org => Apprenda
  * OU => Kismatic
  * Country => US
  * State => NY
  * Locality => Troy
* Expiration: configurable, defaults to 17600h (2 years)
