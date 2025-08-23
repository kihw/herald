// Herald.lol Dashboard E2E Tests

describe('Dashboard Flow', () => {
  beforeEach(() => {
    // Login before each test
    cy.intercept('POST', '**/auth/login', { fixture: 'auth-response.json' }).as('loginRequest');
    cy.intercept('GET', '**/auth/profile', { fixture: 'user-profile.json' }).as('getProfile');
    cy.intercept('GET', '**/health', { fixture: 'health-check.json' }).as('healthCheck');
    
    // Mock dashboard data endpoints
    cy.intercept('GET', '**/api/stats/overview', { fixture: 'dashboard-stats.json' }).as('getDashboardStats');
    cy.intercept('GET', '**/api/matches/recent', { fixture: 'recent-matches.json' }).as('getRecentMatches');
    cy.intercept('GET', '**/api/analytics/performance', { fixture: 'performance-analytics.json' }).as('getPerformanceAnalytics');
    
    // Login as test user
    cy.visit('/login');
    cy.fillLoginForm('dev@herald.lol', 'password123');
    cy.get('[data-testid="login-submit"]').click();
    cy.waitForAPI('@loginRequest');
  });

  it('should display dashboard with user stats', () => {
    cy.url().should('include', '/dashboard');
    
    // Check dashboard layout
    cy.get('[data-testid="dashboard-header"]').should('be.visible');
    cy.get('[data-testid="user-greeting"]').should('contain.text', 'Welcome back, Test Developer');
    
    // Check navigation
    cy.get('[data-testid="nav-dashboard"]').should('have.class', 'active');
    cy.get('[data-testid="nav-analytics"]').should('be.visible');
    cy.get('[data-testid="nav-matches"]').should('be.visible');
    cy.get('[data-testid="nav-champions"]').should('be.visible');
    
    // Wait for data to load
    cy.waitForAPI('@getDashboardStats');
    
    // Check key metrics cards
    cy.get('[data-testid="stat-current-rank"]').should('contain.text', 'Gold II');
    cy.get('[data-testid="stat-total-matches"]').should('contain.text', '150');
    cy.get('[data-testid="stat-win-rate"]').should('be.visible');
    cy.get('[data-testid="stat-kda-ratio"]').should('be.visible');
  });

  it('should display recent matches', () => {
    cy.waitForAPI('@getRecentMatches');
    
    // Check recent matches section
    cy.get('[data-testid="recent-matches"]').should('be.visible');
    cy.get('[data-testid="recent-matches-header"]').should('contain.text', 'Recent Matches');
    
    // Check match items
    cy.get('[data-testid^="match-item-"]').should('have.length.at.least', 1);
    cy.get('[data-testid^="match-item-"]').first().should('contain.text', 'Victory').or('contain.text', 'Defeat');
    
    // Check match details
    cy.get('[data-testid^="match-champion-"]').first().should('be.visible');
    cy.get('[data-testid^="match-kda-"]').first().should('be.visible');
    cy.get('[data-testid^="match-duration-"]').first().should('be.visible');
  });

  it('should navigate to detailed analytics', () => {
    cy.get('[data-testid="nav-analytics"]').click();
    cy.url().should('include', '/analytics');
    
    // Check analytics page elements
    cy.get('[data-testid="analytics-header"]').should('contain.text', 'Performance Analytics');
    cy.get('[data-testid="analytics-timeframe-selector"]').should('be.visible');
    cy.get('[data-testid="analytics-charts"]').should('be.visible');
  });

  it('should display performance analytics', () => {
    cy.waitForAPI('@getPerformanceAnalytics');
    
    // Check performance metrics
    cy.get('[data-testid="performance-overview"]').should('be.visible');
    cy.get('[data-testid="kda-trend"]').should('be.visible');
    cy.get('[data-testid="cs-per-minute"]').should('be.visible');
    cy.get('[data-testid="vision-score"]').should('be.visible');
    
    // Check champion performance
    cy.get('[data-testid="champion-performance"]').should('be.visible');
    cy.get('[data-testid="favorite-champion"]').should('contain.text', 'Jinx');
  });

  it('should handle loading states', () => {
    // Mock slow API response
    cy.intercept('GET', '**/api/stats/overview', { fixture: 'dashboard-stats.json', delay: 2000 }).as('slowDashboardStats');
    
    cy.visit('/dashboard');
    
    // Should show loading indicators
    cy.get('[data-testid="stats-loading"]').should('be.visible');
    cy.get('[data-testid="matches-loading"]').should('be.visible');
    
    // Loading should disappear after API completes
    cy.waitForAPI('@slowDashboardStats');
    cy.get('[data-testid="stats-loading"]').should('not.exist');
  });

  it('should handle API errors gracefully', () => {
    // Mock API error
    cy.intercept('GET', '**/api/stats/overview', {
      statusCode: 500,
      body: { error: 'Internal Server Error', message: 'Failed to fetch dashboard data' }
    }).as('dashboardError');
    
    cy.visit('/dashboard');
    
    // Should show error state
    cy.waitForAPI('@dashboardError');
    cy.get('[data-testid="dashboard-error"]').should('be.visible');
    cy.get('[data-testid="dashboard-error"]').should('contain.text', 'Failed to load dashboard data');
    
    // Should have retry button
    cy.get('[data-testid="retry-dashboard"]').should('be.visible');
  });

  it('should refresh data when retry button clicked', () => {
    // Mock initial error then success
    cy.intercept('GET', '**/api/stats/overview', {
      statusCode: 500,
      body: { error: 'Internal Server Error' }
    }).as('dashboardError');
    
    cy.visit('/dashboard');
    cy.waitForAPI('@dashboardError');
    
    // Click retry after setting up successful response
    cy.intercept('GET', '**/api/stats/overview', { fixture: 'dashboard-stats.json' }).as('dashboardRetry');
    cy.get('[data-testid="retry-dashboard"]').click();
    
    cy.waitForAPI('@dashboardRetry');
    cy.get('[data-testid="dashboard-error"]').should('not.exist');
    cy.get('[data-testid="stat-current-rank"]').should('be.visible');
  });

  it('should logout successfully', () => {
    // Mock logout endpoint
    cy.intercept('POST', '**/auth/logout', { statusCode: 200 }).as('logoutRequest');
    
    // Click user menu
    cy.get('[data-testid="user-menu"]').click();
    cy.get('[data-testid="logout-button"]').click();
    
    cy.waitForAPI('@logoutRequest');
    
    // Should redirect to login
    cy.url().should('include', '/login');
    cy.shouldNotBeAuthenticated();
  });

  it('should work on mobile viewports', () => {
    cy.viewport('iphone-x');
    cy.visit('/dashboard');
    
    // Check mobile layout
    cy.get('[data-testid="mobile-nav-toggle"]').should('be.visible');
    cy.get('[data-testid="dashboard-header"]').should('be.visible');
    
    // Navigation should be collapsible on mobile
    cy.get('[data-testid="mobile-nav-toggle"]').click();
    cy.get('[data-testid="nav-analytics"]').should('be.visible');
  });

  it('should display real-time updates', () => {
    // Mock SSE or WebSocket connection for real-time data
    cy.visit('/dashboard');
    
    // Check for real-time indicators
    cy.get('[data-testid="live-indicator"]').should('be.visible');
    cy.get('[data-testid="last-updated"]').should('be.visible');
  });
});