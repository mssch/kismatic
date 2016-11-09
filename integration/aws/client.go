package aws

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	// Ubuntu1604LTSEast is the AMI for Ubuntu 16.04 LTS
	Ubuntu1604LTSEast = AMI("ami-29f96d3e")
	// CentOS7East is the AMI for CentOS 7
	CentOS7East = AMI("ami-6d1c2007")
	// T2Micro is the T2 Micro instance type
	T2Micro = InstanceType(ec2.InstanceTypeT2Micro)
	// T2Medium is the T2 Medium instance type
	T2Medium = InstanceType(ec2.InstanceTypeT2Medium)
	// UpdateRetry is the number of times will try before failing
	// Exponential backoff for AWS eventual consistency
	UpdateRetry = 3
)

// A Node on AWS
type Node struct {
	PrivateDNSName string
	PrivateIP      string
	PublicIP       string
	SSHUser        string
}

// AMI is the Amazon Machine Image
type AMI string

// InstanceType is the type of the Amazon machine
type InstanceType string

// ClientConfig of the AWS client
type ClientConfig struct {
	Region          string
	SubnetID        string
	Keyname         string
	SecurityGroupID string
}

// Credentials to be used for accessing the AI
type Credentials struct {
	ID     string
	Secret string
}

// Client for provisioning machines on AWS
type Client struct {
	Config      ClientConfig
	Credentials Credentials
	ec2Client   *ec2.EC2
}

func (c *Client) getAPIClient() (*ec2.EC2, error) {
	if c.ec2Client == nil {
		creds := credentials.NewStaticCredentials(c.Credentials.ID, c.Credentials.Secret, "")
		_, err := creds.Get()
		if err != nil {
			return nil, fmt.Errorf("Error with credentials provided: %v", err)
		}
		config := aws.NewConfig().WithRegion(c.Config.Region).WithCredentials(creds).WithMaxRetries(10)
		c.ec2Client = ec2.New(session.New(config))
	}
	return c.ec2Client, nil
}

// CreateNode is for creating a machine on AWS using the given AMI and InstanceType.
// Returns the ID of the newly created machine.
func (c Client) CreateNode(ami AMI, instanceType InstanceType) (string, error) {
	api, err := c.getAPIClient()
	if err != nil {
		return "", err
	}
	req := &ec2.RunInstancesInput{
		ImageId: aws.String(string(ami)),
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sda1"),
				Ebs: &ec2.EbsBlockDevice{
					DeleteOnTermination: aws.Bool(true),
					VolumeSize:          aws.Int64(8),
				},
			},
		},
		InstanceType:     aws.String(string(instanceType)),
		MinCount:         aws.Int64(1),
		MaxCount:         aws.Int64(1),
		SubnetId:         aws.String(c.Config.SubnetID),
		KeyName:          aws.String(c.Config.Keyname),
		SecurityGroupIds: []*string{aws.String(c.Config.SecurityGroupID)},
	}
	res, err := api.RunInstances(req)
	if err != nil {
		return "", err
	}
	instanceID := res.Instances[0].InstanceId
	// Modify the node
	modifyReq := &ec2.ModifyInstanceAttributeInput{
		InstanceId: instanceID,
		SourceDestCheck: &ec2.AttributeBooleanValue{
			Value: aws.Bool(false),
		},
	}
	// Use exponential backoff here due to eventual consistency nature of AWS API.
	var attempts uint
	for {
		_, err = api.ModifyInstanceAttribute(modifyReq)
		if err == nil {
			break
		}
		if err != nil && attempts == UpdateRetry { // we failed to modify the instance attributes...
			if err = c.DestroyNodes([]string{*instanceID}); err != nil {
				fmt.Printf("Failed to modify instance attributes")
				fmt.Printf("AWS NODE %q MUST BE CLEANED UP MANUALLY\n", *instanceID)
			}
			return "", err
		}
		time.Sleep((1 << attempts) * time.Second)
		attempts++
	}
	// Tag the nodes
	thisHost, _ := os.Hostname()
	tagReq := &ec2.CreateTagsInput{
		Resources: []*string{instanceID},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("ApprendaTeam"),
				Value: aws.String("Kismatic"),
			},
			{
				Key:   aws.String("CreatedBy"),
				Value: aws.String(thisHost),
			},
		},
	}
	attempts = 0
	for {
		_, err = api.CreateTags(tagReq)
		if err == nil {
			break
		}
		if err != nil && attempts == UpdateRetry { // Failed to tag the nodes after retrying a couple of times
			if err = c.DestroyNodes([]string{*instanceID}); err != nil {
				fmt.Printf("Failed to tag instance")
				fmt.Printf("AWS NODE %q MUST BE CLEANED UP MANUALLY\n", *instanceID)
			}
			return "", err
		}
		time.Sleep((1 << attempts) * time.Second)
		attempts++
	}
	return *res.Instances[0].InstanceId, nil
}

// GetNode returns information about a specific node. The consumer of this method
// is responsible for checking that the information it needs has been returned
// in the Node. (i.e. it's possible for the hostname, public IP to be empty)
func (c Client) GetNode(id string) (*Node, error) {
	api, err := c.getAPIClient()
	if err != nil {
		return nil, err
	}
	req := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(id)},
	}
	resp, err := api.DescribeInstances(req)
	if err != nil {
		return nil, err
	}
	if len(resp.Reservations) != 1 {
		return nil, fmt.Errorf("Attempted to get a single node, but API returned %d reservations", len(resp.Reservations))
	}
	if len(resp.Reservations[0].Instances) != 1 {
		return nil, fmt.Errorf("Attempted to get a single node, but API returned %d instances", len(resp.Reservations[0].Instances))
	}
	instance := resp.Reservations[0].Instances[0]

	var publicIP string
	if instance.PublicIpAddress != nil {
		publicIP = *instance.PublicIpAddress
	}
	return &Node{
		PrivateDNSName: *instance.PrivateDnsName,
		PrivateIP:      *instance.PrivateIpAddress,
		PublicIP:       publicIP,
		SSHUser:        defaultSSHUserForAMI(AMI(*instance.ImageId)),
	}, nil
}

// DestroyNodes destroys the nodes identified by the ID.
func (c Client) DestroyNodes(nodeIDs []string) error {
	api, err := c.getAPIClient()
	if err != nil {
		return err
	}
	req := &ec2.TerminateInstancesInput{
		InstanceIds: aws.StringSlice(nodeIDs),
	}
	_, err = api.TerminateInstances(req)
	if err != nil {
		return err
	}
	return nil
}

func defaultSSHUserForAMI(ami AMI) string {
	switch ami {
	case Ubuntu1604LTSEast:
		return "ubuntu"
	case CentOS7East:
		return "centos"
	default:
		panic(fmt.Sprintf("unsupported AMI: %q", ami))
	}
}
