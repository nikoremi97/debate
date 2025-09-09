data "terraform_remote_state" "database" {
  backend = "s3"
  config = {
    bucket = "terraform-state-debate-chatbot"
    key    = "terraform/projects/database/terraform.tfstate"
    region = var.region
  }
}

output "database" {
  value = {
    cluster_endpoint        = data.terraform_remote_state.database.outputs.database_cluster_endpoint
    connection_string       = data.terraform_remote_state.database.outputs.database_connection_string
    password_secret_arn     = data.terraform_remote_state.database.outputs.database_password_secret_arn
    redis_cluster_endpoint  = data.terraform_remote_state.database.outputs.redis_cluster_endpoint
    redis_connection_string = data.terraform_remote_state.database.outputs.redis_connection_string
    database_kms_key_arn    = data.terraform_remote_state.database.outputs.database_kms_key_arn
  }
}
