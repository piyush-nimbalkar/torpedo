package aws

import (
	"fmt"
	"github.com/portworx/torpedo/drivers/node"
	//"github.com/portworx/torpedo/drivers/scheduler"
	aws_pkg "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"strings"
)

const (
	// DriverName is the name of the aws driver
	DriverName = "aws"
	// AwsRegion is the region of aws ec2
	AwsRegion = "us-east-1"
)

type aws struct {
	node.Driver
	session     *session.Session
	credentials *credentials.Credentials
	config      *aws_pkg.Config
	svc         *ec2.EC2
	Instances   []*ec2.Instance
}

func (aws *aws) String() string {
	return DriverName
}

func (aws *aws) Init(sched string) error {
	sess := session.Must(session.NewSession())
	aws.session = sess
	creds := credentials.NewEnvCredentials()
	aws.credentials = creds
	config := &aws_pkg.Config{Region: aws_pkg.String(AwsRegion)}
	config.WithCredentials(creds)
	aws.config = config
	svc := ec2.New(sess, config)
	aws.svc = svc
	return nil
}

func (aws *aws) TestConnection(n node.Node, options node.TestConectionOpts) error {
	return nil
}

func (aws *aws) RebootNode(n node.Node, options node.RebootNodeOpts) error {
	addr, err := aws.getAddrToConnect(n)
	if err != nil {
		return &ErrFailedToRebootNode{
			Node:  n,
			Cause: fmt.Sprintf("failed to get node address due to: %v", err),
		}
	}
	fmt.Printf("addr: %s\n", addr)
	fmt.Printf("\naws rebootnode private IP: %s", n.Addresses[0])

	instances := aws.getAllInstancesInEast()
	aws.Instances = instances
	fmt.Printf("Get all instances: +%v\n", aws.Instances)

	instanceID := "i-034154fd986ba74b6"
	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws_pkg.String(strings.Join([]string{instanceID}, "")),
		},
	}
	resp, err := aws.svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances", instanceID, err.Error())
		log.Fatal(err.Error())
	}
	reservations := resp.Reservations
	fmt.Printf("%+v\n", reservations)
	for i, resv := range reservations {
		fmt.Printf("reservation %d: %+v\n", i, resv)
	}
	for i, ins := range reservations[0].Instances {
		fmt.Printf("\ninstance %d: %+v\n", i, ins)
		fmt.Printf("\ninstance %d id: %s\n", i, *(ins.InstanceId))
		fmt.Printf("\ninstance %d PublicIpAddress: %s\n", i, *(ins.PublicIpAddress))
		fmt.Printf("\ninstance %d PrivateIpAddress: %s\n", i, *(ins.PrivateIpAddress))
	}
	return nil
}

func (aws *aws) ShutdownNode(n node.Node, options node.ShutdownNodeOpts) error {
	return nil
}

func (aws *aws) getAddrToConnect(n node.Node) (string, error) {
	if n.Addresses == nil || len(n.Addresses) == 0 {
		return "", fmt.Errorf("no address available to connect")
	}
	addr := n.Addresses[0] // TODO don't stick to first address
	return addr, nil
}

func (aws *aws) getAllInstancesInEast() []*ec2.Instance {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws_pkg.String("tag:Name"),
				Values: []*string{
					aws_pkg.String(strings.Join([]string{"*", "", "*"}, "")),
				},
			},
		},
	}
	resp, err := aws.svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", AwsRegion, err.Error())
		log.Fatal(err.Error())
	}
	reservations := resp.Reservations
	for i, resv := range reservations {
		fmt.Printf("reservation %d: %+v\n", i, resv)
	}
	//If resp is empty, then return empty []
	fmt.Printf("resp: %+v\n", *resp)
	fmt.Printf("Found %d instances in %s", len(reservations[0].Instances), AwsRegion)

	for i, ins := range reservations[0].Instances {
		fmt.Printf("\ninstance %d: %+v\n", i, ins)
		fmt.Printf("\ninstance %d id: %s\n", i, *(ins.InstanceId))
		fmt.Printf("\ninstance %d PublicIpAddress: %s\n", i, *(ins.PublicIpAddress))
		fmt.Printf("\ninstance %d PrivateIpAddress: %s\n", i, *(ins.PrivateIpAddress))
	}
	return reservations[0].Instances
}

func init() {
	a := &aws{
		Driver: node.NotSupportedDriver,
	}
	node.Register(DriverName, a)
}
