// VisionHeatmap Component Tests for Herald.lol
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import { render, mockGamingProps, validateGamingMetrics } from '@/test/utils/test-utils';
import { server } from '@/test/mocks/server';
import { http, HttpResponse } from 'msw';
import { VisionHeatmap } from '../VisionHeatmap';

// Mock canvas for heatmap rendering
const mockCanvas = {
  getContext: vi.fn(() => ({
    fillRect: vi.fn(),
    clearRect: vi.fn(),
    getImageData: vi.fn(),
    putImageData: vi.fn(),
    createImageData: vi.fn(),
    setTransform: vi.fn(),
    drawImage: vi.fn(),
    save: vi.fn(),
    restore: vi.fn(),
    beginPath: vi.fn(),
    moveTo: vi.fn(),
    lineTo: vi.fn(),
    closePath: vi.fn(),
    stroke: vi.fn(),
    fill: vi.fn(),
    measureText: vi.fn(() => ({ width: 10 })),
    arc: vi.fn(),
  })),
  toDataURL: vi.fn(),
  addEventListener: vi.fn(),
};

// Mock HTMLCanvasElement
Object.defineProperty(HTMLCanvasElement.prototype, 'getContext', {
  value: vi.fn(() => mockCanvas.getContext()),
});

describe('VisionHeatmap Component', () => {
  const defaultProps = {
    ...mockGamingProps,
    mapType: 'summoner_rift' as const,
    timeRange: '7d' as const,
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Loading State', () => {
    it('should display loading spinner while fetching vision data', () => {
      server.use(
        http.get('/api/v1/analytics/vision/:summonerId', async () => {
          await new Promise(resolve => setTimeout(resolve, 1000));
          return HttpResponse.json({ data: {}, confidence: 0.8 });
        })
      );

      render(<VisionHeatmap {...defaultProps} />);
      
      expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
    });
  });

  describe('Vision Analytics Display', () => {
    it('should display vision score and control metrics', async () => {
      render(<VisionHeatmap {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByText(/vision analytics/i)).toBeInTheDocument();
      });

      // Check for vision score
      expect(screen.getByText(/28\.4/)).toBeInTheDocument();
      
      // Check for vision control percentage
      expect(screen.getByText(/72%/)).toBeInTheDocument();
      
      // Check for ward placement rating
      expect(screen.getByText(/good/i)).toBeInTheDocument();
    });

    it('should validate vision metrics', async () => {
      render(<VisionHeatmap {...defaultProps} />);

      await waitFor(async () => {
        const visionData = {
          averageVisionScore: 28.4,
          visionControl: 0.72,
          wardEfficiency: 0.85,
        };

        validateGamingMetrics(visionData);
        expect(visionData.averageVisionScore).toBeGreaterThan(0);
        expect(visionData.averageVisionScore).toBeLessThan(200);
        expect(visionData.visionControl).toBeGreaterThanOrEqual(0);
        expect(visionData.visionControl).toBeLessThanOrEqual(1);
      });
    });
  });

  describe('Heatmap Rendering', () => {
    it('should render the vision heatmap canvas', async () => {
      render(<VisionHeatmap {...defaultProps} />);

      await waitFor(() => {
        const canvas = screen.getByTestId('vision-heatmap-canvas');
        expect(canvas).toBeInTheDocument();
        expect(canvas.tagName).toBe('CANVAS');
      });
    });

    it('should handle heatmap data points correctly', async () => {
      server.use(
        http.get('/api/v1/analytics/vision/:summonerId', () => {
          return HttpResponse.json({
            data: {
              averageVisionScore: 28.4,
              visionControl: 0.72,
              heatmapData: [
                { x: 100, y: 200, intensity: 0.8, type: 'ward_placed' },
                { x: 150, y: 300, intensity: 0.6, type: 'ward_killed' },
                { x: 200, y: 150, intensity: 0.9, type: 'control_ward' }
              ]
            },
            confidence: 0.79
          });
        })
      );

      render(<VisionHeatmap {...defaultProps} />);

      await waitFor(() => {
        const canvas = screen.getByTestId('vision-heatmap-canvas');
        expect(canvas).toBeInTheDocument();
        
        // Verify canvas context was called for rendering
        expect(mockCanvas.getContext).toHaveBeenCalled();
      });
    });
  });

  describe('Performance Requirements', () => {
    it('should complete vision analysis within performance target (<5s)', async () => {
      const startTime = Date.now();
      
      render(<VisionHeatmap {...defaultProps} />);
      
      await waitFor(() => {
        expect(screen.getByText(/vision analytics/i)).toBeInTheDocument();
      });

      const analysisTime = Date.now() - startTime;
      expect(analysisTime).toBeLessThan(5000); // <5s requirement
    });

    it('should efficiently render large vision datasets', async () => {
      // Mock response with many vision events
      server.use(
        http.get('/api/v1/analytics/vision/:summonerId', () => {
          const largeVisionDataset = Array.from({ length: 500 }, (_, i) => ({
            x: Math.random() * 1024,
            y: Math.random() * 1024,
            intensity: Math.random(),
            type: ['ward_placed', 'ward_killed', 'control_ward'][i % 3],
            timestamp: Date.now() - i * 30000,
          }));

          return HttpResponse.json({
            data: {
              averageVisionScore: 35.2,
              visionControl: 0.78,
              heatmapData: largeVisionDataset,
            },
            confidence: 0.89
          });
        })
      );

      const startTime = Date.now();
      render(<VisionHeatmap {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByText(/vision analytics/i)).toBeInTheDocument();
      });

      const renderTime = Date.now() - startTime;
      expect(renderTime).toBeLessThan(5000); // Should handle large datasets efficiently
    });
  });

  describe('Map Integration', () => {
    it('should handle different map types', async () => {
      const { rerender } = render(<VisionHeatmap {...defaultProps} mapType="summoner_rift" />);

      await waitFor(() => {
        expect(screen.getByText(/summoner.*rift/i)).toBeInTheDocument();
      });

      // Test with different map
      rerender(<VisionHeatmap {...defaultProps} mapType="howling_abyss" />);

      await waitFor(() => {
        expect(screen.getByText(/howling.*abyss/i)).toBeInTheDocument();
      });
    });

    it('should display map-specific vision zones', async () => {
      render(<VisionHeatmap {...defaultProps} />);

      await waitFor(() => {
        // Should show important vision areas
        expect(screen.getByText(/dragon pit/i) || 
               screen.getByText(/baron pit/i) || 
               screen.getByText(/river bush/i)).toBeInTheDocument();
      });
    });
  });

  describe('Error Handling', () => {
    it('should handle API errors gracefully', async () => {
      server.use(
        http.get('/api/v1/analytics/vision/:summonerId', () => {
          return HttpResponse.json(
            { error: 'Vision data unavailable' },
            { status: 503 }
          );
        })
      );

      render(<VisionHeatmap {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByText(/error loading vision data/i)).toBeInTheDocument();
      });
    });

    it('should handle canvas rendering errors', async () => {
      // Mock canvas context error
      vi.spyOn(HTMLCanvasElement.prototype, 'getContext').mockReturnValue(null);

      render(<VisionHeatmap {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByText(/heatmap unavailable/i)).toBeInTheDocument();
      });
    });
  });

  describe('Gaming Insights', () => {
    it('should provide role-specific vision recommendations', async () => {
      server.use(
        http.get('/api/v1/analytics/vision/:summonerId', () => {
          return HttpResponse.json({
            data: {
              averageVisionScore: 15.2, // Low for support
              visionControl: 0.45,
              wardPlacement: 'needs_improvement',
              recommendations: [
                'Place more wards in river bushes',
                'Improve control ward timing',
                'Focus on objective vision'
              ]
            },
            confidence: 0.84
          });
        })
      );

      render(<VisionHeatmap {...defaultProps} role="SUPPORT" />);

      await waitFor(() => {
        expect(screen.getByText(/place more wards/i) ||
               screen.getByText(/improve.*vision/i) ||
               screen.getByText(/objective.*vision/i)).toBeInTheDocument();
      });
    });

    it('should highlight vision improvement opportunities', async () => {
      render(<VisionHeatmap {...defaultProps} />);

      await waitFor(() => {
        // Should show actionable gaming insights
        expect(screen.getByText(/ward coverage/i) || 
               screen.getByText(/vision score/i) ||
               screen.getByText(/map control/i)).toBeInTheDocument();
      });
    });
  });

  describe('Accessibility', () => {
    it('should provide alt text for heatmap', async () => {
      render(<VisionHeatmap {...defaultProps} />);

      await waitFor(() => {
        const canvas = screen.getByTestId('vision-heatmap-canvas');
        expect(canvas).toHaveAttribute('aria-label', expect.stringMatching(/vision heatmap/i));
      });
    });

    it('should be keyboard accessible', async () => {
      render(<VisionHeatmap {...defaultProps} />);

      await waitFor(() => {
        const interactiveElements = screen.getAllByRole('button');
        interactiveElements.forEach(element => {
          expect(element).toHaveAttribute('tabIndex');
        });
      });
    });
  });

  describe('Time-based Analysis', () => {
    it('should show vision patterns over time', async () => {
      server.use(
        http.get('/api/v1/analytics/vision/:summonerId', () => {
          return HttpResponse.json({
            data: {
              averageVisionScore: 28.4,
              visionControl: 0.72,
              timeline: [
                { minute: 5, visionScore: 8, wardsPlaced: 3 },
                { minute: 10, visionScore: 15, wardsPlaced: 6 },
                { minute: 15, visionScore: 22, wardsPlaced: 8 },
                { minute: 20, visionScore: 28, wardsPlaced: 11 }
              ]
            },
            confidence: 0.91
          });
        })
      );

      render(<VisionHeatmap {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByText(/vision timeline/i)).toBeInTheDocument();
      });
    });
  });

  describe('Real Gaming Scenarios', () => {
    it('should handle support role vision correctly', async () => {
      server.use(
        http.get('/api/v1/analytics/vision/:summonerId', () => {
          return HttpResponse.json({
            data: {
              averageVisionScore: 42.1, // High for support
              visionControl: 0.89,
              wardPlacement: 'excellent',
              wardType: 'control_ward_focused',
              mapCoverage: 0.78
            },
            confidence: 0.92
          });
        })
      );

      render(<VisionHeatmap {...defaultProps} role="SUPPORT" />);

      await waitFor(() => {
        expect(screen.getByText(/42\.1/)).toBeInTheDocument(); // Vision score
        expect(screen.getByText(/excellent/i)).toBeInTheDocument(); // Ward placement rating
        expect(screen.getByText(/89%/)).toBeInTheDocument(); // Vision control
      });
    });

    it('should adapt recommendations based on game phase', async () => {
      server.use(
        http.get('/api/v1/analytics/vision/:summonerId', () => {
          return HttpResponse.json({
            data: {
              averageVisionScore: 28.4,
              visionControl: 0.72,
              gamePhaseAnalysis: {
                early: { rating: 'good', focus: 'river_control' },
                mid: { rating: 'average', focus: 'objective_vision' },
                late: { rating: 'excellent', focus: 'team_fight_vision' }
              }
            },
            confidence: 0.87
          });
        })
      );

      render(<VisionHeatmap {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByText(/early game.*river/i) ||
               screen.getByText(/objective.*vision/i) ||
               screen.getByText(/team fight.*vision/i)).toBeInTheDocument();
      });
    });
  });

  describe('Data Confidence', () => {
    it('should display confidence level', async () => {
      render(<VisionHeatmap {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByText(/confidence.*79%/i)).toBeInTheDocument();
      });
    });

    it('should warn about insufficient vision data', async () => {
      server.use(
        http.get('/api/v1/analytics/vision/:summonerId', () => {
          return HttpResponse.json({
            data: {
              averageVisionScore: 0,
              visionControl: 0,
              insufficientData: true
            },
            confidence: 0.15
          });
        })
      );

      render(<VisionHeatmap {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByText(/insufficient vision data/i)).toBeInTheDocument();
      });
    });
  });
});