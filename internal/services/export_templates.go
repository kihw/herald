package services

import (
	"fmt"
	"time"
)

// ExportTemplate définit un template d'export prédéfini
type ExportTemplate struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Options     ExportOptions `json:"options"`
	CreatedAt   time.Time     `json:"created_at"`
	IsDefault   bool          `json:"is_default"`
}

// TemplateService gère les templates d'export
type TemplateService struct {
	templates map[string]ExportTemplate
}

// NewTemplateService crée une nouvelle instance du service de templates
func NewTemplateService() *TemplateService {
	service := &TemplateService{
		templates: make(map[string]ExportTemplate),
	}
	
	// Charger les templates par défaut
	service.loadDefaultTemplates()
	
	return service
}

// loadDefaultTemplates charge les templates par défaut
func (ts *TemplateService) loadDefaultTemplates() {
	defaultTemplates := []ExportTemplate{
		{
			ID:          "basic_csv",
			Name:        "Export CSV Basic",
			Description: "Export CSV simple avec les statistiques essentielles",
			Options: ExportOptions{
				Format: FormatCSV,
				Filter: ExportFilter{
					RecentFirst:   true,
					IncludeRemake: false,
				},
				Columns: []string{
					"match_id", "game_creation", "champion_name", "win", 
					"kills", "deaths", "assists", "kda", "cs", "gold",
				},
				Compression: false,
				Metadata:    false,
			},
			CreatedAt: time.Now(),
			IsDefault: true,
		},
		{
			ID:          "ranked_detailed",
			Name:        "Parties Classées Détaillées",
			Description: "Export complet des parties classées avec toutes les statistiques",
			Options: ExportOptions{
				Format: FormatExcel,
				Filter: ExportFilter{
					RankedOnly:    &[]bool{true}[0],
					RecentFirst:   true,
					IncludeRemake: false,
				},
				Columns: []string{
					"match_id", "game_creation", "game_duration", "queue_id",
					"champion_name", "role", "lane", "win", "kills", "deaths", 
					"assists", "kda", "cs", "gold", "damage", "vision", 
					"rank", "lp", "mmr",
				},
				Compression: true,
				Metadata:    true,
			},
			CreatedAt: time.Now(),
			IsDefault: true,
		},
		{
			ID:          "win_analysis",
			Name:        "Analyse des Victoires",
			Description: "Export JSON des victoires pour analyse statistique",
			Options: ExportOptions{
				Format: FormatJSON,
				Filter: ExportFilter{
					WinOnly:       &[]bool{true}[0],
					RecentFirst:   true,
					IncludeRemake: false,
				},
				Compression: false,
				Metadata:    true,
			},
			CreatedAt: time.Now(),
			IsDefault: true,
		},
		{
			ID:          "champion_performance",
			Name:        "Performance par Champion",
			Description: "Export détaillé pour analyser les performances par champion",
			Options: ExportOptions{
				Format: FormatCSV,
				Filter: ExportFilter{
					RecentFirst:   true,
					IncludeRemake: false,
				},
				Columns: []string{
					"champion_name", "win", "kills", "deaths", "assists", 
					"kda", "cs", "gold", "damage", "game_duration",
				},
				Compression: false,
				Metadata:    false,
			},
			CreatedAt: time.Now(),
			IsDefault: true,
		},
		{
			ID:          "monthly_summary",
			Name:        "Résumé Mensuel",
			Description: "Export complet du mois en cours avec compression",
			Options: ExportOptions{
				Format: FormatExcel,
				Filter: ExportFilter{
					DateFrom:      &[]time.Time{time.Now().AddDate(0, -1, 0)}[0],
					RecentFirst:   true,
					IncludeRemake: false,
				},
				Compression: true,
				Metadata:    true,
			},
			CreatedAt: time.Now(),
			IsDefault: true,
		},
		{
			ID:          "tournament_prep",
			Name:        "Préparation Tournoi",
			Description: "Export optimisé pour l'analyse de performance en vue d'un tournoi",
			Options: ExportOptions{
				Format: FormatJSON,
				Filter: ExportFilter{
					RankedOnly:    &[]bool{true}[0],
					RecentFirst:   true,
					IncludeRemake: false,
					DateFrom:      &[]time.Time{time.Now().AddDate(0, 0, -30)}[0], // 30 derniers jours
				},
				Columns: []string{
					"match_id", "game_creation", "champion_name", "role", 
					"win", "kills", "deaths", "assists", "kda", "cs", 
					"gold", "damage", "vision", "items", "summoners",
				},
				Compression: true,
				Metadata:    true,
			},
			CreatedAt: time.Now(),
			IsDefault: true,
		},
	}

	for _, template := range defaultTemplates {
		ts.templates[template.ID] = template
	}
}

// GetTemplate récupère un template par son ID
func (ts *TemplateService) GetTemplate(id string) (ExportTemplate, bool) {
	template, exists := ts.templates[id]
	return template, exists
}

// GetAllTemplates retourne tous les templates disponibles
func (ts *TemplateService) GetAllTemplates() []ExportTemplate {
	templates := make([]ExportTemplate, 0, len(ts.templates))
	for _, template := range ts.templates {
		templates = append(templates, template)
	}
	return templates
}

// GetDefaultTemplates retourne uniquement les templates par défaut
func (ts *TemplateService) GetDefaultTemplates() []ExportTemplate {
	var defaults []ExportTemplate
	for _, template := range ts.templates {
		if template.IsDefault {
			defaults = append(defaults, template)
		}
	}
	return defaults
}

// CreateTemplate crée un nouveau template personnalisé
func (ts *TemplateService) CreateTemplate(template ExportTemplate) error {
	if template.ID == "" {
		return fmt.Errorf("l'ID du template est requis")
	}
	
	if _, exists := ts.templates[template.ID]; exists {
		return fmt.Errorf("un template avec cet ID existe déjà")
	}
	
	template.CreatedAt = time.Now()
	template.IsDefault = false
	ts.templates[template.ID] = template
	
	return nil
}

// UpdateTemplate met à jour un template existant
func (ts *TemplateService) UpdateTemplate(id string, template ExportTemplate) error {
	existing, exists := ts.templates[id]
	if !exists {
		return fmt.Errorf("template non trouvé")
	}
	
	// Ne pas permettre la modification des templates par défaut
	if existing.IsDefault {
		return fmt.Errorf("impossible de modifier un template par défaut")
	}
	
	template.ID = id
	template.CreatedAt = existing.CreatedAt
	template.IsDefault = false
	ts.templates[id] = template
	
	return nil
}

// DeleteTemplate supprime un template
func (ts *TemplateService) DeleteTemplate(id string) error {
	template, exists := ts.templates[id]
	if !exists {
		return fmt.Errorf("template non trouvé")
	}
	
	// Ne pas permettre la suppression des templates par défaut
	if template.IsDefault {
		return fmt.Errorf("impossible de supprimer un template par défaut")
	}
	
	delete(ts.templates, id)
	return nil
}

// ApplyTemplate applique un template à des options d'export
func (ts *TemplateService) ApplyTemplate(templateID string, baseOptions ExportOptions) (ExportOptions, error) {
	template, exists := ts.templates[templateID]
	if !exists {
		return baseOptions, fmt.Errorf("template non trouvé")
	}
	
	// Commencer avec les options du template
	result := template.Options
	
	// Permettre de surcharger certains champs depuis les options de base
	if baseOptions.Filename != "" {
		result.Filename = baseOptions.Filename
	}
	
	// Fusionner les filtres de date si spécifiés dans les options de base
	if baseOptions.Filter.DateFrom != nil {
		result.Filter.DateFrom = baseOptions.Filter.DateFrom
	}
	if baseOptions.Filter.DateTo != nil {
		result.Filter.DateTo = baseOptions.Filter.DateTo
	}
	
	return result, nil
}

// GetTemplateByCategory retourne les templates d'une catégorie donnée
func (ts *TemplateService) GetTemplateByCategory(category string) []ExportTemplate {
	var result []ExportTemplate
	
	for _, template := range ts.templates {
		// Logique de catégorisation basée sur les options
		switch category {
		case "analysis":
			if template.Options.Format == FormatJSON || 
			   template.Options.Metadata {
				result = append(result, template)
			}
		case "reporting":
			if template.Options.Format == FormatExcel || 
			   template.Options.Format == FormatCSV {
				result = append(result, template)
			}
		case "ranked":
			if template.Options.Filter.RankedOnly != nil && 
			   *template.Options.Filter.RankedOnly {
				result = append(result, template)
			}
		case "compressed":
			if template.Options.Compression {
				result = append(result, template)
			}
		}
	}
	
	return result
}

// ValidateTemplate valide qu'un template est correct
func (ts *TemplateService) ValidateTemplate(template ExportTemplate) error {
	if template.ID == "" {
		return fmt.Errorf("l'ID du template est requis")
	}
	
	if template.Name == "" {
		return fmt.Errorf("le nom du template est requis")
	}
	
	// Utiliser le service d'export pour valider les options
	exportService := NewExportService("./temp")
	return exportService.ValidateOptions(template.Options)
}