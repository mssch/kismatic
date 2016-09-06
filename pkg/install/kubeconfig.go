package install

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"

	"github.com/apprenda/kismatic-platform/pkg/util"
)

// ConfigOptions sds
type ConfigOptions struct {
	CA      string
	Server  string
	Cluster string
	User    string
	Context string
	Token   string
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
    token: {{.Token}}
    client-certificate-data: {{.Cert}}
    client-key-data: {{.Key}}
`

// GenerateKubeconfig generate a kubeconfig file for a specific user
func GenerateKubeconfig(p *Plan, certPath string) error {
	user := "admin"
	token := p.Cluster.AdminPassword
	server := "https://" + p.Master.LoadBalancedFQDN + ":6443"
	cluster := p.Cluster.Name
	context := p.Cluster.Name + "-" + user

	// Base64 encoded ca
	caEncoded, err := util.Base64String(filepath.Join(certPath, "ca.pem"))
	if err != nil {
		return fmt.Errorf("error reading ca file for kubeconfig: %v", err)
	}
	// Base64 encoded cert
	certEncoded, err := util.Base64String(filepath.Join(certPath, user+".pem"))
	if err != nil {
		return fmt.Errorf("error reading certificate file for kubeconfig: %v", err)
	}
	// Base64 encoded key
	keyEncoded, err := util.Base64String(filepath.Join(certPath, user+"-key.pem"))
	if err != nil {
		return fmt.Errorf("error reading certificate key file for kubeconfig: %v", err)
	}

	// Process template file
	tmpl, err := template.New("kubeconfig").Parse(kubeconfigTemplate)
	if err != nil {
		return fmt.Errorf("error reading config template: %v", err)
	}
	configOptions := ConfigOptions{caEncoded, server, cluster, user, context, token, certEncoded, keyEncoded}
	var kubeconfig bytes.Buffer
	err = tmpl.Execute(&kubeconfig, configOptions)
	if err != nil {
		return fmt.Errorf("error processing config template: %v", err)
	}
	// Write config file
	kubeconfigFile := "config"
	err = ioutil.WriteFile(kubeconfigFile, kubeconfig.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing kubeconfig file: %v", err)
	}

	return nil
}
