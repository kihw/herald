# Herald.lol Gaming Analytics Platform - Terraform Outputs
# Output values for Herald gaming infrastructure

# EKS Cluster Outputs
output "cluster_id" {
  description = "EKS cluster ID"
  value       = module.eks.cluster_id
}

output "cluster_arn" {
  description = "EKS cluster ARN"
  value       = module.eks.cluster_arn
}

output "cluster_endpoint" {
  description = "Endpoint for EKS control plane"
  value       = module.eks.cluster_endpoint
}

output "cluster_version" {
  description = "The Kubernetes version for the EKS cluster"
  value       = module.eks.cluster_version
}

output "cluster_platform_version" {
  description = "Platform version for the EKS cluster"
  value       = module.eks.cluster_platform_version
}

output "cluster_status" {
  description = "Status of the EKS cluster"
  value       = module.eks.cluster_status
}

output "cluster_security_group_id" {
  description = "Security group ID attached to the EKS cluster control plane"
  value       = module.eks.cluster_security_group_id
}

output "cluster_security_group_arn" {
  description = "Amazon Resource Name (ARN) of the cluster security group"
  value       = module.eks.cluster_security_group_arn
}

output "cluster_iam_role_name" {
  description = "IAM role name associated with EKS cluster"
  value       = module.eks.cluster_iam_role_name
}

output "cluster_iam_role_arn" {
  description = "IAM role ARN associated with EKS cluster"
  value       = module.eks.cluster_iam_role_arn
}

output "cluster_oidc_issuer_url" {
  description = "The URL on the EKS cluster for the OpenID Connect identity provider"
  value       = module.eks.cluster_oidc_issuer_url
}

output "cluster_certificate_authority_data" {
  description = "Base64 encoded certificate data required to communicate with the cluster"
  value       = module.eks.cluster_certificate_authority_data
  sensitive   = true
}

# Node Groups Outputs
output "eks_managed_node_groups" {
  description = "Map of attribute maps for all EKS managed node groups created"
  value       = module.eks.eks_managed_node_groups
}

output "eks_managed_node_groups_autoscaling_group_names" {
  description = "List of the autoscaling group names created by EKS managed node groups"
  value       = module.eks.eks_managed_node_groups_autoscaling_group_names
}

# VPC Outputs
output "vpc_id" {
  description = "ID of the VPC where Herald gaming infrastructure is deployed"
  value       = module.vpc.vpc_id
}

output "vpc_arn" {
  description = "The ARN of the VPC"
  value       = module.vpc.vpc_arn
}

output "vpc_cidr_block" {
  description = "The CIDR block of the VPC"
  value       = module.vpc.vpc_cidr_block
}

output "private_subnets" {
  description = "List of IDs of private subnets"
  value       = module.vpc.private_subnets
}

output "public_subnets" {
  description = "List of IDs of public subnets"
  value       = module.vpc.public_subnets
}

output "private_subnet_arns" {
  description = "List of ARNs of private subnets"
  value       = module.vpc.private_subnet_arns
}

output "public_subnet_arns" {
  description = "List of ARNs of public subnets"
  value       = module.vpc.public_subnet_arns
}

output "nat_ids" {
  description = "List of allocation IDs of NAT Gateways"
  value       = module.vpc.nat_ids
}

output "nat_public_ips" {
  description = "List of public Elastic IPs created for AWS NAT Gateway"
  value       = module.vpc.nat_public_ips
}

output "internet_gateway_id" {
  description = "The ID of the Internet Gateway"
  value       = module.vpc.internet_gateway_id
}

output "internet_gateway_arn" {
  description = "The ARN of the Internet Gateway"
  value       = module.vpc.internet_gateway_arn
}

# RDS Outputs
output "rds_cluster_id" {
  description = "The RDS cluster ID"
  value       = aws_rds_cluster.herald_gaming_db.cluster_identifier
}

output "rds_cluster_arn" {
  description = "Amazon Resource Name (ARN) of cluster"
  value       = aws_rds_cluster.herald_gaming_db.arn
}

output "rds_cluster_endpoint" {
  description = "RDS cluster endpoint for gaming database writes"
  value       = aws_rds_cluster.herald_gaming_db.endpoint
}

output "rds_cluster_reader_endpoint" {
  description = "RDS cluster reader endpoint for gaming database reads"
  value       = aws_rds_cluster.herald_gaming_db.reader_endpoint
}

output "rds_cluster_port" {
  description = "The database port"
  value       = aws_rds_cluster.herald_gaming_db.port
}

output "rds_cluster_database_name" {
  description = "The database name"
  value       = aws_rds_cluster.herald_gaming_db.database_name
}

output "rds_cluster_master_username" {
  description = "The database master username"
  value       = aws_rds_cluster.herald_gaming_db.master_username
  sensitive   = true
}

output "rds_cluster_hosted_zone_id" {
  description = "The Route53 Hosted Zone ID of the endpoint"
  value       = aws_rds_cluster.herald_gaming_db.hosted_zone_id
}

output "rds_cluster_instances" {
  description = "A map of cluster instances and their attributes"
  value = {
    primary = {
      identifier   = aws_rds_cluster_instance.herald_gaming_primary[0].identifier
      endpoint     = aws_rds_cluster_instance.herald_gaming_primary[0].endpoint
      arn          = aws_rds_cluster_instance.herald_gaming_primary[0].arn
    }
    readers = [
      for instance in aws_rds_cluster_instance.herald_gaming_readers : {
        identifier   = instance.identifier
        endpoint     = instance.endpoint
        arn          = instance.arn
      }
    ]
  }
}

# ElastiCache Redis Outputs
output "redis_cluster_id" {
  description = "The Redis cluster ID"
  value       = aws_elasticache_replication_group.herald_gaming_redis.replication_group_id
}

output "redis_cluster_arn" {
  description = "The Amazon Resource Name (ARN) of the Redis cluster"
  value       = aws_elasticache_replication_group.herald_gaming_redis.arn
}

output "redis_configuration_endpoint" {
  description = "Redis cluster configuration endpoint for gaming cache"
  value       = aws_elasticache_replication_group.herald_gaming_redis.configuration_endpoint_address
}

output "redis_primary_endpoint" {
  description = "Redis primary endpoint address"
  value       = aws_elasticache_replication_group.herald_gaming_redis.primary_endpoint_address
}

output "redis_reader_endpoint" {
  description = "Redis reader endpoint address"
  value       = aws_elasticache_replication_group.herald_gaming_redis.reader_endpoint_address
}

output "redis_port" {
  description = "Redis port"
  value       = aws_elasticache_replication_group.herald_gaming_redis.port
}

# S3 Outputs
output "s3_gaming_assets_bucket_id" {
  description = "The name of the gaming assets bucket"
  value       = aws_s3_bucket.herald_gaming_assets.id
}

output "s3_gaming_assets_bucket_arn" {
  description = "The ARN of the gaming assets bucket"
  value       = aws_s3_bucket.herald_gaming_assets.arn
}

output "s3_gaming_assets_bucket_domain_name" {
  description = "The bucket domain name for gaming assets"
  value       = aws_s3_bucket.herald_gaming_assets.bucket_domain_name
}

output "s3_gaming_assets_bucket_regional_domain_name" {
  description = "The bucket regional domain name for gaming assets"
  value       = aws_s3_bucket.herald_gaming_assets.bucket_regional_domain_name
}

output "s3_gaming_backups_bucket_id" {
  description = "The name of the gaming backups bucket"
  value       = aws_s3_bucket.herald_gaming_backups.id
}

output "s3_gaming_backups_bucket_arn" {
  description = "The ARN of the gaming backups bucket"
  value       = aws_s3_bucket.herald_gaming_backups.arn
}

# CloudFront Outputs
output "cloudfront_distribution_id" {
  description = "CloudFront distribution ID for gaming CDN"
  value       = aws_cloudfront_distribution.herald_gaming_cdn.id
}

output "cloudfront_distribution_arn" {
  description = "ARN for the CloudFront distribution"
  value       = aws_cloudfront_distribution.herald_gaming_cdn.arn
}

output "cloudfront_distribution_caller_reference" {
  description = "Internal value used by CloudFront to allow future updates"
  value       = aws_cloudfront_distribution.herald_gaming_cdn.caller_reference
}

output "cloudfront_distribution_status" {
  description = "Current status of the distribution"
  value       = aws_cloudfront_distribution.herald_gaming_cdn.status
}

output "cloudfront_distribution_domain_name" {
  description = "CloudFront domain name for global gaming content delivery"
  value       = aws_cloudfront_distribution.herald_gaming_cdn.domain_name
}

output "cloudfront_distribution_etag" {
  description = "Current version of the distribution's information"
  value       = aws_cloudfront_distribution.herald_gaming_cdn.etag
}

output "cloudfront_distribution_hosted_zone_id" {
  description = "CloudFront Route 53 zone ID"
  value       = aws_cloudfront_distribution.herald_gaming_cdn.hosted_zone_id
}

# Security Outputs
output "rds_security_group_id" {
  description = "ID of the RDS security group"
  value       = aws_security_group.rds_gaming.id
}

output "rds_security_group_arn" {
  description = "ARN of the RDS security group"
  value       = aws_security_group.rds_gaming.arn
}

output "elasticache_security_group_id" {
  description = "ID of the ElastiCache security group"
  value       = aws_security_group.elasticache_gaming.id
}

output "elasticache_security_group_arn" {
  description = "ARN of the ElastiCache security group"
  value       = aws_security_group.elasticache_gaming.arn
}

# KMS Outputs
output "kms_key_id" {
  description = "The globally unique identifier for the KMS key"
  value       = aws_kms_key.gaming_db.key_id
}

output "kms_key_arn" {
  description = "The Amazon Resource Name (ARN) of the KMS key"
  value       = aws_kms_key.gaming_db.arn
}

output "kms_alias_arn" {
  description = "The Amazon Resource Name (ARN) of the KMS key alias"
  value       = aws_kms_alias.gaming_db.arn
}

# Gaming Infrastructure Summary
output "gaming_infrastructure_summary" {
  description = "Summary of Herald gaming infrastructure endpoints and configuration"
  value = {
    cluster_name    = var.cluster_name
    cluster_endpoint = module.eks.cluster_endpoint
    database_endpoint = aws_rds_cluster.herald_gaming_db.endpoint
    redis_endpoint  = aws_elasticache_replication_group.herald_gaming_redis.configuration_endpoint_address
    cdn_domain     = aws_cloudfront_distribution.herald_gaming_cdn.domain_name
    region         = var.aws_region
    environment    = var.environment
  }
  sensitive = false
}

# Authentication and Connection Info
output "kubectl_config_command" {
  description = "Command to configure kubectl for Herald gaming cluster"
  value       = "aws eks update-kubeconfig --region ${var.aws_region} --name ${module.eks.cluster_name}"
}

output "aws_region" {
  description = "AWS region where Herald gaming infrastructure is deployed"
  value       = var.aws_region
}

output "environment" {
  description = "Environment name"
  value       = var.environment
}

# Infrastructure Health Checks
output "infrastructure_health_endpoints" {
  description = "Endpoints for monitoring Herald gaming infrastructure health"
  value = {
    eks_health          = "${module.eks.cluster_endpoint}/healthz"
    rds_health_check    = aws_rds_cluster.herald_gaming_db.endpoint
    redis_health_check  = aws_elasticache_replication_group.herald_gaming_redis.configuration_endpoint_address
  }
}