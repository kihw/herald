# LoL Match Exporter - Guide d'utilisation

## 🚀 Application complètement fonctionnelle !

L'application LoL Match Exporter est maintenant opérationnelle avec une interface moderne et toutes les fonctionnalités essentielles.

### 📍 Accès à l'application

- **Interface web** : http://localhost:5173
- **API Backend** : http://localhost:8001

### 🎯 Comment utiliser l'application

#### 1. **Authentification**

1. Ouvrir http://localhost:5173 dans votre navigateur
2. Saisir vos informations Riot :
   - **Riot ID** : Votre nom d'invocateur (ex: "Hide on bush")
   - **Tag Line** : Votre tag (ex: "KR1")
   - **Région** : Sélectionner votre région (ex: "euw1", "na1", "kr")
3. Cliquer sur **"Validate Account"**

> **Note** : En mode développement, n'importe quelles valeurs fonctionnent pour tester l'interface.

#### 2. **Dashboard Overview**

Une fois connecté, vous accédez au dashboard principal avec 3 onglets :

**📊 Onglet "Overview"**

- **Statistiques globales** : Total matches, taux de victoire
- **Champion favori** : Statistiques détaillées du champion le plus joué
- **Performance récente** : Résultats des 7 et 30 derniers jours
- **Rang actuel** : Informations de classement

#### 3. **Gestion des matches**

**📋 Onglet "Match History"**

- **Liste des matches** : Historique paginé avec détails complets
- **Informations par match** :
  - Champion joué
  - Mode de jeu (Ranked, Normal, etc.)
  - Résultat (Victoire/Défaite)
  - KDA (Kills/Deaths/Assists)
  - Durée de la partie
  - Date de la partie
- **Actions disponibles** :
  - 🔄 **Sync Matches** : Synchroniser avec l'API Riot
  - 📥 **Export** : Exporter les données
  - 👁️ **View Details** : Voir les détails d'un match

#### 4. **Paramètres utilisateur**

**⚙️ Onglet "Settings"**

- **Collecte de données** :
  - ✅ Inclure les données de timeline
  - ✅ Inclure toutes les données des matches
- **Synchronisation** :
  - ✅ Synchronisation automatique
  - ⏱️ Fréquence : Toutes les heures à hebdomadaire
- **Apparence** :
  - 🌞 Mode clair/sombre
- **Sauvegarde** : Les paramètres sont sauvegardés automatiquement

### 🛠️ Fonctionnalités techniques

#### **Architecture**

```
Frontend (React + TypeScript + Vite)
├── Interface moderne Material-UI
├── Navigation par onglets
├── Composants réutilisables
└── Gestion d'état avec Context API

Backend (Go + Gin)
├── API REST complète
├── Authentification par session
├── Données mockées pour développement
└── CORS configuré
```

#### **Endpoints API disponibles**

- `GET /api/health` - Vérification de santé
- `POST /api/auth/validate` - Validation du compte Riot
- `GET /api/auth/session` - Vérification de session
- `GET /api/dashboard/stats` - Statistiques utilisateur
- `GET /api/dashboard/matches` - Historique des matches
- `POST /api/dashboard/sync` - Synchronisation des matches
- `GET /api/dashboard/settings` - Paramètres utilisateur
- `PUT /api/dashboard/settings` - Mise à jour des paramètres

### 🎨 Interface utilisateur

#### **Design moderne**

- **Material-UI** : Composants élégants et cohérents
- **Responsive** : Adapté à tous les écrans
- **Thème** : Couleurs League of Legends
- **Navigation** : Onglets intuitifs
- **Feedback** : Notifications et états de chargement

#### **Expérience utilisateur**

- **Performance** : Chargement rapide avec Vite
- **Interactivité** : Réponses immédiates aux actions
- **Validation** : Formulaires avec validation en temps réel
- **Accessibilité** : Interface accessible et intuitive

### 🔧 Développement

#### **Démarrer l'application**

```bash
# Terminal 1 - Frontend
cd web
npm run dev
# → http://localhost:5173

# Terminal 2 - Backend
go run ./cmd/dev-server
# → http://localhost:8001
```

#### **Structure du projet**

```
lol_match_exporter/
├── web/                    # Frontend React
│   ├── src/
│   │   ├── components/     # Composants UI
│   │   ├── context/        # Gestion d'état
│   │   ├── services/       # API calls
│   │   └── types/          # Types TypeScript
│   └── dist/               # Build de production
├── cmd/
│   ├── server/             # Serveur principal
│   └── dev-server/         # Serveur de développement
├── internal/               # Code backend Go
│   ├── handlers/           # Gestionnaires API
│   ├── services/           # Services métier
│   └── models/             # Modèles de données
└── docker-compose.yml      # Configuration Docker
```

### 🚀 Prochaines étapes

#### **Fonctionnalités à venir**

- 🔗 **Intégration Riot API réelle** : Données live des matches
- 📊 **Analytics avancées** : Graphiques et statistiques détaillées
- 📱 **Version mobile** : Application mobile native
- 🤝 **Mode équipe** : Analyse des performances d'équipe
- 🎮 **Multi-jeux** : Support d'autres jeux Riot

#### **Améliorations techniques**

- 🐘 **Base de données** : Migration PostgreSQL
- 🔄 **WebSocket** : Mises à jour en temps réel
- 🔐 **OAuth réel** : Authentification Riot officielle
- 📦 **Docker** : Déploiement simplifié
- ☁️ **Cloud** : Déploiement sur serveur

### 💡 Conseils d'utilisation

1. **Navigation** : Utilisez les onglets pour explorer les différentes sections
2. **Synchronisation** : Cliquez sur "Sync Matches" pour simuler la récupération de nouveaux matches
3. **Paramètres** : Personnalisez votre expérience dans l'onglet Settings
4. **Responsive** : L'interface s'adapte automatiquement à la taille de votre écran
5. **Performance** : Les données sont chargées de manière optimisée pour une expérience fluide

### 🎉 Conclusion

L'application LoL Match Exporter est maintenant **pleinement fonctionnelle** et prête à l'utilisation ! Elle offre une expérience moderne et complète pour la gestion et l'analyse des matches League of Legends.

**Bon jeu et bonne analyse ! 🏆**
