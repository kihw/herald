// Performance utilities for optimizing Herald.lol

// Debounce function for performance-critical operations
export const debounce = <T extends (...args: any[]) => any>(
  func: T,
  wait: number,
  immediate = false
): ((...args: Parameters<T>) => void) => {
  let timeout: NodeJS.Timeout | null = null;
  
  return (...args: Parameters<T>) => {
    const later = () => {
      timeout = null;
      if (!immediate) func(...args);
    };
    
    const callNow = immediate && !timeout;
    
    if (timeout) clearTimeout(timeout);
    timeout = setTimeout(later, wait);
    
    if (callNow) func(...args);
  };
};

// Throttle function for scroll and resize events
export const throttle = <T extends (...args: any[]) => any>(
  func: T,
  limit: number
): ((...args: Parameters<T>) => void) => {
  let inThrottle: boolean = false;
  
  return (...args: Parameters<T>) => {
    if (!inThrottle) {
      func(...args);
      inThrottle = true;
      setTimeout(() => (inThrottle = false), limit);
    }
  };
};

// Intersection Observer for lazy loading
export const createIntersectionObserver = (
  callback: IntersectionObserverCallback,
  options: IntersectionObserverInit = {}
): IntersectionObserver => {
  const defaultOptions: IntersectionObserverInit = {
    rootMargin: '50px',
    threshold: 0.1,
    ...options,
  };
  
  return new IntersectionObserver(callback, defaultOptions);
};

// Memory usage monitoring
export const getMemoryUsage = (): {
  used: number;
  total: number;
  percentage: number;
} | null => {
  if ('memory' in performance) {
    const memory = (performance as any).memory;
    return {
      used: memory.usedJSHeapSize,
      total: memory.totalJSHeapSize,
      percentage: (memory.usedJSHeapSize / memory.totalJSHeapSize) * 100,
    };
  }
  return null;
};

// FPS monitoring
export class FPSMonitor {
  private times: number[] = [];
  private fps = 0;
  
  public getFPS(): number {
    return this.fps;
  }
  
  public update(): void {
    const now = performance.now();
    while (this.times.length > 0 && this.times[0] <= now - 1000) {
      this.times.shift();
    }
    this.times.push(now);
    this.fps = this.times.length;
  }
}

// Performance budget checker
export const checkPerformanceBudget = () => {
  const budgets = {
    firstContentfulPaint: 1500, // 1.5s
    largestContentfulPaint: 2500, // 2.5s
    firstInputDelay: 100, // 100ms
    cumulativeLayoutShift: 0.1, // 0.1
  };
  
  return new Promise((resolve) => {
    if ('PerformanceObserver' in window) {
      const observer = new PerformanceObserver((list) => {
        const entries = list.getEntries();
        const metrics: Record<string, number> = {};
        
        entries.forEach((entry) => {
          if (entry.entryType === 'paint') {
            metrics[entry.name.replace('-', '')] = entry.startTime;
          } else if (entry.entryType === 'largest-contentful-paint') {
            metrics.largestContentfulPaint = entry.startTime;
          } else if (entry.entryType === 'first-input') {
            metrics.firstInputDelay = (entry as any).processingStart - entry.startTime;
          } else if (entry.entryType === 'layout-shift') {
            metrics.cumulativeLayoutShift = (metrics.cumulativeLayoutShift || 0) + (entry as any).value;
          }
        });
        
        const results = Object.entries(budgets).map(([metric, budget]) => ({
          metric,
          value: metrics[metric] || 0,
          budget,
          passing: (metrics[metric] || 0) <= budget,
        }));
        
        resolve(results);
      });
      
      observer.observe({ entryTypes: ['paint', 'largest-contentful-paint', 'first-input', 'layout-shift'] });
      
      // Timeout after 10 seconds
      setTimeout(() => {
        observer.disconnect();
        resolve([]);
      }, 10000);
    } else {
      resolve([]);
    }
  });
};

// Bundle size analyzer
export const analyzeBundleSize = async (): Promise<{
  totalSize: number;
  gzippedSize: number;
  chunks: Array<{ name: string; size: number }>;
}> => {
  try {
    // This would typically connect to a bundle analyzer API
    // For now, we'll return mock data
    return {
      totalSize: 2400000, // 2.4MB
      gzippedSize: 700000, // 700KB
      chunks: [
        { name: 'vendor', size: 800000 },
        { name: 'mui', size: 600000 },
        { name: 'charts', size: 400000 },
        { name: 'groups', size: 300000 },
        { name: 'auth', size: 150000 },
        { name: 'dashboard', size: 200000 },
        { name: 'analytics', size: 250000 },
      ],
    };
  } catch (error) {
    console.error('Failed to analyze bundle size:', error);
    throw error;
  }
};

// Image optimization utilities
export const optimizeImage = (
  src: string,
  options: {
    width?: number;
    height?: number;
    quality?: number;
    format?: 'webp' | 'avif' | 'jpeg' | 'png';
  } = {}
): string => {
  const { width, height, quality = 80, format = 'webp' } = options;
  
  // In a real implementation, this would connect to an image optimization service
  // For now, we'll return the original src with query parameters
  const params = new URLSearchParams();
  if (width) params.set('w', width.toString());
  if (height) params.set('h', height.toString());
  params.set('q', quality.toString());
  params.set('f', format);
  
  return `${src}?${params.toString()}`;
};

// Code splitting utilities
export const preloadRoute = (routeName: string): Promise<void> => {
  const routes: Record<string, () => Promise<any>> = {
    groups: () => import('../components/groups/GroupManagement'),
    charts: () => import('../components/charts/ComparisonCharts'),
    dashboard: () => import('../components/dashboard/MainDashboard'),
    auth: () => import('../components/auth/AuthPage'),
  };
  
  const importFn = routes[routeName];
  if (importFn) {
    return importFn().then(() => {
      console.log(`Route ${routeName} preloaded successfully`);
    });
  }
  
  return Promise.resolve();
};

// Service Worker registration
export const registerServiceWorker = async (): Promise<ServiceWorkerRegistration | null> => {
  if ('serviceWorker' in navigator && process.env.NODE_ENV === 'production') {
    try {
      const registration = await navigator.serviceWorker.register('/sw.js');
      console.log('Service Worker registered successfully:', registration);
      return registration;
    } catch (error) {
      console.error('Service Worker registration failed:', error);
      return null;
    }
  }
  return null;
};

// Critical resource hints
export const addResourceHints = (resources: Array<{ href: string; as: string; type?: string }>) => {
  resources.forEach(({ href, as, type }) => {
    const link = document.createElement('link');
    link.rel = 'preload';
    link.href = href;
    link.as = as;
    if (type) link.type = type;
    
    // Add to head if not already present
    if (!document.querySelector(`link[href="${href}"]`)) {
      document.head.appendChild(link);
    }
  });
};

// Export performance utilities
export default {
  debounce,
  throttle,
  createIntersectionObserver,
  getMemoryUsage,
  FPSMonitor,
  checkPerformanceBudget,
  analyzeBundleSize,
  optimizeImage,
  preloadRoute,
  registerServiceWorker,
  addResourceHints,
};