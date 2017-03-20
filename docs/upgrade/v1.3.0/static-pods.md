# Kubernetes Components as Static Pods

Starting with Kismatic `v1.3.0`, the Kubernetes control plane components are deployed
as static pods instead of systemd services. These include the API server, the controller manager
and the scheduler. The kube-proxy and calico-node components also run as a pod on each node. 
These static pods are deployed into the `kube-system` namespace.

The Kubelet is the only Kubernetes component that continues to be deployed as a systemd service, along with docker and etcd.

The following is a listing of pods that are running in the `kube-system` namespace. As you can see,
the API server, controller manager, scheduler and kube-proxy are all running as pods on the cluster.
Each static pod has the node's hostname appended to it's name (e.g. `kube-proxy-node001` is 
running on `node001`).

```
# kubectl get pods -n kube-system
NAME                                    READY     STATUS    RESTARTS   AGE
calico-node-rxp2b                       2/2       Running   0          2m
default-http-backend-4kpt3              1/1       Running   0          1m
ingress-q65t5                           1/1       Running   0          1m
kube-apiserver-node001                  1/1       Running   0          3m
kube-controller-manager-node001         1/1       Running   0          3m
kube-dns-b0d29                          3/3       Running   0          2m
kube-proxy-node001                      1/1       Running   0          2m
kube-scheduler-node001                  1/1       Running   0          3m
kubernetes-dashboard-1389325151-90dzp   1/1       Running   0          1m
```

The fact that most of the Kubernetes components are now being hosted by Kubernetes itself has
implications when it comes to troubleshooting and getting logs. If the API Server is running,
you may use `kubectl logs` to get logs for a component. 

For example, if you are interested in the controller manager's logs, 
you may use the following command (Replace `node001` with the node's hostname):

```
kubectl logs -n kube-system kube-controller-manager-node001
```

In the case that the API Server is down or unavailable, you must use `docker logs` to get logs from the
components.

For example, if you are interested in diagnosing why the API server is down, you may use the following command:

```
# This gets the container ID of the API server container.
CONTAINER_ID=$(docker ps -a | grep kube-apiserver-amd64 | cut -f 1 -d " ")
docker logs $CONTAINER_ID
```
## Upgrading to Static Pods
When upgrading to `v.1.3.0`, the systemd services that will be replaced by static pods 
are stopped, and the packages removed. Once done, the necessary container images 
are downloaded and static pod manifests are created in `/etc/kubernetes/manifests`.
The kubelet is configured to watch this directory for static pod manifests.
