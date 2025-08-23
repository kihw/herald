import { defineConfig } from 'cypress';

export default defineConfig({
  e2e: {
    baseUrl: 'http://localhost:3000',
    supportFile: 'cypress/support/e2e.ts',
    specPattern: 'cypress/e2e/**/*.cy.{js,jsx,ts,tsx}',
    fixturesFolder: 'cypress/fixtures',
    screenshotsFolder: 'cypress/screenshots',
    videosFolder: 'cypress/videos',
    
    // Viewport settings
    viewportWidth: 1280,
    viewportHeight: 720,
    
    // Herald.lol Gaming Performance Requirements
    defaultCommandTimeout: 10000,
    requestTimeout: 15000,
    responseTimeout: 15000,
    pageLoadTimeout: 30000,
    
    // Gaming Analytics Performance Validation
    taskTimeout: 60000,
    
    // Video and screenshot settings
    video: true,
    screenshotOnRunFailure: true,
    
    setupNodeEvents(on, config) {
      // Gaming performance monitoring
      on('task', {
        log(message) {
          console.log(message);
          return null;
        },
        
        // Custom task for gaming analytics performance validation
        validateGamingPerformance(data) {
          const { startTime, endTime, type } = data;
          const duration = endTime - startTime;
          
          // Herald.lol performance requirements
          const thresholds = {
            analytics: 5000,  // <5s for analytics
            ui: 2000,         // <2s for UI load
            api: 1000,        // <1s for API calls
          };
          
          const threshold = thresholds[type] || 5000;
          
          if (duration > threshold) {
            throw new Error(`Performance requirement failed: ${type} took ${duration}ms (max ${threshold}ms)`);
          }
          
          return { duration, threshold, passed: true };
        },
        
        // Custom task for gaming metrics validation
        validateGamingMetrics(data) {
          const { type, value } = data;
          
          // Gaming metrics validation rules
          const validationRules = {
            kda: (val) => val >= 0 && val <= 50 && !isNaN(val),
            cs: (val) => val >= 0 && val <= 15 && !isNaN(val),
            vision: (val) => val >= 0 && val <= 200 && !isNaN(val),
            damage: (val) => val >= 0 && !isNaN(val),
            gold: (val) => val >= 0 && val <= 1000 && !isNaN(val),
            percentile: (val) => val >= 0 && val <= 100 && Number.isInteger(val),
            efficiency: (val) => val >= 0 && val <= 100 && Number.isInteger(val),
          };
          
          const validator = validationRules[type];
          if (!validator) {
            throw new Error(`Unknown gaming metric type: ${type}`);
          }
          
          const isValid = validator(value);
          if (!isValid) {
            throw new Error(`Invalid ${type} value: ${value}`);
          }
          
          return { type, value, valid: true };
        },
        
        // Task for logging gaming test results
        logGamingTestResult(data) {
          const timestamp = new Date().toISOString();
          console.log(`[${timestamp}] Herald.lol Gaming Test:`, data);
          return null;
        },
      });
      
      // Gaming-specific browser launch options
      on('before:browser:launch', (browser, launchOptions) => {
        // Optimize for gaming analytics testing
        if (browser.name === 'chrome') {
          launchOptions.args.push('--disable-web-security');
          launchOptions.args.push('--disable-features=VizDisplayCompositor');
          launchOptions.args.push('--no-sandbox');
          launchOptions.args.push('--disable-dev-shm-usage');
          
          // Performance optimization for gaming tests
          launchOptions.args.push('--max_old_space_size=4096');
          launchOptions.args.push('--memory-pressure-off');
        }
        
        return launchOptions;
      });
      
      // Gaming test result reporting
      on('after:spec', (spec, results) => {
        if (results && results.stats) {
          console.log(`Herald.lol Gaming Test Complete: ${spec.relative}`);
          console.log(`✅ Passed: ${results.stats.passes}`);
          console.log(`❌ Failed: ${results.stats.failures}`);
          console.log(`⏱️  Duration: ${results.stats.duration}ms`);
        }
      });

      // Code coverage setup (if needed)
      // require('@cypress/code-coverage/task')(on, config);
      
      return config;
    },

    env: {
      // Test user credentials
      TEST_USER_EMAIL: 'test@herald.lol',
      TEST_USER_PASSWORD: 'TestPassword123!',
      TEST_SUMMONER_NAME: 'TestSummoner',
      
      // Gaming performance thresholds
      ANALYTICS_LOAD_TIMEOUT: 5000, // <5s requirement
      UI_LOAD_TIMEOUT: 2000,        // <2s requirement
      
      // API endpoints
      API_BASE_URL: 'http://localhost:8080/api',
      
      // Gaming data validation
      VALIDATE_GAMING_METRICS: true,
      VALIDATE_PERFORMANCE: true,
      
      // Test data
      SAMPLE_MATCH_ID: 'NA1_4567890123',
      SAMPLE_PUUID: 'test-puuid-12345',
    },

    // Browser settings
    chromeWebSecurity: false,
    
    // Retry settings
    retries: {
      runMode: 2,
      openMode: 0,
    },
    
    // Exclude from tests
    excludeSpecPattern: [
      'cypress/e2e/examples/**',
      '**/__snapshots__/*',
      '**/__image_snapshots__/*'
    ],
    
    // Enhanced debugging for gaming tests
    experimentalStudio: true,
    experimentalSourceRewriting: true,
    
    // Gaming performance monitoring
    watchForFileChanges: true,
    numTestsKeptInMemory: 10,
  },

  component: {
    devServer: {
      framework: 'react',
      bundler: 'vite',
    },
    setupNodeEvents(on, config) {
      // component testing setup
    },
    specPattern: 'src/**/*.cy.{js,jsx,ts,tsx}',
    supportFile: 'cypress/support/component.ts',
  },
});