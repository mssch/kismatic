package aws

import (
	"fmt"
	"os"
	"regexp"

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
)

// A Node on AWS
type Node struct {
	Hostname  string
	PrivateIP string
	PublicIP  string
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
		config := aws.NewConfig().WithRegion(c.Config.Region).WithCredentials(creds)
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
	_, err = api.ModifyInstanceAttribute(modifyReq)
	if err != nil {
		if err := c.DestroyNodes([]string{*instanceID}); err != nil {
			fmt.Printf("AWS NODE %q MUST BE CLEANED UP MANUALLY\n", instanceID)
		}
		return "", err
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
	if _, err = api.CreateTags(tagReq); err != nil {
		if err := c.DestroyNodes([]string{*instanceID}); err != nil {
			fmt.Printf("AWS NODE %q MUST BE CLEANED UP MANUALLY\n", instanceID)
		}
		return "", err
	}
	return *res.Instances[0].InstanceId, nil
}

// GetNode returns information about a specific node
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
	re := regexp.MustCompile("[^.]*")
	hostname := re.FindString(*instance.PrivateDnsName)
	if hostname == "" {
		return nil, fmt.Errorf("Failed to get hostname from instance's DNS name %q", *instance.PrivateDnsName)
	}
	var publicIP string
	if instance.PublicIpAddress != nil {
		publicIP = *instance.PublicIpAddress
	}
	return &Node{
		Hostname:  hostname,
		PrivateIP: *instance.PrivateIpAddress,
		PublicIP:  publicIP,
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
