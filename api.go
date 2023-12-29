package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

var awsClient *ec2.Client

func getAwsClient(options AppOptions) *ec2.Client {
	if options.Profile != "" {
		cfg, err := config.LoadDefaultConfig(
			context.Background(),
			config.WithSharedConfigProfile("iamadmin-general"),
		)
		if err != nil {
			log.Fatalf("error setting up AWS client: %s", err.Error())
		}
		return ec2.NewFromConfig(cfg)
	}
	// No profile name was provided - do our best:
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
	)
	if err != nil {
		log.Fatalf("error setting up AWS client: %s", err.Error())
	}
	return ec2.NewFromConfig(cfg)
}

func processEc2Data(data *ec2.DescribeInstancesOutput) []InstanceDetail {

	var details = []InstanceDetail{}

	for _, reservation := range data.Reservations {
		for _, instance := range reservation.Instances {
			// try to get the name of the instance:
			name := "<empty>"
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}
			// get the remaining details:
			details = append(details, InstanceDetail{
				InstanceId:       getTableValue(instance.InstanceId),
				ImageId:          getTableValue(instance.ImageId),
				AZ:               getTableValue(instance.Placement.AvailabilityZone),
				State:            getTableValue((*string)(&instance.State.Name)),
				PrivateIpAddress: getTableValue(instance.PrivateIpAddress),
				PublicIpAddress:  getTableValue(instance.PublicIpAddress),
				Name:             name,
			})
		}
	}
	return details
}

func getTableValue(raw *string) string {
	if raw == nil {
		return "<empty>"
	}
	return *raw
}

func fetchEc2Data() []InstanceDetail {
	if awsClient == nil {
		log.Fatal("AWS client was not intialized prior to fetch.")
	}
	output, err := awsClient.DescribeInstances(
		context.TODO(),
		&ec2.DescribeInstancesInput{},
		func(o *ec2.Options) {
			o.Region = Options.Region
		},
	)
	if err != nil {
		log.Fatalf("error fetching data from AWS: %s", err.Error())
	}
	return processEc2Data(output)
}
