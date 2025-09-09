output "repository_url" {
  description = "The URL of the ECR repository"
  value       = aws_ecr_repository.debate_api.repository_url
}

output "repository_arn" {
  description = "The ARN of the ECR repository"
  value       = aws_ecr_repository.debate_api.arn
}

output "registry_id" {
  description = "The registry ID where the repository was created"
  value       = aws_ecr_repository.debate_api.registry_id
}
