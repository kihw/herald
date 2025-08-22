# Configuration OAuth Google pour Herald.lol

## 1. Créer un projet Google Cloud

1. Aller sur [Google Cloud Console](https://console.cloud.google.com/)
2. Créer un nouveau projet ou sélectionner un projet existant
3. Activer l'API Google+ (ou Google Identity)

## 2. Configurer les identifiants OAuth 2.0

1. Dans le menu de navigation, aller à "APIs & Services" > "Credentials"
2. Cliquer sur "Create credentials" > "OAuth 2.0 Client IDs"
3. Sélectionner "Web application" comme type d'application

### Configuration des URI autorisées :

**JavaScript origins autorisées:**
```
https://herald.lol
http://localhost:3000  (pour le développement local)
```

**URI de redirection autorisées:**
```
https://herald.lol/auth/google/callback
http://localhost:3000/auth/google/callback  (pour le développement local)
```

## 3. Configuration des variables d'environnement

Copier le fichier `.env.example` vers `.env` et remplir les valeurs :

```bash
cp .env.example .env
```

Éditer `.env` avec vos identifiants Google :

```bash
GOOGLE_CLIENT_ID=123456789-abcdefghijklmnopqrstuvwxyz.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-your_client_secret_here
GOOGLE_REDIRECT_URL=https://herald.lol/auth/google/callback
```

## 4. Redémarrer les services

Après avoir configuré les variables d'environnement :

```bash
# Arrêter les services
docker-compose -f docker-compose.production.yml down

# Redémarrer avec les nouvelles variables
docker-compose -f docker-compose.production.yml up -d
```

## 5. Test de l'authentification

1. Aller sur https://herald.lol
2. Cliquer sur "Se connecter avec Google"
3. Autoriser l'application
4. Vous devriez être redirigé vers le dashboard avec vos informations utilisateur

## 6. Dépannage

### Erreur "redirect_uri_mismatch"
- Vérifier que l'URI de redirection dans Google Cloud Console correspond exactement à celle configurée
- S'assurer que HTTPS est utilisé en production

### Erreur "invalid_client"
- Vérifier que GOOGLE_CLIENT_ID et GOOGLE_CLIENT_SECRET sont corrects
- S'assurer que l'API Google+ est activée

### Mode développement (sans OAuth configuré)
- L'application fonctionne en mode mock si les variables OAuth ne sont pas configurées
- Un utilisateur de test sera créé automatiquement

## 7. Sécurité en production

- Ne jamais exposer GOOGLE_CLIENT_SECRET dans le code frontend
- Utiliser HTTPS en production
- Configurer un domaine de confiance dans Google Cloud Console
- Implémenter une gestion de session sécurisée (JWT, Redis, etc.)

## 8. Fonctionnalités OAuth implémentées

✅ Initiation du flow OAuth Google  
✅ Gestion du callback et échange de code  
✅ Récupération des informations utilisateur  
✅ Gestion des erreurs OAuth  
✅ Mode fallback mock pour le développement  
✅ Interface utilisateur intégrée  
✅ Persistance locale de la session  
✅ Déconnexion utilisateur  

## 9. Améliorations possibles

- [ ] Gestion des refresh tokens
- [ ] Stockage sécurisé des sessions (Redis/Base de données)
- [ ] API de gestion des utilisateurs
- [ ] Liaison avec les comptes Riot Games
- [ ] Rôles et permissions utilisateur