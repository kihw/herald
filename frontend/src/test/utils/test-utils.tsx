import { ReactElement } from 'react';
import { render, RenderOptions } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ThemeProvider } from '@mui/material/styles';
import { CssBaseline } from '@mui/material';
import { darkTheme } from '@/theme';
import { vi } from 'vitest';

// Create a test query client
const createTestQueryClient = () => {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        gcTime: 0,
      },
      mutations: {
        retry: false,
      },
    },
  });
};

interface AllTheProvidersProps {
  children: React.ReactNode;
}

const AllTheProviders = ({ children }: AllTheProvidersProps) => {
  const queryClient = createTestQueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider theme={darkTheme}>
        <CssBaseline />
        <BrowserRouter>
          {children}
        </BrowserRouter>
      </ThemeProvider>
    </QueryClientProvider>
  );
};

const customRender = (
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>
) => {
  return render(ui, { wrapper: AllTheProviders, ...options });
};

// Re-export everything
export * from '@testing-library/react';
export { customRender as render };

// Gaming-specific test utilities
export const mockGamingProps = {
  summonerId: 'test-summoner-id',
  summonerName: 'TestSummoner',
  region: 'na1',
  champion: 'Jinx',
  role: 'ADC',
};

export const mockMatchAnalysis = {
  matchId: 'NA1_4567890123',
  kda: 2.34,
  csPerMin: 7.2,
  visionScore: 28.4,
  damageShare: 0.32,
  goldEfficiency: 0.87,
  duration: 1834,
  result: 'victory',
  champion: 'Jinx',
  role: 'ADC',
};

export const mockPerformanceData = {
  overallRating: 75,
  recentTrend: 'improving' as const,
  strengthAreas: ['CS/min', 'Late game'],
  improvementAreas: ['Early game', 'Vision control'],
  confidenceLevel: 0.89,
};

// Gaming calculation test helpers
export const expectValidKDA = (kda: number) => {
  expect(kda).toBeTypeOf('number');
  expect(kda).toBeGreaterThanOrEqual(0);
  expect(Number.isFinite(kda)).toBe(true);
};

export const expectValidCSPerMin = (csPerMin: number) => {
  expect(csPerMin).toBeTypeOf('number');
  expect(csPerMin).toBeGreaterThanOrEqual(0);
  expect(csPerMin).toBeLessThan(15); // Reasonable upper bound
};

export const expectValidPercentage = (percentage: number) => {
  expect(percentage).toBeTypeOf('number');
  expect(percentage).toBeGreaterThanOrEqual(0);
  expect(percentage).toBeLessThanOrEqual(1);
};

export const expectValidDamageShare = (damageShare: number) => {
  expectValidPercentage(damageShare);
  expect(damageShare).toBeLessThan(0.8); // Unrealistic to have >80% damage share
};

// Performance test helper
export const waitForAnalysis = async (maxTime: number = 5000) => {
  const start = Date.now();
  
  return new Promise((resolve, reject) => {
    const checkCompletion = () => {
      const elapsed = Date.now() - start;
      if (elapsed >= maxTime) {
        reject(new Error(`Analysis took longer than ${maxTime}ms`));
      } else {
        // Simulate analysis completion
        setTimeout(resolve, 100);
      }
    };
    
    checkCompletion();
  });
};

// Mock user
export const mockUser = {
  id: '1',
  email: 'test@herald.lol',
  username: 'testuser',
  display_name: 'Test User',
  is_authenticated: true,
};

export const mockLocalStorage = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
};

// Test helpers
export const waitForLoadingToFinish = () => {
  return new Promise(resolve => setTimeout(resolve, 0));
};

export const createMockMatchMedia = (matches: boolean) => {
  return (query: string) => ({
    matches,
    media: query,
    onchange: null,
    addListener: () => {},
    removeListener: () => {},
    addEventListener: () => {},
    removeEventListener: () => {},
    dispatchEvent: () => {},
  });
};

// Gaming metrics validation
export const validateGamingMetrics = (metrics: any) => {
  if (metrics.kda !== undefined) expectValidKDA(metrics.kda);
  if (metrics.csPerMin !== undefined) expectValidCSPerMin(metrics.csPerMin);
  if (metrics.damageShare !== undefined) expectValidDamageShare(metrics.damageShare);
  if (metrics.winRate !== undefined) expectValidPercentage(metrics.winRate);
  if (metrics.visionScore !== undefined) {
    expect(metrics.visionScore).toBeGreaterThanOrEqual(0);
    expect(metrics.visionScore).toBeLessThan(200); // Reasonable upper bound
  }
};

// Gaming analytics test scenarios
export const gamingTestScenarios = {
  perfectGame: {
    kills: 20,
    deaths: 0,
    assists: 15,
    cs: 300,
    duration: 1800,
    visionScore: 45,
    expectedKDA: Infinity,
    expectedCSPerMin: 10.0,
  },
  
  averageGame: {
    kills: 8,
    deaths: 4,
    assists: 12,
    cs: 180,
    duration: 1800,
    visionScore: 28,
    expectedKDA: 5.0,
    expectedCSPerMin: 6.0,
  },
  
  difficultGame: {
    kills: 2,
    deaths: 8,
    assists: 6,
    cs: 120,
    duration: 1200,
    visionScore: 15,
    expectedKDA: 1.0,
    expectedCSPerMin: 6.0,
  },
};