data "aws_kms_secrets" "api_key" {
  secret {
    name    = "debate-chatbot-api-key"
    payload = "AQICAHil1RvAZbITWyvQUuxduerGX6Hb+a8Lw5dOFuyOY6+zXAHxbausDa7pMRIohYG+F+yhAAAAcTBvBgkqhkiG9w0BBwagYjBgAgEAMFsGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQM85dzzQWoQPYrtvthAgEQgC6NvlrB26bMP0LV8w+4f2rMKzQ5bbCruA+ktFKdI5pE5AAO7Dptrq5oSSBzKIH1"
  }
}

# API Key Secret for authentication
resource "aws_secretsmanager_secret" "api_key" {
  name        = "debate-chatbot-api-key"
  description = "API key for debate chatbot authentication"
  tags        = var.tags
}

resource "aws_secretsmanager_secret_version" "api_key" {
  secret_id     = aws_secretsmanager_secret.api_key.id
  secret_string = data.aws_kms_secrets.api_key.plaintext["debate-chatbot-api-key"]
}
