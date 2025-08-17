# LoL Match Exporter - Documentation Unifiée

## Vue d'Ensemble 🎯

**LoL Match Exporter** est une plateforme d'analytics League of Legends complète utilisant l'API Riot Games v5 pour fournir des insights approfondis sur les performances des joueurs.

### Fonctionnalités Principales

- 🗄️ **Persistance Complète**: Base de données SQLite avec sauvegarde automatique
- 🚀 **Performance**: Cache intelligent avec amélioration 3x
- 🤖 **IA Avancée**: Recommandations personnalisées avec scoring de confiance
- ⚡ **Temps Réel**: WebSocket pour notifications instantanées
- 📊 **Analytics**: Statistiques détaillées et tendances de performance
- 🧪 **Tests**: Suite complète avec benchmarks et stress tests
- 📈 **Monitoring**: Métriques système et santé en temps réel

## Architecture Technique 🏗️

### Stack Technologique
- **Backend**: Go 1.23+ avec Gin framework
- **Database**: SQLite avec modernc.org/sqlite (pure Go)
- **Frontend**: React 18 + TypeScript + Vite
- **Cache**: Système intelligent en mémoire avec TTL
- **WebSocket**: Gorilla WebSocket pour temps réel
- **API**: Riot Games API v5 (100% authentique)

### Structure du Projet
```
lol_match_exporter/
├── cmd/real-server/          # Serveur principal Go
│   ├── main.go              # Point d'entrée principal
│   ├── database.go          # Couche persistance SQLite
│   ├── cache.go             # Système cache intelligent
│   ├── websocket.go         # Communication temps réel
│   ├── monitoring.go        # Métriques et monitoring
│   └── testing.go           # Suite de tests automatisés
├── web/                     # Frontend React
│   ├── src/components/      # Composants UI
│   ├── src/services/        # Services API
│   └── src/types/           # Types TypeScript
├── internal/                # Services Go internes
│   ├── services/            # Logique métier
│   └── models/              # Modèles de données
└── docs/                    # Documentation
```

## Installation et Configuration 🚀

### Prérequis
- **Go 1.23+**
- **Node.js 18+**
- **Clé API Riot Games** ([Obtenir ici](https://developer.riotgames.com/))

### Installation Rapide
```bash
# 1. Cloner le repository
git clone <repository-url>
cd lol_match_exporter

# 2. Configuration API Riot
cp config.example.env .env
# Éditer .env avec votre clé API Riot

# 3. Build et lancement du backend
cd cmd/real-server
go build -o server.exe .
PORT=8004 ./server.exe

# 4. Lancement du frontend (nouveau terminal)
cd web
npm install
npm run dev
```

### Accès à l'Application
- **Interface Web**: http://localhost:5173
- **API Backend**: http://localhost:8004/api
- **Health Check**: http://localhost:8004/api/system/health

## Guide d'Utilisation 📋

### 1. Authentification
1. Ouvrir l'interface web
2. Entrer votre Riot ID (ex: "Hide on bush#KR1")
3. Sélectionner votre région
4. Valider avec l'API Riot

### 2. Synchronisation des Matchs
1. Cliquer sur "Synchroniser les matchs"
2. Suivre le progress en temps réel via WebSocket
3. Les matchs sont sauvegardés automatiquement en base

### 3. Analytics et Insights
- **Dashboard**: Vue d'ensemble des performances
- **Champions**: Statistiques par champion joué
- **Tendances**: Évolution des performances dans le temps
- **Recommandations IA**: Suggestions personnalisées

### 4. Export de Données
- **CSV**: Export classique pour Excel/Sheets
- **JSON**: Format structuré pour développeurs
- **Parquet**: Format optimisé pour big data

## API Reference 📡

### Endpoints Principaux
```http
# Authentification
POST /api/auth/validate
Content-Type: application/json
{
  "riot_id": "Hide on bush",
  "riot_tag": "KR1",
  "region": "kr"
}

# Synchronisation matchs
POST /api/matches/sync
Cookie: lol-session=...

# Statistiques utilisateur
GET /api/stats/dashboard
Cookie: lol-session=...

# Recommandations IA
GET /api/ai/recommendations
Cookie: lol-session=...

# Métriques système
GET /api/system/metrics
Cookie: lol-session=...
```

### WebSocket Events
```javascript
// Connexion WebSocket
const ws = new WebSocket('ws://localhost:8004/api/ws');

// Événements disponibles
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  switch(data.type) {
    case 'match_sync':     // Nouveau match synchronisé
    case 'stats_update':   // Statistiques mises à jour
    case 'system_update':  // Mise à jour système
  }
};
```

## Performance et Monitoring 📊

### Métriques Système
Le système collecte automatiquement:
- **Mémoire**: Utilisation RAM et garbage collection
- **Performance**: Latence des requêtes et cache hit ratio
- **Database**: Nombre de matchs, utilisateurs, performance queries
- **WebSocket**: Connexions actives, messages échangés
- **Santé**: Score global et détection d'anomalies

### Endpoints de Monitoring
```http
GET /api/system/health         # État de santé global
GET /api/system/metrics        # Métriques détaillées
GET /api/system/cache          # Performance cache
GET /api/system/websocket      # Statistiques WebSocket
```

## Tests et Validation 🧪

### Lancement des Tests
```bash
# Tests automatisés via API
curl -X POST http://localhost:8004/api/system/test/run
curl http://localhost:8004/api/system/test/results

# Stress test
curl -X POST http://localhost:8004/api/system/test/stress \
  -H "Content-Type: application/json" \
  -d '{"duration_seconds": 30, "concurrent_requests": 10}'
```

### Types de Tests
- **Tests Unitaires**: Database, cache, WebSocket
- **Tests d'Intégration**: Flux complets end-to-end
- **Tests de Performance**: Latence, memory, stress
- **Tests de Santé**: Monitoring et alerting

## Déploiement Production 🌐

### Docker Deployment
```bash
# Build de l'image
docker build -t lol-exporter .

# Lancement avec docker-compose
docker-compose -f docker-compose.prod.yml up -d
```

### Variables d'Environnement
```env
RIOT_API_KEY=RGAPI-your-key-here
GIN_MODE=release
PORT=8004
SESSION_SECRET=your-production-secret
DATABASE_PATH=./production.db
```

## Troubleshooting 🔧

### Problèmes Courants

**Erreur "Aucun match trouvé"**
- Vérifier la clé API Riot
- Confirmer la région sélectionnée
- Vérifier l'historique récent de matchs

**Performance lente**
- Vérifier le cache hit ratio (`/api/system/cache`)
- Monitorer la mémoire (`/api/system/metrics`)
- Consulter les logs de performance

**Erreurs WebSocket**
- Vérifier les CORS autorisés
- Confirmer les ports ouverts
- Tester la connectivité réseau

### Logs et Debugging
```bash
# Logs détaillés
GIN_MODE=debug ./server.exe

# Monitoring en temps réel
curl -s http://localhost:8004/api/system/health | jq .
```

## Contribution 🤝

### Standards de Code
- **Go**: gofmt, golint, vet
- **TypeScript**: ESLint, Prettier
- **Tests**: Coverage >90%
- **Documentation**: Commentaires GoDoc

### Workflow de Développement
1. Fork du repository
2. Branche feature (`git checkout -b feature/nom-feature`)
3. Tests et validation
4. Pull Request avec description détaillée

## Roadmap et TODO 🗺️

Voir [TODO.md](../TODO.md) pour les fonctionnalités planifiées:
- Meta-game analytics avancées
- Prédictions ML de performance
- Export multi-formats
- Dashboard temps réel amélioré
- Déploiement Kubernetes

## Support et Communauté 💬

- **Issues**: [GitHub Issues](repository-url/issues)
- **Documentation**: [Wiki](repository-url/wiki)
- **Discord**: [Serveur communauté](discord-link)

---

**Version**: 2.0.0 (Production Ready)
**Dernière mise à jour**: 2025-08-17
**Statut**: ✅ Stable - Prêt pour production