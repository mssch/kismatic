package packet

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/packethost/packngo"
)

// OS is an operating system supported on Packet
type OS string

const (
	Ubuntu1604LTS = OS("ubuntu_16_04_image")
	CentOS7       = OS("centos_7_image")
)

var token = os.Getenv("PACKET_TOKEN")
var projectID = os.Getenv("PACKET_PROJECT_ID")
var sshKey = os.Getenv("PACKET_SSH_KEY")
var sshUser = "root"

var client *packngo.Client

func init() {
	if token == "" {
		log.Fatal("PACKET_TOKEN environment variable must be set")
	}
	if projectID == "" {
		log.Fatal("PACKET_PROJECT_ID environment variable must be set")
	}
	if sshKey == "" {
		log.Fatal("PACKET_SSH_KEY environment variable must be set")
	}
	client = packngo.NewClient("", token, http.DefaultClient)
}

// CreateNode creates a node in packet with the given hostname and OS
func CreateNode(hostname string, os OS) (*packngo.Device, error) {
	device := &packngo.DeviceCreateRequest{
		HostName:     hostname,
		OS:           string(os),
		Tags:         []string{"integration-test"},
		ProjectID:    projectID,
		Plan:         "baremetal_0",
		BillingCycle: "hourly",
		Facility:     "ewr1",
	}
	dev, _, err := client.Devices.Create(device)
	if err != nil {
		return nil, err
	}
	return dev, nil
}

// DeleteNode deletes the node that matches the given ID
func DeleteNode(deviceID string) error {
	resp, err := client.Devices.Delete(deviceID)
	if err != nil {
		return fmt.Errorf("failed to delete node with ID %q", deviceID)
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete node with ID %q", deviceID)
	}
	return nil
}

// BlockUntilNodeAccessible blocks until the given node is accessible,
// or the timeout is reached.
func BlockUntilNodeAccessible(deviceID string, timeout time.Duration) error {
	timeoutChan := make(chan bool, 1)
	go func() {
		time.Sleep(timeout)
		timeoutChan <- true
	}()
	fmt.Printf("Waiting for node %s to be accessible", deviceID)
	// Loop until we get the node IP
	var nodeIP string
	for {
		select {
		case <-timeoutChan:
			return fmt.Errorf("timed out waiting for node to be accessible")
		default:
			dev, _, err := client.Devices.Get(deviceID)
			if err != nil {
				continue
			}
			nodeIP = getPublicIPv4(dev)
			fmt.Printf("\nGot node's IP: %s\n", nodeIP)
		}
		if nodeIP != "" {
			break
		}
		fmt.Print(".")
		time.Sleep(5 * time.Second)
	}
	// Loop until state is active
	fmt.Print("Waiting for node state to be 'active'")
	active := false
	for {
		select {
		case <-timeoutChan:
			return fmt.Errorf("timedout waiting for node to be active")
		default:
			dev, _, err := client.Devices.Get(deviceID)
			if err != nil {
				continue
			}
			if dev.State == "active" {
				active = true
			}
		}
		if active {
			break
		}
		fmt.Print(".")
		time.Sleep(30 * time.Second)
	}
	// Loop until timeout or ssh is accessible
	fmt.Printf("Waiting for SSH access")
	for {
		select {
		case <-timeoutChan:
			return fmt.Errorf("timed out waiting for node to be accessible")
		default:
			if sshAccessible(nodeIP) {
				fmt.Printf("\nSSH is GO!\n")
				return nil
			}
		}
		fmt.Print(".")
		time.Sleep(5 * time.Second)
	}
}

func getPublicIPv4(device *packngo.Device) string {
	for _, net := range device.Network {
		if net.Public != true || net.AddressFamily != 4 {
			continue
		}
		if net.Address != "" {
			return net.Address
		}
	}
	return ""
}

func sshAccessible(ip string) bool {
	cmd := exec.Command("ssh")
	cmd.Args = append(cmd.Args, "-i", sshKey)
	cmd.Args = append(cmd.Args, "-o", "ConnectTimeout=5")
	cmd.Args = append(cmd.Args, "-o", "BatchMode=yes")
	cmd.Args = append(cmd.Args, "-o", "StrictHostKeyChecking=no")
	cmd.Args = append(cmd.Args, fmt.Sprintf("%s@%s", sshUser, ip), "exit") // just call exit if we are able to connect
	err := cmd.Run()
	return err == nil
}
