# ⚡ Performance Optimization Guide

This guide helps you optimize the LoL Match Exporter for better performance and efficiency.

## 🎯 Quick Wins

### 1. Enable Caching
```typescript
// In ExporterMUI component - keep cache enabled
const [useCache, setUseCache] = useState(true); // ✅ Default to true
```

**Impact**: Reduces API calls by 40-60% for repeated data requests.

### 2. Use Seasonal Filtering
```typescript
// Filter by specific season instead of all-time data
const [season, setSeason] = useState(2024); // ✅ Current season only
```

**Impact**: Reduces data volume and processing time by 70-80%.

### 3. Optimize Export Count
```typescript
// For analysis, 100-500 matches is usually sufficient
const COUNT_MAX = 500; // ✅ Instead of 1000
```

**Impact**: Faster exports, less memory usage, better user experience.

## 🔧 API Optimization

### Rate Limiting Configuration

Adjust based on your API key tier:

```python
# For Personal API Key (20 req/sec, 100 req/2min)
RATE_LIMITER = AdaptiveRateLimiter(
    short_limit=18,      # ✅ Leave some headroom
    short_window=1.0,
    long_limit=90,       # ✅ Leave some headroom
    long_window=120.0,
    burst_allowance=0.1  # ✅ 10% burst for high priority
)

# For Production API Key (Higher limits)
RATE_LIMITER = AdaptiveRateLimiter(
    short_limit=95,      # Adjust based on your limits
    short_window=1.0,
    long_limit=3000,     # Adjust based on your limits
    long_window=120.0,
    burst_allowance=0.2
)
```

### Cache Optimization

```python
# Optimize cache settings for your use case
cache = LRUCache(
    max_size=5000 if production else 1000  # ✅ More cache in production
)

# Cache TTL by data type
CACHE_TTL = {
    "match_data": 300,      # 5 minutes (matches don't change)
    "summoner_data": 3600,  # 1 hour (summoner info stable)
    "champion_data": 86400, # 24 hours (champion info very stable)
    "queue_data": 86400,    # 24 hours (queue info very stable)
}
```

### Batch Processing

```python
# Process matches in batches to reduce memory usage
async def process_matches_batched(match_ids: List[str], batch_size: int = 50):
    results = []
    for i in range(0, len(match_ids), batch_size):
        batch = match_ids[i:i + batch_size]
        batch_results = await process_match_batch(batch)
        results.extend(batch_results)
        
        # Small delay between batches to be API-friendly
        await asyncio.sleep(0.1)
    
    return results
```

## 🖥️ Frontend Optimization

### Component Memoization

```typescript
// Memoize expensive calculations
const championStats = useMemo(() => {
    return data.reduce((acc, row) => {
        // ... expensive calculation
        return acc;
    }, {});
}, [data]); // ✅ Only recalculate when data changes

// Memoize components that don't change often
const MemoizedChampionIcon = React.memo(ChampionIcon);
const MemoizedChartContainer = React.memo(ChartContainer);
```

### Chart Optimization

```typescript
// Reduce data points for better chart performance
const chartData = useMemo(() => {
    const data = processedData;
    
    // For large datasets, sample or aggregate
    if (data.length > 100) {
        return data.slice(0, 100); // ✅ Show top 100
    }
    
    return data;
}, [processedData]);

// Use appropriate chart types
// ✅ BarChart for categorical data (fast)
// ✅ LineChart for time series (fast)
// ⚠️ ScatterChart for correlations (slower with many points)
// ⚠️ RadarChart for comparisons (slower)
```

### Virtual Scrolling for Large Tables

```typescript
// DataGrid already uses virtual scrolling
<DataGrid
    rows={championStats}
    columns={columns}
    density="compact"          // ✅ Compact density
    rowsPerPageOptions={[25, 50, 100]} // ✅ Reasonable page sizes
    initialState={{
        pagination: {
            paginationModel: { pageSize: 25 } // ✅ Start with smaller pages
        }
    }}
/>
```

## 🗄️ Data Management

### Efficient Data Structures

```typescript
// Use Maps for O(1) lookups instead of arrays
const championMap = useMemo(() => {
    const map = new Map();
    data.forEach(row => {
        const champion = row.champion;
        if (!map.has(champion)) {
            map.set(champion, { games: 0, wins: 0, /* ... */ });
        }
        // Update stats
    });
    return map;
}, [data]);

// Use Sets for unique collections
const uniqueChampions = useMemo(() => {
    return new Set(data.map(row => row.champion));
}, [data]);
```

### Lazy Loading

```typescript
// Lazy load heavy components
const ChampionDetails = React.lazy(() => import('./views/ChampionDetails'));
const AdvancedCharts = React.lazy(() => import('./components/AdvancedCharts'));

// Use Suspense for loading states
<Suspense fallback={<CircularProgress />}>
    <ChampionDetails {...props} />
</Suspense>
```

## 📊 Memory Management

### Monitor Memory Usage

```typescript
// Add memory monitoring (development only)
if (process.env.NODE_ENV === 'development') {
    const observer = new PerformanceObserver((list) => {
        list.getEntries().forEach((entry) => {
            console.log(`${entry.name}: ${entry.duration}ms`);
        });
    });
    observer.observe({ entryTypes: ['measure', 'navigation'] });
}
```

### Cleanup

```typescript
// Cleanup effects properly
useEffect(() => {
    const eventSource = new EventSource(url);
    
    return () => {
        eventSource.close(); // ✅ Always cleanup
    };
}, [url]);

// Clear intervals and timeouts
useEffect(() => {
    const interval = setInterval(() => {
        // ... polling logic
    }, 5000);
    
    return () => clearInterval(interval); // ✅ Cleanup
}, []);
```

## 🚀 Production Optimizations

### Bundle Optimization

```typescript
// vite.config.ts
export default defineConfig({
    build: {
        rollupOptions: {
            output: {
                manualChunks: {
                    'vendor': ['react', 'react-dom'],
                    'mui': ['@mui/material', '@mui/icons-material'],
                    'charts': ['recharts'],
                    'export': ['html2canvas', 'xlsx']
                }
            }
        },
        chunkSizeWarningLimit: 1000 // Increase limit for chunks
    }
});
```

### Service Worker for Caching

```typescript
// public/sw.js
const CACHE_NAME = 'lol-exporter-v2.1.0';
const STATIC_CACHE = [
    '/',
    '/static/js/bundle.js',
    '/static/css/main.css'
];

self.addEventListener('install', (event) => {
    event.waitUntil(
        caches.open(CACHE_NAME)
            .then(cache => cache.addAll(STATIC_CACHE))
    );
});
```

## 📈 Performance Monitoring

### Frontend Metrics

```typescript
// Track performance metrics
const trackPerformance = (name: string, fn: () => void) => {
    const start = performance.now();
    fn();
    const duration = performance.now() - start;
    
    if (duration > 100) { // Log slow operations
        console.warn(`Slow operation: ${name} took ${duration}ms`);
    }
};

// Usage
trackPerformance('processChampionData', () => {
    setChampionStats(processChampionData(data));
});
```

### Backend Metrics

```python
# Monitor API performance
import time
from functools import wraps

def monitor_performance(func):
    @wraps(func)
    async def wrapper(*args, **kwargs):
        start = time.time()
        result = await func(*args, **kwargs)
        duration = time.time() - start
        
        if duration > 5.0:  # Log slow API calls
            print(f"Slow API call: {func.__name__} took {duration:.2f}s")
        
        return result
    return wrapper

@monitor_performance
async def get_match_data(match_id: str):
    # ... API call logic
```

## 🎛️ Configuration Presets

### Development Mode
```python
# Fast iteration, more logging
RATE_LIMITER = AdaptiveRateLimiter(short_limit=10, long_limit=50)
CACHE_SIZE = 500
LOG_LEVEL = "DEBUG"
```

### Production Mode
```python
# Optimized for throughput
RATE_LIMITER = AdaptiveRateLimiter(short_limit=18, long_limit=90)
CACHE_SIZE = 5000
LOG_LEVEL = "INFO"
```

### Bulk Export Mode
```python
# For large-scale data exports
RATE_LIMITER = AdaptiveRateLimiter(
    short_limit=15,    # More conservative
    long_limit=80,
    burst_allowance=0.05  # Less burst
)
CACHE_SIZE = 10000  # Larger cache
BATCH_SIZE = 25     # Smaller batches
```

## 🔍 Troubleshooting Performance Issues

### High Memory Usage
1. Check cache size: `cache.get_stats()`
2. Monitor component re-renders with React DevTools
3. Use `useMemo` for expensive calculations
4. Implement pagination for large datasets

### Slow API Responses
1. Check rate limiter stats
2. Monitor cache hit rates
3. Use seasonal filtering
4. Implement request prioritization

### Slow Chart Rendering
1. Reduce data points (sampling)
2. Use appropriate chart types
3. Implement virtual scrolling
4. Debounce user interactions

With these optimizations, your LoL Match Exporter will run smoothly even with large datasets! ⚡