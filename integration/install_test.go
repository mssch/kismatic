package integration

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
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
)

var guidMaker = guid.SimpleGenerator()

const TARGET_REGION = "us-east-1"
const SUBNETID = "subnet-85f111b9"
const KEYNAME = "kismatic-integration-testing"
const SECURITYGROUPID = "sg-d1dc4dab"
const AMIUbuntu1604USEAST = "ami-29f96d3e"

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
		//os.RemoveAll(kisPath)
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
		Context("Using a 1/1/1 Ubtunu 16.05 layout", func() {
			It("should result in a working cluster", func() {
				By("Building a template")
				template, err := template.New("planUbuntuAWSOverlay").Parse(planUbuntuAWSOverlay)
				if err != nil {
					log.Printf("Error parsing template %v", err)
					Fail("Couldn't parse template")
				}

				By("Making infrastructure")
				etcdNode, etcErr := MakeETCDNode()
				if etcErr != nil {
					log.Printf("Error making etcd node %v", etcErr)
					Fail("Error making etcd node")
				}
				masterNode, masterErr := MakeMasterNode()
				if masterErr != nil {
					log.Printf("Error making master node %v", masterErr)
					Fail("Error making master node")
				}
				workerNode, workerErr := MakeWorkerNode()
				if workerErr != nil {
					log.Printf("Error making worker node %v", etcErr)
					Fail("Error making worker node")
				}
				//defer TerminateInstances(etcdNode.Instanceid, masterNode.Instanceid, workerNode.Instanceid)
				descErr := WaitForInstanceToStart(&etcdNode, &masterNode, &workerNode)
				if descErr != nil {
					log.Printf("Error waiting for nodes %v")
					Fail("Error waiting for nodes")
				}
				log.Printf("Created etcd nodes: %v, master nodes %v, workerNodes %v", etcdNode.Instanceid, masterNode.Instanceid, workerNode.Instanceid)

				By("Building a plan to set up an overlay network cluster on this hardware")
				nodes := PlanUbuntuAWS{
					Etcd:                []AWSNodeDeets{etcdNode},
					Master:              []AWSNodeDeets{masterNode},
					Worker:              []AWSNodeDeets{workerNode},
					MasterNodeFQDN:      masterNode.Hostname,
					MasterNodeShortName: masterNode.Hostname,
				}
				f, fileErr := os.Create("kismatic-1-1-1-ubuntu.yaml")
				FailIfError(fileErr, "Error waiting for nodes")
				defer f.Close()
				w := bufio.NewWriter(f)
				execErr := template.Execute(w, &nodes)
				FailIfError(execErr, "Error filling in plan template")
				w.Flush()

				By("Validing our plan")
				// home := os.Getenv("HOME")
				// cop := exec.Command("cp", home+"/.ssh/kismatic-integration-testing.pem", ".")
				// out, copyErr := cop.CombinedOutput()
				// log.Println(string(out))
				// FailIfError(copyErr, "Error copying key file to working directory")

				ver := exec.Command("./kismatic", "install", "validate", "-f", f.Name())
				verbytes, verErr := ver.CombinedOutput()
				verText := string(verbytes)

				FailIfError(verErr, "Error validating plan", verText)

				By("Punch it Chewie!")
				app := exec.Command("./kismatic", "install", "apply", "-f", f.Name())
				appbytes, appErr := app.CombinedOutput()
				appText := string(appbytes)

				FailIfError(appErr, "Error applying plan", appText)
			})
		})
	})
})

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

func MakeETCDNode() (AWSNodeDeets, error) {
	return MakeAWSNode(AMIUbuntu1604USEAST, ec2.InstanceTypeT2Micro)
}

func MakeMasterNode() (AWSNodeDeets, error) {
	return MakeAWSNode(AMIUbuntu1604USEAST, ec2.InstanceTypeT2Micro)
}

func MakeWorkerNode() (AWSNodeDeets, error) {
	return MakeAWSNode(AMIUbuntu1604USEAST, ec2.InstanceTypeT2Medium)
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

	deets := AWSNodeDeets{
		Instanceid: *runResult.Instances[0].InstanceId,
		Privateip:  *runResult.Instances[0].PrivateIpAddress,
		Hostname:   *runResult.Instances[0].PrivateDnsName,
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
		deets.Publicip = "" //not returned immediately; eep!

		for deets.Publicip == "" {
			fmt.Print(".")
			descResult, descErr := svc.DescribeInstances(&ec2.DescribeInstancesInput{
				InstanceIds: []*string{aws.String(deets.Instanceid)},
			})
			if descErr != nil {
				return descErr
			}

			time.Sleep(500 * time.Millisecond)

			if *descResult.Reservations[0].Instances[0].State.Name == ec2.InstanceStateNameRunning &&
				descResult.Reservations[0].Instances[0].PublicIpAddress != nil {
				deets.Publicip = *descResult.Reservations[0].Instances[0].PublicIpAddress
				deets.Hostname = *descResult.Reservations[0].Instances[0].PublicDnsName
			}
		}
	}
	return nil
}
