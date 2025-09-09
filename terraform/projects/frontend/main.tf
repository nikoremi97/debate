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

module "s3" {
  source      = "./s3"
  bucket_name = var.bucket_name
  tags        = local.tags
}

module "cloudfront" {
  source                = "./cloudfront"
  distribution_name     = var.distribution_name
  s3_bucket_name        = module.s3.bucket_name
  s3_bucket_domain_name = module.s3.bucket_domain_name
  s3_bucket_arn         = module.s3.bucket_arn
  tags                  = local.tags
}
