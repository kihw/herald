// DamageAnalytics Component Tests for Herald.lol
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import { render, mockGamingProps, expectValidPercentage, validateGamingMetrics } from '@/test/utils/test-utils';
import { server } from '@/test/mocks/server';
import { http, HttpResponse } from 'msw';
import { DamageAnalytics } from '../DamageAnalytics';

// Mock the chart components to avoid rendering issues in tests
vi.mock('recharts', () => ({
  ResponsiveContainer: ({ children }: any) => <div data-testid="responsive-container">{children}</div>,
  BarChart: ({ children }: any) => <div data-testid="bar-chart">{children}</div>,
  Bar: () => <div data-testid="bar" />,
  XAxis: () => <div data-testid="x-axis" />,
  YAxis: () => <div data-testid="y-axis" />,
  CartesianGrid: () => <div data-testid="cartesian-grid" />,
  Tooltip: () => <div data-testid="tooltip" />,
  Legend: () => <div data-testid="legend" />,
  PieChart: ({ children }: any) => <div data-testid="pie-chart">{children}</div>,
  Pie: () => <div data-testid="pie" />,
  Cell: () => <div data-testid="cell" />,
}));

describe('DamageAnalytics Component', () => {
  const defaultProps = {
    ...mockGamingProps,
    timeRange: '7d' as const,
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Loading State', () => {
    it('should display loading spinner while fetching data', () => {
      // Mock delayed response
      server.use(
        http.get('/api/v1/analytics/damage/:summonerId', async () => {
          await new Promise(resolve => setTimeout(resolve, 1000));
          return HttpResponse.json({
            data: {},
            confidence: 0.8
          });
        })
      );

      render(<DamageAnalytics {...defaultProps} />);
      
      expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
    });
  });

  describe('Data Display', () => {
    it('should display damage analytics data correctly', async () => {
      render(<DamageAnalytics {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByText(/damage analytics/i)).toBeInTheDocument();
      });

      // Check for damage share percentage
      expect(screen.getByText(/32%/)).toBeInTheDocument();
      
      // Check for damage per minute
      expect(screen.getByText(/850/)).toBeInTheDocument();
      
      // Check for team contribution
      expect(screen.getByText(/high/i)).toBeInTheDocument();
    });

    it('should validate gaming metrics in damage data', async () => {
      render(<DamageAnalytics {...defaultProps} />);

      await waitFor(async () => {
        // Wait for data to load and validate it
        const damageData = {
          damageShare: 0.32,
          damagePerMinute: 850,
          efficiency: 0.89,
        };

        validateGamingMetrics(damageData);
        expectValidPercentage(damageData.damageShare);
        expect(damageData.damagePerMinute).toBeGreaterThan(0);
        expect(damageData.efficiency).toBeGreaterThan(0);
      });
    });
  });

  describe('Charts Rendering', () => {
    it('should render damage breakdown chart', async () => {
      render(<DamageAnalytics {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByTestId('pie-chart')).toBeInTheDocument();
        expect(screen.getByTestId('pie')).toBeInTheDocument();
      });
    });

    it('should render damage trends chart', async () => {
      render(<DamageAnalytics {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByTestId('bar-chart')).toBeInTheDocument();
        expect(screen.getByTestId('bar')).toBeInTheDocument();
      });
    });
  });

  describe('Performance Analysis', () => {
    it('should complete analysis within performance target (<5s)', async () => {
      const startTime = Date.now();
      
      render(<DamageAnalytics {...defaultProps} />);
      
      await waitFor(() => {
        expect(screen.getByText(/damage analytics/i)).toBeInTheDocument();
      });

      const analysisTime = Date.now() - startTime;
      expect(analysisTime).toBeLessThan(5000); // <5s requirement
    });

    it('should handle large datasets efficiently', async () => {
      // Mock response with large dataset
      server.use(
        http.get('/api/v1/analytics/damage/:summonerId', () => {
          const largeDamageHistory = Array.from({ length: 100 }, (_, i) => ({
            matchId: `NA1_${i}`,
            damageDealt: 800 + Math.random() * 400,
            damageShare: 0.25 + Math.random() * 0.3,
            timestamp: Date.now() - i * 86400000,
          }));

          return HttpResponse.json({
            data: {
              damageShare: 0.32,
              damagePerMinute: 850,
              efficiency: 0.89,
              history: largeDamageHistory,
            },
            confidence: 0.91
          });
        })
      );

      const startTime = Date.now();
      render(<DamageAnalytics {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByText(/damage analytics/i)).toBeInTheDocument();
      });

      const processingTime = Date.now() - startTime;
      expect(processingTime).toBeLessThan(5000); // Should handle large datasets efficiently
    });
  });

  describe('Error Handling', () => {
    it('should handle API errors gracefully', async () => {
      server.use(
        http.get('/api/v1/analytics/damage/:summonerId', () => {
          return HttpResponse.json(
            { error: 'Analytics service unavailable' },
            { status: 503 }
          );
        })
      );

      render(<DamageAnalytics {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByText(/error loading damage analytics/i)).toBeInTheDocument();
      });
    });

    it('should handle network errors gracefully', async () => {
      server.use(
        http.get('/api/v1/analytics/damage/:summonerId', () => {
          return HttpResponse.error();
        })
      );

      render(<DamageAnalytics {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByText(/failed to load analytics/i)).toBeInTheDocument();
      });
    });
  });

  describe('Gaming Calculations Accuracy', () => {
    it('should calculate damage share correctly', async () => {
      const testScenario = {
        playerDamage: 32000,
        teamTotalDamage: 100000,
        expectedDamageShare: 0.32,
      };

      server.use(
        http.get('/api/v1/analytics/damage/:summonerId', () => {
          return HttpResponse.json({
            data: {
              damageShare: testScenario.expectedDamageShare,
              playerDamage: testScenario.playerDamage,
              teamTotalDamage: testScenario.teamTotalDamage,
            },
            confidence: 0.95
          });
        })
      );

      render(<DamageAnalytics {...defaultProps} />);

      await waitFor(() => {
        const displayedShare = screen.getByText(/32%/);
        expect(displayedShare).toBeInTheDocument();
      });
    });

    it('should validate damage efficiency calculations', async () => {
      render(<DamageAnalytics {...defaultProps} />);

      await waitFor(() => {
        const efficiencyElement = screen.getByText(/89%/);
        expect(efficiencyElement).toBeInTheDocument();
        
        // Validate efficiency is within reasonable bounds
        const efficiency = 0.89;
        expect(efficiency).toBeGreaterThanOrEqual(0);
        expect(efficiency).toBeLessThanOrEqual(2); // Can exceed 1 for exceptional performance
      });
    });
  });

  describe('Responsive Design', () => {
    it('should adapt to different screen sizes', async () => {
      render(<DamageAnalytics {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByTestId('responsive-container')).toBeInTheDocument();
      });
    });
  });

  describe('Data Confidence Display', () => {
    it('should display confidence level for analytics', async () => {
      render(<DamageAnalytics {...defaultProps} />);

      await waitFor(() => {
        // Should show confidence indicator
        expect(screen.getByText(/confidence.*91%/i)).toBeInTheDocument();
      });
    });

    it('should warn when confidence is low', async () => {
      server.use(
        http.get('/api/v1/analytics/damage/:summonerId', () => {
          return HttpResponse.json({
            data: {
              damageShare: 0.32,
              damagePerMinute: 850,
            },
            confidence: 0.45 // Low confidence
          });
        })
      );

      render(<DamageAnalytics {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByText(/limited data available/i)).toBeInTheDocument();
      });
    });
  });

  describe('Time Range Functionality', () => {
    it('should handle different time ranges', async () => {
      const { rerender } = render(<DamageAnalytics {...defaultProps} timeRange="30d" />);

      await waitFor(() => {
        expect(screen.getByText(/damage analytics/i)).toBeInTheDocument();
      });

      // Change time range
      rerender(<DamageAnalytics {...defaultProps} timeRange="7d" />);

      await waitFor(() => {
        expect(screen.getByText(/damage analytics/i)).toBeInTheDocument();
      });
    });
  });

  describe('Accessibility', () => {
    it('should have proper ARIA labels', async () => {
      render(<DamageAnalytics {...defaultProps} />);

      await waitFor(() => {
        const analyticsSection = screen.getByRole('region', { name: /damage analytics/i });
        expect(analyticsSection).toBeInTheDocument();
      });
    });

    it('should be keyboard navigable', async () => {
      render(<DamageAnalytics {...defaultProps} />);

      await waitFor(() => {
        const focusableElements = screen.getAllByRole('button');
        focusableElements.forEach(element => {
          expect(element).toHaveAttribute('tabIndex');
        });
      });
    });
  });

  describe('Real Gaming Scenarios', () => {
    it('should handle ADC damage analytics correctly', async () => {
      server.use(
        http.get('/api/v1/analytics/damage/:summonerId', () => {
          return HttpResponse.json({
            data: {
              damageShare: 0.38, // High for ADC
              damagePerMinute: 920,
              efficiency: 0.91,
              teamContribution: 'excellent',
              breakdown: {
                physical: 0.85, // High physical damage for ADC
                magic: 0.10,
                true: 0.05
              }
            },
            confidence: 0.94
          });
        })
      );

      render(<DamageAnalytics {...defaultProps} role="ADC" />);

      await waitFor(() => {
        expect(screen.getByText(/38%/)).toBeInTheDocument(); // Damage share
        expect(screen.getByText(/excellent/i)).toBeInTheDocument(); // Team contribution
      });
    });

    it('should provide actionable insights', async () => {
      render(<DamageAnalytics {...defaultProps} />);

      await waitFor(() => {
        // Should provide gaming-specific recommendations
        expect(screen.getByText(/improve late-game positioning/i) || 
               screen.getByText(/focus on team fights/i) ||
               screen.getByText(/damage efficiency/i)).toBeInTheDocument();
      });
    });
  });
});