import { useState, useEffect, useCallback, useMemo } from 'react';

// Performance monitoring hook
export const usePerformance = () => {
  const [metrics, setMetrics] = useState({
    loadTime: 0,
    renderTime: 0,
    memoryUsage: 0,
    isSlowDevice: false,
  });

  useEffect(() => {
    // Measure initial load time
    if (performance.timing) {
      const loadTime = performance.timing.loadEventEnd - performance.timing.navigationStart;
      setMetrics(prev => ({ ...prev, loadTime }));
    }

    // Detect slow device based on hardware concurrency
    const isSlowDevice = navigator.hardwareConcurrency <= 2;
    setMetrics(prev => ({ ...prev, isSlowDevice }));

    // Memory usage monitoring (if available)
    if ('memory' in performance) {
      const memoryInfo = (performance as any).memory;
      const memoryUsage = memoryInfo.usedJSHeapSize / memoryInfo.jsHeapSizeLimit;
      setMetrics(prev => ({ ...prev, memoryUsage }));
    }
  }, []);

  // Debounced performance measurement
  const measureRenderTime = useCallback((componentName: string) => {
    const startTime = performance.now();
    
    return () => {
      const endTime = performance.now();
      const renderTime = endTime - startTime;
      
      // Log slow renders (> 16ms for 60fps)
      if (renderTime > 16) {
        console.warn(`Slow render detected: ${componentName} took ${renderTime.toFixed(2)}ms`);
      }
      
      setMetrics(prev => ({ ...prev, renderTime }));
    };
  }, []);

  // Performance optimization helpers
  const shouldReduceAnimations = useMemo(() => {
    return metrics.isSlowDevice || metrics.memoryUsage > 0.8;
  }, [metrics.isSlowDevice, metrics.memoryUsage]);

  const shouldLazyLoad = useMemo(() => {
    return metrics.isSlowDevice || window.navigator.connection?.effectiveType === '2g';
  }, [metrics.isSlowDevice]);

  const getOptimalChunkSize = useCallback((totalItems: number) => {
    if (metrics.isSlowDevice) return Math.min(10, totalItems);
    if (totalItems > 100) return 25;
    if (totalItems > 50) return 15;
    return totalItems;
  }, [metrics.isSlowDevice]);

  return {
    metrics,
    measureRenderTime,
    shouldReduceAnimations,
    shouldLazyLoad,
    getOptimalChunkSize,
  };
};

// Virtual scrolling hook for large lists
export const useVirtualScroll = (items: any[], itemHeight: number, containerHeight: number) => {
  const [scrollTop, setScrollTop] = useState(0);
  
  const visibleItems = useMemo(() => {
    const startIndex = Math.floor(scrollTop / itemHeight);
    const endIndex = Math.min(
      startIndex + Math.ceil(containerHeight / itemHeight) + 1,
      items.length
    );
    
    return {
      startIndex,
      endIndex,
      items: items.slice(startIndex, endIndex),
      totalHeight: items.length * itemHeight,
      offsetY: startIndex * itemHeight,
    };
  }, [items, itemHeight, containerHeight, scrollTop]);

  const handleScroll = useCallback((e: React.UIEvent) => {
    setScrollTop(e.currentTarget.scrollTop);
  }, []);

  return {
    visibleItems,
    handleScroll,
  };
};

// Image lazy loading hook
export const useLazyImage = (src: string, placeholder?: string) => {
  const [imageSrc, setImageSrc] = useState(placeholder || '');
  const [imageRef, setImageRef] = useState<HTMLImageElement | null>(null);
  const [isLoaded, setIsLoaded] = useState(false);

  useEffect(() => {
    let observer: IntersectionObserver;
    
    if (imageRef && src) {
      observer = new IntersectionObserver(
        ([entry]) => {
          if (entry.isIntersecting) {
            setImageSrc(src);
            observer.unobserve(imageRef);
          }
        },
        { threshold: 0.1 }
      );
      
      observer.observe(imageRef);
    }
    
    return () => {
      if (observer && imageRef) {
        observer.unobserve(imageRef);
      }
    };
  }, [imageRef, src]);

  const handleLoad = useCallback(() => {
    setIsLoaded(true);
  }, []);

  return {
    src: imageSrc,
    ref: setImageRef,
    isLoaded,
    onLoad: handleLoad,
  };
};

export default usePerformance;