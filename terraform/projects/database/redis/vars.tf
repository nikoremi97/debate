variable "vpc_id" {
  description = "VPC ID where resources will be created"
  type        = string
}

variable "private_subnet_ids" {
  description = "List of private subnet IDs for ElastiCache"
  type        = list(string)
}


variable "kms_key_id" {
  description = "KMS key ID for encrypting ElastiCache"
  type        = string
}

variable "valkey_auth_token" {
  description = "Auth token for Valkey cluster"
  type        = string
  sensitive   = true
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}

variable "vpc" {
  description = "VPC module"
  default     = {}
}
