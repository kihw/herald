# Herald.lol Gaming Analytics Platform - AWS Infrastructure
# Terraform Configuration for Production Environment

terraform {
  required_version = ">= 1.0"
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
      Team        = "Gaming-Analytics"
      Owner       = "herald-platform"
      CostCenter  = "gaming-analytics"
    }
  }
}

# Data Sources
data "aws_caller_identity" "current" {}
data "aws_availability_zones" "available" {
  state = "available"
}

# Variables
variable "aws_region" {
  description = "AWS region for Herald.lol infrastructure"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

variable "cluster_name" {
  description = "EKS cluster name"
  type        = string
  default     = "herald-gaming-cluster"
}

variable "node_instance_type" {
  description = "EC2 instance type for EKS nodes"
  type        = string
  default     = "c5.2xlarge"
}

variable "min_nodes" {
  description = "Minimum number of nodes"
  type        = number
  default     = 3
}

variable "max_nodes" {
  description = "Maximum number of nodes"
  type        = number
  default     = 100
}

variable "desired_nodes" {
  description = "Desired number of nodes"
  type        = number
  default     = 10
}

# VPC Configuration
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"

  name = "${var.cluster_name}-vpc"
  cidr = "10.0.0.0/16"

  azs             = slice(data.aws_availability_zones.available.names, 0, 3)
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  enable_nat_gateway     = true
  single_nat_gateway     = false
  enable_vpn_gateway     = false
  enable_dns_hostnames   = true
  enable_dns_support     = true

  # Gaming-optimized networking
  enable_flow_log      = true
  flow_log_destination_type = "cloud-watch-logs"

  public_subnet_tags = {
    "kubernetes.io/role/elb"                    = "1"
    "kubernetes.io/cluster/${var.cluster_name}" = "owned"
    Type = "Gaming-Public"
  }

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb"           = "1"
    "kubernetes.io/cluster/${var.cluster_name}" = "owned"
    Type = "Gaming-Private"
  }

  tags = {
    "kubernetes.io/cluster/${var.cluster_name}" = "owned"
    Name = "Herald Gaming VPC"
  }
}

# EKS Cluster
module "eks" {
  source = "terraform-aws-modules/eks/aws"
  version = "~> 19.0"

  cluster_name    = var.cluster_name
  cluster_version = "1.28"

  vpc_id                         = module.vpc.vpc_id
  subnet_ids                     = module.vpc.private_subnets
  cluster_endpoint_public_access = true

  # Gaming-optimized cluster configuration
  cluster_addons = {
    coredns = {
      most_recent = true
    }
    kube-proxy = {
      most_recent = true
    }
    vpc-cni = {
      most_recent = true
      configuration_values = jsonencode({
        env = {
          ENABLE_PREFIX_DELEGATION = "true"
          WARM_PREFIX_TARGET       = "1"
        }
      })
    }
    aws-ebs-csi-driver = {
      most_recent = true
    }
  }

  # EKS Managed Node Groups
  eks_managed_node_groups = {
    gaming-primary = {
      name = "gaming-primary-nodes"
      
      instance_types = [var.node_instance_type]
      capacity_type  = "ON_DEMAND"
      
      min_size     = var.min_nodes
      max_size     = var.max_nodes
      desired_size = var.desired_nodes

      ami_type = "AL2_x86_64"
      
      # Gaming workload optimizations
      pre_bootstrap_user_data = <<-EOT
        #!/bin/bash
        # Gaming-specific optimizations
        echo 'net.core.rmem_max = 134217728' >> /etc/sysctl.conf
        echo 'net.core.wmem_max = 134217728' >> /etc/sysctl.conf
        echo 'net.ipv4.tcp_rmem = 4096 87380 134217728' >> /etc/sysctl.conf
        echo 'net.ipv4.tcp_wmem = 4096 65536 134217728' >> /etc/sysctl.conf
        sysctl -p
      EOT

      labels = {
        WorkloadType = "gaming-analytics"
      }

      tags = {
        Name = "Herald Gaming Nodes"
        NodeGroup = "gaming-primary"
      }
    }

    # High-memory nodes for analytics processing
    analytics-nodes = {
      name = "analytics-processing-nodes"
      
      instance_types = ["r5.4xlarge"]
      capacity_type  = "SPOT"
      
      min_size     = 0
      max_size     = 20
      desired_size = 2

      labels = {
        WorkloadType = "analytics-processing"
        NodeType     = "high-memory"
      }

      taints = [
        {
          key    = "analytics-workload"
          value  = "true"
          effect = "NO_SCHEDULE"
        }
      ]

      tags = {
        Name = "Herald Analytics Nodes"
        NodeGroup = "analytics-processing"
      }
    }
  }

  tags = {
    Name = "Herald Gaming EKS Cluster"
  }
}

# RDS PostgreSQL Cluster for Gaming Data
resource "aws_rds_cluster" "herald_gaming_db" {
  cluster_identifier      = "herald-gaming-cluster"
  engine                 = "aurora-postgresql"
  engine_version         = "15.4"
  database_name          = "herald_gaming"
  master_username        = "herald_admin"
  manage_master_user_password = true

  # Multi-AZ deployment for high availability
  availability_zones = slice(data.aws_availability_zones.available.names, 0, 3)
  
  vpc_security_group_ids = [aws_security_group.rds_gaming.id]
  db_subnet_group_name   = aws_db_subnet_group.gaming.name
  
  # Gaming-optimized settings
  port = 5432
  
  backup_retention_period = 30
  preferred_backup_window = "03:00-04:00"
  preferred_maintenance_window = "sun:04:00-sun:05:00"
  
  deletion_protection = true
  skip_final_snapshot = false
  final_snapshot_identifier = "herald-gaming-final-snapshot-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"
  
  # Performance optimizations for gaming workloads
  db_cluster_parameter_group_name = aws_rds_cluster_parameter_group.gaming_params.name
  
  # Encryption
  storage_encrypted = true
  kms_key_id        = aws_kms_key.gaming_db.arn
  
  # Enable enhanced monitoring
  enabled_cloudwatch_logs_exports = ["postgresql"]
  
  tags = {
    Name = "Herald Gaming Database Cluster"
    Purpose = "Gaming Analytics Data"
  }
}

# RDS Cluster Instances (Read Replicas in 3 zones)
resource "aws_rds_cluster_instance" "herald_gaming_primary" {
  count              = 1
  identifier         = "herald-gaming-primary-${count.index}"
  cluster_identifier = aws_rds_cluster.herald_gaming_db.id
  instance_class     = "db.r6g.2xlarge"
  engine             = aws_rds_cluster.herald_gaming_db.engine
  engine_version     = aws_rds_cluster.herald_gaming_db.engine_version
  
  availability_zone = data.aws_availability_zones.available.names[0]
  
  performance_insights_enabled = true
  monitoring_interval         = 60
  monitoring_role_arn        = aws_iam_role.rds_enhanced_monitoring.arn
  
  tags = {
    Name = "Herald Gaming Primary DB"
    Role = "primary"
  }
}

resource "aws_rds_cluster_instance" "herald_gaming_readers" {
  count              = 2
  identifier         = "herald-gaming-reader-${count.index}"
  cluster_identifier = aws_rds_cluster.herald_gaming_db.id
  instance_class     = "db.r6g.xlarge"
  engine             = aws_rds_cluster.herald_gaming_db.engine
  engine_version     = aws_rds_cluster.herald_gaming_db.engine_version
  
  availability_zone = data.aws_availability_zones.available.names[count.index + 1]
  
  performance_insights_enabled = true
  monitoring_interval         = 60
  monitoring_role_arn        = aws_iam_role.rds_enhanced_monitoring.arn
  
  tags = {
    Name = "Herald Gaming Reader DB ${count.index + 1}"
    Role = "reader"
  }
}

# RDS Parameter Group for Gaming Workloads
resource "aws_rds_cluster_parameter_group" "gaming_params" {
  name   = "herald-gaming-params"
  family = "aurora-postgresql15"

  # Gaming-optimized PostgreSQL parameters
  parameter {
    name  = "shared_preload_libraries"
    value = "pg_stat_statements,auto_explain"
  }

  parameter {
    name  = "max_connections"
    value = "2000"
  }

  parameter {
    name  = "effective_cache_size"
    value = "24GB"
  }

  parameter {
    name  = "random_page_cost"
    value = "1.1"
  }

  parameter {
    name  = "work_mem"
    value = "64MB"
  }

  parameter {
    name  = "maintenance_work_mem"
    value = "2GB"
  }

  tags = {
    Name = "Herald Gaming DB Parameters"
  }
}

# ElastiCache Redis Cluster for Gaming Sessions
resource "aws_elasticache_replication_group" "herald_gaming_redis" {
  replication_group_id         = "herald-gaming-redis"
  description                  = "Herald Gaming Analytics Redis Cluster"
  
  port                         = 6379
  parameter_group_name         = "default.redis7.cluster.on"
  node_type                    = "cache.r7g.2xlarge"
  num_cache_clusters           = 3
  
  # Multi-AZ with automatic failover
  multi_az_enabled             = true
  automatic_failover_enabled   = true
  
  # Clustering mode for horizontal scaling
  num_node_groups              = 3
  replicas_per_node_group      = 1
  
  subnet_group_name            = aws_elasticache_subnet_group.gaming.name
  security_group_ids           = [aws_security_group.elasticache_gaming.id]
  
  # Encryption and security
  at_rest_encryption_enabled   = true
  transit_encryption_enabled   = true
  auth_token                   = random_password.redis_auth.result
  
  # Maintenance and backups
  maintenance_window           = "sun:03:00-sun:04:00"
  snapshot_retention_limit     = 7
  snapshot_window              = "02:00-03:00"
  
  tags = {
    Name = "Herald Gaming Redis Cluster"
    Purpose = "Gaming Session Cache"
  }
}

# S3 Buckets for Gaming Assets and Backups
resource "aws_s3_bucket" "herald_gaming_assets" {
  bucket = "herald-gaming-assets-${random_string.bucket_suffix.result}"
  
  tags = {
    Name = "Herald Gaming Assets"
    Purpose = "Gaming Assets and Static Files"
  }
}

resource "aws_s3_bucket" "herald_gaming_backups" {
  bucket = "herald-gaming-backups-${random_string.bucket_suffix.result}"
  
  tags = {
    Name = "Herald Gaming Backups"
    Purpose = "Database and Application Backups"
  }
}

# CloudFront Distribution for Global Performance
resource "aws_cloudfront_distribution" "herald_gaming_cdn" {
  origin {
    domain_name = aws_s3_bucket.herald_gaming_assets.bucket_regional_domain_name
    origin_id   = "S3-herald-gaming-assets"
    
    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.gaming_oai.cloudfront_access_identity_path
    }
  }

  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"
  
  # Gaming-optimized caching
  default_cache_behavior {
    allowed_methods        = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = "S3-herald-gaming-assets"
    compress               = true
    viewer_protocol_policy = "redirect-to-https"
    
    # Low latency for gaming content
    min_ttl     = 0
    default_ttl = 300
    max_ttl     = 31536000
    
    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }
  }
  
  # Global edge locations for minimum latency
  price_class = "PriceClass_All"
  
  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }
  
  viewer_certificate {
    cloudfront_default_certificate = true
  }
  
  tags = {
    Name = "Herald Gaming CDN"
    Purpose = "Global Content Delivery"
  }
}

# Security Groups
resource "aws_security_group" "rds_gaming" {
  name        = "herald-rds-gaming"
  description = "Security group for Herald Gaming RDS cluster"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description = "PostgreSQL from EKS nodes"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "Herald Gaming RDS Security Group"
  }
}

resource "aws_security_group" "elasticache_gaming" {
  name        = "herald-elasticache-gaming"
  description = "Security group for Herald Gaming ElastiCache cluster"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description = "Redis from EKS nodes"
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "Herald Gaming ElastiCache Security Group"
  }
}

# Random resources for unique naming
resource "random_string" "bucket_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "random_password" "redis_auth" {
  length  = 32
  special = true
}

# Supporting resources
resource "aws_db_subnet_group" "gaming" {
  name       = "herald-gaming-subnet-group"
  subnet_ids = module.vpc.private_subnets

  tags = {
    Name = "Herald Gaming DB Subnet Group"
  }
}

resource "aws_elasticache_subnet_group" "gaming" {
  name       = "herald-gaming-cache-subnet"
  subnet_ids = module.vpc.private_subnets

  tags = {
    Name = "Herald Gaming Cache Subnet Group"
  }
}

resource "aws_kms_key" "gaming_db" {
  description             = "KMS key for Herald Gaming database encryption"
  deletion_window_in_days = 7

  tags = {
    Name = "Herald Gaming DB Encryption Key"
  }
}

resource "aws_kms_alias" "gaming_db" {
  name          = "alias/herald-gaming-db"
  target_key_id = aws_kms_key.gaming_db.key_id
}

resource "aws_iam_role" "rds_enhanced_monitoring" {
  name = "herald-rds-enhanced-monitoring-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "monitoring.rds.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "rds_enhanced_monitoring" {
  role       = aws_iam_role.rds_enhanced_monitoring.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole"
}

resource "aws_cloudfront_origin_access_identity" "gaming_oai" {
  comment = "Herald Gaming Assets OAI"
}

# Output values
output "cluster_endpoint" {
  description = "Endpoint for EKS control plane"
  value       = module.eks.cluster_endpoint
}

output "cluster_security_group_id" {
  description = "Security group ids attached to the cluster control plane"
  value       = module.eks.cluster_security_group_id
}

output "region" {
  description = "AWS region"
  value       = var.aws_region
}

output "cluster_name" {
  description = "Kubernetes Cluster Name"
  value       = module.eks.cluster_name
}

output "rds_cluster_endpoint" {
  description = "RDS cluster endpoint"
  value       = aws_rds_cluster.herald_gaming_db.endpoint
}

output "rds_cluster_reader_endpoint" {
  description = "RDS cluster reader endpoint"
  value       = aws_rds_cluster.herald_gaming_db.reader_endpoint
}

output "redis_configuration_endpoint" {
  description = "Redis cluster configuration endpoint"
  value       = aws_elasticache_replication_group.herald_gaming_redis.configuration_endpoint_address
}

output "cloudfront_distribution_id" {
  description = "CloudFront Distribution ID"
  value       = aws_cloudfront_distribution.herald_gaming_cdn.id
}

output "cloudfront_domain_name" {
  description = "CloudFront Distribution Domain Name"
  value       = aws_cloudfront_distribution.herald_gaming_cdn.domain_name
}