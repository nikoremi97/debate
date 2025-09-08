terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.90"
    }
  }
}

provider "aws" {
  region = var.region
}

locals {
  tags = {
    Environment = "development"
    Project     = "debate-chatbot"
  }
}

module "vpc" {
  source               = "./vpc"
  tags                 = local.tags
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true
  public_subnets = {
    cidr_block = ["10.0.1.0/24", "10.0.2.0/24"]
    zone       = ["us-east-2a", "us-east-2b"]
    endpoint   = []
  }
  private_subnets = {
    cidr_block = ["10.0.4.0/24", "10.0.5.0/24"]
    zone       = ["us-east-2a", "us-east-2b"]
    endpoint   = []
  }
}
