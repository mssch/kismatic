package install

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apprenda/kismatic/pkg/util"
)

const kubeconfigFilename = "kubeconfig"

// ConfigOptions sds
type ConfigOptions struct {
	CA      string
	Server  string
	Cluster string
	User    string
	Context string
	Cert    string
	Key     string
}

var kubeconfigTemplate = `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: {{.CA}}
    server: {{.Server}}
  name: {{.Cluster}}
contexts:
- context:
    cluster: {{.Cluster}}
    user: {{.User}}
  name: {{.Context}}
current-context: {{.Context}}
kind: Config
preferences: {}
users:
- name: {{.User}}
  user:
    client-certificate-data: {{.Cert}}
    client-key-data: {{.Key}}
`

// GenerateKubeconfig generate a kubeconfig file for a specific user
func GenerateKubeconfig(p *Plan, generatedAssetsDir string) error {
	user := "admin"
	server := "https://" + p.Master.LoadBalancedFQDN + ":6443"
	cluster := p.Cluster.Name
	context := p.Cluster.Name + "-" + user

	certsDir := filepath.Join(generatedAssetsDir, "keys")

	// Base64 encoded ca
	caEncoded, err := util.Base64String(filepath.Join(certsDir, "ca.pem"))
	if err != nil {
		return fmt.Errorf("error reading ca file for kubeconfig: %v", err)
	}
	// Base64 encoded cert
	certEncoded, err := util.Base64String(filepath.Join(certsDir, user+".pem"))
	if err != nil {
		return fmt.Errorf("error reading certificate file for kubeconfig: %v", err)
	}
	// Base64 encoded key
	keyEncoded, err := util.Base64String(filepath.Join(certsDir, user+"-key.pem"))
	if err != nil {
		return fmt.Errorf("error reading certificate key file for kubeconfig: %v", err)
	}

	// Process template file
	tmpl, err := template.New("kubeconfig").Parse(kubeconfigTemplate)
	if err != nil {
		return fmt.Errorf("error reading config template: %v", err)
	}
	configOptions := ConfigOptions{caEncoded, server, cluster, user, context, certEncoded, keyEncoded}
	var kubeconfig bytes.Buffer
	err = tmpl.Execute(&kubeconfig, configOptions)
	if err != nil {
		return fmt.Errorf("error processing config template: %v", err)
	}
	// Write config file
	kubeconfigFile := filepath.Join(generatedAssetsDir, kubeconfigFilename)
	err = ioutil.WriteFile(kubeconfigFile, kubeconfig.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing kubeconfig file: %v", err)
	}

	return nil
}

// RegenerateKubeconfig backs up the old kubeconfig file if it exists. Returns
// true if the new kubeconfig file is different than the previous one.
// Otherwise returns false.
func RegenerateKubeconfig(p *Plan, generatedAssetsDir string) (bool, error) {
	kubeconfigFile := filepath.Join(generatedAssetsDir, kubeconfigFilename)
	kubeconfigBackup := filepath.Join(generatedAssetsDir, kubeconfigFilename) + ".bak"

	err := os.Rename(kubeconfigFile, kubeconfigBackup)
	if os.IsNotExist(err) {
		// Nothing else to do as the old kubeconfig does not exist
		return false, GenerateKubeconfig(p, generatedAssetsDir)
	}
	if err != nil {
		return false, fmt.Errorf("error backing up existing kubeconfig file: %v", err)
	}

	if err := GenerateKubeconfig(p, generatedAssetsDir); err != nil {
		return false, err
	}

	// Check if the new kubeconfig is different than the previous one
	old, err := ioutil.ReadFile(kubeconfigBackup)
	if err != nil {
		return false, fmt.Errorf("error reading file %q: %v", kubeconfigBackup, err)
	}
	new, err := ioutil.ReadFile(kubeconfigFile)
	if err != nil {
		return false, fmt.Errorf("error reading file %q: %v", kubeconfigFile, err)
	}
	if bytes.Equal(old, new) {
		os.Remove(kubeconfigBackup)
		return false, nil
	}

	return true, nil
}
