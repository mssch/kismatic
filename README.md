# Kismatic
[![Build Status](https://snap-ci.com/On8xdVQV0xY5VXICf0Fx0Vq7fVMDUAfU6JFc8Wtt94A/build_image)](https://snap-ci.com/apprenda/kismatic-platform/branch/master)

Kismatic is a set of utilities for managing Kubernetes installations.

We're focused on making it straightforward for those who manage their own infrastructure to install secure, highly available Kubernetes clusters.

The Kismatic tools include:

1. `kismatic`
   * A utility for installing and configuring Kubernetes on provisioned infrastructure
2. `kismatic-inspector`
   * A utility for assuring that software and network configuration of a node are correct prior to installing Kubernetes
3. [`kuberang`](https://github.com/apprenda/kuberang)
   * Tests a cluster to be sure that its networking and scaling work as intended. This tool is used to "smoke test" a newly built cluster.
4. Kismatic RPM & DEB packages
   * Packages for installing Kubernetes and its dependendencies, focused on specific roles in an HA cluster
   * With these packages installed on a local repo it is possible to install use Kismatic to Kubernetes to nodes that do not have access to the public internet.

| Dependency | Currently installed version |
| --- | --- |
| Kubernetes | 1.4.3 |
| Docker | 1.11.2 |
| Calicoctl | 0.22.0 |
| Etcd (for Kubernetes) | 3.0.10 |
| Etcd (for Calico) | 2.37 |

# Kismatic Documentation

[Plan & Install a Kubernetes cluster](docs/INSTALL.md)

[Cert Generation](docs/cert_generation.md)

[Roadmap](ROADMAP.md)

# Installation Basics
Use the `kismatic install` command to work through installation of a platform. The installer expects the underlying infrastructure to be accessible via SSH using Public Key Authentication.

The installation consists of three phases:

1. **Plan**: `kismatic install plan` 
   1. The installer will ask basic questions about the intent of your cluster.
   2. The installer will produce a `kismatic-cluster.yaml` to capture your intent.
   3. Review the installation plan, and fill out the information for each node.
   4. Once done, run the `kismatic install validate` command.
2. **Validate**: `kismatic install validate`
   1. The installer will validate the installation plan, and make sure that all fields are valid.
   2. No changes will be made to the cluster at this time.
3. **Install**: `kismatic install apply` 
   1. Every install phase begins by validating the plan and testing the infrastructure referenced within it.
   2. If the installation plan is valid, the installer will execute installation.

# Development documentation

[How to build](BUILDING.md)

[How to contribute](CONTRIBUTING.md)

[Running integration tests](INTEGRATION_TESTING.md)

[Code of Conduct](code-of-conduct.md)
