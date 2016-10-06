#!/bin/bash
K8S_VERSION=1.4.0
DOCKER_VERSION=1.11.2

rm -rf out/DEBs
mkdir -p out/DEBs

# build Kubernetes
# #deb
# master
fpm -s dir -n "kismatic-kubernetes-master" \
-p out/DEBs \
-C source/ \
-d 'kismatic-docker-engine = 1.11.2-0~xenial' \
-d 'bridge-utils' \
-v $K8S_VERSION  \
-a amd64 \
-t deb \
--license "Apache Software License 2.0" \
--maintainer "Apprenda <info@apprenda.com>" \
--vendor "Apprenda" \
--description "Kubernetes master binaries" \
--url "https://apprenda.com/" \
kubernetes/apiserver/bin/kube-apiserver=/usr/bin/kube-apiserver \
kubernetes/kubelet/bin/kubelet=/usr/bin/kubelet \
kubernetes/proxy/bin/kube-proxy=/usr/bin/kube-proxy \
kubernetes/scheduler/bin/kube-scheduler=/usr/bin/kube-scheduler \
kubernetes/controller-manager/bin/kube-controller-manager=/usr/bin/kube-controller-manager \
kubernetes/kubectl/bin/kubectl=/usr/bin/kubectl \
networking/ctl/bin/calicoctl=/usr/bin/calicoctl \
networking/cni/bin/=/opt/cni/ \
images/=/opt/
# worker
fpm -s dir -n "kismatic-kubernetes-node" \
-p out/DEBs \
-C source/ \
-d 'kismatic-docker-engine = 1.11.2-0~xenial' \
-d 'bridge-utils' \
-v $K8S_VERSION  \
-a amd64 \
-t deb \
--license "Apache Software License 2.0" \
--maintainer "Apprenda <info@apprenda.com>" \
--vendor "Apprenda" \
--description "Kubernetes node binaries" \
--url "https://apprenda.com/" \
kubernetes/kubelet/bin/kubelet=/usr/bin/kubelet \
kubernetes/proxy/bin/kube-proxy=/usr/bin/kube-proxy \
networking/ctl/bin/calicoctl=/usr/bin/calicoctl \
networking/cni/bin/=/opt/cni/

# build etcd
#deb
fpm -s dir -n "kismatic-etcd" \
-p out/DEBs \
-C source/ \
-v $K8S_VERSION  \
-a amd64 \
-t deb \
--license "Apache Software License 2.0" \
--maintainer "Apprenda <info@apprenda.com>" \
--vendor "Apprenda" \
--description "Etcd kubernetes and networking binaries" \
--url "https://apprenda.com/" \
etcd/k8s/bin/etcdk8s=/usr/bin/etcdk8s \
etcd/k8s/bin/etcdctlk8s=/usr/bin/etcdctlk8s \
etcd/networking/bin/etcdnetworking=/usr/bin/etcdnetworking \
etcd/networking/bin/etcdctlnetworking=/usr/bin/etcdctlnetworking

# build docker
#deb
fpm -s deb -n "kismatic-docker-engine" \
-p out/DEBs/ \
-v $DOCKER_VERSION \
-a amd64 \
-t deb \
--license "Apache Software License 2.0" \
--maintainer "Apprenda <info@apprenda.com>" \
--vendor "Apprenda" \
--description "Docker and its dependencies" \
--url "https://apprenda.com/" \
source/docker/deb/docker-engine_1.11.2-0~xenial_amd64.deb
