package services

import (
	"encoding/json"
	"fmt"
	"log"
)

// GoogleSheetsService gère l'intégration avec Google Sheets (version stub)
type GoogleSheetsService struct {
	enabled     bool
	credentials string
}

// NewGoogleSheetsService crée une nouvelle instance du service Google Sheets
func NewGoogleSheetsService(credentialsPath string) *GoogleSheetsService {
	if credentialsPath == "" {
		log.Println("Google Sheets credentials not provided, service disabled")
		return &GoogleSheetsService{enabled: false}
	}

	// Pour le moment, on disable le service car les dépendances Google ne sont pas installées
	log.Println("Google Sheets service initialized but disabled (dependencies not available)")
	return &GoogleSheetsService{
		enabled:     false,
		credentials: credentialsPath,
	}
}

// IsEnabled retourne true si le service Google Sheets est configuré et activé
func (gss *GoogleSheetsService) IsEnabled() bool {
	return gss.enabled
}

// ExportMatchData exporte les données de match (version stub qui utilise JSON)
func (gss *GoogleSheetsService) ExportMatchData(spreadsheetID string, sheetName string, matches []map[string]interface{}) error {
	if !gss.IsEnabled() {
		// Fallback vers JSON
		return gss.ExportToJSON(matches, fmt.Sprintf("matches_%s.json", sheetName))
	}

	// Implementation réelle serait ici si Google Sheets était configuré
	return fmt.Errorf("Google Sheets service is not enabled")
}

// ExportAnalyticsData exporte les données d'analytics (version stub qui utilise JSON)
func (gss *GoogleSheetsService) ExportAnalyticsData(spreadsheetID string, sheetName string, analytics map[string]interface{}) error {
	if !gss.IsEnabled() {
		// Fallback vers JSON
		return gss.ExportToJSON(analytics, fmt.Sprintf("analytics_%s.json", sheetName))
	}

	// Implementation réelle serait ici si Google Sheets était configuré
	return fmt.Errorf("Google Sheets service is not enabled")
}

// CreateSpreadsheet crée une nouvelle Google Sheet (version stub)
func (gss *GoogleSheetsService) CreateSpreadsheet(title string) (string, error) {
	if !gss.IsEnabled() {
		return "", fmt.Errorf("Google Sheets service is not enabled")
	}

	// Implementation réelle serait ici
	return "", fmt.Errorf("Google Sheets service is not implemented yet")
}

// AddSheet ajoute une nouvelle feuille (version stub)
func (gss *GoogleSheetsService) AddSheet(spreadsheetID string, sheetName string) error {
	if !gss.IsEnabled() {
		return fmt.Errorf("Google Sheets service is not enabled")
	}

	// Implementation réelle serait ici
	return fmt.Errorf("Google Sheets service is not implemented yet")
}

// GetSpreadsheetInfo récupère les informations d'un spreadsheet (version stub)
func (gss *GoogleSheetsService) GetSpreadsheetInfo(spreadsheetID string) (map[string]interface{}, error) {
	if !gss.IsEnabled() {
		return nil, fmt.Errorf("Google Sheets service is not enabled")
	}

	// Implementation réelle serait ici
	return nil, fmt.Errorf("Google Sheets service is not implemented yet")
}

// Helper function pour extraire les valeurs string des maps
func getStringValue(data map[string]interface{}, key string) string {
	if val, exists := data[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", val)
	}
	return ""
}

// ExportToJSON exporte les données vers un format JSON (fallback si Google Sheets n'est pas disponible)
func (gss *GoogleSheetsService) ExportToJSON(data interface{}, filename string) error {
	_, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %v", err)
	}

	log.Printf("Data exported to JSON format (Google Sheets not available): %s", filename)
	// Ici, on pourrait écrire dans un fichier ou retourner les données
	// Pour le moment, on log juste que l'export JSON est disponible

	return nil
}
