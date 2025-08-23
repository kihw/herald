# Herald.lol Production Environment Configuration

environment = "production"
aws_region  = "us-east-1"

# Network Configuration
vpc_cidr             = "10.0.0.0/16"
public_subnet_cidrs  = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
private_subnet_cidrs = ["10.0.11.0/24", "10.0.12.0/24", "10.0.13.0/24"]

# Kubernetes Configuration
kubernetes_version = "1.28"

# Database Configuration - Production Grade
rds_instance_class = "db.r6g.2xlarge"  # 8 vCPUs, 64 GB RAM
redis_node_type    = "cache.r6g.xlarge" # 4 vCPUs, 26.32 GB RAM

# Gaming Analytics Settings
gaming_analytics_enabled      = true
riot_api_rate_limit           = 100
max_concurrent_users          = 1000000
analytics_processing_timeout  = 5

# Scaling Configuration - Production Scale
min_api_replicas       = 5
max_api_replicas       = 100
target_cpu_utilization = 65

# Monitoring Configuration
enable_monitoring      = true
log_retention_days     = 90
metrics_retention_days = 365

# Security Configuration
enable_waf              = true
enable_guardduty        = true
enable_secrets_rotation = true

# Backup Configuration
backup_retention_days         = 30
enable_point_in_time_recovery = true

# Cost Optimization
enable_spot_instances = false  # Don't use spot for production
spot_max_price        = ""

# Tags
tags = {
  Project     = "Herald.lol"
  Environment = "Production"
  Platform    = "Gaming Analytics"
  Stack       = "Production"
  Performance = "Critical"
  Compliance  = "GDPR"
  SLA         = "99.9"
}