package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type AppOptions struct {
	ShowHelp bool
	Region   string
	Profile  string
}

var Options AppOptions

const helpText = `
  precarious-map: A simple TUI for inspecting AWS EC2 instances.

  Options:
  * --help    - show this text and exit
  * --region  - set the AWS region (defaults to 'us-east-1')
  * --profile - set the AWS CLI profile to use for authentication
`

func HandleOptions(writer io.Writer, showHelp bool, region string, profile string) AppOptions {

	if showHelp {
		fmt.Fprintf(writer, helpText)
	}

	return AppOptions{
		// ...
	}
}

func main() {
	showHelp := flag.Bool("help", false, "show help")
	region := flag.String("region", "us-east-1", "set the AWS region")
	profile := flag.String("profile", "", "set the AWS CLI profile to use for authentication")
	flag.Parse()

	Options = HandleOptions(os.Stdout, *showHelp, *region, *profile)
	if Options.ShowHelp {
		os.Exit(0)
	}
}
