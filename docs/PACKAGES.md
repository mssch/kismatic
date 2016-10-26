# Kismatic packages

We provide our own packages to ensure that all services have versions of their dependencies that will work together. Where possible, our packages simply declare dependencies on a vendor's own packages (Docker, for example).

With Kubernetes, packages are still in incubation; we will continue to produce our own for now and repackage as new Kubernetes releases come out.

By default, Kismatic will install its own repos onto machines and use them to download Kismatic packages. This may not be acceptable, for example, if you want to adopt a "golden image" prior to rolling out a many-node cluster, if you need to install a cluster in a lab where most machines are disconnected from the internet, or if you simply want to save bandwidth. If this is your use case, please view the [instructions below](#synclocal).

## Installing via RPM (Redhat, CentOS, Fedora)

1. Add the Kismatic repo to the machine

`sudo curl https://s3.amazonaws.com/kismatic-rpm/kismatic.repo -o /etc/yum.repos.d/kismatic.repo`

2. Install the RPMs for the type of node you want to create

| Product | install command | etcd | master | worker |
| --- | --- | --- | --- | --- |
| Etcd | `sudo yum -y install kismatic-etcd` | Install | | |
| Docker | `sudo yum -y install kismatic-docker-engine` | | Install first | Install first |
| Kubernetes Master | `sudo yum -y install kismatic-kubernetes-master` | | | Install second |
| Kubernetes Worker | `sudo yum -y install kismatic-kubernetes-node` | | Install second\* | |

\*Note: if your intent is to make a single node that fills all three roles (a "minikube" cluster), you should omit kismatic-kubernetes-node. All the packages needed are installed by kismatic-kubernetes-master 

## Installing via DEB (Ubuntu Xenial)

1. Add the Kismatic repo to the machine
   1. Add the Kismatic public key to apt

`wget -qO - https://kismatic-deb.s3-accelerate.amazonaws.com/public.key | sudo apt-key add -` 

   2. Add the Kismatic repo

`sudo add-apt-repository "deb https://kismatic-deb.s3-accelerate.amazonaws.com xenial main"`

2. Refresh the machine's repo cache

`sudo apt-get update`

3. Install the RPMs for the type of node you want to create

| Product | install command | etcd | master | worker |
| --- | --- | --- | --- | --- |
| Etcd | `sudo apt-get -y install kismatic-etcd` | Install | | |
| Docker | `sudo apt-get -y install kismatic-docker-engine` | | Install first | Install first |
| Kubernetes Master | `sudo apt-get -y install kismatic-kubernetes-master` | | | Install second |
| Kubernetes Worker | `sudo apt-get -y install kismatic-kubernetes-node` | | Install second\* | |

\*Note: if your goal is to make a single node that fills all three roles (a "minikube" cluster), you should omit kismatic-kubernetes-node. All the packages needed are installed by kismatic-kubernetes-master 


# <a name="synclocal"></a>Synchronizing a local repo

If you maintain a package repository, you should not perform step 1 in the instructions above. Instead, you should point machines to your own package repository and keep it in sync with Kismatic.

Each Kismatic package will also have many transitive dependencies. To be able to install nodes fully disconnected from the internet, you will need to synchronize your repo with these packages' repos as well.

Dependencies and their versions will change over time and as such they are not listed here. Instead, they can be derived from any machine that is integrated with our repos using the commands linked below.

To help ensure you've correctly synchronized your repo, you can attempt to a test cluster with 1 node of each role. The Kismatic inspector will check your packages to be sure they installed correctly and will alert if any of them are missing.

Changes to dependencies should be called out in the notes that accompany a release.

## yum

[Listing dependencies of a package](http://stackoverflow.com/questions/4627158/how-to-list-all-dependencies-of-a-package-on-linux)

[Syncing with a repo](http://www.tecmint.com/setup-local-repositories-in-ubuntu/)

## apt

[Listing dependencies of a package](http://serverfault.com/questions/199743/how-to-list-rpm-dependencies)

[Syncing with a repo](http://bencane.com/2013/04/15/creating-a-local-yum-repository/)

