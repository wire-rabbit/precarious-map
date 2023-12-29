# Setting Up Example Servers

**Important:** Depending on your AWS account's free-tier status, running this example **MAY COST MONEY**.

It spins up two small (`t2.micro`) Ubuntu Linux instances on EC2 in the region you specify. Once you're finished, be sure to run:
```bash
terraform destroy
```
To remove the infrastructure and minimize any possible expense.

## Prerequisites

* Install [Terraform](https://www.terraform.io)
* Install the [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) tools
* [Create (or select) a profile](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/configure/index.html) to use for the example

## Using the example:

Once you have everything in place:
```bash
# Move into this directory:
cd example-setup

# Initialize Terraform (this will download the necessary AWS provider plugin):
terraform init

# View the plan Terraform will use to deploy the infrastructure:
terraform plan
```
This step will ask you two questions:
* "AWS profile to use for authentication" - enter the name of the AWS profile you created
* "This example may incur costs on AWS. Are you sure?" - type `yes` to proceed (*CAREFUL!* Deploying infrastructure on AWS is rarely free of charge!) 

If the output looks acceptable, you can deploy the two servers with:
```bash
terraform apply
```
It will ask the same two questions and then actually deploy the infrastructure to your account.

Running `precarious-map` once the deploy has completed, for the `us-east-1` region, will show the two servers (along with others you may have in that region).

Once you're done testing, remember to remove the added servers with:
```bash
terraform destroy
```

