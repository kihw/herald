// Herald.lol Gaming Analytics Performance E2E Tests

describe('Gaming Analytics Performance', () => {
  beforeEach(() => {
    // Mock user authentication
    cy.mockAuthentication();
    
    // Mock Riot account
    cy.mockRiotAccount({
      summonerName: 'TestSummoner',
      region: 'na1',
      isVerified: true
    });

    // Mock analytics API responses
    cy.intercept('GET', '**/analytics/kda/*', { fixture: 'analytics/kda-analysis.json' }).as('getKDA');
    cy.intercept('GET', '**/analytics/cs/*', { fixture: 'analytics/cs-analysis.json' }).as('getCS');
    cy.intercept('GET', '**/analytics/vision/*', { fixture: 'analytics/vision-analysis.json' }).as('getVision');
    cy.intercept('GET', '**/analytics/damage/*', { fixture: 'analytics/damage-analysis.json' }).as('getDamage');
    cy.intercept('GET', '**/analytics/gold/*', { fixture: 'analytics/gold-analysis.json' }).as('getGold');
  });

  it('should complete analytics dashboard load within 5 seconds', () => {
    const startTime = Date.now();
    
    cy.visit('/dashboard');
    
    // Wait for all analytics components to load
    cy.waitForAPI('@getKDA');
    cy.waitForAPI('@getCS');
    cy.waitForAPI('@getVision');
    cy.waitForAPI('@getDamage');
    cy.waitForAPI('@getGold');
    
    // Verify all analytics components are visible
    cy.get('[data-testid="kda-analytics"]').should('be.visible');
    cy.get('[data-testid="cs-analytics"]').should('be.visible');
    cy.get('[data-testid="vision-analytics"]').should('be.visible');
    cy.get('[data-testid="damage-analytics"]').should('be.visible');
    cy.get('[data-testid="gold-analytics"]').should('be.visible');

    // Verify performance requirement
    cy.then(() => {
      const loadTime = Date.now() - startTime;
      expect(loadTime).to.be.lessThan(5000, 'Analytics dashboard must load within 5 seconds');
      cy.log(`Dashboard loaded in ${loadTime}ms (target: <5000ms)`);
    });
  });

  it('should display gaming metrics accurately', () => {
    cy.visit('/dashboard');
    
    // Wait for KDA analytics
    cy.waitForAPI('@getKDA');
    
    // Verify KDA display
    cy.get('[data-testid="kda-value"]').should('contain.text', '2.34');
    cy.get('[data-testid="kda-trend"]').should('contain.text', 'improving');
    cy.get('[data-testid="kda-percentile"]').should('contain.text', '65');
    
    // Verify gaming calculations are reasonable
    cy.get('[data-testid="kda-value"]').should(($el) => {
      const kda = parseFloat($el.text());
      expect(kda).to.be.greaterThan(0);
      expect(kda).to.be.lessThan(20); // Reasonable upper bound
    });
    
    // Wait for CS analytics
    cy.waitForAPI('@getCS');
    
    // Verify CS/min display
    cy.get('[data-testid="cs-per-min"]').should('contain.text', '7.2');
    cy.get('[data-testid="cs-efficiency"]').should('contain.text', '85%');
    
    // Verify CS calculations
    cy.get('[data-testid="cs-per-min"]').should(($el) => {
      const csPerMin = parseFloat($el.text());
      expect(csPerMin).to.be.greaterThan(0);
      expect(csPerMin).to.be.lessThan(15); // Reasonable upper bound
    });
    
    // Wait for Vision analytics
    cy.waitForAPI('@getVision');
    
    // Verify Vision Score
    cy.get('[data-testid="vision-score"]').should('contain.text', '28.4');
    cy.get('[data-testid="vision-control"]').should('contain.text', '72%');
    
    // Verify vision calculations
    cy.get('[data-testid="vision-score"]').should(($el) => {
      const visionScore = parseFloat($el.text());
      expect(visionScore).to.be.greaterThan(0);
      expect(visionScore).to.be.lessThan(200); // Reasonable upper bound
    });
  });

  it('should handle large match datasets efficiently', () => {
    // Mock large dataset response
    cy.intercept('GET', '**/analytics/kda/*', { fixture: 'analytics/large-kda-dataset.json' }).as('getLargeKDA');
    
    const startTime = Date.now();
    cy.visit('/dashboard');
    
    cy.waitForAPI('@getLargeKDA');
    
    // Should still render within performance target
    cy.get('[data-testid="kda-analytics"]').should('be.visible');
    cy.get('[data-testid="match-count"]').should('contain.text', '500+');
    
    cy.then(() => {
      const processTime = Date.now() - startTime;
      expect(processTime).to.be.lessThan(5000, 'Large dataset should process within 5 seconds');
    });
  });

  it('should display confidence levels for analytics', () => {
    cy.visit('/dashboard');
    
    // Wait for all analytics
    cy.waitForAPI('@getKDA');
    cy.waitForAPI('@getCS');
    cy.waitForAPI('@getVision');
    
    // Verify confidence indicators
    cy.get('[data-testid="kda-confidence"]').should('contain.text', '92%');
    cy.get('[data-testid="cs-confidence"]').should('contain.text', '88%');
    cy.get('[data-testid="vision-confidence"]').should('contain.text', '79%');
    
    // High confidence should show green indicator
    cy.get('[data-testid="kda-confidence"]').should('have.class', 'high-confidence');
    
    // Medium confidence should show yellow indicator
    cy.get('[data-testid="vision-confidence"]').should('have.class', 'medium-confidence');
  });

  it('should warn when confidence is low', () => {
    // Mock low confidence response
    cy.intercept('GET', '**/analytics/kda/*', { 
      body: {
        data: { currentKDA: 1.5, trend: 'insufficient_data' },
        confidence: 0.25
      }
    }).as('getLowConfidenceKDA');
    
    cy.visit('/dashboard');
    cy.waitForAPI('@getLowConfidenceKDA');
    
    // Should show warning
    cy.get('[data-testid="confidence-warning"]').should('be.visible');
    cy.get('[data-testid="confidence-warning"]').should('contain.text', 'Limited data available');
    
    // Should suggest more matches
    cy.get('[data-testid="data-suggestion"]').should('contain.text', 'Play more matches');
  });

  it('should update analytics in real-time', () => {
    cy.visit('/dashboard');
    
    // Initial data load
    cy.waitForAPI('@getKDA');
    cy.get('[data-testid="kda-value"]').should('contain.text', '2.34');
    
    // Mock updated data
    cy.intercept('GET', '**/analytics/kda/*', {
      body: {
        data: { currentKDA: 2.67, trend: 'improving' },
        confidence: 0.94
      }
    }).as('getUpdatedKDA');
    
    // Trigger refresh
    cy.get('[data-testid="refresh-analytics"]').click();
    cy.waitForAPI('@getUpdatedKDA');
    
    // Should show updated value
    cy.get('[data-testid="kda-value"]').should('contain.text', '2.67');
  });

  it('should handle analytics errors gracefully', () => {
    // Mock API error
    cy.intercept('GET', '**/analytics/kda/*', {
      statusCode: 503,
      body: { error: 'Analytics service unavailable' }
    }).as('getKDAError');
    
    cy.visit('/dashboard');
    cy.waitForAPI('@getKDAError');
    
    // Should show error state
    cy.get('[data-testid="kda-error"]').should('be.visible');
    cy.get('[data-testid="kda-error"]').should('contain.text', 'Unable to load KDA analytics');
    
    // Should show retry button
    cy.get('[data-testid="retry-kda"]').should('be.visible');
  });

  it('should validate gaming metric ranges', () => {
    cy.visit('/dashboard');
    
    cy.waitForAPI('@getKDA');
    cy.waitForAPI('@getCS');
    cy.waitForAPI('@getDamage');
    
    // KDA should be reasonable
    cy.get('[data-testid="kda-value"]').should(($el) => {
      const kda = parseFloat($el.text());
      expect(kda).to.be.at.least(0);
      expect(kda).to.be.at.most(50); // Even perfect games rarely exceed 50 KDA
    });
    
    // CS/min should be reasonable
    cy.get('[data-testid="cs-per-min"]').should(($el) => {
      const cs = parseFloat($el.text());
      expect(cs).to.be.at.least(0);
      expect(cs).to.be.at.most(15); // 15+ CS/min is unrealistic
    });
    
    // Damage share should be percentage
    cy.get('[data-testid="damage-share"]').should(($el) => {
      const text = $el.text();
      const percentage = parseFloat(text.replace('%', ''));
      expect(percentage).to.be.at.least(0);
      expect(percentage).to.be.at.most(100);
    });
  });

  it('should work across different screen sizes', () => {
    // Test desktop
    cy.viewport(1920, 1080);
    cy.visit('/dashboard');
    cy.waitForAPI('@getKDA');
    cy.get('[data-testid="analytics-grid"]').should('have.class', 'desktop-layout');
    
    // Test tablet
    cy.viewport(768, 1024);
    cy.visit('/dashboard');
    cy.waitForAPI('@getKDA');
    cy.get('[data-testid="analytics-grid"]').should('have.class', 'tablet-layout');
    
    // Test mobile
    cy.viewport(375, 667);
    cy.visit('/dashboard');
    cy.waitForAPI('@getKDA');
    cy.get('[data-testid="analytics-grid"]').should('have.class', 'mobile-layout');
  });

  it('should support keyboard navigation', () => {
    cy.visit('/dashboard');
    cy.waitForAPI('@getKDA');
    
    // Tab through analytics cards
    cy.get('body').tab();
    cy.focused().should('have.attr', 'data-testid', 'kda-analytics');
    
    cy.focused().tab();
    cy.focused().should('have.attr', 'data-testid', 'cs-analytics');
    
    cy.focused().tab();
    cy.focused().should('have.attr', 'data-testid', 'vision-analytics');
    
    // Should be able to interact with focused elements
    cy.focused().type('{enter}');
    cy.get('[data-testid="analytics-details"]').should('be.visible');
  });

  it('should maintain performance with concurrent users simulation', () => {
    // Simulate network delay for high load
    cy.intercept('GET', '**/analytics/**', (req) => {
      req.reply((res) => {
        // Add 500ms delay to simulate high load
        return new Promise(resolve => {
          setTimeout(() => resolve(res.send()), 500);
        });
      });
    });
    
    const startTime = Date.now();
    cy.visit('/dashboard');
    
    // Even with simulated load, should meet performance targets
    cy.get('[data-testid="analytics-container"]').should('be.visible');
    
    cy.then(() => {
      const loadTime = Date.now() - startTime;
      // Should handle simulated load gracefully
      expect(loadTime).to.be.lessThan(10000, 'Should handle high load within 10 seconds');
    });
  });
});