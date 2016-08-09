package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/rds"
)

type Result struct {
	Data interface{} `json:"data"`
}

func main() {
	var (
		discovery = flag.String("discovery", "", "type of discovery. Only ELB and RDS supported right now")
		awsRegion = flag.String("aws.region", "eu-central-1", "AWS region")
	)
	flag.Parse()

	switch *discovery {
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

	default:
		log.Printf("discovery type %s not supported", *discovery)
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

	result := Result{Data: rdsIdentifiers}
	b, err := json.Marshal(result)

	if err != nil {
		return fmt.Errorf("error marshaling", err)
	}

	os.Stdout.Write(b)
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

	r := Result{Data: elbs}
	b, err := json.Marshal(r)

	if err != nil {
		return fmt.Errorf("error marshaling", err)
	}

	os.Stdout.Write(b)
	return nil
}
