locals {
  create_public_subnets  = length(var.public_subnets.cidr_block) > 0 ? 1 : 0
  public_subnets_length  = length(var.public_subnets.cidr_block) > 0 ? length(var.public_subnets.cidr_block) : 0
  public_endpoint_length = length(var.public_subnets.endpoint) > 0 ? length(var.public_subnets.endpoint) : 0

  create_private_subnets  = length(var.private_subnets.cidr_block) > 0 ? 1 : 0
  private_subnets_length  = length(var.private_subnets.cidr_block) > 0 ? length(var.private_subnets.cidr_block) : 0
  private_endpoint_length = length(var.private_subnets.endpoint) > 0 ? length(var.private_subnets.endpoint) : 0
}

resource "aws_vpc" "basic" {
  cidr_block           = var.cidr_block
  enable_dns_support   = var.enable_dns_support
  enable_dns_hostnames = var.enable_dns_hostnames

  tags = var.tags
}

// network acl
resource "aws_default_network_acl" "default" {
  default_network_acl_id = aws_vpc.basic.default_network_acl_id

  egress {
    protocol   = var.protocol["ALL"]
    rule_no    = 100
    action     = "allow"
    cidr_block = "0.0.0.0/0"
    from_port  = 0
    to_port    = 0
  }

  ingress {
    protocol   = var.protocol["ALL"]
    rule_no    = 100
    action     = "allow"
    cidr_block = "0.0.0.0/0"
    from_port  = 0
    to_port    = 0
  }

  tags = var.tags

  subnet_ids = concat(aws_subnet.public[*].id, aws_subnet.private[*].id)
}
