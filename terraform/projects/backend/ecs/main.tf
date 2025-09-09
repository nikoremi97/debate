# ECS Cluster
resource "aws_ecs_cluster" "debate_cluster" {
  name = "debate-chatbot-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = var.tags
}

# ECS Cluster Capacity Providers
resource "aws_ecs_cluster_capacity_providers" "debate_cluster" {
  cluster_name = aws_ecs_cluster.debate_cluster.name

  capacity_providers = ["FARGATE", "FARGATE_SPOT"]

  default_capacity_provider_strategy {
    base              = 1
    weight            = 100
    capacity_provider = "FARGATE"
  }
}

# Application Load Balancer
resource "aws_lb" "debate_alb" {
  name               = "debate-chatbot-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [var.alb_security_group_id]
  subnets            = var.public_subnet_ids

  enable_deletion_protection = false

  tags = var.tags
}

# ALB Target Group
resource "aws_lb_target_group" "debate_api" {
  name        = "debate-api-tg"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "ip"

  health_check {
    enabled             = true
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 5
    interval            = 30
    path                = "/healthz"
    matcher             = "200"
    port                = "traffic-port"
    protocol            = "HTTP"
  }

  tags = var.tags
}

# ALB Listener
resource "aws_lb_listener" "debate_api" {
  load_balancer_arn = aws_lb.debate_alb.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.debate_api.arn
  }
}

# ECS Task Definition
resource "aws_ecs_task_definition" "debate_api" {
  family                   = "debate-api"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 512
  memory                   = 1024
  execution_role_arn       = aws_iam_role.ecs_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name  = "debate-api"
      image = var.ecr_repository_url

      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
          protocol      = "tcp"
        }
      ]

      environment = [
        {
          name  = "PORT"
          value = "8080"
        },
        {
          name  = "OPENAI_MODEL"
          value = "gpt-4o-mini"
        },
        {
          name  = "AWS_REGION"
          value = var.region
        },
        {
          name  = "API_KEY_SECRET_NAME"
          value = "debate-chatbot-api-key"
        }
      ]

      secrets = [
        {
          name      = "OPENAI_API_KEY"
          valueFrom = var.openai_api_key_secret_arn
        },
        {
          name      = "DATABASE_URL"
          valueFrom = var.database_url_secret_arn
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.debate_api.name
          awslogs-region        = var.region
          awslogs-stream-prefix = "ecs"
        }
      }

      essential = true
    }
  ])

  tags = var.tags
}

# CloudWatch Log Group
resource "aws_cloudwatch_log_group" "debate_api" {
  name              = "/ecs/debate-api"
  retention_in_days = 7

  tags = var.tags
}

# ECS Service
resource "aws_ecs_service" "debate_api" {
  name            = "debate-api-service"
  cluster         = aws_ecs_cluster.debate_cluster.id
  task_definition = aws_ecs_task_definition.debate_api.arn
  desired_count   = 2
  launch_type     = "FARGATE"

  network_configuration {
    security_groups  = [var.ecs_tasks_security_group_id]
    subnets          = var.private_subnet_ids
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.debate_api.arn
    container_name   = "debate-api"
    container_port   = 8080
  }

  depends_on = [aws_lb_listener.debate_api]

  tags = var.tags
}

