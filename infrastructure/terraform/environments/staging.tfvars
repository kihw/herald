# Herald.lol Staging Environment Configuration

environment = "staging"
aws_region  = "us-west-2"

# Network Configuration
vpc_cidr             = "10.1.0.0/16"
public_subnet_cidrs  = ["10.1.1.0/24", "10.1.2.0/24", "10.1.3.0/24"]
private_subnet_cidrs = ["10.1.11.0/24", "10.1.12.0/24", "10.1.13.0/24"]

# Kubernetes Configuration
kubernetes_version = "1.28"

# Database Configuration - Staging (Smaller than Production)
rds_instance_class = "db.t4g.large"    # 2 vCPUs, 8 GB RAM
redis_node_type    = "cache.t4g.medium" # 2 vCPUs, 3.09 GB RAM

# Gaming Analytics Settings
gaming_analytics_enabled      = true
riot_api_rate_limit           = 100
max_concurrent_users          = 10000
analytics_processing_timeout  = 5

# Scaling Configuration - Staging Scale
min_api_replicas       = 2
max_api_replicas       = 10
target_cpu_utilization = 70

# Monitoring Configuration
enable_monitoring      = true
log_retention_days     = 30
metrics_retention_days = 90

# Security Configuration
enable_waf              = true
enable_guardduty        = true
enable_secrets_rotation = true

# Backup Configuration
backup_retention_days         = 7
enable_point_in_time_recovery = true

# Cost Optimization
enable_spot_instances = true  # Use spot instances for staging
spot_max_price        = "0.05"

# Tags
tags = {
  Project     = "Herald.lol"
  Environment = "Staging"
  Platform    = "Gaming Analytics"
  Stack       = "Staging"
  Performance = "Standard"
  Purpose     = "Testing"
}