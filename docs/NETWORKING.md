# Kubernetes Networking
Kismatic installs [Calico](https://www.projectcalico.org/) as the networking solution for the cluster. 
Detailed documentation for Calico may be found [here](http://docs.projectcalico.org/).

The main Calico components are running on each master and worker node of the cluster. The components themselves are
running inside a container called `calico-node`. 

## Calicoctl
Calicoctl is the command-line utility for managing the Calico network. This utility is installed on all master nodes of your cluster.
In order to use calicoctl, a number of environment variables must be defined. 
These environment variables can be found in `/etc/network-environment`. 

A quick way to start using calicoctl:
```
for e in `cat /etc/network-environment`; do export $e; done
calicoctl status
```

You may find calicoctl's reference guide [here](http://docs.projectcalico.org/v1.6/reference/calicoctl/)

## Network Policy (Advanced Feature)
Calico supports defining Network Policy in your cluster. This is an advanced feature that is disabled by default in
the installation plan file. 

More detailed documentation on policy can be found [here](http://docs.projectcalico.org/v1.6/getting-started/kubernetes/tutorials/simple-policy)

## Useful links
* Kubernetes + Calico overview: http://docs.projectcalico.org/v1.6/getting-started/kubernetes/
* Troubleshooting: http://docs.projectcalico.org/v1.6/getting-started/kubernetes/troubleshooting
* Logging: http://docs.projectcalico.org/v1.6/usage/troubleshooting/logging
