output "cluster_endpoint" {
  description = "Aurora cluster endpoint"
  value       = aws_rds_cluster.debate_aurora.endpoint
}

output "cluster_reader_endpoint" {
  description = "Aurora cluster reader endpoint"
  value       = aws_rds_cluster.debate_aurora.reader_endpoint
}

output "cluster_identifier" {
  description = "Aurora cluster identifier"
  value       = aws_rds_cluster.debate_aurora.cluster_identifier
}

output "cluster_arn" {
  description = "Aurora cluster ARN"
  value       = aws_rds_cluster.debate_aurora.arn
}

output "database_name" {
  description = "Database name"
  value       = aws_rds_cluster.debate_aurora.database_name
}

output "master_username" {
  description = "Master username"
  value       = aws_rds_cluster.debate_aurora.master_username
}

output "password_secret_arn" {
  description = "ARN of the password secret in Secrets Manager"
  value       = aws_secretsmanager_secret.db_password.arn
}

output "connection_string" {
  description = "Database connection string (without password)"
  value       = "postgresql://${aws_rds_cluster.debate_aurora.master_username}@${aws_rds_cluster.debate_aurora.endpoint}:5432/${aws_rds_cluster.debate_aurora.database_name}"
}
