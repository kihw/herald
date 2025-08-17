# 🔑 Guide : Configuration de l'API Riot Games

## 📋 Étapes pour obtenir votre clé API

### 1. Créer un compte développeur Riot

1. Allez sur https://developer.riotgames.com/
2. Cliquez sur "SIGN IN" en haut à droite
3. Connectez-vous avec votre compte Riot Games existant
4. Si vous n'avez pas de compte, créez-en un sur https://auth.riotgames.com/

### 2. Générer votre clé API

1. Une fois connecté, vous verrez votre tableau de bord développeur
2. Dans la section "Development API Key", vous verrez votre clé personnelle
3. Copiez cette clé (format : `RGAPI-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`)

### 3. Configurer la clé dans votre environnement

#### Windows (PowerShell) - Session actuelle

```powershell
$env:RIOT_API_KEY="RGAPI-votre-clé-ici"
```

#### Windows (PowerShell) - Permanent (utilisateur)

```powershell
[Environment]::SetEnvironmentVariable("RIOT_API_KEY", "RGAPI-votre-clé-ici", "User")
```

#### Windows (CMD)

```cmd
set RIOT_API_KEY=RGAPI-votre-clé-ici
```

#### Linux/Mac

```bash
export RIOT_API_KEY="RGAPI-votre-clé-ici"
```

#### Fichier .env (recommandé pour le développement)

Créez un fichier `.env` dans le répertoire racine :

```
RIOT_API_KEY=RGAPI-votre-clé-ici
```

## 🚨 Limitations importantes

### Clé de développement personnelle

- **Rate Limit** : 100 requêtes toutes les 2 minutes
- **Expire** : Au bout de 24 heures
- **Usage** : Développement et tests uniquement

### Clé de production

Pour une application en production, vous devez :

1. Créer une "Application" sur le portail développeur
2. Décrire votre projet et son usage
3. Obtenir l'approbation de Riot
4. Recevoir une clé avec des limites plus élevées

## 🧪 Tester votre configuration

### Test rapide avec PowerShell

```powershell
# Vérifier que la clé est configurée
$env:RIOT_API_KEY

# Tester avec notre script
.\test-real-api.ps1
```

### Test manuel avec cURL

```bash
curl -H "X-Riot-Token: VOTRE-CLE-ICI" \
"https://europe.api.riotgames.com/riot/account/v1/accounts/by-riot-id/Hide%20on%20bush/KR1"
```

## 🔧 Configuration pour différents environnements

### Développement local

```powershell
# Dans votre terminal avant de lancer l'app
$env:RIOT_API_KEY="RGAPI-votre-clé-dev"
$env:GIN_MODE="debug"
.\real-server.exe
```

### Production Docker

```yaml
# docker-compose.prod.yml
environment:
  - RIOT_API_KEY=${RIOT_API_KEY}
  - GIN_MODE=release
```

### Variables d'environnement complètes

```bash
# Clé API Riot (OBLIGATOIRE)
RIOT_API_KEY=RGAPI-votre-clé-ici

# Configuration serveur
PORT=8001
GIN_MODE=release
SESSION_SECRET=votre-secret-session

# Base de données (futur)
DB_HOST=localhost
DB_PORT=5432
DB_USER=lol_user
DB_PASSWORD=secure_password
DB_NAME=lol_match_db
```

## ⚠️ Sécurité et bonnes pratiques

### ❌ À NE PAS FAIRE

- Commiter votre clé API dans Git
- Exposer votre clé côté client (JavaScript)
- Partager votre clé de développement
- Utiliser une clé de dev en production

### ✅ À FAIRE

- Utiliser des variables d'environnement
- Ajouter `.env` à votre `.gitconfig`
- Renouveler régulièrement votre clé
- Monitorer votre usage de rate limits

## 🚀 Démarrage rapide

### 1. Configuration express

```powershell
# Configurer la clé API
$env:RIOT_API_KEY="RGAPI-votre-clé-ici"

# Construire et tester
go build -o real-server.exe ./cmd/real-server
.\test-real-api.ps1
```

### 2. Démarrer l'application réelle

```powershell
# Au lieu du dev-server mocké
.\real-server.exe

# L'app sera disponible sur http://localhost:8001
```

## 🔍 Résolution des problèmes

### Erreur "API key not configured"

- Vérifiez : `echo $env:RIOT_API_KEY`
- La clé doit commencer par `RGAPI-`

### Erreur 403 "Forbidden"

- Votre clé a expiré (24h pour les clés dev)
- Générez une nouvelle clé sur le portail

### Erreur 429 "Rate Limited"

- Trop de requêtes (100/2min max)
- Attendez avant de retester
- Implémentez un système de rate limiting

### Erreur 404 "Account not found"

- L'utilisateur n'existe pas dans cette région
- Vérifiez le nom d'invocateur et le tag
- Essayez une autre région

## 📚 Documentation API Riot

- **Portal développeur** : https://developer.riotgames.com/
- **Documentation API** : https://developer.riotgames.com/apis
- **Rate limits** : https://developer.riotgames.com/docs/portal#web-apis_rate-limiting
- **Régions** : https://developer.riotgames.com/docs/lol#routing-values

---

💡 **Tip** : Commencez avec la clé de développement pour tester, puis demandez une clé de production quand votre app est prête !
