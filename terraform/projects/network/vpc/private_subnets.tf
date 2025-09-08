locals {
  private_tags = merge(var.tags, { Name = "private" })
}

resource "aws_eip" "nat" {
  count  = var.use_eip ? 1 : 0
  domain = "vpc"
  tags   = var.tags
}

resource "aws_nat_gateway" "gw" {
  count = local.create_public_subnets

  allocation_id = aws_eip.nat[0].id
  subnet_id     = aws_subnet.public[0].id
  tags          = var.tags
}

// subnets
resource "aws_subnet" "private" {
  count             = local.private_subnets_length
  vpc_id            = aws_vpc.basic.id
  cidr_block        = var.private_subnets.cidr_block[count.index]
  availability_zone = var.private_subnets.zone[count.index]
  tags              = local.private_tags
}

// endpoint
resource "aws_vpc_endpoint" "private_endpoint" {
  count = local.private_endpoint_length

  vpc_id       = aws_vpc.basic.id
  service_name = var.private_subnets.endpoint[count.index]
  tags         = local.private_tags
}

// route table
resource "aws_route_table" "private_route_table" {
  count = local.create_public_subnets

  vpc_id = aws_vpc.basic.id

  # Default route to NAT Gateway
  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.gw[0].id
  }

  tags = local.private_tags
}

// route table association
resource "aws_route_table_association" "private_rta" {
  count = local.private_subnets_length

  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private_route_table[0].id
}

resource "aws_vpc_endpoint_route_table_association" "private_endpoint" {
  count = local.private_endpoint_length

  vpc_endpoint_id = aws_vpc_endpoint.private_endpoint[count.index].id
  route_table_id  = aws_route_table.private_route_table[0].id
}
