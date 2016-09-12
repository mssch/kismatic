package integration

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"strconv"
	"text/template"
	"time"

	yaml "gopkg.in/yaml.v2"

	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jmcvetta/guid"
	homedir "github.com/mitchellh/go-homedir"
)

var guidMaker = guid.SimpleGenerator()
var leaveIt = os.Getenv("LEAVE_ARTIFACTS") != ""

const TARGET_REGION = "us-east-1"
const SUBNETID = "subnet-85f111b9"
const KEYNAME = "kismatic-integration-testing"
const SECURITYGROUPID = "sg-d1dc4dab"
const AMIUbuntu1604USEAST = "ami-29f96d3e"
const AMICentos7UsEast = "ami-6d1c2007"

var _ = Describe("Happy Path Installation Tests", func() {
	kisPath := CopyKismaticToTemp()

	BeforeSuite(func() {
		fmt.Println("Unpacking kismatic to", kisPath)
		c := exec.Command("tar", "-zxf", "../out/kismatic.tar.gz", "-C", kisPath)
		tarOut, tarErr := c.CombinedOutput()
		if tarErr != nil {
			log.Fatal("Error unpacking installer", string(tarOut), tarErr)
		}
		os.Chdir(kisPath)
	})

	AfterSuite(func() {
		if !leaveIt {
			os.RemoveAll(kisPath)
		}
	})

	Describe("Calling installer with no input", func() {
		It("should output help text", func() {
			c := exec.Command("./kismatic")
			helpbytes, helperr := c.Output()
			Expect(helperr).To(BeNil())
			helpText := string(helpbytes)
			Expect(helpText).To(ContainSubstring("Usage"))
		})
	})

	Describe("Calling installer with 'install plan'", func() {
		Context("and just hitting enter", func() {
			It("should result in the output of a well formed default plan file", func() {
				By("Outputing a file")
				c := exec.Command("./kismatic", "install", "plan")
				helpbytes, helperr := c.Output()
				Expect(helperr).To(BeNil())
				helpText := string(helpbytes)
				Expect(helpText).To(ContainSubstring("Generating installation plan file with 3 etcd nodes, 2 master nodes and 3 worker nodes"))
				Expect(FileExists("kismatic-cluster.yaml")).To(Equal(true))

				By("Outputing a file with valid YAML")
				yamlBytes, err := ioutil.ReadFile("kismatic-cluster.yaml")
				if err != nil {
					Fail("Could not read cluster file")
				}
				yamlBlob := string(yamlBytes)

				planFromYaml := ClusterPlan{}

				unmarshallErr := yaml.Unmarshal([]byte(yamlBlob), &planFromYaml)
				if unmarshallErr != nil {
					Fail("Could not unmarshall cluster yaml: %v")
				}
			})
		})
	})

	Describe("Calling installer with a plan targetting AWS", func() {
		Context("Using a 1/1/1 Ubtunu 16.04 layout", func() {
			It("should result in a working cluster", func() {
				InstallKismatic(AMIUbuntu1604USEAST, "ubuntu")
			})
		})
		Context("Using a 1/1/1 CentOS 7 layout", func() {
			It("should result in a working cluster", func() {
				InstallKismatic(AMICentos7UsEast, "centos")
			})
		})
	})
})

func InstallKismatic(nodeType string, user string) {
	By("Building a template")
	template, err := template.New("planAWSOverlay").Parse(planAWSOverlay)
	FailIfError(err, "Couldn't parse template")

	By("Making infrastructure")
	etcdNode, etcErr := MakeETCDNode(nodeType)
	FailIfError(etcErr, "Error making etcd node")

	masterNode, masterErr := MakeMasterNode(nodeType)
	FailIfError(masterErr, "Error making master node")

	workerNode, workerErr := MakeWorkerNode(nodeType)
	FailIfError(workerErr, "Error making worker node")

	defer TerminateInstances(etcdNode.Instanceid, masterNode.Instanceid, workerNode.Instanceid)
	descErr := WaitForInstanceToStart(&etcdNode, &masterNode, &workerNode)
	FailIfError(descErr, "Error waiting for nodes")
	log.Printf("Created etcd nodes: %v (%v), master nodes %v (%v), workerNodes %v (%v)",
		etcdNode.Instanceid, etcdNode.Publicip,
		masterNode.Instanceid, masterNode.Publicip,
		workerNode.Instanceid, workerNode.Publicip)

	By("Building a plan to set up an overlay network cluster on this hardware")
	nodes := PlanAWS{
		Etcd:                []AWSNodeDeets{etcdNode},
		Master:              []AWSNodeDeets{masterNode},
		Worker:              []AWSNodeDeets{workerNode},
		MasterNodeFQDN:      masterNode.Hostname,
		MasterNodeShortName: masterNode.Hostname,
		User:                user,
	}
	var hdErr error
	nodes.HomeDirectory, hdErr = homedir.Dir()
	FailIfError(hdErr, "Error getting home directory")

	f, fileErr := os.Create("kismatic-testing.yaml")
	FailIfError(fileErr, "Error waiting for nodes")
	defer f.Close()
	w := bufio.NewWriter(f)
	execErr := template.Execute(w, &nodes)
	FailIfError(execErr, "Error filling in plan template")
	w.Flush()

	By("Validing our plan")
	ver := exec.Command("./kismatic", "install", "validate", "-f", f.Name())
	verbytes, verErr := ver.CombinedOutput()
	verText := string(verbytes)

	FailIfError(verErr, "Error validating plan", verText)

	By("Punch it Chewie!")
	app := exec.Command("./kismatic", "install", "apply", "-f", f.Name())
	appbytes, appErr := app.CombinedOutput()
	appText := string(appbytes)

	FailIfError(appErr, "Error applying plan", appText)
}

func FailIfError(err error, message ...string) {
	if err != nil {
		log.Printf(message[0]+": %v\n%v", err, message[1:])
		Fail(message[0])
	}
}

func CopyKismaticToTemp() string {
	tmpDir := os.TempDir()
	randomness, randomErr := GenerateGUIDString()
	if randomErr != nil {
		log.Fatal("Error making a GUID: ", randomErr)
	}
	kisPath := tmpDir + "/kisint/" + randomness
	err := os.MkdirAll(kisPath, 0777)
	if err != nil {
		log.Fatal("Error making temp dir: ", err)
	}

	return kisPath
}

func GenerateGUIDString() (string, error) {
	randomness, randomErr := guidMaker.NextId()

	if randomErr != nil {
		return "", randomErr
	}

	return strconv.FormatInt(randomness, 16), nil
}

func AssertKismaticDirectory(kisPath string) {
	if FileExists(kisPath + "/kismatic") {
		log.Fatal("Installer unpacked but kismatic wasn't there")
	}
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func MakeETCDNode(nodeType string) (AWSNodeDeets, error) {
	return MakeAWSNode(nodeType, ec2.InstanceTypeT2Micro)
}

func MakeMasterNode(nodeType string) (AWSNodeDeets, error) {
	return MakeAWSNode(nodeType, ec2.InstanceTypeT2Micro)
}

func MakeWorkerNode(nodeType string) (AWSNodeDeets, error) {
	return MakeAWSNode(nodeType, ec2.InstanceTypeT2Medium)
}

func MakeAWSNode(ami string, instanceType string) (AWSNodeDeets, error) {
	svc := ec2.New(session.New(&aws.Config{Region: aws.String(TARGET_REGION)}))
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:          aws.String(ami),
		InstanceType:     aws.String(instanceType),
		MinCount:         aws.Int64(1),
		MaxCount:         aws.Int64(1),
		SubnetId:         aws.String(SUBNETID),
		KeyName:          aws.String(KEYNAME),
		SecurityGroupIds: []*string{aws.String(SECURITYGROUPID)},
	})

	if err != nil {
		return AWSNodeDeets{}, err
	}

	re := regexp.MustCompile("[^.]*")
	hostname := re.FindString(*runResult.Instances[0].PrivateDnsName)

	deets := AWSNodeDeets{
		Instanceid: *runResult.Instances[0].InstanceId,
		Privateip:  *runResult.Instances[0].PrivateIpAddress,
		Hostname:   hostname,
	}

	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{aws.String(deets.Instanceid)},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("ApprendaTeam"),
				Value: aws.String("Kismatic"),
			},
		},
	})
	if errtag != nil {
		return deets, errtag
	}

	return deets, nil
}

func TerminateInstances(instanceids ...string) {
	if leaveIt {
		return
	}
	awsinstanceids := make([]*string, len(instanceids))
	for i, id := range instanceids {
		awsinstanceids[i] = aws.String(id)
	}
	sess, err := session.NewSession()

	if err != nil {
		log.Printf("failed to create session: %v", err)
		return
	}

	svc := ec2.New(sess, &aws.Config{Region: aws.String(TARGET_REGION)})

	params := &ec2.TerminateInstancesInput{
		InstanceIds: awsinstanceids,
	}
	resp, err := svc.TerminateInstances(params)

	if err != nil {
		log.Printf("Could not terminate: %v", resp)
		return
	}
}

func WaitForInstanceToStart(nodes ...*AWSNodeDeets) error {
	sess, err := session.NewSession()

	if err != nil {
		fmt.Println("failed to create session,", err)
		return err
	}

	fmt.Print("Waiting for nodes to come up")
	defer fmt.Println()

	svc := ec2.New(sess, &aws.Config{Region: aws.String(TARGET_REGION)})
	for _, deets := range nodes {
		deets.Publicip = ""

		for deets.Publicip == "" {
			fmt.Print(".")
			descResult, descErr := svc.DescribeInstances(&ec2.DescribeInstancesInput{
				InstanceIds: []*string{aws.String(deets.Instanceid)},
			})
			if descErr != nil {
				return descErr
			}

			if *descResult.Reservations[0].Instances[0].State.Name == ec2.InstanceStateNameRunning &&
				descResult.Reservations[0].Instances[0].PublicIpAddress != nil {
				deets.Publicip = *descResult.Reservations[0].Instances[0].PublicIpAddress
				BlockUntilSSHOpen(deets)
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}
	return nil
}

func BlockUntilSSHOpen(node *AWSNodeDeets) {
	conn, err := net.Dial("tcp", node.Publicip+":22")
	fmt.Print("?")
	if err != nil {
		time.Sleep(5 & time.Second)
		BlockUntilSSHOpen(node)
	} else {
		conn.Close()
		return
	}
}
