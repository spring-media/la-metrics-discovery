package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/rds"
)

type Result struct {
	Data interface{} `json:"data"`
}

func main() {
	var (
		discoveryType = flag.String("type", "", "type of discovery. EC2, ELB, RDS, CloudFront or Lambda")
		awsRegion     = flag.String("aws.region", "eu-central-1", "AWS region")
		list          interface{}
		err           error
	)

	flag.Parse()

	awsSession := session.New(aws.NewConfig().WithRegion(*awsRegion))

	switch *discoveryType {
	case "ELB":
		list, err = getAllElasticLoadBalancers(elb.New(awsSession))
		if err != nil {
			log.Fatalf("Could not descibe load balancers: %v", err)
		}
	case "EC2":
		list, err = getAllEC2Instances(ec2.New(awsSession))
		if err != nil {
			log.Fatalf("Could not get ec2 instances: %v", err)
		}
	case "RDS":
		list, err = getAllDBInstances(rds.New(awsSession))
		if err != nil {
			log.Fatalf("Could not describe db instances: %v", err)
		}
	case "CloudFront":
		list, err = getAllCloudFrontDistributions(cloudfront.New(awsSession))
		if err != nil {
			log.Fatalf("Could not list distributions")
		}
	case "Lambda":
		list, err = getAllLambdas(lambda.New(awsSession))
		if err != nil {
			log.Fatalf("Could not list distributions")
		}
	default:
		log.Fatalf("discovery type %s not supported", *discoveryType)
	}

	err = json.NewEncoder(os.Stdout).Encode(Result{Data: list})

	if err != nil {
		log.Fatal(err)
	}
}

func getAllEC2Instances(ec2Cli interface {
	DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
}) ([]map[string]string, error) {

	resp, err := ec2Cli.DescribeInstances(&ec2.DescribeInstancesInput{})

	if err != nil {
		return nil, fmt.Errorf("getting EC2 instances: %v", err)
	}

	ec2Identifiers := make([]map[string]string, 0, len(resp.Reservations))

	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			var name string
			for _, t := range instance.Tags {
				if *t.Key == "Name" {
					name = *t.Value
					break
				}
			}
			ec2Identifiers = append(ec2Identifiers, map[string]string{
				"{#INSTANCEID}":   *instance.InstanceId,
				"{#INSTANCENAME}": name,
			})
		}

	}
	return ec2Identifiers, nil
}

func getAllDBInstances(rdsCli interface {
	DescribeDBInstances(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error)
}) ([]map[string]string, error) {
	resp, err := rdsCli.DescribeDBInstances(&rds.DescribeDBInstancesInput{})

	if err != nil {
		return nil, fmt.Errorf("getting RDS instances:%v", err)
	}

	rdsIdentifiers := make([]map[string]string, len(resp.DBInstances))

	for ctr, rds := range resp.DBInstances {

		rdsIdentifiers[ctr] = map[string]string{
			"{#RDSIDENTIFIER}": *rds.DBInstanceIdentifier,
		}
	}
	return rdsIdentifiers, nil
}

func getAllCloudFrontDistributions(cloudFrontCli interface {
	ListDistributions(*cloudfront.ListDistributionsInput) (*cloudfront.ListDistributionsOutput, error)
}) ([]map[string]string, error) {

	resp, err := cloudFrontCli.ListDistributions(&cloudfront.ListDistributionsInput{})

	if err != nil {
		return nil, fmt.Errorf("listing CloudFront distributions %v", err)
	}

	dists := make([]map[string]string, len(resp.DistributionList.Items))

	for ctr, dist := range resp.DistributionList.Items {
		dists[ctr] = map[string]string{
			"{#DISTID}":    *dist.Id,
			"{#DISTALIAS}": *dist.Aliases.Items[0],
		}
	}

	return dists, nil
}

func getAllElasticLoadBalancers(elbCli interface {
	DescribeLoadBalancers(*elb.DescribeLoadBalancersInput) (*elb.DescribeLoadBalancersOutput, error)
}) ([]map[string]string, error) {

	resp, err := elbCli.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{})

	if err != nil {
		return nil, fmt.Errorf("reading ELBs:%v", err)
	}

	elbs := make([]map[string]string, len(resp.LoadBalancerDescriptions))

	for ctr, elb := range resp.LoadBalancerDescriptions {
		elbs[ctr] = map[string]string{
			"{#LOADBALANCERNAME}": *elb.LoadBalancerName,
		}
	}

	return elbs, nil
}

func getAllLambdas(lambdaCli interface {
	ListFunctions(*lambda.ListFunctionsInput) (*lambda.ListFunctionsOutput, error)
}) ([]map[string]string, error) {

	resp, err := lambdaCli.ListFunctions(&lambda.ListFunctionsInput{})

	if err != nil {
		return nil, fmt.Errorf("listing lambdas %v", err)
	}

	lambdas := make([]map[string]string, len(resp.Functions))

	for ctr, lambda := range resp.Functions {
		lambdas[ctr] = map[string]string{
			"{#FUNCTIONNAME}": *lambda.FunctionName,
		}
	}

	return lambdas, nil
}
