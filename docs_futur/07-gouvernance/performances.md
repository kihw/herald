# Performance et Optimisation

## Vue d'Ensemble des Performances

Herald.lol implémente une **stratégie de performance globale** qui optimise chaque composant de la plateforme pour offrir une expérience utilisateur exceptionnelle. Cette approche holistique couvre l'infrastructure, l'application, les données et l'expérience utilisateur.

## Objectifs de Performance

### Service Level Objectives (SLOs)

#### Métriques de Performance Critiques
```yaml
performance_targets:
  availability:
    target: 99.9%
    measurement_window: 30d
    downtime_budget: 43.2_minutes_per_month
    
  latency:
    api_endpoints:
      p50: < 100ms
      p95: < 500ms
      p99: < 1000ms
    page_load:
      p50: < 2s
      p95: < 5s
      p99: < 10s
      
  throughput:
    api_requests: 10000_req_per_sec
    concurrent_users: 50000
    data_processing: 1M_matches_per_hour
    
  error_rates:
    api_errors: < 0.1%
    frontend_errors: < 0.05%
    data_sync_failures: < 0.5%
```

#### Core Web Vitals Targets
- **Largest Contentful Paint (LCP)** : < 2.5s
- **First Input Delay (FID)** : < 100ms
- **Cumulative Layout Shift (CLS)** : < 0.1
- **First Contentful Paint (FCP)** : < 1.8s
- **Time to Interactive (TTI)** : < 3.8s

## Optimisation Infrastructure

### Architecture Haute Performance

#### Load Balancing et Distribution
```yaml
# HAProxy Configuration
global:
  maxconn 50000
  log stdout local0
  
defaults:
  mode http
  timeout connect 5000ms
  timeout client 50000ms
  timeout server 50000ms
  option httplog
  
frontend herald_frontend:
  bind *:80
  bind *:443 ssl crt /etc/ssl/herald.pem
  redirect scheme https if !{ ssl_fc }
  
  # Rate limiting
  stick-table type ip size 100k expire 30s store http_req_rate(10s)
  http-request track-sc0 src
  http-request reject if { sc_http_req_rate(0) gt 100 }
  
  # Routing
  use_backend api_servers if { path_beg /api/ }
  default_backend web_servers
  
backend api_servers:
  balance roundrobin
  option httpchk GET /health
  server api1 10.0.1.10:8000 check weight 100
  server api2 10.0.1.11:8000 check weight 100
  server api3 10.0.1.12:8000 check weight 100
  
backend web_servers:
  balance roundrobin
  option httpchk GET /health
  server web1 10.0.2.10:80 check weight 100
  server web2 10.0.2.11:80 check weight 100
```

#### Auto-Scaling Intelligent
```yaml
# Kubernetes HPA avec métriques custom
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: herald-backend-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: herald-backend
  minReplicas: 3
  maxReplicas: 50
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  - type: Pods
    pods:
      metric:
        name: api_request_rate
      target:
        type: AverageValue
        averageValue: "1000"
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 100
        periodSeconds: 15
      - type: Pods
        value: 5
        periodSeconds: 60
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60
```

### CDN et Edge Computing

#### Global Content Distribution
```javascript
// Cloudflare Workers pour edge computing
export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    const cache = caches.default;
    
    // Cache strategy par type de contenu
    if (url.pathname.startsWith('/api/')) {
      return handleAPIRequest(request, env);
    }
    
    if (url.pathname.match(/\.(js|css|png|jpg|jpeg|gif|ico|svg|woff2?)$/)) {
      return handleStaticAsset(request, cache);
    }
    
    return handleDynamicContent(request, env, cache);
  }
};

async function handleStaticAsset(request, cache) {
  const cacheKey = new Request(request.url, request);
  let response = await cache.match(cacheKey);
  
  if (!response) {
    response = await fetch(request);
    
    if (response.ok) {
      // Cache static assets for 1 year
      const newResponse = new Response(response.body, {
        status: response.status,
        statusText: response.statusText,
        headers: {
          ...response.headers,
          'Cache-Control': 'public, max-age=31536000, immutable',
          'CDN-Cache-Control': 'max-age=31536000'
        }
      });
      
      ctx.waitUntil(cache.put(cacheKey, newResponse.clone()));
      return newResponse;
    }
  }
  
  return response;
}

async function handleAPIRequest(request, env) {
  // Géolocalisation intelligente
  const country = request.cf.country;
  const region = getOptimalRegion(country);
  
  // Load balancing régional
  const backend = selectBackend(region, request.url);
  
  return fetch(backend, {
    method: request.method,
    headers: request.headers,
    body: request.body
  });
}
```

## Optimisation Application Backend

### Performance Go Server

#### Connection Pooling Optimisé
```go
// Database connection pool optimization
type DatabaseConfig struct {
    MaxOpenConns    int           `yaml:"max_open_conns"`
    MaxIdleConns    int           `yaml:"max_idle_conns"`
    ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
    ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
}

func NewOptimizedDB(config DatabaseConfig) (*sql.DB, error) {
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        return nil, err
    }
    
    // Optimisation connection pool
    db.SetMaxOpenConns(config.MaxOpenConns)       // 100 for high traffic
    db.SetMaxIdleConns(config.MaxIdleConns)       // 25 idle connections
    db.SetConnMaxLifetime(config.ConnMaxLifetime) // 1 hour
    db.SetConnMaxIdleTime(config.ConnMaxIdleTime) // 15 minutes
    
    return db, nil
}

// Query optimization avec prepared statements
type QueryCache struct {
    statements map[string]*sql.Stmt
    mutex      sync.RWMutex
}

func (qc *QueryCache) GetPreparedStatement(db *sql.DB, query string) (*sql.Stmt, error) {
    qc.mutex.RLock()
    stmt, exists := qc.statements[query]
    qc.mutex.RUnlock()
    
    if exists {
        return stmt, nil
    }
    
    qc.mutex.Lock()
    defer qc.mutex.Unlock()
    
    // Double-check après acquisition du lock
    if stmt, exists := qc.statements[query]; exists {
        return stmt, nil
    }
    
    stmt, err := db.Prepare(query)
    if err != nil {
        return nil, err
    }
    
    qc.statements[query] = stmt
    return stmt, nil
}
```

#### Cache Strategy Multi-Niveau
```go
// L1 Cache: In-memory avec TTL
type MemoryCache struct {
    data   map[string]*CacheEntry
    mutex  sync.RWMutex
    maxTTL time.Duration
}

type CacheEntry struct {
    Value     interface{}
    ExpiresAt time.Time
}

func (mc *MemoryCache) Get(key string) (interface{}, bool) {
    mc.mutex.RLock()
    defer mc.mutex.RUnlock()
    
    entry, exists := mc.data[key]
    if !exists || time.Now().After(entry.ExpiresAt) {
        return nil, false
    }
    
    return entry.Value, true
}

func (mc *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
    mc.mutex.Lock()
    defer mc.mutex.Unlock()
    
    mc.data[key] = &CacheEntry{
        Value:     value,
        ExpiresAt: time.Now().Add(ttl),
    }
}

// L2 Cache: Redis distributed
type DistributedCache struct {
    client      *redis.ClusterClient
    localCache  *MemoryCache
    compression bool
}

func (dc *DistributedCache) Get(ctx context.Context, key string) ([]byte, error) {
    // Try L1 cache first
    if value, found := dc.localCache.Get(key); found {
        return value.([]byte), nil
    }
    
    // Try L2 cache (Redis)
    data, err := dc.client.Get(ctx, key).Bytes()
    if err != nil {
        return nil, err
    }
    
    // Decompress if needed
    if dc.compression {
        data, err = decompress(data)
        if err != nil {
            return nil, err
        }
    }
    
    // Store in L1 cache
    dc.localCache.Set(key, data, 5*time.Minute)
    
    return data, nil
}

func (dc *DistributedCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    // Compress if enabled
    if dc.compression {
        compressed, err := compress(value)
        if err == nil && len(compressed) < len(value) {
            value = compressed
        }
    }
    
    // Store in both caches
    dc.localCache.Set(key, value, min(ttl, 5*time.Minute))
    return dc.client.Set(ctx, key, value, ttl).Err()
}
```

#### Async Processing Optimisé
```go
// Worker pool pour traitement asynchrone
type WorkerPool struct {
    jobs        chan Job
    workers     []*Worker
    wg          sync.WaitGroup
    ctx         context.Context
    cancel      context.CancelFunc
    maxWorkers  int
    maxJobs     int
}

type Job interface {
    Execute(ctx context.Context) error
    Priority() int
    Retry() bool
}

func NewWorkerPool(maxWorkers, maxJobs int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    
    wp := &WorkerPool{
        jobs:       make(chan Job, maxJobs),
        workers:    make([]*Worker, maxWorkers),
        ctx:        ctx,
        cancel:     cancel,
        maxWorkers: maxWorkers,
        maxJobs:    maxJobs,
    }
    
    // Start workers
    for i := 0; i < maxWorkers; i++ {
        worker := &Worker{
            id:   i,
            jobs: wp.jobs,
            ctx:  ctx,
        }
        wp.workers[i] = worker
        wp.wg.Add(1)
        go worker.Start(&wp.wg)
    }
    
    return wp
}

func (wp *WorkerPool) Submit(job Job) error {
    select {
    case wp.jobs <- job:
        return nil
    default:
        return errors.New("job queue full")
    }
}

type Worker struct {
    id   int
    jobs <-chan Job
    ctx  context.Context
}

func (w *Worker) Start(wg *sync.WaitGroup) {
    defer wg.Done()
    
    for {
        select {
        case job := <-w.jobs:
            if err := w.executeWithRetry(job); err != nil {
                log.Printf("Worker %d failed to execute job: %v", w.id, err)
            }
        case <-w.ctx.Done():
            return
        }
    }
}

func (w *Worker) executeWithRetry(job Job) error {
    const maxRetries = 3
    var lastErr error
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        if err := job.Execute(w.ctx); err != nil {
            lastErr = err
            if !job.Retry() {
                break
            }
            
            // Exponential backoff
            backoff := time.Duration(attempt*attempt) * time.Second
            time.Sleep(backoff)
            continue
        }
        
        return nil
    }
    
    return lastErr
}
```

## Optimisation Frontend

### React Performance Optimization

#### Code Splitting et Lazy Loading
```typescript
// Route-based code splitting
import { lazy, Suspense } from 'react';
import { Routes, Route } from 'react-router-dom';
import LoadingSpinner from './components/LoadingSpinner';

// Lazy load components
const Dashboard = lazy(() => import('./pages/Dashboard'));
const Analytics = lazy(() => import('./pages/Analytics'));
const Settings = lazy(() => import('./pages/Settings'));

function App() {
  return (
    <Suspense fallback={<LoadingSpinner />}>
      <Routes>
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/analytics" element={<Analytics />} />
        <Route path="/settings" element={<Settings />} />
      </Routes>
    </Suspense>
  );
}

// Component-level code splitting
const HeavyAnalyticsChart = lazy(() => 
  import('./components/HeavyAnalyticsChart').then(module => ({
    default: module.HeavyAnalyticsChart
  }))
);

// Progressive loading avec intersection observer
const LazyComponent: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [isVisible, setIsVisible] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  
  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true);
          observer.disconnect();
        }
      },
      { threshold: 0.1 }
    );
    
    if (ref.current) {
      observer.observe(ref.current);
    }
    
    return () => observer.disconnect();
  }, []);
  
  return (
    <div ref={ref}>
      {isVisible ? children : <div className="placeholder" />}
    </div>
  );
};
```

#### Optimisation State Management
```typescript
// React Query optimizations
import { useQuery, useInfiniteQuery, useMutation, useQueryClient } from '@tanstack/react-query';

// Optimized query configuration
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutes
      cacheTime: 10 * 60 * 1000, // 10 minutes
      retry: (failureCount, error) => {
        if (error.status === 404) return false;
        return failureCount < 3;
      },
      retryDelay: attemptIndex => Math.min(1000 * 2 ** attemptIndex, 30000),
    },
  },
});

// Optimized data fetching avec parallel queries
function usePlayerDashboard(playerId: string) {
  const playerQuery = useQuery({
    queryKey: ['player', playerId],
    queryFn: () => fetchPlayer(playerId),
    staleTime: 30 * 60 * 1000, // 30 minutes for player data
  });
  
  const matchesQuery = useInfiniteQuery({
    queryKey: ['player-matches', playerId],
    queryFn: ({ pageParam = 0 }) => fetchPlayerMatches(playerId, pageParam),
    getNextPageParam: (lastPage) => lastPage.nextCursor,
    enabled: !!playerQuery.data,
  });
  
  const analyticsQuery = useQuery({
    queryKey: ['player-analytics', playerId],
    queryFn: () => fetchPlayerAnalytics(playerId),
    enabled: !!playerQuery.data,
    staleTime: 15 * 60 * 1000, // 15 minutes for analytics
  });
  
  return {
    player: playerQuery.data,
    matches: matchesQuery.data?.pages.flatMap(page => page.matches) ?? [],
    analytics: analyticsQuery.data,
    isLoading: playerQuery.isLoading || matchesQuery.isLoading || analyticsQuery.isLoading,
    hasMore: matchesQuery.hasNextPage,
    loadMore: matchesQuery.fetchNextPage,
  };
}

// Memoization sophistiquée
const ExpensiveComponent = memo(({ data, filters }: Props) => {
  const processedData = useMemo(() => {
    return processComplexData(data, filters);
  }, [data, filters]);
  
  const expensiveCalculation = useMemo(() => {
    return calculateComplexMetrics(processedData);
  }, [processedData]);
  
  return <div>{/* Render processed data */}</div>;
}, (prevProps, nextProps) => {
  // Custom comparison function
  return (
    prevProps.data === nextProps.data &&
    JSON.stringify(prevProps.filters) === JSON.stringify(nextProps.filters)
  );
});
```

#### Virtual Scrolling pour Large Lists
```typescript
// Virtualized list component pour performance
import { FixedSizeList as List } from 'react-window';

interface VirtualizedMatchListProps {
  matches: Match[];
  height: number;
  itemHeight: number;
}

const VirtualizedMatchList: React.FC<VirtualizedMatchListProps> = ({
  matches,
  height,
  itemHeight
}) => {
  const Row = useCallback(({ index, style }: { index: number; style: React.CSSProperties }) => {
    const match = matches[index];
    
    return (
      <div style={style}>
        <MatchRow match={match} />
      </div>
    );
  }, [matches]);
  
  return (
    <List
      height={height}
      itemCount={matches.length}
      itemSize={itemHeight}
      width="100%"
    >
      {Row}
    </List>
  );
};

// Intersection observer pour infinite scroll
const useInfiniteScroll = (callback: () => void, hasMore: boolean) => {
  const observer = useRef<IntersectionObserver>();
  
  const lastElementRef = useCallback((node: HTMLDivElement) => {
    if (observer.current) observer.current.disconnect();
    
    observer.current = new IntersectionObserver(entries => {
      if (entries[0].isIntersecting && hasMore) {
        callback();
      }
    });
    
    if (node) observer.current.observe(node);
  }, [callback, hasMore]);
  
  return lastElementRef;
};
```

## Optimisation Base de Données

### PostgreSQL Performance Tuning

#### Index Strategy Avancée
```sql
-- Composite indexes pour queries fréquentes
CREATE INDEX CONCURRENTLY idx_matches_user_time_performance 
ON matches (user_id, game_creation DESC, game_duration) 
WHERE game_creation > NOW() - INTERVAL '6 months';

-- Partial indexes pour données filtrées
CREATE INDEX CONCURRENTLY idx_matches_ranked_recent 
ON matches (user_id, rank_tier, game_creation DESC) 
WHERE queue_type = 'RANKED_SOLO_5x5' 
AND game_creation > NOW() - INTERVAL '3 months';

-- Expression indexes pour recherches complexes
CREATE INDEX CONCURRENTLY idx_champions_winrate 
ON match_participants ((stats->>'win')::boolean, champion_id, 
                        (stats->>'kills')::int + (stats->>'assists')::int);

-- GIN indexes pour JSON data
CREATE INDEX CONCURRENTLY idx_match_participants_stats_gin 
ON match_participants USING GIN (stats jsonb_path_ops);

-- Covering indexes pour éviter les lookups
CREATE INDEX CONCURRENTLY idx_users_comprehensive 
ON users (riot_puuid) 
INCLUDE (riot_id, riot_tag, region, created_at, last_sync);
```

#### Query Optimization
```sql
-- Optimized aggregate queries avec window functions
WITH player_stats AS (
  SELECT 
    user_id,
    champion_id,
    COUNT(*) as games_played,
    AVG((stats->>'kills')::int) as avg_kills,
    AVG((stats->>'deaths')::int) as avg_deaths,
    AVG((stats->>'assists')::int) as avg_assists,
    SUM(CASE WHEN (stats->>'win')::boolean THEN 1 ELSE 0 END) as wins,
    ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY COUNT(*) DESC) as champion_rank
  FROM match_participants mp
  JOIN matches m ON mp.match_id = m.match_id
  WHERE m.game_creation > NOW() - INTERVAL '30 days'
    AND m.queue_type = 'RANKED_SOLO_5x5'
  GROUP BY user_id, champion_id
  HAVING COUNT(*) >= 5
),
champion_performance AS (
  SELECT 
    *,
    ROUND((wins::numeric / games_played) * 100, 2) as win_rate,
    ROUND((avg_kills + avg_assists) / NULLIF(avg_deaths, 0), 2) as kda_ratio
  FROM player_stats
  WHERE champion_rank <= 5  -- Top 5 champions par utilisateur
)
SELECT * FROM champion_performance
ORDER BY user_id, win_rate DESC;

-- Materialized views pour analytics lourdes
CREATE MATERIALIZED VIEW mv_daily_player_stats AS
SELECT 
  user_id,
  DATE(game_creation) as game_date,
  COUNT(*) as games_played,
  AVG((stats->>'kills')::int + (stats->>'assists')::int) as avg_kda,
  AVG(CAST(stats->>'goldEarned' AS int)) as avg_gold,
  SUM(CASE WHEN (stats->>'win')::boolean THEN 1 ELSE 0 END) as wins
FROM match_participants mp
JOIN matches m ON mp.match_id = m.match_id
WHERE m.game_creation >= CURRENT_DATE - INTERVAL '90 days'
GROUP BY user_id, DATE(game_creation);

-- Refresh automatique avec pg_cron
SELECT cron.schedule('refresh-daily-stats', '0 2 * * *', 
  'REFRESH MATERIALIZED VIEW CONCURRENTLY mv_daily_player_stats;');
```

#### Partitioning Strategy
```sql
-- Time-based partitioning pour matches table
CREATE TABLE matches_partitioned (
  LIKE matches INCLUDING ALL
) PARTITION BY RANGE (game_creation);

-- Partitions mensuelles
CREATE TABLE matches_2024_01 PARTITION OF matches_partitioned
  FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

CREATE TABLE matches_2024_02 PARTITION OF matches_partitioned
  FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');

-- Automatisation création partitions
CREATE OR REPLACE FUNCTION create_monthly_partition(table_name text, start_date date)
RETURNS void AS $$
DECLARE
  partition_name text;
  end_date date;
BEGIN
  partition_name := table_name || '_' || to_char(start_date, 'YYYY_MM');
  end_date := start_date + interval '1 month';
  
  EXECUTE format('CREATE TABLE %I PARTITION OF %I FOR VALUES FROM (%L) TO (%L)',
    partition_name, table_name, start_date, end_date);
END;
$$ LANGUAGE plpgsql;
```

## Monitoring et Métriques Performance

### Real-Time Performance Monitoring

#### Custom Metrics Collection
```go
// Performance metrics collector
type PerformanceCollector struct {
    registry prometheus.Registerer
    
    // Request metrics
    requestDuration *prometheus.HistogramVec
    requestCount    *prometheus.CounterVec
    
    // Database metrics
    dbConnections   *prometheus.GaugeVec
    dbQueryDuration *prometheus.HistogramVec
    
    // Cache metrics
    cacheHitRate    *prometheus.GaugeVec
    cacheOperations *prometheus.CounterVec
    
    // Business metrics
    activeUsers     *prometheus.GaugeVec
    dataProcessed   *prometheus.CounterVec
}

func NewPerformanceCollector() *PerformanceCollector {
    pc := &PerformanceCollector{
        requestDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "herald_request_duration_seconds",
                Help: "Request duration in seconds",
                Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
            },
            []string{"method", "endpoint", "status_code"},
        ),
        
        dbQueryDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "herald_db_query_duration_seconds",
                Help: "Database query duration in seconds",
                Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
            },
            []string{"query_type", "table"},
        ),
        
        cacheHitRate: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "herald_cache_hit_rate",
                Help: "Cache hit rate percentage",
            },
            []string{"cache_type", "cache_level"},
        ),
    }
    
    // Register metrics
    prometheus.MustRegister(
        pc.requestDuration,
        pc.requestCount, 
        pc.dbConnections,
        pc.dbQueryDuration,
        pc.cacheHitRate,
        pc.cacheOperations,
        pc.activeUsers,
        pc.dataProcessed,
    )
    
    return pc
}
```

#### Performance Dashboards
```yaml
# Grafana dashboard configuration pour performance
dashboard:
  title: "Herald.lol - Performance Overview"
  panels:
    - title: "Request Latency Distribution"
      type: "heatmap"
      targets:
        - expr: "rate(herald_request_duration_seconds_bucket[5m])"
          legendFormat: "{{le}}"
    
    - title: "Throughput"
      type: "graph"
      targets:
        - expr: "rate(herald_request_count[5m])"
          legendFormat: "{{method}} {{endpoint}}"
          
    - title: "Error Rate"
      type: "stat"
      targets:
        - expr: "rate(herald_request_count{status_code=~'5..'}[5m]) / rate(herald_request_count[5m]) * 100"
          legendFormat: "Error Rate %"
          
    - title: "Database Performance"
      type: "graph"
      targets:
        - expr: "rate(herald_db_query_duration_seconds_sum[5m]) / rate(herald_db_query_duration_seconds_count[5m])"
          legendFormat: "Avg Query Time"
        - expr: "herald_db_connections{state='active'}"
          legendFormat: "Active Connections"
          
    - title: "Cache Efficiency"
      type: "graph"
      targets:
        - expr: "herald_cache_hit_rate"
          legendFormat: "{{cache_type}} Hit Rate"
```

Cette stratégie de performance complète garantit que Herald.lol maintient des performances exceptionnelles à tous les niveaux, offrant une expérience utilisateur fluide même sous forte charge.