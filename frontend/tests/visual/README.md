# Herald.lol Visual Regression Testing

Comprehensive visual regression testing suite for Herald.lol gaming analytics platform using Playwright and Percy.

## 🎮 Gaming Visual Testing Overview

This directory contains visual regression tests specifically designed for Herald.lol's gaming analytics components, ensuring consistent UI/UX across all gaming features.

### Gaming Components Tested
- **Analytics Dashboard** - KDA, CS/min, Vision Score, Damage, Gold widgets
- **Match Analysis** - Post-game analysis interface and recommendations
- **Gaming Components** - Team composition, counter-picks, skill progression
- **Gaming Charts** - Performance visualizations and gaming metrics
- **Responsive Design** - Gaming-optimized layouts across all devices

## 🚀 Quick Start

### Prerequisites

1. **Install dependencies:**
   ```bash
   npm install
   npm run playwright:install
   npm run playwright:install-deps
   ```

2. **Start Herald.lol services:**
   ```bash
   # Backend (in separate terminal)
   cd ../backend && go run main.go
   
   # Frontend (in separate terminal)  
   npm run dev
   ```

### Running Visual Tests

```bash
# Run all visual regression tests
npm run test:visual

# Run with UI for debugging
npm run test:visual:ui

# Debug mode (step through tests)
npm run test:visual:debug

# Update visual baselines (after approved changes)
npm run test:visual:update

# Generate and view reports
npm run test:visual:report
```

## 📊 Test Structure

### Gaming Test Files

#### `gaming-analytics-dashboard.visual.test.ts`
Comprehensive testing of the gaming analytics dashboard:
- **Desktop/Tablet/Mobile layouts** - Responsive gaming design
- **KDA Widget states** - Default, hover, loading states
- **CS/min Widget** - Default and expanded views  
- **Vision Widget** - Default and heatmap visualizations
- **Damage Widget** - Default and chart views
- **Gold Widget** - Default and timeline views
- **Dark theme** - Gaming theme variations
- **Error states** - API failures and loading errors
- **Performance indicators** - Gaming metric confidence levels

#### `match-analysis.visual.test.ts`
Complete match analysis interface testing:
- **Match overview** - Desktop and mobile layouts
- **Match header** - Victory/defeat styling
- **Gaming metrics grid** - All metric cards (KDA, CS, Vision, Damage, Gold)
- **Match timeline** - Event visualization
- **Performance charts** - Multiple chart views
- **Recommendations panel** - Gaming improvement suggestions
- **Comparison views** - Match-to-match comparisons
- **Gaming visualizations** - Heatmaps, patterns, movement tracking
- **Responsive breakpoints** - 7 different viewport sizes
- **Theme variations** - Light and dark gaming themes

#### `gaming-components.visual.test.ts`
Individual gaming component testing:
- **Team Composition Optimizer** - Strategy selection and optimization
- **Counter-Pick Analyzer** - Champion selection and recommendations
- **Skill Progression Tracker** - Individual skill categories and progress
- **Gaming Charts** - KDA trends, CS performance, vision radar, damage distribution
- **Gaming UI Cards** - Rank, winrate, matches, mastery cards
- **Gaming Forms** - Summoner search, champion filters
- **Gaming Navigation** - Main nav, breadcrumbs, tabs
- **Gaming Modals** - Settings, match details, confirmations
- **Loading/Empty States** - Spinners, skeletons, no-data states
- **Error Components** - Error boundaries, API errors, network failures

### Configuration Files

#### `playwright.config.ts`
Main Playwright configuration:
- **Gaming viewports** - Desktop (1920x1080), tablet (1024x768), mobile (390x844)
- **Herald.lol performance** - <5s timeout for analytics loading
- **Visual thresholds** - 20% page threshold, 25% component threshold
- **Percy integration** - Visual diff service configuration
- **Multi-browser** - Chrome, Firefox, Safari testing
- **Animation handling** - Disabled for consistent screenshots

#### `global-setup.ts`
Pre-test environment setup:
- **Herald.lol services check** - Backend/frontend health validation
- **Gaming test data** - Pre-loaded mock analytics data
- **Visual consistency** - Font rendering and animation configuration
- **Gaming API mocks** - Riot API response mocking
- **Performance monitoring** - Gaming analytics load time tracking

#### `global-teardown.ts`
Post-test cleanup and reporting:
- **Visual test summary** - Comprehensive markdown report
- **Gaming performance metrics** - Analytics load time validation
- **Screenshot organization** - File management and cleanup
- **Percy integration** - Build summary and results

## 🎯 Gaming-Specific Features

### Performance Validation
- **Analytics Load Time** - <5 second Herald.lol requirement
- **UI Response Time** - <2 second interface requirement
- **Gaming API Mocking** - <100ms consistent response times

### Gaming Metrics Testing
All Herald.lol gaming metrics are visually validated:
- **KDA Analysis** - Kill/Death/Assist calculations and trends
- **CS/min Analysis** - Creep Score farming efficiency
- **Vision Score** - Map control and ward placement analysis
- **Damage Analysis** - Team fight contribution and patterns
- **Gold Efficiency** - Economic performance and optimization

### Responsive Gaming Design
Testing across gaming-relevant viewports:
- **Gaming Desktop** - 1920x1080 (primary gaming resolution)
- **Gaming Laptop** - 1440x900, 1280x720 (portable gaming)
- **Gaming Tablet** - 1024x768 (analysis on tablet)
- **Gaming Mobile** - 414x896, 375x667 (mobile analytics)

### Gaming Theme Support
- **Light Theme** - Default gaming interface
- **Dark Theme** - Gaming-preferred dark mode
- **High Contrast** - Accessibility gaming mode (if implemented)

## 🔧 Configuration

### Environment Variables
```bash
# Herald.lol service URLs
export FRONTEND_URL="http://localhost:3000"
export API_BASE_URL="http://localhost:8080"

# Percy visual testing (optional)
export PERCY_TOKEN="your-percy-token"
export PERCY_PROJECT="herald-lol"
export PERCY_BRANCH="main"

# Performance testing
export ANALYTICS_TIMEOUT="5000"  # 5s Herald.lol requirement
export UI_TIMEOUT="2000"         # 2s UI requirement
```

### Visual Regression Thresholds
```typescript
// Component-level testing
toHaveScreenshot: {
  threshold: 0.25,        // 25% pixel difference allowed
  maxDiffPixels: 500,     // Max 500 different pixels
  animations: 'disabled', // No animations for consistency
}

// Page-level testing  
toMatchSnapshot: {
  threshold: 0.2,         // 20% pixel difference allowed
  maxDiffPixels: 1000,    // Max 1000 different pixels
}
```

### Gaming Mock Data
Pre-configured gaming analytics data for consistent visuals:
```javascript
const GAMING_TEST_DATA = {
  kda: { currentKDA: 2.34, trend: 'improving', percentile: 65 },
  cs: { currentCSPerMin: 7.2, efficiency: 85, percentile: 72 },
  vision: { averageVisionScore: 28.4, efficiency: 78, percentile: 58 },
  damage: { damageShare: 0.32, efficiency: 89, percentile: 71 },
  gold: { goldPerMinute: 425, efficiency: 87, percentile: 68 }
};
```

## 📊 Results Analysis

### Generated Files
- **Screenshots** - `test-results/*.png` (baseline screenshots)
- **Visual Diffs** - `test-results/*-diff.png` (difference highlights)
- **Test Results** - `test-results/visual-results.json` (detailed results)
- **HTML Report** - `test-results/index.html` (interactive report)
- **Summary Report** - `test-results/visual-test-summary.md` (comprehensive summary)

### Percy Integration
If Percy is configured, visual diffs are also available in the Percy dashboard:
- **Build Comparisons** - Branch-to-branch visual comparisons
- **Review Workflow** - Approve/reject visual changes
- **Team Collaboration** - Share visual feedback
- **Historical Tracking** - Visual regression history

### Performance Metrics
Visual tests also validate Herald.lol performance requirements:
- **Analytics Load Time** - Must be <5 seconds
- **Component Render Time** - Must be <2 seconds
- **Screenshot Generation** - Optimized for gaming components
- **Gaming API Response** - Mocked for <100ms consistency

## 🎮 Gaming Test Scenarios

### Dashboard Analytics
- ✅ All gaming widgets load and display correctly
- ✅ Responsive layout adapts to gaming viewports
- ✅ Dark theme maintains gaming aesthetics
- ✅ Loading states provide clear feedback
- ✅ Error states gracefully handle failures

### Match Analysis
- ✅ Post-game analysis displays all gaming metrics
- ✅ Timeline visualization shows key game events
- ✅ Performance charts render gaming statistics
- ✅ Recommendations provide actionable gaming advice
- ✅ Heatmaps and visualizations enhance understanding

### Gaming Components
- ✅ Team composition optimizer suggests optimal picks
- ✅ Counter-pick analyzer shows champion matchups
- ✅ Skill progression tracks gaming improvement
- ✅ Gaming forms enable efficient summoner search
- ✅ Navigation supports gaming workflow

### Error Handling
- ✅ Riot API failures display helpful messages
- ✅ Network errors maintain gaming interface
- ✅ Loading failures provide retry options
- ✅ Data validation ensures gaming metric accuracy

## 🛠️ Development Workflow

### Adding New Visual Tests
1. Create test file in appropriate category directory
2. Use Herald.lol gaming mock data for consistency
3. Test across all gaming-relevant viewports
4. Include gaming theme variations
5. Validate performance requirements

### Updating Visual Baselines
```bash
# After approved UI changes
npm run test:visual:update

# Review changes before committing
git diff test-results/
```

### Debugging Visual Failures
```bash
# Run with UI for interactive debugging
npm run test:visual:ui

# Run in headed mode to see browser
npm run test:visual:headed

# Debug specific test
npm run test:visual:debug -- --grep "Gaming Analytics Dashboard"
```

## 🤝 Contributing

### Gaming Visual Standards
- Follow Herald.lol gaming design system
- Maintain <5s analytics load performance
- Test across all gaming viewports
- Include dark theme variants
- Validate gaming metric accuracy

### Pull Request Checklist
- [ ] All visual tests pass
- [ ] Performance requirements met (<5s analytics)
- [ ] Gaming components tested across viewports
- [ ] Dark theme variations included
- [ ] Percy builds (if configured) reviewed and approved

## 📚 Resources

- [Playwright Testing Documentation](https://playwright.dev/docs/test-intro)
- [Percy Visual Testing](https://docs.percy.io/)
- [Herald.lol Design System](../../docs/design-system.md)
- [Gaming Performance Requirements](../../CLAUDE.md)

---

**Herald.lol Visual Regression Testing** - Ensuring consistent gaming analytics experience across all platforms 🎮