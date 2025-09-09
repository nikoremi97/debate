terraform {
  backend "s3" {
    bucket         = "terraform-state-debate-chatbot"
    region         = "us-east-2"
    key            = "terraform/projects/frontend/terraform.tfstate"
    dynamodb_table = "terraform-state"
  }
}
