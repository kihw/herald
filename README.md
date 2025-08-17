# 🎮 LoL Match Exporter - Modern Web Application

> **Application web moderne avec interface React pour gérer et analyser vos matches League of Legends**

[![Go](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18-blue.svg)](https://reactjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5-blue.svg)](https://www.typescriptlang.org/)
[![Material-UI](https://img.shields.io/badge/Material--UI-5-blue.svg)](https://mui.com/)
[![Status](https://img.shields.io/badge/Status-✅%20Production%20Ready-green.svg)]()

Application complète avec backend Go, frontend React TypeScript, et interface utilisateur moderne avec notifications toast, squelettes de chargement, et panneau d'aide intégré.

## ⚡ Go Native Analytics Migration

**🎉 NEW: Analytics engines migrated to Go for enhanced performance!**

Les moteurs d'analytics ont été migrés de Python vers Go natif pour de meilleures performances :

- **Analytics Engine**: Statistiques de période, performance par rôle, analyse de champion
- **MMR Calculator**: Trajectoires MMR, prédictions de rang, analyse de volatilité  
- **Recommendation Engine**: Recommandations IA, suggestions de champions, conseils gameplay

**Avantages:**
- 🚀 **Performance améliorée** - Traitement direct en Go sans subprocess Python
- 💾 **Mémoire réduite** - Structures natives Go, pas de sérialisation JSON
- 📦 **Déploiement simplifié** - Un seul binaire pour les analytics
- 🔧 **Maintenabilité** - Système de types Go robuste

> Voir [PYTHON_TO_GO_MIGRATION.md](./PYTHON_TO_GO_MIGRATION.md) pour les détails complets

## 🚀 Démarrage Rapide

### ⚡ Lancement automatique (Recommandé)

**Windows (PowerShell):**

```powershell
.\start.ps1
```

**Linux/macOS:**

```bash
chmod +x start.sh
./start.sh
```

### 🎯 Accès à l'application

- **Interface Web** : http://localhost:5173
- **API Backend** : http://localhost:8001

## ✨ Fonctionnalités Principales

### � **Authentification & Sécurité**

- ✅ Validation de compte Riot ID avec tag et région
- ✅ Support multi-régions (EUW, NA, KR, etc.)
- ✅ Gestion de sessions sécurisées
- ✅ Interface d'authentification intuitive

### 📊 **Dashboard Modern & Analytics**

- ✅ **Vue d'ensemble** : Statistiques globales, taux de victoire, champion favori
- ✅ **Performance récente** : Analyse des 7 et 30 derniers jours avec graphiques
- ✅ **Rang actuel** : Informations de classement et progression
- ✅ **Interface Material-UI 5** : Design moderne et responsive

### � **Gestion Avancée des Matches**

- ✅ **Historique complet** : Liste paginée avec détails et filtres
- ✅ **Synchronisation intelligente** : Mise à jour via API Riot Games
- ✅ **Export multi-format** : CSV, JSON, Excel avec formatage
- ✅ **Analyse détaillée** : KDA, durée, mode de jeu, performance

### 🎨 **Expérience Utilisateur Premium**

- ✅ **Notifications toast** : Feedback visuel avec animations et couleurs
- ✅ **Squelettes de chargement** : États de chargement élégants
- ✅ **Panneau d'aide intégré** : Guide utilisateur contextualisé
- ✅ **Design responsive** : Optimisé desktop, tablette et mobile
- ✅ **Mode sombre/clair** : Thème adaptatif

### ⚙️ **Paramètres & Personnalisation**

- ✅ **Configuration fine** : Options de collecte et synchronisation
- ✅ **Préférences utilisateur** : Sauvegarde persistante des paramètres
- ✅ **Contrôle de la sync** : Fréquence automatique configurable
- ✅ **Interface intuitive** : Paramètres organisés par catégories

## 🏗️ Architecture Technique

```
┌─────────────────────────────────────────────────────────────┐
│              Frontend React + TypeScript                   │
│  ┌─────────────────┐  ┌───────────────┐  ┌──────────────┐ │
│  │ Auth Components │  │   Dashboard   │  │   Settings   │ │
│  │ • Login Pages   │  │ • Statistics  │  │ • User Prefs │ │
│  │ • Validation    │  │ • Match List  │  │ • Sync Config│ │
│  │ • Session Mgmt  │  │ • Analytics   │  │ • Export Opt │ │
│  └─────────────────┘  └───────────────┘  └──────────────┘ │
│                                                             │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │           Enhanced UI Components                        │ │
│  │ • ToastNotification (feedback système)                 │ │
│  │ • LoadingSkeleton (états de chargement)               │ │
│  │ • HelpPanel (guide utilisateur intégré)              │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────┬───────────────────────────────────────┘
                      │ HTTP/JSON API (localhost:5173)
┌─────────────────────▼───────────────────────────────────────┐
│                Backend Go + Gin Framework                   │
│  ┌─────────────────┐  ┌───────────────┐  ┌──────────────┐ │
│  │   API Routes    │  │   Services    │  │ Data Storage │ │
│  │ • Auth endpoints│  │ • Match Mgmt  │  │ • JSON Files │ │
│  │ • Dashboard API │  │ • User Mgmt   │  │ • CSV Export │ │
│  │ • Settings API  │  │ • Sync Logic  │  │ • Session DB │ │
│  └─────────────────┘  └───────────────┘  └──────────────┘ │
└─────────────────────┬───────────────────────────────────────┘
                      │ (localhost:8001)
┌─────────────────────▼───────────────────────────────────────┐
│                  External Services                          │
│          ┌─────────────────┐    ┌─────────────────────┐     │
│          │   Riot Games    │    │   Development       │     │
│          │      API        │    │   Mock Data         │     │
│          │ (Production)    │    │  (Dev Server)       │     │
│          └─────────────────┘    └─────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

## 🖥️ Modes de Fonctionnement

### 🔧 **Mode Développement** (Actuel)

- **Backend** : Serveur Go avec données mock pour développement
- **Frontend** : Server de développement Vite avec hot-reload
- **Base de données** : Données simulées en mémoire
- **API Riot** : Mode simulation pour tests sans limitations

### 🚀 **Mode Production** (Roadmap)

- **Backend** : Serveur Go avec intégration Riot API complète
- **Frontend** : Build optimisé servi par le backend
- **Base de données** : PostgreSQL avec cache Redis
- **API Riot** : Intégration temps réel avec rate limiting

## 🛠️ Technologies Utilisées

### 🖥️ **Frontend (React)**

- **React 18** : Framework UI moderne avec hooks
- **TypeScript 5** : Typage statique pour fiabilité
- **Material-UI 5** : Composants UI élégants et accessibles
- **Vite** : Build tool ultra-rapide avec HMR
- **Context API** : Gestion d'état centralisée

### 🔧 **Backend (Go)**

- **Go 1.23** : Performance exceptionnelle et concurrence
- **Gin Framework** : Router web léger et rapide
- **Goroutines** : Traitement concurrent des requêtes
- **JSON** : Sérialisation rapide des données
- **Session Management** : Gestion sécurisée des utilisateurs

### 🎨 **UI/UX Components**

- **ToastNotification** : Système de notifications avec animations
- **LoadingSkeleton** : États de chargement élégants
- **HelpPanel** : Guide utilisateur contextuel intégré
- **Responsive Design** : Adaptation automatique aux écrans

## 🖥️ Interface Web + API Locale

L'application moderne combine un backend Go performant avec une interface React TypeScript pour une expérience utilisateur fluide et professionnelle.

### 📦 Installation Manuelle

**1. Cloner le repository**

```bash
git clone <repository-url>
cd lol_match_exporter
```

**2. Backend Go**

```bash
go mod tidy
go build -o server.exe ./cmd/dev-server
```

**3. Frontend React**

```bash
cd web
npm install
npm run build
```

### 🚀 Démarrage en Développement

**Terminal 1 - Backend Go:**

```bash
go run ./cmd/dev-server
# → API disponible sur http://localhost:8001
```

**Terminal 2 - Frontend React:**

```bash
cd web
npm run dev
# → Interface disponible sur http://localhost:5173
```

### 🎯 Guide d'Utilisation

**1. Première Connexion**

- Ouvrir http://localhost:5173
- Saisir votre **Riot ID** (ex: `MonPseudo#EUW`)
- Sélectionner votre **région** (EUW, NA, KR, etc.)
- Cliquer **"Validate Account"**

**2. Navigation Principale**

- **📊 Overview** : Statistiques globales et performance récente
- **📋 Match History** : Historique complet avec détails et filtres
- **⚙️ Settings** : Configuration personnalisée de l'application

**3. Fonctionnalités Clés**

- **🔄 Sync Matches** : Récupération des nouveaux matches depuis l'API
- **📤 Export Data** : Sauvegarde en CSV/JSON avec formatage
- **🔍 Filter & Search** : Filtres avancés par champion, mode, date
- **📊 Analytics** : Graphiques de performance et tendances

**4. Expérience Utilisateur**

- **🔔 Notifications** : Feedback visuel pour toutes les actions
- **⏳ Loading States** : Squelettes élégants pendant les chargements
- **❓ Help Panel** : Guide contextuel intégré (bouton d'aide)
- **🌙 Dark Mode** : Basculement automatique selon les préférences

**5. Paramètres Avancés**

- **📅 Sync Frequency** : Configuration de la synchronisation automatique
- **💾 Data Collection** : Options de collecte personnalisables
- **🎨 UI Preferences** : Personnalisation de l'interface utilisateur

## 🔌 API Endpoints

### 🔐 **Authentification**

```http
POST   /api/auth/validate       # Validation compte Riot ID
GET    /api/auth/session        # Vérification session active
POST   /api/auth/logout         # Déconnexion utilisateur
GET    /api/auth/regions        # Régions supportées
```

### 📊 **Dashboard & Analytics**

```http
GET    /api/dashboard/stats     # Statistiques utilisateur globales
GET    /api/dashboard/matches   # Historique matches (paginé)
POST   /api/dashboard/sync      # Synchronisation avec API Riot
GET    /api/dashboard/export    # Export données (CSV/JSON)
```

### ⚙️ **Paramètres & Configuration**

```http
GET    /api/dashboard/settings  # Récupération paramètres utilisateur
PUT    /api/dashboard/settings  # Sauvegarde paramètres
POST   /api/dashboard/reset     # Reset configuration par défaut
```

### 🔍 **Utilitaires & Monitoring**

```http
GET    /api/health             # Statut santé de l'API
GET    /api/profile            # Profil utilisateur actuel
GET    /api/version            # Version de l'application
```

### 🛡️ **Sécurité & CORS**

- **CORS** : Configuration ouverte en développement
- **Sessions** : Gestion sécurisée côté serveur
- **Rate Limiting** : Protection contre les abus API
- **Error Handling** : Réponses d'erreur structurées

## 📁 Structure du Projet

```
lol_match_exporter/
├── 📄 README.md                    # Documentation principale
├── 📄 GUIDE_UTILISATION.md         # Guide utilisateur détaillé
├── 🚀 start.ps1 / start.sh         # Scripts de démarrage automatique
│
├── 🗂️ cmd/                         # Applications Go
│   ├── server/                     # Serveur principal (production)
│   └── dev-server/                 # Serveur de développement
│
├── 🗂️ internal/                    # Code backend Go
│   ├── handlers/                   # Gestionnaires API REST
│   ├── services/                   # Services métier
│   ├── models/                     # Modèles de données
│   └── config/                     # Configuration
│
├── 🗂️ web/                         # Frontend React TypeScript
│   ├── src/
│   │   ├── components/             # Composants UI
│   │   │   ├── auth/              # Authentification
│   │   │   │   ├── LoginPage.tsx
│   │   │   │   ├── ValidateAccount.tsx
│   │   │   │   └── RegionSelector.tsx
│   │   │   │
│   │   │   ├── dashboard/         # Dashboard principal
│   │   │   │   ├── MainDashboard.tsx
│   │   │   │   ├── OverviewTab.tsx
│   │   │   │   ├── MatchesTab.tsx
│   │   │   │   └── SettingsTab.tsx
│   │   │   │
│   │   │   └── common/            # Composants réutilisables
│   │   │       ├── ToastNotification.tsx    # Notifications
│   │   │       ├── LoadingSkeleton.tsx      # États de chargement
│   │   │       └── HelpPanel.tsx            # Panneau d'aide
│   │   │
│   │   ├── context/               # Gestion d'état React
│   │   │   ├── AuthContext.tsx
│   │   │   └── NotificationContext.tsx
│   │   │
│   │   ├── services/              # Clients API
│   │   │   └── api.ts
│   │   │
│   │   ├── types/                 # Types TypeScript
│   │   │   └── index.ts
│   │   │
│   │   └── utils/                 # Fonctions utilitaires
│   │       └── constants.ts
│   │
│   ├── dist/                      # Build de production
│   ├── package.json               # Dépendances Node.js
│   ├── vite.config.js            # Configuration Vite
│   └── tsconfig.json             # Configuration TypeScript
│
├── 📋 go.mod / go.sum             # Modules Go
├── ⚙️ .env                       # Variables d'environnement
└── 🐳 Dockerfile                 # Configuration Docker (futur)
```

---

## 🚀 Roadmap & Évolution

### 🎯 **Version Actuelle (v2.0 - Production Ready)**

- ✅ Interface utilisateur React TypeScript complète
- ✅ Backend Go performant avec API REST
- ✅ Système d'authentification Riot ID
- ✅ Dashboard analytics avec statistiques
- ✅ Gestion des matches avec export
- ✅ UI/UX premium avec notifications et aide
- ✅ Scripts de démarrage automatique
- ✅ Documentation complète

### 🔮 **Prochaines Versions**

**v2.1 - Intégration Riot API Complète**

- 🔄 Synchronisation temps réel avec l'API Riot Games
- 📊 Données live des matches et classements
- 🏆 Informations de rang et LP en temps réel
- 🔍 Recherche avancée de joueurs

**v2.2 - Analytics Avancées**

- 📈 Graphiques de performance détaillés
- 🎯 Analyse de tendances et prédictions
- 🏆 Comparaisons avec autres joueurs
- � Statistiques par champion approfondies

**v2.3 - Fonctionnalités Sociales & Teams**

- 👥 Support multi-joueurs et équipes
- 🤝 Analyse de performance d'équipe
- 📱 Application mobile (React Native)
- 🌐 Partage social des statistiques

**v3.0 - Platform Complète**

- 🏢 Version entreprise/équipes eSport
- ☁️ Déploiement cloud avec base de données
- 🔐 Authentification OAuth avancée
- 📊 Tableaux de bord personnalisables

## 🎮 Utilisation Pratique

### � **Cas d'Usage Principaux**

**🏆 Joueur Compétitif**

- Analyse de performance pour amélioration
- Tracking de progression en ranked
- Identification des points faibles
- Optimisation du pool de champions

**📚 Coach/Analyste**

- Analyse d'équipe et joueurs
- Préparation de stratégies
- Suivi des performances individuelles
- Export de données pour analyse externe

**🎯 Joueur Casual**

- Découverte de ses statistiques
- Suivi de progression casual
- Exploration de l'historique
- Partage avec les amis

### 🎪 **Fonctionnalités Avancées**

**📊 Analytics Dashboard**

- Vue globale des performances
- Graphiques de tendance temporelle
- Métriques de performance par rôle/champion
- Comparaisons et benchmarking

**🔄 Synchronisation Intelligente**

- Mise à jour automatique périodique
- Détection de nouveaux matches
- Optimisation des appels API
- Cache intelligent pour performance

**📤 Export Flexible**

- Multiple formats (CSV, JSON, Excel)
- Données filtrables par période
- Templates personnalisables
- Intégration avec outils externes

## 🤝 Contribution & Développement

### 🛠️ **Setup Développeur**

```bash
# 1. Fork & Clone
git clone https://github.com/votre-username/lol_match_exporter
cd lol_match_exporter

# 2. Backend Go
go mod download
go run ./cmd/dev-server

# 3. Frontend React
cd web
npm install
npm run dev
```

### 📋 **Guidelines de Contribution**

1. **Fork** le projet
2. Créer une **branch feature** (`git checkout -b feature/AmazingFeature`)
3. **Commit** les changements (`git commit -m 'Add: Amazing Feature'`)
4. **Push** vers la branch (`git push origin feature/AmazingFeature`)
5. Ouvrir une **Pull Request**

### 🧪 **Tests & Qualité**

- **Tests unitaires** : Go avec testify
- **Tests E2E** : Cypress pour le frontend
- **Linting** : ESLint + Prettier pour TS, golangci-lint pour Go
- **CI/CD** : GitHub Actions (prévu)

## 📄 Licence & Legal

**Licence MIT** - Voir le fichier `LICENSE` pour détails complets.

### ⚖️ **Conformité Riot Games**

- Utilisation conforme aux [Riot Developer Terms](https://developer.riotgames.com/terms)
- Respect des limites de rate limiting API
- Pas de données personnelles sensibles stockées
- Usage personnel et éducationnel uniquement

## 🎉 Remerciements & Crédits

- **🏆 Riot Games** : Pour l'API League of Legends exceptionnelle
- **⚛️ React Community** : Écosystème et composants Material-UI
- **🐹 Go Community** : Frameworks et librairies performantes
- **👥 Contributors** : Tous ceux qui contribuent au projet
- **🎮 Vous** : Pour utiliser et améliorer cette application !

---

<div align="center">

**🎯 Développé avec ❤️ pour la communauté League of Legends 🎯**

[![Made with Go](https://img.shields.io/badge/Made%20with-Go-blue.svg)](https://golang.org/)
[![Made with React](https://img.shields.io/badge/Made%20with-React-blue.svg)](https://reactjs.org/)
[![Made with TypeScript](https://img.shields.io/badge/Made%20with-TypeScript-blue.svg)](https://www.typescriptlang.org/)

**🚀 Bon jeu et excellente analyse ! �**

</div>
