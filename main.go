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
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/rds"
)

type Result struct {
	Data interface{} `json:"data"`
}

func main() {
	var (
		discoveryType = flag.String("type", "", "type of discovery. Only ELB, RDS and CloudFront supported right now")
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
	default:
		log.Fatalf("discovery type %s not supported", *discoveryType)
	}

	err = json.NewEncoder(os.Stdout).Encode(Result{Data: list})

	if err != nil {
		log.Fatal(err)
	}
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
			"{#DISTID}": *dist.Id,
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
