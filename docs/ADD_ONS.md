# Add-ons

Once the core Kubernetes cluster components are up and running, KET installs
additional components that are required to produce a fully operational cluster.
These extra components are collectively referred to as add-ons, and they include
the networking plugin, cluster DNS, heapster monitoring, and others.

Add-on installation and configuration can be controlled via the plan file under 
the `add-ons` section.

## List of add-ons

### Heapster
Heapster is a monitoring solution that enables container monitoring throughout
the cluster. When heapster is running on the cluster, `kubectl` and the Kubernetes 
dashboard surface resource utilization metrics for all pods.

**Important:** If you wish to persist the gathered metrics, you must set the `add_ons.heapster.options.influxdb.pvc_name` option.

Plan file options:

| Plan file key | Description |
|---------------|-------------|
| add_ons.heapster.disable | Set to true if heapster should not be deployed during installation |
| add_ons.heapster.options.heapster.replicas  | Number of replicas for the heapster deployment |
| add_ons.heapster.options.influxdb.pvc_name | Name of a persistent volume claim that will be used by the influxdb databse for persistence. This PVC must be manually created after installation. |

### Package Manager
Helm is the official package manager for Kubernetes. KET includes the `helm` client-side binary in the distribution package. KET also installs the server-side agent, Tiller, on the cluster during installation. 

Plan file options:

| Plan file key | Description |
|---------------|-------------|
|add_ons.package_manager.disable | Set to true if the package manager should not be deployed during installation |
| add_ons.package_manager.provider | The package manager that should be deployed. Supported options: `helm`.|
