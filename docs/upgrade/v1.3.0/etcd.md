# Etcd Upgrade

Prior to Kismatic `v1.3.0`, all cluster deployments were composed of two versions of etcd:

- Kubernetes etcd with version `v3.0.x`
- Calico etcd with version `v2.3.x`

With Kismatic `v1.3.0`, both versions of etcd will be upgraded to etcd `v.3.1.x`
which is not backwards compatible with etcd `v2.3.x` or `v3.0.x`.

## Upgrading Calico etcd `v2.3.x` to `v.3.0.x`

On an etcd cluster that has more than one member, it is not possible to directly 
upgrade from `v2.3.x` to `v.3.1.x`. All nodes in the cluster must be upgraded to `v3.0.x`,
and then upgraded to `v.3.1.x`.

The following is performed on each node, one node at a time:

1. Install `transition-etcd` package containing the etcd `v3.0` binaries: `etcd_v3_0`, `etcdctl_v3_0`
1. Backup etcd data to `/etc/etcd_networking/backup/$timestamp`
1. Restart the `etcd_networking` service
1. Wait until the node has caught up with the rest of the cluster
1. Verify cluster with `etcdctl_v3_0 --endpoint='https://127.0.0.1:6666 cluster-health`

## Upgrading Kubernetes and Calico etcd `v3.0.x` to `v.3.1.x`

The following is performed on each node, one node at a time:

1. Remove old `kismatic-etcd` and `transition-etcd` packages
1. Install new `etcd` package, which contains the following binaries: `etcd_k8s`, `etcd_networking` and `etcdctl`
1. Backup etcd data to `/etc/etcd_k8s/backup/$timestamp`
1. Restart the `etcd_k8s` service
1. Wait until the member has caught up with rest of the cluster
1. Verify Kubernetes etcd cluster with `etcdctl --endpoint='https://127.0.0.1:2379 cluster-health`
1. Backup etcd data to `/etc/etcd_networking/backup/$timestamp`
1. Restart the `etcd_networking` service
1. Wait until the member has caught up with the rest of the cluster
1. Verify Calico etcd cluster with `etcdctl --endpoint='https://127.0.0.1:6666 cluster-health`

Note: Downgrading etcd is not possible once all members are upgraded to `v3.1.x`.