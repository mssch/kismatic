// +build aws_client_integration

package aws

import (
	"fmt"
	"os"
	"testing"
)

// This isn't really a test, but just a method to exercise
// the client. Verification is performed by looking at the AWS
// dashboard.
func TestClient(t *testing.T) {
	// Creata a node
	c := Client{
		Credentials: Credentials{
			ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
			Secret: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
		Config: ClientConfig{
			Region:          os.Getenv("AWS_DEFAULT_REGION"),
			SubnetID:        os.Getenv("AWS_SUBNET_ID"),
			Keyname:         os.Getenv("AWS_KEYNAME"),
			SecurityGroupID: os.Getenv("AWS_SECURITY_GROUP_ID"),
		},
	}
	fmt.Println("Creating node")
	nodeID, err := c.CreateNode(Ubuntu1604LTSEast, T2Micro)
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
	fmt.Printf("Created a node with ID %q\n", nodeID)
	node, err := c.GetNode(nodeID)
	if err != nil {
		t.Fatalf("Failed to get node details: %v", err)
	}
	fmt.Println(node)
	if err := c.DestroyNodes([]string{nodeID}); err != nil {
		t.Fatalf("Failed to destroy node: %v", err)
	}
}
