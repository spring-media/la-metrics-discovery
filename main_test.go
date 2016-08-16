package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/rds"
)

type mockRdsCli struct {
	DescribeDBInstancesFn func(*rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error)
}

func (m *mockRdsCli) DescribeDBInstances(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
	return m.DescribeDBInstancesFn(input)
}

func TestGetAllDBInstances(t *testing.T) {
	testCases := map[string]struct {
		returnedInstances []*rds.DBInstance
		expectedSize      int
	}{
		"nilReturn": {
			returnedInstances: nil,
			expectedSize:      0,
		},
		"emptyReturn": {
			returnedInstances: []*rds.DBInstance{},
			expectedSize:      0,
		},
	}

	for testName, testCase := range testCases {
		dbInstances, err := getAllDBInstances(&mockRdsCli{
			DescribeDBInstancesFn: func(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
				return &rds.DescribeDBInstancesOutput{
					DBInstances: testCase.returnedInstances,
				}, nil
			},
		})
		if err != nil {
			t.Errorf("%s %v", testName, err)
		}
		if len(dbInstances) != testCase.expectedSize {
			t.Errorf("%s expected %d instances, got %d", testName, testCase.expectedSize, len(dbInstances))
		}
	}

}
