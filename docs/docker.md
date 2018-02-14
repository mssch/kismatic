# Docker Configuration

Kismatic supports docker as the underlying container runtime for your cluster. During the installation of your cluster, 
Kismatic will download and install docker on all nodes.

Starting with KET `v1.8.0` it is possible to use nodes with Docker already installed. When setting `docker.disable` to `true` KET will not try to install docker, and instead will validate that the `docker` command is available.

## Storage in RHEL/CentOS (KET v1.3.1 - v1.7.1)
The default storage driver that gets installed on RHEL/CentOS is `devicemapper` in `loop-lvm` mode. However, 
this is not recommended for production setups. Instead, `devicemapper` must be setup in `direct-lvm` mode. 

The `direct-lvm` mode requires a block storage device that will be used exclusively for docker's storage needs.
For this reason, this is an opt-in feature that can be enabled in the plan file:

``` yaml
docker:
  storage:
    direct_lvm:                          # Configure devicemapper in direct-lvm mode (RHEL/CentOS only).
      enabled: false
      block_device: ""                   # Path to the block device that will be used for direct-lvm mode. This device will be wiped and used exclusively by docker.
      enable_deferred_deletion: false    # Set to true if you want to enable deferred deletion when using direct-lvm mode.
```

When `direct-lvm` is enabled, Kismatic will create a logical volume configured as a thin pool to use as the backing storage
for docker. Creating a thin pool logical volume requires a "data" logical volume and a "metadata" logical volume.

The provided block storage device will be split 95% for data, and 1% for metadata, as recommended by Docker. The remaining 
space allows for auto-extending the logical volumes when they start to fill up as a temporary stop gap. 
The threshold for auto-extension is set to 80% of the size of the volume.

For more information on `devicemapper` see https://docs.docker.com/engine/userguide/storagedriver/device-mapper-driver/#configure-direct-lvm-mode-for-production

## Storage (KET v1.8.0+)

Starting in KET `v1.8.0` the docker storage options were modified to be more flexible. The old configuration will still work however it is recommended to start using the new format.

The new options `docker.storage.driver` and `docker.storage.opts` are equivalent to `"storage-driver"` and `"storage-opts"` in docker's `daemon.json`.

When `docker.storage.driver` is left empty, the docker daemon will determine the most suitable driver to use for the specific OS.

When setting `docker.storage.driver` to `devicemapper`, it is still possible for Kismatic to create the required storage device for `direct-lvm` mode. To do so, set `docker.storage.direct_lvm_block_device.path` to the absolute path of the block device. If the path is left empty with `devicemapper` driver, docker will be configured in `loop-lvm`.

``` yaml
docker:
  # Set to true if docker is already installed and configured.
  disable: false
  storage:
    # Leave empty to have docker automatically select the driver.
    driver: ""
    opts: {}
    # Used for setting up Device Mapper storage driver in direct-lvm mode.
    direct_lvm_block_device:
      # Absolute path to the block device that will be used for direct-lvm mode.
      # This device will be wiped and used exclusively by docker.
      path: ""
      thinpool_percent: "95"
      thinpool_metapercent: "1"
      thinpool_autoextend_threshold: "80"
      thinpool_autoextend_percent: "20"
```

### Migrating to the new format

**Old config (Deprecated)**

``` yaml
docker:
  storage:
    direct_lvm:
      enabled: true
      block_device: "/dev/xvdb"
```

**New Equivalent Config**
```yaml
docker:
  storage:
    driver: devicemapper
    opts:
      dm.thinpooldev: "/dev/mapper/docker-thinpool"
      dm.use_deferred_deletion: "false"
      dm.use_deferred_removal: "true"
    direct_lvm_block_device:
      path: "/dev/xvdb"
      thinpool_percent: "95"
      thinpool_metapercent: "1"
      thinpool_autoextend_threshold: "80"
      thinpool_autoextend_percent: "20"
```

## Log Driver (KET v1.7.0+)
``` yaml
docker
  # Set to true if docker is already installed and configured.
  disable: false
  logs:
    driver: "json-file"
    opts:
      max-file: "1"
      max-size: "50m"
```

In the plan file the `docker.logs` field allows to set the docker daemon log driver [options](https://docs.docker.com/engine/admin/logging/overview/). The specified options from the plan file get set in `/etc/docker/daemon.json`. 

This is an advanced feature and the values provided will not be validated, please refer to the documentation for your specific driver for the valid options.