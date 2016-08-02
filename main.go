package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

func main() {
	var (
		discovery = flag.String("discovery", "", "type of discovery. Only ELB supported right now")
		awsRegion = flag.String("aws-region", "eu-central-1", "AWS region")
	)
	flag.Parse()

	switch *discovery {
	case "ELB":
		err := getAllElasticLoadBalancers(*awsRegion)
		if err != nil {
			log.Printf("Could not descibe load balancers: %v", err)
		}

	default:
		log.Printf("discovery type %s not supported", *discovery)
	}
}

func getAllElasticLoadBalancers(awsRegion string) error {
	svc := elb.New(session.New(), aws.NewConfig().WithRegion(awsRegion))
	params := &elb.DescribeLoadBalancersInput{}
	resp, err := svc.DescribeLoadBalancers(params)
	if err != nil {
		return fmt.Errorf("reading ELBs in region %q :%v", awsRegion, err)
	}

	for _, elb := range resp.LoadBalancerDescriptions {
		fmt.Println(*elb.LoadBalancerName)
	}
	return nil
}
