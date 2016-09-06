package install

// AnsibleVars holds all "--extra-vars"
type AnsibleVars struct {
	ClusterName                   string `json:"kubernetes_cluster_name"`
	AdminPassword                 string `json:"kubernetes_admin_password"`
	TLSDirectory                  string `json:"tls_directory"`
	KubernetesServicesCIDR        string `json:"kubernetes_services_cidr"`
	KubernetesPodsCIDR            string `json:"kubernetes_pods_cidr"`
	KubernetesDNSServiceIP        string `json:"kubernetes_dns_service_ip"`
	CalicoNetworkType             string `json:"calico_network_type"`
	LocalRepository               string `json:"local_repoository_path,omitempty"` //Optional
	ForceEtcdRestart              string `json:"force_etcd_restart"`
	ForceApiserverRestart         string `json:"force_apiserver_restart"`
	ForceControllerManagerRestart string `json:"force_controller_manager_restart"`
	ForceSchedulerRestart         string `json:"force_scheduler_restart"`
	ForceProxyRestart             string `json:"force_proxy_restart"`
	ForceKubeletRestart           string `json:"force_kubelet_restart"`
	ForceCalicoRestart            string `json:"force_calico_node_restart"`
	ForceDockerRestart            string `json:"force_docker_restart"`
}
