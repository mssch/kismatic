# High Level Roadmap

- Kubeadm
   * Requiring a centralized tool to add or change the node layout of the cluster is suboptimal and a barrier to infrastructure elasticity. Kubeadm solves some of these issues, but it remains too new for us to include in 1.0.0. We will continue to support the development of Kubeadm upstream and investigate how to best integrate it with KET over the coming releases.
   
- Cluster Upgrades
   * Shortly after the release of Kubernetes 1.5, we expect to have a fully-functional release of Kismatic that will allow users to upgrade Kubernetes custers with zero-downtime to running workloads from minor and patch versions.
   
- Storage
   * We want to make it easy to add durable storage to a cluster via Kismatic. This will also allow us to set aside some amount of durable storage for management tools like Prometheus.
   
-  Monitoring and Logging
   * Built-in packaging and installation of the ELK/Elastic Stack and Fluentd.
   
-  Alerting
   * Built-in packaging and installation of Prometheus.
   
If you could like to add to this roadmap or provide any feedback on the items currently listed, please file an issue!
