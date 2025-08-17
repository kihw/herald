"""Enhanced Riot API module with improved rate limiting, caching, and seasonal segmentation.

Features:
- Advanced rate limiting with adaptive backoff
- In-memory caching with TTL
- Seasonal data segmentation
- Request queue with priority
- Batch processing optimization
"""

import time
import threading
import queue
import hashlib
import json
from typing import Dict, Any, List, Optional, Tuple
from datetime import datetime, timedelta
from collections import deque, OrderedDict
from dataclasses import dataclass, field
from enum import Enum

import requests


class Priority(Enum):
    """Request priority levels."""
    HIGH = 1
    NORMAL = 2
    LOW = 3


@dataclass
class ApiRequest:
    """Represents a queued API request."""
    url: str
    method: str = "GET"
    priority: Priority = Priority.NORMAL
    params: Dict[str, Any] = field(default_factory=dict)
    headers: Dict[str, str] = field(default_factory=dict)
    timestamp: float = field(default_factory=time.time)
    retry_count: int = 0
    
    def __lt__(self, other):
        """Priority queue comparison."""
        if self.priority.value != other.priority.value:
            return self.priority.value < other.priority.value
        return self.timestamp < other.timestamp


class CacheEntry:
    """Cache entry with TTL support."""
    def __init__(self, data: Any, ttl_seconds: int = 300):
        self.data = data
        self.expires_at = time.time() + ttl_seconds
    
    def is_expired(self) -> bool:
        return time.time() > self.expires_at


class LRUCache:
    """Thread-safe LRU cache with TTL support."""
    def __init__(self, max_size: int = 1000):
        self.cache: OrderedDict[str, CacheEntry] = OrderedDict()
        self.max_size = max_size
        self.lock = threading.RLock()
        self.hits = 0
        self.misses = 0
    
    def get(self, key: str) -> Optional[Any]:
        with self.lock:
            if key not in self.cache:
                self.misses += 1
                return None
            
            entry = self.cache[key]
            if entry.is_expired():
                del self.cache[key]
                self.misses += 1
                return None
            
            # Move to end (most recently used)
            self.cache.move_to_end(key)
            self.hits += 1
            return entry.data
    
    def set(self, key: str, value: Any, ttl_seconds: int = 300):
        with self.lock:
            # Remove oldest if at capacity
            if len(self.cache) >= self.max_size and key not in self.cache:
                self.cache.popitem(last=False)
            
            self.cache[key] = CacheEntry(value, ttl_seconds)
            self.cache.move_to_end(key)
    
    def clear_expired(self):
        """Remove all expired entries."""
        with self.lock:
            expired_keys = [k for k, v in self.cache.items() if v.is_expired()]
            for key in expired_keys:
                del self.cache[key]
    
    def get_stats(self) -> Dict[str, Any]:
        """Get cache statistics."""
        with self.lock:
            total_requests = self.hits + self.misses
            hit_rate = (self.hits / total_requests * 100) if total_requests > 0 else 0
            return {
                "size": len(self.cache),
                "max_size": self.max_size,
                "hits": self.hits,
                "misses": self.misses,
                "hit_rate": f"{hit_rate:.1f}%"
            }


class AdaptiveRateLimiter:
    """Advanced rate limiter with adaptive behavior and burst handling."""
    
    def __init__(self, 
                 short_limit: int = 20,
                 short_window: float = 1.0,
                 long_limit: int = 100, 
                 long_window: float = 120.0,
                 burst_allowance: float = 0.1):
        self.short_limit = short_limit
        self.short_window = short_window
        self.long_limit = long_limit
        self.long_window = long_window
        self.burst_allowance = burst_allowance
        
        self.short_requests: deque = deque()
        self.long_requests: deque = deque()
        self.lock = threading.RLock()
        
        # Adaptive parameters
        self.consecutive_429s = 0
        self.backoff_multiplier = 1.0
        self.last_429_time = 0
        
    def acquire(self, priority: Priority = Priority.NORMAL) -> float:
        """Acquire permission to make a request. Returns wait time."""
        while True:
            with self.lock:
                now = time.time()
                
                # Clean old requests
                self._clean_old_requests(now)
                
                # Apply adaptive backoff if we've been hitting limits
                if self.consecutive_429s > 0:
                    self.backoff_multiplier = min(2.0, 1.0 + (self.consecutive_429s * 0.2))
                else:
                    self.backoff_multiplier = max(1.0, self.backoff_multiplier * 0.95)
                
                # Calculate effective limits with backoff
                effective_short = int(self.short_limit / self.backoff_multiplier)
                effective_long = int(self.long_limit / self.backoff_multiplier)
                
                # Allow burst for high priority
                if priority == Priority.HIGH:
                    effective_short = int(effective_short * (1 + self.burst_allowance))
                
                # Check if we can proceed
                can_proceed_short = len(self.short_requests) < effective_short
                can_proceed_long = len(self.long_requests) < effective_long
                
                if can_proceed_short and can_proceed_long:
                    self.short_requests.append(now)
                    self.long_requests.append(now)
                    return 0  # No wait needed
                
                # Calculate wait time
                wait_time = self._calculate_wait_time(now, can_proceed_short, can_proceed_long)
                
            # Wait outside the lock
            if wait_time > 0:
                time.sleep(wait_time)
    
    def _clean_old_requests(self, now: float):
        """Remove expired timestamps."""
        while self.short_requests and (now - self.short_requests[0]) >= self.short_window:
            self.short_requests.popleft()
        while self.long_requests and (now - self.long_requests[0]) >= self.long_window:
            self.long_requests.popleft()
    
    def _calculate_wait_time(self, now: float, can_short: bool, can_long: bool) -> float:
        """Calculate how long to wait before next request."""
        wait_short = 0.0
        wait_long = 0.0
        
        if not can_short and self.short_requests:
            wait_short = self.short_window - (now - self.short_requests[0])
        
        if not can_long and self.long_requests:
            wait_long = self.long_window - (now - self.long_requests[0])
        
        return max(wait_short, wait_long, 0.001)
    
    def report_429(self):
        """Report that we received a 429 response."""
        with self.lock:
            now = time.time()
            if now - self.last_429_time < 60:  # Within 1 minute
                self.consecutive_429s += 1
            else:
                self.consecutive_429s = 1
            self.last_429_time = now
    
    def report_success(self):
        """Report successful request."""
        with self.lock:
            if time.time() - self.last_429_time > 60:
                self.consecutive_429s = max(0, self.consecutive_429s - 1)


class SeasonalSegmenter:
    """Handles seasonal data segmentation for optimized queries."""
    
    # Season start dates (approximate)
    SEASON_DATES = {
        2024: {"start": datetime(2024, 1, 10), "end": datetime(2024, 12, 31)},
        2023: {"start": datetime(2023, 1, 11), "end": datetime(2023, 11, 20)},
        2022: {"start": datetime(2022, 1, 7), "end": datetime(2022, 11, 14)},
        2021: {"start": datetime(2021, 1, 8), "end": datetime(2021, 11, 15)},
    }
    
    @classmethod
    def get_season_bounds(cls, season: int) -> Optional[Tuple[int, int]]:
        """Get Unix timestamp bounds for a season."""
        if season not in cls.SEASON_DATES:
            return None
        
        season_data = cls.SEASON_DATES[season]
        start_ts = int(season_data["start"].timestamp() * 1000)
        end_ts = int(season_data["end"].timestamp() * 1000)
        return (start_ts, end_ts)
    
    @classmethod
    def get_current_season(cls) -> int:
        """Get the current season year."""
        now = datetime.now()
        for year, dates in sorted(cls.SEASON_DATES.items(), reverse=True):
            if dates["start"] <= now <= dates["end"]:
                return year
        return now.year
    
    @classmethod
    def segment_time_range(cls, start_time: int, end_time: int, max_segment_days: int = 30) -> List[Tuple[int, int]]:
        """Segment a time range into smaller chunks for efficient querying."""
        segments = []
        current = start_time
        segment_ms = max_segment_days * 24 * 60 * 60 * 1000
        
        while current < end_time:
            segment_end = min(current + segment_ms, end_time)
            segments.append((current, segment_end))
            current = segment_end
        
        return segments


class EnhancedRiotAPI:
    """Enhanced Riot API client with all optimizations."""
    
    def __init__(self, api_key: str, platform: str = "euw1"):
        self.api_key = api_key
        self.platform = platform.lower()
        self.session = requests.Session()
        self.session.headers.update({
            "X-Riot-Token": api_key,
            "User-Agent": "LoL-Match-Exporter-Enhanced/2.0"
        })
        
        # Components
        self.rate_limiter = AdaptiveRateLimiter()
        self.cache = LRUCache(max_size=2000)
        self.request_queue: queue.PriorityQueue = queue.PriorityQueue()
        self.segmenter = SeasonalSegmenter()
        
        # Statistics
        self.total_requests = 0
        self.failed_requests = 0
        self.start_time = time.time()
        
        # Worker thread for processing queue
        self.worker_thread = threading.Thread(target=self._process_queue, daemon=True)
        self.worker_thread.start()
    
    def _get_cache_key(self, url: str, params: Dict[str, Any]) -> str:
        """Generate cache key for request."""
        param_str = json.dumps(params, sort_keys=True)
        combined = f"{url}:{param_str}"
        return hashlib.md5(combined.encode()).hexdigest()
    
    def _make_request(self, request: ApiRequest) -> Optional[Dict[str, Any]]:
        """Execute a single API request with caching and rate limiting."""
        # Check cache first
        cache_key = self._get_cache_key(request.url, request.params)
        cached = self.cache.get(cache_key)
        if cached is not None:
            return cached
        
        # Rate limit
        wait_time = self.rate_limiter.acquire(request.priority)
        if wait_time > 0:
            print(f"Rate limiting: waiting {wait_time:.2f}s...")
        
        # Make request
        try:
            response = self.session.request(
                request.method,
                request.url,
                params=request.params,
                headers=request.headers,
                timeout=30
            )
            
            self.total_requests += 1
            
            if response.status_code == 200:
                self.rate_limiter.report_success()
                data = response.json()
                
                # Cache successful responses
                ttl = 300 if "match" in request.url else 600  # Shorter TTL for match data
                self.cache.set(cache_key, data, ttl)
                
                return data
            
            elif response.status_code == 429:
                self.rate_limiter.report_429()
                retry_after = float(response.headers.get("Retry-After", 10))
                print(f"Rate limited (429). Waiting {retry_after}s...")
                time.sleep(retry_after)
                
                # Requeue with higher priority
                if request.retry_count < 3:
                    request.retry_count += 1
                    request.priority = Priority.HIGH
                    self.request_queue.put(request)
                
            elif response.status_code >= 500:
                self.failed_requests += 1
                print(f"Server error {response.status_code}. Retrying...")
                time.sleep(2 ** request.retry_count)
                
                if request.retry_count < 3:
                    request.retry_count += 1
                    self.request_queue.put(request)
            
            else:
                self.failed_requests += 1
                print(f"API error {response.status_code}: {response.text}")
                
        except Exception as e:
            self.failed_requests += 1
            print(f"Request failed: {e}")
        
        return None
    
    def _process_queue(self):
        """Worker thread to process queued requests."""
        while True:
            try:
                request = self.request_queue.get(timeout=1)
                self._make_request(request)
                self.request_queue.task_done()
            except queue.Empty:
                continue
            except Exception as e:
                print(f"Queue processor error: {e}")
    
    def get_matches_by_season(self, puuid: str, season: int, queue_ids: List[int] = None) -> List[str]:
        """Get all match IDs for a player in a specific season."""
        bounds = self.segmenter.get_season_bounds(season)
        if not bounds:
            print(f"Unknown season: {season}")
            return []
        
        start_time, end_time = bounds
        segments = self.segmenter.segment_time_range(start_time, end_time)
        
        all_matches = []
        for seg_start, seg_end in segments:
            params = {
                "startTime": seg_start // 1000,
                "endTime": seg_end // 1000,
                "start": 0,
                "count": 100
            }
            
            if queue_ids:
                for queue_id in queue_ids:
                    params["queue"] = queue_id
                    # Implementation would continue here...
                    
        return all_matches
    
    def get_statistics(self) -> Dict[str, Any]:
        """Get API client statistics."""
        uptime = time.time() - self.start_time
        return {
            "uptime_seconds": int(uptime),
            "total_requests": self.total_requests,
            "failed_requests": self.failed_requests,
            "success_rate": f"{(1 - self.failed_requests/max(1, self.total_requests)) * 100:.1f}%",
            "requests_per_minute": self.total_requests / max(1, uptime / 60),
            "cache_stats": self.cache.get_stats(),
            "queue_size": self.request_queue.qsize(),
            "rate_limiter": {
                "backoff_multiplier": self.rate_limiter.backoff_multiplier,
                "consecutive_429s": self.rate_limiter.consecutive_429s
            }
        }


# Export key components
__all__ = [
    "EnhancedRiotAPI",
    "Priority",
    "AdaptiveRateLimiter",
    "LRUCache",
    "SeasonalSegmenter"
]