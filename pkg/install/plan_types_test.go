package install

import (
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestCanReadAPIServerOverrides(t *testing.T) {
	d, _ := ioutil.ReadFile("test/cluster-config.yaml")
	p := &Plan{}
	yaml.Unmarshal(d, p)

	assertEqual(t, p.Cluster.APIServerOptions.Overrides["runtime-config"], "beta/v2api=true,alpha/v1api=true")
}

func TestClusterAddress(t *testing.T) {
	tests := []struct {
		plan  Plan
		valid bool
		host  string
		port  string
	}{
		{
			plan: Plan{
				Master: MasterNodeGroup{
					LoadBalancer: "lb:6443",
				},
			},
			valid: true,
			host:  "lb",
			port:  "6443",
		},
		{
			plan: Plan{
				Master: MasterNodeGroup{
					LoadBalancer: "lb:443",
				},
			},
			valid: true,
			host:  "lb",
			port:  "443",
		},
		{
			plan: Plan{
				Master: MasterNodeGroup{
					LoadBalancer: "lb",
				},
			},
			valid: false,
		},
		{
			plan: Plan{
				Master: MasterNodeGroup{
					LoadBalancer: "",
				},
			},
			valid: false,
		},
	}
	for _, test := range tests {
		host, port, err := test.plan.ClusterAddress()
		if test.valid != (err == nil) {
			t.Fatalf("expected err to be %q, instead got %q", test.valid, err)
		}
		if test.host != host {
			t.Errorf("expected host to be %q, instead got %q", test.host, host)
		}
		if test.port != port {
			t.Errorf("expected host to be %q, instead got %q", test.port, port)
		}
	}

}
