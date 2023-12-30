package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

var awsClient *ec2.Client

func getAwsClient(options AppOptions) *ec2.Client {
	if options.Profile != "" {
		cfg, err := config.LoadDefaultConfig(
			context.Background(),
			config.WithSharedConfigProfile(options.Profile),
		)
		if err != nil {
			fmt.Printf("error setting up AWS client: %s\n", err)
			os.Exit(1)
		}
		return ec2.NewFromConfig(cfg)
	}
	// No profile name was provided - do our best:
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
	)
	if err != nil {
		fmt.Printf("error setting up AWS client: %s\n", err)
		os.Exit(1)
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

			instanceJson, err := json.MarshalIndent(instance, "", "  ")
			if err != nil {
				fmt.Printf("error marshaling data: %s\n", err)
				os.Exit(1)
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
				JSON:             string(instanceJson),
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
		fmt.Println("AWS client was not intialized prior to fetch.")
		os.Exit(1)
	}
	output, err := awsClient.DescribeInstances(
		context.TODO(),
		&ec2.DescribeInstancesInput{},
		func(o *ec2.Options) {
			o.Region = Options.Region
		},
	)
	if err != nil {
		fmt.Printf("error fetching data from AWS: %s\n", err)
		os.Exit(1)
	}
	return processEc2Data(output)
}
