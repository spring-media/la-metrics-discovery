package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

var discovery = flag.String("discovery", "", "type of discovery. Only ELB supported right now")
var awsRegion = flag.String("aws-region", "eu-central-1", "AWS region")

func main() {
	flag.Parse()

	if *discovery == "ELB" {
		getAllElasticLoadBalancers()
	} else {
		log.Fatal("Only ELB discovery supported right now")
	}
}

func getAllElasticLoadBalancers() {
	svc := elb.New(session.New(), aws.NewConfig().WithRegion(*awsRegion))
	params := &elb.DescribeLoadBalancersInput{}
	resp, err := svc.DescribeLoadBalancers(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, elb := range resp.LoadBalancerDescriptions {
		fmt.Println(*elb.LoadBalancerName)
	}
}
