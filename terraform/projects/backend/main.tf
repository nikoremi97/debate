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
    Service = "backend"
  }
}

# Dependencies
module "dependencies" {
  source = "./dependencies"
}

# KMS Key for encryption
module "kms" {
  source = "./kms"
  tags   = local.tags
}

# ECR Repository
module "ecr" {
  source = "./ecr"
  tags   = local.tags
}

# Secrets Manager
module "secrets" {
  source = "./secrets"

  tags = local.tags
}

# ECS Infrastructure
module "ecs" {
  source = "./ecs"

  vpc_id                      = module.dependencies.vpc.vpc.id
  public_subnet_ids           = module.dependencies.vpc.public_subnets_ids
  private_subnet_ids          = module.dependencies.vpc.private_subnets_ids
  ecr_repository_url          = module.ecr.repository_url
  openai_api_key_secret_arn   = module.secrets.openai_api_key_secret_arn
  database_url_secret_arn     = module.dependencies.database.password_secret_arn
  api_key_secret_arn          = module.secrets.api_key_secret_arn
  debate_chatbot_kms_key_arn  = module.kms.key_arn
  debate_database_kms_key_arn = module.dependencies.database.database_kms_key_arn
  alb_security_group_id       = module.dependencies.vpc.alb_security_group.id
  ecs_tasks_security_group_id = module.dependencies.vpc.ecs_tasks_security_group.id
  region                      = var.region
  tags                        = local.tags
}

# CloudFront Distribution for HTTPS
module "cloudfront" {
  source = "./cloudfront"

  alb_dns_name = module.ecs.load_balancer_dns
  tags         = local.tags
}
