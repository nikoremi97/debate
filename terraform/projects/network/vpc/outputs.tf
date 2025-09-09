output "vpc" {
  value = aws_vpc.basic
}

output "public_route_table" {
  value = aws_route_table.public_route_table
}

output "private_route_table" {
  value = aws_route_table.private_route_table
}

output "default_network_acl" {
  value = aws_default_network_acl.default
}

output "internet_gateway" {
  value = aws_internet_gateway.internet_gw
}

output "nat_gateway" {
  value = aws_nat_gateway.gw
}

output "subnet_ids" {
  value = concat(aws_subnet.public[*].id, aws_subnet.private[*].id)
}

output "private_subnets_ids" {
  value = aws_subnet.private[*].id
}

output "public_subnets_ids" {
  value = aws_subnet.public[*].id
}

output "ecs_tasks_security_group" {
  value = aws_security_group.ecs_tasks
}

output "alb_security_group" {
  value = aws_security_group.alb
}

output "valkey_security_group" {
  value = aws_security_group.valkey
}

output "aurora_security_group" {
  value = aws_security_group.aurora
}

output "db_subnet_group" {
  value = aws_db_subnet_group.debate_db
}
