// Herald.lol Visual Regression Testing Configuration
import { defineConfig, devices } from '@playwright/test';

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
  testDir: './visual-tests',
  
  /* Run tests in files in parallel */
  fullyParallel: true,
  
  /* Fail the build on CI if you accidentally left test.only in the source code. */
  forbidOnly: !!process.env.CI,
  
  /* Retry on CI only */
  retries: process.env.CI ? 2 : 0,
  
  /* Opt out of parallel tests on CI. */
  workers: process.env.CI ? 1 : undefined,
  
  /* Reporter to use. See https://playwright.dev/docs/test-reporters */
  reporter: [
    ['html'],
    ['json', { outputFile: 'test-results/visual-results.json' }],
    ['@percy/playwright', { 
      percy: {
        snapshot: {
          widths: [375, 768, 1280, 1920], // Gaming responsive breakpoints
          minHeight: 1024,
        }
      }
    }]
  ],
  
  /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
  use: {
    /* Base URL to use in actions like `await page.goto('/')`. */
    baseURL: process.env.FRONTEND_URL || 'http://localhost:3000',
    
    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'on-first-retry',
    
    /* Take screenshot on failures */
    screenshot: 'only-on-failure',
    
    /* Video recording */
    video: 'retain-on-failure',
    
    /* Herald.lol gaming-specific settings */
    viewport: { width: 1280, height: 720 }, // Gaming standard viewport
    ignoreHTTPSErrors: true,
    
    /* Gaming analytics performance timeout */
    actionTimeout: 5000, // Herald.lol <5s requirement
    navigationTimeout: 10000,
  },
  
  /* Configure projects for major browsers */
  projects: [
    {
      name: 'chromium-desktop',
      use: { 
        ...devices['Desktop Chrome'],
        viewport: { width: 1920, height: 1080 }, // Gaming desktop
      },
    },
    
    {
      name: 'firefox-desktop', 
      use: { 
        ...devices['Desktop Firefox'],
        viewport: { width: 1920, height: 1080 },
      },
    },
    
    {
      name: 'webkit-desktop',
      use: { 
        ...devices['Desktop Safari'],
        viewport: { width: 1920, height: 1080 },
      },
    },
    
    /* Gaming tablet viewports */
    {
      name: 'tablet',
      use: { 
        ...devices['iPad Pro'],
        viewport: { width: 1024, height: 768 }, // Gaming tablet
      },
    },
    
    /* Gaming mobile viewports */
    {
      name: 'mobile-chrome',
      use: { 
        ...devices['Pixel 5'],
        viewport: { width: 393, height: 851 }, // Gaming mobile
      },
    },
    
    {
      name: 'mobile-safari',
      use: { 
        ...devices['iPhone 12'],
        viewport: { width: 390, height: 844 },
      },
    },
  ],
  
  /* Run your local dev server before starting the tests */
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
    timeout: 120000, // 2 minutes for Herald.lol startup
  },
  
  /* Visual regression specific settings */
  expect: {
    // Herald.lol visual regression thresholds
    toHaveScreenshot: { 
      threshold: 0.2,        // 20% threshold for gaming UI changes
      maxDiffPixels: 1000,   // Max different pixels allowed
      animations: 'disabled', // Disable animations for consistency
    },
    toMatchSnapshot: {
      threshold: 0.25,       // Gaming component visual threshold
      maxDiffPixels: 500,
    },
  },
  
  /* Global test timeout */
  timeout: 30000, // 30s for gaming analytics loading
  
  /* Test match patterns */
  testMatch: [
    '**/*.visual.test.{js,ts}',
    '**/*.screenshot.test.{js,ts}',
    '**/visual-regression/**/*.test.{js,ts}',
  ],
  
  /* Global setup for gaming visual tests */
  globalSetup: './global-setup.ts',
  globalTeardown: './global-teardown.ts',
});