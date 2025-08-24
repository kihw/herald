# Herald.lol Gaming Analytics Platform - Terraform Variables
# Variable definitions for AWS infrastructure deployment

variable "aws_region" {
  description = "AWS region for Herald.lol gaming infrastructure"
  type        = string
  default     = "us-east-1"
  
  validation {
    condition = can(regex("^[a-z0-9-]+$", var.aws_region))
    error_message = "AWS region must be a valid region name."
  }
}

variable "environment" {
  description = "Deployment environment (production, staging, development)"
  type        = string
  default     = "production"
  
  validation {
    condition     = contains(["production", "staging", "development"], var.environment)
    error_message = "Environment must be one of: production, staging, development."
  }
}

variable "cluster_name" {
  description = "EKS cluster name for Herald gaming platform"
  type        = string
  default     = "herald-gaming-cluster"
  
  validation {
    condition     = can(regex("^[a-zA-Z][a-zA-Z0-9-]*$", var.cluster_name))
    error_message = "Cluster name must start with a letter and contain only alphanumeric characters and hyphens."
  }
}

variable "cluster_version" {
  description = "Kubernetes version for EKS cluster"
  type        = string
  default     = "1.28"
}

# Node Group Configuration
variable "node_instance_type" {
  description = "EC2 instance type for EKS gaming nodes"
  type        = string
  default     = "c5.2xlarge"
  
  validation {
    condition = can(regex("^[a-z][0-9][a-z]*\\.[a-z0-9]+$", var.node_instance_type))
    error_message = "Instance type must be a valid EC2 instance type."
  }
}

variable "min_nodes" {
  description = "Minimum number of nodes in gaming cluster"
  type        = number
  default     = 3
  
  validation {
    condition     = var.min_nodes >= 1 && var.min_nodes <= 100
    error_message = "Minimum nodes must be between 1 and 100."
  }
}

variable "max_nodes" {
  description = "Maximum number of nodes in gaming cluster"
  type        = number
  default     = 100
  
  validation {
    condition     = var.max_nodes >= var.min_nodes && var.max_nodes <= 1000
    error_message = "Maximum nodes must be at least equal to min_nodes and not exceed 1000."
  }
}

variable "desired_nodes" {
  description = "Desired number of nodes in gaming cluster"
  type        = number
  default     = 10
  
  validation {
    condition     = var.desired_nodes >= var.min_nodes && var.desired_nodes <= var.max_nodes
    error_message = "Desired nodes must be between min_nodes and max_nodes."
  }
}

# Analytics Node Configuration
variable "analytics_node_instance_type" {
  description = "EC2 instance type for analytics processing nodes"
  type        = string
  default     = "r5.4xlarge"
}

variable "analytics_min_nodes" {
  description = "Minimum number of analytics nodes"
  type        = number
  default     = 0
}

variable "analytics_max_nodes" {
  description = "Maximum number of analytics nodes"
  type        = number
  default     = 20
}

variable "analytics_desired_nodes" {
  description = "Desired number of analytics nodes"
  type        = number
  default     = 2
}

# Database Configuration
variable "db_instance_class" {
  description = "RDS instance class for gaming database"
  type        = string
  default     = "db.r6g.2xlarge"
}

variable "db_engine_version" {
  description = "PostgreSQL engine version"
  type        = string
  default     = "15.4"
}

variable "db_backup_retention_days" {
  description = "Number of days to retain database backups"
  type        = number
  default     = 30
  
  validation {
    condition     = var.db_backup_retention_days >= 7 && var.db_backup_retention_days <= 35
    error_message = "Backup retention must be between 7 and 35 days."
  }
}

variable "db_multi_az" {
  description = "Enable Multi-AZ deployment for RDS"
  type        = bool
  default     = true
}

variable "db_deletion_protection" {
  description = "Enable deletion protection for RDS cluster"
  type        = bool
  default     = true
}

# Redis Configuration
variable "redis_node_type" {
  description = "ElastiCache node type for gaming cache"
  type        = string
  default     = "cache.r7g.2xlarge"
}

variable "redis_num_cache_clusters" {
  description = "Number of cache clusters for Redis"
  type        = number
  default     = 3
  
  validation {
    condition     = var.redis_num_cache_clusters >= 1 && var.redis_num_cache_clusters <= 10
    error_message = "Number of cache clusters must be between 1 and 10."
  }
}

variable "redis_num_node_groups" {
  description = "Number of node groups for Redis cluster mode"
  type        = number
  default     = 3
}

variable "redis_replicas_per_node_group" {
  description = "Number of replica nodes per node group"
  type        = number
  default     = 1
}

variable "redis_snapshot_retention_limit" {
  description = "Number of days to retain Redis snapshots"
  type        = number
  default     = 7
}

# Network Configuration
variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
  
  validation {
    condition     = can(cidrhost(var.vpc_cidr, 0))
    error_message = "VPC CIDR must be a valid CIDR block."
  }
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks for private subnets"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
}

variable "public_subnet_cidrs" {
  description = "CIDR blocks for public subnets"
  type        = list(string)
  default     = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
}

variable "enable_nat_gateway" {
  description = "Enable NAT Gateway for private subnets"
  type        = bool
  default     = true
}

variable "single_nat_gateway" {
  description = "Use single NAT Gateway for all private subnets"
  type        = bool
  default     = false
}

# Monitoring and Logging
variable "enable_flow_log" {
  description = "Enable VPC Flow Logs"
  type        = bool
  default     = true
}

variable "enable_performance_insights" {
  description = "Enable Performance Insights for RDS"
  type        = bool
  default     = true
}

variable "monitoring_interval" {
  description = "Enhanced monitoring interval for RDS"
  type        = number
  default     = 60
  
  validation {
    condition     = contains([0, 1, 5, 10, 15, 30, 60], var.monitoring_interval)
    error_message = "Monitoring interval must be one of: 0, 1, 5, 10, 15, 30, 60."
  }
}

# Security
variable "enable_encryption" {
  description = "Enable encryption at rest for storage"
  type        = bool
  default     = true
}

variable "enable_transit_encryption" {
  description = "Enable encryption in transit"
  type        = bool
  default     = true
}

variable "kms_key_deletion_window" {
  description = "KMS key deletion window in days"
  type        = number
  default     = 7
  
  validation {
    condition     = var.kms_key_deletion_window >= 7 && var.kms_key_deletion_window <= 30
    error_message = "KMS key deletion window must be between 7 and 30 days."
  }
}

# CloudFront
variable "cloudfront_price_class" {
  description = "CloudFront distribution price class"
  type        = string
  default     = "PriceClass_All"
  
  validation {
    condition = contains([
      "PriceClass_100", 
      "PriceClass_200", 
      "PriceClass_All"
    ], var.cloudfront_price_class)
    error_message = "Price class must be one of: PriceClass_100, PriceClass_200, PriceClass_All."
  }
}

variable "cloudfront_default_ttl" {
  description = "Default TTL for CloudFront cache"
  type        = number
  default     = 300
}

variable "cloudfront_max_ttl" {
  description = "Maximum TTL for CloudFront cache"
  type        = number
  default     = 31536000
}

# Gaming-specific Configuration
variable "gaming_performance_mode" {
  description = "Enable gaming performance optimizations"
  type        = bool
  default     = true
}

variable "low_latency_mode" {
  description = "Enable low latency optimizations for real-time gaming"
  type        = bool
  default     = true
}

variable "high_throughput_mode" {
  description = "Enable high throughput optimizations for analytics"
  type        = bool
  default     = true
}

# Tags
variable "common_tags" {
  description = "Common tags to apply to all resources"
  type        = map(string)
  default = {
    Project     = "Herald.lol"
    Team        = "Gaming-Analytics"
    Owner       = "herald-platform"
    CostCenter  = "gaming-analytics"
    Terraform   = "true"
  }
}

variable "additional_tags" {
  description = "Additional tags to apply to resources"
  type        = map(string)
  default     = {}
}