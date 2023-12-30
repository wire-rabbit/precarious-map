# Precarious Map - A TUI View Into AWS EC2 Instances

Written in Go and built with the [bubbletea](https://github.com/charmbracelet/bubbletea) framework, this little TUI app gives a quick overview of the state of all EC2 instances in a given region. When selected, the full JSON from the underlying `DescribeInstances` API call is shown in another panel.

Usage:
```bash
# Include the name of the AWS profile to use for authentication:
./precarious-map --profile devs

# Optionally supply a region. Defaults to us-east-1:
./precarious-map --profile devs --region us-west-2
```

An example Terraform deployment to test the tool is included under `example-setup`. It deploys two small `t2.micro` ubuntu servers. (Note, of course, that running servers on AWS - even for a very short time - may not be free of charge.)

