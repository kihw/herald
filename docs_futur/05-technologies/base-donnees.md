# Architecture Base de Données et Stockage

## Vue d'Ensemble de l'Architecture Data

Herald.lol implémente une **architecture de données polyglotte** qui optimise chaque type de données selon ses caractéristiques et patterns d'usage. Cette approche garantit performance, scalabilité et cohérence à travers l'ensemble de la plateforme.

## Architecture Polyglotte Multi-Base

### Bases de Données Transactionnelles

#### PostgreSQL Cluster (Primary)
- **Version** : PostgreSQL 15+ avec extensions avancées
- **Configuration** : Master-Slave avec réplication synchrone
- **Extensions** : PostGIS pour données géographiques, pg_cron pour jobs
- **Usage** : Données transactionnelles, profils utilisateur, configurations

##### Schema Design Optimisé
```sql
-- Core User Management
users (id, riot_puuid, riot_id, riot_tag, region, preferences)
user_sessions (session_id, user_id, created_at, expires_at)
user_settings (user_id, platform, queue_types, sync_preferences)

-- Gaming Data
matches (match_id, platform, game_creation, game_duration, raw_data)
match_participants (match_id, user_id, champion_id, stats, performance_metrics)
match_timeline (match_id, timestamp, event_type, event_data)

-- Analytics Computed
user_analytics (user_id, period, metrics, computed_at)
champion_performance (user_id, champion_id, stats_aggregated)
ranking_history (user_id, tier, rank, lp, timestamp)
```

##### Performance Optimizations
- **Partitioning Strategy** : Partitionnement par temps pour matches
- **Index Strategy** : Index composites pour requêtes fréquentes
- **Query Optimization** : Requêtes optimisées avec EXPLAIN ANALYZE
- **Connection Pooling** : PgBouncer pour gestion connexions

#### SQLite pour Développement
- **Version** : SQLite 3.42+ avec WAL mode
- **Configuration** : Mode WAL pour concurrence
- **Extensions** : JSON extensions pour données semi-structurées
- **Usage** : Développement local, déploiements edge, prototypage

### Cache et Session Store

#### Redis Cluster Configuration
- **Version** : Redis 7+ avec modules Redis Stack
- **Topology** : Cluster 6 nœuds (3 masters, 3 replicas)
- **Persistence** : RDB + AOF pour durabilité
- **Modules** : RedisJSON, RedisTimeSeries, RedisBloom

##### Usage Patterns Optimisés
```redis
# Session Management
SET user:session:{session_id} {user_data} EX 604800  # 7 days

# Caching Strategy
SET cache:user:{user_id}:stats {stats_json} EX 3600  # 1 hour
SET cache:match:{match_id} {match_data} EX 86400     # 24 hours

# Real-Time Features
ZADD live:matches {timestamp} {match_id}
PUBLISH user:{user_id}:notifications {notification_json}

# Rate Limiting
INCR rate_limit:{api_key}:{endpoint} EX 60
```

##### Advanced Redis Features
- **Redis Streams** : Event streaming pour analytics temps réel
- **Pub/Sub** : Notifications temps réel utilisateurs
- **Lua Scripts** : Scripts atomiques pour opérations complexes
- **RedisJSON** : Stockage JSON performant pour données semi-structurées

### Time-Series et Metrics

#### InfluxDB pour Métriques Performance
- **Version** : InfluxDB 2.x avec Flux query language
- **Configuration** : Cluster haute disponibilité
- **Retention Policies** : Politiques rétention granulaires
- **Usage** : Métriques performance, analytics temps réel

##### Schema Time-Series
```flux
// Performance Metrics
measurement: "player_performance"
tags: user_id, champion_id, queue_type, region
fields: kda, cs_per_min, gold_per_min, damage_share
time: match_timestamp

// System Metrics
measurement: "api_metrics"
tags: endpoint, method, status_code
fields: response_time, throughput, error_rate
time: request_timestamp

// Analytics Events
measurement: "user_events"
tags: user_id, event_type, platform
fields: duration, success_rate, engagement_score
time: event_timestamp
```

##### InfluxDB Optimizations
- **Downsampling** : Agrégation automatique données anciennes
- **Continuous Queries** : Requêtes continues pour pré-agrégations
- **Kapacitor Integration** : Alerting basé sur métriques
- **Telegraf Collection** : Collection métriques système automatisée

### Document Store et Semi-Structured Data

#### MongoDB pour Données Flexibles
- **Version** : MongoDB 6+ avec sharding
- **Configuration** : Replica Set avec sharding horizontal
- **Indexes** : Index composites et text search
- **Usage** : Configurations dynamiques, templates, logs structurés

##### Collections Design
```javascript
// Dynamic Configurations
game_configs: {
  _id: ObjectId,
  game: "league_of_legends",
  version: "13.24",
  champions: [...],
  items: [...],
  meta_data: {...}
}

// User Templates
analytics_templates: {
  _id: ObjectId,
  user_id: "uuid",
  template_name: "My Custom Dashboard",
  widgets: [...],
  layout: {...},
  shared: boolean
}

// Event Logs
audit_logs: {
  _id: ObjectId,
  timestamp: Date,
  user_id: "uuid",
  action: "data_export",
  details: {...},
  ip_address: "string"
}
```

##### MongoDB Features
- **Aggregation Pipeline** : Agrégations complexes pour analytics
- **Change Streams** : Real-time notifications changements
- **GridFS** : Stockage fichiers volumineux
- **Atlas Search** : Recherche full-text avancée

## Data Warehouse et Analytics

### Data Lake Architecture

#### Amazon S3 Data Lake
- **Structure** : Structure médaillée (Bronze/Silver/Gold)
- **Formats** : Parquet pour analytics, JSON pour raw data
- **Partitioning** : Partitionnement par date et région
- **Lifecycle** : Politiques lifecycle automatiques

##### Data Lake Structure
```
s3://herald-data-lake/
├── bronze/          # Raw data ingestion
│   ├── riot_api/
│   ├── user_events/
│   └── system_logs/
├── silver/          # Cleaned and validated
│   ├── matches/
│   ├── users/
│   └── analytics/
└── gold/           # Business-ready aggregated
    ├── dashboards/
    ├── reports/
    └── ml_features/
```

##### Data Processing Pipeline
- **Apache Airflow** : Orchestration pipelines ETL
- **AWS Glue** : ETL serverless pour transformations
- **Apache Spark** : Processing big data distribué
- **dbt** : Transformations SQL avec version control

#### Snowflake Data Warehouse (Roadmap)
- **Architecture** : Compute et storage séparés
- **Auto-Scaling** : Scaling automatique selon charge
- **Time Travel** : Historique modifications données
- **Data Sharing** : Partage sécurisé avec partenaires

### OLAP et Business Intelligence

#### Apache Druid pour Analytics Temps Réel
- **Version** : Druid 25+ avec native batch ingestion
- **Configuration** : Cluster distribué multi-tenant
- **Indexing** : Index bitmap pour queries rapides
- **Usage** : Dashboards temps réel, analytics interactives

##### Druid Datasources
```json
{
  "dataSource": "gaming_analytics",
  "dimensions": [
    "user_id", "champion_id", "queue_type", 
    "region", "rank_tier", "game_mode"
  ],
  "metrics": [
    "win_rate", "kda_avg", "cs_per_min", 
    "damage_share", "vision_score"
  ],
  "granularity": "HOUR"
}
```

##### Real-Time Analytics Features
- **Real-Time Ingestion** : Ingestion Kafka temps réel
- **Historical + Real-Time** : Lambda architecture hybride
- **Fast Aggregations** : Agrégations sub-seconde
- **SQL Interface** : Interface SQL standard pour BI tools

## Data Migration et Evolution

### Schema Evolution Strategy

#### Database Migrations Framework
- **Flyway Migrations** : Migrations versionnées pour PostgreSQL
- **Backward Compatibility** : Compatibilité versions antérieures
- **Blue-Green Deployments** : Déploiements sans downtime
- **Testing Migrations** : Tests automatisés migrations

##### Migration Best Practices
```sql
-- V001__Initial_schema.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    riot_puuid VARCHAR(78) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- V002__Add_user_preferences.sql
ALTER TABLE users ADD COLUMN preferences JSONB DEFAULT '{}';
CREATE INDEX idx_users_preferences ON users USING GIN (preferences);

-- V003__Partition_matches_table.sql
CREATE TABLE matches_partitioned (
    LIKE matches INCLUDING ALL
) PARTITION BY RANGE (game_creation);
```

#### Data Format Evolution
- **Schema Registry** : Registre schémas Avro/JSON Schema
- **Backward/Forward Compatibility** : Compatibilité bidirectionnelle
- **Version Management** : Gestion versions données
- **Migration Scripts** : Scripts migration automatisés

### Backup et Disaster Recovery

#### Multi-Level Backup Strategy
- **PostgreSQL** : Backup continu avec point-in-time recovery
- **Redis** : Snapshots RDB + journaux AOF
- **MongoDB** : Replica sets avec backup automatisé
- **S3** : Versioning et cross-region replication

##### Recovery Procedures
```bash
# PostgreSQL Point-in-Time Recovery
pg_basebackup -h master -D /recovery -U postgres -P -W
pg_ctl start -D /recovery -o "-c recovery_target_time='2025-01-01 12:00:00'"

# Redis Backup and Restore
redis-cli --rdb /backup/dump.rdb
redis-cli --pipe < /backup/commands.txt

# MongoDB Replica Set Recovery
mongodump --host replica-set/host1:27017,host2:27017,host3:27017
mongorestore --drop /backup/dump/
```

##### Disaster Recovery Metrics
- **RTO (Recovery Time Objective)** : < 4 heures
- **RPO (Recovery Point Objective)** : < 15 minutes
- **Cross-Region Failover** : Automatique avec Route 53
- **Data Consistency Checks** : Vérifications automatisées

## Performance et Optimization

### Query Optimization

#### PostgreSQL Performance Tuning
```sql
-- Connection et Memory Settings
shared_buffers = '2GB'
effective_cache_size = '6GB'
work_mem = '256MB'
maintenance_work_mem = '512MB'

-- Query Planning
random_page_cost = 1.1
seq_page_cost = 1.0
cpu_tuple_cost = 0.01
cpu_index_tuple_cost = 0.005

-- Logging pour Optimization
log_min_duration_statement = 1000
log_statement = 'all'
log_checkpoints = on
```

#### Index Strategy Optimization
```sql
-- Composite Indexes pour Queries Fréquentes
CREATE INDEX idx_matches_user_time ON matches (user_id, game_creation DESC);
CREATE INDEX idx_participants_champion ON match_participants (champion_id, win);

-- Partial Indexes pour Données Filtrées
CREATE INDEX idx_users_active ON users (last_active) WHERE is_active = true;

-- Expression Indexes pour Queries Complexes
CREATE INDEX idx_users_region_rank ON users ((preferences->>'region'), rank);
```

### Caching Strategy Multi-Niveau

#### L1 Cache (Application)
- **In-Memory Caching** : Cache applicatif avec TTL
- **LRU Eviction** : Éviction LRU pour mémoire limitée
- **Cache Warming** : Préchauffage cache critique
- **Hit Rate Monitoring** : Monitoring taux hit cache

#### L2 Cache (Redis)
- **Distributed Caching** : Cache distribué cross-instances
- **Cache Aside Pattern** : Pattern cache-aside optimisé
- **Pub/Sub Invalidation** : Invalidation pub/sub coordonnée
- **Compression** : Compression données cache volumineuses

#### L3 Cache (CDN)
- **Edge Caching** : Cache edge global pour assets
- **Dynamic Content Caching** : Cache contenu dynamique
- **Cache Headers Optimization** : Headers cache optimisés
- **Purge API** : API purge cache programmable

### Monitoring et Observability

#### Database Monitoring
```yaml
# Prometheus Metrics pour PostgreSQL
pg_up: Connection status
pg_stat_database_tup_returned_total: Rows returned
pg_stat_database_tup_fetched_total: Rows fetched
pg_stat_database_tup_inserted_total: Rows inserted
pg_locks_count: Lock count by mode

# Redis Monitoring
redis_connected_clients: Connected clients
redis_used_memory_bytes: Memory usage
redis_keyspace_hits_total: Cache hits
redis_keyspace_misses_total: Cache misses
```

#### Performance Dashboards
- **Query Performance** : Dashboards performance requêtes
- **Resource Utilization** : Utilisation ressources temps réel
- **Cache Hit Rates** : Taux hit cache multi-niveau
- **Replication Lag** : Lag réplication cross-region

#### Alerting Rules
```yaml
# PostgreSQL Alerts
- alert: PostgreSQLDown
  expr: pg_up == 0
  for: 5m

- alert: PostgreSQLHighConnections
  expr: pg_stat_database_numbackends > 80

# Redis Alerts  
- alert: RedisHighMemoryUsage
  expr: redis_memory_used_bytes / redis_memory_max_bytes > 0.9

- alert: RedisLowCacheHitRate
  expr: rate(redis_keyspace_hits_total[5m]) / rate(redis_keyspace_misses_total[5m]) < 0.8
```

Cette architecture de données robuste et évolutive permet à Herald.lol de gérer efficacement des téraoctets de données gaming tout en maintenant des performances exceptionnelles pour des millions d'utilisateurs.