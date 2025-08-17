/**
 * Playwright Configuration for LoL Match Exporter E2E Tests
 * 
 * This configuration covers all E2E testing scenarios including:
 * - Multiple browsers and devices
 * - Different environments (dev, staging, prod)
 * - Performance testing
 * - Visual regression testing
 * - API testing
 */

const { defineConfig, devices } = require('@playwright/test');

module.exports = defineConfig({
  // Test directory
  testDir: './playwright',
  
  // Run tests in files in parallel
  fullyParallel: true,
  
  // Fail the build on CI if you accidentally left test.only in the source code
  forbidOnly: !!process.env.CI,
  
  // Retry on CI only
  retries: process.env.CI ? 2 : 0,
  
  // Opt out of parallel tests on CI
  workers: process.env.CI ? 1 : undefined,
  
  // Reporter to use
  reporter: [
    ['html', { outputFolder: '../../reports/playwright-html' }],
    ['json', { outputFile: '../../reports/playwright-results.json' }],
    ['junit', { outputFile: '../../reports/playwright-junit.xml' }],
    ['line']
  ],
  
  // Shared settings for all the projects below
  use: {
    // Base URL to use in actions like `await page.goto('/')`
    baseURL: process.env.TEST_BASE_URL || 'http://localhost:5173',
    
    // Collect trace when retrying the failed test
    trace: 'on-first-retry',
    
    // Capture screenshot after each test failure
    screenshot: 'only-on-failure',
    
    // Record video on failure
    video: 'retain-on-failure',
    
    // Global test timeout
    actionTimeout: 30000,
    navigationTimeout: 30000,
    
    // Ignore HTTPS errors
    ignoreHTTPSErrors: true,
    
    // Accept downloads
    acceptDownloads: true,
    
    // Locale for testing
    locale: 'en-US',
    timezoneId: 'Europe/Paris',
    
    // Extra HTTP headers
    extraHTTPHeaders: {
      'Accept-Language': 'en-US,en;q=0.9'
    }
  },

  // Test timeout
  timeout: 60000,
  
  // Global setup and teardown
  globalSetup: require.resolve('./global-setup.js'),
  globalTeardown: require.resolve('./global-teardown.js'),

  // Configure projects for major browsers
  projects: [
    // Desktop browsers
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
      testMatch: '**/*.spec.js'
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
      testMatch: '**/*.spec.js'
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
      testMatch: '**/*.spec.js'
    },
    
    // Mobile devices
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] },
      testMatch: '**/mobile/*.spec.js'
    },
    {
      name: 'Mobile Safari',
      use: { ...devices['iPhone 12'] },
      testMatch: '**/mobile/*.spec.js'
    },
    
    // Tablet devices
    {
      name: 'iPad',
      use: { ...devices['iPad Pro'] },
      testMatch: '**/tablet/*.spec.js'
    },
    
    // High DPI displays
    {
      name: 'High DPI',
      use: {
        ...devices['Desktop Chrome'],
        deviceScaleFactor: 2,
        viewport: { width: 1920, height: 1080 }
      },
      testMatch: '**/visual/*.spec.js'
    },
    
    // API testing
    {
      name: 'api',
      use: {
        // No browser needed for API tests
        browserName: undefined,
        baseURL: process.env.TEST_API_URL || 'http://localhost:8001'
      },
      testMatch: '**/api/*.spec.js'
    },
    
    // Performance testing
    {
      name: 'performance',
      use: {
        ...devices['Desktop Chrome'],
        // Enable performance metrics
        launchOptions: {
          args: ['--enable-precise-memory-info']
        }
      },
      testMatch: '**/performance/*.spec.js'
    },
    
    // Authentication tests (run first)
    {
      name: 'setup',
      testMatch: '**/auth.setup.js',
      teardown: 'cleanup'
    },
    
    // Cleanup after tests
    {
      name: 'cleanup',
      testMatch: '**/cleanup.spec.js'
    }
  ],

  // Dependencies between projects
  dependencies: [
    { name: 'setup', soft: true }
  ],

  // Web server configuration
  webServer: [
    {
      command: 'npm run dev',
      cwd: '../../web',
      port: 5173,
      reuseExistingServer: !process.env.CI,
      stdout: 'pipe',
      stderr: 'pipe',
      timeout: 120000
    },
    {
      command: 'go run analytics_server_standalone.go',
      cwd: '../..',
      port: 8001,
      reuseExistingServer: !process.env.CI,
      env: {
        PORT: '8001',
        GIN_MODE: 'debug'
      },
      timeout: 60000
    }
  ],

  // Global test configuration
  expect: {
    // Threshold for visual comparisons
    threshold: 0.2,
    // Animation handling
    toHaveScreenshot: { 
      threshold: 0.2, 
      animations: 'disabled' 
    },
    toMatchSnapshot: { 
      threshold: 0.2 
    }
  },

  // Test metadata
  metadata: {
    'Test Suite': 'LoL Match Exporter E2E Tests',
    'Environment': process.env.TEST_ENV || 'development',
    'Browser Versions': 'Latest',
    'Test Categories': [
      'Authentication',
      'Analytics Dashboard', 
      'Real-time Notifications',
      'Performance',
      'Mobile Responsiveness',
      'API Integration'
    ]
  }
});