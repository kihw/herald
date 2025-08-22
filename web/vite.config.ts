import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          // Vendor chunks
          vendor: ['react', 'react-dom'],
          mui: ['@mui/material', '@mui/icons-material', '@mui/lab'],
          charts: ['chart.js', 'react-chartjs-2'],
          
          // Feature chunks
          auth: [
            './src/components/auth/AuthPage',
            './src/components/auth/GoogleAuth',
            './src/components/auth/RiotValidationForm',
            './src/context/AuthContext',
          ],
          groups: [
            './src/components/groups/GroupManagement',
            './src/components/groups/CreateGroupDialog',
            './src/components/groups/JoinGroupDialog',
            './src/components/groups/GroupDetailsDialog',
            './src/components/groups/ComparisonManager',
            './src/services/groupApi',
          ],
          charts: [
            './src/components/charts/ComparisonCharts',
            './src/components/charts/GroupStatsCharts',
            './src/components/charts/PlayerPerformanceWidget',
            './src/components/charts/ChartConfig',
          ],
          dashboard: [
            './src/components/dashboard/MainDashboard',
            './src/components/dashboard/AnalyticsDashboard',
            './src/components/dashboard/EnhancedDashboard',
          ],
          analytics: [
            './src/components/ChampionsAnalytics',
            './src/components/MMRAnalytics',
            './src/components/ExporterMUI',
          ],
        },
      },
    },
    // Optimize chunk size
    chunkSizeWarningLimit: 1000,
    // Enable source maps for debugging
    sourcemap: false,
    // Minify for production
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true,
      },
    },
  },
  // Optimize dev server
  server: {
    port: 5173,
    host: true,
    hmr: {
      overlay: false,
    },
  },
  // Optimize dependencies
  optimizeDeps: {
    include: [
      'react',
      'react-dom',
      '@mui/material',
      '@mui/icons-material',
      'chart.js',
      'react-chartjs-2',
    ],
    exclude: [
      // Exclude large dependencies that should be loaded separately
    ],
  },
  // Performance hints
  define: {
    // Enable performance profiling in development
    __DEV__: JSON.stringify(process.env.NODE_ENV === 'development'),
  },
});