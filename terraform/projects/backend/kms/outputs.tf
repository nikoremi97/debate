output "key_id" {
  description = "The globally unique identifier for the key"
  value       = aws_kms_key.debate_secrets.key_id
}

output "key_arn" {
  description = "The Amazon Resource Name (ARN) of the key"
  value       = aws_kms_key.debate_secrets.arn
}

output "alias_name" {
  description = "The display name of the alias"
  value       = aws_kms_alias.debate_secrets.name
}
