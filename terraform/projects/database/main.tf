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
    Project = "debate-chatbot"
    Service = "database"
  }
}

# Dependencies
module "dependencies" {
  source = "./dependencies"
}

# KMS Key for database encryption
resource "aws_kms_key" "database" {
  description             = "KMS key for database encryption"
  deletion_window_in_days = 7
  enable_key_rotation     = true

  tags = local.tags
}

resource "aws_kms_alias" "database" {
  name          = "alias/debate-chatbot-database"
  target_key_id = aws_kms_key.database.key_id
}

# Random auth token for Valkey
resource "random_password" "valkey_auth_token" {
  length  = 32
  special = false
}

module "database" {
  source = "./database"

  vpc_id             = module.dependencies.vpc.vpc.id
  private_subnet_ids = module.dependencies.vpc.private_subnets_ids
  kms_key            = aws_kms_key.database
  tags               = local.tags
  vpc                = module.dependencies.vpc
}

module "redis" {
  source = "./redis"

  vpc_id             = module.dependencies.vpc.vpc.id
  private_subnet_ids = module.dependencies.vpc.private_subnets_ids
  kms_key_id         = aws_kms_key.database.key_id
  valkey_auth_token  = random_password.valkey_auth_token.result
  tags               = local.tags
  vpc                = module.dependencies.vpc
}
