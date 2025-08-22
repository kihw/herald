# 🚀 Guide de Déploiement - Herald.lol

## 📋 Pré-requis de Production

### Infrastructure Requise
- **Serveur** : Linux (Ubuntu 20.04+ recommandé)
- **RAM** : Minimum 2GB, recommandé 4GB
- **CPU** : 2 cores minimum
- **Stockage** : 20GB SSD
- **Bande passante** : Illimitée
- **Domaine** : herald-lol.com (ou similaire)

### Services Externes
- **Google OAuth** : Clés API configurées
- **Riot Games API** : Clé développeur active
- **SSL Certificate** : Let's Encrypt ou commercial
- **Monitoring** : Optionnel (Grafana, Prometheus)

---

## 🔧 Configuration de Production

### 1. Variables d'Environnement
Créer `/home/debian/herald/.env.production` :

```bash
# Base de données
DB_PATH="/data/herald.db"
DB_DRIVER="sqlite3"

# Authentification
GOOGLE_CLIENT_ID="your_google_client_id"
GOOGLE_CLIENT_SECRET="your_google_client_secret"
GOOGLE_REDIRECT_URL="https://herald-lol.com/api/auth/callback"

# Riot API
RIOT_API_KEY="your_riot_api_key"
RIOT_API_BASE_URL="https://euw1.api.riotgames.com"

# JWT
JWT_SECRET="your_super_secure_jwt_secret_key_here"
JWT_EXPIRES="24h"

# Serveur
PORT="8000"
ENV="production"
HOST="0.0.0.0"

# CORS
ALLOWED_ORIGINS="https://herald-lol.com"

# Cache
REDIS_URL="redis://localhost:6379"
CACHE_TTL="3600"

# Logging
LOG_LEVEL="info"
LOG_FILE="/var/log/herald/app.log"
```

### 2. Configuration Docker Production
Mettre à jour `docker-compose.production.yml` :

```yaml
version: '3.8'

services:
  herald-app:
    build:
      context: .
      dockerfile: Dockerfile
      target: production
    container_name: herald-app
    restart: unless-stopped
    ports:
      - "8000:8000"
    volumes:
      - herald-data:/data
      - herald-logs:/var/log/herald
    environment:
      - ENV=production
    env_file:
      - .env.production
    depends_on:
      - redis
    networks:
      - herald-network

  redis:
    image: redis:7-alpine
    container_name: herald-redis
    restart: unless-stopped
    volumes:
      - redis-data:/data
    networks:
      - herald-network
    command: redis-server --appendonly yes

  nginx:
    image: nginx:alpine
    container_name: herald-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.production.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
      - ./web/dist:/usr/share/nginx/html
    depends_on:
      - herald-app
    networks:
      - herald-network

volumes:
  herald-data:
  herald-logs:
  redis-data:

networks:
  herald-network:
    driver: bridge
```

---

## 🌐 Configuration Nginx Production

### `/home/debian/herald/nginx/nginx.production.conf`

```nginx
worker_processes auto;
worker_connections 1024;

events {
    use epoll;
    multi_accept on;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;
    
    # Logs
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                   '$status $body_bytes_sent "$http_referer" '
                   '"$http_user_agent" "$http_x_forwarded_for"';
                   
    access_log /var/log/nginx/access.log main;
    error_log /var/log/nginx/error.log warn;
    
    # Performance
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    client_max_body_size 10M;
    
    # Compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_comp_level 6;
    gzip_types
        text/plain
        text/css
        text/javascript
        text/xml
        application/json
        application/javascript
        application/xml+rss
        application/atom+xml
        image/svg+xml;
    
    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=login:10m rate=1r/s;
    
    # SSL Configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    
    # Security headers
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self' https://api.riotgames.com; frame-ancestors 'none';" always;
    
    # Redirect HTTP to HTTPS
    server {
        listen 80;
        server_name herald-lol.com www.herald-lol.com;
        return 301 https://herald-lol.com$request_uri;
    }
    
    # Main HTTPS server
    server {
        listen 443 ssl http2;
        server_name herald-lol.com;
        
        # SSL certificates
        ssl_certificate /etc/nginx/ssl/herald-lol.com.pem;
        ssl_certificate_key /etc/nginx/ssl/herald-lol.com.key;
        
        # Root directory for static files
        root /usr/share/nginx/html;
        index index.html;
        
        # API proxy
        location /api/ {
            limit_req zone=api burst=20 nodelay;
            
            proxy_pass http://herald-app:8000;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_cache_bypass $http_upgrade;
            
            # Timeouts
            proxy_connect_timeout 30s;
            proxy_send_timeout 30s;
            proxy_read_timeout 30s;
        }
        
        # Auth endpoints with stricter rate limiting
        location /api/auth/ {
            limit_req zone=login burst=5 nodelay;
            proxy_pass http://herald-app:8000;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
        
        # Health check
        location /health {
            proxy_pass http://herald-app:8000;
            access_log off;
        }
        
        # Static files with caching
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
            add_header Vary Accept-Encoding;
        }
        
        # Service Worker
        location /sw.js {
            expires 0;
            add_header Cache-Control "no-cache, no-store, must-revalidate";
        }
        
        # SPA routing - serve index.html for all routes
        location / {
            try_files $uri $uri/ /index.html;
            
            # Cache HTML files for short time
            location ~* \.html$ {
                expires 10m;
                add_header Cache-Control "no-cache";
            }
        }
        
        # Deny access to hidden files
        location ~ /\. {
            deny all;
        }
    }
    
    # www redirect
    server {
        listen 443 ssl http2;
        server_name www.herald-lol.com;
        
        ssl_certificate /etc/nginx/ssl/herald-lol.com.pem;
        ssl_certificate_key /etc/nginx/ssl/herald-lol.com.key;
        
        return 301 https://herald-lol.com$request_uri;
    }
}
```

---

## 🔒 Configuration SSL/HTTPS

### Option 1: Let's Encrypt (Gratuit)

```bash
# Installation Certbot
sudo apt update
sudo apt install certbot

# Génération certificat
sudo certbot certonly --standalone -d herald-lol.com -d www.herald-lol.com

# Copie certificats
sudo cp /etc/letsencrypt/live/herald-lol.com/fullchain.pem /home/debian/herald/ssl/herald-lol.com.pem
sudo cp /etc/letsencrypt/live/herald-lol.com/privkey.pem /home/debian/herald/ssl/herald-lol.com.key
sudo chown -R $USER:$USER /home/debian/herald/ssl/

# Renouvellement automatique
echo "0 2 * * * certbot renew --quiet && docker-compose -f /home/debian/herald/docker-compose.production.yml restart nginx" | sudo crontab -
```

### Option 2: Certificat Commercial
1. Acheter certificat SSL chez un CA reconnu
2. Placer les fichiers dans `/home/debian/herald/ssl/`
3. Configurer le renouvellement selon les instructions du CA

---

## 🚀 Commandes de Déploiement

### Déploiement Initial

```bash
# 1. Cloner et préparer
cd /home/debian/herald
git pull origin main

# 2. Configuration
cp .env.example .env.production
# Éditer .env.production avec vos vraies valeurs

# 3. Build et déploiement
docker-compose -f docker-compose.production.yml build
docker-compose -f docker-compose.production.yml up -d

# 4. Vérifications
docker-compose -f docker-compose.production.yml logs -f
curl https://herald-lol.com/health
```

### Scripts de Maintenance

#### `/home/debian/herald/scripts/deploy.sh`
```bash
#!/bin/bash
set -e

echo "🚀 Déploiement Herald.lol"

# Pull latest code
git pull origin main

# Build frontend
cd web && npm ci && npm run build && cd ..

# Deploy with zero downtime
docker-compose -f docker-compose.production.yml build
docker-compose -f docker-compose.production.yml up -d --no-deps herald-app
docker-compose -f docker-compose.production.yml up -d nginx

echo "✅ Déploiement terminé"
```

#### `/home/debian/herald/scripts/backup.sh`
```bash
#!/bin/bash
BACKUP_DIR="/backup/herald"
DATE=$(date +%Y%m%d_%H%M%S)

# Backup database
mkdir -p $BACKUP_DIR
docker-compose -f docker-compose.production.yml exec herald-app cp /data/herald.db /tmp/
docker cp herald-app:/tmp/herald.db $BACKUP_DIR/herald_$DATE.db

# Backup logs
docker-compose -f docker-compose.production.yml logs > $BACKUP_DIR/logs_$DATE.txt

# Cleanup old backups (keep 30 days)
find $BACKUP_DIR -name "herald_*.db" -mtime +30 -delete
find $BACKUP_DIR -name "logs_*.txt" -mtime +7 -delete

echo "✅ Backup terminé: $BACKUP_DIR/herald_$DATE.db"
```

---

## 📊 Monitoring de Production

### Health Checks
```bash
# Status des services
curl https://herald-lol.com/health

# Métriques détaillées
curl https://herald-lol.com/api/metrics

# Logs en temps réel
docker-compose -f docker-compose.production.yml logs -f --tail=100
```

### Alertes Importantes
1. **Espace disque < 10%**
2. **Mémoire > 90%**
3. **API response time > 2s**
4. **SSL expiration < 30 jours**
5. **Erreur rate > 5%**

---

## 🔧 Dépannage Production

### Problèmes Courants

#### 1. Service ne démarre pas
```bash
# Vérifier logs
docker-compose -f docker-compose.production.yml logs herald-app

# Vérifier configuration
docker-compose -f docker-compose.production.yml config
```

#### 2. Performance lente
```bash
# Vérifier ressources
docker stats

# Optimiser base de données
docker-compose -f docker-compose.production.yml exec herald-app sqlite3 /data/herald.db "VACUUM;"
```

#### 3. SSL invalide
```bash
# Renouveler certificat
sudo certbot renew
docker-compose -f docker-compose.production.yml restart nginx
```

---

## ✅ Checklist de Mise en Production

### Avant le Déploiement
- [ ] Domaine configuré et pointant vers le serveur
- [ ] Certificats SSL générés et installés
- [ ] Variables d'environnement configurées
- [ ] Google OAuth configuré avec bons domaines
- [ ] Clé Riot API active et testée
- [ ] Tests d'intégration passés (100%)

### Après le Déploiement
- [ ] Application accessible via HTTPS
- [ ] Tous les endpoints API fonctionnels
- [ ] Authentification Google opérationnelle
- [ ] Interface responsive sur mobile/desktop
- [ ] Performance optimale (< 2s load time)
- [ ] Monitoring et logs configurés
- [ ] Backup automatique programmé

---

## 🎯 Performance Attendue

### Métriques Cibles
- **Temps de chargement** : < 2 secondes
- **First Contentful Paint** : < 1.5s
- **API Response Time** : < 500ms
- **Uptime** : 99.9%
- **Concurrent Users** : 1000+
- **Bundle Size** : 703KB (gzipped)

### Capacité
- **Utilisateurs simultanés** : 1000+
- **Groupes maximum** : Illimité
- **Comparaisons/heure** : 10,000+
- **API Calls/minute** : 6,000+

---

**Herald.lol est prêt pour la production ! 🚀**

*Guide mis à jour le $(date)*