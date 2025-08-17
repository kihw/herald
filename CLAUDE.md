# Claude Code Context

## Project: LoL Match Exporter

### Overview
Application hybride Go/Python/TypeScript pour l'export de données League of Legends avec interface web dashboard.

### Architecture
- **Analytics Backend**: Go natif pour les moteurs d'analytics (MMR, recommandations, statistiques)
- **Web Backend**: Python + FastAPI pour l'API REST et l'orchestration des jobs
- **Frontend**: React + TypeScript + Vite pour l'interface utilisateur
- **Data**: Export CSV/Parquet des statistiques de matchs via l'API Riot Games

### ⚡ Go Native Analytics Migration
Les moteurs d'analytics ont été migrés de Python vers Go pour de meilleures performances :
- `analytics_engine.py` → `internal/services/analytics_engine_service.go`
- `mmr_calculator.py` → `internal/services/mmr_calculation_service.go`
- `recommendation_engine.py` → `internal/services/recommendation_engine_service.go`

### Structure des dossiers
```
lol_match_exporter/
├── main.go                 # Serveur Go principal avec analytics natifs
├── internal/               # Services Go natifs
│   ├── services/
│   │   ├── analytics_engine_service.go      # Analytics natif Go
│   │   ├── mmr_calculation_service.go       # MMR Calculator natif Go
│   │   ├── recommendation_engine_service.go # Recommendations natives Go
│   │   └── analytics_service.go             # Interface unifiée
│   ├── models/             # Modèles de données Go
│   └── db/                 # Gestionnaire de base de données Go
├── lol_match_exporter.py   # CLI principal d'export (Python)
├── server.py               # Serveur FastAPI pour l'interface web
├── view_csv.py            # Utilitaire de visualisation CSV
├── requirements.txt       # Dépendances Python (réduites)
├── go.mod                 # Dépendances Go
├── config.example.env     # Configuration d'exemple
├── web/                   # Frontend React
│   ├── src/
│   │   ├── main.tsx
│   │   └── ui/
│   │       ├── App.tsx
│   │       ├── Dashboard.tsx
│   │       ├── Exporter.tsx
│   │       └── util.ts
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.js
├── jobs/                  # Dossier des jobs d'export (généré)
└── ranked_timelines/      # Données de match (généré)
```

### Technologies principales
- **Go 1.23+**: Services natifs pour analytics, MMR, et recommandations
- **Python 3.8+**: pandas, requests, FastAPI, uvicorn, rich (web interface)
- **Node.js/npm**: React 18, TypeScript, Vite, Recharts
- **Database**: PostgreSQL/SQLite avec gestionnaire Go
- **API**: Riot Games API v5 (League of Legends)

### Scripts de développement

#### Go Analytics Backend
```bash
# Construire le serveur Go avec analytics natifs
go build -o analytics-server.exe main.go

# Lancer le serveur Go
./analytics-server.exe  # http://localhost:8001

# Tests Go natifs
go test -v ./internal/services
```

#### Python Web Backend
```bash
# Installer les dépendances
pip install -r requirements.txt

# Lancer le serveur FastAPI
python server.py  # http://localhost:8000
```

#### Frontend
```bash
cd web

# Installer les dépendances
npm install

# Développement
npm run dev       # http://localhost:5173

# Production
npm run build
npm run preview

# Qualité de code
npm run lint
npm run type-check
```

#### Raccourcis Makefile
```bash
make install      # Installe toutes les dépendances
make dev         # Lance les serveurs de dev
make build       # Build de production
make lint        # Linting Python + TypeScript
make clean       # Nettoie les artifacts
```

### Configuration requise
1. **Clé API Riot Games**: Obtenir sur https://developer.riotgames.com/
2. **Fichier .env**: Copier `config.example.env` vers `.env` et configurer
3. **Python 3.8+** et **Node.js 18+**

### Usage typique
1. Configurer la clé API Riot dans `.env`
2. Lancer le backend: `python server.py`
3. Lancer le frontend: `cd web && npm run dev`
4. Ouvrir http://localhost:5173
5. Entrer un Riot ID et lancer l'export
6. Analyser les données dans le dashboard

### Notes de sécurité
- Les clés API ne doivent jamais être commitées
- Protection optionnelle par clé API locale
- CORS ouvert en développement (à restreindre en production)

### Limites connues
- Rate limiting API Riot géré automatiquement
- Certains modes de jeu ne renseignent pas tous les champs
- Nécessite une connexion Internet pour Data Dragon

### Support
Consulter le README.md pour les détails d'utilisation et les exemples.