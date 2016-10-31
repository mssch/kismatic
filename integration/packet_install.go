package integration

//
// import (
// 	"bufio"
// 	"fmt"
// 	"html/template"
// 	"os"
// 	"os/exec"
// 	"time"
//
// 	"github.com/apprenda/kismatic-platform/integration/packet"
//
// 	. "github.com/onsi/ginkgo"
// )
//
// func packetInstallKismatic(nodeOS packet.OS, nodeCount NodeCount) {
// 	totalNodes := int(nodeCount.Total())
// 	createdNodes := []*packet.Node{}
// 	By("Making infrastructure")
// 	defer func() {
// 		if leaveIt() {
// 			return
// 		}
// 		for _, n := range createdNodes {
// 			if err := packet.DeleteNode(n.ID); err != nil {
// 				fmt.Printf("FAILED TO DELETE NODE %q - MUST BE DELETED MANUALLY\n", n.ID)
// 			}
// 		}
// 	}()
// 	// Create all the infrastructure we need
// 	testTime := time.Now()
// 	for i := 0; i < totalNodes; i++ {
// 		name := fmt.Sprintf("integration-%d-%d", i, testTime.Unix())
// 		node, err := packet.CreateNode(name, nodeOS)
// 		FailIfError(err, "error creating node")
// 		fmt.Println("Created node", node.ID)
// 		createdNodes = append(createdNodes, node)
// 	}
// 	// Wait until all infra is ready to go
// 	for _, n := range createdNodes {
// 		err := packet.BlockUntilNodeAccessible(n.ID, 10*time.Minute)
// 		FailIfError(err, "node %q did not come up in time: %v", n.ID)
// 	}
// 	if len(createdNodes) != totalNodes {
// 		Fail("created nodes not equal total requested nodes...")
// 	}
//
// 	By("Building a plan file")
// 	etcdNodes := []AWSNodeDeets{}
// 	currentNode := 0
// 	for i := 0; i < int(nodeCount.Etcd); i++ {
// 		node, err := packet.GetNode(createdNodes[currentNode].ID)
// 		FailIfError(err, "failed to get node %q", createdNodes[currentNode].ID)
// 		deets := AWSNodeDeets{
// 			Hostname:  node.Host,
// 			Publicip:  node.PublicIPv4,
// 			Privateip: node.PrivateIPv4,
// 		}
// 		etcdNodes = append(etcdNodes, deets)
// 		currentNode++
// 	}
//
// 	masterNodes := []AWSNodeDeets{}
// 	for i := 0; i < int(nodeCount.Master); i++ {
// 		node, err := packet.GetNode(createdNodes[currentNode].ID)
// 		FailIfError(err, "failed to get node %q", createdNodes[currentNode].ID)
// 		deets := AWSNodeDeets{
// 			Hostname:  node.Host,
// 			Publicip:  node.PublicIPv4,
// 			Privateip: node.PrivateIPv4,
// 		}
// 		masterNodes = append(masterNodes, deets)
// 		currentNode++
// 	}
//
// 	workerNodes := []AWSNodeDeets{}
// 	for i := 0; i < int(nodeCount.Worker); i++ {
// 		node, err := packet.GetNode(createdNodes[currentNode].ID)
// 		FailIfError(err, "failed to get node %q", createdNodes[currentNode].ID)
// 		deets := AWSNodeDeets{
// 			Hostname:  node.Host,
// 			Publicip:  node.PublicIPv4,
// 			Privateip: node.PrivateIPv4,
// 		}
// 		workerNodes = append(workerNodes, deets)
// 		currentNode++
// 	}
//
// 	nodes := PlanAWS{
// 		Etcd:                etcdNodes,
// 		Master:              masterNodes,
// 		Worker:              workerNodes,
// 		MasterNodeFQDN:      masterNodes[0].Hostname,
// 		MasterNodeShortName: masterNodes[0].Hostname,
// 		SSHKeyFile:          packet.SSHKey,
// 		SSHUser:             "root",
// 	}
//
// 	// Create template
// 	template, err := template.New("planOverlay").Parse(planAWSOverlay)
// 	FailIfError(err, "Couldn't parse template")
// 	f, fileErr := os.Create("kismatic-testing.yaml")
// 	FailIfError(fileErr, "Error waiting for nodes")
// 	defer f.Close()
// 	w := bufio.NewWriter(f)
// 	execErr := template.Execute(w, &nodes)
// 	FailIfError(execErr, "Error filling in plan template")
// 	w.Flush()
//
// 	By("Validing our plan")
// 	ver := exec.Command("./kismatic", "install", "validate", "-f", f.Name())
// 	verbytes, verErr := ver.CombinedOutput()
// 	verText := string(verbytes)
// 	FailIfError(verErr, "Error validating plan", verText)
// 	if bailBeforeAnsible() {
// 		return
// 	}
//
// 	By("Punch it Chewie!")
// 	app := exec.Command("./kismatic", "install", "apply", "-f", f.Name())
// 	app.Stdout = os.Stdout
// 	app.Stderr = os.Stderr
// 	appErr := app.Run()
// 	FailIfError(appErr, "Error applying plan")
// }
//
// func packetInstallKismaticMini(nodeOS packet.OS) {
// 	By("Building a template")
// 	template, err := template.New("planOverlay").Parse(planAWSOverlay)
// 	FailIfError(err, "Couldn't parse template")
//
// 	By("Making infrastructure")
// 	name := "minikube01.integ.test"
// 	node, err := packet.CreateNode(name, nodeOS)
// 	FailIfError(err, "failed to create node")
//
// 	if !leaveIt() {
// 		defer packet.DeleteNode(node.ID)
// 	}
//
// 	err = packet.BlockUntilNodeAccessible(node.ID, 10*time.Minute)
// 	FailIfError(err, "node did not become accessible")
//
// 	By("Building a plan to set up an overlay network cluster on this hardware")
// 	node, err = packet.GetNode(node.ID)
// 	FailIfError(err, "failed to get node details")
// 	nodeDeets := AWSNodeDeets{
// 		Publicip:  node.PublicIPv4,
// 		Privateip: node.PrivateIPv4,
// 		Hostname:  node.Host,
// 	}
// 	nodes := PlanAWS{
// 		Etcd:                []AWSNodeDeets{nodeDeets},
// 		Master:              []AWSNodeDeets{nodeDeets},
// 		Worker:              []AWSNodeDeets{nodeDeets},
// 		MasterNodeFQDN:      nodeDeets.Hostname,
// 		MasterNodeShortName: nodeDeets.Hostname,
// 		SSHKeyFile:          packet.SSHKey,
// 		SSHUser:             "root",
// 	}
//
// 	f, fileErr := os.Create("kismatic-testing.yaml")
// 	FailIfError(fileErr, "Error waiting for nodes")
// 	defer f.Close()
// 	w := bufio.NewWriter(f)
// 	execErr := template.Execute(w, &nodes)
// 	FailIfError(execErr, "Error filling in plan template")
// 	w.Flush()
//
// 	By("Validing our plan")
// 	ver := exec.Command("./kismatic", "install", "validate", "-f", f.Name())
// 	verbytes, verErr := ver.CombinedOutput()
// 	verText := string(verbytes)
//
// 	FailIfError(verErr, "Error validating plan", verText)
//
// 	if bailBeforeAnsible() {
// 		return
// 	}
//
// 	By("Punch it Chewie!")
// 	app := exec.Command("./kismatic", "install", "apply", "-f", f.Name())
// 	app.Stdout = os.Stdout
// 	app.Stderr = os.Stderr
// 	appErr := app.Run()
//
// 	FailIfError(appErr, "Error applying plan")
// }
