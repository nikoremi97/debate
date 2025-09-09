variable "bucket_name" {
  description = "Name of the S3 bucket for the dashboard"
  type        = string
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}
