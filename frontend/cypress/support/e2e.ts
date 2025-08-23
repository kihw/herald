// Cypress E2E Support File for Herald.lol

import './commands';

// Import Cypress commands
/// <reference types="cypress" />

// Configure Cypress behavior
Cypress.on('uncaught:exception', (err, runnable) => {
  // Ignore specific errors that don't affect functionality
  if (err.message.includes('ResizeObserver loop limit exceeded')) {
    return false;
  }
  
  if (err.message.includes('Non-Error promise rejection captured')) {
    return false;
  }
  
  // Don't fail tests on gaming analytics errors that might be expected
  if (err.message.includes('Analytics not loaded') || 
      err.message.includes('Riot API rate limit') ||
      err.message.includes('Match not found')) {
    return false;
  }
  
  // Log gaming-specific errors
  if (runnable) {
    cy.task('logGamingTestResult', {
      error: err.message,
      test: runnable.title,
      type: 'uncaught_exception'
    });
  }
  
  // Let other errors fail the test
  return true;
});

// Global test setup
beforeEach(() => {
  // Clear localStorage and sessionStorage before each test
  cy.clearLocalStorage();
  cy.clearCookies();
  
  // Set viewport for consistent testing
  cy.viewport(1280, 720);
  
  // Intercept and stub common API calls to avoid external dependencies
  cy.intercept('GET', '**/health', { fixture: 'health-check.json' });
  
  // Stub rate limit status to avoid Riot API calls
  cy.intercept('GET', '**/riot/rate-limit', { fixture: 'rate-limit-status.json' });
  
  // Set up common interceptors for gaming analytics
  if (Cypress.env('VALIDATE_GAMING_METRICS')) {
    cy.mockGamingAnalytics();
  }
  
  // Performance monitoring setup
  if (Cypress.env('VALIDATE_PERFORMANCE')) {
    cy.window().then((win) => {
      // Mark test start time for performance validation
      win.testStartTime = Date.now();
    });
  }
});

// Global after hook
afterEach(() => {
  // Performance monitoring cleanup
  if (Cypress.env('VALIDATE_PERFORMANCE')) {
    cy.window().then((win) => {
      if (win.testStartTime) {
        const duration = Date.now() - win.testStartTime;
        cy.task('logGamingTestResult', {
          test: Cypress.currentTest.title,
          duration,
          status: Cypress.currentTest.state
        });
      }
    });
  }
  
  // Clean up after each test
  cy.clearLocalStorage();
  cy.clearCookies();
});

// Gaming test configuration
Cypress.config('defaultCommandTimeout', Cypress.env('ANALYTICS_LOAD_TIMEOUT') || 5000);
Cypress.config('requestTimeout', 15000);
Cypress.config('responseTimeout', 15000);

// Gaming performance assertions
Cypress.Commands.add('assertGamingPerformance', (type: string, startTime: number) => {
  const endTime = Date.now();
  const duration = endTime - startTime;
  
  cy.task('validateGamingPerformance', { startTime, endTime, type })
    .then((result) => {
      expect(result.passed).to.be.true;
      cy.log(`✅ Gaming Performance: ${type} took ${result.duration}ms (threshold: ${result.threshold}ms)`);
    });
});

// Gaming metrics validation command
Cypress.Commands.add('assertGamingMetric', (type: string, value: number) => {
  cy.task('validateGamingMetrics', { type, value })
    .then((result) => {
      expect(result.valid).to.be.true;
      cy.log(`✅ Gaming Metric: ${result.type} = ${result.value}`);
    });
});

declare global {
  namespace Cypress {
    interface Chainable {
      assertGamingPerformance(type: string, startTime: number): Chainable<Element>;
      assertGamingMetric(type: string, value: number): Chainable<Element>;
    }
  }
}

// Console logging for debugging gaming tests
cy.task('logGamingTestResult', {
  message: 'Herald.lol E2E Testing Environment Initialized',
  timestamp: new Date().toISOString(),
  config: {
    analyticsTimeout: Cypress.env('ANALYTICS_LOAD_TIMEOUT'),
    uiTimeout: Cypress.env('UI_LOAD_TIMEOUT'),
    validateMetrics: Cypress.env('VALIDATE_GAMING_METRICS'),
    validatePerformance: Cypress.env('VALIDATE_PERFORMANCE'),
  }
});

// Silence specific warnings
const resizeObserverLoopErrRe = /^[^(ResizeObserver loop limit exceeded)]/;
Cypress.on('uncaught:exception', (err) => {
  if (resizeObserverLoopErrRe.test(err.message)) {
    return false;
  }
});