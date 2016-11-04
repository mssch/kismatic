# High Level Roadmap

1. Storage
   * We want to make it easy to add durable storage to a cluster via Kismatic. This will also allow us to set aside some amount of durable storage for management tools like Prometheus and allow our customers to run reliable stateful applications on Kubernetes clustered deployed by KET.
2. Ingress
   * As it stands, some networking configurations that are possible with Kismatic are not possible until an Ingress server and controller is present. It would be valuable to be able to manage this configuration within Kismatic.
3. Kubeadm
   * Requiring a centralized tool to add or change the node layout of the cluster is suboptimal and a barrier to infrastructure elasticity. Kubeadm solves a lot of these issues but it remains too new for us to include in this release. We will continue to support the development of Kubeadm and investigate how to best integrate it with KET over the coming releases.
4. Platform Upgrades
   * Shortly after the release of Kubernetes 1.5, we expect to have a release of Kismatic that will allow upgrades from 1.4 -> 1.5.
