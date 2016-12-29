package integration

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
)

func testGlusterCluster(aws infrastructureProvisioner, distro linuxDistro) {
	WithInfrastructure(NodeCount{Worker: 2}, distro, aws, func(nodes provisionedNodes, sshKey string) {
		By("Setting up a plan file with storage nodes")
		plan := PlanAWS{
			Etcd:                     nodes.worker,
			Master:                   nodes.worker,
			Worker:                   nodes.worker,
			Storage:                  nodes.worker,
			MasterNodeFQDN:           nodes.worker[0].Hostname,
			MasterNodeShortName:      nodes.worker[0].Hostname,
			AllowPackageInstallation: true,
			SSHKeyFile:               sshKey,
			SSHUser:                  nodes.worker[0].SSHUser,
		}

		By("Writing plan file out to disk")
		template, err := template.New("planAWSOverlay").Parse(planAWSOverlay)
		FailIfError(err, "Couldn't parse template")
		f, err := os.Create("kismatic-testing.yaml")
		FailIfError(err, "Error waiting for nodes")
		defer f.Close()
		err = template.Execute(f, &plan)
		FailIfError(err, "Error filling in plan template")

		if distro == Ubuntu1604LTS { // Ubuntu doesn't have python installed
			By("Running the all play with the plan")
			cmd := exec.Command("./kismatic", "install", "step", "_all.yaml", "-f", f.Name())
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			FailIfError(err, "Error running all play")
		}

		By("Running the storage play with the plan")
		cmd := exec.Command("./kismatic", "install", "step", "_storage.yaml", "-f", f.Name())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		FailIfError(err, "Error running storage play")

		By("Setting up a gluster volume on the nodes")
		err = runViaSSH([]string{"sudo mkdir -p /data/test-volume"}, nodes.worker, sshKey, 30*time.Second)
		FailIfError(err, "Error creating test volume directory")
		create := fmt.Sprintf("sudo gluster volume create gv0 replica 2 %s:/data/test-volume %s:/data/test-volume force", nodes.worker[0].Hostname, nodes.worker[1].Hostname)
		cmds := []string{create, "sudo gluster volume start gv0", "sudo gluster volume info"}
		err = runViaSSH(cmds, nodes.worker[0:1], sshKey, 1*time.Minute)
		FailIfError(err, "Error creating gluster volume")

		By("Mounting the volume on one of the nodes, and writing a file")
		mount := fmt.Sprintf("sudo mount -t glusterfs %s:/gv0 /mnt", nodes.worker[0].Hostname)
		err = runViaSSH([]string{mount, "sudo touch /mnt/test-file"}, nodes.worker[0:1], sshKey, 30*time.Second)
		FailIfError(err, "Error mounting gluster volume")

		time.Sleep(3 * time.Second)
		By("Verifying file is on the other node")
		err = runViaSSH([]string{"sudo cat /data/test-volume/test-file"}, nodes.worker[1:2], sshKey, 30*time.Second)
		FailIfError(err, "Error verifying that the test file is in the gluster volume")
	})
}
