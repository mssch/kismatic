package integration

type PlanUbuntuAWS struct {
	Etcd                []AWSNodeDeets
	Master              []AWSNodeDeets
	Worker              []AWSNodeDeets
	MasterNodeFQDN      string
	MasterNodeShortName string
}

type AWSNodeDeets struct {
	Instanceid string
	Publicip   string
	Privateip  string
	Hostname   string
}

const planUbuntuAWSOverlay = `cluster:
  name: kubernetes
  admin_password: "abbazabba"
  networking:
    type: overlay
    pod_cidr_block: 172.16.0.0/16
    service_cidr_block: 172.17.0.0/16
  certificates:
    expiry: 17520h
    location_city: Troy
    location_state: New York
    location_country: US
  ssh:
    user: ubuntu
    ssh_key: /Users/mmiller/.ssh/kismatic-integration-testing.pem
    ssh_port: 22
etcd:
  expected_count: {{len .Etcd}} 
  nodes:
  {{range .Etcd}}- host: {{.Hostname}}
    ip: {{.Publicip}}
    internalip: {{.Privateip}}{{end}}
master:
  nodegroup:
    expected_count: {{len .Master}}
    nodes:
    {{range .Master}}- host: {{.Hostname}}
      ip: {{.Publicip}}
      internalip: {{.Privateip}}{{end}}
  load_balanced_fqdn: {{.MasterNodeFQDN}}
  load_balanced_short_name: {{.MasterNodeShortName}}
worker:
  expected_count: {{len .Worker}}
  nodes:
  {{range .Worker}}- host: {{.Hostname}}
    ip: {{.Publicip}}
    internalip: {{.Privateip}}{{end}}
`
