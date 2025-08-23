// Herald.lol Custom Cypress Commands

/// <reference types="cypress" />

declare global {
  namespace Cypress {
    interface Chainable {
      /**
       * Custom command to login as a test user
       * @example cy.loginAsTestUser()
       */
      loginAsTestUser(): Chainable<Element>;

      /**
       * Custom command to register a new test user
       * @example cy.registerTestUser('test@example.com', 'testuser', 'password123')
       */
      registerTestUser(email: string, username: string, password: string): Chainable<Element>;

      /**
       * Custom command to wait for API response
       * @example cy.waitForAPI('@getProfile')
       */
      waitForAPI(alias: string): Chainable<Element>;

      /**
       * Custom command to check if user is authenticated
       * @example cy.shouldBeAuthenticated()
       */
      shouldBeAuthenticated(): Chainable<Element>;

      /**
       * Custom command to check if user is not authenticated
       * @example cy.shouldNotBeAuthenticated()
       */
      shouldNotBeAuthenticated(): Chainable<Element>;

      /**
       * Custom command to fill login form
       * @example cy.fillLoginForm('test@example.com', 'password123')
       */
      fillLoginForm(email: string, password: string): Chainable<Element>;

      /**
       * Custom command to fill registration form
       * @example cy.fillRegistrationForm('test@example.com', 'testuser', 'password123')
       */
      fillRegistrationForm(email: string, username: string, password: string, displayName?: string): Chainable<Element>;

      /**
       * Custom command to navigate to dashboard
       * @example cy.goToDashboard()
       */
      goToDashboard(): Chainable<Element>;

      /**
       * Custom command to mock Riot API responses
       * @example cy.mockRiotAPI()
       */
      mockRiotAPI(): Chainable<Element>;

      /**
       * Custom command to check loading state
       * @example cy.shouldBeLoading()
       */
      shouldBeLoading(): Chainable<Element>;

      /**
       * Custom command to wait for loading to complete
       * @example cy.waitForLoadingToComplete()
       */
      waitForLoadingToComplete(): Chainable<Element>;

      /**
       * Gaming-specific custom commands for Herald.lol
       */

      /**
       * Custom command to wait for analytics to load (max 5s)
       * @example cy.waitForAnalytics()
       */
      waitForAnalytics(): Chainable<Element>;

      /**
       * Custom command to validate gaming metrics
       * @example cy.validateGamingMetrics('kda')
       */
      validateGamingMetrics(metricType: string): Chainable<Element>;

      /**
       * Custom command to check analytics performance
       * @example cy.checkAnalyticsPerformance(5000)
       */
      checkAnalyticsPerformance(maxTimeMs: number): Chainable<Element>;

      /**
       * Custom command to navigate to analytics dashboard
       * @example cy.goToAnalyticsDashboard()
       */
      goToAnalyticsDashboard(): Chainable<Element>;

      /**
       * Custom command to select gaming metric
       * @example cy.selectGamingMetric('kda')
       */
      selectGamingMetric(metric: string): Chainable<Element>;

      /**
       * Custom command to validate KDA display
       * @example cy.validateKDADisplay()
       */
      validateKDADisplay(): Chainable<Element>;

      /**
       * Custom command to validate CS/min display
       * @example cy.validateCSDisplay()
       */
      validateCSDisplay(): Chainable<Element>;

      /**
       * Custom command to validate vision score
       * @example cy.validateVisionScore()
       */
      validateVisionScore(): Chainable<Element>;

      /**
       * Custom command to validate damage metrics
       * @example cy.validateDamageMetrics()
       */
      validateDamageMetrics(): Chainable<Element>;

      /**
       * Custom command to validate gold efficiency
       * @example cy.validateGoldEfficiency()
       */
      validateGoldEfficiency(): Chainable<Element>;

      /**
       * Custom command to mock gaming analytics data
       * @example cy.mockGamingAnalytics()
       */
      mockGamingAnalytics(): Chainable<Element>;

      /**
       * Custom command to simulate match analysis
       * @example cy.simulateMatchAnalysis()
       */
      simulateMatchAnalysis(): Chainable<Element>;

      /**
       * Custom command to validate gaming performance
       * @example cy.validateGamingPerformance()
       */
      validateGamingPerformance(): Chainable<Element>;
    }
  }
}

// Login as test user
Cypress.Commands.add('loginAsTestUser', () => {
  const email = Cypress.env('TEST_USER_EMAIL');
  const password = Cypress.env('TEST_USER_PASSWORD');

  // Intercept login API call
  cy.intercept('POST', '**/auth/login', { fixture: 'auth-response.json' }).as('loginRequest');
  cy.intercept('GET', '**/auth/profile', { fixture: 'user-profile.json' }).as('getProfile');

  cy.visit('/login');
  cy.fillLoginForm(email, password);
  cy.get('[data-testid="login-submit"]').click();
  
  cy.waitForAPI('@loginRequest');
  cy.shouldBeAuthenticated();
});

// Register test user
Cypress.Commands.add('registerTestUser', (email: string, username: string, password: string) => {
  // Intercept registration API call
  cy.intercept('POST', '**/auth/register', { fixture: 'auth-response.json' }).as('registerRequest');
  cy.intercept('GET', '**/auth/profile', { fixture: 'user-profile.json' }).as('getProfile');

  cy.visit('/register');
  cy.fillRegistrationForm(email, username, password);
  cy.get('[data-testid="register-submit"]').click();
  
  cy.waitForAPI('@registerRequest');
  cy.shouldBeAuthenticated();
});

// Wait for API response
Cypress.Commands.add('waitForAPI', (alias: string) => {
  cy.wait(alias, { timeout: 10000 });
});

// Check if authenticated
Cypress.Commands.add('shouldBeAuthenticated', () => {
  cy.url().should('not.contain', '/login');
  cy.get('[data-testid="user-menu"]', { timeout: 10000 }).should('exist');
});

// Check if not authenticated
Cypress.Commands.add('shouldNotBeAuthenticated', () => {
  cy.get('[data-testid="login-form"]').should('exist');
  cy.url().should('contain', '/login');
});

// Fill login form
Cypress.Commands.add('fillLoginForm', (email: string, password: string) => {
  cy.get('[data-testid="email-input"]').type(email);
  cy.get('[data-testid="password-input"]').type(password);
});

// Fill registration form
Cypress.Commands.add('fillRegistrationForm', (email: string, username: string, password: string, displayName?: string) => {
  cy.get('[data-testid="email-input"]').type(email);
  cy.get('[data-testid="username-input"]').type(username);
  if (displayName) {
    cy.get('[data-testid="display-name-input"]').type(displayName);
  }
  cy.get('[data-testid="password-input"]').type(password);
  cy.get('[data-testid="confirm-password-input"]').type(password);
});

// Navigate to dashboard
Cypress.Commands.add('goToDashboard', () => {
  cy.intercept('GET', '**/auth/profile', { fixture: 'user-profile.json' }).as('getProfile');
  cy.intercept('GET', '**/analytics/**', { fixture: 'analytics-stats.json' }).as('getAnalytics');
  
  cy.visit('/dashboard');
  cy.waitForAPI('@getProfile');
});

// Mock Riot API responses
Cypress.Commands.add('mockRiotAPI', () => {
  // Mock common Riot API endpoints
  cy.intercept('GET', '**/riot/accounts', { fixture: 'riot-accounts.json' }).as('getRiotAccounts');
  cy.intercept('POST', '**/riot/link', { fixture: 'riot-account-link.json' }).as('linkRiotAccount');
  cy.intercept('GET', '**/riot/matches**', { fixture: 'match-history.json' }).as('getMatches');
  cy.intercept('POST', '**/riot/summoner', { fixture: 'summoner-info.json' }).as('getSummonerInfo');
  cy.intercept('POST', '**/riot/ranked', { fixture: 'ranked-info.json' }).as('getRankedInfo');
});

// Check loading state
Cypress.Commands.add('shouldBeLoading', () => {
  cy.get('[data-testid="loading-spinner"]').should('exist');
});

// Wait for loading to complete
Cypress.Commands.add('waitForLoadingToComplete', () => {
  cy.get('[data-testid="loading-spinner"]').should('not.exist', { timeout: 15000 });
});

// Gaming-specific command implementations for Herald.lol

// Wait for analytics to load (Herald.lol <5s requirement)
Cypress.Commands.add('waitForAnalytics', () => {
  const startTime = Date.now();
  
  cy.get('[data-testid="analytics-dashboard"]', { timeout: 5000 }).should('be.visible').then(() => {
    const loadTime = Date.now() - startTime;
    expect(loadTime).to.be.lessThan(5000, `Analytics should load in <5s (took ${loadTime}ms)`);
  });
  
  // Wait for all analytics components to be ready
  cy.get('[data-testid="kda-widget"]').should('be.visible');
  cy.get('[data-testid="cs-widget"]').should('be.visible');
  cy.get('[data-testid="vision-widget"]').should('be.visible');
  cy.get('[data-testid="damage-widget"]').should('be.visible');
  cy.get('[data-testid="gold-widget"]').should('be.visible');
});

// Validate gaming metrics
Cypress.Commands.add('validateGamingMetrics', (metricType: string) => {
  switch (metricType) {
    case 'kda':
      cy.validateKDADisplay();
      break;
    case 'cs':
      cy.validateCSDisplay();
      break;
    case 'vision':
      cy.validateVisionScore();
      break;
    case 'damage':
      cy.validateDamageMetrics();
      break;
    case 'gold':
      cy.validateGoldEfficiency();
      break;
    default:
      throw new Error(`Unknown metric type: ${metricType}`);
  }
});

// Check analytics performance (Herald.lol <5s requirement)
Cypress.Commands.add('checkAnalyticsPerformance', (maxTimeMs: number = 5000) => {
  const startTime = Date.now();
  
  cy.window().then((win) => {
    // Check if performance timing is available
    if (win.performance && win.performance.timing) {
      const loadTime = win.performance.timing.loadEventEnd - win.performance.timing.navigationStart;
      expect(loadTime).to.be.lessThan(maxTimeMs, `Page load should be <${maxTimeMs}ms (was ${loadTime}ms)`);
    }
  });
  
  // Check that analytics widgets render within time limit
  cy.get('[data-testid="analytics-dashboard"]').should('be.visible').then(() => {
    const totalTime = Date.now() - startTime;
    expect(totalTime).to.be.lessThan(maxTimeMs, `Analytics should render in <${maxTimeMs}ms (took ${totalTime}ms)`);
  });
});

// Navigate to analytics dashboard
Cypress.Commands.add('goToAnalyticsDashboard', () => {
  // Mock analytics API endpoints
  cy.mockGamingAnalytics();
  
  cy.visit('/analytics');
  cy.waitForAnalytics();
});

// Select gaming metric
Cypress.Commands.add('selectGamingMetric', (metric: string) => {
  cy.get(`[data-testid="${metric}-tab"]`).click();
  cy.get(`[data-testid="${metric}-widget"]`).should('be.visible');
});

// Validate KDA display
Cypress.Commands.add('validateKDADisplay', () => {
  cy.get('[data-testid="kda-widget"]').within(() => {
    // Check KDA value is displayed and valid
    cy.get('[data-testid="kda-value"]').should('be.visible').invoke('text').then((text) => {
      const kda = parseFloat(text);
      expect(kda).to.be.a('number');
      expect(kda).to.be.at.least(0);
      expect(kda).to.be.at.most(50); // Reasonable upper bound
    });
    
    // Check trend indicator
    cy.get('[data-testid="kda-trend"]').should('be.visible');
    
    // Check percentile
    cy.get('[data-testid="kda-percentile"]').should('be.visible').invoke('text').then((text) => {
      const percentile = parseInt(text.replace('%', ''));
      expect(percentile).to.be.at.least(0);
      expect(percentile).to.be.at.most(100);
    });
    
    // Check confidence level
    cy.get('[data-testid="confidence-level"]').should('be.visible');
  });
});

// Validate CS/min display
Cypress.Commands.add('validateCSDisplay', () => {
  cy.get('[data-testid="cs-widget"]').within(() => {
    // Check CS/min value is displayed and valid
    cy.get('[data-testid="cs-value"]').should('be.visible').invoke('text').then((text) => {
      const cs = parseFloat(text);
      expect(cs).to.be.a('number');
      expect(cs).to.be.at.least(0);
      expect(cs).to.be.at.most(15); // Reasonable upper bound for CS/min
    });
    
    // Check efficiency rating
    cy.get('[data-testid="cs-efficiency"]').should('be.visible').invoke('text').then((text) => {
      const efficiency = parseInt(text.replace('%', ''));
      expect(efficiency).to.be.at.least(0);
      expect(efficiency).to.be.at.most(100);
    });
    
    // Check benchmarks
    cy.get('[data-testid="cs-benchmarks"]').should('be.visible');
    
    // Check recommendations
    cy.get('[data-testid="cs-recommendations"]').should('be.visible');
  });
});

// Validate vision score
Cypress.Commands.add('validateVisionScore', () => {
  cy.get('[data-testid="vision-widget"]').within(() => {
    // Check vision score value
    cy.get('[data-testid="vision-score"]').should('be.visible').invoke('text').then((text) => {
      const score = parseFloat(text);
      expect(score).to.be.a('number');
      expect(score).to.be.at.least(0);
      expect(score).to.be.at.most(200); // Reasonable upper bound
    });
    
    // Check vision control percentage
    cy.get('[data-testid="vision-control"]').should('be.visible');
    
    // Check heatmap
    cy.get('[data-testid="vision-heatmap"]').should('be.visible');
    
    // Check warding patterns
    cy.get('[data-testid="warding-patterns"]').should('be.visible');
  });
});

// Validate damage metrics
Cypress.Commands.add('validateDamageMetrics', () => {
  cy.get('[data-testid="damage-widget"]').within(() => {
    // Check damage share
    cy.get('[data-testid="damage-share"]').should('be.visible').invoke('text').then((text) => {
      const share = parseFloat(text.replace('%', '')) / 100;
      expect(share).to.be.at.least(0);
      expect(share).to.be.at.most(1);
    });
    
    // Check damage per minute
    cy.get('[data-testid="damage-per-minute"]').should('be.visible').invoke('text').then((text) => {
      const dpm = parseFloat(text);
      expect(dpm).to.be.a('number');
      expect(dpm).to.be.at.least(0);
    });
    
    // Check damage breakdown
    cy.get('[data-testid="damage-breakdown"]').should('be.visible');
    
    // Check team fight contribution
    cy.get('[data-testid="teamfight-contribution"]').should('be.visible');
  });
});

// Validate gold efficiency
Cypress.Commands.add('validateGoldEfficiency', () => {
  cy.get('[data-testid="gold-widget"]').within(() => {
    // Check gold per minute
    cy.get('[data-testid="gold-per-minute"]').should('be.visible').invoke('text').then((text) => {
      const gpm = parseFloat(text);
      expect(gpm).to.be.a('number');
      expect(gpm).to.be.at.least(0);
      expect(gpm).to.be.at.most(1000); // Reasonable upper bound
    });
    
    // Check gold efficiency percentage
    cy.get('[data-testid="gold-efficiency"]').should('be.visible').invoke('text').then((text) => {
      const efficiency = parseInt(text.replace('%', ''));
      expect(efficiency).to.be.at.least(0);
      expect(efficiency).to.be.at.most(100);
    });
    
    // Check power spikes
    cy.get('[data-testid="power-spikes"]').should('be.visible');
    
    // Check backing efficiency
    cy.get('[data-testid="backing-efficiency"]').should('be.visible');
  });
});

// Mock gaming analytics data
Cypress.Commands.add('mockGamingAnalytics', () => {
  // Mock all gaming analytics endpoints with fixtures
  cy.intercept('GET', '**/analytics/kda**', { fixture: 'analytics/kda-analysis.json' }).as('getKDAAnalytics');
  cy.intercept('GET', '**/analytics/cs**', { fixture: 'analytics/cs-analysis.json' }).as('getCSAnalytics');
  cy.intercept('GET', '**/analytics/vision**', { fixture: 'analytics/vision-analysis.json' }).as('getVisionAnalytics');
  cy.intercept('GET', '**/analytics/damage**', { fixture: 'analytics/damage-analysis.json' }).as('getDamageAnalytics');
  cy.intercept('GET', '**/analytics/gold**', { fixture: 'analytics/gold-analysis.json' }).as('getGoldAnalytics');
  
  // Mock match analysis endpoints
  cy.intercept('GET', '**/matches/**', { fixture: 'match-analysis.json' }).as('getMatchAnalysis');
  cy.intercept('POST', '**/analytics/analyze-match', { fixture: 'match-analysis-result.json' }).as('analyzeMatch');
  
  // Mock team composition endpoints
  cy.intercept('GET', '**/team-composition/**', { fixture: 'team-composition.json' }).as('getTeamComposition');
  cy.intercept('POST', '**/team-composition/optimize', { fixture: 'team-composition-result.json' }).as('optimizeComposition');
});

// Simulate match analysis
Cypress.Commands.add('simulateMatchAnalysis', () => {
  cy.mockGamingAnalytics();
  
  // Visit match analysis page
  cy.visit('/match/NA1_4567890123');
  
  // Wait for match data to load
  cy.get('[data-testid="match-analysis-container"]', { timeout: 5000 }).should('be.visible');
  
  // Validate all gaming metrics are displayed
  cy.validateGamingMetrics('kda');
  cy.validateGamingMetrics('cs');
  cy.validateGamingMetrics('vision');
  cy.validateGamingMetrics('damage');
  cy.validateGamingMetrics('gold');
  
  // Check that recommendations are provided
  cy.get('[data-testid="match-recommendations"]').should('be.visible');
  
  // Verify performance requirements
  cy.checkAnalyticsPerformance(5000);
});

// Validate gaming performance
Cypress.Commands.add('validateGamingPerformance', () => {
  // Check that the app meets Herald.lol performance requirements
  cy.window().then((win) => {
    // Verify analytics load time
    cy.get('[data-testid="analytics-dashboard"]').should('be.visible');
    
    // Check memory usage (should be reasonable for gaming analytics)
    if (win.performance && win.performance.memory) {
      const memUsage = win.performance.memory.usedJSHeapSize / 1024 / 1024; // MB
      expect(memUsage).to.be.lessThan(100, `Memory usage should be <100MB (was ${memUsage.toFixed(2)}MB)`);
    }
    
    // Validate that all gaming widgets are responsive
    cy.viewport('macbook-13');
    cy.get('[data-testid="analytics-dashboard"]').should('be.visible');
    
    cy.viewport('ipad-2');
    cy.get('[data-testid="analytics-dashboard"]').should('be.visible');
    
    cy.viewport('iphone-x');
    cy.get('[data-testid="analytics-dashboard"]').should('be.visible');
  });
});