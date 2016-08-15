# metrics-discovery
Can be used in a monitoring systems like nagios or zabbix to discover items on aws 

#### Installing
	go get github.com/weltn24/metrics-discovery

#### Discover ELBs
	
	metrics-discovery -aws.region eu-central-1 -type ELB

#### Discover RDS Instances

	metrics-discovery -aws.region eu-central-1 -type RDS

#### Discover CloudFront Distributions

	metrics-discovery -aws.region us-east-1 -type CloudFront