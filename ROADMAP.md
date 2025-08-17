# LoL Analytics Platform - Roadmap de développement

## 🎯 **Vision**
Transformer le LoL Match Exporter en plateforme d'analytics complète avec :
- Dashboard intelligent multi-temporel
- Suggestions personnalisées par champion/rôle
- Analyse MMR et évolution des performances
- Prédictions et recommandations adaptatives

## 📋 **Architecture complète**

### **1. Frontend - Suppression + Nouvelles pages**
- **Supprimer** : Section "Fichiers générés" (ExporterMUI.tsx:531-564)
- **Ajouter** : 
  - `pages/Dashboard.tsx` - Analytics multi-temporel
  - `pages/Champions.tsx` - Suggestions par champion/rôle
  - `pages/MMRAnalytics.tsx` - Évolution MMR et prédictions
  - `components/StatsSummary.tsx` - Widgets de statistiques
  - `components/ChampionRecommendations.tsx` - IA suggestions

### **2. Base de données étendue**
```sql
-- Tables existantes
users, matches, scan_history

-- Nouvelles tables analytics
champion_stats (
  id INTEGER PRIMARY KEY,
  user_id INTEGER,
  champion_id INTEGER,
  role TEXT,
  season INTEGER,
  games_played INTEGER,
  win_rate REAL,
  avg_kda REAL,
  avg_cs_per_min REAL,
  performance_score REAL,
  FOREIGN KEY(user_id) REFERENCES users(id)
)

role_performance (
  id INTEGER PRIMARY KEY,
  user_id INTEGER,
  role TEXT,
  time_period TEXT, -- 'today', 'week', 'month', 'season'
  games INTEGER,
  win_rate REAL,
  avg_performance REAL,
  trend_direction TEXT, -- 'improving', 'declining', 'stable'
  FOREIGN KEY(user_id) REFERENCES users(id)
)

mmr_history (
  id INTEGER PRIMARY KEY,
  user_id INTEGER,
  match_id TEXT,
  estimated_mmr INTEGER,
  mmr_change INTEGER,
  confidence_score REAL,
  game_date TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id)
)

performance_insights (
  id INTEGER PRIMARY KEY,
  user_id INTEGER,
  insight_type TEXT, -- 'suggestion', 'warning', 'achievement'
  category TEXT, -- 'champion', 'role', 'gameplay', 'mmr'
  title TEXT,
  description TEXT,
  priority INTEGER,
  time_period TEXT,
  is_active BOOLEAN,
  created_at TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id)
)
```

### **3. Backend - Modules d'intelligence**

#### **Nouveau : `analytics_engine.py`**
```python
class AnalyticsEngine:
    def generate_period_stats(user_id, period): # today/week/month/season
    def calculate_role_performance(user_id, role, period):
    def analyze_champion_mastery(user_id, champion_id):
    def generate_improvement_suggestions(user_id):
    def calculate_performance_trends(user_id):
```

#### **Nouveau : `mmr_calculator.py`**
```python
class MMRAnalyzer:
    def estimate_mmr_from_match(match_data, opponent_ranks):
    def calculate_mmr_trajectory(user_id):
    def predict_rank_changes(user_id):
    def analyze_mmr_volatility(user_id):
    def calculate_skill_ceiling(user_id, role):
```

#### **Nouveau : `recommendation_engine.py`**
```python
class RecommendationEngine:
    def suggest_champions_for_role(user_id, role, meta_data):
    def analyze_champion_performance_gaps(user_id):
    def generate_gameplay_tips(user_id, champion_id):
    def recommend_ban_priorities(user_id, current_meta):
    def suggest_build_optimizations(user_id, champion_id):
```

### **4. Pages et fonctionnalités**

#### **Dashboard Analytics (`/dashboard`)**
- **Période sélectionnable** : Aujourd'hui / Cette semaine / Ce mois / Cette saison
- **Widgets dynamiques** :
  ```typescript
  interface PeriodStats {
    winRate: number;
    avgKDA: number;
    gamesPlayed: number;
    rolePerformance: RoleStats[];
    topChampions: ChampionPerf[];
    mmrTrend: 'up' | 'down' | 'stable';
    suggestions: Suggestion[];
  }
  ```
- **Suggestions contextuelles** :
  - **Aujourd'hui** : "Essaie Jinx ADC, +73% winrate récent"
  - **Semaine** : "Améliore ton early game, -15% CS@10 vs moyenne"
  - **Mois** : "Focus Support, ton meilleur rôle (+67% WR)"
  - **Saison** : "Objectif: Diamond 4 atteignable (+180 LP estimés)"

#### **Analytics Champions (`/champions`)**
- **Filtrés par rôle** : Top, Jungle, Mid, ADC, Support
- **Métriques par champion** :
  - Performance score (algorithme composite)
  - Tendance (amélioration/déclin)
  - Potentiel de climb (+/- LP estimé)
  - Recommendations build/runes
- **Suggestions intelligentes** :
  ```javascript
  {
    champion: "Thresh",
    role: "Support", 
    status: "declining",
    suggestion: "Focus sur le warding (+23% vision score needed)",
    priority: "high",
    expectedGain: "+12% winrate"
  }
  ```

#### **MMR Analytics (`/mmr`)**
- **Graphique évolution MMR** avec prédictions
- **Analyse des gains/pertes** :
  - MMR moyen estimé
  - Volatilité (stable vs. streaky)
  - Potentiel de montée par rôle
  - Prédiction rank fin de saison
- **Métriques avancées** :
  - Cote de confiance MMR (A+ à F)
  - Analyse comparative (vs. peers)
  - Identification de skill gaps

### **5. Intelligence artificielle et suggestions**

#### **Algorithmes de suggestions**
- **Performance trending** : Détection automatique des améliorations/déclins
- **Meta adaptation** : Suggestions basées sur patch notes + winrates globaux
- **Personnalisation** : Adaptées au style de jeu et niveau du joueur
- **Temporal context** : Différentes selon la période analysée

#### **Système de scoring**
```python
def calculate_performance_score(champion_stats):
    base_score = (win_rate * 0.4 + 
                  normalized_kda * 0.3 + 
                  cs_efficiency * 0.2 + 
                  objective_participation * 0.1)
    
    trend_modifier = calculate_recent_trend()
    meta_modifier = get_meta_strength(champion_id)
    
    return base_score * trend_modifier * meta_modifier
```

### **6. API Extensions**

#### **Nouveaux endpoints**
```python
# Analytics endpoints
GET /api/users/{puuid}/dashboard/{period}  # today|week|month|season
GET /api/users/{puuid}/champions/{role}
GET /api/users/{puuid}/mmr/history
GET /api/users/{puuid}/suggestions
GET /api/users/{puuid}/predictions

# Real-time insights
GET /api/users/{puuid}/insights/live
POST /api/users/{puuid}/feedback  # ML feedback loop
```

### **7. Architecture technique**

#### **Frontend Navigation**
```typescript
const routes = [
  { path: '/', component: ExporterPage },
  { path: '/dashboard', component: AnalyticsDashboard },
  { path: '/champions', component: ChampionsAnalytics },
  { path: '/mmr', component: MMRAnalytics },
  { path: '/profile', component: UserProfile }
]
```

#### **Data Pipeline**
1. **Scraping** → **Raw matches** (BDD)
2. **Processing** → **Computed stats** (champion_stats, role_performance)
3. **Analysis** → **Insights generation** (performance_insights)
4. **Presentation** → **Dashboard widgets**

### **8. Calculs MMR et prédictions**

#### **Algorithme MMR estimation**
```python
def estimate_mmr(match_data):
    # Facteurs: ranks adversaires, team MMR moyenne, résultat
    base_mmr = calculate_average_opponent_mmr(match_data)
    performance_modifier = analyze_individual_performance(match_data)
    confidence = calculate_confidence_score(match_history)
    
    return {
        'estimated_mmr': base_mmr + performance_modifier,
        'confidence': confidence,
        'change': calculate_mmr_change(previous_mmr, result)
    }
```

#### **Prédictions et tendances**
- **Court terme** (jour/semaine) : Performance immédiate, form récente
- **Moyen terme** (mois) : Progression skill, adaptation meta
- **Long terme** (saison) : Potentiel maximum, objectifs réalistes

## ✅ **Impact utilisateur**

### **Expérience transformée**
- **Export simple** → **Plateforme analytics complète**
- **Données brutes** → **Insights actionnables**
- **Statique** → **Prédictif et adaptatif**

### **Valeur ajoutée**
- **Coaching IA** : Suggestions personnalisées pour améliorer
- **Optimisation** : Focus sur champions/rôles les plus profitables
- **Motivation** : Objectifs clairs et atteignables
- **Stratégie** : Décisions éclairées (bans, picks, builds)

## 🚀 **Roadmap d'implémentation**

### **Phase 1 : Foundation (2-3 jours)**
- [ ] Créer le schéma de base de données étendu
- [ ] Implémenter les migrations SQLite
- [ ] Modifier le scraper pour usage incrémental
- [ ] Supprimer l'interface "Fichiers générés"

### **Phase 2 : Analytics Engine (2-3 jours)**
- [ ] Module `analytics_engine.py`
- [ ] Calculs de statistiques par période
- [ ] Algorithmes de détection de tendances
- [ ] Dashboard frontend basique

### **Phase 3 : MMR & Prédictions (2 jours)**
- [ ] Module `mmr_calculator.py`
- [ ] Estimation MMR par match
- [ ] Trajectoire et prédictions
- [ ] Page MMR Analytics

### **Phase 4 : Recommendation Engine (3-4 jours)**
- [ ] Module `recommendation_engine.py`
- [ ] Algorithmes de suggestions
- [ ] Page Champions Analytics
- [ ] Système de scoring composite

### **Phase 5 : Polish & Optimisation (1-2 jours)**
- [ ] Tests et debugging
- [ ] Optimisation des performances
- [ ] UX/UI refinements
- [ ] Documentation

## 📊 **Métriques de succès**
- **Performance** : Temps de chargement < 2s
- **Précision** : MMR estimation ±150 pts
- **Engagement** : Suggestions actionnables > 80%
- **Utility** : Amélioration winrate mesurable

---

*Dernière mise à jour : 16 août 2025*