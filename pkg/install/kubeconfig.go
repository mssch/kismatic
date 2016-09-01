package install

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
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
	cafile, err := ioutil.ReadFile(certPath + "/ca.pem")
	if err != nil {
		return fmt.Errorf("error reading ca file for kubeconfig: %v", err)
	}
	caEncoded := b64.StdEncoding.EncodeToString(cafile)
	// Base64 encoded cert
	certfile, err := ioutil.ReadFile(certPath + "/" + user + ".pem")
	if err != nil {
		return fmt.Errorf("error reading certificate file for kubeconfig: %v", err)
	}
	certEncoded := b64.StdEncoding.EncodeToString(certfile)
	// Base64 encoded key
	keyfile, err := ioutil.ReadFile(certPath + "/" + user + "-key.pem")
	if err != nil {
		return fmt.Errorf("error reading certificate key file for kubeconfig: %v", err)
	}
	keyEncoded := b64.StdEncoding.EncodeToString(keyfile)

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
