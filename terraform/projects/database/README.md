# Debate Chatbot Database Infrastructure

This Terraform configuration creates the database infrastructure for the Debate Chatbot application using Aurora PostgreSQL (Serverless v2) and ElastiCache Redis.

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   ECS Tasks     │    │   Aurora        │    │   ElastiCache   │
│   (Backend)     │◄──►│   PostgreSQL    │    │   Redis         │
│                 │    │   Serverless v2 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Private       │    │   Private       │    │   Private       │
│   Subnets       │    │   Subnets       │    │   Subnets       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │   Secrets       │
                       │   Manager       │
                       │   + KMS         │
                       └─────────────────┘
```

## Components

### 🗄️ **Aurora PostgreSQL Serverless v2**
- **Engine**: Aurora PostgreSQL 15.4
- **Scaling**: Auto-scaling from 0.5 to 2 ACUs
- **High Availability**: Multi-AZ with automatic failover
- **Encryption**: At-rest and in-transit encryption with KMS
- **Backup**: 7-day retention with automated backups
- **Monitoring**: Enhanced monitoring and logging

### ⚡ **ElastiCache Redis**
- **Engine**: Redis 7
- **Node Type**: cache.t4g.micro (ARM-based, cost-effective)
- **High Availability**: Multi-AZ with automatic failover
- **Encryption**: At-rest and in-transit encryption
- **Auth**: Token-based authentication
- **Backup**: 5-day snapshot retention

### 🔐 **Security**
- **KMS Encryption**: Dedicated KMS key for database encryption
- **Secrets Manager**: Secure storage of passwords and auth tokens
- **Security Groups**: Restrictive access only from ECS tasks
- **Private Subnets**: All databases in private subnets only
- **VPC Isolation**: No public internet access

## Prerequisites

1. **VPC Infrastructure**: The network project must be deployed first
2. **ECS Security Groups**: Backend infrastructure must be deployed to get ECS security group IDs
3. **AWS Credentials**: Configured AWS CLI or environment variables

## Usage

### 1. Deploy Backend Infrastructure First

The database project requires ECS security group IDs from the backend project:

```bash
# Deploy backend first
cd ../backend
terraform apply
```

### 2. Get ECS Security Group IDs

```bash
# Get the ECS security group ID from backend outputs
cd ../backend
terraform output -json | jq -r '.ecs_security_group_id'
```

### 3. Configure Database Project

Edit `config.auto.tfvars`:

```hcl
region = "us-east-2"
ecs_security_group_ids = ["sg-xxxxxxxxx"]  # From backend output
```

### 4. Deploy Database Infrastructure

```bash
# Initialize Terraform
terraform init

# Plan the deployment
terraform plan

# Apply the configuration
terraform apply
```

## Configuration

### Aurora PostgreSQL Settings

- **Database Name**: `debate`
- **Master Username**: `postgres`
- **Port**: `5432`
- **Scaling**: 0.5-2 ACUs (auto-scaling)
- **Backup Window**: 07:00-09:00 UTC
- **Maintenance Window**: Sunday 05:00-07:00 UTC

### Redis Settings

- **Port**: `6379`
- **Node Type**: `cache.t4g.micro`
- **Num Nodes**: 2 (primary + replica)
- **Auth**: Token-based authentication
- **Memory Policy**: `allkeys-lru`
- **Timeout**: 300 seconds

### Security Groups

**Aurora Security Group**:
- Ingress: Port 5432 from ECS security groups only
- Egress: All traffic

**Redis Security Group**:
- Ingress: Port 6379 from ECS security groups only
- Egress: All traffic

## Connection Information

### Database Connection

```bash
# Get connection details
terraform output database_connection_string
# Output: postgresql://postgres@<endpoint>:5432/debate

# Get password from Secrets Manager
aws secretsmanager get-secret-value \
  --secret-id $(terraform output -raw database_password_secret_arn) \
  --query SecretString --output text | jq -r .password
```

### Redis Connection

```bash
# Get connection details
terraform output redis_connection_string
# Output: <endpoint>:6379

# Get auth token from Secrets Manager
aws secretsmanager get-secret-value \
  --secret-id $(terraform output -raw redis_auth_token_secret_arn) \
  --query SecretString --output text | jq -r .auth_token
```

## Environment Variables for Application

Your application should use these environment variables:

```bash
# Database
DATABASE_URL=postgresql://postgres:<password>@<endpoint>:5432/debate

# Redis
REDIS_URL=redis://:<auth_token>@<endpoint>:6379
```

## Cost Optimization

### Aurora Serverless v2
- ✅ Auto-scaling (0.5-2 ACUs)
- ✅ Pay-per-use pricing
- ✅ Automatic pause when idle
- ✅ No minimum charges

### ElastiCache Redis
- ✅ ARM-based t4g.micro instances
- ✅ Reserved instances available
- ✅ Multi-AZ for high availability
- ✅ Automated backups

## Monitoring

### CloudWatch Metrics
- **Aurora**: CPU, connections, read/write latency
- **Redis**: CPU, memory, cache hits/misses
- **Custom**: Application-specific metrics

### Logs
- **Aurora**: Query logs, error logs
- **Redis**: Slow log, general log
- **Application**: Via CloudWatch Logs

## Backup and Recovery

### Aurora
- **Automated Backups**: 7-day retention
- **Point-in-Time Recovery**: Available
- **Snapshot**: Manual snapshots supported
- **Cross-Region**: Replication available

### Redis
- **Automated Snapshots**: 5-day retention
- **Manual Snapshots**: Supported
- **Backup Window**: 03:00-05:00 UTC

## Security Best Practices

- ✅ **Encryption at Rest**: KMS encryption enabled
- ✅ **Encryption in Transit**: TLS/SSL enabled
- ✅ **Network Isolation**: Private subnets only
- ✅ **Access Control**: Security groups restrict access
- ✅ **Secrets Management**: Passwords in Secrets Manager
- ✅ **Key Rotation**: KMS key rotation enabled
- ✅ **Audit Logging**: CloudTrail integration

## Troubleshooting

### Common Issues

1. **Connection Timeouts**
   - Check security group rules
   - Verify subnet routing
   - Ensure ECS tasks are in correct subnets

2. **Authentication Failures**
   - Verify secrets in Secrets Manager
   - Check IAM permissions for secrets access
   - Ensure auth tokens are correct

3. **Performance Issues**
   - Monitor Aurora ACU usage
   - Check Redis memory usage
   - Review slow query logs

### Useful Commands

```bash
# Check Aurora cluster status
aws rds describe-db-clusters --db-cluster-identifier debate-chatbot-aurora

# Check Redis cluster status
aws elasticache describe-cache-clusters --cache-cluster-id debate-chatbot-redis

# View Aurora logs
aws logs describe-log-groups --log-group-name-prefix /aws/rds/cluster/debate-chatbot-aurora

# Test database connection
psql "postgresql://postgres:<password>@<endpoint>:5432/debate"

# Test Redis connection
redis-cli -h <endpoint> -p 6379 -a <auth_token> ping
```

## Cleanup

To destroy the database infrastructure:

```bash
terraform destroy
```

**Warning**: This will delete all data in Aurora and Redis. Make sure to backup any important data first.

## Integration with Backend

After deploying the database infrastructure, update your backend application to use the new database endpoints:

1. Update the `database_url` in your backend's `config.auto.tfvars`
2. Update the Redis connection string in your application
3. Redeploy the backend application

The database endpoints and credentials are automatically stored in Secrets Manager for secure access by your ECS tasks.
