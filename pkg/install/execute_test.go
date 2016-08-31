package install

import "testing"

func TestBuildNodeInventory(t *testing.T) {
	p := &Plan{}
	p.Cluster.SSH.Key = "id_rsa"
	p.Cluster.SSH.Port = 2222
	p.Cluster.SSH.User = "alice"

	p.Etcd.Nodes = []Node{
		{
			Host:       "etcd01",
			IP:         "10.0.0.1",
			InternalIP: "192.168.0.11",
		},
	}
	p.Master.Nodes = []Node{
		{
			Host:       "master01",
			IP:         "10.0.0.2",
			InternalIP: "192.168.0.12",
		},
	}
	p.Worker.Nodes = []Node{
		{
			Host:       "worker01",
			IP:         "10.0.0.3",
			InternalIP: "192.168.0.13",
		},
		{
			Host:       "worker02",
			IP:         "10.0.0.4",
			InternalIP: "192.168.0.14",
		},
	}

	expected := `[etcd]
etcd01 ansible_host=10.0.0.1 internal_ipv4=192.168.0.11 ansible_ssh_private_key_file=id_rsa ansible_port=2222 ansible_user=alice
[master]
master01 ansible_host=10.0.0.2 internal_ipv4=192.168.0.12 ansible_ssh_private_key_file=id_rsa ansible_port=2222 ansible_user=alice
[worker]
worker01 ansible_host=10.0.0.3 internal_ipv4=192.168.0.13 ansible_ssh_private_key_file=id_rsa ansible_port=2222 ansible_user=alice
worker02 ansible_host=10.0.0.4 internal_ipv4=192.168.0.14 ansible_ssh_private_key_file=id_rsa ansible_port=2222 ansible_user=alice
`

	inv := buildNodeInventory(p)

	if expected != string(inv) {
		t.Errorf("expected inventory: \n%s\n but got: \n%s\n", expected, string(inv))
	}

}
