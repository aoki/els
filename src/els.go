package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
)

var (
	Version  string
	Revision string
)

var (
	v        *bool
	noHeader *bool
)

func parseFlags() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
 %s [OPTIONS] ARGS...
Options
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	v = flag.Bool("v", false, "Display the command version")
	noHeader = flag.Bool("no-header", false, "Hide the header")
	flag.Parse()
}

func main() {
	parseFlags()
	if *v {
		fmt.Println(Version)
		os.Exit(0)
	}
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
	if !*noHeader {
		table.SetHeader([]string{
			"Environment", "Role", "Name", "InstanceId", "InstanceType", "AZ",
			"PrivateIP", "PublicIP", "Status"})
	}
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
			var publicIpAddress, privateIpAddress string
			if i.PublicIpAddress == nil {
				publicIpAddress = "-"
			} else {
				publicIpAddress = *i.PublicIpAddress
			}
			if i.PrivateIpAddress == nil {
				privateIpAddress = "-"
			} else {
				privateIpAddress = *i.PrivateIpAddress
			}
			data = append(data, []string{
				tagEnvironment, tagRole, tagName,
				*i.InstanceId, *i.InstanceType, *i.Placement.AvailabilityZone,
				privateIpAddress, publicIpAddress, *i.State.Name})
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
