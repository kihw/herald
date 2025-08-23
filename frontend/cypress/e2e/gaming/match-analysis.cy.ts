// Herald.lol Match Analysis E2E Tests

describe('Gaming Match Analysis', () => {
  beforeEach(() => {
    cy.mockAuthentication();
    cy.mockRiotAccount();

    // Mock match data
    cy.intercept('GET', '**/matches', { fixture: 'matches/recent-matches.json' }).as('getMatches');
    cy.intercept('GET', '**/matches/*/analysis', { fixture: 'matches/match-analysis.json' }).as('getMatchAnalysis');
    cy.intercept('POST', '**/analytics/predict/match', { fixture: 'analytics/match-prediction.json' }).as('predictMatch');
  });

  it('should complete match analysis within 5 seconds', () => {
    cy.visit('/matches');
    cy.waitForAPI('@getMatches');
    
    const startTime = Date.now();
    
    // Click on a match to analyze
    cy.get('[data-testid="match-card"]').first().click();
    cy.waitForAPI('@getMatchAnalysis');
    
    // Verify analysis components are loaded
    cy.get('[data-testid="match-overview"]').should('be.visible');
    cy.get('[data-testid="performance-metrics"]').should('be.visible');
    cy.get('[data-testid="timeline-chart"]').should('be.visible');
    cy.get('[data-testid="recommendations"]').should('be.visible');
    
    cy.then(() => {
      const analysisTime = Date.now() - startTime;
      expect(analysisTime).to.be.lessThan(5000, 'Match analysis must complete within 5 seconds');
      cy.log(`Match analysis completed in ${analysisTime}ms`);
    });
  });

  it('should display accurate gaming statistics', () => {
    cy.visit('/matches');
    cy.waitForAPI('@getMatches');
    
    cy.get('[data-testid="match-card"]').first().click();
    cy.waitForAPI('@getMatchAnalysis');
    
    // Verify match overview
    cy.get('[data-testid="champion-name"]').should('contain.text', 'Jinx');
    cy.get('[data-testid="game-result"]').should('contain.text', 'Victory');
    cy.get('[data-testid="game-duration"]').should('contain.text', '30:34');
    cy.get('[data-testid="queue-type"]').should('contain.text', 'Ranked Solo/Duo');
    
    // Verify KDA display
    cy.get('[data-testid="kills"]').should('contain.text', '12');
    cy.get('[data-testid="deaths"]').should('contain.text', '3');
    cy.get('[data-testid="assists"]').should('contain.text', '18');
    cy.get('[data-testid="kda-ratio"]').should('contain.text', '10.0');
    
    // Verify CS and vision
    cy.get('[data-testid="creep-score"]').should('contain.text', '245');
    cy.get('[data-testid="cs-per-min"]').should('contain.text', '8.0');
    cy.get('[data-testid="vision-score"]').should('contain.text', '32');
    
    // Verify damage and gold
    cy.get('[data-testid="total-damage"]').should('contain.text', '45,678');
    cy.get('[data-testid="damage-share"]').should('contain.text', '32%');
    cy.get('[data-testid="gold-earned"]').should('contain.text', '16,234');
  });

  it('should validate gaming metric calculations', () => {
    cy.visit('/matches');
    cy.waitForAPI('@getMatches');
    
    cy.get('[data-testid="match-card"]').first().click();
    cy.waitForAPI('@getMatchAnalysis');
    
    // Validate KDA calculation: (12 + 18) / 3 = 10.0
    cy.get('[data-testid="kda-ratio"]').should(($el) => {
      const kda = parseFloat($el.text());
      expect(kda).to.equal(10.0);
    });
    
    // Validate CS/min calculation: 245 CS in 30.57 minutes ≈ 8.0
    cy.get('[data-testid="cs-per-min"]').should(($el) => {
      const csPerMin = parseFloat($el.text());
      expect(csPerMin).to.be.closeTo(8.0, 0.1);
    });
    
    // Validate damage share percentage
    cy.get('[data-testid="damage-share"]').should(($el) => {
      const percentage = parseFloat($el.text().replace('%', ''));
      expect(percentage).to.be.greaterThan(0);
      expect(percentage).to.be.lessThan(100);
      expect(percentage).to.equal(32);
    });
  });

  it('should display performance compared to rank average', () => {
    cy.visit('/matches');
    cy.waitForAPI('@getMatches');
    
    cy.get('[data-testid="match-card"]').first().click();
    cy.waitForAPI('@getMatchAnalysis');
    
    // Should show rank comparison
    cy.get('[data-testid="rank-comparison"]').should('be.visible');
    cy.get('[data-testid="kda-vs-rank"]').should('contain.text', 'Above average');
    cy.get('[data-testid="cs-vs-rank"]').should('contain.text', 'Above average');
    cy.get('[data-testid="vision-vs-rank"]').should('contain.text', 'Average');
    
    // Should show percentiles
    cy.get('[data-testid="kda-percentile"]').should('contain.text', '78th');
    cy.get('[data-testid="cs-percentile"]').should('contain.text', '65th');
    cy.get('[data-testid="vision-percentile"]').should('contain.text', '52nd');
  });

  it('should show actionable improvement recommendations', () => {
    cy.visit('/matches');
    cy.waitForAPI('@getMatches');
    
    cy.get('[data-testid="match-card"]').first().click();
    cy.waitForAPI('@getMatchAnalysis');
    
    // Should have recommendations section
    cy.get('[data-testid="recommendations"]').should('be.visible');
    cy.get('[data-testid="improvement-tips"]').should('exist');
    
    // Should have specific gaming recommendations
    cy.get('[data-testid="recommendation-item"]').should('have.length.greaterThan', 0);
    
    // Should categorize recommendations
    cy.get('[data-testid="mechanical-tips"]').should('exist');
    cy.get('[data-testid="tactical-tips"]').should('exist');
    cy.get('[data-testid="strategic-tips"]').should('exist');
    
    // Recommendations should be actionable
    cy.get('[data-testid="recommendation-item"]').first().should('contain', 'positioning');
  });

  it('should display match timeline correctly', () => {
    cy.visit('/matches');
    cy.waitForAPI('@getMatches');
    
    cy.get('[data-testid="match-card"]').first().click();
    cy.waitForAPI('@getMatchAnalysis');
    
    // Timeline chart should be visible
    cy.get('[data-testid="timeline-chart"]').should('be.visible');
    
    // Should show key events
    cy.get('[data-testid="timeline-events"]').should('exist');
    cy.get('[data-testid="first-blood"]').should('exist');
    cy.get('[data-testid="tower-kills"]').should('exist');
    cy.get('[data-testid="dragon-kills"]').should('exist');
    cy.get('[data-testid="baron-kills"]').should('exist');
    
    // Should show gold advantage over time
    cy.get('[data-testid="gold-chart"]').should('be.visible');
    
    // Should be interactive
    cy.get('[data-testid="timeline-chart"]').trigger('mouseover');
    cy.get('[data-testid="timeline-tooltip"]').should('be.visible');
  });

  it('should handle different champion roles correctly', () => {
    // Mock different role data
    cy.intercept('GET', '**/matches/*/analysis', { fixture: 'matches/support-analysis.json' }).as('getSupportAnalysis');
    
    cy.visit('/matches');
    cy.waitForAPI('@getMatches');
    
    cy.get('[data-testid="match-card"]').eq(1).click(); // Support match
    cy.waitForAPI('@getSupportAnalysis');
    
    // Should show role-specific metrics
    cy.get('[data-testid="role"]').should('contain.text', 'Support');
    cy.get('[data-testid="champion-name"]').should('contain.text', 'Thresh');
    
    // Support-specific metrics should be emphasized
    cy.get('[data-testid="vision-score"]').should('have.class', 'emphasized');
    cy.get('[data-testid="assists"]').should('have.class', 'emphasized');
    
    // Should show support-specific recommendations
    cy.get('[data-testid="recommendations"]').should('contain', 'ward placement');
    cy.get('[data-testid="recommendations"]').should('contain', 'vision control');
  });

  it('should predict match outcome accurately', () => {
    cy.visit('/matches');
    cy.waitForAPI('@getMatches');
    
    // Click on live/upcoming match for prediction
    cy.get('[data-testid="predict-match"]').first().click();
    cy.waitForAPI('@predictMatch');
    
    // Should show prediction results
    cy.get('[data-testid="win-probability"]').should('be.visible');
    cy.get('[data-testid="win-probability"]').should('contain.text', '67%');
    
    // Should show key factors
    cy.get('[data-testid="prediction-factors"]').should('exist');
    cy.get('[data-testid="factor-item"]').should('contain', 'champion composition');
    cy.get('[data-testid="factor-item"]').should('contain', 'recent form');
    cy.get('[data-testid="factor-item"]').should('contain', 'skill matchup');
    
    // Should show confidence level
    cy.get('[data-testid="prediction-confidence"]').should('contain.text', '73%');
  });

  it('should handle perfect and poor game scenarios', () => {
    // Test perfect game
    cy.intercept('GET', '**/matches/*/analysis', { fixture: 'matches/perfect-game.json' }).as('getPerfectGame');
    
    cy.visit('/matches');
    cy.waitForAPI('@getMatches');
    
    cy.get('[data-testid="match-card"]').first().click();
    cy.waitForAPI('@getPerfectGame');
    
    // Should show perfect game indicators
    cy.get('[data-testid="perfect-kda"]').should('be.visible');
    cy.get('[data-testid="deaths"]').should('contain.text', '0');
    cy.get('[data-testid="performance-badge"]').should('contain.text', 'Flawless');
    
    // Test poor game
    cy.intercept('GET', '**/matches/*/analysis', { fixture: 'matches/poor-game.json' }).as('getPoorGame');
    
    cy.go('back');
    cy.get('[data-testid="match-card"]').eq(1).click();
    cy.waitForAPI('@getPoorGame');
    
    // Should show areas for improvement
    cy.get('[data-testid="improvement-areas"]').should('be.visible');
    cy.get('[data-testid="recommendation-item"]').should('have.length.greaterThan', 3);
  });

  it('should export match data', () => {
    cy.visit('/matches');
    cy.waitForAPI('@getMatches');
    
    cy.get('[data-testid="match-card"]').first().click();
    cy.waitForAPI('@getMatchAnalysis');
    
    // Should have export functionality
    cy.get('[data-testid="export-match"]').should('be.visible');
    cy.get('[data-testid="export-match"]').click();
    
    // Should show export options
    cy.get('[data-testid="export-options"]').should('be.visible');
    cy.get('[data-testid="export-pdf"]').should('exist');
    cy.get('[data-testid="export-json"]').should('exist');
    cy.get('[data-testid="export-csv"]').should('exist');
  });

  it('should work offline with cached data', () => {
    // Load match first
    cy.visit('/matches');
    cy.waitForAPI('@getMatches');
    
    cy.get('[data-testid="match-card"]').first().click();
    cy.waitForAPI('@getMatchAnalysis');
    
    // Simulate offline
    cy.intercept('GET', '**/matches/**', { forceNetworkError: true }).as('offlineRequest');
    
    // Should still show cached data
    cy.reload();
    cy.get('[data-testid="match-overview"]').should('be.visible');
    cy.get('[data-testid="offline-indicator"]').should('be.visible');
  });

  it('should handle real-time match updates', () => {
    // Mock live match
    cy.intercept('GET', '**/matches/live/*', { fixture: 'matches/live-match.json' }).as('getLiveMatch');
    
    cy.visit('/matches/live/123');
    cy.waitForAPI('@getLiveMatch');
    
    // Should show live indicators
    cy.get('[data-testid="live-indicator"]').should('be.visible');
    cy.get('[data-testid="game-timer"]').should('be.visible');
    
    // Should update in real-time
    cy.get('[data-testid="gold-difference"]').should('contain.text', '+1,500');
    
    // Simulate update
    cy.intercept('GET', '**/matches/live/*', { 
      body: { goldDifference: 2000, gameTime: 1200 } 
    }).as('getUpdatedMatch');
    
    // Should show updated values
    cy.wait(1000); // Simulate real-time update
    cy.get('[data-testid="gold-difference"]').should('contain.text', '+2,000');
  });
});