package main

import (
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
)

func main() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := ec2.New(sess, &aws.Config{Region: aws.String("ap-northeast-1")})

	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	var data [][]string
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{
		"Environment", "Role", "Name", "InstanceId", "InstanceType", "AZ",
		"PrivateIP", "PublicIP", "Status"})
	for _, r := range resp.Reservations {
		for _, i := range r.Instances {
			// fmt.Println(i)
			var tagEnvironment, tagName, tagRole string
			for _, t := range i.Tags {
				if *t.Key == "Environment" {
					tagEnvironment = *t.Value
				}
				if *t.Key == "Role" {
					tagRole = *t.Value
				}
				if *t.Key == "Name" {
					tagName = *t.Value
				}
			}
			var publicIpAddress string
			if i.PublicIpAddress == nil {
				publicIpAddress = "-"
			} else {
				publicIpAddress = *i.PublicIpAddress
			}
			data = append(data, []string{
				tagEnvironment, tagRole, tagName,
				*i.InstanceId, *i.InstanceType, *i.Placement.AvailabilityZone,
				*i.PrivateIpAddress, publicIpAddress, *i.State.Name})
		}
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i][0] < data[j][0]
	})
	table.AppendBulk(data)
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.Render()
}
