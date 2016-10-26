# Kismatic -- Keep IT Simple: Make a Tidy infrastructure Cluster
[![Build Status](https://snap-ci.com/On8xdVQV0xY5VXICf0Fx0Vq7fVMDUAfU6JFc8Wtt94A/build_image)](https://snap-ci.com/apprenda/kismatic-platform/branch/master)

Kismatic is a set of utilities for managing Kubernetes installations.

We're focused on making it straightforward for those who manage their own infrastructure to install secure, highly available Kubernetes clusters.

The Kismatic tools include:

1. `kismatic`
   * A utility for installing and configuring Kubernetes on provisioned infrastructure
2. [`kismatic-inspector`](cmd/kismatic-inspector/README.md)
   * A utility for assuring that software and network configuration of a node are correct prior to installing Kubernetes
3. [`kuberang`](https://github.com/apprenda/kuberang)
   * Tests a cluster to be sure that its networking and scaling work as intended. This tool is used to "smoke test" a newly built cluster.
4. [Kismatic RPM & DEB packages](docs/PACKAGES.md)
   * Packages for installing Kubernetes and its dependendencies, focused on specific roles in an HA cluster
   * With these packages installed on a local repo it is possible to install use Kismatic to Kubernetes to nodes that do not have access to the public internet.

| Dependency | Current version |
| --- | --- |
| Kubernetes | 1.4.3 |
| Docker | 1.11.2 |
| Calicoctl | 0.22.0 |
| Etcd (for Kubernetes) | 3.0.10 |
| Etcd (for Calico) | 2.37 |

[Download latest install tarball (OSX)](https://kismatic-installer.s3-accelerate.amazonaws.com/kismatic-installer/latest-darwin/kismatic.tar.gz)

[Download latest install tarball (Linux)](https://kismatic-installer.s3-accelerate.amazonaws.com/kismatic-installer/latest/kismatic.tar.gz)

# Kismatic Documentation

[Plan & Install a Kubernetes cluster](docs/INSTALL.md)

[Cert Generation](docs/cert_generation.md)

[Roadmap](ROADMAP.md)

# Dangerously Basic Installation instructions
Use the `kismatic install` command to work through installation of a platform. The installer expects the underlying infrastructure to be accessible via SSH using Public Key Authentication.

The installation consists of three phases:

1. **Plan**: `kismatic install plan` 
   1. The installer will ask basic questions about the intent of your cluster.
   2. The installer will produce a `kismatic-cluster.yaml` which you will edit to capture your intent.
2. **Provision** 
   1. Provision machines
      1. Allocate hardware (bare metal machines, VMs, EC2 instances)
      2. Open access to the installer
      3. Optionally install software packages
   2. Provision networking
      1. Find convenient CIDR blocks for Kubernetes components
      2. Decide on network type (overlay or routed)
      3. Open up ports where necessary.
   3. Review the installation plan in `kismatic-cluster.yaml` and add information for each node.
3. **Install**: `kismatic install apply` 
   1. Every install phase begins by validating the plan and testing the infrastructure referenced within it.
   2. If the installation plan is valid, the installer will build you a cluster.
      1. Validate that nodes are provisioned correctly
      2. Install (or validate) software packages including Docker and Kubernetes
      3. Generate TLS certificates and keys for intra cluster communications
      4. Configure the cluster 
      5. After configuration, run a smoke test to ensure that scaling and pod networking are working as prescribed.

# Development documentation

[How to build](BUILDING.md)

[How to contribute](CONTRIBUTING.md)

[Running integration tests](INTEGRATION_TESTING.md)

[Code of Conduct](code-of-conduct.md)
