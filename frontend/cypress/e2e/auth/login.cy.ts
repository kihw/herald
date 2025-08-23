// Herald.lol Authentication E2E Tests

describe('Login Flow', () => {
  beforeEach(() => {
    // Reset state before each test
    cy.clearLocalStorage();
    
    // Mock API responses
    cy.intercept('POST', '**/auth/login', { fixture: 'auth-response.json' }).as('loginRequest');
    cy.intercept('GET', '**/auth/profile', { fixture: 'user-profile.json' }).as('getProfile');
  });

  it('should display login form correctly', () => {
    cy.visit('/login');
    
    // Check page title and heading
    cy.title().should('include', 'Herald.lol');
    cy.get('h1, h4').should('contain.text', 'Welcome Back');
    
    // Check form elements exist
    cy.get('[data-testid="email-input"]').should('be.visible');
    cy.get('[data-testid="password-input"]').should('be.visible');
    cy.get('[data-testid="login-submit"]').should('be.visible');
    cy.get('[data-testid="login-submit"]').should('contain.text', 'Sign In');
    
    // Check links
    cy.get('a[href="/register"]').should('contain.text', 'Sign up');
    cy.get('a').should('contain.text', 'Forgot password');
  });

  it('should show validation errors for empty fields', () => {
    cy.visit('/login');
    
    // Try to submit empty form
    cy.get('[data-testid="login-submit"]').click();
    
    // Browser validation should prevent submission
    cy.get('[data-testid="email-input"]:invalid').should('exist');
  });

  it('should toggle password visibility', () => {
    cy.visit('/login');
    
    const passwordInput = cy.get('[data-testid="password-input"]');
    const toggleButton = cy.get('[aria-label="toggle password visibility"]');
    
    // Initially password should be hidden
    passwordInput.should('have.attr', 'type', 'password');
    
    // Click toggle to show password
    toggleButton.click();
    passwordInput.should('have.attr', 'type', 'text');
    
    // Click toggle to hide password again
    toggleButton.click();
    passwordInput.should('have.attr', 'type', 'password');
  });

  it('should login successfully with valid credentials', () => {
    cy.visit('/login');
    
    // Fill in form
    cy.fillLoginForm('dev@herald.lol', 'password123');
    
    // Submit form
    cy.get('[data-testid="login-submit"]').click();
    
    // Should show loading state
    cy.get('[data-testid="login-submit"]').should('contain.text', 'Signing In...');
    cy.get('[data-testid="login-submit"]').should('be.disabled');
    
    // Wait for API call
    cy.waitForAPI('@loginRequest');
    
    // Should redirect to dashboard
    cy.url().should('include', '/dashboard');
    cy.shouldBeAuthenticated();
  });

  it('should show error message for invalid credentials', () => {
    // Mock failed login response
    cy.intercept('POST', '**/auth/login', {
      statusCode: 401,
      body: {
        error: 'Invalid credentials',
        message: 'Email or password is incorrect'
      }
    }).as('failedLoginRequest');
    
    cy.visit('/login');
    
    // Fill in form with invalid credentials
    cy.fillLoginForm('invalid@test.com', 'wrongpassword');
    cy.get('[data-testid="login-submit"]').click();
    
    // Wait for failed API call
    cy.waitForAPI('@failedLoginRequest');
    
    // Should show error message
    cy.get('[role="alert"]').should('contain.text', 'Email or password is incorrect');
    
    // Should remain on login page
    cy.url().should('include', '/login');
    cy.shouldNotBeAuthenticated();
  });

  it('should clear error message when user starts typing', () => {
    // Mock failed login response
    cy.intercept('POST', '**/auth/login', {
      statusCode: 401,
      body: {
        error: 'Invalid credentials', 
        message: 'Email or password is incorrect'
      }
    }).as('failedLoginRequest');
    
    cy.visit('/login');
    
    // Trigger error first
    cy.fillLoginForm('invalid@test.com', 'wrongpassword');
    cy.get('[data-testid="login-submit"]').click();
    cy.waitForAPI('@failedLoginRequest');
    cy.get('[role="alert"]').should('be.visible');
    
    // Error should disappear when user starts typing
    cy.get('[data-testid="email-input"]').clear().type('n');
    cy.get('[role="alert"]').should('not.exist');
  });

  it('should navigate to register page', () => {
    cy.visit('/login');
    
    cy.get('a[href="/register"]').click();
    cy.url().should('include', '/register');
    cy.get('h1, h4').should('contain.text', 'Join Herald.lol');
  });

  it('should handle network errors gracefully', () => {
    // Mock network error
    cy.intercept('POST', '**/auth/login', {
      statusCode: 500,
      body: {
        error: 'Network Error',
        message: 'An unexpected error occurred'
      }
    }).as('networkErrorRequest');
    
    cy.visit('/login');
    
    cy.fillLoginForm('dev@herald.lol', 'password123');
    cy.get('[data-testid="login-submit"]').click();
    
    cy.waitForAPI('@networkErrorRequest');
    
    // Should show generic error message
    cy.get('[role="alert"]').should('contain.text', 'An unexpected error occurred');
  });

  it('should redirect to intended page after login', () => {
    // Try to access protected route first
    cy.visit('/dashboard');
    
    // Should redirect to login with return URL
    cy.url().should('include', '/login');
    
    // Login
    cy.fillLoginForm('dev@herald.lol', 'password123');
    cy.get('[data-testid="login-submit"]').click();
    cy.waitForAPI('@loginRequest');
    
    // Should redirect back to originally requested page
    cy.url().should('include', '/dashboard');
  });

  it('should work on mobile viewports', () => {
    cy.viewport('iphone-x');
    cy.visit('/login');
    
    // Check mobile layout
    cy.get('[data-testid="email-input"]').should('be.visible');
    cy.get('[data-testid="password-input"]').should('be.visible');
    cy.get('[data-testid="login-submit"]').should('be.visible');
    
    // Form should still work
    cy.fillLoginForm('dev@herald.lol', 'password123');
    cy.get('[data-testid="login-submit"]').click();
    cy.waitForAPI('@loginRequest');
    cy.shouldBeAuthenticated();
  });
});