# Guide de configuration Jenkins pour LoL Match Exporter

## Installation et configuration initiale

### 1. Installation de Jenkins

```bash
# Sur Ubuntu/Debian
curl -fsSL https://pkg.jenkins.io/debian-stable/jenkins.io.key | sudo tee /usr/share/keyrings/jenkins-keyring.asc > /dev/null
echo deb [signed-by=/usr/share/keyrings/jenkins-keyring.asc] https://pkg.jenkins.io/debian-stable binary/ | sudo tee /etc/apt/sources.list.d/jenkins.list > /dev/null
sudo apt update
sudo apt install jenkins openjdk-11-jdk

# Démarrer Jenkins
sudo systemctl start jenkins
sudo systemctl enable jenkins
```

### 2. Configuration initiale

1. Accéder à Jenkins : `http://your-server:8080`
2. Récupérer le mot de passe initial : `sudo cat /var/lib/jenkins/secrets/initialAdminPassword`
3. Installer les plugins recommandés
4. Créer un utilisateur admin

### 3. Plugins requis

Installer les plugins suivants via **Manage Jenkins > Plugin Manager** :

- **Pipeline** : Support des pipelines Jenkins
- **Git Pipeline** : Intégration Git
- **Docker Pipeline** : Support Docker dans les pipelines
- **SSH Pipeline Steps** : Exécution de commandes SSH
- **Slack Notification** : Notifications Slack (optionnel)
- **Email Extension** : Notifications email avancées
- **Blue Ocean** : Interface moderne (optionnel)
- **Pipeline Utility Steps** : Utilitaires pour pipelines
- **Timestamper** : Timestamps dans les logs

## Configuration des credentials

### 1. Clé SSH pour le serveur de production

```bash
# Générer une clé SSH pour Jenkins
sudo -u jenkins ssh-keygen -t rsa -b 4096 -C "jenkins@yourdomain.com"

# Copier la clé publique sur le serveur de production
sudo -u jenkins ssh-copy-id debian@51.178.17.78
```

Dans Jenkins :

1. **Manage Jenkins > Credentials**
2. **Add Credentials**
3. Type : **SSH Username with private key**
4. ID : `production-server-ssh`
5. Username : `debian`
6. Private Key : Coller la clé privée de `/var/lib/jenkins/.ssh/id_rsa`

### 2. GitHub Token (si repo privé)

1. **GitHub > Settings > Developer Settings > Personal Access Tokens**
2. Générer un token avec permissions `repo`
3. Dans Jenkins : **Credentials > Add > Secret text**
4. ID : `github-token`

### 3. Slack Webhook (optionnel)

1. **Slack > Apps > Jenkins CI**
2. Configurer le webhook
3. Dans Jenkins : **Credentials > Add > Secret text**
4. ID : `slack-webhook`

## Configuration du job Jenkins

### 1. Créer un nouveau Pipeline

1. **New Item > Pipeline**
2. Nom : `lol-match-exporter-pipeline`
3. **Pipeline > Definition : Pipeline script from SCM**
4. SCM : Git
5. Repository URL : `https://github.com/your-username/lol_match_exporter.git`
6. Credentials : Sélectionner si repo privé
7. Branch : `*/main`
8. Script Path : `Jenkinsfile`

### 2. Configuration des triggers

**Build Triggers :**

- ☑️ **GitHub hook trigger for GITScm polling** (si webhook configuré)
- ☑️ **Poll SCM** : `H/5 * * * *` (vérifie toutes les 5 minutes)

### 3. Configuration avancée

**Pipeline :**

- Lightweight checkout : ☑️
- **Build Triggers > Advanced > Allowed branches** : `main develop feature/*`

## Variables d'environnement globales

**Manage Jenkins > Configure System > Global Properties > Environment variables :**

```
DEPLOY_HOST=51.178.17.78
DEPLOY_USER=debian
DOCKER_REGISTRY_URL=
SLACK_WEBHOOK_URL=[if configured]
NOTIFICATION_EMAIL=admin@yourdomain.com
```

## Configuration GitHub Webhook

### 1. Dans GitHub Repository Settings

1. **Settings > Webhooks > Add webhook**
2. **Payload URL** : `http://your-jenkins-server:8080/github-webhook/`
3. **Content type** : `application/json`
4. **Events** :
   - ☑️ Push events
   - ☑️ Pull request events
5. **Active** : ☑️

### 2. Sécurité Jenkins

Si Jenkins est exposé publiquement, configurer la sécurité :

```bash
# Configuration nginx pour Jenkins (optionnel)
sudo nano /etc/nginx/sites-available/jenkins

server {
    listen 80;
    server_name jenkins.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Stratégie de branchement

Le Jenkinsfile est configuré pour :

- **main** : Déploiement automatique en production (avec confirmation)
- **develop** : Déploiement automatique en staging
- **feature/** : Build et tests seulement

## Notifications

### 1. Configuration Slack

```groovy
// Dans le Jenkinsfile, section post
post {
    success {
        slackSend(
            channel: '#deployments',
            color: 'good',
            message: "✅ Deployment successful: ${env.JOB_NAME} - ${env.BUILD_NUMBER}"
        )
    }
    failure {
        slackSend(
            channel: '#deployments',
            color: 'danger',
            message: "❌ Deployment failed: ${env.JOB_NAME} - ${env.BUILD_NUMBER}"
        )
    }
}
```

### 2. Configuration Email

**Manage Jenkins > Configure System > Extended E-mail Notification :**

- **SMTP Server** : `smtp.gmail.com`
- **Port** : `587`
- **Username/Password** : Credentials email
- **Use SSL** : ☑️

## Monitoring et logs

### 1. Accès aux logs

- **Jenkins Dashboard > Build History**
- **Blue Ocean** : Interface moderne
- **Console Output** : Logs détaillés

### 2. Métriques

Installer **Metrics Plugin** pour :

- Durée des builds
- Taux de succès/échec
- Utilisation des resources

## Maintenance

### 1. Nettoyage automatique

Configuration dans le Jenkinsfile :

```groovy
options {
    buildDiscarder(logRotator(
        numToKeepStr: '10',
        artifactNumToKeepStr: '5'
    ))
}
```

### 2. Backup Jenkins

```bash
# Script de backup quotidien
#!/bin/bash
BACKUP_DIR="/backup/jenkins"
mkdir -p $BACKUP_DIR
tar -czf $BACKUP_DIR/jenkins-backup-$(date +%Y%m%d).tar.gz -C /var/lib/jenkins .
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
```

## Troubleshooting

### Problèmes courants

1. **Permission denied SSH** : Vérifier les clés SSH et permissions
2. **Docker not found** : Installer Docker sur Jenkins server
3. **Build timeout** : Augmenter timeout dans options du pipeline
4. **Out of disk space** : Configurer nettoyage automatique

### Commandes utiles

```bash
# Redémarrer Jenkins
sudo systemctl restart jenkins

# Vérifier les logs
sudo journalctl -u jenkins -f

# Vérifier l'espace disque
df -h /var/lib/jenkins

# Nettoyer les builds anciens
find /var/lib/jenkins/jobs/*/builds -mtime +30 -delete
```

## Sécurité

### 1. Configuration recommandée

- **Manage Jenkins > Configure Global Security**
- **Enable security** : ☑️
- **Security Realm** : Jenkins' own user database
- **Authorization** : Matrix-based security
- **Prevent Cross Site Request Forgery** : ☑️

### 2. Utilisateurs et permissions

Créer des rôles :

- **Admin** : Toutes permissions
- **Developer** : Build, cancel, read
- **Viewer** : Read seulement

### 3. Audit

Installer **Audit Trail Plugin** pour tracer :

- Connexions
- Modifications de configuration
- Exécutions de builds

## Performance

### 1. Optimisation

- **Manage Jenkins > Configure System > # of executors** : Ajuster selon CPU
- **Pipeline > Parallel stages** : Utiliser la parallélisation
- **Docker > Reuse containers** : Éviter les rebuilds

### 2. Monitoring

```bash
# Monitoring des resources Jenkins
htop
iotop
docker stats
```

Cette configuration Jenkins permet un déploiement automatisé, sécurisé et monitoré de votre application LoL Match Exporter.
