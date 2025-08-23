// Herald.lol Gaming Components Visual Regression Tests
import { test, expect, Page } from '@playwright/test';

test.describe('Gaming Components Visual Regression', () => {
  let page: Page;
  
  test.beforeEach(async ({ page: testPage }) => {
    page = testPage;
    
    // Mock component APIs
    await page.route('**/api/**', async (route) => {
      const url = route.request().url();
      
      // Mock team composition data
      if (url.includes('/team-composition')) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            recommendations: [
              { champion: 'Jinx', role: 'ADC', synergy: 0.92, confidence: 0.89 },
              { champion: 'Thresh', role: 'Support', synergy: 0.87, confidence: 0.85 },
              { champion: 'Graves', role: 'Jungle', synergy: 0.81, confidence: 0.82 },
            ],
            strategy: 'meta_optimal',
            winRate: 0.734,
            confidence: 0.91
          })
        });
      }
      
      // Mock counter-pick data
      else if (url.includes('/counter-picks')) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            counters: [
              { champion: 'Caitlyn', effectiveness: 0.78, matchup: 'favorable', confidence: 0.85 },
              { champion: 'Draven', effectiveness: 0.71, matchup: 'slightly_favorable', confidence: 0.79 },
              { champion: 'Ezreal', effectiveness: 0.65, matchup: 'neutral', confidence: 0.74 },
            ],
            targetChampion: 'Jinx',
            role: 'ADC',
            confidence: 0.82
          })
        });
      }
      
      // Mock skill progression data
      else if (url.includes('/skill-progression')) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            skills: {
              mechanical: { current: 73, trend: 'improving', target: 80 },
              tactical: { current: 68, trend: 'stable', target: 75 },
              strategic: { current: 71, trend: 'improving', target: 78 },
              mental: { current: 65, trend: 'declining', target: 72 },
            },
            overallProgress: 0.74,
            nextMilestone: 'Gold I',
            confidence: 0.87
          })
        });
      }
      
      // Default mock response
      else {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ message: 'Mock response' })
        });
      }
    });
    
    // Navigate to components showcase (create if doesn't exist)
    await page.goto('/components-showcase');
    await page.waitForLoadState('networkidle');
  });
  
  test('Team Composition Component', async () => {
    // Mock team composition component page
    await page.goto('/team-composition');
    await page.waitForSelector('[data-testid="team-composition-optimizer"]', { timeout: 5000 });
    
    const teamComp = page.locator('[data-testid="team-composition-optimizer"]');
    
    // Default state
    await expect(teamComp).toHaveScreenshot('team-composition-default.png');
    
    // With recommendations
    const optimizeButton = page.locator('[data-testid="optimize-composition"]');
    if (await optimizeButton.isVisible()) {
      await optimizeButton.click();
      await page.waitForTimeout(2000); // Wait for optimization
      await expect(teamComp).toHaveScreenshot('team-composition-optimized.png');
    }
    
    // Different strategies
    const strategies = ['meta_optimal', 'synergy_focused', 'balanced', 'comfort_picks'];
    
    for (const strategy of strategies) {
      const strategySelect = page.locator(`[data-testid="strategy-${strategy}"]`);
      if (await strategySelect.isVisible()) {
        await strategySelect.click();
        await page.waitForTimeout(1000);
        await expect(teamComp).toHaveScreenshot(`team-composition-${strategy}.png`);
      }
    }
  });
  
  test('Counter-Pick Analysis Component', async () => {
    await page.goto('/counter-picks');
    await page.waitForSelector('[data-testid="counter-pick-analyzer"]', { timeout: 5000 });
    
    const counterPick = page.locator('[data-testid="counter-pick-analyzer"]');
    
    // Default state
    await expect(counterPick).toHaveScreenshot('counter-pick-default.png');
    
    // Champion selection
    const championSelect = page.locator('[data-testid="champion-select"]');
    if (await championSelect.isVisible()) {
      await championSelect.click();
      await page.locator('[data-testid="champion-option-jinx"]').click();
      await page.waitForTimeout(1000);
      await expect(counterPick).toHaveScreenshot('counter-pick-jinx-selected.png');
    }
    
    // Counter recommendations display
    const recommendations = page.locator('[data-testid="counter-recommendations"]');
    if (await recommendations.isVisible()) {
      await expect(recommendations).toHaveScreenshot('counter-pick-recommendations.png');
      
      // Hover states for counter champions
      const counterChampions = recommendations.locator('[data-testid="counter-champion"]');
      const champCount = await counterChampions.count();
      
      for (let i = 0; i < Math.min(champCount, 3); i++) {
        await counterChampions.nth(i).hover();
        await page.waitForTimeout(300);
        await expect(counterChampions.nth(i)).toHaveScreenshot(`counter-champion-hover-${i}.png`);
      }
    }
  });
  
  test('Skill Progression Component', async () => {
    await page.goto('/skill-progression');
    await page.waitForSelector('[data-testid="skill-progression-tracker"]', { timeout: 5000 });
    
    const skillTracker = page.locator('[data-testid="skill-progression-tracker"]');
    
    // Default progression view
    await expect(skillTracker).toHaveScreenshot('skill-progression-default.png');
    
    // Individual skill categories
    const skillCategories = ['mechanical', 'tactical', 'strategic', 'mental'];
    
    for (const category of skillCategories) {
      const skillCard = page.locator(`[data-testid="skill-${category}"]`);
      if (await skillCard.isVisible()) {
        await expect(skillCard).toHaveScreenshot(`skill-${category}-card.png`);
        
        // Expanded view
        await skillCard.click();
        await page.waitForTimeout(500);
        await expect(skillCard).toHaveScreenshot(`skill-${category}-expanded.png`);
      }
    }
    
    // Overall progress chart
    const progressChart = page.locator('[data-testid="overall-progress-chart"]');
    if (await progressChart.isVisible()) {
      await page.waitForTimeout(2000); // Wait for chart rendering
      await expect(progressChart).toHaveScreenshot('overall-progress-chart.png');
    }
  });
  
  test('Gaming Charts and Visualizations', async () => {
    // KDA Trend Chart
    const kdaChart = page.locator('[data-testid="kda-trend-chart"]');
    if (await kdaChart.isVisible()) {
      await page.waitForTimeout(2000);
      await expect(kdaChart).toHaveScreenshot('kda-trend-chart.png');
    }
    
    // CS/min Performance Chart
    const csChart = page.locator('[data-testid="cs-performance-chart"]');
    if (await csChart.isVisible()) {
      await page.waitForTimeout(2000);
      await expect(csChart).toHaveScreenshot('cs-performance-chart.png');
    }
    
    // Vision Control Radar Chart
    const visionRadar = page.locator('[data-testid="vision-radar-chart"]');
    if (await visionRadar.isVisible()) {
      await page.waitForTimeout(2000);
      await expect(visionRadar).toHaveScreenshot('vision-radar-chart.png');
    }
    
    // Damage Distribution Pie Chart
    const damagePie = page.locator('[data-testid="damage-distribution-chart"]');
    if (await damagePie.isVisible()) {
      await page.waitForTimeout(2000);
      await expect(damagePie).toHaveScreenshot('damage-distribution-chart.png');
    }
    
    // Gold Timeline Chart
    const goldTimeline = page.locator('[data-testid="gold-timeline-chart"]');
    if (await goldTimeline.isVisible()) {
      await page.waitForTimeout(2000);
      await expect(goldTimeline).toHaveScreenshot('gold-timeline-chart.png');
    }
  });
  
  test('Gaming UI Cards and Widgets', async () => {
    // Metric Cards
    const metricCards = [
      'rank-card',
      'winrate-card', 
      'recent-matches-card',
      'champion-mastery-card',
      'performance-summary-card'
    ];
    
    for (const cardType of metricCards) {
      const card = page.locator(`[data-testid="${cardType}"]`);
      if (await card.isVisible()) {
        await expect(card).toHaveScreenshot(`${cardType}.png`);
        
        // Hover state
        await card.hover();
        await page.waitForTimeout(300);
        await expect(card).toHaveScreenshot(`${cardType}-hover.png`);
      }
    }
  });
  
  test('Gaming Form Components', async () => {
    // Summoner Search Form
    const summonerSearch = page.locator('[data-testid="summoner-search-form"]');
    if (await summonerSearch.isVisible()) {
      await expect(summonerSearch).toHaveScreenshot('summoner-search-form.png');
      
      // With input
      const searchInput = summonerSearch.locator('input[type="text"]');
      await searchInput.fill('TestSummoner');
      await expect(summonerSearch).toHaveScreenshot('summoner-search-form-filled.png');
      
      // With suggestions
      await page.waitForTimeout(500);
      const suggestions = page.locator('[data-testid="search-suggestions"]');
      if (await suggestions.isVisible()) {
        await expect(suggestions).toHaveScreenshot('summoner-search-suggestions.png');
      }
    }
    
    // Champion Filter Form
    const championFilter = page.locator('[data-testid="champion-filter-form"]');
    if (await championFilter.isVisible()) {
      await expect(championFilter).toHaveScreenshot('champion-filter-form.png');
      
      // With filters applied
      const roleFilter = championFilter.locator('[data-testid="role-filter"]');
      if (await roleFilter.isVisible()) {
        await roleFilter.selectOption('ADC');
        await expect(championFilter).toHaveScreenshot('champion-filter-form-adc.png');
      }
    }
  });
  
  test('Gaming Navigation Components', async () => {
    // Main Navigation
    const mainNav = page.locator('[data-testid="main-navigation"]');
    if (await mainNav.isVisible()) {
      await expect(mainNav).toHaveScreenshot('main-navigation.png');
      
      // Active states
      const navItems = mainNav.locator('[data-testid="nav-item"]');
      const navCount = await navItems.count();
      
      for (let i = 0; i < navCount; i++) {
        await navItems.nth(i).click();
        await page.waitForTimeout(300);
        await expect(mainNav).toHaveScreenshot(`main-navigation-active-${i}.png`);
      }
    }
    
    // Gaming Breadcrumbs
    const breadcrumbs = page.locator('[data-testid="breadcrumbs"]');
    if (await breadcrumbs.isVisible()) {
      await expect(breadcrumbs).toHaveScreenshot('gaming-breadcrumbs.png');
    }
    
    // Tab Navigation
    const tabNav = page.locator('[data-testid="tab-navigation"]');
    if (await tabNav.isVisible()) {
      await expect(tabNav).toHaveScreenshot('tab-navigation.png');
      
      const tabs = tabNav.locator('[data-testid="tab"]');
      const tabCount = await tabs.count();
      
      for (let i = 0; i < tabCount; i++) {
        await tabs.nth(i).click();
        await page.waitForTimeout(300);
        await expect(tabNav).toHaveScreenshot(`tab-navigation-${i}.png`);
      }
    }
  });
  
  test('Gaming Modal and Dialog Components', async () => {
    // Settings Modal
    const settingsButton = page.locator('[data-testid="settings-button"]');
    if (await settingsButton.isVisible()) {
      await settingsButton.click();
      await page.waitForTimeout(500);
      
      const settingsModal = page.locator('[data-testid="settings-modal"]');
      if (await settingsModal.isVisible()) {
        await expect(settingsModal).toHaveScreenshot('settings-modal.png');
      }
    }
    
    // Match Details Dialog
    const matchButton = page.locator('[data-testid="match-details-button"]');
    if (await matchButton.isVisible()) {
      await matchButton.click();
      await page.waitForTimeout(500);
      
      const matchDialog = page.locator('[data-testid="match-details-dialog"]');
      if (await matchDialog.isVisible()) {
        await expect(matchDialog).toHaveScreenshot('match-details-dialog.png');
      }
    }
    
    // Confirmation Dialog
    const deleteButton = page.locator('[data-testid="delete-button"]');
    if (await deleteButton.isVisible()) {
      await deleteButton.click();
      await page.waitForTimeout(300);
      
      const confirmDialog = page.locator('[data-testid="confirm-dialog"]');
      if (await confirmDialog.isVisible()) {
        await expect(confirmDialog).toHaveScreenshot('confirmation-dialog.png');
      }
    }
  });
  
  test('Gaming Loading and Empty States', async () => {
    // Loading states
    const loadingSpinner = page.locator('[data-testid="loading-spinner"]');
    if (await loadingSpinner.isVisible()) {
      await expect(loadingSpinner).toHaveScreenshot('loading-spinner.png');
    }
    
    const loadingSkeleton = page.locator('[data-testid="loading-skeleton"]');
    if (await loadingSkeleton.isVisible()) {
      await expect(loadingSkeleton).toHaveScreenshot('loading-skeleton.png');
    }
    
    // Empty states
    const emptyMatches = page.locator('[data-testid="empty-matches"]');
    if (await emptyMatches.isVisible()) {
      await expect(emptyMatches).toHaveScreenshot('empty-matches-state.png');
    }
    
    const noData = page.locator('[data-testid="no-data"]');
    if (await noData.isVisible()) {
      await expect(noData).toHaveScreenshot('no-data-state.png');
    }
  });
  
  test('Gaming Error Components', async () => {
    // Error boundary
    const errorBoundary = page.locator('[data-testid="error-boundary"]');
    if (await errorBoundary.isVisible()) {
      await expect(errorBoundary).toHaveScreenshot('error-boundary.png');
    }
    
    // API error message
    const apiError = page.locator('[data-testid="api-error"]');
    if (await apiError.isVisible()) {
      await expect(apiError).toHaveScreenshot('api-error-message.png');
    }
    
    // Network error
    const networkError = page.locator('[data-testid="network-error"]');
    if (await networkError.isVisible()) {
      await expect(networkError).toHaveScreenshot('network-error.png');
    }
  });
  
  test('Gaming Responsive Component Behavior', async () => {
    const components = [
      'team-composition-optimizer',
      'counter-pick-analyzer', 
      'skill-progression-tracker',
      'analytics-dashboard'
    ];
    
    const breakpoints = [
      { name: 'mobile', width: 375, height: 667 },
      { name: 'tablet', width: 768, height: 1024 },
      { name: 'desktop', width: 1920, height: 1080 }
    ];
    
    for (const component of components) {
      const element = page.locator(`[data-testid="${component}"]`);
      
      if (await element.isVisible()) {
        for (const breakpoint of breakpoints) {
          await page.setViewportSize({ width: breakpoint.width, height: breakpoint.height });
          await page.waitForTimeout(1000);
          
          await expect(element).toHaveScreenshot(`${component}-${breakpoint.name}.png`);
        }
      }
    }
  });
});