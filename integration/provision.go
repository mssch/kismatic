package integration

import (
	"fmt"
	"os"
	"time"

	"github.com/apprenda/kismatic-platform/integration/aws"
)

const (
	Ubuntu1604LTS = linuxDistro("ubuntu1604LTS")
	CentOS7       = linuxDistro("centos7")
)

type linuxDistro string

type NodeCount struct {
	Etcd   uint16
	Master uint16
	Worker uint16
}

type provisionedNodes struct {
	etcd   []AWSNodeDeets
	master []AWSNodeDeets
	worker []AWSNodeDeets
}

func (p provisionedNodes) allNodes() []AWSNodeDeets {
	n := []AWSNodeDeets{}
	n = append(n, p.etcd...)
	n = append(n, p.master...)
	n = append(n, p.worker...)
	return n
}

type AWSNodeDeets struct {
	id        string
	Hostname  string
	PublicIP  string
	PrivateIP string
}

func (nc NodeCount) Total() uint16 {
	return nc.Etcd + nc.Master + nc.Worker
}

const (
	AWSTargetRegion     = "us-east-1"
	AWSSubnetID         = "subnet-25e13d08"
	AWSKeyName          = "kismatic-integration-testing"
	AWSSecurityGroupID  = "sg-d1dc4dab"
	AMIUbuntu1604USEAST = "ami-29f96d3e"
	AMICentos7UsEast    = "ami-6d1c2007"
)

func awsClientFromEnvironment() (*aws.Client, bool) {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if accessKeyID == "" || secretAccessKey == "" {
		return nil, false
	}
	c := aws.Client{
		Config: aws.ClientConfig{
			Region:          AWSTargetRegion,
			SubnetID:        AWSSubnetID,
			Keyname:         AWSKeyName,
			SecurityGroupID: AWSSecurityGroupID,
		},
		Credentials: aws.Credentials{
			ID:     accessKeyID,
			Secret: secretAccessKey,
		},
	}
	overrideRegion := os.Getenv("AWS_TARGET_REGION")
	if overrideRegion != "" {
		c.Config.Region = overrideRegion
	}
	overrideSubnet := os.Getenv("AWS_SUBNET_ID")
	if overrideSubnet != "" {
		c.Config.SubnetID = overrideSubnet
	}
	overrideKeyName := os.Getenv("AWS_KEY_NAME")
	if overrideKeyName != "" {
		c.Config.Keyname = overrideKeyName
	}
	overrideSecGroup := os.Getenv("AWS_SECURITY_GROUP_ID")
	if overrideSecGroup != "" {
		c.Config.SecurityGroupID = overrideSecGroup
	}
	return &c, true
}

func provisionAWSNodes(client *aws.Client, nodeCount NodeCount, distro linuxDistro) (provisionedNodes, error) {
	var ami aws.AMI
	switch distro {
	case Ubuntu1604LTS:
		ami = aws.Ubuntu1604LTSEast
	case CentOS7:
		ami = aws.CentOS7East
	default:
		panic(fmt.Sprintf("Used an unsupported distribution: %s", distro))
	}
	provisioned := provisionedNodes{}
	var i uint16
	for i = 0; i < nodeCount.Etcd; i++ {
		nodeID, err := client.CreateNode(ami, aws.T2Micro)
		if err != nil {
			return provisioned, err
		}
		provisioned.etcd = append(provisioned.etcd, AWSNodeDeets{id: nodeID})
	}
	for i = 0; i < nodeCount.Master; i++ {
		nodeID, err := client.CreateNode(ami, aws.T2Micro)
		if err != nil {
			return provisioned, err
		}
		provisioned.master = append(provisioned.master, AWSNodeDeets{id: nodeID})
	}
	for i = 0; i < nodeCount.Worker; i++ {
		nodeID, err := client.CreateNode(ami, aws.T2Medium)
		if err != nil {
			return provisioned, err
		}
		provisioned.worker = append(provisioned.worker, AWSNodeDeets{id: nodeID})
	}
	// Wait until all instances have their public IPs assigned
	for i := range provisioned.etcd {
		etcd := &provisioned.etcd[i]
		node, err := waitForPublicIP(client, etcd.id)
		if err != nil {
			return provisioned, err
		}
		etcd.Hostname = node.Hostname
		etcd.PrivateIP = node.PrivateIP
		etcd.PublicIP = node.PublicIP
	}
	for i := range provisioned.master {
		master := &provisioned.master[i]
		node, err := waitForPublicIP(client, master.id)
		if err != nil {
			return provisioned, err
		}
		master.Hostname = node.Hostname
		master.PrivateIP = node.PrivateIP
		master.PublicIP = node.PublicIP
	}
	for i := range provisioned.worker {
		worker := &provisioned.worker[i]
		node, err := waitForPublicIP(client, worker.id)
		if err != nil {
			return provisioned, err
		}
		worker.Hostname = node.Hostname
		worker.PrivateIP = node.PrivateIP
		worker.PublicIP = node.PublicIP
	}
	return provisioned, nil
}

func waitForPublicIP(client *aws.Client, nodeID string) (*aws.Node, error) {
	for {
		fmt.Print(".")
		node, err := client.GetNode(nodeID)
		if err != nil {
			return nil, err
		}
		if node.PublicIP != "" {
			fmt.Println()
			return node, nil
		}
		time.Sleep(5 * time.Second)
	}
}

func terminateNodes(client *aws.Client, runningNodes provisionedNodes) error {
	nodes := runningNodes.allNodes()
	nodeIDs := []string{}
	for _, n := range nodes {
		nodeIDs = append(nodeIDs, n.id)
	}
	return client.DestroyNodes(nodeIDs)
}

func waitForSSH(provisionedNodes provisionedNodes, sshUser, sshKey string) error {
	nodes := provisionedNodes.allNodes()
	for _, n := range nodes {
		BlockUntilSSHOpen(n.PublicIP, sshUser, sshKey)
	}
	return nil
}
