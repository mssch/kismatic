# Kismatic -- The Enterprise Toolkit for Kubernetes Reliability Engineering (KRE)
[![Build Status](https://snap-ci.com/On8xdVQV0xY5VXICf0Fx0Vq7fVMDUAfU6JFc8Wtt94A/build_image)](https://snap-ci.com/apprenda/kismatic/branch/master)

Kismatic is a set of tools for designing, deploying and operating enterprise-tuned production ready Kubernetes clusters.

We're focused on making it straightforward for organizations who fully manage their own infrastructure to deploy secure, highly available and Kubernetes installations with built-in sane defaults for networking, distributed routing and tracing, cluster health, security and more!

The Kismatic tools include:

1. [`kismatic`](docs/INSTALL.md)
   * A utility for installing and configuring Kubernetes on provisioned infrastructure
2. [`kismatic-inspector`](cmd/kismatic-inspector/README.md)
   * A utility for assuring that software and network configuration of a node are correct prior to installing Kubernetes
3. [`kuberang`](https://github.com/apprenda/kuberang)
   * Tests a cluster to be sure that its networking and scaling work as intended. This tool is used to "smoke test" a newly built cluster.
4. [Kismatic RPM & DEB packages](docs/PACKAGES.md)
   * Packages for installing Kubernetes and its dependencies, focused on specific roles in an HA cluster
   * With these packages installed on a local repo it is possible to use Kismatic to install Kubernetes to nodes that do not have access to the public internet.

## Dependencies
| Dependency | Current version |
| --- | --- |
| Kubernetes | 1.4.5 |
| Docker | 1.11.2 |
| Calico | 1.6 |
| Etcd (for Kubernetes) | 3.0.13 |
| Etcd (for Calico) | 2.37 |

[Download latest install tarball (OSX)](https://kismatic-installer.s3-accelerate.amazonaws.com/latest-darwin/kismatic.tar.gz)

[Download latest install tarball (Linux)](https://kismatic-installer.s3-accelerate.amazonaws.com/latest/kismatic.tar.gz)

# Kismatic Documentation

[Kismatic CLI](https://github.com/apprenda/kismatic/tree/master/kismatic-cli-docs)

[What you can build with Kismatic](docs/INTENT.md)

[Plan & Install a Kubernetes cluster](docs/INSTALL.md)

[Cert Generation](docs/CERT_GENERATION.md)

[Roadmap](ROADMAP.md)

# Dangerously Basic Installation instructions
Use the `kismatic install` command to work through installation of a cluster. The installer expects the underlying infrastructure to be accessible via SSH using Public Key Authentication.

The installation consists of three phases:

1. **Plan**: `kismatic install plan`
   1. The installer will ask basic questions about the intent of your cluster.
   2. The installer will produce a `kismatic-cluster.yaml` which you will edit to capture your intent.
2. **Provision**
   1. You provision your own machines
   2. You tweak your network
   3. Review the installation plan in `kismatic-cluster.yaml` and add information for each node.
3. **Install**: `kismatic install apply`
   1. The installer checks your provisioned infrastructure against your intent.
   2. If the installation plan is valid, Kismatic will build you a cluster.
   3. After installation, Kismatic performs a basic test of scaling and networking on the cluster

###Using your cluster

The installer automatically configures and deploys [Kubernetes Dashboard](http://kubernetes.io/docs/user-guide/ui/) in the cluster, open the link provided at the end of the installation in your browser to use it.    
It will be in the form of `https://%load_balanced_fqdn%:6443/ui`, using `%load_balanced_fqdn%`(from your `kismatic-cluster.yaml` file).      
You will also be prompted for credentials, use `admin` for the **User Name** and `%admin_password%` (from your `kismatic-cluster.yaml` file) for the **Password**.

The installer also generates a [kubeconfig file](http://kubernetes.io/docs/user-guide/kubeconfig-file/) required for [kubectl](http://kubernetes.io/docs/user-guide/kubectl-overview/), just follow the instructions provided at the end of the installation to use it.   

# Development documentation

[How to build](BUILDING.md)

[How to contribute](CONTRIBUTING.md)

[Running integration tests](INTEGRATION_TESTING.md)

[Code of Conduct](code-of-conduct.md)

[Release process](RELEASE.md)
