# CloudFront Distribution for API HTTPS
resource "aws_cloudfront_distribution" "api" {
  origin {
    domain_name = var.alb_dns_name
    origin_id   = "api-alb"

    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "http-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  enabled             = true
  is_ipv6_enabled     = true
  comment             = "Debate Chatbot API - HTTPS"
  default_root_object = ""

  default_cache_behavior {
    allowed_methods        = ["HEAD", "DELETE", "POST", "GET", "OPTIONS", "PUT", "PATCH"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = "api-alb"
    compress               = true
    viewer_protocol_policy = "redirect-to-https"

    forwarded_values {
      query_string = true
      headers      = ["X-API-Key", "Authorization", "Content-Type", "Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"]
      cookies {
        forward = "none"
      }
    }

    min_ttl     = 0
    default_ttl = 0
    max_ttl     = 0
  }

  # Cache behavior for health endpoints (no caching)
  ordered_cache_behavior {
    path_pattern     = "/health*"
    allowed_methods  = ["HEAD", "GET", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "api-alb"
    compress         = true

    forwarded_values {
      query_string = true
      headers      = ["X-API-Key", "Authorization", "Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Access-Control-Allow-Origin"]
      cookies {
        forward = "none"
      }
    }

    min_ttl     = 0
    default_ttl = 0
    max_ttl     = 0

    viewer_protocol_policy = "redirect-to-https"
  }

  # Cache behavior for ready endpoints (no caching)
  ordered_cache_behavior {
    path_pattern     = "/ready*"
    allowed_methods  = ["HEAD", "GET", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "api-alb"
    compress         = true

    forwarded_values {
      query_string = true
      headers      = ["X-API-Key", "Authorization", "Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Access-Control-Allow-Origin"]
      cookies {
        forward = "none"
      }
    }

    min_ttl     = 0
    default_ttl = 0
    max_ttl     = 0

    viewer_protocol_policy = "redirect-to-https"
  }

  # Cache behavior for chat endpoints (no caching)
  ordered_cache_behavior {
    path_pattern     = "/chat*"
    allowed_methods  = ["HEAD", "DELETE", "POST", "GET", "OPTIONS", "PUT", "PATCH"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "api-alb"
    compress         = true

    forwarded_values {
      query_string = true
      headers      = ["X-API-Key", "Authorization", "Content-Type", "Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Access-Control-Allow-Origin"]
      cookies {
        forward = "none"
      }
    }

    min_ttl     = 0
    default_ttl = 0
    max_ttl     = 0

    viewer_protocol_policy = "redirect-to-https"
  }

  # Cache behavior for conversations endpoints (no caching)
  ordered_cache_behavior {
    path_pattern     = "/conversations*"
    allowed_methods  = ["HEAD", "DELETE", "POST", "GET", "OPTIONS", "PUT", "PATCH"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "api-alb"
    compress         = true

    forwarded_values {
      query_string = true
      headers      = ["X-API-Key", "Authorization", "Content-Type", "Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"]
      cookies {
        forward = "none"
      }
    }

    min_ttl     = 0
    default_ttl = 0
    max_ttl     = 0

    viewer_protocol_policy = "redirect-to-https"
  }

  price_class = "PriceClass_100"

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  tags = var.tags
}

# Output the CloudFront domain name
output "cloudfront_domain_name" {
  description = "CloudFront distribution domain name"
  value       = aws_cloudfront_distribution.api.domain_name
}

output "cloudfront_url" {
  description = "CloudFront distribution URL"
  value       = "https://${aws_cloudfront_distribution.api.domain_name}"
}
