data "terraform_remote_state" "vpc" {
  backend = "s3"
  config = {
    bucket = "terraform-state-debate-chatbot"
    key    = "terraform/projects/vpc/terraform.tfstate"
    region = var.region
  }
}

variable "region" {
  type    = string
  default = "us-east-2"
}

output "vpc" {
  value = data.terraform_remote_state.vpc.outputs.vpc
}
