// Herald.lol System Health E2E Tests

describe('System Health & Configuration', () => {
  it('should validate Cypress configuration', () => {
    // Test that Cypress is properly configured
    expect(Cypress.config('baseUrl')).to.include('localhost');
    expect(Cypress.config('viewportWidth')).to.equal(1280);
    expect(Cypress.config('viewportHeight')).to.equal(720);
  });

  it('should connect to frontend application', () => {
    cy.visit('/');
    cy.get('body').should('be.visible');
    cy.title().should('include', 'Herald.lol');
  });

  it('should validate API mocking works', () => {
    cy.intercept('GET', '**/health', { fixture: 'health-check.json' }).as('healthCheck');
    
    cy.request('GET', '/api/health').then((response) => {
      expect(response.status).to.equal(200);
    });
  });

  it('should validate custom commands work', () => {
    cy.visit('/login');
    
    // Test fillLoginForm command
    cy.fillLoginForm('test@herald.lol', 'password123');
    cy.get('[data-testid="email-input"]').should('have.value', 'test@herald.lol');
    cy.get('[data-testid="password-input"]').should('have.value', 'password123');
  });

  it('should validate authentication state commands', () => {
    // Test shouldNotBeAuthenticated command
    cy.shouldNotBeAuthenticated();
    
    // Mock login and test shouldBeAuthenticated command
    cy.intercept('POST', '**/auth/login', { fixture: 'auth-response.json' }).as('loginRequest');
    cy.fillLoginForm('dev@herald.lol', 'password123');
    cy.get('[data-testid="login-submit"]').click();
    cy.waitForAPI('@loginRequest');
    
    cy.shouldBeAuthenticated();
  });

  it('should validate fixture data loading', () => {
    cy.fixture('auth-response.json').then((authResponse) => {
      expect(authResponse).to.have.property('token');
      expect(authResponse).to.have.property('user');
      expect(authResponse.user).to.have.property('email', 'dev@herald.lol');
    });
    
    cy.fixture('dashboard-stats.json').then((stats) => {
      expect(stats).to.have.property('current_rank');
      expect(stats).to.have.property('total_matches', 150);
      expect(stats.current_rank).to.have.property('tier', 'GOLD');
    });
  });

  it('should validate performance targets', () => {
    const startTime = Date.now();
    
    cy.visit('/');
    
    cy.then(() => {
      const loadTime = Date.now() - startTime;
      // Herald.lol target: <2s UI load time
      expect(loadTime).to.be.lessThan(2000, 'Page should load in under 2 seconds');
    });
  });

  it('should validate error handling setup', () => {
    // Test 404 handling
    cy.visit('/non-existent-route', { failOnStatusCode: false });
    cy.get('[data-testid="error-404"]').should('be.visible');
    
    // Test API error handling
    cy.intercept('GET', '**/api/stats/overview', {
      statusCode: 500,
      body: { error: 'Internal Server Error' }
    }).as('apiError');
    
    cy.visit('/dashboard', { failOnStatusCode: false });
    // Error handling should prevent app crash
    cy.get('body').should('be.visible');
  });

  it('should validate responsive design breakpoints', () => {
    // Desktop
    cy.viewport(1280, 720);
    cy.visit('/');
    cy.get('[data-testid="desktop-nav"]').should('be.visible');
    
    // Tablet
    cy.viewport('ipad-2');
    cy.visit('/');
    cy.get('[data-testid="mobile-nav-toggle"]').should('be.visible');
    
    // Mobile
    cy.viewport('iphone-x');
    cy.visit('/');
    cy.get('[data-testid="mobile-nav-toggle"]').should('be.visible');
  });

  it('should validate accessibility features', () => {
    cy.visit('/');
    
    // Check for proper ARIA labels
    cy.get('[aria-label]').should('exist');
    cy.get('[role]').should('exist');
    
    // Check for keyboard navigation
    cy.get('a, button, input, select').first().focus();
    cy.focused().should('be.visible');
  });

  it('should validate gaming-specific features', () => {
    cy.intercept('POST', '**/auth/login', { fixture: 'auth-response.json' }).as('loginRequest');
    cy.visit('/login');
    cy.fillLoginForm('dev@herald.lol', 'password123');
    cy.get('[data-testid="login-submit"]').click();
    cy.waitForAPI('@loginRequest');
    
    // Check gaming-specific UI elements
    cy.get('[data-testid="current-rank"]').should('contain.text', 'Gold II');
    cy.get('[data-testid="favorite-champion"]').should('contain.text', 'Jinx');
    cy.get('[data-testid="main-role"]').should('contain.text', 'ADC');
  });
});