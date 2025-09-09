variable "region" {
  type    = string
  default = "us-east-2"
}

variable "bucket_name" {
  description = "Name of the S3 bucket for the dashboard"
  type        = string
  default     = "debate-chatbot-dashboard"
}

variable "distribution_name" {
  description = "Name of the CloudFront distribution"
  type        = string
  default     = "debate-chatbot-dashboard"
}
