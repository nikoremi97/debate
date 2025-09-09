variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-2"
}

variable "api_key" {
  description = "API key for authentication"
  type        = string
  sensitive   = true
}


