package install

// NetworkConfig describes the cluster's networking configuration
type NetworkConfig struct {
	Type             string
	PodCIDRBlock     string `yaml:"pod_cidr_block"`
	ServiceCIDRBlock string `yaml:"service_cidr_block"`
}

// CertsConfig describes the cluster's trust and certificate configuration
type CertsConfig struct {
	Expiry          string
	LocationCity    string `yaml:"location_city"`
	LocationState   string `yaml:"location_state"`
	LocationCountry string `yaml:"location_country"`
}

// SSHConfig describes the cluster's SSH configuration for accessing nodes
type SSHConfig struct {
	User string
	Key  string `yaml:"ssh_key"`
	Port int    `yaml:"ssh_port"`
}

// Cluster describes a Kismatic cluster
type Cluster struct {
	Name            string
	AdminPassword   string `yaml:"admin_password"`
	LocalRepository string `yaml:"local_repository"`
	HostsFileDNS    bool   `yaml:"hosts_file_dns"`
	Networking      NetworkConfig
	Certificates    CertsConfig
	SSH             SSHConfig
}

// A Node is a compute unit, virtual or physical, that is part of the cluster
type Node struct {
	Host       string
	IP         string
	InternalIP string
	Labels     []string
}

// A NodeGroup is a collection of nodes
type NodeGroup struct {
	ExpectedCount int `yaml:"expected_count"`
	Nodes         []Node
}

// MasterNodeGroup is the collection of master nodes
type MasterNodeGroup struct {
	NodeGroup
	LoadBalancedFQDN      string `yaml:"load_balanced_fqdn"`
	LoadBalancedShortName string `yaml:"load_balanced_short_name"`
}

// Plan is the installation plan that the user intends to execute
type Plan struct {
	Cluster Cluster
	Etcd    NodeGroup
	Master  MasterNodeGroup
	Worker  NodeGroup
}
