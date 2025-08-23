# Herald.lol AWS Infrastructure - Main Configuration
# High-performance gaming analytics platform infrastructure

terraform {
  required_version = ">= 1.5.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.23"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.11"
    }
  }
  
  backend "s3" {
    bucket         = "herald-terraform-state"
    key            = "production/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "herald-terraform-locks"
  }
}

# AWS Provider Configuration
provider "aws" {
  region = var.aws_region
  
  default_tags {
    tags = {
      Project     = "Herald.lol"
      Environment = var.environment
      ManagedBy   = "Terraform"
      Platform    = "Gaming Analytics"
      Stack       = "Production"
    }
  }
}

# Data Sources
data "aws_availability_zones" "available" {
  state = "available"
}

data "aws_caller_identity" "current" {}

# VPC Module - Multi-AZ Network Infrastructure
module "vpc" {
  source = "./modules/vpc"
  
  name               = "herald-${var.environment}"
  cidr               = var.vpc_cidr
  availability_zones = data.aws_availability_zones.available.names
  environment        = var.environment
  
  enable_nat_gateway = true
  enable_vpn_gateway = false
  enable_flow_logs   = true
  
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
  
  tags = {
    "kubernetes.io/cluster/herald-${var.environment}" = "shared"
  }
}

# EKS Cluster - Kubernetes for Container Orchestration
module "eks" {
  source = "./modules/eks"
  
  cluster_name    = "herald-${var.environment}"
  cluster_version = var.kubernetes_version
  
  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnet_ids
  
  # Node Groups for Gaming Workloads
  node_groups = {
    # General purpose nodes for API and web services
    general = {
      desired_capacity = 3
      min_capacity     = 3
      max_capacity     = 10
      
      instance_types = ["t3.large"]
      
      k8s_labels = {
        Environment = var.environment
        Workload    = "general"
      }
    }
    
    # High-performance nodes for gaming analytics
    analytics = {
      desired_capacity = 2
      min_capacity     = 2
      max_capacity     = 8
      
      instance_types = ["c5.2xlarge"]
      
      k8s_labels = {
        Environment = var.environment
        Workload    = "analytics"
        Gaming      = "true"
      }
      
      taints = [{
        key    = "analytics"
        value  = "true"
        effect = "NO_SCHEDULE"
      }]
    }
    
    # Memory-optimized nodes for caching
    cache = {
      desired_capacity = 2
      min_capacity     = 1
      max_capacity     = 4
      
      instance_types = ["r5.xlarge"]
      
      k8s_labels = {
        Environment = var.environment
        Workload    = "cache"
      }
    }
  }
  
  # OIDC for IRSA (IAM Roles for Service Accounts)
  enable_irsa = true
  
  # Cluster addons
  cluster_addons = {
    coredns = {
      most_recent = true
    }
    kube-proxy = {
      most_recent = true
    }
    vpc-cni = {
      most_recent = true
    }
    aws-ebs-csi-driver = {
      most_recent = true
    }
  }
  
  tags = {
    Platform = "Herald.lol"
    Gaming   = "true"
  }
}

# RDS PostgreSQL - Primary Database for Gaming Data
module "rds" {
  source = "./modules/rds"
  
  identifier = "herald-${var.environment}"
  
  engine         = "postgres"
  engine_version = "15.4"
  instance_class = var.rds_instance_class
  
  allocated_storage     = 100
  max_allocated_storage = 1000
  storage_encrypted     = true
  storage_type          = "gp3"
  
  database_name = "herald"
  username      = "herald_admin"
  
  vpc_id                  = module.vpc.vpc_id
  subnet_ids              = module.vpc.private_subnet_ids
  allowed_security_groups = [module.eks.cluster_security_group_id]
  
  # Multi-AZ for high availability
  multi_az               = true
  backup_retention_days  = 30
  backup_window          = "03:00-04:00"
  maintenance_window     = "sun:04:00-sun:05:00"
  
  # Performance Insights for gaming metrics monitoring
  performance_insights_enabled = true
  performance_insights_retention_period = 7
  
  # Read replicas for analytics workloads
  create_read_replicas = true
  read_replica_count   = 2
  
  tags = {
    Gaming = "true"
    Critical = "true"
  }
}

# ElastiCache Redis - High-Performance Caching for Gaming Data
module "elasticache" {
  source = "./modules/elasticache"
  
  cluster_id = "herald-${var.environment}"
  
  engine         = "redis"
  engine_version = "7.0"
  node_type      = var.redis_node_type
  
  # Cluster mode for horizontal scaling
  cluster_mode_enabled = true
  num_cache_clusters   = 3
  
  vpc_id                  = module.vpc.vpc_id
  subnet_ids              = module.vpc.private_subnet_ids
  allowed_security_groups = [module.eks.cluster_security_group_id]
  
  # Automatic failover for high availability
  automatic_failover_enabled = true
  multi_az_enabled           = true
  
  # Snapshots for data persistence
  snapshot_retention_limit = 5
  snapshot_window          = "03:00-05:00"
  
  # Gaming-specific Redis configuration
  parameter_group_family = "redis7"
  parameters = [
    {
      name  = "maxmemory-policy"
      value = "allkeys-lru"
    },
    {
      name  = "timeout"
      value = "300"
    }
  ]
  
  tags = {
    Gaming = "true"
    Performance = "critical"
  }
}

# S3 Buckets - Object Storage for Gaming Assets
module "s3" {
  source = "./modules/s3"
  
  buckets = {
    # Static assets (champion images, icons)
    assets = {
      name = "herald-${var.environment}-assets"
      versioning = true
      lifecycle_rules = [
        {
          id      = "archive-old-assets"
          enabled = true
          transition = {
            days          = 90
            storage_class = "STANDARD_IA"
          }
        }
      ]
    }
    
    # Match replays and analytics data
    analytics = {
      name = "herald-${var.environment}-analytics"
      versioning = true
      lifecycle_rules = [
        {
          id      = "delete-old-analytics"
          enabled = true
          expiration = {
            days = 365
          }
        }
      ]
    }
    
    # Backups
    backups = {
      name = "herald-${var.environment}-backups"
      versioning = true
      lifecycle_rules = [
        {
          id      = "transition-old-backups"
          enabled = true
          transition = {
            days          = 30
            storage_class = "GLACIER"
          }
        }
      ]
    }
  }
  
  tags = {
    Gaming = "true"
  }
}

# CloudFront CDN - Global Content Delivery for Gaming Assets
module "cloudfront" {
  source = "./modules/cloudfront"
  
  origins = {
    s3_assets = {
      domain_name = module.s3.bucket_regional_domain_names["assets"]
      origin_id   = "s3-assets"
      s3_origin   = true
    }
    
    api = {
      domain_name = module.alb.dns_name
      origin_id   = "api"
      custom_origin = {
        http_port              = 80
        https_port             = 443
        origin_protocol_policy = "https-only"
      }
    }
  }
  
  default_cache_behavior = {
    target_origin_id = "s3-assets"
    viewer_protocol_policy = "redirect-to-https"
    
    allowed_methods = ["GET", "HEAD", "OPTIONS"]
    cached_methods  = ["GET", "HEAD"]
    
    cache_policy_id            = "658327ea-f89d-4fab-a63d-7e88639e58f6" # Managed-CachingOptimized
    origin_request_policy_id   = "88a5eaf4-2fd4-4709-b370-b4c650ea3fcf" # Managed-CORS-S3Origin
    response_headers_policy_id = "5cc3b908-e619-46e2-b9c2-d499de5e7e17" # Managed-CORS-With-Preflight
  }
  
  ordered_cache_behaviors = [
    {
      path_pattern     = "/api/*"
      target_origin_id = "api"
      
      viewer_protocol_policy = "https-only"
      allowed_methods        = ["GET", "HEAD", "OPTIONS", "PUT", "POST", "PATCH", "DELETE"]
      cached_methods         = ["GET", "HEAD"]
      
      cache_policy_id          = "4135ea2d-6df8-44a3-9df3-4b5a84be39ad" # Managed-CachingDisabled
      origin_request_policy_id = "b689b0a8-53d0-40ab-baf2-68738e2966ac" # Managed-AllViewerExceptCookie
    }
  ]
  
  price_class = "PriceClass_200" # US, Canada, Europe, Asia, Middle East, Africa
  
  tags = {
    Gaming = "true"
    Global = "true"
  }
}

# Application Load Balancer - Traffic Distribution
module "alb" {
  source = "./modules/alb"
  
  name = "herald-${var.environment}"
  
  vpc_id          = module.vpc.vpc_id
  subnet_ids      = module.vpc.public_subnet_ids
  security_groups = [aws_security_group.alb.id]
  
  target_groups = {
    api = {
      port     = 8080
      protocol = "HTTP"
      health_check = {
        enabled             = true
        path                = "/health"
        healthy_threshold   = 2
        unhealthy_threshold = 2
        timeout             = 5
        interval            = 30
      }
    }
    
    frontend = {
      port     = 3000
      protocol = "HTTP"
      health_check = {
        enabled             = true
        path                = "/"
        healthy_threshold   = 2
        unhealthy_threshold = 2
        timeout             = 5
        interval            = 30
      }
    }
  }
  
  listeners = {
    http = {
      port     = 80
      protocol = "HTTP"
      default_action = {
        type = "redirect"
        redirect = {
          port        = "443"
          protocol    = "HTTPS"
          status_code = "HTTP_301"
        }
      }
    }
    
    https = {
      port            = 443
      protocol        = "HTTPS"
      certificate_arn = var.acm_certificate_arn
      default_action = {
        type             = "forward"
        target_group_arn = module.alb.target_group_arns["frontend"]
      }
      
      rules = [
        {
          priority = 100
          condition = {
            path_pattern = ["/api/*"]
          }
          action = {
            type             = "forward"
            target_group_arn = module.alb.target_group_arns["api"]
          }
        }
      ]
    }
  }
  
  tags = {
    Gaming = "true"
  }
}

# Security Group for ALB
resource "aws_security_group" "alb" {
  name_prefix = "herald-alb-"
  vpc_id      = module.vpc.vpc_id
  
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  ingress {
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
  
  tags = {
    Name = "herald-alb-${var.environment}"
  }
}

# WAF - Web Application Firewall for Gaming Platform Protection
module "waf" {
  source = "./modules/waf"
  
  name = "herald-${var.environment}"
  
  # Associate with ALB
  resource_arn = module.alb.arn
  
  # Gaming-specific rate limiting
  rate_limit_rules = {
    api = {
      name     = "api-rate-limit"
      priority = 1
      limit    = 2000
      window   = 300 # 5 minutes
    }
    
    riot_api = {
      name     = "riot-api-rate-limit"
      priority = 2
      limit    = 100
      window   = 120 # 2 minutes (Riot API limits)
    }
  }
  
  # IP reputation lists
  ip_reputation_lists = [
    "AWSManagedRulesAmazonIpReputationList",
    "AWSManagedRulesAnonymousIPList"
  ]
  
  tags = {
    Gaming   = "true"
    Security = "critical"
  }
}

# Outputs
output "eks_cluster_endpoint" {
  description = "EKS cluster endpoint"
  value       = module.eks.cluster_endpoint
}

output "rds_endpoint" {
  description = "RDS PostgreSQL endpoint"
  value       = module.rds.endpoint
  sensitive   = true
}

output "redis_endpoint" {
  description = "ElastiCache Redis endpoint"
  value       = module.elasticache.endpoint
  sensitive   = true
}

output "alb_dns_name" {
  description = "ALB DNS name"
  value       = module.alb.dns_name
}

output "cloudfront_domain" {
  description = "CloudFront distribution domain"
  value       = module.cloudfront.domain_name
}