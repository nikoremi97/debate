variable "alb_dns_name" {
  description = "DNS name of the Application Load Balancer"
  type        = string
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}
