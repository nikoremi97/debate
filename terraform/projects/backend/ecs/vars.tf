variable "vpc_id" {
  description = "VPC ID where resources will be created"
  type        = string
}

variable "public_subnet_ids" {
  description = "List of public subnet IDs for the load balancer"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "List of private subnet IDs for ECS tasks"
  type        = list(string)
}

variable "ecr_repository_url" {
  description = "ECR repository URL for the container image"
  type        = string
}

variable "openai_api_key_secret_arn" {
  description = "ARN of the OpenAI API key secret"
  type        = string
}

variable "database_url_secret_arn" {
  description = "ARN of the database URL secret"
  type        = string
}

variable "api_key_secret_arn" {
  description = "ARN of the API key secret"
  type        = string
}

variable "debate_chatbot_kms_key_arn" {
  description = "ARN of the debate chatbot KMS key"
  type        = string
}

variable "debate_database_kms_key_arn" {
  description = "ARN of the debate database KMS key"
  type        = string
}

variable "region" {
  description = "AWS region"
  type        = string
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}

variable "alb_security_group_id" {
  description = "ALB security group ID from VPC project"
  type        = string
}

variable "ecs_tasks_security_group_id" {
  description = "ECS tasks security group ID from VPC project"
  type        = string
}
