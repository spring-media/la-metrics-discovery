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
	)
	flag.Parse()

	switch *discoveryType {
	case "ELB":
		err := getAllElasticLoadBalancers(*awsRegion)
		if err != nil {
			log.Printf("Could not descibe load balancers: %v", err)
		}
	case "RDS":
		err := getAllDBInstances(*awsRegion)
		if err != nil {
			log.Printf("Could not describe db instances: %v", err)
		}
	case "CloudFront":
		err := getAllCloudFrontDistributions(*awsRegion)
		if err != nil {
			log.Printf("Could not list distributions")
		}

	default:
		log.Printf("discovery type %s not supported", *discoveryType)
	}
}

func getAllDBInstances(awsRegion string) error {
	svc := rds.New(session.New(), aws.NewConfig().WithRegion(awsRegion))
	params := &rds.DescribeDBInstancesInput{}
	resp, err := svc.DescribeDBInstances(params)

	if err != nil {
		return fmt.Errorf("getting RDS instances in region %q:%v", awsRegion, err)
	}

	rdsIdentifiers := [](map[string]string){}

	for _, rds := range resp.DBInstances {
		rdsIdentifier := map[string]string{
			"{RDSIDENTIFIER}": *rds.DBInstanceIdentifier,
		}
		rdsIdentifiers = append(rdsIdentifiers, rdsIdentifier)
	}

	err = toJson(&rdsIdentifiers)

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func getAllCloudFrontDistributions(awsRegion string) error {
	svc := cloudfront.New(session.New(), aws.NewConfig().WithRegion(awsRegion))

	params := &cloudfront.ListDistributionsInput{}

	resp, err := svc.ListDistributions(params)

	if err != nil {
		return fmt.Errorf("listing CloudFront distributions %v", err)
	}

	dists := [](map[string]string){}

	for _, dist := range resp.DistributionList.Items {
		distId := map[string]string{
			"{DISTID}": *dist.Id,
		}

		dists = append(dists, distId)
	}

	err = toJson(&dists)

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func getAllElasticLoadBalancers(awsRegion string) error {
	svc := elb.New(session.New(), aws.NewConfig().WithRegion(awsRegion))
	params := &elb.DescribeLoadBalancersInput{}
	resp, err := svc.DescribeLoadBalancers(params)

	if err != nil {
		return fmt.Errorf("reading ELBs in region %q:%v", awsRegion, err)
	}

	elbs := [](map[string]string){}

	for _, elb := range resp.LoadBalancerDescriptions {
		elbName := map[string]string{
			"{#LOADBALANCERNAME}": *elb.LoadBalancerName,
		}

		elbs = append(elbs, elbName)
	}

	err = toJson(&elbs)

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func toJson(descriptions *[]map[string]string) error {
	res := Result{Data: descriptions}
	b, err := json.Marshal(res)

	if err != nil {
		return fmt.Errorf("error marshaling", err)
	}

	os.Stdout.Write(b)
	return nil
}
