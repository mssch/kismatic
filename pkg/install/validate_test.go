package install

import (
	"fmt"
	"testing"
)

func validPlan() Plan {
	return Plan{
		Cluster: Cluster{
			Name:    "test",
			Version: "v1.10.11",
			Networking: NetworkConfig{
				Type:             "overlay",
				PodCIDRBlock:     "172.16.0.0/16",
				ServiceCIDRBlock: "172.20.0.0/16",
			},
			Certificates: CertsConfig{
				Expiry: "17250h",
			},
			SSH: SSHConfig{
				User: "root",
				Key:  "/bin/sh",
				Port: 22,
			},
		},
		AdditionalFiles: []AdditionalFile{
			{
				Source:      "/bin/sh",
				Destination: "/bin/sh",
				Hosts:       []string{"master01"},
			},
		},
		AddOns: AddOns{
			CNI: &CNI{
				Provider: "calico",
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode:     "overlay",
						LogLevel: "info",
					},
				},
			},
			DNS: DNS{
				Provider: "kubedns",
			},
			Dashboard: Dashboard{
				Options: DashboardOptions{
					ServiceType: "ClusterIP",
				},
			},
			HeapsterMonitoring: &HeapsterMonitoring{
				Options: HeapsterOptions{
					Heapster: Heapster{
						Replicas:    2,
						ServiceType: "ClusterIP",
					},
				},
			},
		},
		Etcd: NodeGroup{
			ExpectedCount: 1,
			Nodes: []Node{
				{
					Host: "etcd01",
					IP:   "192.168.205.10",
				},
			},
		},
		Master: MasterNodeGroup{
			ExpectedCount: 1,
			Nodes: []Node{
				{
					Host: "master01",
					IP:   "192.168.205.11",
				},
			},
			LoadBalancer: "test:6443",
		},
		Worker: NodeGroup{
			ExpectedCount: 1,
			Nodes: []Node{
				{
					Host: "worker01",
					IP:   "192.168.205.12",
				},
			},
		},
		Ingress: OptionalNodeGroup{
			ExpectedCount: 1,
			Nodes: []Node{
				{
					Host: "etcd01",
					IP:   "192.168.205.10",
				},
			},
		},
		NFS: &NFS{
			Volumes: []NFSVolume{
				{
					Host: "10.10.2.20",
					Path: "/",
				},
			},
		},
	}
}

func assertInvalidPlan(t *testing.T, p Plan) {
	valid, _ := ValidatePlan(&p)
	if valid {
		t.Errorf("expected invalid, but got valid")
	}
}

func TestValidateBlankPlan(t *testing.T) {
	p := Plan{}
	assertInvalidPlan(t, p)
}

func TestValidateValidPlan(t *testing.T) {
	p := validPlan()
	valid, errs := ValidatePlan(&p)
	if !valid {
		t.Errorf("expected valid, but got invalid")
	}
	fmt.Println(errs)
}

func TestClusterVersion(t *testing.T) {
	tests := []struct {
		c     Cluster
		valid bool
	}{
		{c: Cluster{
			Name:    "test",
			Version: "v1.10.11",
			Networking: NetworkConfig{
				Type:             "overlay",
				PodCIDRBlock:     "172.16.0.0/16",
				ServiceCIDRBlock: "172.20.0.0/16",
			},
			Certificates: CertsConfig{
				Expiry: "17250h",
			},
			SSH: SSHConfig{
				User: "root",
				Key:  "/bin/sh",
				Port: 22,
			},
		},
			valid: true,
		},
		{c: Cluster{
			Name:    "test",
			Version: "v1.10.11",
			Networking: NetworkConfig{
				Type:             "overlay",
				PodCIDRBlock:     "172.16.0.0/16",
				ServiceCIDRBlock: "172.20.0.0/16",
			},
			Certificates: CertsConfig{
				Expiry: "17250h",
			},
			SSH: SSHConfig{
				User: "root",
				Key:  "/bin/sh",
				Port: 22,
			},
		},
			valid: true,
		},
		{c: Cluster{
			Name:    "test",
			Version: "foo",
			Networking: NetworkConfig{
				Type:             "overlay",
				PodCIDRBlock:     "172.16.0.0/16",
				ServiceCIDRBlock: "172.20.0.0/16",
			},
			Certificates: CertsConfig{
				Expiry: "17250h",
			},
			SSH: SSHConfig{
				User: "root",
				Key:  "/bin/sh",
				Port: 22,
			},
		},
			valid: false,
		},
		{c: Cluster{
			Name:    "test",
			Version: "v1.10.110",
			Networking: NetworkConfig{
				Type:             "overlay",
				PodCIDRBlock:     "172.16.0.0/16",
				ServiceCIDRBlock: "172.20.0.0/16",
			},
			Certificates: CertsConfig{
				Expiry: "17250h",
			},
			SSH: SSHConfig{
				User: "root",
				Key:  "/bin/sh",
				Port: 22,
			},
		},
			valid: false,
		},
		{c: Cluster{
			Name:    "test",
			Version: "v1.8.0",
			Networking: NetworkConfig{
				Type:             "overlay",
				PodCIDRBlock:     "172.16.0.0/16",
				ServiceCIDRBlock: "172.20.0.0/16",
			},
			Certificates: CertsConfig{
				Expiry: "17250h",
			},
			SSH: SSHConfig{
				User: "root",
				Key:  "/bin/sh",
				Port: 22,
			},
		},
			valid: false,
		},
		{c: Cluster{
			Name:    "test",
			Version: "v1.20.0",
			Networking: NetworkConfig{
				Type:             "overlay",
				PodCIDRBlock:     "172.16.0.0/16",
				ServiceCIDRBlock: "172.20.0.0/16",
			},
			Certificates: CertsConfig{
				Expiry: "17250h",
			},
			SSH: SSHConfig{
				User: "root",
				Key:  "/bin/sh",
				Port: 22,
			},
		},
			valid: false,
		},
		{c: Cluster{
			Name:                     "test",
			Version:                  "v1.8.0",
			DisconnectedInstallation: true,
			Networking: NetworkConfig{
				Type:             "overlay",
				PodCIDRBlock:     "172.16.0.0/16",
				ServiceCIDRBlock: "172.20.0.0/16",
			},
			Certificates: CertsConfig{
				Expiry: "17250h",
			},
			SSH: SSHConfig{
				User: "root",
				Key:  "/bin/sh",
				Port: 22,
			},
		},
			valid: false,
		},
		{c: Cluster{
			Name:                     "test",
			Version:                  "v1.20.0",
			DisconnectedInstallation: true,
			Networking: NetworkConfig{
				Type:             "overlay",
				PodCIDRBlock:     "172.16.0.0/16",
				ServiceCIDRBlock: "172.20.0.0/16",
			},
			Certificates: CertsConfig{
				Expiry: "17250h",
			},
			SSH: SSHConfig{
				User: "root",
				Key:  "/bin/sh",
				Port: 22,
			},
		},
			valid: false,
		},
	}
	for n, test := range tests {
		if valid, _ := test.c.validate(); valid != test.valid {
			t.Errorf("%d: expected %t with %+v, but got %t", n, test.valid, test.c, valid)
		}
	}
}

func TestValidatePlanEmptyPodCIDR(t *testing.T) {
	p := validPlan()
	p.Cluster.Networking.PodCIDRBlock = ""
	assertInvalidPlan(t, p)
}

func TestValidatePlanInvalidPodCIDR(t *testing.T) {
	p := validPlan()
	p.Cluster.Networking.PodCIDRBlock = "foo"
	assertInvalidPlan(t, p)
}

func TestValidatePlanEmptyServicesCIDR(t *testing.T) {
	p := validPlan()
	p.Cluster.Networking.ServiceCIDRBlock = ""
	assertInvalidPlan(t, p)
}

func TestValidatePlanInvalidServicesCIDR(t *testing.T) {
	p := validPlan()
	p.Cluster.Networking.ServiceCIDRBlock = "foo"
	assertInvalidPlan(t, p)
}

func TestValidatePlanEmptyCertificatesExpiry(t *testing.T) {
	p := validPlan()
	p.Cluster.Certificates.Expiry = ""
	assertInvalidPlan(t, p)
}

func TestValidatePlanInvalidCertExpiry(t *testing.T) {
	p := validPlan()
	p.Cluster.Certificates.Expiry = "foo"
	assertInvalidPlan(t, p)
}

func TestValidatePlanEmptyCACertExpiryIsValid(t *testing.T) {
	p := validPlan()
	p.Cluster.Certificates.CAExpiry = ""
	valid, _ := p.validate()
	if !valid {
		t.Errorf("plan was found invalid")
	}
}

func TestValidatePlanInvalidCACertificatesExpiry(t *testing.T) {
	p := validPlan()
	p.Cluster.Certificates.CAExpiry = "foo"
	assertInvalidPlan(t, p)
}

func TestValidatePlanEmptySSHUser(t *testing.T) {
	p := validPlan()
	p.Cluster.SSH.User = ""
	assertInvalidPlan(t, p)
}

func TestValidatePlanEmptySSHKey(t *testing.T) {
	p := validPlan()
	p.Cluster.SSH.Key = ""
	assertInvalidPlan(t, p)
}

func TestValidatePlanNonExistentSSHKey(t *testing.T) {
	p := validPlan()
	p.Cluster.SSH.Key = "/foo"
	assertInvalidPlan(t, p)
}

func TestValidatePlanNegativeSSHPort(t *testing.T) {
	p := validPlan()
	p.Cluster.SSH.Port = -1
	assertInvalidPlan(t, p)
}

func TestValidatePlanEmptyLoadBalancer(t *testing.T) {
	p := validPlan()
	p.Master.LoadBalancer = ""
	assertInvalidPlan(t, p)
}

func TestValidatePlanNoEtcdNodes(t *testing.T) {
	p := validPlan()
	p.Etcd.ExpectedCount = 0
	p.Etcd.Nodes = []Node{}
	assertInvalidPlan(t, p)
}

func TestValidatePlanNoMasterNodes(t *testing.T) {
	p := validPlan()
	p.Master.ExpectedCount = 0
	p.Master.Nodes = []Node{}
	assertInvalidPlan(t, p)
}

func TestValidatePlanNoWorkerNodes(t *testing.T) {
	p := validPlan()
	p.Worker.ExpectedCount = 0
	p.Worker.Nodes = []Node{}
	assertInvalidPlan(t, p)
}

func TestValidatePlanEtcdNodesMismatch(t *testing.T) {
	p := validPlan()
	p.Etcd.ExpectedCount = 100
	assertInvalidPlan(t, p)
}

func TestValidatePlanMasterNodesMismatch(t *testing.T) {
	p := validPlan()
	p.Master.ExpectedCount = 100
	assertInvalidPlan(t, p)
}

func TestValidatePlanWorkerNodesMismatch(t *testing.T) {
	p := validPlan()
	p.Worker.ExpectedCount = 100
	assertInvalidPlan(t, p)
}

func TestValidatePlanUnexpectedEtcdNodes(t *testing.T) {
	p := validPlan()
	p.Etcd.ExpectedCount = 1
	p.Etcd.Nodes = []Node{
		{
			Host: "etcd01",
			IP:   "192.168.205.10",
		},
		{
			Host: "etcd02",
			IP:   "192.168.205.11",
		},
	}
	assertInvalidPlan(t, p)
}

func TestValidatePlanUnexpectedMasterNodes(t *testing.T) {
	p := validPlan()
	p.Master.ExpectedCount = 1
	p.Master.Nodes = []Node{
		{
			Host: "master01",
			IP:   "192.168.205.10",
		},
		{
			Host: "master02",
			IP:   "192.168.205.11",
		},
	}
	assertInvalidPlan(t, p)
}

func TestValidatePlanUnexpectedWorkerNodes(t *testing.T) {
	p := validPlan()
	p.Worker.ExpectedCount = 1
	p.Worker.Nodes = []Node{
		{
			Host: "worker01",
			IP:   "192.168.205.10",
		},
		{
			Host: "worker02",
			IP:   "192.168.205.11",
		},
	}
	assertInvalidPlan(t, p)
}

func TestValidatePlanNoIngress(t *testing.T) {
	p := validPlan()
	p.Ingress.ExpectedCount = 0
	p.Ingress.Nodes = []Node{}
	valid, _ := ValidatePlan(&p)
	if !valid {
		t.Errorf("expected valid, but got invalid")
	}
}

func TestValidatePlanIngressExpected(t *testing.T) {
	p := validPlan()
	p.Ingress.ExpectedCount = 1
	p.Ingress.Nodes = []Node{}
	assertInvalidPlan(t, p)
}

func TestValidatePlanIngressProvidedNotExpected(t *testing.T) {
	p := validPlan()
	p.Ingress.ExpectedCount = 0
	p.Ingress.Nodes = []Node{
		{
			Host: "ingress",
			IP:   "192.168.205.10",
		},
	}
	assertInvalidPlan(t, p)
}

func TestValidateStorageVolume(t *testing.T) {
	tests := []struct {
		sv    StorageVolume
		valid bool
	}{
		{
			sv: StorageVolume{
				Name:              "foo",
				SizeGB:            100,
				DistributionCount: 2,
				ReplicateCount:    2,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: true,
		},
		{
			sv: StorageVolume{
				Name:              "foo",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: true,
		},
		{
			sv: StorageVolume{
				Name:              "foo",
				SizeGB:            100,
				DistributionCount: 0,
				ReplicateCount:    1,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "foo",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    0,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "bad name with spaces",
				SizeGB:            100,
				DistributionCount: 2,
				ReplicateCount:    2,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "bad:name2",
				SizeGB:            100,
				DistributionCount: 2,
				ReplicateCount:    2,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            0,
				DistributionCount: 2,
				ReplicateCount:    2,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            -1,
				DistributionCount: 2,
				ReplicateCount:    2,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: -1,
				ReplicateCount:    2,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 2,
				ReplicateCount:    -1,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv:    StorageVolume{},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "foo",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: true,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "Recycle",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: true,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "Delete",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: true,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "Delete",
				AccessModes:       []string{"ReadWriteOnce"},
			},
			valid: true,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "Delete",
				AccessModes:       []string{"ReadOnlyMany"},
			},
			valid: true,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "Delete",
				AccessModes:       []string{"ReadWriteMany", "ReadWriteOnce", "ReadOnlyMany"},
			},
			valid: true,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "Delete",
				AccessModes:       []string{"someBadAccessMode"},
			},
			valid: false,
		},
	}
	for _, test := range tests {
		if valid, _ := test.sv.validate(); valid != test.valid {
			t.Errorf("expected %v with %+v, but got %v", test.valid, test.sv, !test.valid)
		}
	}
}

func TestValidateAllowAddress(t *testing.T) {
	tests := []struct {
		address string
		valid   bool
	}{
		{"192.168.205.10", true},
		{"192.168.205.*", true},
		{"192.168.*.*", true},
		{"192.*.*.*", true},
		{"*.168.205.10", true},
		{"*.*.205.10", true},
		{"*.*.*.10", true},
		{"*.*.*.*", true},
		{"-1.-1.-1.-1", false},
		{"-1.*.*.*", false},
		{"*.-1.*.*", false},
		{"*.*.-1.*", false},
		{"*.*.*.-1", false},
		{"256.256.256.256", false},
		{"256.*.*.*", false},
		{"*.256.*.*", false},
		{"*.*.256.*", false},
		{"*.*.*.256", false},
		{"a.a.a.a", false},
		{"*.*.*.a", false},
		{"*.*.a.*", false},
		{"*.a.*.*", false},
		{"a.*.*.*", false},
		{"", false},
		{"foo", false},
		{"192", false},
		{"192.168", false},
		{"192.168.205", false},
		{"...", false},
		{"192...", false},
		{"192.168..", false},
		{"192.168.205.", false},
	}
	for _, test := range tests {
		if validateAllowedAddress(test.address) != test.valid {
			t.Errorf("expected %v with address %q, but got %v", test.valid, test.address, !test.valid)
		}
	}
}

func TestValidatePlanNFSDupes(t *testing.T) {
	p := validPlan()

	p.NFS.Volumes = append(p.NFS.Volumes, NFSVolume{
		Host: "10.10.2.20",
		Path: "/",
	})

	assertInvalidPlan(t, p)
}

func TestValidateNFSVolume(t *testing.T) {
	tests := []struct {
		host  string
		path  string
		valid bool
	}{
		{
			host:  "10.10.2.10",
			path:  "/foo",
			valid: true,
		},
		{
			host:  "10.10.2.10",
			path:  "",
			valid: false,
		},
		{
			host:  "10.10.2.10",
			path:  "../someRelativePath",
			valid: false,
		},
		{
			host:  "",
			path:  "/foo",
			valid: false,
		},
	}
	for _, test := range tests {
		v := NFSVolume{
			Host: test.host,
			Path: test.path,
		}
		if valid, _ := v.validate(); valid != test.valid {
			t.Errorf("Expected valid = %v, but got %v", test.valid, valid)
		}
	}
}

func TestValidatePlanCerts(t *testing.T) {
	p := validPlan()

	pki := getPKI(t)
	defer cleanup(pki.GeneratedCertsDirectory, t)
	ca, err := pki.GenerateClusterCA(&p)
	if err != nil {
		t.Fatalf("error generating CA for test: %v", err)
	}
	proxyClientCA, err := pki.GenerateProxyClientCA(&p)
	if err != nil {
		t.Fatalf("error generating proxy-client CA for test: %v", err)
	}
	if err := pki.GenerateClusterCertificates(&p, ca, proxyClientCA); err != nil {
		t.Fatalf("failed to generate certs: %v", err)
	}

	valid, errs := ValidateCertificates(&p, &pki)
	if !valid {
		t.Errorf("expected valid, but got invalid")
		fmt.Println(errs)
	}
}

func TestValidatePlanBadCerts(t *testing.T) {
	p := validPlan()

	pki := getPKI(t)
	defer cleanup(pki.GeneratedCertsDirectory, t)

	ca, err := pki.GenerateClusterCA(&p)
	if err != nil {
		t.Fatalf("error generating CA for test: %v", err)
	}
	proxyClientCA, err := pki.GenerateProxyClientCA(&p)
	if err != nil {
		t.Fatalf("error generating proxy-client CA for test: %v", err)
	}
	if err := pki.GenerateClusterCertificates(&p, ca, proxyClientCA); err != nil {
		t.Fatalf("failed to generate certs: %v", err)
	}
	p.Master.Nodes[0] = Node{
		Host:       "master01",
		IP:         "11.12.13.14",
		InternalIP: "22.33.44.55",
	}

	valid, _ := ValidateCertificates(&p, &pki)
	if valid {
		t.Errorf("expected an error, but got valid")
	}
}

func TestValidatePlanMissingCerts(t *testing.T) {
	p := validPlan()

	pki := getPKI(t)
	defer cleanup(pki.GeneratedCertsDirectory, t)

	valid, errs := ValidateCertificates(&p, &pki)
	if !valid {
		t.Errorf("expected valid, but got invalid")
		fmt.Println(errs)
	}
}

func TestValidatePlanMissingSomeCerts(t *testing.T) {
	p := validPlan()

	pki := getPKI(t)
	defer cleanup(pki.GeneratedCertsDirectory, t)

	ca, err := pki.GenerateClusterCA(&p)
	if err != nil {
		t.Fatalf("error generating CA for test: %v", err)
	}
	proxyClientCA, err := pki.GenerateProxyClientCA(&p)
	if err != nil {
		t.Fatalf("error generating proxy-client CA for test: %v", err)
	}
	if err := pki.GenerateClusterCertificates(&p, ca, proxyClientCA); err != nil {
		t.Fatalf("failed to generate certs: %v", err)
	}

	newNode := Node{
		Host:       "master2",
		IP:         "11.12.13.14",
		InternalIP: "22.33.44.55",
	}
	p.Master.Nodes = append(p.Master.Nodes, newNode)

	valid, errs := ValidateCertificates(&p, &pki)
	if !valid {
		t.Errorf("expected valid, but got invalid")
		fmt.Println(errs)
	}
}

func TestValidateNodeListDuplicate(t *testing.T) {
	tests := []struct {
		nl    nodeList
		valid bool
	}{
		{
			nl: nodeList{
				[]Node{
					{
						Host: "host1",
						IP:   "10.0.0.1",
					},
				},
			},
			valid: true,
		},
		{
			nl: nodeList{
				[]Node{
					{
						Host: "host1",
						IP:   "10.0.0.1",
					},
					{
						Host: "host1",
						IP:   "10.0.0.1",
					},
				},
			},
			valid: true,
		},
		{
			nl: nodeList{
				[]Node{
					{
						Host: "host1",
						IP:   "10.0.0.1",
					},
					{
						Host: "host2",
						IP:   "10.0.0.2",
					},
				},
			},
			valid: true,
		},
		{
			nl: nodeList{
				[]Node{
					{
						Host: "host1",
						IP:   "10.0.0.1",
					},
					{
						Host: "host1",
						IP:   "10.0.0.2",
					},
				},
			},
			valid: false,
		},
		{
			nl: nodeList{
				[]Node{
					{
						Host: "host1",
						IP:   "10.0.0.2",
					},
					{
						Host: "host2",
						IP:   "10.0.0.2",
					},
				},
			},
			valid: false,
		},
		{
			nl: nodeList{
				[]Node{
					{
						Host:       "host1",
						IP:         "10.0.0.1",
						InternalIP: "192.168.205.10",
					},
					{
						Host:       "host2",
						IP:         "10.0.0.2",
						InternalIP: "192.168.205.10",
					},
				},
			},
			valid: false,
		},
	}
	for i, test := range tests {
		ok, _ := test.nl.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}

func TestValidatePlanDisconnectedInstallationFailsDueToMissingRegistry(t *testing.T) {
	plan := validPlan()
	plan.Cluster.DisconnectedInstallation = true
	ok, errs := plan.validate()
	if ok {
		t.Errorf("expected validation to fail due to missing external registry information in plan")
	}
	var found bool
	for _, err := range errs {
		if err.Error() == "A container image registry is required when disconnected_installation is true" {
			found = true
		}
	}
	if !found {
		t.Errorf("validation did not return the expected failure message")
	}
}

func TestValidatePlanDisconnectedInstallationSucceeds(t *testing.T) {
	plan := validPlan()
	plan.Cluster.DisconnectedInstallation = true
	plan.DockerRegistry.Server = "localhost:5000"
	if ok, errs := plan.validate(); !ok {
		t.Error("expected validation to succeed, but it failed")
		t.Logf("errors were: %v\n", errs)
	}
}

func TestDockerRegistry(t *testing.T) {
	tests := []struct {
		d     DockerRegistry
		valid bool
	}{
		{
			d:     DockerRegistry{},
			valid: true,
		},
		{
			d: DockerRegistry{
				Server: "172.0.0.1",
				CAPath: "/bin/sh",
			},
			valid: true,
		},
		{
			d: DockerRegistry{
				Server:   "172.0.0.1",
				Username: "user",
				Password: "password",
			},
			valid: true,
		},
		{
			d: DockerRegistry{
				Address: "172.0.0.1",
				CAPath:  "/bin/sh",
			},
			valid: true,
		},
		{
			d: DockerRegistry{
				Address:  "172.0.0.1",
				Username: "user",
				Password: "password",
			},
			valid: true,
		},
		{
			d: DockerRegistry{
				CAPath: "user",
			},
			valid: false,
		},
		{
			d: DockerRegistry{
				Username: "user",
			},
			valid: false,
		},
		{
			d: DockerRegistry{
				Password: "password",
			},
			valid: false,
		},
	}
	for i, test := range tests {
		ok, _ := test.d.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}

func TestValidateDockerStorage(t *testing.T) {
	tests := []struct {
		storage DockerStorage
		valid   bool
	}{
		{
			storage: DockerStorage{
				Driver: "devicemapper",
			},
			valid: true,
		},
		{
			storage: DockerStorage{
				Driver: "overlay2",
			},
			valid: true,
		},
		{
			storage: DockerStorage{
				Driver: "foo",
			},
			valid: true,
		},
		{
			storage: DockerStorage{
				Driver: "devicemapper",
				DirectLVMBlockDevice: DirectLVMBlockDevice{
					Path: "",
				},
			},
			valid: true,
		},
		{
			storage: DockerStorage{
				Driver: "devicemapper",
				DirectLVMBlockDevice: DirectLVMBlockDevice{
					Path: "foo",
				},
			},
			valid: false,
		},
		{
			storage: DockerStorage{
				Driver: "devicemapper",
				DirectLVMBlockDevice: DirectLVMBlockDevice{
					Path: "/foo/bar",
				},
			},
			valid: true,
		},
		{
			storage: DockerStorage{
				Driver: "foo",
				DirectLVMBlockDevice: DirectLVMBlockDevice{
					Path: "/foo/bar",
				},
			},
			valid: false,
		},
	}
	for i, test := range tests {
		ok, _ := test.storage.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect valid, but got invalid", i)
		}
	}
}

func TestValidateDockerStorageDirectLVM(t *testing.T) {
	tests := []struct {
		config DockerStorageDirectLVMDeprecated
		valid  bool
	}{
		{
			config: DockerStorageDirectLVMDeprecated{
				Enabled: false,
			},
			valid: true,
		},
		{
			config: DockerStorageDirectLVMDeprecated{
				Enabled: true,
			},
			valid: false,
		},
		{
			config: DockerStorageDirectLVMDeprecated{
				Enabled:     true,
				BlockDevice: "foo",
			},
			valid: false,
		},
		{
			config: DockerStorageDirectLVMDeprecated{
				Enabled:     true,
				BlockDevice: "/dev/sdb",
			},
			valid: true,
		},
	}
	for i, test := range tests {
		ok, _ := test.config.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect valid, but got invalid", i)
		}
	}
}

func TestCNIAddOn(t *testing.T) {
	tests := []struct {
		n     CNI
		valid bool
	}{
		{
			n: CNI{
				Provider: "calico",
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode: "overlay",
					},
				},
			},
			valid: true,
		},
		{
			n: CNI{
				Provider: "calico",
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode: "routed",
					},
				},
			},
			valid: true,
		},
		{
			n: CNI{
				Provider: "weave",
			},
			valid: true,
		},
		{
			n: CNI{
				Provider: "contiv",
			},
			valid: true,
		},
		{
			n: CNI{
				Provider: "foo",
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode: "overlay",
					},
				},
			},
			valid: false,
		},
		{
			n: CNI{
				Provider: "calico",
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode: "foo",
					},
				},
			},
			valid: false,
		},
		{
			n: CNI{
				Provider: "foo",
				Disable:  true,
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode: "overlay",
					},
				},
			},
			valid: true,
		},
		{
			n: CNI{
				Provider: "calico",
				Disable:  true,
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode: "foo",
					},
				},
			},
			valid: true,
		},
		{
			n: CNI{
				Provider: "calico",
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode:     "overlay",
						LogLevel: "",
					},
				},
			},
			valid: true,
		},
		{
			n: CNI{
				Provider: "calico",
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode:     "overlay",
						LogLevel: "info",
					},
				},
			},
			valid: true,
		},
		{
			n: CNI{
				Provider: "calico",
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode:     "overlay",
						LogLevel: "warning",
					},
				},
			},
			valid: true,
		},
		{
			n: CNI{
				Provider: "calico",
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode:     "overlay",
						LogLevel: "debug",
					},
				},
			},
			valid: true,
		},
		{
			n: CNI{
				Provider: "calico",
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode:     "overlay",
						LogLevel: "INFO",
					},
				},
			},
			valid: false,
		},
		{
			n: CNI{
				Provider: "calico",
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode:     "overlay",
						LogLevel: "foo",
					},
				},
			},
			valid: false,
		},
	}
	for i, test := range tests {
		ok, _ := test.n.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}

func TestDNSProvider(t *testing.T) {
	tests := []struct {
		d     DNS
		valid bool
	}{
		{
			d: DNS{
				Provider: "kubedns",
			},
			valid: true,
		},
		{
			d: DNS{
				Provider: "coredns",
			},
			valid: true,
		},
		{
			d: DNS{
				Disable:  true,
				Provider: "foo",
			},
			valid: true,
		},
		{
			d: DNS{
				Provider: "foo",
			},
			valid: false,
		},
	}
	for i, test := range tests {
		ok, _ := test.d.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}

func TestHeapsterAddOn(t *testing.T) {
	tests := []struct {
		h     HeapsterMonitoring
		valid bool
	}{
		{
			h: HeapsterMonitoring{
				Options: HeapsterOptions{
					Heapster: Heapster{
						Replicas:    0,
						ServiceType: "ClusterIP",
					},
				},
			},
			valid: false,
		},
		{
			h: HeapsterMonitoring{
				Options: HeapsterOptions{
					Heapster: Heapster{
						Replicas:    1,
						ServiceType: "Foo",
					},
				},
			},
			valid: false,
		},
		{
			h: HeapsterMonitoring{
				Options: HeapsterOptions{
					Heapster: Heapster{
						Replicas:    -1,
						ServiceType: "ClusterIP",
					},
				},
			},
			valid: false,
		},
		{
			h: HeapsterMonitoring{
				Options: HeapsterOptions{
					Heapster: Heapster{
						Replicas:    1,
						ServiceType: "ClusterIP",
					},
				},
			},
			valid: true,
		},
	}
	for i, test := range tests {
		ok, _ := test.h.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}

func TestDashbordAddOn(t *testing.T) {
	tests := []struct {
		d     Dashboard
		valid bool
	}{
		{
			d: Dashboard{
				Options: DashboardOptions{
					ServiceType: "ClusterIP",
				},
			},
			valid: true,
		},
		{
			d: Dashboard{
				Options: DashboardOptions{
					ServiceType: "Foo",
				},
			},
			valid: false,
		},
		{
			d: Dashboard{
				Options: DashboardOptions{
					ServiceType: "ClusterIP",
					NodePort:    "32500",
				},
			},
			valid: false,
		},
		{
			d: Dashboard{
				Options: DashboardOptions{
					ServiceType: "NodePort",
					NodePort:    "32500",
				},
			},
			valid: true,
		},
	}
	for i, test := range tests {
		ok, _ := test.d.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}

func TestPackageManagerAddOn(t *testing.T) {
	tests := []struct {
		p     PackageManager
		valid bool
	}{
		{
			p: PackageManager{
				Disable:  false,
				Provider: "helm",
			},
			valid: true,
		},
		{
			p: PackageManager{
				Disable:  true,
				Provider: "",
			},
			valid: true,
		},
		{
			p: PackageManager{
				Disable:  true,
				Provider: "foo",
			},
			valid: true,
		},
		{
			p: PackageManager{
				Disable:  false,
				Provider: "",
			},
			valid: true,
		},
		{
			p: PackageManager{
				Disable:  false,
				Provider: "foo",
			},
			valid: false,
		},
	}
	for i, test := range tests {
		ok, _ := test.p.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}

func TestCloudProvider(t *testing.T) {
	tests := []struct {
		c     CloudProvider
		valid bool
	}{
		{
			c: CloudProvider{
				Provider: "",
			},
			valid: true,
		},
		{
			c: CloudProvider{
				Provider: "aws",
			},
			valid: true,
		},
		{
			c: CloudProvider{
				Provider: "awss",
			},
			valid: false,
		},
		{
			c: CloudProvider{
				Provider: "gce",
				Config:   "/bin/sh",
			},
			valid: true,
		},
		{
			c: CloudProvider{
				Provider: "gce",
				Config:   "/bin/foo",
			},
			valid: false,
		},
		{
			c: CloudProvider{
				Provider: "gce",
				Config:   "foo",
			},
			valid: false,
		},
	}
	for i, test := range tests {
		ok, _ := test.c.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}

func TestNodeLabels(t *testing.T) {
	tests := []struct {
		n     Node
		valid bool
	}{
		{
			n: Node{
				Host: "foo",
				IP:   "192.1.1.1",
			},
			valid: true,
		},
		{
			n: Node{
				Host:   "foo",
				IP:     "192.1.1.1",
				Labels: map[string]string{},
			},
			valid: true,
		},
		{
			n: Node{
				Host:   "foo",
				IP:     "192.1.1.1",
				Labels: map[string]string{"com.foo/bar": ""},
			},
			valid: true,
		},
		{
			n: Node{
				Host:   "foo",
				IP:     "192.1.1.1",
				Labels: map[string]string{"com.foo/bar": "foobar"},
			},
			valid: true,
		},
		{
			n: Node{
				Host:   "foo",
				IP:     "192.1.1.1",
				Labels: map[string]string{"com.foo/bar": "foobar", "com.foo/xyz": "xyz"},
			},
			valid: true,
		},
		{
			n: Node{
				Host:   "foo",
				IP:     "192.1.1.1",
				Labels: map[string]string{"kismatic/foo": "bar"},
			},
			valid: false,
		},
		{
			n: Node{
				Host:   "foo",
				IP:     "192.1.1.1",
				Labels: map[string]string{"com.foo/kismatic-version": "v1.0.0"},
			},
			valid: true,
		},
		{
			n: Node{
				Host:   "foo",
				IP:     "192.1.1.1",
				Labels: map[string]string{"": "com.foo/worker"},
			},
			valid: false,
		},
		{
			n: Node{
				Host:   "foo",
				IP:     "192.1.1.1",
				Labels: map[string]string{"": ""},
			},
			valid: false,
		},
		{
			n: Node{
				Host:   "foo",
				IP:     "192.1.1.1",
				Labels: map[string]string{"node-type:test": "test"},
			},
			valid: false,
		},
		{
			n: Node{
				Host:   "foo",
				IP:     "192.1.1.1",
				Labels: map[string]string{"com.foo/invalid": ":test"},
			},
			valid: false,
		},
		{
			n: Node{
				Host:   "foo",
				IP:     "192.1.1.1",
				Labels: map[string]string{"node-type:test": ":test"},
			},
			valid: false,
		},
	}
	for i, test := range tests {
		ok, _ := test.n.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}

func TestNodeTaints(t *testing.T) {
	tests := []struct {
		n     Node
		valid bool
	}{
		{
			n: Node{
				Host: "foo",
				IP:   "192.1.1.1",
			},
			valid: true,
		},
		{
			n: Node{
				Host:   "foo",
				IP:     "192.1.1.1",
				Taints: []Taint{},
			},
			valid: true,
		},
		{
			n: Node{
				Host: "foo",
				IP:   "192.1.1.1",
				Taints: []Taint{
					Taint{
						Key:    "com.foo/bar",
						Value:  "",
						Effect: "NoSchedule",
					},
				},
			},
			valid: true,
		},
		{
			n: Node{
				Host: "foo",
				IP:   "192.1.1.1",
				Taints: []Taint{
					Taint{
						Key:    "com.foo/bar",
						Value:  "foobar",
						Effect: "NoSchedule",
					},
				},
			},
			valid: true,
		},
		{
			n: Node{
				Host: "foo",
				IP:   "192.1.1.1",
				Taints: []Taint{
					Taint{
						Key:    "com.foo/bar",
						Value:  "foobar",
						Effect: "NoSchedule",
					},
					Taint{
						Key:    "com.foo/xyz",
						Value:  "xyz",
						Effect: "NoSchedule",
					},
				},
			},
			valid: true,
		},
		{
			n: Node{
				Host: "foo",
				IP:   "192.1.1.1",
				Taints: []Taint{
					Taint{
						Key:    "kismatic/foo",
						Value:  "bar",
						Effect: "NoSchedule",
					},
				},
			},
			valid: false,
		},
		{
			n: Node{
				Host: "foo",
				IP:   "192.1.1.1",
				Taints: []Taint{
					Taint{
						Key:    "com.foo/kismatic-version",
						Value:  "v1.0.0",
						Effect: "NoSchedule",
					},
				},
			},
			valid: true,
		},
		{
			n: Node{
				Host: "foo",
				IP:   "192.1.1.1",
				Taints: []Taint{
					Taint{
						Key:    "",
						Value:  "v1.0.0",
						Effect: "NoSchedule",
					},
				},
			},
			valid: false,
		},
		{
			n: Node{
				Host: "foo",
				IP:   "192.1.1.1",
				Taints: []Taint{
					Taint{
						Key:    "",
						Value:  "",
						Effect: "NoSchedule",
					},
				},
				Labels: map[string]string{"": ""},
			},
			valid: false,
		},
		{
			n: Node{
				Host: "foo",
				IP:   "192.1.1.1",
				Taints: []Taint{
					Taint{
						Key:    "",
						Value:  "",
						Effect: "",
					},
				},
				Labels: map[string]string{"": ""},
			},
			valid: false,
		},
		{
			n: Node{
				Host: "foo",
				IP:   "192.1.1.1",
				Taints: []Taint{
					Taint{
						Key:    "node-type:test",
						Value:  "test",
						Effect: "NoSchedule",
					},
				},
				Labels: map[string]string{"node-type:test": "test"},
			},
			valid: false,
		},
		{
			n: Node{
				Host: "foo",
				IP:   "192.1.1.1",
				Taints: []Taint{
					Taint{
						Key:    "com.foo/invalid",
						Value:  ":test",
						Effect: "NoSchedule",
					},
				},
			},
			valid: false,
		},
		{
			n: Node{
				Host: "foo",
				IP:   "192.1.1.1",
				Taints: []Taint{
					Taint{
						Key:    "node-type:test",
						Value:  ":test",
						Effect: "NoSchedule",
					},
				},
			},
			valid: false,
		},
		{
			n: Node{
				Host: "foo",
				IP:   "192.1.1.1",
				Taints: []Taint{
					Taint{
						Key:    "kismatic/foo",
						Value:  "bar",
						Effect: "",
					},
				},
			},
			valid: false,
		},
		{
			n: Node{
				Host: "foo",
				IP:   "192.1.1.1",
				Taints: []Taint{
					Taint{
						Key:    "kismatic/foo",
						Value:  "bar",
						Effect: "Foo",
					},
				},
			},
			valid: false,
		},
	}
	for i, test := range tests {
		ok, _ := test.n.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}

func TestNodeKubeletOptions(t *testing.T) {
	tests := []struct {
		nl    nodeList
		valid bool
	}{
		{
			nl: nodeList{
				[]Node{
					{
						Host: "host1",
						IP:   "10.0.0.1",
						KubeletOptions: KubeletOptions{
							Overrides: map[string]string{
								"v": "2",
							},
						},
					},
					{
						Host: "host2",
						IP:   "10.0.0.2",
						KubeletOptions: KubeletOptions{
							Overrides: map[string]string{
								"v": "2",
							},
						},
					},
				},
			},
			valid: true,
		},
		{
			nl: nodeList{
				[]Node{
					{
						Host: "host1",
						IP:   "10.0.0.1",
						KubeletOptions: KubeletOptions{
							Overrides: map[string]string{
								"v": "2",
							},
						},
					},
					{
						Host: "host1",
						IP:   "10.0.0.1",
						KubeletOptions: KubeletOptions{
							Overrides: map[string]string{
								"v": "2",
							},
						},
					},
				},
			},
			valid: true,
		},
		{
			nl: nodeList{
				[]Node{
					{
						Host: "host1",
						IP:   "10.0.0.1",
						KubeletOptions: KubeletOptions{
							Overrides: map[string]string{
								"v": "2",
							},
						},
					},
					{
						Host: "host1",
						IP:   "10.0.0.1",
						KubeletOptions: KubeletOptions{
							Overrides: map[string]string{
								"foo": "bar",
							},
						},
					},
				},
			},
			valid: false,
		},
		{
			nl: nodeList{
				[]Node{
					{
						Host: "host1",
						IP:   "10.0.0.1",
						KubeletOptions: KubeletOptions{
							Overrides: map[string]string{
								"v": "2",
							},
						},
					},
					{
						Host: "host1",
						IP:   "10.0.0.1",
						KubeletOptions: KubeletOptions{
							Overrides: map[string]string{
								"v": "3",
							},
						},
					},
				},
			},
			valid: false,
		},
		{
			nl: nodeList{
				[]Node{
					{
						Host: "host1",
						IP:   "10.0.0.1",
						KubeletOptions: KubeletOptions{
							Overrides: map[string]string{
								"v": "2",
							},
						},
					},
					{
						Host: "host1",
						IP:   "10.0.0.1",
						KubeletOptions: KubeletOptions{
							Overrides: map[string]string{
								"v":   "2",
								"foo": "bar",
							},
						},
					},
				},
			},
			valid: false,
		},
	}
	for i, test := range tests {
		ok, _ := test.nl.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}

func TestValidateFile(t *testing.T) {
	tests := []struct {
		srcPath        string
		destPath       string
		hosts          []string
		skipValidation bool
		valid          bool
	}{
		{
			srcPath:  "/bin/sh",
			destPath: "/tmp/copy_xa.yaml",
			hosts:    []string{"master01"},
			valid:    true,
		},
		{
			srcPath:  "/bin/sh",
			destPath: "/tmp/copy_xa.yaml",
			hosts:    []string{"worker01"},
			valid:    true,
		},
		{
			srcPath:  "/bin/sh",
			destPath: "/tmp/copy_xa.yaml",
			hosts:    []string{"etcd01"},
			valid:    true,
		},
		{
			srcPath:        "/tmp/xa.yaml",
			destPath:       "/tmp/copy_xa.yaml",
			hosts:          []string{"master01"},
			skipValidation: true,
			valid:          true,
		},
		{
			srcPath:  "/bin/sh",
			destPath: "/tmp/copy_xa.yaml",
			hosts:    []string{"worker"},
			valid:    true,
		},
		{
			srcPath:  "/bin/sh",
			destPath: "/tmp/copy_xa.yaml",
			hosts:    []string{"all"},
			valid:    true,
		},
		{
			srcPath:  "/bin/sh",
			destPath: "/tmp/copy_xa.yaml",
			hosts:    []string{"master", "worker"},
			valid:    true,
		},
		{
			srcPath:  "/bin/sh",
			destPath: "/tmp/copy_xa.yaml",
			hosts:    []string{"master", "worker", "worker100"},
			valid:    false,
		},
		{
			srcPath:  "",
			destPath: "",
			hosts:    []string{"master01"},
			valid:    false,
		},
		{
			srcPath:  "",
			destPath: "/tmp/copy_xa.yaml",
			hosts:    []string{"master01"},
			valid:    false,
		},
		{
			srcPath:  "/bin/sh",
			destPath: "",
			hosts:    []string{"master01"},
			valid:    false,
		},
		{
			srcPath:  "../someRelativePath",
			destPath: "../someRelativePath",
			hosts:    []string{"master01"},
			valid:    false,
		},
		{
			srcPath:  "/bin/sh",
			destPath: "../someRelativePath",
			hosts:    []string{"master01"},
			valid:    false,
		},
		{
			srcPath:  "/bin/sh",
			destPath: "/bin/sh",
			hosts:    []string{},
			valid:    false,
		},
		{
			srcPath:  "/bin/sh",
			destPath: "/bin/sh",
			hosts:    []string{"master02"},
			valid:    false,
		},
		{
			srcPath:  "/bin/sh",
			destPath: "/bin/sh",
			hosts:    []string{"foo"},
			valid:    false,
		},
	}
	for n, test := range tests {
		v := []AdditionalFile{
			AdditionalFile{
				Source:         test.srcPath,
				Destination:    test.destPath,
				Hosts:          test.hosts,
				SkipValidation: test.skipValidation,
			},
		}
		plan := validPlan()
		fg := additionalFilesGroup{AdditionalFiles: v, Plan: &plan}
		if valid, errs := fg.validate(); valid != test.valid {
			t.Errorf("test %d: expect valid = %t, but got %t", n, test.valid, valid)
			if !valid {
				t.Errorf("error %v", errs)
			}
		}
	}
}
