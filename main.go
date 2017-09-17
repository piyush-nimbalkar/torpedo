package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println("vim-go")
	//List aws ec2 instances
	//
	sess := session.Must(session.NewSession())
	nameFilter := os.Args[1]
	awsRegion := "us-east-1"
	config := &aws.Config{Region: aws.String(awsRegion)}
	config.WithCredentialsChainVerboseErrors(true)

	//Add credentials
	//
	creds := credentials.NewEnvCredentials()

	// Retrieve the credentials value
	credValue, err := creds.Get()
	fmt.Printf("credValue: %v\n", credValue)
	if err != nil {
		// handle error
	}
	config.WithCredentials(creds)
	svc := ec2.New(sess, config)
	fmt.Printf("svc.config: %v\n", svc.Config)
	fmt.Printf("svc.config.cred: %v", *svc.Config.Credentials)
	fmt.Printf("listing instances with tag %v in: %v\n", nameFilter, awsRegion)
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(strings.Join([]string{"*", nameFilter, "*"}, "")),
				},
			},
		},
	}
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", awsRegion, err.Error())
		log.Fatal(err.Error())
	}
	fmt.Printf("%+v\n", *resp)
}
