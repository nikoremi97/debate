#!/bin/bash

# Frontend Deployment Script
# This script builds the Next.js frontend and deploys it to S3 with CloudFront invalidation

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
FRONTEND_DIR="frontend"
OUTPUT_DIR="out"
S3_BUCKET="debate-chatbot-dashboard"
REGION="us-east-2"
DISTRIBUTION_ID=""

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if required tools are installed
check_dependencies() {
    print_status "Checking dependencies..."

    if ! command -v node &> /dev/null; then
        print_error "Node.js is not installed. Please install Node.js first."
        exit 1
    fi

    if ! command -v npm &> /dev/null; then
        print_error "npm is not installed. Please install npm first."
        exit 1
    fi

    if ! command -v aws &> /dev/null; then
        print_error "AWS CLI is not installed. Please install AWS CLI first."
        exit 1
    fi

    print_success "All dependencies are available"
}

# Function to get CloudFront distribution ID
get_distribution_id() {
    print_status "Getting CloudFront distribution ID..."

    # Try to get distribution ID from Terraform outputs
    if command -v terraform &> /dev/null; then
        cd terraform/projects/frontend
        if [ -f "terraform.tfstate" ]; then
            DISTRIBUTION_ID=$(terraform output -raw cloudfront_distribution_id 2>/dev/null || echo "")
            cd - > /dev/null
        else
            cd - > /dev/null
        fi
    fi

    # If not found in Terraform, try to find by bucket name
    if [ -z "$DISTRIBUTION_ID" ]; then
        print_warning "Could not get distribution ID from Terraform. Trying to find by bucket name..."
        DISTRIBUTION_ID=$(aws cloudfront list-distributions --query "DistributionList.Items[?Origins.Items[0].DomainName=='${S3_BUCKET}.s3.${REGION}.amazonaws.com'].Id" --output text 2>/dev/null || echo "")
    fi

    if [ -z "$DISTRIBUTION_ID" ]; then
        print_error "Could not find CloudFront distribution ID. Please provide it manually:"
        read -p "Enter CloudFront Distribution ID: " DISTRIBUTION_ID
    fi

    if [ -z "$DISTRIBUTION_ID" ]; then
        print_error "Distribution ID is required for CloudFront invalidation"
        exit 1
    fi

    print_success "Found CloudFront distribution ID: $DISTRIBUTION_ID"
}

# Function to build the frontend
build_frontend() {
    print_status "Building frontend application..."

    cd $FRONTEND_DIR

    # Install dependencies if node_modules doesn't exist
    if [ ! -d "node_modules" ]; then
        print_status "Installing dependencies..."
        npm install
    fi

    # Build the application
    print_status "Building Next.js application..."
    npm run build

    # Check if build was successful
    if [ ! -d "$OUTPUT_DIR" ]; then
        print_error "Build failed - output directory not found"
        exit 1
    fi

    cd - > /dev/null
    print_success "Frontend build completed successfully"
}

# Function to deploy to S3
deploy_to_s3() {
    print_status "Deploying to S3 bucket: $S3_BUCKET"

    # Check if bucket exists
    if ! aws s3 ls "s3://$S3_BUCKET" &> /dev/null; then
        print_error "S3 bucket '$S3_BUCKET' does not exist. Please create it first with Terraform."
        exit 1
    fi

    # Sync files to S3
    aws s3 sync "$FRONTEND_DIR/$OUTPUT_DIR/" "s3://$S3_BUCKET/" --delete

    print_success "Files uploaded to S3 successfully"
}

# Function to invalidate CloudFront
invalidate_cloudfront() {
    print_status "Invalidating CloudFront distribution: $DISTRIBUTION_ID"

    # Create invalidation
    INVALIDATION_ID=$(aws cloudfront create-invalidation \
        --distribution-id "$DISTRIBUTION_ID" \
        --paths "/*" \
        --query 'Invalidation.Id' \
        --output text)

    if [ $? -eq 0 ]; then
        print_success "CloudFront invalidation created successfully"
        print_status "Invalidation ID: $INVALIDATION_ID"
        print_status "You can check the invalidation status in the AWS Console"
    else
        print_error "Failed to create CloudFront invalidation"
        exit 1
    fi
}

# Function to show deployment summary
show_summary() {
    print_success "Deployment completed successfully!"
    echo
    print_status "Deployment Summary:"
    echo "  - S3 Bucket: $S3_BUCKET"
    echo "  - CloudFront Distribution: $DISTRIBUTION_ID"
    echo "  - Region: $REGION"
    echo
    print_status "Your dashboard should be available at:"
    echo "  https://$(aws cloudfront get-distribution --id $DISTRIBUTION_ID --query 'Distribution.DomainName' --output text 2>/dev/null || echo 'distribution-domain')"
    echo
    print_warning "Note: It may take a few minutes for the CloudFront invalidation to complete"
}

# Main execution
main() {
    echo "=========================================="
    echo "  Frontend Deployment Script"
    echo "=========================================="
    echo

    check_dependencies
    get_distribution_id
    build_frontend
    deploy_to_s3
    invalidate_cloudfront
    show_summary
}

# Run main function
main "$@"
