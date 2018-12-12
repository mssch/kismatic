# Plan File Reference
## Index
* [cluster](#cluster)
  * [name](#clustername)
  * [version](#clusterversion)
  * [admin_password _(deprecated)_](#clusteradmin_password-deprecated)
  * [disable_package_installation](#clusterdisable_package_installation)
  * [allow_package_installation _(deprecated)_](#clusterallow_package_installation-deprecated)
  * [disconnected_installation](#clusterdisconnected_installation)
  * [networking](#clusternetworking)
    * [type _(deprecated)_](#clusternetworkingtype-deprecated)
    * [pod_cidr_block](#clusternetworkingpod_cidr_block)
    * [service_cidr_block](#clusternetworkingservice_cidr_block)
    * [update_hosts_files](#clusternetworkingupdate_hosts_files)
    * [http_proxy](#clusternetworkinghttp_proxy)
    * [https_proxy](#clusternetworkinghttps_proxy)
    * [no_proxy](#clusternetworkingno_proxy)
  * [certificates](#clustercertificates)
    * [expiry](#clustercertificatesexpiry)
    * [ca_expiry](#clustercertificatesca_expiry)
    * [apiserver_cert_extra_sans](#clustercertificatesapiserver_cert_extra_sans)
  * [ssh](#clusterssh)
    * [user](#clustersshuser)
    * [ssh_key](#clustersshssh_key)
    * [ssh_port](#clustersshssh_port)
  * [kube_apiserver](#clusterkube_apiserver)
    * [option_overrides](#clusterkube_apiserveroption_overrides)
  * [kube_controller_manager](#clusterkube_controller_manager)
    * [option_overrides](#clusterkube_controller_manageroption_overrides)
  * [kube_scheduler](#clusterkube_scheduler)
    * [option_overrides](#clusterkube_scheduleroption_overrides)
  * [kube_proxy](#clusterkube_proxy)
    * [option_overrides](#clusterkube_proxyoption_overrides)
  * [kubelet](#clusterkubelet)
    * [option_overrides](#clusterkubeletoption_overrides)
  * [cloud_provider](#clustercloud_provider)
    * [provider](#clustercloud_providerprovider)
    * [config](#clustercloud_providerconfig)
* [docker](#docker)
  * [disable](#dockerdisable)
  * [logs](#dockerlogs)
    * [driver](#dockerlogsdriver)
    * [opts](#dockerlogsopts)
  * [storage](#dockerstorage)
    * [driver](#dockerstoragedriver)
    * [opts](#dockerstorageopts)
    * [direct_lvm_block_device](#dockerstoragedirect_lvm_block_device)
      * [path](#dockerstoragedirect_lvm_block_devicepath)
      * [thinpool_percent](#dockerstoragedirect_lvm_block_devicethinpool_percent)
      * [thinpool_metapercent](#dockerstoragedirect_lvm_block_devicethinpool_metapercent)
      * [thinpool_autoextend_threshold](#dockerstoragedirect_lvm_block_devicethinpool_autoextend_threshold)
      * [thinpool_autoextend_percent](#dockerstoragedirect_lvm_block_devicethinpool_autoextend_percent)
    * [direct_lvm _(deprecated)_](#dockerstoragedirect_lvm-deprecated)
      * [enabled](#dockerstoragedirect_lvmenabled)
      * [block_device](#dockerstoragedirect_lvmblock_device)
      * [enable_deferred_deletion](#dockerstoragedirect_lvmenable_deferred_deletion)
* [docker_registry](#docker_registry)
  * [server](#docker_registryserver)
  * [address _(deprecated)_](#docker_registryaddress-deprecated)
  * [port _(deprecated)_](#docker_registryport-deprecated)
  * [CA](#docker_registryCA)
  * [username](#docker_registryusername)
  * [password](#docker_registrypassword)
* [additional_files](#additional_files)
  * [hosts](#additional_fileshosts)
  * [source](#additional_filessource)
  * [destination](#additional_filesdestination)
  * [skip_validation](#additional_filesskip_validation)
* [add_ons](#add_ons)
  * [cni](#add_onscni)
    * [disable](#add_onscnidisable)
    * [provider](#add_onscniprovider)
    * [options](#add_onscnioptions)
      * [portmap](#add_onscnioptionsportmap)
        * [disable](#add_onscnioptionsportmapdisable)
      * [calico](#add_onscnioptionscalico)
        * [mode](#add_onscnioptionscalicomode)
        * [log_level](#add_onscnioptionscalicolog_level)
        * [workload_mtu](#add_onscnioptionscalicoworkload_mtu)
        * [felix_input_mtu](#add_onscnioptionscalicofelix_input_mtu)
        * [ip_autodetection_method](#add_onscnioptionscalicoip_autodetection_method)
      * [weave](#add_onscnioptionsweave)
        * [password](#add_onscnioptionsweavepassword)
  * [dns](#add_onsdns)
    * [disable](#add_onsdnsdisable)
    * [provider](#add_onsdnsprovider)
    * [options](#add_onsdnsoptions)
      * [replicas](#add_onsdnsoptionsreplicas)
  * [heapster](#add_onsheapster)
    * [disable](#add_onsheapsterdisable)
    * [options](#add_onsheapsteroptions)
      * [heapster](#add_onsheapsteroptionsheapster)
        * [replicas](#add_onsheapsteroptionsheapsterreplicas)
        * [service_type](#add_onsheapsteroptionsheapsterservice_type)
        * [sink](#add_onsheapsteroptionsheapstersink)
      * [influxdb](#add_onsheapsteroptionsinfluxdb)
        * [pvc_name](#add_onsheapsteroptionsinfluxdbpvc_name)
      * [heapster_replicas _(deprecated)_](#add_onsheapsteroptionsheapster_replicas-deprecated)
      * [influxdb_pvc_name _(deprecated)_](#add_onsheapsteroptionsinfluxdb_pvc_name-deprecated)
  * [metrics_server](#add_onsmetrics_server)
    * [disable](#add_onsmetrics_serverdisable)
  * [dashboard](#add_onsdashboard)
    * [disable](#add_onsdashboarddisable)
    * [options](#add_onsdashboardoptions)
      * [service_type](#add_onsdashboardoptionsservice_type)
      * [node_port](#add_onsdashboardoptionsnode_port)
  * [package_manager](#add_onspackage_manager)
    * [disable](#add_onspackage_managerdisable)
    * [provider](#add_onspackage_managerprovider)
    * [options](#add_onspackage_manageroptions)
      * [helm](#add_onspackage_manageroptionshelm)
        * [namespace](#add_onspackage_manageroptionshelmnamespace)
  * [rescheduler](#add_onsrescheduler)
    * [disable](#add_onsreschedulerdisable)
* [features _(deprecated)_](#features-deprecated)
  * [package_manager _(deprecated)_](#featurespackage_manager-deprecated)
    * [enabled _(deprecated)_](#featurespackage_managerenabled-deprecated)
* [etcd](#etcd)
  * [expected_count](#etcdexpected_count)
  * [nodes](#etcdnodes)
    * [host](#etcdnodeshost)
    * [ip](#etcdnodesip)
    * [internalip](#etcdnodesinternalip)
    * [labels](#etcdnodeslabels)
    * [taints](#etcdnodestaints)
      * [key](#etcdnodestaintskey)
      * [value](#etcdnodestaintsvalue)
      * [effect](#etcdnodestaintseffect)
    * [kubelet](#etcdnodeskubelet)
      * [option_overrides](#etcdnodeskubeletoption_overrides)
* [master](#master)
  * [load_balancer](#masterload_balancer)
  * [expected_count](#masterexpected_count)
  * [load_balanced_fqdn _(deprecated)_](#masterload_balanced_fqdn-deprecated)
  * [load_balanced_short_name _(deprecated)_](#masterload_balanced_short_name-deprecated)
  * [nodes](#masternodes)
    * [host](#masternodeshost)
    * [ip](#masternodesip)
    * [internalip](#masternodesinternalip)
    * [labels](#masternodeslabels)
    * [taints](#masternodestaints)
      * [key](#masternodestaintskey)
      * [value](#masternodestaintsvalue)
      * [effect](#masternodestaintseffect)
    * [kubelet](#masternodeskubelet)
      * [option_overrides](#masternodeskubeletoption_overrides)
* [worker](#worker)
  * [expected_count](#workerexpected_count)
  * [nodes](#workernodes)
    * [host](#workernodeshost)
    * [ip](#workernodesip)
    * [internalip](#workernodesinternalip)
    * [labels](#workernodeslabels)
    * [taints](#workernodestaints)
      * [key](#workernodestaintskey)
      * [value](#workernodestaintsvalue)
      * [effect](#workernodestaintseffect)
    * [kubelet](#workernodeskubelet)
      * [option_overrides](#workernodeskubeletoption_overrides)
* [ingress](#ingress)
  * [expected_count](#ingressexpected_count)
  * [nodes](#ingressnodes)
    * [host](#ingressnodeshost)
    * [ip](#ingressnodesip)
    * [internalip](#ingressnodesinternalip)
    * [labels](#ingressnodeslabels)
    * [taints](#ingressnodestaints)
      * [key](#ingressnodestaintskey)
      * [value](#ingressnodestaintsvalue)
      * [effect](#ingressnodestaintseffect)
    * [kubelet](#ingressnodeskubelet)
      * [option_overrides](#ingressnodeskubeletoption_overrides)
* [storage](#storage)
  * [expected_count](#storageexpected_count)
  * [nodes](#storagenodes)
    * [host](#storagenodeshost)
    * [ip](#storagenodesip)
    * [internalip](#storagenodesinternalip)
    * [labels](#storagenodeslabels)
    * [taints](#storagenodestaints)
      * [key](#storagenodestaintskey)
      * [value](#storagenodestaintsvalue)
      * [effect](#storagenodestaintseffect)
    * [kubelet](#storagenodeskubelet)
      * [option_overrides](#storagenodeskubeletoption_overrides)
* [nfs](#nfs)
  * [nfs_volume](#nfsnfs_volume)
    * [nfs_host](#nfsnfs_volumenfs_host)
    * [mount_path](#nfsnfs_volumemount_path)
##  cluster

 Kubernetes cluster configuration 

###  cluster.name

 Name of the cluster to be used when generating assets that require a cluster name, such as kubeconfig files and certificates. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  cluster.version

 The Kubernetes version to install. If left blank will be set to the latest tested version. Only a single Minor version is supported with. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `v1.10.11` | 

###  cluster.admin_password _(deprecated)_

 The password for the admin user. If provided, ABAC will be enabled in the cluster. This field will be removed completely in a future release. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  cluster.disable_package_installation

 Whether KET should install the packages on the cluster nodes. When true, KET will not install the required packages. Instead, it will verify that the packages have been installed by the operator. 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

###  cluster.allow_package_installation _(deprecated)_

 Whether KET should install the packages on the cluster nodes. Use DisablePackageInstallation instead. 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

###  cluster.disconnected_installation

 Whether the cluster nodes are disconnected from the internet. When set to `true`, internal package repositories and a container image registry are required for installation. 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

###  cluster.networking

 The Networking configuration for the cluster. 

###  cluster.networking.type _(deprecated)_

 The datapath technique that should be configured in Calico. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `overlay` | 
| **Options** |  `overlay`, `routed`

###  cluster.networking.pod_cidr_block

 The pod network's CIDR block. For example: `172.16.0.0/16` 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  cluster.networking.service_cidr_block

 The Kubernetes service network's CIDR block. For example: `172.20.0.0/16` 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  cluster.networking.update_hosts_files

 Whether the /etc/hosts file should be updated on the cluster nodes. When set to true, KET will update the hosts file on all nodes to include entries for all other nodes in the cluster. 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

###  cluster.networking.http_proxy

 The URL of the proxy that should be used for HTTP connections. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  cluster.networking.https_proxy

 The URL of the proxy that should be used for HTTPS connections. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  cluster.networking.no_proxy

 Comma-separated list of host names and/or IPs for which connections should not go through a proxy. All nodes' 'host' and 'IPs' are always set. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  cluster.certificates

 The Certificates configuration for the cluster. 

###  cluster.certificates.expiry

 The length of time that the generated certificates should be valid for. For example: "17520h" for 2 years. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  cluster.certificates.ca_expiry

 The length of time that the generated Certificate Authority should be valid for. For example: "17520h" for 2 years. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  cluster.certificates.apiserver_cert_extra_sans

 Comma-separated list of Subject Alternative Names (SANs) to use for the API Server serving certificate. Can be both IP addresses and DNS names. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  cluster.ssh

 The SSH configuration for the cluster nodes. 

###  cluster.ssh.user

 The user for accessing the cluster nodes via SSH. This user requires sudo elevation privileges on the cluster nodes. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  cluster.ssh.ssh_key

 The absolute path of the SSH key that should be used for accessing the cluster nodes via SSH. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  cluster.ssh.ssh_port

 The port number on which cluster nodes are listening for SSH connections. 

| | |
|----------|-----------------|
| **Kind** |  int |
| **Required** |  Yes |
| **Default** | ` ` | 

###  cluster.kube_apiserver

 Kubernetes API Server configuration. 

###  cluster.kube_apiserver.option_overrides

 Listing of option overrides that are to be applied to the Kubernetes API server configuration. This is an advanced feature that can prevent the API server from starting up if invalid configuration is provided. 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

###  cluster.kube_controller_manager

 Kubernetes Controller Manager configuration. 

###  cluster.kube_controller_manager.option_overrides

 Listing of option overrides that are to be applied to the Kubernetes Controller Manager configuration. This is an advanced feature that can prevent the Controller Manager from starting up if invalid configuration is provided. 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

###  cluster.kube_scheduler

 Kubernetes Scheduler configuration. 

###  cluster.kube_scheduler.option_overrides

 Listing of option overrides that are to be applied to the Kubernetes Scheduler configuration. This is an advanced feature that can prevent the Scheduler from starting up if invalid configuration is provided. 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

###  cluster.kube_proxy

 Kubernetes Proxy configuration. 

###  cluster.kube_proxy.option_overrides

 Listing of option overrides that are to be applied to the Kubernetes Proxy configuration. This is an advanced feature that can prevent the Proxy from starting up if invalid configuration is provided. 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

###  cluster.kubelet

 Kubelet configuration applied to all nodes. 

###  cluster.kubelet.option_overrides

 Listing of option overrides that are to be applied to the Kubelet configurations. This is an advanced feature that can prevent the Kubelet from starting up if invalid configuration is provided. 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

###  cluster.cloud_provider

 The CloudProvider configuration for the cluster. 

###  cluster.cloud_provider.provider

 The cloud provider that should be set in the Kubernetes components 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 
| **Options** |  `aws`, `azure`, `cloudstack`, `fake`, `gce`, `mesos`, `openstack`, `ovirt`, `photon`, `rackspace`, `vsphere`

###  cluster.cloud_provider.config

 Path to the cloud provider config file. This will be copied to all the machines in the cluster 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

##  docker

 Configuration for the docker engine installed by KET 

###  docker.disable

 Set to true to disable the installation of docker container runtime on the nodes. The installer will validate that docker is installed and running prior to proceeding. Use this option if a different version of docker from the included one is required. 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

###  docker.logs

 Log configuration for the docker engine. 

###  docker.logs.driver

 Docker logging driver, more details https://docs.docker.com/engine/admin/logging/overview/. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `json-file` | 

###  docker.logs.opts

 Driver specific options. 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

###  docker.storage

 Storage configuration for the docker engine. 

###  docker.storage.driver

 Docker storage driver, more details https://docs.docker.com/engine/userguide/storagedriver/. Leave empty to have docker automatically select the driver. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `'empty'` | 

###  docker.storage.opts

 Driver specific options 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

###  docker.storage.direct_lvm_block_device

 DirectLVMBlockDevice is the configuration required for setting up Device Mapper storage driver in direct-lvm mode. Refer to https://docs.docker.com/v17.03/engine/userguide/storagedriver/device-mapper-driver/#manage-devicemapper docs. 

###  docker.storage.direct_lvm_block_device.path

 The path to the block device. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  docker.storage.direct_lvm_block_device.thinpool_percent

 The percentage of space to use for storage from the passed in block device. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `95` | 

###  docker.storage.direct_lvm_block_device.thinpool_metapercent

 The percentage of space to for metadata storage from the passed in block device. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `1` | 

###  docker.storage.direct_lvm_block_device.thinpool_autoextend_threshold

 The threshold for when lvm should automatically extend the thin pool as a percentage of the total storage space. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `80` | 

###  docker.storage.direct_lvm_block_device.thinpool_autoextend_percent

 The percentage to increase the thin pool by when an autoextend is triggered. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `20` | 

###  docker.storage.direct_lvm _(deprecated)_

 DirectLVM is the configuration required for setting up device mapper in direct-lvm mode. 

###  docker.storage.direct_lvm.enabled

 Whether the direct_lvm mode of the devicemapper storage driver should be enabled. When set to true, a dedicated block storage device must be available on each cluster node. 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

###  docker.storage.direct_lvm.block_device

 The path to the block storage device that will be used by the devicemapper storage driver. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  docker.storage.direct_lvm.enable_deferred_deletion

 Whether deferred deletion should be enabled when using devicemapper in direct_lvm mode. 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

##  docker_registry

 Docker registry configuration 

###  docker_registry.server

 The hostname or IP address and port of a private container image registry. Do not include http or https. When performing a disconnected installation, this registry will be used to fetch all the required container images. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  docker_registry.address _(deprecated)_

 The hostname or IP address of a private container image registry. When performing a disconnected installation, this registry will be used to fetch all the required container images. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  docker_registry.port _(deprecated)_

 The port on which the private container image registry is listening on. 

| | |
|----------|-----------------|
| **Kind** |  int |
| **Required** |  No |
| **Default** | ` ` | 

###  docker_registry.CA

 The absolute path of the Certificate Authority that should be installed on all cluster nodes that have a docker daemon. This is required to establish trust between the daemons and the private registry when the registry is using a self-signed certificate. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  docker_registry.username

 The username that should be used when connecting to a registry that has authentication enabled. Otherwise leave blank for unauthenticated access. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  docker_registry.password

 The password that should be used when connecting to a registry that has authentication enabled. Otherwise leave blank for unauthenticated access. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

##  additional_files

 A set of files or directories to copy from the local machine to any of the nodes in the cluster. 

###  additional_files.hosts

 Hostname or role where additional files or directories will be copied. 

###  additional_files.source

 Path to the file or directory on local machine. Must be an absolute path. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  additional_files.destination

 Path to the file or directory on remote machine, where file will be copied. Must be an absolute path. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  additional_files.skip_validation

 Set to true if validation will be run before the file exists on the local machine. Useful for files generated at install time, ie. assets in generated/ directory. 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

##  add_ons

 Add on configuration 

###  add_ons.cni

 The Container Networking Interface (CNI) add-on configuration. 

###  add_ons.cni.disable

 Whether the CNI add-on is disabled. When set to true, CNI will not be installed on the cluster. Furthermore, the smoke test and any validation that depends on a functional pod network will be skipped. 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

###  add_ons.cni.provider

 The CNI provider that should be installed on the cluster. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `calico` | 
| **Options** |  `calico`, `weave`, `contiv`, `custom`

###  add_ons.cni.options

 The CNI options that can be configured for each CNI provider. 

###  add_ons.cni.options.portmap

 The options that can be configured for the Portmap CNI provider. 

###  add_ons.cni.options.portmap.disable

 Disable the portmap CNI plugin 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

###  add_ons.cni.options.calico

 The options that can be configured for the Calico CNI provider. 

###  add_ons.cni.options.calico.mode

 The datapath technique that should be configured in Calico. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `overlay` | 
| **Options** |  `overlay`, `routed`

###  add_ons.cni.options.calico.log_level

 The logging level for the CNI plugin 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `info` | 
| **Options** |  `warning`, `info`, `debug`

###  add_ons.cni.options.calico.workload_mtu

 MTU for the workload interface, configures the CNI config. 

| | |
|----------|-----------------|
| **Kind** |  int |
| **Required** |  No |
| **Default** | `1500` | 

###  add_ons.cni.options.calico.felix_input_mtu

 MTU for the tunnel device used if IPIP is enabled. 

| | |
|----------|-----------------|
| **Kind** |  int |
| **Required** |  No |
| **Default** | `1440` | 

###  add_ons.cni.options.calico.ip_autodetection_method

 IPAutodetectionMethod is used to detect the IPv4 address of the host. The value gets set in IP_AUTODETECTION_METHOD variable in the pod. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `first-found` | 

###  add_ons.cni.options.weave

 The options that can be configured for the Weave CNI provider. 

###  add_ons.cni.options.weave.password

 The password to use for network traffic encryption. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  add_ons.dns

 The DNS add-on configuration. 

###  add_ons.dns.disable

 Whether the DNS add-on should be disabled. When set to true, no DNS solution will be deployed on the cluster. 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

###  add_ons.dns.provider

 This property indicates the in-cluster DNS provider. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | `kubedns` | 
| **Options** |  `kubedns`, `coredns`

###  add_ons.dns.options

 The options that can be configured for the cluster DNS add-on 

###  add_ons.dns.options.replicas

 Number of cluster DNS replicas that should be scheduled on the cluster. 

| | |
|----------|-----------------|
| **Kind** |  int |
| **Required** |  No |
| **Default** | `2` | 

###  add_ons.heapster

 The Heapster Monitoring add-on configuration. 

###  add_ons.heapster.disable

 Whether the Heapster add-on should be disabled. When set to true, Heapster and InfluxDB will not be deployed on the cluster. 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

###  add_ons.heapster.options

 The options that can be configured for the Heapster add-on 

###  add_ons.heapster.options.heapster

 The Heapster configuration options. 

###  add_ons.heapster.options.heapster.replicas

 Number of Heapster replicas that should be scheduled on the cluster. 

| | |
|----------|-----------------|
| **Kind** |  int |
| **Required** |  No |
| **Default** | `2` | 

###  add_ons.heapster.options.heapster.service_type

 Kubernetes service type of the Heapster service. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `ClusterIP` | 
| **Options** |  `ClusterIP`, `NodePort`, `LoadBalancer`, `ExternalName`

###  add_ons.heapster.options.heapster.sink

 URL of the backend store that will be used as the Heapster sink. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `influxdb:http://heapster-influxdb.kube-system.svc:8086` | 

###  add_ons.heapster.options.influxdb

 The InfluxDB configuration options. 

###  add_ons.heapster.options.influxdb.pvc_name

 Name of the Persistent Volume Claim that will be used by InfluxDB. This PVC must be created after the installation. If not set, InfluxDB will be configured with ephemeral storage. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  add_ons.heapster.options.heapster_replicas _(deprecated)_

 Number of Heapster replicas that should be scheduled on the cluster. 

| | |
|----------|-----------------|
| **Kind** |  int |
| **Required** |  No |
| **Default** | ` ` | 

###  add_ons.heapster.options.influxdb_pvc_name _(deprecated)_

 Name of the Persistent Volume Claim that will be used by InfluxDB. When set, this PVC must be created after the installation. If not set, InfluxDB will be configured with ephemeral storage. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  add_ons.metrics_server

 Metrics Server add-on configuration. A cluster-wide aggregator of resource usage data. Required for Horizontal Pod Autoscaler to function properly. 

###  add_ons.metrics_server.disable

 Whether the metrics-server add-on should be disabled. When set to true, metrics-server will not be deployed on the cluster. 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

###  add_ons.dashboard

 The Dashboard add-on configuration. 

###  add_ons.dashboard.disable

 Whether the dashboard add-on should be disabled. When set to true, the Kubernetes Dashboard will not be installed on the cluster. 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

###  add_ons.dashboard.options

 The options that can be configured for the Dashboard add-on 

###  add_ons.dashboard.options.service_type

 Kubernetes service type of the Dashboard service. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `ClusterIP` | 
| **Options** |  `ClusterIP`, `NodePort`, `LoadBalancer`, `ExternalName`

###  add_ons.dashboard.options.node_port

 When using NodePort set the port to use. When left empty Kubernetes will allocate a random port. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `''` | 

###  add_ons.package_manager

 The PackageManager add-on configuration. 

###  add_ons.package_manager.disable

 Whether the package manager add-on should be disabled. When set to true, the package manager will not be installed on the cluster. 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

###  add_ons.package_manager.provider

 This property indicates the package manager provider. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 
| **Options** |  `helm`

###  add_ons.package_manager.options

 The PackageManager options. 

###  add_ons.package_manager.options.helm

 Helm PackageManager options 

###  add_ons.package_manager.options.helm.namespace

 Namespace to deploy tiller 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | `kube-system` | 

###  add_ons.rescheduler

 The Rescheduler add-on configuration. Because the Rescheduler does not have leader election and therefore can only run as a single instance in a cluster, it will be deployed as a static pod on the first master. More information about the Rescheduler can be found here: https://kubernetes.io/docs/tasks/administer-cluster/guaranteed-scheduling-critical-addon-pods/ 

###  add_ons.rescheduler.disable

 Whether the pod rescheduler add-on should be disabled. When set to true, the rescheduler will not be installed on the cluster. 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

##  features _(deprecated)_

 Feature configuration 

###  features.package_manager _(deprecated)_

 The PackageManager feature configuration. 

###  features.package_manager.enabled _(deprecated)_

 Whether the package manager add-on should be enabled. 

| | |
|----------|-----------------|
| **Kind** |  bool |
| **Required** |  No |
| **Default** | `false` | 

##  etcd

 Etcd nodes of the cluster 

###  etcd.expected_count

 Number of nodes. 

| | |
|----------|-----------------|
| **Kind** |  int |
| **Required** |  Yes |
| **Default** | ` ` | 

###  etcd.nodes

 List of nodes. 

###  etcd.nodes.host

 The hostname of the node. The hostname is verified in the validation phase of the installation. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  etcd.nodes.ip

 The IP address of the node. This is the IP address that will be used to connect to the node over SSH. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  etcd.nodes.internalip

 The internal (or private) IP address of the node. If set, this IP will be used when configuring cluster components. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  etcd.nodes.labels

 Labels to add when installing the node in the cluster. If a node is defined under multiple roles, the labels for that node will be merged. If a label is repeated for the same node, only one will be used in this order: etcd,master,worker,ingress,storage roles where 'storage' has the highest precedence. It is recommended to use reverse-DNS notation to avoid collision with other labels. 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

###  etcd.nodes.taints

 Taints to add when installing the node in the cluster. If a node is defined under multiple roles, the taints for that node will be merged. If a taint is repeated for the same node, only one will be used in this order: etcd,master,worker,ingress,storage roles where 'storage' has the highest precedence. 

###  etcd.nodes.taints.key

 Key for the taint 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  etcd.nodes.taints.value

 Value for the taint 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  etcd.nodes.taints.effect

 Effect for the taint 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 
| **Options** |  `NoSchedule`, `PreferNoSchedule`, `NoExecute`

###  etcd.nodes.kubelet

 Kubelet configuration applied to this node. If a node is repeated for multiple roles, the overrides cannot be different. 

###  etcd.nodes.kubelet.option_overrides

 Listing of option overrides that are to be applied to the Kubelet configurations. This is an advanced feature that can prevent the Kubelet from starting up if invalid configuration is provided. 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

##  master

 Master nodes of the cluster 

###  master.load_balancer

 The IP or DNS and Port of the load balancer that is fronting multiple master nodes. In the case where there no load balancer this can be set to the IP address of the master node with port '6443'. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  master.expected_count

 Number of master nodes that are part of the cluster. 

| | |
|----------|-----------------|
| **Kind** |  int |
| **Required** |  Yes |
| **Default** | ` ` | 

###  master.load_balanced_fqdn _(deprecated)_

 The FQDN of the load balancer that is fronting multiple master nodes. In the case where there is only one master node, this can be set to the IP address of the master node. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  master.load_balanced_short_name _(deprecated)_

 The short name of the load balancer that is fronting multiple master nodes. In the case where there is only one master node, this can be set to the IP address of the master nodes. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  master.nodes

 List of master nodes that are part of the cluster. 

###  master.nodes.host

 The hostname of the node. The hostname is verified in the validation phase of the installation. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  master.nodes.ip

 The IP address of the node. This is the IP address that will be used to connect to the node over SSH. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  master.nodes.internalip

 The internal (or private) IP address of the node. If set, this IP will be used when configuring cluster components. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  master.nodes.labels

 Labels to add when installing the node in the cluster. If a node is defined under multiple roles, the labels for that node will be merged. If a label is repeated for the same node, only one will be used in this order: etcd,master,worker,ingress,storage roles where 'storage' has the highest precedence. It is recommended to use reverse-DNS notation to avoid collision with other labels. 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

###  master.nodes.taints

 Taints to add when installing the node in the cluster. If a node is defined under multiple roles, the taints for that node will be merged. If a taint is repeated for the same node, only one will be used in this order: etcd,master,worker,ingress,storage roles where 'storage' has the highest precedence. 

###  master.nodes.taints.key

 Key for the taint 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  master.nodes.taints.value

 Value for the taint 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  master.nodes.taints.effect

 Effect for the taint 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 
| **Options** |  `NoSchedule`, `PreferNoSchedule`, `NoExecute`

###  master.nodes.kubelet

 Kubelet configuration applied to this node. If a node is repeated for multiple roles, the overrides cannot be different. 

###  master.nodes.kubelet.option_overrides

 Listing of option overrides that are to be applied to the Kubelet configurations. This is an advanced feature that can prevent the Kubelet from starting up if invalid configuration is provided. 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

##  worker

 Worker nodes of the cluster 

###  worker.expected_count

 Number of nodes. 

| | |
|----------|-----------------|
| **Kind** |  int |
| **Required** |  Yes |
| **Default** | ` ` | 

###  worker.nodes

 List of nodes. 

###  worker.nodes.host

 The hostname of the node. The hostname is verified in the validation phase of the installation. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  worker.nodes.ip

 The IP address of the node. This is the IP address that will be used to connect to the node over SSH. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  worker.nodes.internalip

 The internal (or private) IP address of the node. If set, this IP will be used when configuring cluster components. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  worker.nodes.labels

 Labels to add when installing the node in the cluster. If a node is defined under multiple roles, the labels for that node will be merged. If a label is repeated for the same node, only one will be used in this order: etcd,master,worker,ingress,storage roles where 'storage' has the highest precedence. It is recommended to use reverse-DNS notation to avoid collision with other labels. 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

###  worker.nodes.taints

 Taints to add when installing the node in the cluster. If a node is defined under multiple roles, the taints for that node will be merged. If a taint is repeated for the same node, only one will be used in this order: etcd,master,worker,ingress,storage roles where 'storage' has the highest precedence. 

###  worker.nodes.taints.key

 Key for the taint 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  worker.nodes.taints.value

 Value for the taint 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  worker.nodes.taints.effect

 Effect for the taint 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 
| **Options** |  `NoSchedule`, `PreferNoSchedule`, `NoExecute`

###  worker.nodes.kubelet

 Kubelet configuration applied to this node. If a node is repeated for multiple roles, the overrides cannot be different. 

###  worker.nodes.kubelet.option_overrides

 Listing of option overrides that are to be applied to the Kubelet configurations. This is an advanced feature that can prevent the Kubelet from starting up if invalid configuration is provided. 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

##  ingress

 Ingress nodes of the cluster 

###  ingress.expected_count

 Number of nodes. 

| | |
|----------|-----------------|
| **Kind** |  int |
| **Required** |  Yes |
| **Default** | ` ` | 

###  ingress.nodes

 List of nodes. 

###  ingress.nodes.host

 The hostname of the node. The hostname is verified in the validation phase of the installation. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  ingress.nodes.ip

 The IP address of the node. This is the IP address that will be used to connect to the node over SSH. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  ingress.nodes.internalip

 The internal (or private) IP address of the node. If set, this IP will be used when configuring cluster components. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  ingress.nodes.labels

 Labels to add when installing the node in the cluster. If a node is defined under multiple roles, the labels for that node will be merged. If a label is repeated for the same node, only one will be used in this order: etcd,master,worker,ingress,storage roles where 'storage' has the highest precedence. It is recommended to use reverse-DNS notation to avoid collision with other labels. 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

###  ingress.nodes.taints

 Taints to add when installing the node in the cluster. If a node is defined under multiple roles, the taints for that node will be merged. If a taint is repeated for the same node, only one will be used in this order: etcd,master,worker,ingress,storage roles where 'storage' has the highest precedence. 

###  ingress.nodes.taints.key

 Key for the taint 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  ingress.nodes.taints.value

 Value for the taint 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  ingress.nodes.taints.effect

 Effect for the taint 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 
| **Options** |  `NoSchedule`, `PreferNoSchedule`, `NoExecute`

###  ingress.nodes.kubelet

 Kubelet configuration applied to this node. If a node is repeated for multiple roles, the overrides cannot be different. 

###  ingress.nodes.kubelet.option_overrides

 Listing of option overrides that are to be applied to the Kubelet configurations. This is an advanced feature that can prevent the Kubelet from starting up if invalid configuration is provided. 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

##  storage

 Storage nodes of the cluster. 

###  storage.expected_count

 Number of nodes. 

| | |
|----------|-----------------|
| **Kind** |  int |
| **Required** |  Yes |
| **Default** | ` ` | 

###  storage.nodes

 List of nodes. 

###  storage.nodes.host

 The hostname of the node. The hostname is verified in the validation phase of the installation. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  storage.nodes.ip

 The IP address of the node. This is the IP address that will be used to connect to the node over SSH. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  storage.nodes.internalip

 The internal (or private) IP address of the node. If set, this IP will be used when configuring cluster components. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  storage.nodes.labels

 Labels to add when installing the node in the cluster. If a node is defined under multiple roles, the labels for that node will be merged. If a label is repeated for the same node, only one will be used in this order: etcd,master,worker,ingress,storage roles where 'storage' has the highest precedence. It is recommended to use reverse-DNS notation to avoid collision with other labels. 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

###  storage.nodes.taints

 Taints to add when installing the node in the cluster. If a node is defined under multiple roles, the taints for that node will be merged. If a taint is repeated for the same node, only one will be used in this order: etcd,master,worker,ingress,storage roles where 'storage' has the highest precedence. 

###  storage.nodes.taints.key

 Key for the taint 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  storage.nodes.taints.value

 Value for the taint 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 

###  storage.nodes.taints.effect

 Effect for the taint 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  No |
| **Default** | ` ` | 
| **Options** |  `NoSchedule`, `PreferNoSchedule`, `NoExecute`

###  storage.nodes.kubelet

 Kubelet configuration applied to this node. If a node is repeated for multiple roles, the overrides cannot be different. 

###  storage.nodes.kubelet.option_overrides

 Listing of option overrides that are to be applied to the Kubelet configurations. This is an advanced feature that can prevent the Kubelet from starting up if invalid configuration is provided. 

| | |
|----------|-----------------|
| **Kind** |  map[string]string |
| **Required** |  No |
| **Default** | ` ` | 

##  nfs

 NFS volumes of the cluster. 

###  nfs.nfs_volume

 List of NFS volumes that should be attached to the cluster during the installation. 

###  nfs.nfs_volume.nfs_host

 The hostname or IP of the NFS volume. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

###  nfs.nfs_volume.mount_path

 The path where the NFS volume should be mounted. 

| | |
|----------|-----------------|
| **Kind** |  string |
| **Required** |  Yes |
| **Default** | ` ` | 

