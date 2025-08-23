// Herald.lol Match Analysis Visual Regression Tests
import { test, expect, Page } from '@playwright/test';

test.describe('Match Analysis Visual Regression', () => {
  let page: Page;
  
  test.beforeEach(async ({ page: testPage }) => {
    page = testPage;
    
    // Mock match analysis API responses
    await page.route('**/api/matches/**', async (route) => {
      const url = route.request().url();
      
      if (url.includes('/analyze')) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            matchId: 'NA1_4567890123',
            analysis: {
              kda: { kills: 8, deaths: 3, assists: 12, ratio: 6.67 },
              cs: { total: 180, perMinute: 7.2, efficiency: 85 },
              vision: { score: 32, wardsPlaced: 15, wardsCleared: 8 },
              damage: { total: 25600, share: 0.32, efficiency: 89 },
              gold: { total: 12750, perMinute: 425, efficiency: 87 },
              performance: 'excellent',
              recommendations: [
                'Maintain aggressive playstyle in early game',
                'Improve ward coverage in enemy jungle',
                'Focus on objective control timing'
              ]
            },
            timeline: {
              phases: ['early', 'mid', 'late'],
              events: [
                { time: '3:45', type: 'first_blood', impact: 'high' },
                { time: '12:30', type: 'dragon', impact: 'medium' },
                { time: '28:15', type: 'baron', impact: 'game_winning' }
              ]
            },
            confidence: 0.94
          })
        });
      } else {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            matchId: 'NA1_4567890123',
            gameData: {
              champion: 'Jinx',
              role: 'ADC',
              duration: 1845, // 30:45
              result: 'victory',
              rank: 'Gold III',
              timestamp: '2024-01-20T15:30:00Z'
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
            username: 'VisualTestUser'
          },
          token: 'visual-test-token'
        })
      });
    });
    
    // Navigate to match analysis page
    await page.goto('/match/NA1_4567890123');
    
    // Wait for match analysis to load (Herald.lol <5s requirement)
    await page.waitForSelector('[data-testid="match-analysis-container"]', { timeout: 5000 });
    await page.waitForLoadState('networkidle');
  });
  
  test('Match Analysis Overview - Desktop', async () => {
    await page.setViewportSize({ width: 1920, height: 1080 });
    
    // Wait for all analysis components
    await expect(page.locator('[data-testid="match-header"]')).toBeVisible();
    await expect(page.locator('[data-testid="match-metrics"]')).toBeVisible();
    await expect(page.locator('[data-testid="match-recommendations"]')).toBeVisible();
    
    await expect(page).toHaveScreenshot('match-analysis-overview-desktop.png', {
      fullPage: true,
      animations: 'disabled',
    });
  });
  
  test('Match Analysis Overview - Mobile', async () => {
    await page.setViewportSize({ width: 390, height: 844 });
    await page.waitForTimeout(1000);
    
    await expect(page).toHaveScreenshot('match-analysis-overview-mobile.png', {
      fullPage: true,
      animations: 'disabled',
    });
  });
  
  test('Match Header Component', async () => {
    const matchHeader = page.locator('[data-testid="match-header"]');
    
    // Default header state
    await expect(matchHeader).toHaveScreenshot('match-header-default.png');
    
    // Victory state styling
    await expect(matchHeader.locator('[data-testid="match-result"]')).toContainText('Victory');
    await expect(matchHeader).toHaveScreenshot('match-header-victory.png');
  });
  
  test('Gaming Metrics Grid', async () => {
    const metricsGrid = page.locator('[data-testid="match-metrics"]');
    
    // Full metrics grid
    await expect(metricsGrid).toHaveScreenshot('match-metrics-grid.png');
    
    // Individual metric cards
    await expect(page.locator('[data-testid="kda-metric"]')).toHaveScreenshot('kda-metric-card.png');
    await expect(page.locator('[data-testid="cs-metric"]')).toHaveScreenshot('cs-metric-card.png');
    await expect(page.locator('[data-testid="vision-metric"]')).toHaveScreenshot('vision-metric-card.png');
    await expect(page.locator('[data-testid="damage-metric"]')).toHaveScreenshot('damage-metric-card.png');
    await expect(page.locator('[data-testid="gold-metric"]')).toHaveScreenshot('gold-metric-card.png');
  });
  
  test('Match Timeline Visualization', async () => {
    const timeline = page.locator('[data-testid="match-timeline"]');
    
    if (await timeline.isVisible()) {
      // Timeline default state
      await expect(timeline).toHaveScreenshot('match-timeline-default.png');
      
      // Click on timeline event
      const firstEvent = timeline.locator('[data-testid="timeline-event"]').first();
      if (await firstEvent.isVisible()) {
        await firstEvent.click();
        await page.waitForTimeout(500);
        await expect(timeline).toHaveScreenshot('match-timeline-event-selected.png');
      }
    }
  });
  
  test('Performance Analysis Chart', async () => {
    const performanceChart = page.locator('[data-testid="performance-chart"]');
    
    if (await performanceChart.isVisible()) {
      // Wait for chart rendering
      await page.waitForTimeout(2000);
      
      await expect(performanceChart).toHaveScreenshot('performance-chart.png');
      
      // Different chart views
      const chartTabs = page.locator('[data-testid="chart-tab"]');
      const tabCount = await chartTabs.count();
      
      for (let i = 0; i < tabCount; i++) {
        await chartTabs.nth(i).click();
        await page.waitForTimeout(1000);
        const tabName = await chartTabs.nth(i).textContent();
        await expect(performanceChart).toHaveScreenshot(`performance-chart-${tabName?.toLowerCase()}.png`);
      }
    }
  });
  
  test('Gaming Recommendations Panel', async () => {
    const recommendations = page.locator('[data-testid="match-recommendations"]');
    
    // Default recommendations panel
    await expect(recommendations).toHaveScreenshot('recommendations-panel.png');
    
    // Expand recommendation details
    const firstRecommendation = recommendations.locator('[data-testid="recommendation-item"]').first();
    if (await firstRecommendation.isVisible()) {
      await firstRecommendation.click();
      await page.waitForTimeout(500);
      await expect(recommendations).toHaveScreenshot('recommendations-expanded.png');
    }
  });
  
  test('Match Comparison View', async () => {
    const compareButton = page.locator('[data-testid="compare-matches"]');
    
    if (await compareButton.isVisible()) {
      await compareButton.click();
      await page.waitForTimeout(1000);
      
      // Comparison interface
      await expect(page.locator('[data-testid="match-comparison"]')).toBeVisible();
      await expect(page).toHaveScreenshot('match-comparison-view.png', {
        fullPage: true,
      });
    }
  });
  
  test('Gaming Heatmaps and Visualizations', async () => {
    // Vision heatmap
    const visionHeatmap = page.locator('[data-testid="vision-heatmap"]');
    if (await visionHeatmap.isVisible()) {
      await page.waitForTimeout(2000); // Wait for heatmap rendering
      await expect(visionHeatmap).toHaveScreenshot('vision-heatmap.png');
    }
    
    // Damage patterns
    const damagePatterns = page.locator('[data-testid="damage-patterns"]');
    if (await damagePatterns.isVisible()) {
      await expect(damagePatterns).toHaveScreenshot('damage-patterns.png');
    }
    
    // Movement tracking
    const movementMap = page.locator('[data-testid="movement-map"]');
    if (await movementMap.isVisible()) {
      await page.waitForTimeout(2000);
      await expect(movementMap).toHaveScreenshot('movement-map.png');
    }
  });
  
  test('Match Analysis Loading States', async () => {
    // Mock delayed API response
    await page.route('**/api/matches/*/analyze', async (route) => {
      await new Promise(resolve => setTimeout(resolve, 3000));
      await route.continue();
    });
    
    await page.reload();
    
    // Loading state screenshot
    await expect(page.locator('[data-testid="analysis-loading"]')).toBeVisible();
    await expect(page).toHaveScreenshot('match-analysis-loading.png', {
      fullPage: true,
    });
  });
  
  test('Match Analysis Error States', async () => {
    // Mock API error
    await page.route('**/api/matches/**', async (route) => {
      await route.fulfill({
        status: 404,
        contentType: 'application/json',
        body: JSON.stringify({
          error: 'Match not found',
          message: 'The requested match could not be found or analyzed.'
        })
      });
    });
    
    await page.reload();
    await page.waitForTimeout(2000);
    
    // Error state screenshot
    await expect(page).toHaveScreenshot('match-analysis-error.png', {
      fullPage: true,
    });
  });
  
  test('Gaming Performance Indicators', async () => {
    // Performance badges and indicators
    const performanceIndicators = page.locator('[data-testid="performance-indicators"]');
    
    if (await performanceIndicators.isVisible()) {
      await expect(performanceIndicators).toHaveScreenshot('performance-indicators.png');
      
      // Hover states for performance badges
      const badges = performanceIndicators.locator('[data-testid="performance-badge"]');
      const badgeCount = await badges.count();
      
      for (let i = 0; i < badgeCount; i++) {
        await badges.nth(i).hover();
        await page.waitForTimeout(300);
        await expect(badges.nth(i)).toHaveScreenshot(`performance-badge-hover-${i}.png`);
      }
    }
  });
  
  test('Match Analysis Responsive Breakpoints', async () => {
    const breakpoints = [
      { name: 'desktop-large', width: 1920, height: 1080 },
      { name: 'desktop-medium', width: 1440, height: 900 },
      { name: 'desktop-small', width: 1280, height: 720 },
      { name: 'tablet-large', width: 1024, height: 768 },
      { name: 'tablet-small', width: 768, height: 1024 },
      { name: 'mobile-large', width: 414, height: 896 },
      { name: 'mobile-small', width: 375, height: 667 },
    ];
    
    for (const breakpoint of breakpoints) {
      await page.setViewportSize({ width: breakpoint.width, height: breakpoint.height });
      await page.waitForTimeout(1000);
      
      await expect(page).toHaveScreenshot(`match-analysis-${breakpoint.name}.png`, {
        fullPage: true,
        animations: 'disabled',
      });
    }
  });
  
  test('Gaming Theme Variations', async () => {
    // Light theme (default)
    await expect(page).toHaveScreenshot('match-analysis-light-theme.png', {
      fullPage: true,
    });
    
    // Dark theme (if available)
    const themeToggle = page.locator('[data-testid="theme-toggle"]');
    if (await themeToggle.isVisible()) {
      await themeToggle.click();
      await page.waitForTimeout(500);
      
      await expect(page).toHaveScreenshot('match-analysis-dark-theme.png', {
        fullPage: true,
      });
    }
  });
});