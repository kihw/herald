# Architecture LoL Match Manager

## 🎯 Vision

Transformation d'un exporteur de matchs vers une plateforme complète de gestion de données LoL avec authentification, synchronisation automatique et interface utilisateur moderne.

## 📊 Architecture Technique

### Backend (Go)

```
/cmd/
  /server/     # Point d'entrée serveur
/internal/
  /auth/       # JWT, bcrypt, sessions
  /models/     # Structures de données
  /handlers/   # Controllers HTTP
  /services/   # Logique métier
  /db/         # Accès base de données
  /scheduler/  # Tâches cron
  /riot/       # Client API Riot
/pkg/
  /utils/      # Utilitaires partagés
/migrations/   # Scripts SQL
```

### Frontend (React + TypeScript)

```
/src/
  /components/  # Composants réutilisables
  /pages/      # Pages principales
    - Login.tsx
    - Register.tsx
    - Dashboard.tsx
    - Profile.tsx
  /hooks/      # Custom hooks
  /services/   # API calls
  /store/      # State management
  /guards/     # Route protection
```

### Base de Données (PostgreSQL)

```sql
-- Tables principales
users
user_settings
matches
match_participants
sync_jobs
system_config
```

## 🔐 Sécurité

- JWT tokens avec refresh
- Bcrypt pour passwords
- Rate limiting
- CORS configuré
- Sessions sécurisées

## 🔄 Synchronisation

- Cron job quotidien 00:00
- Sync incrémentale par utilisateur
- Gestion des rate limits Riot
- Retry automatique avec backoff

## 📱 UX/UI

- Dashboard moderne avec métriques
- Configuration profil intuitive
- Actualisation manuelle (cooldown 2min)
- Notifications temps réel
- Responsive design
