output "bucket_name" {
  description = "Name of the S3 bucket"
  value       = aws_s3_bucket.dashboard.bucket
}

output "bucket_arn" {
  description = "ARN of the S3 bucket"
  value       = aws_s3_bucket.dashboard.arn
}

output "bucket_domain_name" {
  description = "Domain name of the S3 bucket"
  value       = aws_s3_bucket.dashboard.bucket_domain_name
}

output "bucket_website_endpoint" {
  description = "Website endpoint of the S3 bucket"
  value       = aws_s3_bucket_website_configuration.dashboard.website_endpoint
}

output "bucket_website_domain" {
  description = "Website domain of the S3 bucket"
  value       = aws_s3_bucket_website_configuration.dashboard.website_domain
}
