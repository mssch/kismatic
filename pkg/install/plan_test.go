package install

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestWritePlanTemplate(t *testing.T) {
	tests := []struct {
		golden   string
		template PlanTemplateOptions
	}{
		{
			golden: "./test/plan-template.golden.yaml",
			template: PlanTemplateOptions{
				EtcdNodes:       3,
				MasterNodes:     2,
				WorkerNodes:     3,
				IngressNodes:    2,
				StorageNodes:    0,
				AdditionalFiles: 1,
			},
		},
		{
			golden: "./test/plan-template-with-storage.golden.yaml",
			template: PlanTemplateOptions{
				EtcdNodes:       3,
				MasterNodes:     2,
				WorkerNodes:     3,
				IngressNodes:    2,
				StorageNodes:    2,
				AdditionalFiles: 1,
			},
		},
	}
	for _, test := range tests {
		expected, err := ioutil.ReadFile(test.golden)
		if err != nil {
			t.Fatalf("error reading golden file: %v", err)
		}
		tmp, err := ioutil.TempDir("", "ket-test-write-plan-template")
		if err != nil {
			t.Fatalf("error creating temp dir: %v", err)
		}
		file := filepath.Join(tmp, "kismatic-cluster.yaml")
		fp := &FilePlanner{file}
		if err = WritePlanTemplate(test.template, fp); err != nil {
			t.Fatalf("error writing plan template: %v", err)
		}
		wrote, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatalf("error reading plan file template: %v", err)
		}
		if !bytes.Equal(wrote, expected) {
			t.Errorf("the resulting plan file did not equal the expected plan file (%s)", test.golden)
			if _, err := exec.LookPath("diff"); err == nil {
				cmd := exec.Command("diff", test.golden, file)
				fmt.Println(file)
				cmd.Stdout = os.Stdout
				cmd.Run()
			}
		}
	}
}

func TestReadWithDeprecated(t *testing.T) {
	pm := &DeprecatedPackageManager{
		Enabled: true,
	}
	p := &Plan{}
	p.Features = &Features{
		PackageManager: pm,
	}
	b := false
	p.Cluster.AllowPackageInstallation = &b
	readDeprecatedFields(p)

	// features.package_manager should be set to add_ons.package_manager
	if p.AddOns.PackageManager.Disable || p.AddOns.PackageManager.Provider != "helm" {
		t.Errorf("Expected add_ons.package_manager to be read from features.package_manager")
	}
	// cluster.disable_package_installation shoule be set to cluster.allow_package_installation
	if p.Cluster.DisablePackageInstallation != true {
		t.Errorf("Expected cluster.allow_package_installation to be read from cluster.disable_package_installation")
	}
}

func TestReadWithNil(t *testing.T) {
	p := &Plan{}
	setDefaults(p)

	if p.AddOns.CNI.Provider != "calico" {
		t.Errorf("Expected add_ons.cni.provider to equal 'calico', instead got %s", p.AddOns.CNI.Provider)
	}
	if p.AddOns.CNI.Options.Calico.Mode != "overlay" {
		t.Errorf("Expected add_ons.cni.options.calico.mode to equal 'overlay', instead got %s", p.AddOns.CNI.Options.Calico.Mode)
	}

	if p.AddOns.HeapsterMonitoring.Options.Heapster.Replicas != 2 {
		t.Errorf("Expected add_ons.heapster.options.heapster.replicas to equal 2, instead got %d", p.AddOns.HeapsterMonitoring.Options.Heapster.Replicas)
	}

	if p.AddOns.HeapsterMonitoring.Options.Heapster.ServiceType != "ClusterIP" {
		t.Errorf("Expected add_ons.heapster.options.heapster.service_type to equal ClusterIP, instead got %s", p.AddOns.HeapsterMonitoring.Options.Heapster.ServiceType)
	}

	if p.AddOns.HeapsterMonitoring.Options.Heapster.Sink != "influxdb:http://heapster-influxdb.kube-system.svc:8086" {
		t.Errorf("Expected add_ons.heapster.options.heapster.service_type to equal 'influxdb:http://heapster-influxdb.kube-system.svc:8086', instead got %s", p.AddOns.HeapsterMonitoring.Options.Heapster.Sink)
	}

	if p.Cluster.Certificates.CAExpiry != defaultCAExpiry {
		t.Errorf("expected ca cert expiry to be %s, but got %s", defaultCAExpiry, p.Cluster.Certificates.CAExpiry)
	}

	if p.Cluster.Version != kubernetesVersionString {
		t.Errorf("expected cluster version to be %s, but got %s", kubernetesVersionString, p.Cluster.Version)
	}
}
