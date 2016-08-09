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
)

type Result struct {
	Data interface{} `json:"data"`
}

func main() {
	var (
		discovery = flag.String("discovery", "", "type of discovery. Only ELB supported right now")
		awsRegion = flag.String("aws.region", "eu-central-1", "AWS region")
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

	elbs := []string{}

	for _, elb := range resp.LoadBalancerDescriptions {
		elbs = append(elbs, "{#LOADBALANCERNAME}:"+(*elb.LoadBalancerName))
	}

	r := Result{Data: elbs}
	b, err := json.Marshal(r)

	if err != nil {
		return fmt.Errorf("error marshaling", err)
	}

	os.Stdout.Write(b)
	return nil
}
