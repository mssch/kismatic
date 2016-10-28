// +build packet_integration

package packet

import (
	"os"
	"testing"
	"time"
)

func TestNode(t *testing.T) {
	// Create node
	c := Client{
		Token:     os.Getenv("PACKET_TOKEN"),
		ProjectID: os.Getenv("PACKET_PROJECT_ID"),
	}

	hostname := "testNode"
	osImage := CentOS7
	dev, err := c.CreateNode(hostname, osImage)
	if err != nil {
		t.Errorf("failed to create node: %v", err)
	}
	// Block until ssh is up
	deviceID := dev.ID
	timeout := 10 * time.Minute
	if err := c.BlockUntilNodeAccessible(deviceID, timeout, os.Getenv("PACKET_SSH_KEY"), "root"); err != nil {
		t.Errorf("node did not become accessible")
	}
	// Delete node
	time.Sleep(5 * time.Second)
	if err := c.DeleteNode(deviceID); err != nil {
		t.Errorf("node %q was not deleted. MANUALLY CLEAN UP NODE IN PACKET.NET", deviceID)
	}
}
