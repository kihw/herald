# Changelog

All notable changes to LoL Match Exporter will be documented in this file.

## [2.1.0] - 2025-01-16

### 🎨 Major UI Overhaul

#### Added
- **Material-UI v5 Integration**
  - Complete redesign with professional League of Legends theme
  - Custom color palette (gold/blue) matching LoL aesthetics
  - Responsive layout with mobile-first design
  - Dark/light theme with system preference detection

- **Hierarchical Navigation System**
  - Drill-down navigation: Overview → Roles → Champions → Details
  - Breadcrumb navigation for easy orientation
  - Mobile-responsive drawer menu
  - Contextual filtering at each level

- **Advanced Data Visualization**
  - Multiple chart types per view (Bar, Pie, Radar, Scatter, Line)
  - Parametrable charts with dynamic data
  - Responsive containers for all screen sizes
  - Interactive tooltips and legends

- **Comprehensive Export System**
  - PNG export using html2canvas
  - Excel export with xlsx library
  - Combined export (charts + tables)
  - Global export button in AppBar
  - Per-view export options

- **Riot API Optimizations**
  - Enhanced rate limiter with adaptive backoff
  - In-memory LRU cache with TTL support
  - Seasonal data segmentation (2021-2024)
  - Request priority queue system
  - Automatic retry with exponential backoff
  - Cache control toggle in UI
  - Season filtering in export options

#### Changed
- Migrated from basic HTML/CSS to Material-UI components
- Replaced Dashboard.tsx with modular view components
- Improved TypeScript typing throughout the application
- Enhanced accessibility with ARIA labels
- Optimized bundle size with code splitting

#### Removed
- Legacy file import functionality (CSV upload)
- Old UI components (ui/App.tsx, ui/Dashboard.tsx)
- Unused CSS files
- Dead code from previous versions

### 🚀 Performance Improvements
- Reduced API calls by up to 40% with intelligent caching
- Faster data loading with memoized calculations
- Optimized chart rendering with React.memo
- Efficient table virtualization for large datasets

### ♿ Accessibility
- Added ARIA labels to all interactive elements
- Keyboard navigation support
- Screen reader compatibility
- High contrast mode support

### 📦 Dependencies Updated
- React 18
- Material-UI v5
- TypeScript 4.9+
- Recharts (latest)
- html2canvas (latest)
- xlsx (latest)

## [2.0.0] - 2024-12-15

### Added
- FastAPI server for web interface
- React dashboard with Vite
- Server-sent events (SSE) for real-time logs
- Automatic CSV loading after export
- Light CSV generation option
- ZIP bundle creation for all export files

### Changed
- Updated for Riot API 2024-2025 changes
- Improved rate limiting with backoff
- Enhanced error handling

## [1.0.0] - 2024-01-01

### Initial Release
- Basic CLI tool for match export
- CSV and Parquet output formats
- Support for multiple regions
- Queue filtering
- Data Dragon integration for champion names