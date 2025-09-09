output "cluster_id" {
  description = "ElastiCache cluster ID"
  value       = aws_elasticache_cluster.redis.cluster_id
}

output "cluster_endpoint" {
  description = "ElastiCache cluster endpoint"
  value       = aws_elasticache_cluster.redis.cache_nodes[0].address
}

output "port" {
  description = "ElastiCache port"
  value       = aws_elasticache_cluster.redis.port
}

output "connection_string" {
  description = "Redis connection string"
  value       = "${aws_elasticache_cluster.redis.cache_nodes[0].address}:${aws_elasticache_cluster.redis.port}"
}
