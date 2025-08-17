# Guide d'Utilisation des Tests - LoL Match Exporter

## 📋 Vue d'ensemble

Ce guide vous explique comment utiliser la suite de tests complète du projet LoL Match Exporter. La stratégie de tests couvre tous les aspects de l'application, des tests unitaires aux tests de bout en bout.

## 🚀 Démarrage Rapide

### Prérequis

- **Go 1.21+**
- **Node.js 18+** 
- **Python 3.11+**
- **PostgreSQL 15+** (pour les tests d'intégration)
- **Docker** (pour les tests de déploiement)
- **k6** (pour les tests de performance)

### Installation des outils de test

```bash
# Outils Go
go install github.com/stretchr/testify@latest
go install github.com/DATA-DOG/go-sqlmock@latest
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Outils Python
pip install pytest pytest-cov pytest-mock pytest-xdist

# Outils Node.js (dans le dossier web/)
cd web && npm install

# Playwright (dans tests/e2e/)
cd tests/e2e && npm install && npx playwright install

# k6 (Windows)
choco install k6
# ou Linux
sudo apt-get install k6
```

## 🔧 Exécution des Tests

### Script Principal

Le moyen le plus simple d'exécuter les tests est d'utiliser le script PowerShell principal :

```powershell
# Tous les tests
.\run_tests.ps1

# Tests spécifiques
.\run_tests.ps1 -TestType unit
.\run_tests.ps1 -TestType integration
.\run_tests.ps1 -TestType e2e
.\run_tests.ps1 -TestType performance
.\run_tests.ps1 -TestType security

# Avec coverage
.\run_tests.ps1 -Coverage

# Mode verbose
.\run_tests.ps1 -Verbose

# Sans génération de rapport
.\run_tests.ps1 -Report:$false
```

### Tests Individuels

#### Tests Go (Unitaires)

```bash
# Tous les tests Go
go test ./tests/unit/go/... -v

# Avec coverage
go test ./tests/unit/go/... -v -cover -coverprofile=coverage.out

# Tests d'un service spécifique
go test ./tests/unit/go/ -run TestAnalyticsService -v

# Tests en parallèle
go test ./tests/unit/go/... -v -parallel 4

# Tests avec race detection
go test ./tests/unit/go/... -v -race
```

#### Tests Go (Intégration)

```bash
# Setup de la base de test (requis)
createdb lol_analytics_test
psql lol_analytics_test < migrations/001_initial_schema.sql

# Tests d'intégration
go test ./tests/integration/... -v

# Avec variables d'environnement
DATABASE_URL="postgres://user:pass@localhost/lol_analytics_test" go test ./tests/integration/... -v
```

#### Tests Python

```bash
# Tests unitaires Python
python -m pytest tests/unit/python/ -v

# Avec coverage
python -m pytest tests/unit/python/ -v --cov=. --cov-report=html

# Tests spécifiques
python -m pytest tests/unit/python/test_analytics_engine.py::TestAnalyticsEngine::test_calculate_period_stats_success -v

# Tests en parallèle
python -m pytest tests/unit/python/ -v -n 4

# Mode debug
python -m pytest tests/unit/python/ -v -s --pdb
```

#### Tests React

```bash
cd web

# Tests unitaires React
npm test

# Avec coverage
npm test -- --coverage

# Tests spécifiques
npm test -- --testNamePattern="Dashboard"

# Mode watch
npm test -- --watch

# TypeScript type check
npm run type-check

# Linting
npm run lint
```

#### Tests End-to-End

```bash
cd tests/e2e

# Installer Playwright
npm ci
npx playwright install

# Tous les tests E2E
npx playwright test

# Tests sur navigateur spécifique
npx playwright test --project=chromium

# Mode debug
npx playwright test --debug

# Mode headed (avec interface graphique)
npx playwright test --headed

# Tests spécifiques
npx playwright test notification_flow.spec.js

# Générer rapport HTML
npx playwright show-report
```

#### Tests de Performance

```bash
# Démarrer le serveur (requis)
go run analytics_server_standalone.go &

# Tests de charge basiques
k6 run tests/performance/load_test.js

# Test de stress
k6 run tests/performance/load_test.js --env SCENARIO=stress

# Test de pic de charge
k6 run tests/performance/load_test.js --env SCENARIO=spike

# Avec rapport JSON
k6 run tests/performance/load_test.js --out json=results.json
```

#### Tests de Sécurité

```bash
# Scan Go avec gosec
gosec ./...

# Audit npm
cd web && npm audit

# Scan complet avec rapport
gosec -fmt json -out security-report.json ./...
```

## 📊 Types de Tests

### 1. Tests Unitaires

**Objectif :** Tester les composants individuels en isolation

**Couverture :**
- Services Go (AnalyticsService, NotificationService, SyncService)
- Engines Python (analytics_engine, mmr_calculator, recommendation_engine)
- Composants React (Dashboard, Charts, Forms)

**Commandes :**
```bash
# Go
go test ./tests/unit/go/... -v -cover

# Python  
python -m pytest tests/unit/python/ -v --cov=.

# React
cd web && npm test -- --coverage
```

### 2. Tests d'Intégration

**Objectif :** Tester l'interaction entre les composants

**Couverture :**
- API handlers avec base de données
- Communication Python-Go
- Frontend-Backend API calls

**Commandes :**
```bash
# API Integration
go test ./tests/integration/api/... -v

# Database Integration  
go test ./tests/integration/database/... -v

# Python-Go Integration
go test ./tests/integration/python-go/... -v
```

### 3. Tests End-to-End

**Objectif :** Tester les parcours utilisateur complets

**Couverture :**
- Authentification complète
- Navigation dans l'interface
- Fonctionnalités analytics
- Notifications temps réel

**Commandes :**
```bash
cd tests/e2e
npx playwright test

# Tests spécifiques
npx playwright test auth.spec.js
npx playwright test notification_flow.spec.js
```

### 4. Tests de Performance

**Objectif :** Valider les performances sous charge

**Couverture :**
- Endpoints API sous charge
- Système de notifications en temps réel
- Performance base de données
- Temps de réponse frontend

**Commandes :**
```bash
# Load testing
k6 run tests/performance/load_test.js

# Stress testing
k6 run tests/performance/load_test.js --env SCENARIO=stress
```

### 5. Tests de Sécurité

**Objectif :** Identifier les vulnérabilités de sécurité

**Couverture :**
- Vulnérabilités code Go
- Dépendances npm vulnérables
- Configuration sécurité
- Tests injection/XSS

**Commandes :**
```bash
# Go security scan
gosec ./...

# npm vulnerabilities
cd web && npm audit

# OWASP dependency check
dependency-check --project lol-match-exporter --scan .
```

## 🔍 Debugging des Tests

### Tests Go

```bash
# Verbose output
go test ./tests/unit/go/... -v

# Avec debugging
go test ./tests/unit/go/... -v -debug

# Tests spécifiques
go test ./tests/unit/go/ -run TestSpecificFunction -v

# Profiling
go test ./tests/unit/go/... -cpuprofile=cpu.prof -memprofile=mem.prof
```

### Tests Python

```bash
# Mode debug avec pdb
python -m pytest tests/unit/python/ -v -s --pdb

# Logs détaillés
python -m pytest tests/unit/python/ -v -s --log-cli-level=DEBUG

# Tests spécifiques avec print
python -m pytest tests/unit/python/test_analytics_engine.py::test_function -v -s
```

### Tests E2E

```bash
# Mode debug interactif
npx playwright test --debug

# Avec browser visible
npx playwright test --headed

# Screenshots sur échec
npx playwright test --screenshot=only-on-failure

# Traces détaillées
npx playwright test --trace=on
```

## 📈 Rapports et Métriques

### Coverage Reports

```bash
# Go coverage
go test ./tests/unit/go/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Python coverage
python -m pytest tests/unit/python/ --cov=. --cov-report=html

# React coverage
cd web && npm test -- --coverage
```

### Test Reports

Les rapports sont générés dans le dossier `test-reports/` :

- `go-coverage.html` - Coverage Go
- `python-coverage/` - Coverage Python
- `react-coverage/` - Coverage React  
- `playwright-report/` - Rapports E2E
- `performance-results.json` - Résultats de performance
- `security-report.json` - Rapport de sécurité
- `comprehensive-report.md` - Rapport global

### Métriques de Qualité

**Objectifs de Coverage :**
- Go Services: ≥85%
- Python Engines: ≥90%
- React Components: ≥80%
- Integration: ≥70%

**Seuils de Performance :**
- API Response Time: p95 <500ms
- Frontend Load Time: <3s
- Database Queries: <100ms avg
- Error Rate: <1%

## 🚀 CI/CD et Automatisation

### GitHub Actions

Le workflow `.github/workflows/comprehensive-tests.yml` exécute automatiquement :

- ✅ Tests unitaires (Go, Python, React)
- ✅ Tests d'intégration  
- ✅ Tests E2E
- ✅ Tests de performance
- ✅ Tests de sécurité
- ✅ Tests Docker
- ✅ Génération de rapports

**Déclencheurs :**
- Push sur `main` ou `develop`
- Pull requests vers `main`
- Schedule quotidien (2h UTC)
- Manuel via workflow dispatch

### Local CI Simulation

```powershell
# Simuler le pipeline CI en local
.\run_tests.ps1 -TestType all -Coverage -Report

# Tests rapides (comme dans PR)
.\run_tests.ps1 -TestType unit -Coverage

# Tests complets (comme dans main)  
.\run_tests.ps1 -TestType all -Coverage -Verbose
```

## 🛠️ Configuration Environnements

### Développement

```bash
# Variables d'environnement
export TESTING=true
export DATABASE_URL=postgres://user:pass@localhost/lol_analytics_test
export REDIS_URL=redis://localhost:6379
export RIOT_API_KEY=test_key
```

### CI/CD

Variables configurées dans GitHub Secrets :
- `DATABASE_URL`
- `REDIS_URL`  
- `RIOT_API_KEY`
- `SLACK_WEBHOOK_URL`

### Production

Tests de smoke uniquement :
```bash
.\run_tests.ps1 -TestType smoke -Environment production
```

## 📚 Bonnes Pratiques

### Écriture de Tests

1. **Nommage clair :** `Test_FunctionName_Scenario`
2. **Structure AAA :** Arrange, Act, Assert
3. **Tests isolés :** Pas de dépendances entre tests
4. **Mocks appropriés :** Isoler les dépendances externes
5. **Données de test :** Utiliser les fixtures dans `testing-utils/`

### Maintenance

1. **Tests parallèles :** Utiliser `-parallel` pour Go, `-n` pour Python
2. **Cleanup :** Nettoyer les ressources après tests
3. **Timeouts :** Définir des timeouts appropriés
4. **Retry :** Retry automatique pour tests flaky
5. **Monitoring :** Surveiller les métriques de tests

### Debugging

1. **Logs détaillés :** Utiliser `-v` et modes debug
2. **Screenshots :** Activer pour tests E2E
3. **Traces :** Capturer les traces d'exécution
4. **Profiling :** Analyser les performances
5. **Coverage :** Identifier le code non testé

## 🔧 Dépannage

### Problèmes Courants

**Tests Go qui échouent :**
```bash
# Vérifier les dépendances
go mod tidy
go mod download

# Nettoyer le cache
go clean -testcache
```

**Tests Python qui échouent :**
```bash
# Réinstaller les dépendances
pip install -r requirements.txt --force-reinstall

# Vérifier la version Python
python --version
```

**Tests E2E qui échouent :**
```bash
# Réinstaller Playwright
cd tests/e2e
npx playwright install --force
```

**Tests de performance qui échouent :**
```bash
# Vérifier que le serveur est démarré
curl http://localhost:8001/api/health

# Ajuster les seuils dans load_test.js
```

### Support

- 📖 Documentation : `TESTING_STRATEGY.md`
- 🐛 Issues : GitHub Issues
- 💬 Discussion : GitHub Discussions
- 📧 Contact : Équipe de développement

---

## 📝 Historique des Tests

| Version | Date | Ajouts |
|---------|------|--------|
| 1.0 | 2025-01 | Tests Go unitaires et intégration |
| 1.1 | 2025-01 | Tests Python analytics engines |
| 1.2 | 2025-01 | Tests React composants |
| 1.3 | 2025-01 | Tests E2E Playwright |
| 1.4 | 2025-01 | Tests performance k6 |
| 1.5 | 2025-01 | Tests sécurité et CI/CD |

---

*Dernière mise à jour : Janvier 2025*