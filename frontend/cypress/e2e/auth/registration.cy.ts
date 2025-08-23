// Herald.lol Registration E2E Tests

describe('Registration Flow', () => {
  beforeEach(() => {
    cy.clearLocalStorage();
    
    // Mock API responses
    cy.intercept('POST', '**/auth/register', { fixture: 'auth-response.json' }).as('registerRequest');
    cy.intercept('GET', '**/auth/profile', { fixture: 'user-profile.json' }).as('getProfile');
  });

  it('should display registration form correctly', () => {
    cy.visit('/register');
    
    // Check page elements
    cy.title().should('include', 'Herald.lol');
    cy.get('h1, h4').should('contain.text', 'Join Herald.lol');
    
    // Check form fields
    cy.get('[data-testid="email-input"]').should('be.visible');
    cy.get('[data-testid="username-input"]').should('be.visible');
    cy.get('[data-testid="display-name-input"]').should('be.visible');
    cy.get('[data-testid="password-input"]').should('be.visible');
    cy.get('[data-testid="confirm-password-input"]').should('be.visible');
    cy.get('[data-testid="register-submit"]').should('be.visible');
    
    // Check helper texts
    cy.get('input[name="username"]').parent().should('contain.text', 'Minimum 3 characters');
    cy.get('input[name="display_name"]').parent().should('contain.text', 'How others will see your name');
    cy.get('input[name="password"]').parent().should('contain.text', 'Minimum 6 characters');
    
    // Check link to login
    cy.get('a[href="/login"]').should('contain.text', 'Already have an account? Sign in');
  });

  it('should show validation errors for empty required fields', () => {
    cy.visit('/register');
    
    // Try to submit empty form
    cy.get('[data-testid="register-submit"]').click();
    
    // Required fields should show browser validation
    cy.get('[data-testid="email-input"]:invalid').should('exist');
    cy.get('[data-testid="username-input"]:invalid').should('exist');
    cy.get('[data-testid="password-input"]:invalid').should('exist');
  });

  it('should validate password confirmation', () => {
    cy.visit('/register');
    
    // Fill form with mismatched passwords
    cy.get('[data-testid="email-input"]').type('test@herald.lol');
    cy.get('[data-testid="username-input"]').type('testuser');
    cy.get('[data-testid="password-input"]').type('password123');
    cy.get('[data-testid="confirm-password-input"]').type('differentpassword');
    
    cy.get('[data-testid="register-submit"]').click();
    
    // Should show password mismatch error
    cy.get('[role="alert"]').should('contain.text', 'Passwords do not match');
  });

  it('should validate minimum password length', () => {
    cy.visit('/register');
    
    // Fill form with short password
    cy.get('[data-testid="email-input"]').type('test@herald.lol');
    cy.get('[data-testid="username-input"]').type('testuser');
    cy.get('[data-testid="password-input"]').type('123');
    cy.get('[data-testid="confirm-password-input"]').type('123');
    
    cy.get('[data-testid="register-submit"]').click();
    
    // Should show password length error
    cy.get('[role="alert"]').should('contain.text', 'Password must be at least 6 characters long');
  });

  it('should validate minimum username length', () => {
    cy.visit('/register');
    
    // Fill form with short username
    cy.get('[data-testid="email-input"]').type('test@herald.lol');
    cy.get('[data-testid="username-input"]').type('ab'); // Only 2 characters
    cy.get('[data-testid="password-input"]').type('password123');
    cy.get('[data-testid="confirm-password-input"]').type('password123');
    
    cy.get('[data-testid="register-submit"]').click();
    
    // Should show username length error
    cy.get('[role="alert"]').should('contain.text', 'Username must be at least 3 characters long');
  });

  it('should toggle password visibility', () => {
    cy.visit('/register');
    
    // Test password field
    const passwordInput = cy.get('[data-testid="password-input"]');
    const passwordToggle = passwordInput.parent().find('[aria-label="toggle password visibility"]').first();
    
    passwordInput.should('have.attr', 'type', 'password');
    passwordToggle.click();
    passwordInput.should('have.attr', 'type', 'text');
    
    // Test confirm password field
    const confirmPasswordInput = cy.get('[data-testid="confirm-password-input"]');
    const confirmPasswordToggle = confirmPasswordInput.parent().find('[aria-label="toggle password visibility"]');
    
    confirmPasswordInput.should('have.attr', 'type', 'password');
    confirmPasswordToggle.click();
    confirmPasswordInput.should('have.attr', 'type', 'text');
  });

  it('should register successfully with valid data', () => {
    cy.visit('/register');
    
    // Fill registration form
    cy.fillRegistrationForm('newuser@herald.lol', 'newuser', 'password123', 'New User');
    
    // Submit form
    cy.get('[data-testid="register-submit"]').click();
    
    // Should show loading state
    cy.get('[data-testid="register-submit"]').should('contain.text', 'Creating Account...');
    cy.get('[data-testid="register-submit"]').should('be.disabled');
    
    // Wait for API call
    cy.waitForAPI('@registerRequest');
    
    // Should redirect to dashboard
    cy.url().should('include', '/dashboard');
    cy.shouldBeAuthenticated();
  });

  it('should handle duplicate email error', () => {
    // Mock duplicate email response
    cy.intercept('POST', '**/auth/register', {
      statusCode: 409,
      body: {
        error: 'User already exists',
        message: 'A user with this email or username already exists'
      }
    }).as('duplicateUserRequest');
    
    cy.visit('/register');
    
    cy.fillRegistrationForm('existing@herald.lol', 'existinguser', 'password123');
    cy.get('[data-testid="register-submit"]').click();
    
    cy.waitForAPI('@duplicateUserRequest');
    
    // Should show error message
    cy.get('[role="alert"]').should('contain.text', 'A user with this email or username already exists');
    
    // Should remain on registration page
    cy.url().should('include', '/register');
  });

  it('should handle weak password error from server', () => {
    // Mock weak password response
    cy.intercept('POST', '**/auth/register', {
      statusCode: 400,
      body: {
        error: 'Weak password',
        message: 'Password must be at least 6 characters long'
      }
    }).as('weakPasswordRequest');
    
    cy.visit('/register');
    
    cy.fillRegistrationForm('test@herald.lol', 'testuser', 'weak');
    cy.get('[data-testid="register-submit"]').click();
    
    cy.waitForAPI('@weakPasswordRequest');
    
    // Should show error message
    cy.get('[role="alert"]').should('contain.text', 'Password must be at least 6 characters long');
  });

  it('should clear errors when user starts typing', () => {
    cy.visit('/register');
    
    // Trigger validation error first
    cy.get('[data-testid="email-input"]').type('test@herald.lol');
    cy.get('[data-testid="username-input"]').type('testuser');
    cy.get('[data-testid="password-input"]').type('password123');
    cy.get('[data-testid="confirm-password-input"]').type('differentpassword');
    cy.get('[data-testid="register-submit"]').click();
    
    // Error should appear
    cy.get('[role="alert"]').should('be.visible');
    
    // Error should disappear when user starts typing
    cy.get('[data-testid="email-input"]').type('x');
    cy.get('[role="alert"]').should('not.exist');
  });

  it('should navigate to login page', () => {
    cy.visit('/register');
    
    cy.get('a[href="/login"]').click();
    cy.url().should('include', '/login');
    cy.get('h1, h4').should('contain.text', 'Welcome Back');
  });

  it('should work without display name (optional field)', () => {
    cy.visit('/register');
    
    // Fill form without display name
    cy.get('[data-testid="email-input"]').type('nodisplay@herald.lol');
    cy.get('[data-testid="username-input"]').type('nodisplay');
    cy.get('[data-testid="password-input"]').type('password123');
    cy.get('[data-testid="confirm-password-input"]').type('password123');
    
    cy.get('[data-testid="register-submit"]').click();
    cy.waitForAPI('@registerRequest');
    
    // Should still succeed
    cy.shouldBeAuthenticated();
  });

  it('should handle network errors gracefully', () => {
    // Mock network error
    cy.intercept('POST', '**/auth/register', {
      statusCode: 500,
      body: {
        error: 'Network Error',
        message: 'An error occurred during registration'
      }
    }).as('networkErrorRequest');
    
    cy.visit('/register');
    
    cy.fillRegistrationForm('test@herald.lol', 'testuser', 'password123');
    cy.get('[data-testid="register-submit"]').click();
    
    cy.waitForAPI('@networkErrorRequest');
    
    // Should show error message
    cy.get('[role="alert"]').should('contain.text', 'An error occurred during registration');
  });

  it('should work on mobile viewports', () => {
    cy.viewport('iphone-x');
    cy.visit('/register');
    
    // Check mobile layout
    cy.get('[data-testid="email-input"]').should('be.visible');
    cy.get('[data-testid="username-input"]').should('be.visible');
    cy.get('[data-testid="register-submit"]').should('be.visible');
    
    // Form should still work
    cy.fillRegistrationForm('mobile@herald.lol', 'mobileuser', 'password123');
    cy.get('[data-testid="register-submit"]').click();
    cy.waitForAPI('@registerRequest');
    cy.shouldBeAuthenticated();
  });
});