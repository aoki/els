package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

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
	v          = flag.Bool("v", false, "Display the command version")
	showStatus = flag.Bool("s", false, "Display the instance status column")
	showId     = flag.Bool("id", false, "Display the instance ID column")
	noHeader   = flag.Bool("no-header", false, "Hide the header")
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
	flag.Parse()
}

func spinner(delay time.Duration) {
	fmt.Print("\033[?25l")
	for {
		for _, r := range `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func main() {
	parseFlags()
	if *v {
		fmt.Println(Version)
		os.Exit(0)
	}

	go spinner(50 * time.Millisecond)

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

	var header []string
	header = append(header, []string{"Environment", "Role", "Name"}...)
	if *showId {
		header = append(header, "InstanceId")
	}
	header = append(header, []string{"InstanceType", "AZ", "PrivateIP", "PublicIP"}...)
	if *showStatus {
		header = append(header, "Status")
	}

	if !*noHeader {
		table.SetHeader(header)
	}

	for _, r := range resp.Reservations {
		for _, i := range r.Instances {
			var record []string

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

			record = append(record, []string{tagEnvironment, tagRole, tagName}...)
			if *showId {
				record = append(record, *i.InstanceId)
			}
			var privateIpAddress, publicIpAddress string
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
			record = append(record, []string{*i.InstanceType, *i.Placement.AvailabilityZone,
				privateIpAddress, publicIpAddress}...)
			if *showStatus {
				record = append(record, *i.State.Name)
			}

			data = append(data, record)
		}
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i][0] < data[j][0]
	})
	table.AppendBulk(data)
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	fmt.Print("\033[2K")
	table.Render()
}
