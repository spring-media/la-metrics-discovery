# metrics-discovery
Can be used in a monitoring systems like nagios or zabbix to discover items on aws 

#### Installing
	go get github.com/mrsn/metrics-discovery

#### Discover ELBs
	
	metrics-discovery -aws-region eu-central-1 -discovery ELB
