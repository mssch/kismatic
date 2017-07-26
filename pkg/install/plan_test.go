package install

import "testing"

func TestGenerateAlphaNumericPassword(t *testing.T) {
	_, err := generateAlphaNumericPassword()
	if err != nil {
		t.Error(err)
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
}
