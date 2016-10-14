// +build packet_integration

package packet

import (
	"testing"
	"time"
)

func TestNode(t *testing.T) {
	// Create node
	hostname := "testNode"
	os := CentOS7
	dev, err := CreateNode(hostname, os)
	if err != nil {
		t.Errorf("failed to create node: %v", err)
	}
	// Block until ssh is up
	deviceID := dev.ID
	timeout := 10 * time.Minute
	if err := BlockUntilNodeAccessible(deviceID, timeout); err != nil {
		t.Errorf("node did not become accessible")
	}
	// Delete node
	time.Sleep(5 * time.Second)
	if err := DeleteNode(deviceID); err != nil {
		t.Errorf("node %q was not deleted. MANUALLY CLEAN UP NODE IN PACKET.NET", deviceID)
	}
}
