package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/efs"
)

type EFSDisk struct {
	Name         string
	IpAddress    string
	FileSystemId string
}

//NOTE: Building EFS file systems and mount points takes time and the only cost is storage.
//So rather than create a new one each time and destroy it, as we do with instances, we
//create if it's not there and reuse if it is.
func (c Client) buildEFSDisk(which string) *EFSDisk {
	disk := EFSDisk{
		Name: which,
	}
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return nil
	}

	svc := efs.New(sess)

	params := &efs.CreateFileSystemInput{
		CreationToken: aws.String(which),
	}
	_, err = svc.CreateFileSystem(params)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	fsId := blockForFileSystemId(which, svc, 30)

	if fsId == nil {
		return nil
	}
	disk.FileSystemId = *fsId

	createMountParams := &efs.CreateMountTargetInput{
		FileSystemId: aws.String(disk.FileSystemId),
		SubnetId:     aws.String(c.Config.SubnetID),
		SecurityGroups: []*string{
			aws.String(c.Config.SecurityGroupID),
		},
	}
	createMountResp, err := svc.CreateMountTarget(createMountParams)

	if err != nil {
		return nil
	}

	ipAdd := blockForMountTargetIpAddress(*createMountResp.MountTargetId, svc, 30)

	if ipAdd == nil {
		return nil
	}

	disk.IpAddress = *ipAdd

	return &disk
}

func blockForFileSystemId(which string, svc *efs.EFS, maxTimeOutSecs int) *string {
	if maxTimeOutSecs == 0 {
		return nil
	}
	descParams := &efs.DescribeFileSystemsInput{
		CreationToken: aws.String(which),
	}

	desResp, err := svc.DescribeFileSystems(descParams)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	if len(desResp.FileSystems) > 0 && *desResp.FileSystems[0].LifeCycleState == "created" {
		return desResp.FileSystems[0].FileSystemId
	}
	fmt.Print("#")
	time.Sleep(1 * time.Second)
	return blockForFileSystemId(which, svc, maxTimeOutSecs-1)
}

func blockForMountTargetIpAddress(mountTarget string, svc *efs.EFS, maxTimeOutSecs int) *string {
	if maxTimeOutSecs == 0 {
		return nil
	}
	descParams := &efs.DescribeMountTargetsInput{
		MountTargetId: aws.String(mountTarget),
	}

	desResp, err := svc.DescribeMountTargets(descParams)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	if len(desResp.MountTargets) > 0 && *desResp.MountTargets[0].LifeCycleState == "available" {
		return desResp.MountTargets[0].IpAddress
	}
	fmt.Print("@")
	time.Sleep(1 * time.Second)
	return blockForMountTargetIpAddress(mountTarget, svc, maxTimeOutSecs-1)
}
