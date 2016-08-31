package install

import b64 "encoding/base64"
import (
	"fmt"
	"io"
	"io/ioutil"
)

// GenerateKubeconfig generate a kubeconfig file for a specific user
func GenerateKubeconfig(p *Plan, user string, certPath string, out io.Writer) error {

	// base64 encode ca and cert/key pem files
	cafile, err := ioutil.ReadFile(certPath + "/ca.pem")
	if err != nil {
		return fmt.Errorf("error reading ca file for kubeconfig: %v", err)
	}
	caEncoded := b64.StdEncoding.EncodeToString(cafile)

	certfile, err := ioutil.ReadFile(certPath + "/" + user + ".pem")
	if err != nil {
		return fmt.Errorf("error reading certificate file for kubeconfig: %v", err)
	}
	certEncoded := b64.StdEncoding.EncodeToString(certfile)

	keyfile, err := ioutil.ReadFile(certPath + "/" + user + "-key.pem")
	if err != nil {
		return fmt.Errorf("error reading certificate key file for kubeconfig: %v", err)
	}
	keyEncoded := b64.StdEncoding.EncodeToString(keyfile)

	kubeconfigTemplate := `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: %s
    server: https://%s:6443
  name: %s
contexts:
- context:
    cluster: %s
    user: %s
  name: %s
current-context: %s
kind: Config
preferences: {}
users:
- name: %s
  user:
    token: admin_password
    client-certificate-data: %s
    client-key-data: %s
`

	// Print kubeconfig
	fmt.Fprintf(
		out,
		kubeconfigTemplate,
		caEncoded,
		p.Master.LoadBalancedFQDN,
		p.Cluster.Name,
		p.Cluster.Name,
		user,
		p.Cluster.Name+"-"+user,
		p.Cluster.Name+"-"+user,
		user,
		certEncoded,
		keyEncoded)

	return nil
}
