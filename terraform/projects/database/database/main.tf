# Random password for the database
resource "random_password" "master_password" {
  length  = 16
  special = true
}

# Store the password in Secrets Manager
resource "aws_secretsmanager_secret" "db_password" {
  name                    = "debate-chatbot/database-password"
  description             = "Database master password for Debate Chatbot"
  recovery_window_in_days = 7
  kms_key_id              = var.kms_key.id

  tags = var.tags
}

resource "aws_secretsmanager_secret_version" "db_password" {
  secret_id = aws_secretsmanager_secret.db_password.id
  secret_string = jsonencode({
    username = "postgres"
    password = random_password.master_password.result
  })
}

# Aurora Serverless v2 Cluster
resource "aws_rds_cluster" "debate_aurora" {
  cluster_identifier              = "debate-chatbot-aurora"
  engine                          = "aurora-postgresql"
  engine_mode                     = "provisioned"
  engine_version                  = "15.4"
  database_name                   = "debate"
  master_username                 = "postgres"
  master_password                 = random_password.master_password.result
  backup_retention_period         = 7
  preferred_backup_window         = "07:00-09:00"
  preferred_maintenance_window    = "sun:05:00-sun:07:00"
  db_subnet_group_name            = var.vpc.db_subnet_group.name
  db_cluster_parameter_group_name = aws_rds_cluster_parameter_group.debate_aurora.name
  vpc_security_group_ids          = [var.vpc.aurora_security_group.id]
  storage_encrypted               = true
  kms_key_id                      = var.kms_key.arn
  deletion_protection             = false
  skip_final_snapshot             = true

  serverlessv2_scaling_configuration {
    max_capacity = 2
    min_capacity = 0.5
  }

  tags = var.tags
}

# Aurora Serverless v2 Instance
resource "aws_rds_cluster_instance" "debate_aurora" {
  count              = 1
  identifier         = "debate-chatbot-aurora-${count.index + 1}"
  cluster_identifier = aws_rds_cluster.debate_aurora.id
  instance_class     = "db.serverless"
  engine             = aws_rds_cluster.debate_aurora.engine
  engine_version     = aws_rds_cluster.debate_aurora.engine_version

  tags = var.tags
}

# Parameter Group for Aurora
resource "aws_rds_cluster_parameter_group" "debate_aurora" {
  family = "aurora-postgresql15"
  name   = "debate-chatbot-aurora-params"

  parameter {
    name  = "log_statement"
    value = "all"
  }

  parameter {
    name  = "log_min_duration_statement"
    value = "1000"
  }

  tags = var.tags
}

