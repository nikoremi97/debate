# Frontend Deployment Guide

This guide explains how to deploy the debate chatbot frontend to AWS S3 and CloudFront.

## Prerequisites

1. **AWS CLI configured** with appropriate permissions
2. **Node.js and npm** installed
3. **Terraform infrastructure** deployed (S3 bucket and CloudFront distribution)

## Quick Deployment

### Option 1: Using Makefile (Recommended)

```bash
# Deploy everything (build + upload + invalidate)
make deploy-frontend

# Or just build the frontend
make build-frontend
```

### Option 2: Using Deployment Scripts

```bash
# Full deployment with automatic CloudFront detection
./scripts/deploy-frontend.sh

# Simple deployment (requires distribution ID)
./scripts/deploy.sh <distribution-id>
```

## Manual Deployment Steps

If you prefer to run the steps manually:

### 1. Build the Frontend

```bash
cd frontend
npm install
npm run build
```

This creates a static export in the `frontend/out/` directory.

### 2. Upload to S3

```bash
aws s3 sync frontend/out/ s3://debate-chatbot-dashboard/ --delete
```

### 3. Invalidate CloudFront

```bash
# Get distribution ID
DISTRIBUTION_ID=$(aws cloudfront list-distributions \
    --query "DistributionList.Items[?Origins.Items[0].DomainName=='debate-chatbot-dashboard.s3.us-east-2.amazonaws.com'].Id" \
    --output text)

# Create invalidation
aws cloudfront create-invalidation \
    --distribution-id $DISTRIBUTION_ID \
    --paths "/*"
```

## Configuration

### Next.js Configuration

The frontend is configured for static export in `frontend/next.config.ts`:

```typescript
const nextConfig: NextConfig = {
  output: 'export',           // Static export
  trailingSlash: true,        // S3-friendly URLs
  images: {
    unoptimized: true         // No server-side image optimization
  }
};
```

### S3 Bucket Configuration

- **Bucket Name**: `debate-chatbot-dashboard`
- **Region**: `us-east-2`
- **Public Access**: Enabled for static website hosting
- **Website Configuration**: `index.html` as default document

### CloudFront Configuration

- **Origin Access Control (OAC)**: Secure access to S3
- **Cache Behaviors**: Optimized for static assets and SPA routing
- **Custom Error Pages**: 404/403 redirect to `index.html`
- **HTTPS Enforcement**: All traffic redirected to HTTPS

## Troubleshooting

### Common Issues

1. **Build Fails**
   ```bash
   cd frontend
   rm -rf node_modules package-lock.json
   npm install
   npm run build
   ```

2. **S3 Upload Fails**
   - Check AWS credentials: `aws sts get-caller-identity`
   - Verify bucket exists: `aws s3 ls s3://debate-chatbot-dashboard`

3. **CloudFront Invalidation Fails**
   - Get distribution ID: `aws cloudfront list-distributions`
   - Check distribution status in AWS Console

4. **Website Not Loading**
   - Check CloudFront invalidation status
   - Verify S3 bucket policy allows public read access
   - Check browser console for errors

### Getting Distribution ID

```bash
# Method 1: From Terraform outputs
cd terraform/projects/frontend
terraform output cloudfront_distribution_id

# Method 2: From AWS CLI
aws cloudfront list-distributions \
    --query "DistributionList.Items[?Origins.Items[0].DomainName=='debate-chatbot-dashboard.s3.us-east-2.amazonaws.com'].Id" \
    --output text
```

## Deployment Workflow

1. **Development**: Make changes to the frontend code
2. **Build**: Run `make build-frontend` or `npm run build`
3. **Deploy**: Run `make deploy-frontend` or `./scripts/deploy-frontend.sh`
4. **Verify**: Check the CloudFront URL to ensure deployment is successful

## Monitoring

- **CloudFront Metrics**: Available in AWS CloudWatch
- **S3 Access Logs**: Enable in S3 bucket settings
- **Distribution Status**: Check in AWS CloudFront Console

## Security Notes

- S3 bucket is configured for public read access (required for static hosting)
- CloudFront uses OAC for secure S3 access
- All traffic is forced to HTTPS
- CORS is configured for web application requests

## Cost Optimization

- CloudFront uses Price Class 100 (North America and Europe)
- Static assets cached for 1 year
- HTML files cached for 1 hour
- Gzip compression enabled
