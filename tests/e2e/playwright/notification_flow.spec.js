/**
 * End-to-End Tests for Notification System
 * 
 * This file contains comprehensive E2E tests for the real-time notification
 * system using Playwright. Tests cover the complete user journey from
 * authentication to receiving real-time insights.
 */

const { test, expect } = require('@playwright/test');

// Test configuration
const BASE_URL = process.env.TEST_BASE_URL || 'http://localhost:5173';
const API_BASE_URL = process.env.TEST_API_URL || 'http://localhost:8001';

// Test data
const TEST_USER = {
  riotId: 'TestUser',
  riotTag: 'EUW',
  region: 'euw1'
};

test.describe('Notification System E2E Tests', () => {
  
  test.beforeEach(async ({ page }) => {
    // Setup: Navigate to the application
    await page.goto(BASE_URL);
    
    // Wait for the page to load
    await page.waitForLoadState('networkidle');
  });

  test.describe('Authentication and Setup', () => {
    
    test('should authenticate user and access notifications', async ({ page }) => {
      // Step 1: Authenticate with test credentials
      await page.fill('[data-testid="riot-id-input"]', TEST_USER.riotId);
      await page.fill('[data-testid="riot-tag-input"]', TEST_USER.riotTag);
      await page.selectOption('[data-testid="region-select"]', TEST_USER.region);
      
      await page.click('[data-testid="validate-account-button"]');
      
      // Step 2: Wait for authentication success
      await expect(page.locator('[data-testid="auth-success"]')).toBeVisible();
      
      // Step 3: Navigate to analytics dashboard
      await page.click('[data-testid="analytics-nav-link"]');
      await expect(page.locator('[data-testid="analytics-dashboard"]')).toBeVisible();
      
      // Step 4: Verify notification icon is present
      await expect(page.locator('[data-testid="notification-icon"]')).toBeVisible();
    });

    test('should handle authentication failure gracefully', async ({ page }) => {
      // Try to authenticate with invalid credentials
      await page.fill('[data-testid="riot-id-input"]', 'InvalidUser');
      await page.fill('[data-testid="riot-tag-input"]', 'XXX');
      await page.selectOption('[data-testid="region-select"]', 'euw1');
      
      await page.click('[data-testid="validate-account-button"]');
      
      // Should show error message
      await expect(page.locator('[data-testid="auth-error"]')).toBeVisible();
      
      // Notification system should not be accessible
      await expect(page.locator('[data-testid="notification-icon"]')).not.toBeVisible();
    });
  });

  test.describe('Insights Display', () => {
    
    test.beforeEach(async ({ page }) => {
      // Authenticate before each test
      await authenticateUser(page);
    });

    test('should display user insights', async ({ page }) => {
      // Navigate to insights page
      await page.click('[data-testid="notification-icon"]');
      await expect(page.locator('[data-testid="insights-panel"]')).toBeVisible();
      
      // Should show insights list
      const insightsList = page.locator('[data-testid="insights-list"]');
      await expect(insightsList).toBeVisible();
      
      // Should have insight items
      const insights = page.locator('[data-testid="insight-item"]');
      await expect(insights).toHaveCountGreaterThan(0);
      
      // Verify insight structure
      const firstInsight = insights.first();
      await expect(firstInsight.locator('[data-testid="insight-title"]')).toBeVisible();
      await expect(firstInsight.locator('[data-testid="insight-message"]')).toBeVisible();
      await expect(firstInsight.locator('[data-testid="insight-timestamp"]')).toBeVisible();
    });

    test('should filter unread insights', async ({ page }) => {
      await page.click('[data-testid="notification-icon"]');
      
      // Toggle unread filter
      await page.click('[data-testid="unread-filter-toggle"]');
      
      // Should only show unread insights
      const insights = page.locator('[data-testid="insight-item"]');
      const count = await insights.count();
      
      for (let i = 0; i < count; i++) {
        const insight = insights.nth(i);
        await expect(insight.locator('[data-testid="unread-indicator"]')).toBeVisible();
      }
    });

    test('should mark insights as read', async ({ page }) => {
      await page.click('[data-testid="notification-icon"]');
      
      // Get first unread insight
      const unreadInsight = page.locator('[data-testid="insight-item"]:has([data-testid="unread-indicator"])').first();
      
      if (await unreadInsight.count() > 0) {
        // Click on the insight
        await unreadInsight.click();
        
        // Should mark as read
        await expect(unreadInsight.locator('[data-testid="unread-indicator"]')).not.toBeVisible();
        
        // Notification badge should update
        const notificationBadge = page.locator('[data-testid="notification-badge"]');
        const initialCount = await notificationBadge.textContent();
        
        // Mark another insight as read
        await page.click('[data-testid="mark-all-read-button"]');
        
        // Badge count should decrease or disappear
        if (parseInt(initialCount) > 1) {
          await expect(notificationBadge).toHaveText((parseInt(initialCount) - 1).toString());
        } else {
          await expect(notificationBadge).not.toBeVisible();
        }
      }
    });

    test('should navigate to analytics page from insight action', async ({ page }) => {
      await page.click('[data-testid="notification-icon"]');
      
      // Click on insight with action URL
      const insightWithAction = page.locator('[data-testid="insight-item"]:has([data-testid="insight-action"])').first();
      
      if (await insightWithAction.count() > 0) {
        await insightWithAction.locator('[data-testid="insight-action"]').click();
        
        // Should navigate to analytics page
        await expect(page.locator('[data-testid="analytics-page"]')).toBeVisible();
        
        // URL should reflect the navigation
        expect(page.url()).toContain('/analytics');
      }
    });
  });

  test.describe('Real-time Notifications', () => {
    
    test.beforeEach(async ({ page }) => {
      await authenticateUser(page);
    });

    test('should receive real-time insights via SSE', async ({ page }) => {
      // Navigate to analytics dashboard
      await page.click('[data-testid="analytics-nav-link"]');
      
      // Setup SSE listener
      let sseReceived = false;
      page.on('response', response => {
        if (response.url().includes('/api/notifications/stream')) {
          sseReceived = true;
        }
      });
      
      // Enable real-time notifications
      await page.click('[data-testid="enable-realtime-notifications"]');
      
      // Should establish SSE connection
      await page.waitForTimeout(2000);
      expect(sseReceived).toBe(true);
      
      // Trigger a test insight
      await page.click('[data-testid="trigger-test-insight"]');
      
      // Should receive notification
      await expect(page.locator('[data-testid="new-insight-toast"]')).toBeVisible({ timeout: 10000 });
      
      // Notification badge should update
      await expect(page.locator('[data-testid="notification-badge"]')).toBeVisible();
    });

    test('should handle SSE connection errors gracefully', async ({ page }) => {
      // Mock network failure
      await page.route('**/api/notifications/stream', route => {
        route.abort();
      });
      
      await page.click('[data-testid="analytics-nav-link"]');
      await page.click('[data-testid="enable-realtime-notifications"]');
      
      // Should show connection error
      await expect(page.locator('[data-testid="connection-error"]')).toBeVisible();
      
      // Should provide retry option
      await expect(page.locator('[data-testid="retry-connection"]')).toBeVisible();
    });

    test('should reconnect after connection loss', async ({ page }) => {
      await page.click('[data-testid="analytics-nav-link"]');
      await page.click('[data-testid="enable-realtime-notifications"]');
      
      // Wait for initial connection
      await page.waitForTimeout(1000);
      
      // Simulate connection loss
      await page.route('**/api/notifications/stream', route => {
        route.abort();
      });
      
      // Should attempt to reconnect
      await expect(page.locator('[data-testid="reconnecting-indicator"]')).toBeVisible();
      
      // Restore connection
      await page.unroute('**/api/notifications/stream');
      
      // Click retry or wait for auto-reconnect
      await page.click('[data-testid="retry-connection"]');
      
      // Should reconnect successfully
      await expect(page.locator('[data-testid="connected-indicator"]')).toBeVisible();
    });
  });

  test.describe('Notification Types', () => {
    
    test.beforeEach(async ({ page }) => {
      await authenticateUser(page);
    });

    test('should display performance insights correctly', async ({ page }) => {
      await page.click('[data-testid="notification-icon"]');
      
      // Look for performance insights
      const performanceInsights = page.locator('[data-testid="insight-item"][data-type="performance"]');
      
      if (await performanceInsights.count() > 0) {
        const insight = performanceInsights.first();
        
        // Should have performance-specific styling
        await expect(insight).toHaveClass(/performance/);
        
        // Should have appropriate icon
        await expect(insight.locator('[data-testid="performance-icon"]')).toBeVisible();
        
        // Should contain performance metrics
        await expect(insight.locator('[data-testid="insight-message"]')).toContainText(['improvement', 'performance', 'score']);
      }
    });

    test('should display streak insights correctly', async ({ page }) => {
      await page.click('[data-testid="notification-icon"]');
      
      const streakInsights = page.locator('[data-testid="insight-item"][data-type="streak"]');
      
      if (await streakInsights.count() > 0) {
        const insight = streakInsights.first();
        
        // Should have streak-specific styling
        await expect(insight).toHaveClass(/streak/);
        
        // Should contain streak information
        await expect(insight.locator('[data-testid="insight-message"]')).toContainText(['streak', 'game', 'win']);
      }
    });

    test('should display MMR insights correctly', async ({ page }) => {
      await page.click('[data-testid="notification-icon"]');
      
      const mmrInsights = page.locator('[data-testid="insight-item"][data-type="mmr"]');
      
      if (await mmrInsights.count() > 0) {
        const insight = mmrInsights.first();
        
        // Should contain MMR information
        await expect(insight.locator('[data-testid="insight-message"]')).toContainText(['MMR', 'rank', 'points']);
        
        // Should have action to view MMR page
        await expect(insight.locator('[data-testid="insight-action"]')).toHaveAttribute('href', /mmr/);
      }
    });

    test('should display recommendation insights correctly', async ({ page }) => {
      await page.click('[data-testid="notification-icon"]');
      
      const recommendationInsights = page.locator('[data-testid="insight-item"][data-type="recommendation"]');
      
      if (await recommendationInsights.count() > 0) {
        const insight = recommendationInsights.first();
        
        // Should contain recommendation information
        await expect(insight.locator('[data-testid="insight-message"]')).toContainText(['recommendation', 'improvement']);
        
        // Should have action to view recommendations
        await expect(insight.locator('[data-testid="insight-action"]')).toHaveAttribute('href', /recommendations/);
      }
    });
  });

  test.describe('Notification Statistics', () => {
    
    test.beforeEach(async ({ page }) => {
      await authenticateUser(page);
    });

    test('should display notification statistics', async ({ page }) => {
      await page.click('[data-testid="notification-icon"]');
      await page.click('[data-testid="notification-stats-tab"]');
      
      // Should show statistics
      await expect(page.locator('[data-testid="total-insights-stat"]')).toBeVisible();
      await expect(page.locator('[data-testid="unread-count-stat"]')).toBeVisible();
      await expect(page.locator('[data-testid="recent-count-stat"]')).toBeVisible();
      
      // Should show breakdown by type
      await expect(page.locator('[data-testid="insights-by-type-chart"]')).toBeVisible();
      
      // Should show breakdown by level
      await expect(page.locator('[data-testid="insights-by-level-chart"]')).toBeVisible();
    });

    test('should update statistics when insights are marked as read', async ({ page }) => {
      await page.click('[data-testid="notification-icon"]');
      
      // Get initial unread count
      const initialUnreadCount = await page.locator('[data-testid="notification-badge"]').textContent();
      
      if (initialUnreadCount && parseInt(initialUnreadCount) > 0) {
        // Mark one insight as read
        await page.locator('[data-testid="insight-item"]:has([data-testid="unread-indicator"])').first().click();
        
        // Check statistics
        await page.click('[data-testid="notification-stats-tab"]');
        
        const updatedUnreadStat = await page.locator('[data-testid="unread-count-stat"]').textContent();
        expect(parseInt(updatedUnreadStat)).toBe(parseInt(initialUnreadCount) - 1);
      }
    });
  });

  test.describe('Performance and Reliability', () => {
    
    test('should load insights quickly', async ({ page }) => {
      await authenticateUser(page);
      
      const startTime = Date.now();
      await page.click('[data-testid="notification-icon"]');
      await page.waitForSelector('[data-testid="insights-list"]');
      const loadTime = Date.now() - startTime;
      
      // Should load within 2 seconds
      expect(loadTime).toBeLessThan(2000);
    });

    test('should handle large numbers of insights', async ({ page }) => {
      await authenticateUser(page);
      
      // Mock API to return many insights
      await page.route('**/api/notifications/insights', route => {
        const insights = Array.from({ length: 100 }, (_, i) => ({
          id: i + 1,
          type: 'performance',
          level: 'info',
          title: `Test Insight ${i + 1}`,
          message: `This is test insight number ${i + 1}`,
          is_read: i % 3 === 0,
          created_at: new Date(Date.now() - i * 60000).toISOString()
        }));
        
        route.fulfill({
          json: {
            insights,
            total: insights.length,
            unread_count: insights.filter(i => !i.is_read).length
          }
        });
      });
      
      await page.click('[data-testid="notification-icon"]');
      
      // Should handle large list without performance issues
      await expect(page.locator('[data-testid="insights-list"]')).toBeVisible();
      await expect(page.locator('[data-testid="insight-item"]')).toHaveCountGreaterThan(50);
      
      // Scrolling should work smoothly
      await page.locator('[data-testid="insights-list"]').scrollIntoView();
      await expect(page.locator('[data-testid="insight-item"]').last()).toBeVisible();
    });

    test('should handle API errors gracefully', async ({ page }) => {
      await authenticateUser(page);
      
      // Mock API error
      await page.route('**/api/notifications/insights', route => {
        route.fulfill({ status: 500 });
      });
      
      await page.click('[data-testid="notification-icon"]');
      
      // Should show error message
      await expect(page.locator('[data-testid="insights-error"]')).toBeVisible();
      
      // Should provide retry option
      await expect(page.locator('[data-testid="retry-insights"]')).toBeVisible();
    });
  });

  test.describe('Mobile Responsiveness', () => {
    
    test.beforeEach(async ({ page }) => {
      // Set mobile viewport
      await page.setViewportSize({ width: 375, height: 667 });
      await authenticateUser(page);
    });

    test('should display notifications properly on mobile', async ({ page }) => {
      await page.click('[data-testid="notification-icon"]');
      
      // Should open as full-screen modal on mobile
      await expect(page.locator('[data-testid="insights-modal"]')).toBeVisible();
      
      // Should have mobile-optimized layout
      const insightItem = page.locator('[data-testid="insight-item"]').first();
      await expect(insightItem).toBeVisible();
      
      // Touch interactions should work
      await insightItem.tap();
      await expect(insightItem.locator('[data-testid="unread-indicator"]')).not.toBeVisible();
    });

    test('should handle swipe gestures for insight actions', async ({ page }) => {
      await page.click('[data-testid="notification-icon"]');
      
      const insightItem = page.locator('[data-testid="insight-item"]').first();
      
      // Swipe right to mark as read
      await insightItem.hover();
      await page.mouse.down();
      await page.mouse.move(100, 0);
      await page.mouse.up();
      
      // Should show action buttons
      await expect(insightItem.locator('[data-testid="swipe-actions"]')).toBeVisible();
    });
  });
});

// Helper functions

async function authenticateUser(page) {
  await page.fill('[data-testid="riot-id-input"]', TEST_USER.riotId);
  await page.fill('[data-testid="riot-tag-input"]', TEST_USER.riotTag);
  await page.selectOption('[data-testid="region-select"]', TEST_USER.region);
  await page.click('[data-testid="validate-account-button"]');
  await expect(page.locator('[data-testid="auth-success"]')).toBeVisible();
}

// Test configuration and setup
test.beforeAll(async () => {
  // Ensure test environment is ready
  console.log('Starting E2E tests for notification system...');
  console.log(`Base URL: ${BASE_URL}`);
  console.log(`API URL: ${API_BASE_URL}`);
});

test.afterAll(async () => {
  console.log('E2E tests completed.');
});