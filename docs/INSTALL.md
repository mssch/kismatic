Apprenda Kubernetes Distribution Platform

Matthew M. Miller, Alex Brand, Dimitri Koshkin, Joseph Jacks

**What the Kismatic installer will do:**

* Use key-based SSH to manage machines in a Kubernetes cluster

* Verify that machines and network have been properly provisioned for a Kubernetes installation

* Generate SSL certificates for internal Kubernetes traffic

* Install a software defined network that Kubernetes will use for Pod and Service traffic

* Install Kubernetes components

* Allow for the addition or removal of nodes

* Maintain a record of the original intent of the cluster to help inform upgrades

(installer-workflow.png)

**Plan**

Setting up a proper cluster takes a little forethought. Depending on your intent, you may need to engage multiple teams within your organization to correctly provision the required infrastructure. Planning will help you identify provisioning tasks and report the results of them to the installer.

Planning involves three main areas of concern:

* The machines we'll be installing Kubernetes on

* The network we'll be installing Kubernetes on

* The services we'll be connecting with

You can stand up a small cluster in AWS or virtualized on a personal computer if you just want to get started with Kubernetes.

## Compute resources

<table>
  <tr>
    <td>Etcd Nodes
Suggested: 3</td>
    <td>1      3     5     7</td>
  </tr>
  <tr>
    <td>Master Nodes
Suggested: 2</td>
    <td>1      2 </td>
  </tr>
</table>


Kubernetes is installed on multiple physical or virtual machines and may be provisioned on or on most public infrastructure clouds. These machines become **nodes** of the Kubernetes cluster. 

In our installation of Kubernetes, machines (which are called **nodes**) are specialized to one of three distinct roles within the cluster: **etcd**, **master** or **worker**.

**etcd: **These nodes provide data storage for the master.

**master: **These nodes provide API endpoints and manage the Pods installed on workers.

**worker: **These nodes are where your pods are instantiated.

Nodes within a cluster should have latencies between them of 10ms or lower to prevent instability. If you would like to host workloads at multiple data centers, or in a hybrid cloud scenario, you should expect to set up at least one cluster in each zone.

### Hardware & Operating System

Infrastructure providers supported:

* bare metal

* virtual machines

* AWS

* Packet.net

If using VMs or IaaS, your best experience will be had using dedicated CPU resources for all node types rather than time sharing, as opposed to over-provisioning.

Operating Systems supported:

* RHEL 7

* Centos 7

* Ubuntu 16.04

Minimum hardware requirements:

<table>
  <tr>
    <td>Node Role</td>
    <td>CPU</td>
    <td>Ram</td>
    <td>Disk</td>
  </tr>
  <tr>
    <td>etcd</td>
    <td>1 CPU Core, 2 GHz</td>
    <td>1 GB</td>
    <td>50 GB</td>
  </tr>
  <tr>
    <td>master</td>
    <td>2 CPU Cores, 2 GHz</td>
    <td>2 GB</td>
    <td>50 GB</td>
  </tr>
  <tr>
    <td>worker</td>
    <td>1 CPU Core, 2 GHz</td>
    <td>1 GB</td>
    <td>200 GB</td>
  </tr>
</table>


### Planning for etcd nodes:

Each node receives all the data for a cluster to help protect against data loss in the event that something happens to one of the nodes. A Kubernetes cluster is able to operate as long as more than 50% of its etcd nodes are online. Because of this, we always use an odd number of etcd nodes. Furthermore, adding too many etcd nodes can decrease storage performance. Therefore, count of etcd nodes is primarily an availability concern.

<table>
  <tr>
    <td>Node Count</td>
    <td>Safe for</td>
  </tr>
  <tr>
    <td>1</td>
    <td>Unsafe. Use only for small development clusters.</td>
  </tr>
  <tr>
    <td>3</td>
    <td>Failure of any one node.</td>
  </tr>
  <tr>
    <td>5</td>
    <td>Simultaneous failure of two nodes</td>
  </tr>
  <tr>
    <td>7</td>
    <td>Simultaneous failure of three nodes</td>
  </tr>
</table>


### Planning for master nodes:

Master nodes provide control behaviors for Kubernetes. A Kubernetes cluster is able to operate as long as one of its master nodes is online. Because of this, we suggest at least two master nodes for availability.

<table>
  <tr>
    <td>Node Count</td>
    <td>Safe for</td>
  </tr>
  <tr>
    <td>1</td>
    <td>Unsafe. Use only for small development clusters.</td>
  </tr>
  <tr>
    <td>2</td>
    <td>Failure of any one node.</td>
  </tr>
</table>


Both users of Kubernetes and Kubernetes itself occasionally attempt to communicate with master via a URL. With two or more masters, we strongly suggest introducing a load balanced url (via a virtual IP or DNS CNAME) is required to allow clients and components with Kubernetes to balance between them and continue to operate in the event of a partial failure. Otherwise, failover of a master may require clients to change their settings to continue working with a cluster.

### Planning for worker nodes:

Worker nodes are where your applications will run. You can add more as necessary after the initial setup without interrupting operation of the cluster; your initial worker count should be large enough to hold all the workloads you intend to deploy to it plus enough slack to handle a partial failure.

For example, if you would like to absorb the failure of a VM host, you would need at least one host's worth of unused worker capacity on top of what you would otherwise need.

## Network 

<table>
  <tr>
    <td>Networking Technique</td>
    <td>Routed
Overlay</td>
  </tr>
  <tr>
    <td>How hostnames will be resolved for nodes</td>
    <td>Use DNS
Let Kismatic Manage Hosts Files on nodes</td>
  </tr>
  <tr>
    <td>Network Policy Control</td>
    <td>No network policy
Calico-managed network policy</td>
  </tr>
  <tr>
    <td>Pod Network CIDR Block</td>
    <td></td>
  </tr>
  <tr>
    <td>Services Network CIDR Block</td>
    <td></td>
  </tr>
  <tr>
    <td>Load-Balanced URL for Master Nodes</td>
    <td></td>
  </tr>
</table>


Kubernetes allocates a unique IP address for every Pod created on a cluster. Within a cluster, all Pods are visible to all other Pods and directly addressable by IP, simplifying point to point communications.

Similarly, Kubernetes uses a special network for Services, allowing them to talk to each other via an address that is stable even as the underlying cluster topology changes.

For this to work, Kubernetes makes use of technologies built in to Docker and Linux (including iptables, bridge networking and the Container Networking Interface). We tie these together with a network technology from Tigera networks called Calico.

### Pod and Service CIDR blocks

To provide these behaviors, Kubernetes needs to be able to issue IP addresses from two IP ranges: a **pod network **and** a services network**. This is in addition to the IP addresses nodes will be assigned on their **local network****.**

The pod and service network ranges each need to be assigned a single contiguous CIDR block large enough to handle your workloads and any future scaling. Worker and Master nodes will reserve IP addresses in blocks of 64, so the pod network must be sized so that:

Pod Network IP Block Size >= (Worker Node Count + Master Node Count) * 64

Our default CIDR block for a pod is **172.16.0.0/16**, which would allow for a maximum of roughly 65k pods and 1000 nodes.

Similarly, the service network needs to be large enough to handle all of the Services that might be created on the cluster. Our default is **172.17.0.0/16**, which would allow for 65k services and that ought to be enough for anybody.

Care should be taken that the IP addresses under management by Kubernetes do not collide with IP addresses on the local network, including omitting these ranges from control of  DHCP.

### Pod Networking

There are two techniques we support for pod networking on Kubernetes: **overlay** and **routed**.

In an **overlay** network, communications between pods happen on a virtual network that is only visible to machines that are running an agent. This agent communicates with other agents via the node's local network and establishes IP-over-IP tunnels through which Kubernetes platform traffic is routed.

In this model, no work has to be done to allow pods to communicate with each other (other than ensuring that you are not blocking IP-over-IP traffic). Two or more Kubernetes clusters might even operate on the same pod and services IP ranges, without being able to see each others’ traffic.

However, work does need to be done to expose pods to the local network. This role is usually filled by a Kubernetes Ingress Controller.

Overlay networks work best for development clusters.

In a **routed** network, communications between pods happen on a network that is accessible to all machines on their local network. In this model, each node acts as a router for the ip ranges of the pods that it hosts. The cluster communicates with existing network routers via BGP to establish the responsibility of nodes for routing these addresses. Once routing is in place, a request to a pod or service IP is treated the same as any other request on the network. There is no tunnel or IP wrapping involved. This may also make it easier to inspect traffic via tools like wire shark and tcpdump.

In a routed model, cluster communications often work out of the box. Sometimes routers need to be configured to expect and acknowledge BGP messages from a cluster.

Sometimes, it is valuable to peer nodes in the cluster with a network router that is physically near to them. For this purpose, the cluster announces its BGP messages with an **AS Number** that may be specified when Kubernetes is installed. Our default AS Number is 64511.

Network Policy
By default, Pods can talk to any port and any other Pod, Service or machine on its network. Pod to pod network access is a requirement of Kubernetes, but this degree of openness is not.

When policy is enabled, access to all Pods is restricted and managed in part by Kubernetes and the Calico networking plugin. When adding new Pods, any ports that are identified within the definition will be made accessible to other pods. Access can be further opened or closed using the Calico command line tools installed on every Master node -- for example, you may grant access to a pod, or a namespace of pods, to a developer’s machine.

Network policy gets very advanced and can make prototyping the cluster more difficult; therefore it’s turned off by default.

### DNS & Load Balancing

All nodes in the cluster will need a short name with which they can communicate with each other. DNS is one way to provide this.

It's also valuable to have a load balanced alias for the master servers in your cluster, allowing for transparent failover if a master node goes offline. This can be performed either via DNS load balancing or via a Virtual IP if your network has a load balancer already. Pick a FQDN and short name for this alias to master that defines your cluster's intent -- for example, if this is the only Kubernetes cluster on your network, [kubernetes.yourdomain.com](http://kubernetes.yourdomain.com)** **would be ideal.

If you do not wish to run DNS, you may optionally allow the Kismatic installer to manage hosts files on all of your nodes. Be aware that this option will not scale beyond a few dozen nodes, as adding or removing nodes through the installer will force a hosts file update to all nodes on the cluster.

### Firewall Rules

Kubernetes must be allowed to manage network policy for any IP range it manages.

Network policies for the local network on which nodes reside will need to be set up prior to construction of the cluster, or installation will fail.

<table>
  <tr>
    <td>Purpose</td>
    <td>Target node types</td>
    <td>Source IP range</td>
    <td>Allow Rules</td>
  </tr>
  <tr>
    <td>To allow API server</td>
    <td>worker</td>
    <td>0.0.0.0/0 OR
Only those IP ranges that you want to be able to manage Kubernetes workloads PLUS Kubernetes nodes</td>
    <td>tcp:6443
tcp:8080</td>
  </tr>
  <tr>
    <td>To allow ICMP</td>
    <td>all</td>
    <td>0.0.0.0/0 OR
Only those IP ranges that you want to be able to manage Kubernetes workloads PLUS Kubernetes nodes</td>
    <td>icmp</td>
  </tr>
  <tr>
    <td>To allow all internal traffic between Kubernetes nodes</td>
    <td>all</td>
    <td>All nodes in the Kubernetes cluster</td>
    <td>tcp:0-65535
udp:0-65535</td>
  </tr>
  <tr>
    <td>To allow SSH</td>
    <td>all</td>
    <td>0.0.0.0/0 OR
Only those IP ranges you want to manage the Kubernetes nodes</td>
    <td>tcp:22</td>
  </tr>
</table>


## Certificates and Keys

<table>
  <tr>
    <td>Expiration period for certificates
default 17520h</td>
    <td></td>
  </tr>
  <tr>
    <td>Location City</td>
    <td></td>
  </tr>
  <tr>
    <td>Location State</td>
    <td></td>
  </tr>
  <tr>
    <td>Location Country
default US</td>
    <td></td>
  </tr>
</table>


Kismatic will manage generation and installation of TLS certificates and keys used for intra-platform security. It does this using the open source CloudFlare SSL library and information provided in the Plan file. These certificates and keys are exclusively used to encrypt and authorize traffic between Kubernetes components; they are not presented to end-users.

The default expiry period for certificates is **17520h (**2 years). Certificates must be updated prior to expiry or the platform will cease to operate without warning. Replacing certificates will cause momentary downtime with Kubernetes as of version 1.3; future versions should allow for certificate "rolling" without downtime.

**Provision**

## Get the distribution

You will need to run the installer either from a Linux machine or from a Darwin (OSX) machine. This machine will need to be able to access via SSH all the machines that will become nodes of the Kubernetes cluster.

The installer can run from a machine that will become a node on the cluster, but since the installer's machine holds secrets (such as SSH/SSL keys and certificates), it's best to run from a machine with limited user access and an encrypted disk.

The machine the installer is run from should be available for future modifications to the platform (adding and removing nodes, upgrading Kubernetes).

## Download the installer

### To install from Linux

From an ssh session, type:


curl -L [https](https://is.gd/kismaticlinux)[://is.gd/kismaticlinux](https://is.gd/kismaticlinux) | tar -zx

### To install from Darwin (Mac OSX)

From a terminal window, type

curl -L https://is.gd/kismaticdarwin | tar -zx

## Generate Plan File

From the machine you installed Kismatic to, run the following:

./kismatic install plan

You will be asked a few questions regarding the decisions you made in the Plan section above. The kismatic installer will then generate a **kismatic-cluster.yaml** file.

As machines are being provisioned, you must record their identity and credentials in this file. You should also give your cluster a name and provide an administrative password.

## Compute

### Accessing nodes via the Installer

To install Kismatic, you will need a user with remote sudo access and an ssh public key added to each node. The same username and keypair must be used for all nodes. This account should only be used by the kismatic installer.

We suggest a default user **kismaticuser**, with a corresponding private key **kismaticuser.key **added to the directory you'll be installing kismatic from. This would be the ideal spot to generate the keypair via

ssh-keygen -t rsa -b 4096 -f kismaticuser.key -P ""

The resulting **kismaticuser.pub** will then need to be copied to each node. ssh-copy-id can be convenient for this.

There are four pieces of information we will need to be able to address each node:

<table>
  <tr>
    <td>hostname</td>
    <td>This is a short name that machines can access each other with. If you opt for Kismatic to manage host files for your cluster, this name will be copied to host files.</td>
  </tr>
  <tr>
    <td>ip</td>
    <td>This is the ip that the installer should connect to your node with. If you don't specify a separate internal_ip for the node, the ip will be used for platform traffic as well</td>
  </tr>
  <tr>
    <td>internal_ip (optional)</td>
    <td>In many cases nodes will have more than one physical network card or more than one ip address. Specifying an internal IP address allows you to route traffic over a specific network. It's best for Kubernetes components to communicate with each other via a local network, rather than over the internet.</td>
  </tr>
  <tr>
    <td>labels (optional)</td>
    <td>With worker nodes, labels allow you to identify details of the hardware that you may want to be available to Kubernetes to aid in scheduling decisions. For example, if you have worker nodes with GPUs and worker nodes without, you may want to tag the nodes with GPUs.</td>
  </tr>
</table>


### Configuration

*Note: for the moment, all the software requirements beyond the base OS will be downloaded,  installed and configured by the Kismatic installer. This will take a while over public internet. Soon, we will be offering the option to download and host omnibus RPM & deb packages on an internal repository, allowing for** internet-free installs** and much faster deployment times.*

*The good news is that there is really nothing to do other than set up a base box with a user and public key SSH access.*

<table>
  <tr>
    <td>Requirement</td>
    <td>Required for</td>
    <td>etcd</td>
    <td>master</td>
    <td>worker</td>
  </tr>
  <tr>
    <td>/etc/ssh/sshd_config contains PubkeyAuthentication yes</td>
    <td>Access from kismatic to manage node</td>
    <td>yes</td>
    <td>yes</td>
    <td>yes</td>
  </tr>
  <tr>
    <td>Access to DNS</td>
    <td>Retrieving binaries over the internet during installation</td>
    <td>yes</td>
    <td>yes</td>
    <td>yes</td>
  </tr>
  <tr>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>iptables</td>
    <td>pod networking</td>
    <td></td>
    <td>yes*</td>
    <td>yes*</td>
  </tr>
  <tr>
    <td>iptables-save</td>
    <td>pod networking</td>
    <td></td>
    <td>yes*</td>
    <td>yes*</td>
  </tr>
  <tr>
    <td>iptables-restore</td>
    <td>pod networking</td>
    <td></td>
    <td>yes*</td>
    <td>yes*</td>
  </tr>
  <tr>
    <td>ip</td>
    <td>pod networking</td>
    <td></td>
    <td>yes*</td>
    <td>yes*</td>
  </tr>
  <tr>
    <td>nsenter</td>
    <td>container access</td>
    <td></td>
    <td>yes*</td>
    <td>yes*</td>
  </tr>
  <tr>
    <td>mount</td>
    <td>hosting containers</td>
    <td></td>
    <td>yes*</td>
    <td>yes*</td>
  </tr>
  <tr>
    <td>umount</td>
    <td>hosting containers</td>
    <td></td>
    <td>yes*</td>
    <td>yes*</td>
  </tr>
  <tr>
    <td>glibc</td>
    <td>glibness library</td>
    <td></td>
    <td>yes*</td>
    <td>yes*</td>
  </tr>
  <tr>
    <td>Python 2.7</td>
    <td>kismatic management of nodes</td>
    <td>yes*</td>
    <td>yes*</td>
    <td>yes*</td>
  </tr>
  <tr>
    <td>Docker 1.11 or 1.12</td>
    <td>hosting containers</td>
    <td></td>
    <td>yes*</td>
    <td>yes*</td>
  </tr>
</table>


* *Installer will attempt to install if missing.*

### Inspector

To double check that your nodes are fit for purpose, you can run the kismatic inspector. This tool will be run on each node as part of validating your platform and network fitness prior to installation.

## Networking

Enter your network settings in the plan file, including 

* pod networking technique (**routed** or **overlay**)

* CIDR ranges for pod and services networks

* whether the kismatic installer should manage hosts files for your cluster

Create your DNS CNAME or load balancer alias for your Kubernetes master nodes based on their hostnames.

**Validate**

Having updated your plan, from your installation machine run

./kismatic install validate

This will cause the installer to validate the structure and content of your plan, as well as the readiness of your nodes and network for installation.  Any errors detected will be written to standard out.

This step will result in the copying of the kismatic-inspector to each node via ssh. You should expect it to fail if all your nodes are not yet set up to be accessed via ssh.

**Apply**

Having a valid plan, from your installation machine run 

./kismatic install apply

Kismatic will connect to each of your machines, install necessary software and prove its correctness. Any errors detected will be written to standard out.

Congratulations! You've got a Kubernetes cluster.

