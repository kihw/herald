import React, { useEffect, useRef } from 'react';

interface PerformanceMetrics {
  renderTime: number;
  memoryUsage: number;
  apiCalls: number;
  timestamp: number;
}

interface PerformanceMonitorProps {
  componentName: string;
  onMetricsUpdate?: (metrics: PerformanceMetrics) => void;
  enableLogging?: boolean;
}

/**
 * Performance monitoring component to track render times and resource usage
 * Only active in development mode for performance optimization insights
 */
const PerformanceMonitor: React.FC<PerformanceMonitorProps> = ({
  componentName,
  onMetricsUpdate,
  enableLogging = process.env.NODE_ENV === 'development'
}) => {
  const startTimeRef = useRef<number>(0);
  const apiCallCountRef = useRef<number>(0);
  const renderCountRef = useRef<number>(0);
  const lastLogTimeRef = useRef<number>(0);

  useEffect(() => {
    if (!enableLogging) return;

    startTimeRef.current = performance.now();
    renderCountRef.current += 1;

    // Monitor memory usage (if available)
    const getMemoryUsage = (): number => {
      if ('memory' in performance) {
        const memory = (performance as any).memory;
        return memory.usedJSHeapSize / 1024 / 1024; // MB
      }
      return 0;
    };

    // Track performance metrics
    const recordMetrics = () => {
      const endTime = performance.now();
      const renderTime = endTime - startTimeRef.current;
      const memoryUsage = getMemoryUsage();
      const now = Date.now();
      
      const metrics: PerformanceMetrics = {
        renderTime,
        memoryUsage,
        apiCalls: apiCallCountRef.current,
        timestamp: now
      };

      // Throttle logging to max once per 2 seconds to prevent spam
      const timeSinceLastLog = now - lastLogTimeRef.current;
      const shouldLog = timeSinceLastLog > 2000;

      // Log slow renders (throttled)
      if (renderTime > 16 && shouldLog) {
        console.warn(
          `🐌 Slow render detected in ${componentName}: ${renderTime.toFixed(2)}ms`
        );
        lastLogTimeRef.current = now;
      }

      // Log high memory usage (throttled)
      if (memoryUsage > 100 && shouldLog) {
        console.warn(
          `🧠 High memory usage in ${componentName}: ${memoryUsage.toFixed(2)}MB`
        );
        lastLogTimeRef.current = now;
      }

      // Call callback if provided (but don't call too frequently)
      if (shouldLog && onMetricsUpdate) {
        onMetricsUpdate(metrics);
      }
    };

    // Record metrics after render
    const timeoutId = setTimeout(recordMetrics, 0);

    return () => {
      clearTimeout(timeoutId);
    };
  }, [componentName, enableLogging]); // Removed onMetricsUpdate to prevent infinite loop

  // Intercept fetch calls to count API usage
  useEffect(() => {
    if (!enableLogging) return;

    const originalFetch = window.fetch;
    
    window.fetch = async (...args) => {
      apiCallCountRef.current += 1;
      console.log(`🔌 API call #${apiCallCountRef.current} from ${componentName}: ${args[0]}`);
      
      const startTime = performance.now();
      const response = await originalFetch(...args);
      const endTime = performance.now();
      
      console.log(
        `📡 API response (${response.status}) in ${(endTime - startTime).toFixed(2)}ms`
      );
      
      return response;
    };

    return () => {
      window.fetch = originalFetch;
    };
  }, [componentName, enableLogging]);

  // Performance observer for navigation timing
  useEffect(() => {
    if (!enableLogging) return;

    const observer = new PerformanceObserver((list) => {
      list.getEntries().forEach((entry) => {
        if (entry.entryType === 'navigation') {
          const navEntry = entry as PerformanceNavigationTiming;
          console.log(`🚀 Page load metrics for ${componentName}:`, {
            domContentLoaded: navEntry.domContentLoadedEventEnd - navEntry.domContentLoadedEventStart,
            loadComplete: navEntry.loadEventEnd - navEntry.loadEventStart,
            firstPaint: navEntry.responseEnd - navEntry.fetchStart
          });
        }
      });
    });

    observer.observe({ entryTypes: ['navigation', 'measure'] });

    return () => {
      observer.disconnect();
    };
  }, [componentName, enableLogging]);

  // This component doesn't render anything
  return null;
};

// Utility hook for tracking component performance
export const usePerformanceTracker = (componentName: string) => {
  const startTime = useRef<number>(0);
  const renderCount = useRef<number>(0);

  useEffect(() => {
    startTime.current = performance.now();
    renderCount.current += 1;
    
    return () => {
      if (process.env.NODE_ENV === 'development') {
        const endTime = performance.now();
        const renderTime = endTime - startTime.current;
        
        if (renderTime > 16) {
          console.warn(
            `⏱️ ${componentName} render #${renderCount.current}: ${renderTime.toFixed(2)}ms`
          );
        }
      }
    };
  });

  return {
    trackOperation: (name: string, fn: () => void) => {
      if (process.env.NODE_ENV !== 'development') {
        fn();
        return;
      }

      const start = performance.now();
      fn();
      const duration = performance.now() - start;
      
      if (duration > 10) {
        console.warn(`🔍 Slow operation in ${componentName}: ${name} took ${duration.toFixed(2)}ms`);
      }
    }
  };
};

export default PerformanceMonitor;
