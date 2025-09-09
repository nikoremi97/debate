#!/bin/bash

# Simple deployment script for the debate chatbot frontend
# Usage: ./scripts/deploy.sh [distribution-id]

set -e

# Configuration
S3_BUCKET="debate-chatbot-dashboard"
REGION="us-east-2"
DISTRIBUTION_ID="${1:-}"

echo "🚀 Deploying frontend to S3 and CloudFront..."

# Build the frontend
echo "📦 Building frontend..."
cd frontend
npm install
npm run build
cd ..

# Upload to S3
echo "☁️  Uploading to S3..."
aws s3 sync frontend/out/ s3://$S3_BUCKET/ --delete

# Get distribution ID if not provided
if [ -z "$DISTRIBUTION_ID" ]; then
    echo "🔍 Finding CloudFront distribution..."
    DISTRIBUTION_ID=$(aws cloudfront list-distributions \
        --query "DistributionList.Items[?Origins.Items[0].DomainName=='${S3_BUCKET}.s3.${REGION}.amazonaws.com'].Id" \
        --output text)
fi

if [ -z "$DISTRIBUTION_ID" ]; then
    echo "❌ Could not find CloudFront distribution. Please provide it manually:"
    echo "Usage: ./scripts/deploy.sh <distribution-id>"
    exit 1
fi

# Invalidate CloudFront
echo "🔄 Invalidating CloudFront cache..."
aws cloudfront create-invalidation \
    --distribution-id $DISTRIBUTION_ID \
    --paths "/*"

echo "✅ Deployment completed!"
echo "🌐 Your dashboard is available at:"
aws cloudfront get-distribution --id $DISTRIBUTION_ID --query 'Distribution.DomainName' --output text | sed 's/^/https:\/\//'
