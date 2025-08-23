# Herald.lol Infrastructure Variables

variable "aws_region" {
  description = "AWS region for resources"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name (staging/production)"
  type        = string
  validation {
    condition     = contains(["staging", "production"], var.environment)
    error_message = "Environment must be staging or production"
  }
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  description = "CIDR blocks for public subnets"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks for private subnets"
  type        = list(string)
  default     = ["10.0.11.0/24", "10.0.12.0/24", "10.0.13.0/24"]
}

variable "kubernetes_version" {
  description = "Kubernetes version for EKS cluster"
  type        = string
  default     = "1.28"
}

variable "rds_instance_class" {
  description = "RDS instance class for PostgreSQL"
  type        = string
  default     = "db.r6g.xlarge"
}

variable "redis_node_type" {
  description = "ElastiCache Redis node type"
  type        = string
  default     = "cache.r6g.large"
}

variable "acm_certificate_arn" {
  description = "ACM certificate ARN for HTTPS"
  type        = string
}

# Gaming-specific variables
variable "gaming_analytics_enabled" {
  description = "Enable gaming analytics features"
  type        = bool
  default     = true
}

variable "riot_api_rate_limit" {
  description = "Riot API rate limit per 2 minutes"
  type        = number
  default     = 100
}

variable "max_concurrent_users" {
  description = "Maximum concurrent users target"
  type        = number
  default     = 1000000
}

variable "analytics_processing_timeout" {
  description = "Maximum time for post-game analysis (seconds)"
  type        = number
  default     = 5
}

# Scaling variables
variable "min_api_replicas" {
  description = "Minimum API service replicas"
  type        = number
  default     = 3
}

variable "max_api_replicas" {
  description = "Maximum API service replicas"
  type        = number
  default     = 50
}

variable "target_cpu_utilization" {
  description = "Target CPU utilization for autoscaling"
  type        = number
  default     = 70
}

# Monitoring variables
variable "enable_monitoring" {
  description = "Enable comprehensive monitoring"
  type        = bool
  default     = true
}

variable "log_retention_days" {
  description = "CloudWatch log retention in days"
  type        = number
  default     = 30
}

variable "metrics_retention_days" {
  description = "Metrics retention in days"
  type        = number
  default     = 90
}

# Security variables
variable "enable_waf" {
  description = "Enable Web Application Firewall"
  type        = bool
  default     = true
}

variable "enable_guardduty" {
  description = "Enable AWS GuardDuty threat detection"
  type        = bool
  default     = true
}

variable "enable_secrets_rotation" {
  description = "Enable automatic secrets rotation"
  type        = bool
  default     = true
}

# Backup variables
variable "backup_retention_days" {
  description = "Database backup retention in days"
  type        = number
  default     = 30
}

variable "enable_point_in_time_recovery" {
  description = "Enable point-in-time recovery for databases"
  type        = bool
  default     = true
}

# Cost optimization variables
variable "enable_spot_instances" {
  description = "Use spot instances for non-critical workloads"
  type        = bool
  default     = false
}

variable "spot_max_price" {
  description = "Maximum price for spot instances"
  type        = string
  default     = "0.10"
}

# Tags
variable "tags" {
  description = "Additional tags for resources"
  type        = map(string)
  default = {
    Project     = "Herald.lol"
    Platform    = "Gaming Analytics"
    Stack       = "Production"
    Performance = "Critical"
  }
}