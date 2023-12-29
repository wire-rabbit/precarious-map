terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.31.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.6.0"
    }
  }
}

provider "aws" {
  region  = var.aws_region
  profile = var.aws_profile
}

variable "confirm_possible_expense" {
  type        = string
  description = "This example may incur costs on AWS. Are you sure?"
}

variable "aws_profile" {
  type        = string
  description = "AWS profile to use for authentication"
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

resource "random_string" "unique_part" {
  length  = 5
  special = false
  upper   = false
}

data "aws_ami" "ubuntu_22_04" {
  most_recent = true
  owners      = ["amazon"]
  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }
}

resource "aws_instance" "web_server" {
  count         = 2
  ami           = data.aws_ami.ubuntu_22_04.id
  instance_type = "t2.micro"

  tags = {
    Name      = "demo_ec2_${random_string.unique_part.result}_${count.index}"
    Terraform = "true"
  }

  lifecycle {
    precondition {
      condition     = var.confirm_possible_expense == "yes"
      error_message = "Use of AWS does incur costs. Confirmation not supplied. Exiting."
    }
  }
}


