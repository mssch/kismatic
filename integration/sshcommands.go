package integration

const CopyKismaticYumRepo = `sudo curl https://s3.amazonaws.com/kismatic-rpm/kismatic.repo -o /etc/yum.repos.d/kismatic.repo`

const InstallEtcdYum = `sudo yum -y install kismatic-etcd`
const InstallDockerEngineYum = `sudo yum -y install kismatic-docker-engine`
const InstallKubernetesMasterYum = `sudo yum -y install kismatic-kubernetes-master`
const InstallKubernetesYum = `sudo yum -y install kismatic-kubernetes-node`

const CopyKismaticKeyDeb = `wget -qO - https://kismatic-deb.s3.amazonaws.com/public.key | sudo apt-key add - `
const CopyKismaticRepoDeb = `sudo add-apt-repository "deb https://kismatic-deb.s3.amazonaws.com/ xenial main"`
const UpdateAptGet = `sudo apt-get update`

const InstallEtcdApt = `sudo apt-get -y install kismatic-etcd`
const InstallDockerApt = `sudo apt-get -y install kismatic-docker-engine`
const InstallKubernetesMasterApt = `sudo apt-get -y install kismatic-kubernetes-master`
const InstallKubernetesApt = `sudo apt-get -y install kismatic-kubernetes-node`
