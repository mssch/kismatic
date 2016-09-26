package inspector

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

// SSHCheck ensures that the target node is accessible via SSH using the provided user and authentication private key.
type SSHCheck struct {
	user                  string
	port                  int
	ipAddress             string
	clientAuthKey         string
	clientAuthKeyPassword string
}

// Check returns nil if SSH connection is successful. Otherwise returns an error message indicating the problem.
func (c *SSHCheck) Check() error {
	// verify key file exists
	if _, err := os.Stat(c.clientAuthKey); os.IsNotExist(err) {
		return fmt.Errorf("the client authentication key does not exist at %q", c.clientAuthKey)
	}

	keyBytes, err := ioutil.ReadFile(c.clientAuthKey)
	if err != nil {
		return fmt.Errorf("error reading authentication key file %q: %v", c.clientAuthKey, err)
	}

	// Handle encrypted client auth keys
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return fmt.Errorf("no key was found in %q", c.clientAuthKey)
	}
	if x509.IsEncryptedPEMBlock(block) {
		if c.clientAuthKeyPassword == "" {
			return fmt.Errorf("the client authentication key %q is encrypted, and no password was provided", c.clientAuthKey)
		}
		block.Bytes, err = x509.DecryptPEMBlock(block, []byte(c.clientAuthKeyPassword))
		if err != nil {
			return fmt.Errorf("failed to decrypt the client authentication key %q with the given password. Error was: %v", c.clientAuthKey, err)
		}
		keyBytes = pem.EncodeToMemory(block)
	}

	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return fmt.Errorf("error parsing authentication key from file %q: %v", c.clientAuthKey, err)
	}

	// create SSH config
	config := &ssh.ClientConfig{
		User: c.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	// Open SSH connection
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.ipAddress, c.port), config)
	if err != nil {
		return fmt.Errorf("unable to connect to %q on port %d. Error was: %v", c.ipAddress, c.port, err)
	}
	defer client.Close()

	return nil
}
