# 🤝 Contributing to LoL Match Exporter

Thank you for your interest in contributing! This guide will help you get started with development.

## 🚀 Quick Start

### Prerequisites
- Python 3.8+
- Node.js 16+
- Git
- Valid Riot Games API Key

### Development Setup

1. **Fork and Clone**
```bash
git clone https://github.com/yourusername/lol_match_exporter.git
cd lol_match_exporter
```

2. **Backend Setup**
```bash
# Create virtual environment
python -m venv .venv
source .venv/bin/activate  # or .venv\Scripts\activate on Windows

# Install dependencies
pip install -r requirements.txt

# Create .env file
echo "RIOT_API_KEY=your_api_key_here" > .env
echo "EXPORTER_API_KEY=dev_secret_key" >> .env
```

3. **Frontend Setup**
```bash
cd web
npm install
```

4. **Start Development Servers**
```bash
# Terminal 1 - Backend (from root directory)
python server.py

# Terminal 2 - Frontend (from web directory)
npm run dev
```

Visit `http://localhost:5173` for the frontend and `http://localhost:8000` for the API.

## 📁 Project Structure

```
lol_match_exporter/
├── server.py                      # FastAPI server
├── lol_match_exporter.py          # Core export logic
├── riot_api_enhanced.py           # API optimization module
├── requirements.txt               # Python dependencies
├── web/                           # React frontend
│   ├── src/
│   │   ├── NewApp.tsx            # Main app component
│   │   ├── components/           # Reusable UI components
│   │   │   ├── layout/          # Layout components
│   │   │   ├── tables/          # Data table components
│   │   │   └── ExporterMUI.tsx  # Export interface
│   │   ├── views/               # Main application views
│   │   │   ├── Overview.tsx     # Dashboard overview
│   │   │   ├── RolesView.tsx    # Roles analysis
│   │   │   ├── ChampionsView.tsx # Champions analysis
│   │   │   └── ChampionDetails.tsx # Individual champion
│   │   ├── hooks/               # Custom React hooks
│   │   ├── services/            # API and export services
│   │   ├── theme/               # MUI theme configuration
│   │   ├── types/               # TypeScript type definitions
│   │   └── utils/               # Utility functions
│   ├── package.json             # Frontend dependencies
│   └── vite.config.ts           # Vite configuration
├── jobs/                        # Export job files (generated)
├── README.md                    # Main documentation
├── CHANGELOG.md                 # Version history
├── DEPLOYMENT.md                # Deployment guide
└── PERFORMANCE.md               # Performance optimization
```

## 🛠️ Development Guidelines

### Code Style

#### Python
- Follow PEP 8
- Use type hints
- Add docstrings for public functions
- Maximum line length: 100 characters

```python
async def get_match_data(match_id: str, use_cache: bool = True) -> Optional[Dict[str, Any]]:
    """Retrieve match data from Riot API with optional caching.
    
    Args:
        match_id: Riot match identifier
        use_cache: Whether to use cached data if available
        
    Returns:
        Match data dictionary or None if not found
    """
    # Implementation here
```

#### TypeScript/React
- Use TypeScript strict mode
- Prefer functional components with hooks
- Use Material-UI components consistently
- Follow React best practices

```typescript
interface ChampionStatsProps {
  data: Row[];
  selectedChampion?: string;
  onChampionSelect: (champion: string) => void;
}

export const ChampionStats: React.FC<ChampionStatsProps> = ({
  data,
  selectedChampion,
  onChampionSelect
}) => {
  // Implementation here
};
```

### Git Workflow

1. **Create Feature Branch**
```bash
git checkout -b feature/your-feature-name
```

2. **Make Changes**
- Write tests for new functionality
- Update documentation if needed
- Follow commit message conventions

3. **Commit Messages**
```bash
git commit -m "feat: add champion mastery integration"
git commit -m "fix: resolve rate limiting edge case"
git commit -m "docs: update API endpoints documentation"
```

Use prefixes:
- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation updates
- `style:` Code style changes
- `refactor:` Code refactoring
- `test:` Test additions/updates
- `chore:` Maintenance tasks

4. **Push and Create PR**
```bash
git push origin feature/your-feature-name
```

## 🧪 Testing

### Backend Tests
```bash
# Install test dependencies
pip install pytest pytest-asyncio httpx

# Run tests
pytest tests/
```

### Frontend Tests
```bash
cd web
npm test
```

### Manual Testing Checklist
- [ ] Export functionality works with test data
- [ ] Charts render correctly with different data sizes
- [ ] Navigation works smoothly between views
- [ ] Export buttons function properly
- [ ] Theme switching works
- [ ] Mobile responsive design
- [ ] Error handling displays appropriate messages

## 🐛 Bug Reports

When reporting bugs, please include:

1. **Environment Information**
   - OS and version
   - Python version
   - Node.js version
   - Browser (for frontend issues)

2. **Steps to Reproduce**
   - Clear, numbered steps
   - Expected vs actual behavior
   - Screenshots if applicable

3. **Error Messages**
   - Full error messages
   - Browser console errors
   - Server logs if relevant

## 💡 Feature Requests

For new features:

1. **Check Existing Issues** - Make sure it hasn't been requested
2. **Describe the Problem** - What problem does this solve?
3. **Propose a Solution** - How should it work?
4. **Consider Alternatives** - Are there other ways to solve this?

## 📋 Areas for Contribution

### High Priority
- [ ] Additional chart types (Heatmaps, Sankey diagrams)
- [ ] Advanced filtering options (date ranges, specific champions)
- [ ] Performance optimizations for large datasets
- [ ] Mobile app version (React Native)
- [ ] Real-time match tracking

### Medium Priority
- [ ] Multi-language support (i18n)
- [ ] User preferences/settings persistence
- [ ] Data export to more formats (PDF, JSON)
- [ ] Integration with other Riot games
- [ ] Team analysis features

### Low Priority
- [ ] Offline mode support
- [ ] Data visualization improvements
- [ ] Accessibility enhancements
- [ ] Documentation improvements
- [ ] Test coverage improvements

## 🔧 Development Tips

### Hot Reload Setup
Both frontend and backend support hot reload:
- Frontend: Vite automatically reloads on changes
- Backend: Use `uvicorn --reload` for auto-restart

### Debugging
```python
# Backend debugging
import logging
logging.basicConfig(level=logging.DEBUG)

# Or use debugger
import pdb; pdb.set_trace()
```

```typescript
// Frontend debugging
console.log('Debug data:', data);

// React Developer Tools
// Chrome extension for component inspection
```

### API Testing
```bash
# Test API endpoints directly
curl -X POST http://localhost:8000/export \
  -H "Content-Type: application/json" \
  -d '{"platform":"euw1","riotId":"test#EUW","count":10}'
```

### Database for Development
For testing with larger datasets, consider using a local database:

```python
# Optional: SQLite for job persistence
import sqlite3

def init_db():
    conn = sqlite3.connect('jobs.db')
    conn.execute('''
        CREATE TABLE IF NOT EXISTS jobs (
            id TEXT PRIMARY KEY,
            created REAL,
            status TEXT,
            data TEXT
        )
    ''')
    conn.close()
```

## 📚 Learning Resources

### Material-UI
- [MUI Documentation](https://mui.com/)
- [MUI Component Examples](https://mui.com/components/)

### React & TypeScript
- [React Documentation](https://react.dev/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)

### Recharts
- [Recharts Documentation](https://recharts.org/)
- [Chart Examples](https://recharts.org/en-US/examples)

### FastAPI
- [FastAPI Documentation](https://fastapi.tiangolo.com/)
- [Python Type Hints](https://docs.python.org/3/library/typing.html)

## 🎉 Recognition

Contributors will be:
- Listed in the project README
- Credited in release notes
- Given appropriate GitHub repository permissions

## 📞 Getting Help

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and general discussion
- **Discord**: [Community server link] (if available)

Thank you for contributing to LoL Match Exporter! 🚀