# Database Outputs
output "database_cluster_endpoint" {
  description = "Aurora cluster endpoint"
  value       = module.database.cluster_endpoint
}

output "database_cluster_reader_endpoint" {
  description = "Aurora cluster reader endpoint"
  value       = module.database.cluster_reader_endpoint
}

output "database_connection_string" {
  description = "Database connection string (without password)"
  value       = module.database.connection_string
}

output "database_password_secret_arn" {
  description = "ARN of the database password secret"
  value       = module.database.password_secret_arn
}

# Redis Outputs
output "redis_cluster_id" {
  description = "Redis cluster ID"
  value       = module.redis.cluster_id
}

output "redis_cluster_endpoint" {
  description = "Redis cluster endpoint"
  value       = module.redis.cluster_endpoint
}

output "redis_connection_string" {
  description = "Redis connection string"
  value       = module.redis.connection_string
}

# KMS Outputs
output "database_kms_key_id" {
  description = "Database KMS key ID"
  value       = aws_kms_key.database.key_id
}

output "database_kms_key_arn" {
  description = "Database KMS key ARN"
  value       = aws_kms_key.database.arn
}
