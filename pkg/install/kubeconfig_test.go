package install

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func createTempDirForRegenerateKubeconfigTests(t *testing.T) string {
	path, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("error creating temp dir: %v", err)
	}
	if err := os.Mkdir(filepath.Join(path, "keys"), 0744); err != nil {
		t.Fatalf("error setting up test: error creating dir: %v", err)
	}
	files := []string{"ca.pem", "admin.pem", "admin-key.pem"}
	for _, f := range files {
		if err := ioutil.WriteFile(filepath.Join(path, "keys", f), []byte(f), 0644); err != nil {
			t.Fatalf("error setting up test: error creating file: %v", err)
		}
	}
	return path
}
func TestRegenerateKubeconfigPreviousDoesNotExist(t *testing.T) {
	tempDir := createTempDirForRegenerateKubeconfigTests(t)
	defer os.Remove(tempDir)

	p := &Plan{}
	p.Cluster.Name = "test"
	p.Master.LoadBalancer = "test:6443"

	backup, err := RegenerateKubeconfig(p, tempDir)
	if err != nil {
		t.Fatalf("unexected error: %v", err)
	}

	if backup {
		t.Error("returned true when no backup was expected")
	}
	if _, err := os.Stat(filepath.Join(tempDir, kubeconfigFilename+".bak")); err == nil {
		t.Error("found unexpected kubeconfig backup")
	}
	if _, err := os.Stat(filepath.Join(tempDir, kubeconfigFilename)); os.IsNotExist(err) {
		t.Errorf("did not find expected kubeconfig file: %v", err)
	}
}

func TestRegenerateKubeconfigDifferentIsBackedUp(t *testing.T) {
	path := createTempDirForRegenerateKubeconfigTests(t)
	defer os.Remove(path)

	p := &Plan{}
	p.Cluster.Name = "old"
	p.Master.LoadBalancer = "old:6443"

	if err := GenerateKubeconfig(p, path); err != nil {
		t.Fatalf("error creating pre-existing kubeconfig file: %v", err)
	}

	p.Cluster.Name = "new"
	p.Master.LoadBalancer = "new:6443"
	backup, err := RegenerateKubeconfig(p, path)
	if err != nil {
		t.Fatalf("unexected error: %v", err)
	}

	if !backup {
		t.Error("expected true, but returned false")
	}
	if _, err := os.Stat(filepath.Join(path, kubeconfigFilename+".bak")); os.IsNotExist(err) {
		t.Error("did not find expected kubeconfig backup")
	}
	if _, err := os.Stat(filepath.Join(path, kubeconfigFilename)); os.IsNotExist(err) {
		t.Error("did not find expected kubeconfig file")
	}
}

func TestRegenerateKubeconfigSameIsNotBackedUp(t *testing.T) {
	path := createTempDirForRegenerateKubeconfigTests(t)
	defer os.Remove(path)

	p := &Plan{}
	p.Cluster.Name = "old"
	p.Master.LoadBalancer = "old:6443"

	if err := GenerateKubeconfig(p, path); err != nil {
		t.Fatalf("error creating pre-existing kubeconfig file: %v", err)
	}

	backup, err := RegenerateKubeconfig(p, path)
	if err != nil {
		t.Fatalf("unexected error: %v", err)
	}

	if backup {
		t.Error("expected false, but returned true")
	}
	if _, err := os.Stat(filepath.Join(path, kubeconfigFilename+".bak")); !os.IsNotExist(err) {
		t.Error("found unexpected backup")
	}
	if _, err := os.Stat(filepath.Join(path, kubeconfigFilename)); os.IsNotExist(err) {
		t.Error("did not find expected kubeconfig file")
	}
}
