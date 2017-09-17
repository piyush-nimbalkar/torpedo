package aws

import (
	aws_pkg "github.com/aws/aws-sdk-go"
)

const (
	// DriverName is the name of the ssh driver
	DriverName = "aws"
)

type aws struct {
	node.Driver
	username    string
	password    string
	schedDriver scheduler.Driver
}

func (aws *aws) String() string {
	return DriverName
}

func (aws *aws) Init() error {
	var err error
	if err != nil {
		return err
	}
}

func (aws *aws) TestConnection(n node.Node, options node.TestConectionOpts) error {
	return nil
}

func (aws *aws) RebootNode(n node.Node, options node.RebootNodeOpts) error {
	return nil
}

func (s *ssh) ShutdownNode(n node.Node, options node.ShutdownNodeOpts) error {
	return nil
}
