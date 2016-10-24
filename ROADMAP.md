# Kismatic is focused on simplifying the management of privately hosted Kubernetes by making best practices simpler to implement.

1. Storage
   * Stateful workloads are a fact of life and running them on Kubernetes, with its PersistentVolume awareness, makes sense. We want to make it easy to add durable storage to a cluster via Kismatic. This will also allow us to set aside some amount of durable storage for management tools like Prometheus.  
2. Ingress
   * As it stands, some networking configurations that are possible with Kismatic aren't terribly useful until you also install an Ingress server and controller. It would be valuable to be able to manage this configuration within Kismatic.
3. Kubeadm
   * Requiring a centralized tool to add or change the node layout of the cluster is suboptimal and a barrier to infrastructure elasticity. Kubeadm solves a lot of these issues but it remains too new for us to include in this release.
4. Platform Upgrades
   * Shortly after the release of Kubernetes 1.5, we expect to have a release of Kismatic that will allow upgrades from 1.4 -> 1.5