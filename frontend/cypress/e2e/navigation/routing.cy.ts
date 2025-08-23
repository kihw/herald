// Herald.lol Navigation & Routing E2E Tests

describe('Navigation & Routing', () => {
  beforeEach(() => {
    cy.clearLocalStorage();
  });

  describe('Public Routes', () => {
    it('should navigate to home page', () => {
      cy.visit('/');
      cy.title().should('include', 'Herald.lol');
      cy.get('[data-testid="hero-section"]').should('be.visible');
      cy.get('[data-testid="cta-register"]').should('contain.text', 'Get Started');
    });

    it('should navigate to login page', () => {
      cy.visit('/login');
      cy.get('h1, h4').should('contain.text', 'Welcome Back');
      cy.get('[data-testid="email-input"]').should('be.visible');
    });

    it('should navigate to register page', () => {
      cy.visit('/register');
      cy.get('h1, h4').should('contain.text', 'Join Herald.lol');
      cy.get('[data-testid="username-input"]').should('be.visible');
    });

    it('should navigate between login and register', () => {
      cy.visit('/login');
      cy.get('a[href="/register"]').click();
      cy.url().should('include', '/register');
      
      cy.get('a[href="/login"]').click();
      cy.url().should('include', '/login');
    });
  });

  describe('Protected Routes', () => {
    beforeEach(() => {
      // Setup authenticated user
      cy.intercept('POST', '**/auth/login', { fixture: 'auth-response.json' }).as('loginRequest');
      cy.intercept('GET', '**/auth/profile', { fixture: 'user-profile.json' }).as('getProfile');
      
      cy.visit('/login');
      cy.fillLoginForm('dev@herald.lol', 'password123');
      cy.get('[data-testid="login-submit"]').click();
      cy.waitForAPI('@loginRequest');
    });

    it('should navigate to dashboard after login', () => {
      cy.url().should('include', '/dashboard');
      cy.get('[data-testid="dashboard-header"]').should('be.visible');
    });

    it('should navigate between protected routes', () => {
      // Dashboard to Analytics
      cy.get('[data-testid="nav-analytics"]').click();
      cy.url().should('include', '/analytics');
      cy.get('[data-testid="analytics-header"]').should('be.visible');
      
      // Analytics to Matches
      cy.get('[data-testid="nav-matches"]').click();
      cy.url().should('include', '/matches');
      cy.get('[data-testid="matches-header"]').should('be.visible');
      
      // Back to Dashboard
      cy.get('[data-testid="nav-dashboard"]').click();
      cy.url().should('include', '/dashboard');
    });

    it('should maintain active navigation state', () => {
      cy.get('[data-testid="nav-dashboard"]').should('have.class', 'active');
      
      cy.get('[data-testid="nav-analytics"]').click();
      cy.get('[data-testid="nav-analytics"]').should('have.class', 'active');
      cy.get('[data-testid="nav-dashboard"]').should('not.have.class', 'active');
    });
  });

  describe('Route Protection', () => {
    it('should redirect to login when accessing protected route without auth', () => {
      cy.visit('/dashboard');
      cy.url().should('include', '/login');
      cy.get('[data-testid="login-required-message"]').should('contain.text', 'Please sign in to continue');
    });

    it('should redirect to intended page after login', () => {
      // Try to access analytics page
      cy.visit('/analytics');
      cy.url().should('include', '/login');
      
      // Login
      cy.intercept('POST', '**/auth/login', { fixture: 'auth-response.json' }).as('loginRequest');
      cy.fillLoginForm('dev@herald.lol', 'password123');
      cy.get('[data-testid="login-submit"]').click();
      cy.waitForAPI('@loginRequest');
      
      // Should redirect to originally requested page
      cy.url().should('include', '/analytics');
    });

    it('should prevent access to auth pages when already logged in', () => {
      // Login first
      cy.intercept('POST', '**/auth/login', { fixture: 'auth-response.json' }).as('loginRequest');
      cy.visit('/login');
      cy.fillLoginForm('dev@herald.lol', 'password123');
      cy.get('[data-testid="login-submit"]').click();
      cy.waitForAPI('@loginRequest');
      
      // Try to visit login page again
      cy.visit('/login');
      cy.url().should('include', '/dashboard');
      
      // Try to visit register page
      cy.visit('/register');
      cy.url().should('include', '/dashboard');
    });
  });

  describe('404 Error Handling', () => {
    it('should show 404 page for non-existent routes', () => {
      cy.visit('/non-existent-page', { failOnStatusCode: false });
      cy.get('[data-testid="error-404"]').should('be.visible');
      cy.get('[data-testid="error-404"]').should('contain.text', '404');
      cy.get('[data-testid="back-to-home"]').should('be.visible');
    });

    it('should navigate back from 404 page', () => {
      cy.visit('/non-existent-page', { failOnStatusCode: false });
      cy.get('[data-testid="back-to-home"]').click();
      cy.url().should('eq', Cypress.config().baseUrl + '/');
    });
  });

  describe('Browser Navigation', () => {
    beforeEach(() => {
      cy.intercept('POST', '**/auth/login', { fixture: 'auth-response.json' }).as('loginRequest');
      cy.visit('/login');
      cy.fillLoginForm('dev@herald.lol', 'password123');
      cy.get('[data-testid="login-submit"]').click();
      cy.waitForAPI('@loginRequest');
    });

    it('should support browser back/forward buttons', () => {
      // Navigate to analytics
      cy.get('[data-testid="nav-analytics"]').click();
      cy.url().should('include', '/analytics');
      
      // Use browser back button
      cy.go('back');
      cy.url().should('include', '/dashboard');
      
      // Use browser forward button  
      cy.go('forward');
      cy.url().should('include', '/analytics');
    });

    it('should maintain page state on refresh', () => {
      cy.get('[data-testid="nav-analytics"]').click();
      cy.url().should('include', '/analytics');
      
      cy.reload();
      cy.url().should('include', '/analytics');
      cy.get('[data-testid="analytics-header"]').should('be.visible');
    });
  });

  describe('Mobile Navigation', () => {
    beforeEach(() => {
      cy.viewport('iphone-x');
      cy.intercept('POST', '**/auth/login', { fixture: 'auth-response.json' }).as('loginRequest');
      cy.visit('/login');
      cy.fillLoginForm('dev@herald.lol', 'password123');
      cy.get('[data-testid="login-submit"]').click();
      cy.waitForAPI('@loginRequest');
    });

    it('should show mobile navigation menu', () => {
      cy.get('[data-testid="mobile-nav-toggle"]').should('be.visible');
      cy.get('[data-testid="mobile-nav-toggle"]').click();
      
      cy.get('[data-testid="mobile-nav-menu"]').should('be.visible');
      cy.get('[data-testid="nav-analytics"]').should('be.visible');
    });

    it('should close mobile menu after navigation', () => {
      cy.get('[data-testid="mobile-nav-toggle"]').click();
      cy.get('[data-testid="nav-analytics"]').click();
      
      cy.url().should('include', '/analytics');
      cy.get('[data-testid="mobile-nav-menu"]').should('not.be.visible');
    });
  });
});