# Cluster Intent

| Why do you want to install Kubernetes? | Example Cluster |
| --- | --- |
| I want to run a small prototype cluster on my laptop| [Minikube-style](#minikube) |
| I want to run a small development cluster in AWS | [Developer](#junior) |
| I want to run a reliable, permanent labratory cluster to host mostly services for my team. | [Skunkworks](#skunkworks) |
| I want to run a reliable, permanent production cluster to host services for my team and others. | [Production](#production) | 

# <a name="minikube"></a>Minikube style

*Just heard about Kubernetes. Love Docker and want to see what the fuss is about.*

![](minikube.jpg)

## What you need:

* A developer's laptop (Running OSX or Linux)
* A desktop virtualizer (such as Oracle VirtualBox)
* A VM lifecycle tool (such as Vagrant) can dramatically speed up set up time for virtual machines

## How you install:

1. Download the [latest Kismatic](PROVISION.md#get) for your OS and unpack it somewhere convenient.
2. Create a [small Linux VM](PLAN.md#compute)
   * You will need to know its IP address and short name.
   * You will need to create a user with passwordless auth, ssh capability and a [public key](PROVISION.md#access).
3. Run `./kismatic install plan`; you will be creating 1 of each node
4. Open the `kismatic-cluster.yaml` file in a text editor
   * Set `allow_package_installation: true`
   * Update your user and public key if necessary.
   * Enter the machine's short hostname as the `load_balanced_fqdn` and `load_balanced_short_name`
   * Add your one machine's IP and short hostname as a node in the section for each type of node -- etcd, worker and master.
5. Run `./kismatic install apply`
6. Congratulations! You have your first cluster.

# <a name="junior"></a>1/1/1+ Developer's Cluster

*I'm loving the Kubernetes and want to share it with the rest of my team, maybe give them each their own little cluster to play with.*

![](dev.jpg)

## What you need:

* A developer's laptop (Running OSX or Linux)
* An AWS account

## How you install:

1. Download the [latest Kismatic](PROVISION.md#get) for your OS and unpack it to `~/kismatic`
2. Make a new [Keypair](https://aws.console.aws.amazon.com/ec2/v2/home?#KeyPairs:sort=keyName) for your cluster. Download its private key to `~/kismatic/kismaticuser.key` 
3. `chmod 0600 kistmaticuser.key` to allow its use for SSH 
3. Make 2 small (micro) instances using Ubuntu 16.04 or RHEL. One of these will be your etcd and the other will be your master.
   * Assign public IPs -- we can always remove them later, but this will ease installation
   * Add the new instances to a security group that will allow access to any port from your machine
   * Be sure to select your new keypair on creation
4. Make 1 or more medium or larger instances using Ubuntu 16.04 or RHEL. This is your worker pool. 
   * Assign public IPs -- we can always remove them later, but this will ease installation
   * Add the new instances to a security group that will allow access to any port from your machine
   * Be sure to select your new keypair on creation
5. Run `./kismatic install plan`; you will be creating 1 of each node (more if using more workers)
6. Open the `kismatic-cluster.yaml` file in a text editor
   * Set `allow_package_installation: true`
   * Update your user (`ubuntu` for Ubuntu, `ec2-user` for RHEL)
   * Add your nodes to the sections for etcd, master and worker(s), being sure to set up the medium or larger instance(s) as worker(s) 
     * `host:` Use the Private DNS Name (these will be in the form `ip-XX-XX-XX-XX`)
     * `ip:` Use the Public IP Address
     * `internal_ip:` Use the first Private IP
   * Enter the master's short hostname as the `load_balanced_fqdn` and `load_balanced_short_name`
7. Run `./kismatic install apply`
8. Congratulations! You have your first cluster.

# <a name="skunkworks"></a>3+/2/2+ Skunkworks Cluster

![](skunkworks.jpg)

*I would like to build and grow a big cluster for my team to share and work with. I don't need to share much with other people in the company, want to avoid introducing complexity to our network or I prefer the security of not having all of my Kubernetes pods addressable automatically by anybody with access to the network. This is a production environment, but we can probably survive a major disaster so long as it is low risk.*

## What you need:

* An AWS account, bare metal machines or a bunch of VMs
* An engineer who understands how to provision hardware

## How you install:

### Hardware Provisioner
1. Make an install machine. This is a small (1 CPU, 1 GB ram, <8 GB hard drive) Linux VM with an encrypted disk and very limited access -- just potential cluster administrators.

### Cluster Operator
1. Download the [latest Kismatic](PROVISION.md#get) for Linux and unpack it somewhere convenient on the install machine.

# <a name="production"></a>5+/2+/2+ Production Cluster

![](production.png)

*I want to build and grow a big cluster for my team to share their apps with the rest of the company. Security is a secondary concern to access. Also, this is production: we can't take any chances.*

## What you need:

* Bare metal machines or a bunch of VMs
* An engineer who understands your network
* An engineer who understands how to provision hardware

## How you install:

1. Take a meeting between the parties responsible for the network infrastructure, machine infrastructure and the cluster you are building
2. Discuss the size of the cluster -- both initial and considering six months to a year of growth
3.  

### Hardware Provisioner
1. Make an install machine. This is a small (1 CPU, 1 GB ram, <8 GB hard drive) Linux VM with an encrypted disk and very limited access -- just potential cluster administrators.


### Cluster Operator
1. Download the [latest Kismatic](PROVISION.md#get) for Linux and unpack it somewhere convenient on the install machine.
