# Herald.lol Gaming Analytics - Infrastructure Drift Detection
# Terraform configuration for gaming platform infrastructure monitoring

terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }

  backend "s3" {
    bucket         = "herald-gaming-terraform-state"
    key            = "drift-detection/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "herald-terraform-locks"
  }
}

# Variables for Herald.lol gaming infrastructure
variable "gaming_environment" {
  description = "Gaming environment (blue/green/production)"
  type        = string
  default     = "production"
}

variable "gaming_performance_target_ms" {
  description = "Gaming analytics performance target in milliseconds"
  type        = number
  default     = 5000
}

variable "gaming_concurrent_users" {
  description = "Target concurrent users for gaming platform"
  type        = number
  default     = 1000000
}

variable "aws_region" {
  description = "AWS region for gaming infrastructure"
  type        = string
  default     = "us-east-1"
}

variable "notification_email" {
  description = "Email for gaming infrastructure drift alerts"
  type        = string
}

variable "slack_webhook_url" {
  description = "Slack webhook for gaming alerts"
  type        = string
  sensitive   = true
}

# Data sources
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

# Local values for gaming infrastructure
locals {
  gaming_tags = {
    Project             = "Herald.lol"
    Environment         = var.gaming_environment
    Component           = "gaming-analytics"
    PerformanceTarget   = "${var.gaming_performance_target_ms}ms"
    ConcurrentUsers     = var.gaming_concurrent_users
    ManagedBy          = "terraform"
    DriftMonitoring    = "enabled"
  }

  drift_detection_name = "herald-gaming-drift-detection"
}

# SNS Topic for gaming infrastructure alerts
resource "aws_sns_topic" "gaming_infrastructure_alerts" {
  name = "herald-gaming-infrastructure-alerts"

  tags = merge(local.gaming_tags, {
    Name        = "Herald Gaming Infrastructure Alerts"
    Description = "SNS topic for gaming platform drift detection alerts"
  })
}

# SNS Topic Subscription for email alerts
resource "aws_sns_topic_subscription" "gaming_email_alerts" {
  topic_arn = aws_sns_topic.gaming_infrastructure_alerts.arn
  protocol  = "email"
  endpoint  = var.notification_email
}

# CloudWatch Log Group for drift detection logs
resource "aws_cloudwatch_log_group" "drift_detection_logs" {
  name              = "/aws/lambda/${local.drift_detection_name}"
  retention_in_days = 30

  tags = merge(local.gaming_tags, {
    Name        = "Herald Gaming Drift Detection Logs"
    LogType     = "infrastructure-drift"
  })
}

# IAM Role for drift detection Lambda
resource "aws_iam_role" "drift_detection_role" {
  name = "${local.drift_detection_name}-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = local.gaming_tags
}

# IAM Policy for drift detection
resource "aws_iam_policy" "drift_detection_policy" {
  name        = "${local.drift_detection_name}-policy"
  description = "Policy for Herald.lol gaming infrastructure drift detection"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "${aws_cloudwatch_log_group.drift_detection_logs.arn}:*"
      },
      {
        Effect = "Allow"
        Action = [
          "ec2:Describe*",
          "rds:Describe*",
          "elasticache:Describe*",
          "elb:Describe*",
          "elbv2:Describe*",
          "autoscaling:Describe*",
          "cloudformation:Describe*",
          "cloudformation:List*",
          "eks:Describe*",
          "eks:List*"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject"
        ]
        Resource = [
          "arn:aws:s3:::herald-gaming-terraform-state/*",
          "arn:aws:s3:::herald-gaming-infrastructure-backups/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "sns:Publish"
        ]
        Resource = aws_sns_topic.gaming_infrastructure_alerts.arn
      },
      {
        Effect = "Allow"
        Action = [
          "ssm:GetParameter",
          "ssm:GetParameters",
          "ssm:GetParametersByPath"
        ]
        Resource = "arn:aws:ssm:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:parameter/herald/gaming/*"
      }
    ]
  })
}

# Attach policy to role
resource "aws_iam_role_policy_attachment" "drift_detection_policy_attachment" {
  role       = aws_iam_role.drift_detection_role.name
  policy_arn = aws_iam_policy.drift_detection_policy.arn
}

# S3 Bucket for storing infrastructure state snapshots
resource "aws_s3_bucket" "infrastructure_snapshots" {
  bucket = "herald-gaming-infrastructure-snapshots"

  tags = merge(local.gaming_tags, {
    Name        = "Herald Gaming Infrastructure Snapshots"
    Purpose     = "drift-detection-storage"
  })
}

resource "aws_s3_bucket_versioning" "infrastructure_snapshots_versioning" {
  bucket = aws_s3_bucket.infrastructure_snapshots.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "infrastructure_snapshots_encryption" {
  bucket = aws_s3_bucket.infrastructure_snapshots.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# Lambda function for drift detection (placeholder - actual code would be deployed separately)
resource "aws_lambda_function" "drift_detection" {
  filename         = "drift_detection.zip"
  function_name    = local.drift_detection_name
  role            = aws_iam_role.drift_detection_role.arn
  handler         = "main.handler"
  runtime         = "python3.11"
  timeout         = 300

  environment {
    variables = {
      GAMING_ENVIRONMENT           = var.gaming_environment
      GAMING_PERFORMANCE_TARGET_MS = var.gaming_performance_target_ms
      GAMING_CONCURRENT_USERS      = var.gaming_concurrent_users
      SNS_TOPIC_ARN               = aws_sns_topic.gaming_infrastructure_alerts.arn
      S3_SNAPSHOTS_BUCKET         = aws_s3_bucket.infrastructure_snapshots.bucket
      SLACK_WEBHOOK_URL           = var.slack_webhook_url
    }
  }

  tags = merge(local.gaming_tags, {
    Name        = "Herald Gaming Drift Detection"
    Function    = "infrastructure-monitoring"
  })

  depends_on = [aws_cloudwatch_log_group.drift_detection_logs]
}

# EventBridge rule for scheduled drift detection
resource "aws_cloudwatch_event_rule" "drift_detection_schedule" {
  name                = "herald-gaming-drift-detection-schedule"
  description         = "Schedule drift detection for Herald.lol gaming infrastructure"
  schedule_expression = "rate(4 hours)"  # Run every 4 hours for gaming infrastructure

  tags = local.gaming_tags
}

# EventBridge target for Lambda
resource "aws_cloudwatch_event_target" "drift_detection_target" {
  rule      = aws_cloudwatch_event_rule.drift_detection_schedule.name
  target_id = "DriftDetectionTarget"
  arn       = aws_lambda_function.drift_detection.arn
}

# Lambda permission for EventBridge
resource "aws_lambda_permission" "allow_eventbridge" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.drift_detection.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.drift_detection_schedule.arn
}

# CloudWatch Dashboard for gaming infrastructure monitoring
resource "aws_cloudwatch_dashboard" "gaming_infrastructure_dashboard" {
  dashboard_name = "Herald-Gaming-Infrastructure-Drift"

  dashboard_body = jsonencode({
    widgets = [
      {
        type   = "metric"
        x      = 0
        y      = 0
        width  = 12
        height = 6

        properties = {
          metrics = [
            ["AWS/Lambda", "Duration", "FunctionName", local.drift_detection_name],
            [".", "Errors", ".", "."],
            [".", "Invocations", ".", "."]
          ]
          view    = "timeSeries"
          stacked = false
          region  = data.aws_region.current.name
          title   = "Herald Gaming Drift Detection Performance"
          period  = 300
        }
      },
      {
        type   = "log"
        x      = 0
        y      = 6
        width  = 24
        height = 6

        properties = {
          query   = "SOURCE '/aws/lambda/${local.drift_detection_name}'\n| fields @timestamp, @message\n| filter @message like /DRIFT_DETECTED|GAMING_IMPACT/\n| sort @timestamp desc\n| limit 100"
          region  = data.aws_region.current.name
          title   = "Herald Gaming Infrastructure Drift Events"
        }
      }
    ]
  })

  tags = local.gaming_tags
}

# SSM Parameters for gaming infrastructure configuration
resource "aws_ssm_parameter" "gaming_performance_target" {
  name  = "/herald/gaming/performance/target_ms"
  type  = "String"
  value = var.gaming_performance_target_ms

  tags = merge(local.gaming_tags, {
    Name = "Gaming Performance Target"
  })
}

resource "aws_ssm_parameter" "gaming_concurrent_users" {
  name  = "/herald/gaming/capacity/concurrent_users"
  type  = "String"
  value = var.gaming_concurrent_users

  tags = merge(local.gaming_tags, {
    Name = "Gaming Concurrent Users Target"
  })
}

# Outputs
output "drift_detection_function_arn" {
  description = "ARN of the gaming infrastructure drift detection Lambda function"
  value       = aws_lambda_function.drift_detection.arn
}

output "sns_topic_arn" {
  description = "ARN of the SNS topic for gaming infrastructure alerts"
  value       = aws_sns_topic.gaming_infrastructure_alerts.arn
}

output "dashboard_url" {
  description = "URL of the gaming infrastructure monitoring dashboard"
  value       = "https://${data.aws_region.current.name}.console.aws.amazon.com/cloudwatch/home?region=${data.aws_region.current.name}#dashboards:name=${aws_cloudwatch_dashboard.gaming_infrastructure_dashboard.dashboard_name}"
}

output "snapshots_bucket" {
  description = "S3 bucket for gaming infrastructure snapshots"
  value       = aws_s3_bucket.infrastructure_snapshots.bucket
}