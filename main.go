package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type AppOptions struct {
	Help    bool
	Region  string
	Profile string
}

var Options AppOptions

const helpText = `
  precarious-map: A simple TUI for inspecting AWS EC2 instances.

  Options:
  * --help    - show this text and exit
  * --region  - set the AWS region (defaults to 'us-east-1')
  * --profile - set the AWS CLI profile to use for authentication
`

func HandleOptions(writer io.Writer, help bool, region string, profile string) AppOptions {

	if help {
		fmt.Fprintf(writer, helpText)
		return AppOptions{Help: true}
	}

	return AppOptions{
		Help:    help,
		Region:  region,
		Profile: profile,
	}
}

func main() {
	help := flag.Bool("help", false, "show help")
	region := flag.String("region", "us-east-1", "set the AWS region")
	profile := flag.String("profile", "", "set the AWS CLI profile to use for authentication")
	flag.Parse()

	Options = HandleOptions(os.Stdout, *help, *region, *profile)
	if Options.Help {
		os.Exit(0)
	}
}
