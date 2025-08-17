# Checklist Configuration Jenkins Freestyle

## ✅ Prérequis Jenkins

- [ ] Jenkins installé et accessible
- [ ] Plugins installés :
  - [ ] Git Plugin
  - [ ] SSH Pipeline Steps
  - [ ] Email Extension Plugin
  - [ ] Slack Notification Plugin (optionnel)
  - [ ] Timestamper Plugin
- [ ] Go installé sur le serveur Jenkins
- [ ] Docker installé sur le serveur Jenkins

## ✅ Configuration du projet

- [ ] Projet freestyle créé : `lol-match-exporter-freestyle`
- [ ] Repository Git configuré
- [ ] Branche `main` sélectionnée
- [ ] Build triggers configurés (Poll SCM)
- [ ] Clean workspace activé

## ✅ Build Steps configurés

- [ ] Étape 1 : Préparation et vérifications
- [ ] Étape 2 : Tests et build Go
- [ ] Étape 3 : Build Docker
- [ ] Étape 4 : Tests de l'image Docker  
- [ ] Étape 5 : Déploiement conditionnel

## ✅ Credentials et sécurité

- [ ] Clé SSH créée pour le serveur de production
- [ ] Credential SSH ajouté dans Jenkins (`production-server-ssh`)
- [ ] Test de connexion SSH réussi
- [ ] Variables d'environnement définies

## ✅ Post-build actions

- [ ] Archive des artifacts configuré
- [ ] Notifications email configurées
- [ ] Notifications Slack configurées (optionnel)
- [ ] Rotation des builds configurée (garder 10 builds)

## ✅ Tests de validation

- [ ] Build de test manuel réussi
- [ ] Logs de build vérifiés
- [ ] Déploiement de test réussi
- [ ] Health check post-déploiement OK

## ✅ Webhooks et automatisation

- [ ] Webhook GitHub configuré (si applicable)
- [ ] Déclenchement automatique testé
- [ ] Notifications fonctionnelles

## 🔧 Commandes de test

```bash
# Tester la configuration
chmod +x scripts/jenkins-test.sh
./scripts/jenkins-test.sh

# Tester le build local
make ci-test

# Tester le déploiement
make deploy-staging
```

## 📞 Dépannage

### Build échoue sur "go: command not found"
- Installer Go sur le serveur Jenkins
- Ajouter Go au PATH dans Jenkins (Manage Jenkins > System Configuration)

### Erreur SSH "Permission denied"
- Vérifier que la clé SSH est correcte
- Tester manuellement : `ssh debian@51.178.17.78`
- Vérifier les permissions de la clé

### Docker build échoue
- Vérifier que Docker est installé et démarré
- Ajouter l'utilisateur Jenkins au groupe docker : `sudo usermod -aG docker jenkins`
- Redémarrer Jenkins

### Déploiement échoue
- Vérifier la connectivité réseau vers le serveur
- Vérifier que Docker est installé sur le serveur de production
- Tester les commandes de déploiement manuellement

## 🎯 URLs importantes

- Jenkins Dashboard: http://your-jenkins:8080
- Projet: http://your-jenkins:8080/job/lol-match-exporter-freestyle/
- Console Output: http://your-jenkins:8080/job/lol-match-exporter-freestyle/lastBuild/console
- Production App: https://herald.lol
- Health Check: https://herald.lol/api/health
