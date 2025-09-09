# ECS Tasks Security Group
resource "aws_security_group" "ecs_tasks" {
  name_prefix = "debate-ecs-tasks"
  vpc_id      = aws_vpc.basic.id

  ingress {
    description     = "HTTP from ALB"
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "debate-chatbot-ecs-tasks-sg"
  })
}

# ALB Security Group
resource "aws_security_group" "alb" {
  name_prefix = "debate-alb"
  vpc_id      = aws_vpc.basic.id

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "debate-chatbot-alb-sg"
  })
}


# Security Group for ElastiCache
resource "aws_security_group" "valkey" {
  name_prefix = "debate-valkey-"
  vpc_id      = aws_vpc.basic.id

  ingress {
    description     = "Valkey from ECS"
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs_tasks.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "debate-chatbot-valkey-sg"
  })
}


# DB Subnet Group
resource "aws_db_subnet_group" "debate_db" {
  name       = "debate-chatbot-db-subnet-group"
  subnet_ids = aws_subnet.private[*].id

  tags = merge(var.tags, {
    Name = "debate-chatbot-db-subnet-group"
  })
}

# Security Group for Aurora
resource "aws_security_group" "aurora" {
  name_prefix = "debate-aurora-"
  vpc_id      = aws_vpc.basic.id

  ingress {
    description     = "PostgreSQL from ECS"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs_tasks.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "debate-chatbot-aurora-sg"
  })
}
