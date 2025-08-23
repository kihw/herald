// Herald.lol Gaming Analytics Dashboard Visual Regression Tests
import { test, expect, Page } from '@playwright/test';

// Herald.lol gaming analytics visual testing
test.describe('Gaming Analytics Dashboard Visual Regression', () => {
  let page: Page;
  
  test.beforeEach(async ({ page: testPage }) => {
    page = testPage;
    
    // Mock gaming analytics API responses
    await page.route('**/api/analytics/**', async (route) => {
      const url = route.request().url();
      
      if (url.includes('/kda')) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            data: {
              currentKDA: 2.34,
              previousKDA: 2.1,
              trend: 'improving',
              percentile: 65,
              rankComparison: 'above_average',
              confidence: 0.92,
            }
          })
        });
      } else if (url.includes('/cs')) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            data: {
              currentCSPerMin: 7.2,
              previousCSPerMin: 6.8,
              trend: 'improving',
              efficiency: 85,
              roleAverage: 6.9,
              percentile: 72,
              confidence: 0.88,
            }
          })
        });
      } else if (url.includes('/vision')) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            data: {
              averageVisionScore: 28.4,
              wardPlacement: 'good',
              visionControl: 72,
              visionEfficiency: 78,
              percentile: 58,
              confidence: 0.79,
            }
          })
        });
      } else if (url.includes('/damage')) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            data: {
              damageShare: 0.32,
              damagePerMinute: 850,
              efficiency: 89,
              teamContribution: 'high',
              percentile: 71,
              confidence: 0.91,
            }
          })
        });
      } else if (url.includes('/gold')) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            data: {
              goldPerMinute: 425,
              goldEfficiency: 87,
              economicRating: 'good',
              percentile: 68,
              confidence: 0.84,
            }
          })
        });
      }
    });
    
    // Mock authentication
    await page.route('**/api/auth/**', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          user: {
            id: 1,
            email: 'visualtest@herald.lol',
            username: 'VisualTestUser',
            displayName: 'Visual Test User'
          },
          token: 'visual-test-token'
        })
      });
    });
    
    // Navigate to analytics dashboard
    await page.goto('/analytics');
    
    // Wait for gaming analytics to load (Herald.lol <5s requirement)
    await page.waitForSelector('[data-testid="analytics-dashboard"]', { timeout: 5000 });
    await page.waitForLoadState('networkidle');
  });
  
  test('Gaming Analytics Dashboard - Desktop View', async () => {
    // Set desktop viewport for gaming
    await page.setViewportSize({ width: 1920, height: 1080 });
    
    // Wait for all gaming widgets to load
    await expect(page.locator('[data-testid="kda-widget"]')).toBeVisible();
    await expect(page.locator('[data-testid="cs-widget"]')).toBeVisible();
    await expect(page.locator('[data-testid="vision-widget"]')).toBeVisible();
    await expect(page.locator('[data-testid="damage-widget"]')).toBeVisible();
    await expect(page.locator('[data-testid="gold-widget"]')).toBeVisible();
    
    // Take full dashboard screenshot
    await expect(page).toHaveScreenshot('gaming-analytics-dashboard-desktop.png', {
      fullPage: true,
      animations: 'disabled',
    });
  });
  
  test('Gaming Analytics Dashboard - Tablet View', async () => {
    // Set tablet viewport for gaming
    await page.setViewportSize({ width: 1024, height: 768 });
    
    // Wait for responsive layout
    await page.waitForTimeout(1000);
    
    // Verify gaming widgets are visible in tablet layout
    await expect(page.locator('[data-testid="analytics-dashboard"]')).toBeVisible();
    
    // Take tablet dashboard screenshot
    await expect(page).toHaveScreenshot('gaming-analytics-dashboard-tablet.png', {
      fullPage: true,
      animations: 'disabled',
    });
  });
  
  test('Gaming Analytics Dashboard - Mobile View', async () => {
    // Set mobile viewport for gaming
    await page.setViewportSize({ width: 390, height: 844 });
    
    // Wait for mobile responsive layout
    await page.waitForTimeout(1000);
    
    // Verify mobile gaming layout
    await expect(page.locator('[data-testid="analytics-dashboard"]')).toBeVisible();
    
    // Take mobile dashboard screenshot
    await expect(page).toHaveScreenshot('gaming-analytics-dashboard-mobile.png', {
      fullPage: true,
      animations: 'disabled',
    });
  });
  
  test('KDA Widget Visual States', async () => {
    const kdaWidget = page.locator('[data-testid="kda-widget"]');
    
    // Default KDA widget state
    await expect(kdaWidget).toHaveScreenshot('kda-widget-default.png');
    
    // Hover state for KDA widget
    await kdaWidget.hover();
    await expect(kdaWidget).toHaveScreenshot('kda-widget-hover.png');
    
    // Loading state (if applicable)
    await page.route('**/api/analytics/kda**', async (route) => {
      // Delay response to show loading state
      await new Promise(resolve => setTimeout(resolve, 2000));
      await route.continue();
    });
    
    await page.reload();
    await expect(page.locator('[data-testid="loading-spinner"]')).toBeVisible();
    await expect(kdaWidget).toHaveScreenshot('kda-widget-loading.png');
  });
  
  test('CS/min Widget Visual States', async () => {
    const csWidget = page.locator('[data-testid="cs-widget"]');
    
    // Default CS widget state
    await expect(csWidget).toHaveScreenshot('cs-widget-default.png');
    
    // Click to expand detailed view
    await csWidget.click();
    await page.waitForTimeout(500); // Wait for animation
    await expect(csWidget).toHaveScreenshot('cs-widget-expanded.png');
  });
  
  test('Vision Score Widget Visual States', async () => {
    const visionWidget = page.locator('[data-testid="vision-widget"]');
    
    // Default vision widget state
    await expect(visionWidget).toHaveScreenshot('vision-widget-default.png');
    
    // Heatmap visualization (if available)
    const heatmapToggle = visionWidget.locator('[data-testid="heatmap-toggle"]');
    if (await heatmapToggle.isVisible()) {
      await heatmapToggle.click();
      await page.waitForTimeout(1000); // Wait for heatmap rendering
      await expect(visionWidget).toHaveScreenshot('vision-widget-heatmap.png');
    }
  });
  
  test('Damage Analysis Widget Visual States', async () => {
    const damageWidget = page.locator('[data-testid="damage-widget"]');
    
    // Default damage widget state  
    await expect(damageWidget).toHaveScreenshot('damage-widget-default.png');
    
    // Damage breakdown chart
    const chartToggle = damageWidget.locator('[data-testid="chart-toggle"]');
    if (await chartToggle.isVisible()) {
      await chartToggle.click();
      await page.waitForTimeout(1000); // Wait for chart rendering
      await expect(damageWidget).toHaveScreenshot('damage-widget-chart.png');
    }
  });
  
  test('Gold Efficiency Widget Visual States', async () => {
    const goldWidget = page.locator('[data-testid="gold-widget"]');
    
    // Default gold widget state
    await expect(goldWidget).toHaveScreenshot('gold-widget-default.png');
    
    // Timeline view (if available)
    const timelineToggle = goldWidget.locator('[data-testid="timeline-toggle"]');
    if (await timelineToggle.isVisible()) {
      await timelineToggle.click();
      await page.waitForTimeout(1000); // Wait for timeline rendering
      await expect(goldWidget).toHaveScreenshot('gold-widget-timeline.png');
    }
  });
  
  test('Gaming Analytics Dark Theme', async () => {
    // Toggle dark theme (if available)
    const themeToggle = page.locator('[data-testid="theme-toggle"]');
    if (await themeToggle.isVisible()) {
      await themeToggle.click();
      await page.waitForTimeout(500); // Wait for theme transition
      
      // Take dark theme screenshot
      await expect(page).toHaveScreenshot('gaming-analytics-dashboard-dark.png', {
        fullPage: true,
        animations: 'disabled',
      });
    }
  });
  
  test('Gaming Analytics Error States', async () => {
    // Mock API error for visual testing
    await page.route('**/api/analytics/kda**', async (route) => {
      await route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({
          error: 'Gaming analytics service unavailable'
        })
      });
    });
    
    await page.reload();
    await page.waitForTimeout(2000);
    
    // Take error state screenshot
    await expect(page).toHaveScreenshot('gaming-analytics-error-state.png', {
      fullPage: true,
    });
  });
  
  test('Gaming Analytics Performance Indicators', async () => {
    // Check performance indicators are visible
    await expect(page.locator('[data-testid="confidence-level"]')).toBeVisible();
    await expect(page.locator('[data-testid="trend-indicator"]')).toBeVisible();
    
    // Take screenshot with performance indicators highlighted
    await expect(page).toHaveScreenshot('gaming-analytics-performance-indicators.png', {
      fullPage: true,
      animations: 'disabled',
    });
  });
});