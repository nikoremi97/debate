# Debate Chatbot Backend Infrastructure

This Terraform configuration creates the AWS infrastructure for the Debate Chatbot application using ECS Fargate with a public load balancer.

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Internet      │    │   ALB           │    │   ECS Fargate   │
│                 │◄──►│   (Public)      │◄──►│   Tasks         │
│                 │    │                 │    │   (Backend API) │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                        │
                                │                        │
                                ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Public        │    │   Private       │
                       │   Subnets       │    │   Subnets       │
                       │   (ALB)         │    │   (ECS Tasks)   │
                       └─────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
                                              ┌─────────────────┐
                                              │   ECR           │
                                              │   (Docker       │
                                              │   Images)       │
                                              └─────────────────┘
                                                        │
                                                        ▼
                                              ┌─────────────────┐
                                              │   Secrets       │
                                              │   Manager       │
                                              │   + KMS         │
                                              │   (Encrypted    │
                                              │   Secrets)      │
                                              └─────────────────┘
                                                        │
                                                        ▼
                                              ┌─────────────────┐
                                              │   Database      │
                                              │   (Aurora +     │
                                              │   Redis)        │
                                              └─────────────────┘
```

## Components

### 🔐 **KMS (Key Management Service)**
- Creates a dedicated KMS key for encrypting secrets
- Enables automatic key rotation for enhanced security
- Used to encrypt all secrets stored in Secrets Manager
- Provides audit trail for key usage

### 🐳 **ECR (Elastic Container Registry)**
- Repository for storing Docker images
- Lifecycle policies to manage image retention
- Scans images for vulnerabilities

### 🔒 **Secrets Manager**
- Stores OpenAI API key securely (encrypted with KMS)
- Stores database connection URL (encrypted with KMS)
- Stores API authentication key (encrypted with KMS)
- Automatic rotation capabilities
- IAM-based access control

### 🚀 **ECS (Elastic Container Service)**
- Fargate cluster for serverless containers
- Auto-scaling service with 2 desired tasks
- Health checks and monitoring
- IAM roles for secure access to secrets and databases
- Container image pulled from ECR

### ⚖️ **Application Load Balancer**
- Public-facing load balancer
- Health checks on `/health` and `/ready` endpoints
- Routes traffic to ECS tasks
- SSL/TLS termination support
- Cross-zone load balancing

### 🛡️ **Security Groups**
- ALB security group (ports 80, 443) - allows public internet access
- ECS tasks security group (port 8080) - allows traffic only from ALB
- Database security groups - allows access only from ECS tasks
- Proper ingress/egress rules with least privilege access

## Prerequisites

1. **VPC Infrastructure**: The network project must be deployed first
2. **AWS Credentials**: Configured AWS CLI or environment variables
3. **Docker Image**: Your application image must be built and pushed to ECR

## Usage

### 1. Configure Variables

Edit `config.auto.tfvars`:

```hcl
region = "us-east-2"
openai_api_key = "your-openai-api-key-here"  # Will be encrypted using KMS
database_url   = "your-database-url-here"    # Will be encrypted using KMS
api_key        = "your-secure-api-key-here"  # Will be encrypted using KMS
```

### 2. Deploy Infrastructure

```bash
# Initialize Terraform
terraform init

# Plan the deployment
terraform plan

# Apply the configuration
terraform apply
```

### 3. Build and Push Docker Image

```bash
# Get ECR login token
aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin <ECR_REPOSITORY_URL>

# Build your image
docker build -t debate-api .

# Tag the image
docker tag debate-api:latest <ECR_REPOSITORY_URL>:latest

# Push the image
docker push <ECR_REPOSITORY_URL>:latest
```

### 4. Access the Application

After deployment, you'll get the load balancer DNS name:

```bash
# Get the application URL
terraform output application_url
```

## Configuration

### Environment Variables

The ECS tasks are configured with:

- `PORT=8080` - Application port
- `OPENAI_MODEL=gpt-4o-mini` - Default OpenAI model
- `AWS_REGION=us-east-2` - AWS region for services
- `API_KEY_SECRET_NAME=debate-chatbot-api-key` - Secret name for API key
- `OPENAI_API_KEY` - Retrieved from Secrets Manager (KMS encrypted)
- `DATABASE_URL` - Retrieved from Secrets Manager (KMS encrypted)

### Health Checks

- **Path**: `/health` and `/ready`
- **Interval**: 30 seconds
- **Timeout**: 5 seconds
- **Healthy threshold**: 2
- **Unhealthy threshold**: 2

### Scaling

- **Desired tasks**: 2
- **CPU**: 512 units (0.5 vCPU)
- **Memory**: 1024 MB (1 GB)

## Security

- ✅ **KMS Encryption**: All secrets encrypted with dedicated KMS key
- ✅ **Network Isolation**: ECS tasks in private subnets, ALB in public subnets
- ✅ **API Authentication**: X-API-Key header authentication required
- ✅ **Secrets Management**: Secure storage in AWS Secrets Manager
- ✅ **IAM Roles**: Least privilege access for ECS tasks
- ✅ **Security Groups**: Minimal access rules between components
- ✅ **Container Security**: ECR image scanning enabled
- ✅ **Audit Trail**: CloudTrail logging for all API calls

## Monitoring

- CloudWatch logs for application logs
- ECS service metrics
- ALB metrics and access logs
- Container insights enabled

## Cost Optimization

- Fargate Spot instances available
- ECR lifecycle policies to clean old images
- CloudWatch log retention (7 days)
- Minimal resource allocation

## Troubleshooting

### Common Issues

1. **Health Check Failures**
   - Ensure your app responds to `/health` and `/ready`
   - Check security group rules
   - Verify container port configuration

2. **Secrets Access Issues**
   - Verify IAM role permissions for Secrets Manager
   - Check KMS key permissions for decryption
   - Ensure secrets exist in Secrets Manager
   - Verify secret names match environment variables

3. **API Authentication Issues**
   - Verify X-API-Key header is included in requests
   - Check API key is correctly stored in Secrets Manager
   - Ensure API key matches the expected value

4. **Image Pull Errors**
   - Verify ECR repository exists
   - Check IAM permissions for ECR
   - Ensure image is pushed to correct repository

### Useful Commands

```bash
# Check ECS service status
aws ecs describe-services --cluster debate-chatbot-cluster --services debate-api-service

# View application logs
aws logs tail /ecs/debate-api --follow

# Check ALB target health
aws elbv2 describe-target-health --target-group-arn <TARGET_GROUP_ARN>

# Test API authentication
curl -H "X-API-Key: your-api-key" https://your-alb-url/health

# Check secrets in Secrets Manager
aws secretsmanager list-secrets --query 'SecretList[?contains(Name, `debate-chatbot`)]'
```

## Cleanup

To destroy the infrastructure:

```bash
terraform destroy
```

**Note**: This will delete all resources including the ECR repository and stored images.
