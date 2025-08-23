// Herald.lol Visual Regression Testing Global Setup
import { chromium, FullConfig } from '@playwright/test';

async function globalSetup(config: FullConfig) {
  console.log('🎮 Herald.lol Visual Regression Testing Setup');
  console.log('============================================');
  
  // Launch browser for setup
  const browser = await chromium.launch();
  const context = await browser.newContext();
  const page = await context.newPage();
  
  try {
    // Check if Herald.lol backend is running
    const backendUrl = process.env.API_BASE_URL || 'http://localhost:8080';
    console.log(`🔍 Checking Herald.lol backend at: ${backendUrl}`);
    
    try {
      const response = await page.goto(`${backendUrl}/health`, { 
        waitUntil: 'networkidle',
        timeout: 10000 
      });
      
      if (response?.ok()) {
        console.log('✅ Herald.lol backend is running');
      } else {
        console.log('⚠️  Herald.lol backend health check failed');
      }
    } catch (error) {
      console.log('⚠️  Could not connect to Herald.lol backend');
      console.log('   Make sure the backend is running on port 8080');
    }
    
    // Check if Herald.lol frontend is running
    const frontendUrl = process.env.FRONTEND_URL || 'http://localhost:3000';
    console.log(`🔍 Checking Herald.lol frontend at: ${frontendUrl}`);
    
    try {
      const response = await page.goto(frontendUrl, { 
        waitUntil: 'networkidle',
        timeout: 10000 
      });
      
      if (response?.ok()) {
        console.log('✅ Herald.lol frontend is running');
      } else {
        console.log('⚠️  Herald.lol frontend check failed');
      }
    } catch (error) {
      console.log('⚠️  Could not connect to Herald.lol frontend');
      console.log('   Make sure the frontend is running on port 3000');
    }
    
    // Setup visual test data
    console.log('📊 Setting up visual test data...');
    
    // Create test user for visual testing
    try {
      await page.request.post(`${backendUrl}/api/auth/register`, {
        data: {
          email: 'visualtest@herald.lol',
          username: 'VisualTestUser',
          password: 'VisualTest123!',
          displayName: 'Visual Test User'
        }
      });
      console.log('✅ Visual test user created');
    } catch (error) {
      console.log('ℹ️  Visual test user may already exist');
    }
    
    // Pre-warm gaming analytics data
    console.log('🎮 Pre-warming gaming analytics data...');
    
    const testData = {
      kda: { currentKDA: 2.34, trend: 'improving', percentile: 65 },
      cs: { currentCSPerMin: 7.2, efficiency: 85, percentile: 72 },
      vision: { averageVisionScore: 28.4, efficiency: 78, percentile: 58 },
      damage: { damageShare: 0.32, efficiency: 89, percentile: 71 },
      gold: { goldPerMinute: 425, efficiency: 87, percentile: 68 }
    };
    
    // Store test data in local storage for consistent visual testing
    await context.addInitScript(testData => {
      window.localStorage.setItem('herald_visual_test_data', JSON.stringify(testData));
      window.localStorage.setItem('herald_visual_test_mode', 'true');
    }, testData);
    
    console.log('✅ Gaming analytics test data pre-loaded');
    
    // Configure visual testing environment
    console.log('🎯 Configuring visual testing environment...');
    
    // Set consistent font rendering
    await context.addInitScript(() => {
      // Disable font smoothing variations
      const style = document.createElement('style');
      style.textContent = `
        * {
          -webkit-font-smoothing: antialiased !important;
          -moz-osx-font-smoothing: grayscale !important;
          text-rendering: optimizeLegibility !important;
        }
        
        /* Disable animations for consistent screenshots */
        *, *::before, *::after {
          animation-duration: 0s !important;
          animation-delay: 0s !important;
          transition-duration: 0s !important;
          transition-delay: 0s !important;
        }
      `;
      document.head.appendChild(style);
    });
    
    // Disable network requests to external services during visual testing
    await context.route('**/*', async (route) => {
      const url = route.request().url();
      
      // Block external analytics/tracking
      if (url.includes('google-analytics.com') || 
          url.includes('googletagmanager.com') ||
          url.includes('hotjar.com') ||
          url.includes('mixpanel.com')) {
        await route.abort();
        return;
      }
      
      // Block external fonts (use system fonts for consistency)
      if (url.includes('fonts.googleapis.com') || 
          url.includes('fonts.gstatic.com')) {
        await route.abort();
        return;
      }
      
      // Continue with other requests
      await route.continue();
    });
    
    console.log('✅ Visual testing environment configured');
    
    // Test screenshot capability
    console.log('📸 Testing screenshot capability...');
    
    try {
      await page.goto(frontendUrl);
      await page.screenshot({ 
        path: 'test-results/setup-screenshot.png',
        fullPage: false 
      });
      console.log('✅ Screenshot capability verified');
    } catch (error) {
      console.log('⚠️  Screenshot test failed:', error);
    }
    
    // Setup Percy integration (if available)
    if (process.env.PERCY_TOKEN) {
      console.log('🎨 Percy visual testing integration detected');
      console.log('   Project:', process.env.PERCY_PROJECT || 'herald-lol');
      console.log('   Branch:', process.env.PERCY_BRANCH || 'main');
    } else {
      console.log('ℹ️  Percy integration not configured (PERCY_TOKEN not set)');
    }
    
    // Gaming-specific setup
    console.log('🎮 Gaming-specific visual setup...');
    
    // Mock Riot API responses for consistent visuals
    await context.route('**/api/riot/**', async (route) => {
      const mockRiotResponse = {
        summoner: {
          name: 'VisualTestSummoner',
          level: 147,
          profileIcon: 1234,
        },
        rank: {
          tier: 'GOLD',
          rank: 'III',
          leaguePoints: 67,
          winRate: 0.643,
        },
        matches: [
          { matchId: 'NA1_4567890123', result: 'victory', champion: 'Jinx' },
          { matchId: 'NA1_4567890124', result: 'victory', champion: 'Caitlyn' },
          { matchId: 'NA1_4567890125', result: 'defeat', champion: 'Ezreal' },
        ]
      };
      
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockRiotResponse)
      });
    });
    
    console.log('✅ Gaming API mocks configured');
    
    // Performance monitoring for visual tests
    console.log('⚡ Performance monitoring setup...');
    
    await context.addInitScript(() => {
      // Track visual test performance
      window.visualTestMetrics = {
        startTime: Date.now(),
        analyticsLoadTime: null,
        screenshotCount: 0,
      };
      
      // Monitor Herald.lol performance during visual tests
      const observer = new PerformanceObserver((list) => {
        list.getEntries().forEach((entry) => {
          if (entry.name.includes('analytics') && !window.visualTestMetrics.analyticsLoadTime) {
            window.visualTestMetrics.analyticsLoadTime = entry.duration;
          }
        });
      });
      
      observer.observe({ entryTypes: ['navigation', 'resource'] });
    });
    
    console.log('✅ Performance monitoring configured');
    
    console.log('');
    console.log('🎮 Herald.lol Visual Regression Testing Ready!');
    console.log('============================================');
    console.log('📊 Gaming Analytics: Ready for visual testing');
    console.log('🎯 Performance Target: <5s analytics load');
    console.log('📸 Screenshot Engine: Playwright');
    console.log('🎨 Visual Diff Engine: Percy (if configured)');
    console.log('');
    
  } catch (error) {
    console.error('❌ Visual testing setup failed:', error);
    throw error;
  } finally {
    await context.close();
    await browser.close();
  }
}

export default globalSetup;