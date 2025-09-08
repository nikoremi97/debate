
locals {
  public_tags = merge(var.tags, { Name = "public" })
}
// internet gateway
resource "aws_internet_gateway" "internet_gw" {
  count  = local.create_public_subnets
  vpc_id = aws_vpc.basic.id
  tags   = var.tags
}

// subnets
resource "aws_subnet" "public" {
  count                   = local.public_subnets_length
  vpc_id                  = aws_vpc.basic.id
  cidr_block              = var.public_subnets.cidr_block[count.index]
  map_public_ip_on_launch = var.map_public_ip_on_launch
  availability_zone       = var.public_subnets.zone[count.index]
  tags                    = local.public_tags
}

// endpoint
resource "aws_vpc_endpoint" "public_endpoint" {
  count = local.public_endpoint_length

  vpc_id       = aws_vpc.basic.id
  service_name = var.public_subnets.endpoint[count.index]
  tags         = local.public_tags
}

// route table
resource "aws_route_table" "public_route_table" {
  count = local.create_public_subnets

  vpc_id = aws_vpc.basic.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.internet_gw[0].id
  }

  tags = local.public_tags
}

// route table association
resource "aws_route_table_association" "public_rta" {
  count = local.public_subnets_length

  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public_route_table[0].id
}

resource "aws_vpc_endpoint_route_table_association" "public_endpoint" {
  count = local.public_endpoint_length

  vpc_endpoint_id = aws_vpc_endpoint.public_endpoint[count.index].id
  route_table_id  = aws_route_table.public_route_table[0].id
}

