#!/usr/bin/env python3
"""
Herald.lol Gaming Analytics - Infrastructure Drift Detection
AWS Lambda function for monitoring gaming infrastructure drift
"""

import json
import boto3
import os
import logging
import hashlib
from datetime import datetime, timezone
from typing import Dict, List, Any, Optional
import requests

# Configure logging for gaming platform
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - [HERALD-GAMING] - %(message)s'
)
logger = logging.getLogger(__name__)

class HeraldGamingDriftDetector:
    """Drift detection optimized for Herald.lol gaming infrastructure"""
    
    def __init__(self):
        self.gaming_environment = os.environ.get('GAMING_ENVIRONMENT', 'production')
        self.performance_target_ms = int(os.environ.get('GAMING_PERFORMANCE_TARGET_MS', '5000'))
        self.concurrent_users = int(os.environ.get('GAMING_CONCURRENT_USERS', '1000000'))
        self.sns_topic_arn = os.environ.get('SNS_TOPIC_ARN')
        self.snapshots_bucket = os.environ.get('S3_SNAPSHOTS_BUCKET')
        self.slack_webhook_url = os.environ.get('SLACK_WEBHOOK_URL')
        
        # AWS clients
        self.ec2 = boto3.client('ec2')
        self.rds = boto3.client('rds')
        self.elasticache = boto3.client('elasticache')
        self.elbv2 = boto3.client('elbv2')
        self.autoscaling = boto3.client('autoscaling')
        self.eks = boto3.client('eks')
        self.s3 = boto3.client('s3')
        self.sns = boto3.client('sns')
        self.ssm = boto3.client('ssm')
        
        logger.info(f"🎮 Herald.lol Drift Detector initialized for {self.gaming_environment}")
        logger.info(f"🎯 Gaming performance target: <{self.performance_target_ms}ms")
        logger.info(f"👥 Concurrent users target: {self.concurrent_users:,}")

    def get_gaming_infrastructure_state(self) -> Dict[str, Any]:
        """Capture current state of gaming infrastructure"""
        logger.info("📊 Capturing Herald.lol gaming infrastructure state...")
        
        state = {
            'timestamp': datetime.now(timezone.utc).isoformat(),
            'gaming_environment': self.gaming_environment,
            'performance_target_ms': self.performance_target_ms,
            'concurrent_users_target': self.concurrent_users,
            'infrastructure': {}
        }
        
        try:
            # Gaming EKS clusters
            state['infrastructure']['eks_clusters'] = self._get_gaming_eks_state()
            
            # Gaming databases (PostgreSQL for analytics)
            state['infrastructure']['rds_instances'] = self._get_gaming_rds_state()
            
            # Gaming cache (Redis for performance)
            state['infrastructure']['elasticache'] = self._get_gaming_cache_state()
            
            # Gaming load balancers
            state['infrastructure']['load_balancers'] = self._get_gaming_elb_state()
            
            # Gaming auto scaling groups
            state['infrastructure']['auto_scaling_groups'] = self._get_gaming_asg_state()
            
            # Gaming EC2 instances
            state['infrastructure']['ec2_instances'] = self._get_gaming_ec2_state()
            
            logger.info("✅ Gaming infrastructure state captured successfully")
            
        except Exception as e:
            logger.error(f"❌ Failed to capture gaming infrastructure state: {e}")
            raise
            
        return state

    def _get_gaming_eks_state(self) -> List[Dict[str, Any]]:
        """Get EKS clusters for gaming platform"""
        clusters = []
        
        try:
            response = self.eks.list_clusters()
            
            for cluster_name in response['clusters']:
                if 'herald' in cluster_name.lower() or 'gaming' in cluster_name.lower():
                    cluster_detail = self.eks.describe_cluster(name=cluster_name)
                    cluster = cluster_detail['cluster']
                    
                    clusters.append({
                        'name': cluster_name,
                        'status': cluster['status'],
                        'version': cluster['version'],
                        'endpoint': cluster['endpoint'],
                        'platform_version': cluster['platformVersion'],
                        'gaming_optimized': True,
                        'performance_critical': True
                    })
                    
        except Exception as e:
            logger.warning(f"⚠️ Failed to get EKS clusters: {e}")
            
        return clusters

    def _get_gaming_rds_state(self) -> List[Dict[str, Any]]:
        """Get RDS instances for gaming analytics"""
        instances = []
        
        try:
            response = self.rds.describe_db_instances()
            
            for instance in response['DBInstances']:
                db_id = instance['DBInstanceIdentifier']
                
                if 'herald' in db_id.lower() or 'gaming' in db_id.lower():
                    instances.append({
                        'identifier': db_id,
                        'engine': instance['Engine'],
                        'engine_version': instance['EngineVersion'],
                        'instance_class': instance['DBInstanceClass'],
                        'status': instance['DBInstanceStatus'],
                        'multi_az': instance['MultiAZ'],
                        'storage_encrypted': instance['StorageEncrypted'],
                        'gaming_analytics_db': True,
                        'performance_optimized': instance['DBInstanceClass'].startswith('r5') or instance['DBInstanceClass'].startswith('r6')
                    })
                    
        except Exception as e:
            logger.warning(f"⚠️ Failed to get RDS instances: {e}")
            
        return instances

    def _get_gaming_cache_state(self) -> List[Dict[str, Any]]:
        """Get ElastiCache clusters for gaming performance"""
        clusters = []
        
        try:
            response = self.elasticache.describe_cache_clusters(ShowCacheNodeInfo=True)
            
            for cluster in response['CacheClusters']:
                cluster_id = cluster['CacheClusterId']
                
                if 'herald' in cluster_id.lower() or 'gaming' in cluster_id.lower():
                    clusters.append({
                        'cluster_id': cluster_id,
                        'engine': cluster['Engine'],
                        'engine_version': cluster['EngineVersion'],
                        'node_type': cluster['CacheNodeType'],
                        'num_cache_nodes': cluster['NumCacheNodes'],
                        'status': cluster['CacheClusterStatus'],
                        'gaming_cache': True,
                        'performance_critical': True
                    })
                    
        except Exception as e:
            logger.warning(f"⚠️ Failed to get ElastiCache clusters: {e}")
            
        return clusters

    def _get_gaming_elb_state(self) -> List[Dict[str, Any]]:
        """Get load balancers for gaming traffic"""
        load_balancers = []
        
        try:
            response = self.elbv2.describe_load_balancers()
            
            for lb in response['LoadBalancers']:
                lb_name = lb['LoadBalancerName']
                
                if 'herald' in lb_name.lower() or 'gaming' in lb_name.lower():
                    load_balancers.append({
                        'name': lb_name,
                        'arn': lb['LoadBalancerArn'],
                        'type': lb['Type'],
                        'scheme': lb['Scheme'],
                        'state': lb['State']['Code'],
                        'availability_zones': len(lb['AvailabilityZones']),
                        'gaming_traffic': True,
                        'high_availability': len(lb['AvailabilityZones']) >= 2
                    })
                    
        except Exception as e:
            logger.warning(f"⚠️ Failed to get load balancers: {e}")
            
        return load_balancers

    def _get_gaming_asg_state(self) -> List[Dict[str, Any]]:
        """Get Auto Scaling Groups for gaming workloads"""
        asgs = []
        
        try:
            response = self.autoscaling.describe_auto_scaling_groups()
            
            for asg in response['AutoScalingGroups']:
                asg_name = asg['AutoScalingGroupName']
                
                if 'herald' in asg_name.lower() or 'gaming' in asg_name.lower():
                    asgs.append({
                        'name': asg_name,
                        'min_size': asg['MinSize'],
                        'max_size': asg['MaxSize'],
                        'desired_capacity': asg['DesiredCapacity'],
                        'instances': len(asg['Instances']),
                        'health_check_type': asg['HealthCheckType'],
                        'gaming_workload': True,
                        'auto_scaling_enabled': asg['MaxSize'] > asg['MinSize']
                    })
                    
        except Exception as e:
            logger.warning(f"⚠️ Failed to get Auto Scaling Groups: {e}")
            
        return asgs

    def _get_gaming_ec2_state(self) -> List[Dict[str, Any]]:
        """Get EC2 instances for gaming platform"""
        instances = []
        
        try:
            response = self.ec2.describe_instances(
                Filters=[
                    {'Name': 'instance-state-name', 'Values': ['running']},
                    {'Name': 'tag:Project', 'Values': ['Herald.lol', 'herald', 'gaming']}
                ]
            )
            
            for reservation in response['Reservations']:
                for instance in reservation['Instances']:
                    instance_name = ''
                    gaming_component = ''
                    
                    for tag in instance.get('Tags', []):
                        if tag['Key'] == 'Name':
                            instance_name = tag['Value']
                        elif tag['Key'] == 'Component':
                            gaming_component = tag['Value']
                    
                    instances.append({
                        'instance_id': instance['InstanceId'],
                        'instance_type': instance['InstanceType'],
                        'state': instance['State']['Name'],
                        'name': instance_name,
                        'gaming_component': gaming_component,
                        'gaming_optimized': instance['InstanceType'].startswith(('c5', 'c6', 'm5', 'm6', 'r5', 'r6')),
                        'performance_instance': instance['InstanceType'].startswith(('c5', 'c6'))
                    })
                    
        except Exception as e:
            logger.warning(f"⚠️ Failed to get EC2 instances: {e}")
            
        return instances

    def save_state_snapshot(self, state: Dict[str, Any]) -> str:
        """Save infrastructure state snapshot to S3"""
        timestamp = datetime.now(timezone.utc).strftime('%Y%m%d_%H%M%S')
        key = f"gaming-infrastructure/{self.gaming_environment}/state_{timestamp}.json"
        
        try:
            self.s3.put_object(
                Bucket=self.snapshots_bucket,
                Key=key,
                Body=json.dumps(state, indent=2, default=str),
                ContentType='application/json',
                Metadata={
                    'gaming-environment': self.gaming_environment,
                    'performance-target-ms': str(self.performance_target_ms),
                    'snapshot-type': 'infrastructure-state'
                }
            )
            
            logger.info(f"📸 Gaming infrastructure snapshot saved: s3://{self.snapshots_bucket}/{key}")
            return key
            
        except Exception as e:
            logger.error(f"❌ Failed to save state snapshot: {e}")
            raise

    def get_previous_state(self) -> Optional[Dict[str, Any]]:
        """Get the most recent previous state snapshot"""
        try:
            response = self.s3.list_objects_v2(
                Bucket=self.snapshots_bucket,
                Prefix=f"gaming-infrastructure/{self.gaming_environment}/",
                MaxKeys=2  # Get current and previous
            )
            
            if 'Contents' not in response or len(response['Contents']) < 2:
                logger.info("📋 No previous state snapshot found for comparison")
                return None
                
            # Sort by last modified, get second most recent
            objects = sorted(response['Contents'], key=lambda x: x['LastModified'], reverse=True)
            previous_key = objects[1]['Key']
            
            # Get previous state
            obj = self.s3.get_object(Bucket=self.snapshots_bucket, Key=previous_key)
            previous_state = json.loads(obj['Body'].read())
            
            logger.info(f"📋 Retrieved previous gaming infrastructure state: {previous_key}")
            return previous_state
            
        except Exception as e:
            logger.warning(f"⚠️ Failed to get previous state: {e}")
            return None

    def detect_gaming_drift(self, current_state: Dict[str, Any], previous_state: Dict[str, Any]) -> Dict[str, Any]:
        """Detect drift in gaming infrastructure with gaming-specific impact analysis"""
        logger.info("🔍 Analyzing gaming infrastructure drift...")
        
        drift_report = {
            'timestamp': datetime.now(timezone.utc).isoformat(),
            'gaming_environment': self.gaming_environment,
            'drift_detected': False,
            'gaming_impact': 'none',
            'critical_changes': [],
            'performance_impact': [],
            'gaming_specific_changes': [],
            'detailed_changes': {}
        }
        
        # Compare each infrastructure component
        for component in ['eks_clusters', 'rds_instances', 'elasticache', 'load_balancers', 'auto_scaling_groups', 'ec2_instances']:
            current_component = current_state['infrastructure'].get(component, [])
            previous_component = previous_state['infrastructure'].get(component, [])
            
            component_changes = self._compare_component(component, current_component, previous_component)
            
            if component_changes['changes_detected']:
                drift_report['drift_detected'] = True
                drift_report['detailed_changes'][component] = component_changes
                
                # Analyze gaming impact
                gaming_impact = self._analyze_gaming_impact(component, component_changes)
                
                if gaming_impact['critical']:
                    drift_report['critical_changes'].extend(gaming_impact['critical_changes'])
                    drift_report['gaming_impact'] = 'critical'
                    
                if gaming_impact['performance_impact']:
                    drift_report['performance_impact'].extend(gaming_impact['performance_changes'])
                    if drift_report['gaming_impact'] != 'critical':
                        drift_report['gaming_impact'] = 'performance'
                        
                if gaming_impact['gaming_specific']:
                    drift_report['gaming_specific_changes'].extend(gaming_impact['gaming_changes'])
                    if drift_report['gaming_impact'] == 'none':
                        drift_report['gaming_impact'] = 'moderate'
        
        # Calculate overall gaming risk score
        drift_report['gaming_risk_score'] = self._calculate_gaming_risk_score(drift_report)
        
        if drift_report['drift_detected']:
            logger.warning(f"🚨 Gaming infrastructure drift detected! Impact: {drift_report['gaming_impact']}")
            logger.warning(f"🎯 Gaming risk score: {drift_report['gaming_risk_score']}/10")
        else:
            logger.info("✅ No gaming infrastructure drift detected")
            
        return drift_report

    def _compare_component(self, component_name: str, current: List[Dict], previous: List[Dict]) -> Dict[str, Any]:
        """Compare infrastructure component states"""
        changes = {
            'changes_detected': False,
            'added': [],
            'removed': [],
            'modified': []
        }
        
        # Create lookup dictionaries
        current_dict = {}
        previous_dict = {}
        
        # Use appropriate key for each component type
        key_mapping = {
            'eks_clusters': 'name',
            'rds_instances': 'identifier',
            'elasticache': 'cluster_id',
            'load_balancers': 'name',
            'auto_scaling_groups': 'name',
            'ec2_instances': 'instance_id'
        }
        
        key_field = key_mapping.get(component_name, 'name')
        
        for item in current:
            current_dict[item[key_field]] = item
            
        for item in previous:
            previous_dict[item[key_field]] = item
        
        # Find added items
        for key in current_dict:
            if key not in previous_dict:
                changes['added'].append(current_dict[key])
                changes['changes_detected'] = True
        
        # Find removed items
        for key in previous_dict:
            if key not in current_dict:
                changes['removed'].append(previous_dict[key])
                changes['changes_detected'] = True
        
        # Find modified items
        for key in current_dict:
            if key in previous_dict:
                if self._hash_object(current_dict[key]) != self._hash_object(previous_dict[key]):
                    changes['modified'].append({
                        'key': key,
                        'current': current_dict[key],
                        'previous': previous_dict[key]
                    })
                    changes['changes_detected'] = True
        
        return changes

    def _hash_object(self, obj: Dict[str, Any]) -> str:
        """Create hash of object for comparison"""
        return hashlib.md5(json.dumps(obj, sort_keys=True).encode()).hexdigest()

    def _analyze_gaming_impact(self, component: str, changes: Dict[str, Any]) -> Dict[str, Any]:
        """Analyze gaming-specific impact of infrastructure changes"""
        impact = {
            'critical': False,
            'performance_impact': False,
            'gaming_specific': False,
            'critical_changes': [],
            'performance_changes': [],
            'gaming_changes': []
        }
        
        # Gaming-critical components
        gaming_critical = ['eks_clusters', 'elasticache', 'load_balancers']
        performance_critical = ['rds_instances', 'auto_scaling_groups', 'ec2_instances']
        
        if component in gaming_critical:
            impact['gaming_specific'] = True
            
            # Check for critical gaming changes
            for added in changes['added']:
                if added.get('gaming_optimized') or added.get('performance_critical'):
                    impact['critical'] = True
                    impact['critical_changes'].append(f"New gaming-critical {component}: {added.get('name', 'unnamed')}")
            
            for removed in changes['removed']:
                if removed.get('gaming_optimized') or removed.get('performance_critical'):
                    impact['critical'] = True
                    impact['critical_changes'].append(f"Removed gaming-critical {component}: {removed.get('name', 'unnamed')}")
            
            for modified in changes['modified']:
                current = modified['current']
                previous = modified['previous']
                
                # Check for performance-impacting changes
                if component == 'eks_clusters' and current.get('version') != previous.get('version'):
                    impact['performance_impact'] = True
                    impact['performance_changes'].append(f"EKS version change: {previous.get('version')} → {current.get('version')}")
                
                if component == 'elasticache' and current.get('node_type') != previous.get('node_type'):
                    impact['performance_impact'] = True
                    impact['performance_changes'].append(f"Cache node type change: {previous.get('node_type')} → {current.get('node_type')}")
        
        if component in performance_critical:
            impact['gaming_specific'] = True
            
            for modified in changes['modified']:
                current = modified['current']
                previous = modified['previous']
                
                # Check for capacity changes that could impact gaming performance
                if component == 'auto_scaling_groups':
                    if current.get('max_size') != previous.get('max_size'):
                        impact['performance_impact'] = True
                        impact['performance_changes'].append(f"ASG max size change: {previous.get('max_size')} → {current.get('max_size')}")
                
                if component == 'rds_instances':
                    if current.get('instance_class') != previous.get('instance_class'):
                        impact['performance_impact'] = True
                        impact['performance_changes'].append(f"RDS instance class change: {previous.get('instance_class')} → {current.get('instance_class')}")
        
        return impact

    def _calculate_gaming_risk_score(self, drift_report: Dict[str, Any]) -> int:
        """Calculate gaming-specific risk score (0-10)"""
        score = 0
        
        # Critical changes = high risk
        score += len(drift_report['critical_changes']) * 3
        
        # Performance impact = medium risk
        score += len(drift_report['performance_impact']) * 2
        
        # Gaming-specific changes = low risk
        score += len(drift_report['gaming_specific_changes']) * 1
        
        # Cap at 10
        return min(score, 10)

    def send_gaming_alert(self, drift_report: Dict[str, Any]) -> None:
        """Send gaming-specific drift alert"""
        if not drift_report['drift_detected']:
            return
        
        # Prepare alert message
        alert_message = self._create_gaming_alert_message(drift_report)
        
        # Send SNS notification
        if self.sns_topic_arn:
            try:
                self.sns.publish(
                    TopicArn=self.sns_topic_arn,
                    Subject=f"🎮 Herald.lol Gaming Infrastructure Drift Alert - {drift_report['gaming_impact'].upper()}",
                    Message=alert_message
                )
                logger.info("📧 SNS gaming drift alert sent")
            except Exception as e:
                logger.error(f"❌ Failed to send SNS alert: {e}")
        
        # Send Slack notification
        if self.slack_webhook_url:
            try:
                self._send_slack_gaming_alert(drift_report)
                logger.info("💬 Slack gaming drift alert sent")
            except Exception as e:
                logger.error(f"❌ Failed to send Slack alert: {e}")

    def _create_gaming_alert_message(self, drift_report: Dict[str, Any]) -> str:
        """Create gaming-specific alert message"""
        message = f"""
🎮 HERALD.LOL GAMING INFRASTRUCTURE DRIFT DETECTED

Environment: {drift_report['gaming_environment']}
Impact Level: {drift_report['gaming_impact'].upper()}
Gaming Risk Score: {drift_report['gaming_risk_score']}/10
Detected: {drift_report['timestamp']}

🎯 GAMING PERFORMANCE IMPACT:
Target: <{self.performance_target_ms}ms analytics response
Concurrent Users: {self.concurrent_users:,} target capacity

"""
        
        if drift_report['critical_changes']:
            message += "🚨 CRITICAL GAMING CHANGES:\n"
            for change in drift_report['critical_changes']:
                message += f"- {change}\n"
            message += "\n"
        
        if drift_report['performance_impact']:
            message += "⚡ PERFORMANCE IMPACT:\n"
            for change in drift_report['performance_impact']:
                message += f"- {change}\n"
            message += "\n"
        
        if drift_report['gaming_specific_changes']:
            message += "🎮 GAMING-SPECIFIC CHANGES:\n"
            for change in drift_report['gaming_specific_changes']:
                message += f"- {change}\n"
            message += "\n"
        
        message += """
🛠️ RECOMMENDED ACTIONS:
1. Review infrastructure changes immediately
2. Test gaming analytics performance (<5s target)
3. Validate Riot API integration functionality
4. Check real-time gaming connections (WebSocket/gRPC)
5. Monitor gaming metrics for performance impact

🎮 Gaming Platform: Herald.lol
📊 Dashboard: CloudWatch Gaming Infrastructure Dashboard
"""
        
        return message

    def _send_slack_gaming_alert(self, drift_report: Dict[str, Any]) -> None:
        """Send Slack alert for gaming infrastructure drift"""
        
        # Determine color based on impact
        color_mapping = {
            'critical': '#ff0000',  # Red
            'performance': '#ffa500',  # Orange
            'moderate': '#ffff00',  # Yellow
            'none': '#00ff00'  # Green
        }
        
        color = color_mapping.get(drift_report['gaming_impact'], '#808080')
        
        slack_message = {
            "attachments": [
                {
                    "color": color,
                    "title": f"🎮 Herald.lol Gaming Infrastructure Drift Alert",
                    "title_link": f"https://console.aws.amazon.com/cloudwatch/home#dashboards:name=Herald-Gaming-Infrastructure-Drift",
                    "fields": [
                        {
                            "title": "Environment",
                            "value": drift_report['gaming_environment'],
                            "short": True
                        },
                        {
                            "title": "Impact Level",
                            "value": drift_report['gaming_impact'].upper(),
                            "short": True
                        },
                        {
                            "title": "Gaming Risk Score",
                            "value": f"{drift_report['gaming_risk_score']}/10",
                            "short": True
                        },
                        {
                            "title": "Performance Target",
                            "value": f"<{self.performance_target_ms}ms",
                            "short": True
                        }
                    ],
                    "footer": "Herald.lol Gaming Platform",
                    "ts": int(datetime.now(timezone.utc).timestamp())
                }
            ]
        }
        
        if drift_report['critical_changes']:
            slack_message["attachments"][0]["fields"].append({
                "title": "🚨 Critical Changes",
                "value": "\n".join(drift_report['critical_changes'][:3]),  # Limit to 3
                "short": False
            })
        
        response = requests.post(self.slack_webhook_url, json=slack_message)
        response.raise_for_status()


def handler(event, context):
    """AWS Lambda handler for gaming infrastructure drift detection"""
    logger.info("🎮 Starting Herald.lol gaming infrastructure drift detection...")
    
    try:
        detector = HeraldGamingDriftDetector()
        
        # Get current infrastructure state
        current_state = detector.get_gaming_infrastructure_state()
        
        # Save current state snapshot
        snapshot_key = detector.save_state_snapshot(current_state)
        
        # Get previous state for comparison
        previous_state = detector.get_previous_state()
        
        if previous_state:
            # Detect drift
            drift_report = detector.detect_gaming_drift(current_state, previous_state)
            
            # Send alerts if drift detected
            if drift_report['drift_detected']:
                detector.send_gaming_alert(drift_report)
            
            return {
                'statusCode': 200,
                'body': json.dumps({
                    'message': 'Gaming infrastructure drift detection completed',
                    'drift_detected': drift_report['drift_detected'],
                    'gaming_impact': drift_report['gaming_impact'],
                    'gaming_risk_score': drift_report['gaming_risk_score'],
                    'snapshot_key': snapshot_key
                })
            }
        else:
            logger.info("📸 Initial gaming infrastructure snapshot created")
            return {
                'statusCode': 200,
                'body': json.dumps({
                    'message': 'Initial gaming infrastructure snapshot created',
                    'snapshot_key': snapshot_key
                })
            }
            
    except Exception as e:
        logger.error(f"❌ Gaming infrastructure drift detection failed: {e}")
        return {
            'statusCode': 500,
            'body': json.dumps({
                'error': 'Gaming infrastructure drift detection failed',
                'details': str(e)
            })
        }


if __name__ == "__main__":
    # For local testing
    test_event = {}
    test_context = {}
    
    result = handler(test_event, test_context)
    print(json.dumps(result, indent=2))