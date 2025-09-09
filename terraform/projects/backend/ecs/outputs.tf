output "cluster_id" {
  description = "ID of the ECS cluster"
  value       = aws_ecs_cluster.debate_cluster.id
}

output "cluster_arn" {
  description = "ARN of the ECS cluster"
  value       = aws_ecs_cluster.debate_cluster.arn
}

output "service_name" {
  description = "Name of the ECS service"
  value       = aws_ecs_service.debate_api.name
}

output "service_arn" {
  description = "ARN of the ECS service"
  value       = aws_ecs_service.debate_api.id
}

output "load_balancer_dns" {
  description = "DNS name of the load balancer"
  value       = aws_lb.debate_alb.dns_name
}

output "load_balancer_zone_id" {
  description = "Zone ID of the load balancer"
  value       = aws_lb.debate_alb.zone_id
}

output "target_group_arn" {
  description = "ARN of the target group"
  value       = aws_lb_target_group.debate_api.arn
}
