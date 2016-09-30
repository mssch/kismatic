#!/bin/bash
K8S_VERSION=1.4.0
DOCKER_VERSION=1.11.2

# build Kubernetes
rm -rf out/kubernetes/
mkdir -p out/kubernetes/rpm/
mkdir -p out/kubernetes/deb/
# RPMs
# master
fpm -s dir -n "kubernetes-master" \
-p out/kubernetes/rpm/ \
-C source/ \
-d 'docker-engine = 1.11.2' \
-d 'bridge-utils' \
-v $K8S_VERSION  \
-a x86_64 \
-t rpm --rpm-os linux \
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
fpm -s dir -n "kubernetes-node" \
-p out/kubernetes/rpm/ \
-C source/ \
-d 'docker-engine = 1.11.2' \
-d 'bridge-utils' \
-v $K8S_VERSION  \
-a x86_64 \
-t rpm --rpm-os linux \
--license "Apache Software License 2.0" \
--maintainer "Apprenda <info@apprenda.com>" \
--vendor "Apprenda" \
--description "Kubernetes node binaries" \
--url "https://apprenda.com/" \
kubernetes/kubelet/bin/kubelet=/usr/bin/kubelet \
kubernetes/proxy/bin/kube-proxy=/usr/bin/kube-proxy \
networking/ctl/bin/calicoctl=/usr/bin/calicoctl \
networking/cni/bin/=/opt/cni/

#deb
# master
fpm -s dir -n "kubernetes-master" \
-p out/kubernetes/deb/ \
-C source/ \
-d 'docker-engine = 1.11.2-0~xenial' \
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
fpm -s dir -n "kubernetes-node" \
-p out/kubernetes/deb/ \
-C source/ \
-d 'docker-engine = 1.11.2-0~xenial' \
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
rm -rf out/etcd/
mkdir -p out/etcd/rpm/
mkdir -p out/etcd/deb/
# RPMs
fpm -s dir -n "etcd" \
-p out/etcd/rpm/ \
-C source/ \
-v $K8S_VERSION \
-a x86_64 \
-t rpm --rpm-os linux \
--license "Apache Software License 2.0" \
--maintainer "Apprenda <info@apprenda.com>" \
--vendor "Apprenda" \
--description "Etcd kubernetes and networking binaries" \
--url "https://apprenda.com/" \
etcd/k8s/bin/etcdk8s=/usr/bin/etcdk8s \
etcd/k8s/bin/etcdctlk8s=/usr/bin/etcdctlk8s \
etcd/networking/bin/etcdnetworking=/usr/bin/etcdnetworking \
etcd/networking/bin/etcdctlnetworking=/usr/bin/etcdctlnetworking

#deb
fpm -s dir -n "etcd" \
-p out/etcd/deb/ \
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
#RPMS
rm -rf out/docker/
mkdir -p out/docker/rpm/
fpm -s rpm -n "docker-engine" \
-p out/docker/rpm/ \
-v $DOCKER_VERSION \
-a x86_64 \
-t rpm --rpm-os linux \
--license "Apache Software License 2.0" \
--maintainer "Apprenda <info@apprenda.com>" \
--vendor "Apprenda" \
--description "Docker and its dependencies" \
--url "https://apprenda.com/" \
source/docker/rpm/docker-engine-selinux-1.11.2-1.el7.centos.noarch.rpm \
source/docker/rpm/docker-engine-1.11.2-1.el7.centos.x86_64.rpm
#deb
mkdir -p out/docker/deb/
cp source/docker/deb/* out/docker/deb/
