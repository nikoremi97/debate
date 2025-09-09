data "aws_kms_secrets" "open_ai_api_key" {
  secret {
    name    = "openai-secret-api-key"
    payload = "AQICAHil1RvAZbITWyvQUuxduerGX6Hb+a8Lw5dOFuyOY6+zXAGCgMBHwKBQOsitE0Jc+l20AAABCDCCAQQGCSqGSIb3DQEHBqCB9jCB8wIBADCB7QYJKoZIhvcNAQcBMB4GCWCGSAFlAwQBLjARBAx+uvnIWrQfVLE21eMCARCAgb/tQBLgitT/FrzJ8Y76JFzCMLqYI+txuZ9gCoZJZKxeLuZKCMe+iOYYMwfJKFrI9Er4s93gCqw6/i0oZifx7A5qBbF4M6PoGJwQxf+Plp2HzISaJ1huD4jZgnmKXjpGgNozo/2ZW5eu/v2LanPKjr57x8VgSDVLE+4PBMmtW2TO/MEulC0FOdNxQ0DCr6crNHSqXu2ckKTgPwimK/jtV0AF5w5BR5i+lcOnm31Zy8+7jKDA/z847yefTngk5fzJoQ=="
  }
}

resource "aws_secretsmanager_secret" "open_ai_api_key" {
  name        = "openai-secret-api-key"
  description = "OpenAI API key"
  tags        = var.tags
}

resource "aws_secretsmanager_secret_version" "open_ai_api_key" {
  secret_id     = aws_secretsmanager_secret.open_ai_api_key.id
  secret_string = data.aws_kms_secrets.open_ai_api_key.plaintext["openai-secret-api-key"]
}
