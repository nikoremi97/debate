resource "aws_kms_key" "debate_secrets" {
  description             = "KMS key for Debate Chatbot secrets encryption"
  deletion_window_in_days = 7
  enable_key_rotation     = true

  tags = var.tags
}

resource "aws_kms_alias" "debate_secrets" {
  name          = "alias/debate-chatbot-secrets"
  target_key_id = aws_kms_key.debate_secrets.key_id
}
