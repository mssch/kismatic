# Pluggable CNI

To support a larger set of environments and networks KET needs to support different CNI solutions.

Currently only Calico is supported. All clusters setup with KET with <`v1.4.x` used Calico as its CNI add-on, users were not given an option.

# Motivation 

Initial support for 2 additional CNI plugins will be added to KET.

[Weave Net](https://github.com/weaveworks/weave): 
* adds support for the [Microsoft Azure](https://azure.microsoft.com/) 
* multicast 
[Contiv](https://github.com/contiv)
* adds support for Cisco ACI 

# Specification
A new `add_ons.cni` section will be added
```
add_ons:
  cni:
    disabled: false
    provider: calico    #Options: calico, weave, contiv, custom
    options: #TBD
  heapster:
    disable: false
    options:
      heapster_replicas: 2
      influxdb_pvc_name: ""            
  package_manager:
    disable: false
    provider: helm                       
```

The `disabled:` field is valid, a user can choose not to install a CNI plugin.
When `disabled:` is `true` the CNI flags will not be set in the Kubernetes components and all pod validation for the other add-ons and the smoketest will be skipped. 
The CNI binaries and conf files will not be configured on the cluster.

When `provider: custom` is set, the CNI flags will be set in the Kubernetes components, however all pod validation for the other add-ons and the smoketest will be skipped.
The CNI binaries will be configured on the cluster as that is a common component for all CNI plugins. The conf file(if there is one) will needs to be placed by the user.

* Kubelet flags:
```
  --cni-bin-dir=/opt/cni/bin \                      # Do not set when cni.disabled == true
  --cni-conf-dir={{ network_plugin_dir }} \         # Do not set when cni.disabled == true
  --network-plugin=cni \                            # Do not set when cni.disabled == true
  --network-plugin-dir=${NETWORK_PLUGIN_DIR} \      # REMOVE, flag is no longer used
```

## Other Considerations
* A mechanism to install required [CNI binaries](https://github.com/containernetworking/cni) is required for certain network solutions. 
This can be a docker container that contains the required binaries, contains a shared volume on the host machine and copies them on the machine. (similar to Calico's approach)
Or a new package containing the CNI binaries(similar to kubeadm `kubernetes-cni` package)

## Plan File Changes
`cluster.networking.type` will be moved to `add_ons.cni.options.calico_mode`, KET will need to support the old flag and print a deprecation warning until a future release. 

# Upgrades
Some tests will be added using Weave, however we will rely on the CNI spec to provide parity between the different providers:
* "skunkworks" cluster
* minikube with the different supported OS
* upgrades in the future releases 

**Switching between CNI providers during upgrades will not be supported at this time.**