# Docker Registry

By default the Kismatic tools is configured to use the public [Docker Hub](https://hub.docker.com/) when pulling images during installation.
If the cluster does not have access to the Internet or the Hub, there are 2 options available in the plan to accommodate file when performing the initial Kubernetes cluster setup with Kismatic.

## User Configured Docker Registry
```
# plan file
docker_registry:                         
    address: myprivateregistry.com                               
    port: 8443              
    CA: /certs/ca.rt    
```
If you already have a Docker registry running and want to use that in the cluster, it is easiest to let the Kismatic tool configure all of the Docker daemons in that cluster.   

It is recommended to follow the [official docs](https://docs.docker.com/registry/deploying/#/running-a-domain-registry) when setting up your own docker registry.  

You will need to provide 3 pieces of information in the plan file:
* `address` the **reachable**(from all of the Kubernetes nodes) URL or IP where the Docker registry is running
* `port` number of the Docker registry
* `CA` that was used to sign the TLS certificates provided when starting the Docker registry

---
**If either option is selected, the Kismatic tool will attempt to push all of the images used during the install to the registry;
this allows for an installation that is completely self-contained and does not require an internet connection.**
