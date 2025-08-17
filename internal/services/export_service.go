package services

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExportFormat définit les formats d'export supportés
type ExportFormat string

const (
	FormatCSV     ExportFormat = "csv"
	FormatJSON    ExportFormat = "json"
	FormatParquet ExportFormat = "parquet"
	FormatExcel   ExportFormat = "xlsx"
)

// ExportFilter définit les filtres pour l'export
type ExportFilter struct {
	DateFrom      *time.Time `json:"date_from,omitempty"`
	DateTo        *time.Time `json:"date_to,omitempty"`
	Queues        []int      `json:"queues,omitempty"`
	Champions     []int      `json:"champions,omitempty"`
	MinDuration   *int       `json:"min_duration,omitempty"`
	MaxDuration   *int       `json:"max_duration,omitempty"`
	WinOnly       *bool      `json:"win_only,omitempty"`
	RankedOnly    *bool      `json:"ranked_only,omitempty"`
	RecentFirst   bool       `json:"recent_first"`
	IncludeRemake bool       `json:"include_remake"`
}

// ExportOptions définit les options d'export
type ExportOptions struct {
	Format      ExportFormat  `json:"format"`
	Filter      ExportFilter  `json:"filter"`
	Columns     []string      `json:"columns,omitempty"`
	Filename    string        `json:"filename"`
	Compression bool          `json:"compression"`
	Metadata    bool          `json:"metadata"`
}

// MatchData représente les données d'un match pour l'export
type MatchData struct {
	MatchID       string    `json:"match_id"`
	GameCreation  time.Time `json:"game_creation"`
	GameDuration  int       `json:"game_duration"`
	QueueID       int       `json:"queue_id"`
	GameMode      string    `json:"game_mode"`
	GameType      string    `json:"game_type"`
	ChampionID    int       `json:"champion_id"`
	ChampionName  string    `json:"champion_name"`
	Role          string    `json:"role"`
	Lane          string    `json:"lane"`
	Win           bool      `json:"win"`
	Kills         int       `json:"kills"`
	Deaths        int       `json:"deaths"`
	Assists       int       `json:"assists"`
	KDA           float64   `json:"kda"`
	CS            int       `json:"cs"`
	Gold          int       `json:"gold"`
	Damage        int       `json:"damage"`
	Vision        int       `json:"vision"`
	Items         []int     `json:"items"`
	Summoners     []int     `json:"summoners"`
	Rank          string    `json:"rank,omitempty"`
	LP            *int      `json:"lp,omitempty"`
	MMR           *int      `json:"mmr,omitempty"`
}

// ExportService gère les exports avancés
type ExportService struct {
	outputDir string
}

// NewExportService crée une nouvelle instance du service d'export
func NewExportService(outputDir string) *ExportService {
	return &ExportService{
		outputDir: outputDir,
	}
}

// ExportMatches exporte les données de matchs selon les options spécifiées
func (s *ExportService) ExportMatches(matches []MatchData, options ExportOptions) (string, error) {
	// Créer le répertoire de sortie si nécessaire
	if err := os.MkdirAll(s.outputDir, 0755); err != nil {
		return "", fmt.Errorf("impossible de créer le répertoire de sortie: %w", err)
	}

	// Filtrer les données
	filteredMatches := s.filterMatches(matches, options.Filter)

	// Générer le nom de fichier
	filename := s.generateFilename(options)
	filepath := filepath.Join(s.outputDir, filename)

	// Exporter selon le format
	var err error
	switch options.Format {
	case FormatCSV:
		err = s.exportToCSV(filteredMatches, filepath, options)
	case FormatJSON:
		err = s.exportToJSON(filteredMatches, filepath, options)
	case FormatParquet:
		err = s.exportToParquet(filteredMatches, filepath, options)
	case FormatExcel:
		err = s.exportToExcel(filteredMatches, filepath, options)
	default:
		return "", fmt.Errorf("format d'export non supporté: %s", options.Format)
	}

	if err != nil {
		return "", err
	}

	// Compresser si demandé
	if options.Compression {
		compressedPath, err := s.compressFile(filepath)
		if err != nil {
			return filepath, nil // Retourner le fichier non compressé en cas d'erreur
		}
		return compressedPath, nil
	}

	return filepath, nil
}

// filterMatches applique les filtres aux données de match
func (s *ExportService) filterMatches(matches []MatchData, filter ExportFilter) []MatchData {
	var filtered []MatchData

	for _, match := range matches {
		// Filtre par date
		if filter.DateFrom != nil && match.GameCreation.Before(*filter.DateFrom) {
			continue
		}
		if filter.DateTo != nil && match.GameCreation.After(*filter.DateTo) {
			continue
		}

		// Filtre par queue
		if len(filter.Queues) > 0 {
			found := false
			for _, queueID := range filter.Queues {
				if match.QueueID == queueID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filtre par champion
		if len(filter.Champions) > 0 {
			found := false
			for _, championID := range filter.Champions {
				if match.ChampionID == championID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filtre par durée
		if filter.MinDuration != nil && match.GameDuration < *filter.MinDuration {
			continue
		}
		if filter.MaxDuration != nil && match.GameDuration > *filter.MaxDuration {
			continue
		}

		// Filtre win only
		if filter.WinOnly != nil && *filter.WinOnly && !match.Win {
			continue
		}

		// Filtre ranked only
		if filter.RankedOnly != nil && *filter.RankedOnly && !s.isRankedQueue(match.QueueID) {
			continue
		}

		// Exclure les remakes si nécessaire
		if !filter.IncludeRemake && match.GameDuration < 300 { // < 5 minutes
			continue
		}

		filtered = append(filtered, match)
	}

	return filtered
}

// exportToCSV exporte les données au format CSV
func (s *ExportService) exportToCSV(matches []MatchData, filepath string, options ExportOptions) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("impossible de créer le fichier CSV: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Déterminer les colonnes à exporter
	columns := options.Columns
	if len(columns) == 0 {
		columns = s.getDefaultColumns()
	}

	// Écrire l'en-tête
	if err := writer.Write(columns); err != nil {
		return fmt.Errorf("impossible d'écrire l'en-tête CSV: %w", err)
	}

	// Écrire les données
	for _, match := range matches {
		row := s.matchToCSVRow(match, columns)
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("impossible d'écrire la ligne CSV: %w", err)
		}
	}

	return nil
}

// exportToJSON exporte les données au format JSON
func (s *ExportService) exportToJSON(matches []MatchData, filepath string, options ExportOptions) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("impossible de créer le fichier JSON: %w", err)
	}
	defer file.Close()

	// Créer la structure d'export avec métadonnées
	exportData := map[string]interface{}{
		"matches": matches,
		"count":   len(matches),
	}

	if options.Metadata {
		exportData["metadata"] = map[string]interface{}{
			"exported_at": time.Now().UTC(),
			"format":      options.Format,
			"filter":      options.Filter,
		}
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	if err := encoder.Encode(exportData); err != nil {
		return fmt.Errorf("impossible d'encoder le JSON: %w", err)
	}

	return nil
}

// exportToParquet exporte les données au format Parquet (placeholder)
func (s *ExportService) exportToParquet(matches []MatchData, filepath string, options ExportOptions) error {
	// Pour l'instant, nous exportons en JSON compressé comme alternative
	// Une vraie implémentation Parquet nécessiterait une bibliothèque spécialisée
	return s.exportToJSON(matches, filepath, options)
}

// exportToExcel exporte les données au format Excel
func (s *ExportService) exportToExcel(matches []MatchData, filepath string, options ExportOptions) error {
	// Créer un nouveau fichier Excel
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Erreur lors de la fermeture du fichier Excel: %v\n", err)
		}
	}()

	// Déterminer les colonnes à exporter
	columns := options.Columns
	if len(columns) == 0 {
		columns = s.getDefaultColumns()
	}

	// Créer la feuille principale
	sheetName := "Matches"
	f.SetSheetName("Sheet1", sheetName)

	// Style pour l'en-tête
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E6E6FA"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("impossible de créer le style d'en-tête: %w", err)
	}

	// Style pour les données
	dataStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "CCCCCC", Style: 1},
			{Type: "top", Color: "CCCCCC", Style: 1},
			{Type: "bottom", Color: "CCCCCC", Style: 1},
			{Type: "right", Color: "CCCCCC", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("impossible de créer le style de données: %w", err)
	}

	// Écrire l'en-tête
	for i, column := range columns {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, s.formatColumnName(column))
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Écrire les données
	for rowIndex, match := range matches {
		row := s.matchToCSVRow(match, columns) // Réutiliser la logique de conversion
		for colIndex, value := range row {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			
			// Convertir les types appropriés
			if colIndex < len(columns) {
				switch columns[colIndex] {
				case "game_duration", "queue_id", "champion_id", "kills", "deaths", "assists", "cs", "gold", "damage", "vision", "lp", "mmr":
					if numValue, err := strconv.Atoi(value); err == nil {
						f.SetCellValue(sheetName, cell, numValue)
					} else {
						f.SetCellValue(sheetName, cell, value)
					}
				case "kda":
					if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
						f.SetCellValue(sheetName, cell, floatValue)
					} else {
						f.SetCellValue(sheetName, cell, value)
					}
				case "win":
					f.SetCellValue(sheetName, cell, value == "true")
				default:
					f.SetCellValue(sheetName, cell, value)
				}
			} else {
				f.SetCellValue(sheetName, cell, value)
			}
			
			f.SetCellStyle(sheetName, cell, cell, dataStyle)
		}
	}

	// Ajuster la largeur des colonnes
	for i := range columns {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheetName, col, col, 15)
	}

	// Ajouter une feuille de métadonnées si demandé
	if options.Metadata {
		metaSheet := "Metadata"
		f.NewSheet(metaSheet)
		
		metadata := [][]interface{}{
			{"Export Date", time.Now().Format("2006-01-02 15:04:05")},
			{"Format", "Excel"},
			{"Total Matches", len(matches)},
			{"Columns", len(columns)},
		}
		
		for i, row := range metadata {
			for j, value := range row {
				cell, _ := excelize.CoordinatesToCellName(j+1, i+1)
				f.SetCellValue(metaSheet, cell, value)
				if j == 0 {
					f.SetCellStyle(metaSheet, cell, cell, headerStyle)
				}
			}
		}
	}

	// Sauvegarder le fichier
	if err := f.SaveAs(filepath); err != nil {
		return fmt.Errorf("impossible de sauvegarder le fichier Excel: %w", err)
	}

	return nil
}

// formatColumnName formate le nom de colonne pour l'affichage
func (s *ExportService) formatColumnName(column string) string {
	// Convertir snake_case en Title Case
	parts := strings.Split(column, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, " ")
}

// generateFilename génère un nom de fichier basé sur les options
func (s *ExportService) generateFilename(options ExportOptions) string {
	if options.Filename != "" {
		return options.Filename
	}

	timestamp := time.Now().Format("20060102_150405")
	extension := string(options.Format)
	
	return fmt.Sprintf("lol_export_%s.%s", timestamp, extension)
}

// getDefaultColumns retourne les colonnes par défaut pour l'export CSV
func (s *ExportService) getDefaultColumns() []string {
	return []string{
		"match_id", "game_creation", "game_duration", "queue_id", "game_mode",
		"champion_name", "role", "lane", "win", "kills", "deaths", "assists",
		"kda", "cs", "gold", "damage", "vision", "rank", "lp", "mmr",
	}
}

// matchToCSVRow convertit un match en ligne CSV
func (s *ExportService) matchToCSVRow(match MatchData, columns []string) []string {
	row := make([]string, len(columns))
	
	for i, column := range columns {
		switch column {
		case "match_id":
			row[i] = match.MatchID
		case "game_creation":
			row[i] = match.GameCreation.Format("2006-01-02 15:04:05")
		case "game_duration":
			row[i] = strconv.Itoa(match.GameDuration)
		case "queue_id":
			row[i] = strconv.Itoa(match.QueueID)
		case "game_mode":
			row[i] = match.GameMode
		case "game_type":
			row[i] = match.GameType
		case "champion_id":
			row[i] = strconv.Itoa(match.ChampionID)
		case "champion_name":
			row[i] = match.ChampionName
		case "role":
			row[i] = match.Role
		case "lane":
			row[i] = match.Lane
		case "win":
			row[i] = strconv.FormatBool(match.Win)
		case "kills":
			row[i] = strconv.Itoa(match.Kills)
		case "deaths":
			row[i] = strconv.Itoa(match.Deaths)
		case "assists":
			row[i] = strconv.Itoa(match.Assists)
		case "kda":
			row[i] = fmt.Sprintf("%.2f", match.KDA)
		case "cs":
			row[i] = strconv.Itoa(match.CS)
		case "gold":
			row[i] = strconv.Itoa(match.Gold)
		case "damage":
			row[i] = strconv.Itoa(match.Damage)
		case "vision":
			row[i] = strconv.Itoa(match.Vision)
		case "items":
			items := make([]string, len(match.Items))
			for j, item := range match.Items {
				items[j] = strconv.Itoa(item)
			}
			row[i] = strings.Join(items, ",")
		case "summoners":
			summoners := make([]string, len(match.Summoners))
			for j, summoner := range match.Summoners {
				summoners[j] = strconv.Itoa(summoner)
			}
			row[i] = strings.Join(summoners, ",")
		case "rank":
			row[i] = match.Rank
		case "lp":
			if match.LP != nil {
				row[i] = strconv.Itoa(*match.LP)
			} else {
				row[i] = ""
			}
		case "mmr":
			if match.MMR != nil {
				row[i] = strconv.Itoa(*match.MMR)
			} else {
				row[i] = ""
			}
		default:
			row[i] = ""
		}
	}
	
	return row
}

// isRankedQueue vérifie si une queue est classée
func (s *ExportService) isRankedQueue(queueID int) bool {
	rankedQueues := []int{420, 440, 470} // Solo/Duo, Flex 5v5, Flex 3v3
	for _, ranked := range rankedQueues {
		if queueID == ranked {
			return true
		}
	}
	return false
}

// GetSupportedFormats retourne la liste des formats supportés
func (s *ExportService) GetSupportedFormats() []ExportFormat {
	return []ExportFormat{FormatCSV, FormatJSON, FormatParquet, FormatExcel}
}

// ValidateOptions valide les options d'export
func (s *ExportService) ValidateOptions(options ExportOptions) error {
	// Vérifier le format
	supportedFormats := s.GetSupportedFormats()
	formatSupported := false
	for _, format := range supportedFormats {
		if options.Format == format {
			formatSupported = true
			break
		}
	}
	if !formatSupported {
		return fmt.Errorf("format non supporté: %s", options.Format)
	}

	// Vérifier les dates
	if options.Filter.DateFrom != nil && options.Filter.DateTo != nil {
		if options.Filter.DateFrom.After(*options.Filter.DateTo) {
			return fmt.Errorf("la date de début doit être antérieure à la date de fin")
		}
	}

	// Vérifier les durées
	if options.Filter.MinDuration != nil && options.Filter.MaxDuration != nil {
		if *options.Filter.MinDuration > *options.Filter.MaxDuration {
			return fmt.Errorf("la durée minimum doit être inférieure à la durée maximum")
		}
	}

	return nil
}

// compressFile compresse un fichier au format ZIP
func (s *ExportService) compressFile(filepath string) (string, error) {
	// Nom du fichier ZIP
	zipPath := filepath + ".zip"
	
	// Créer le fichier ZIP
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("impossible de créer le fichier ZIP: %w", err)
	}
	defer zipFile.Close()

	// Créer le writer ZIP
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Ouvrir le fichier source
	sourceFile, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("impossible d'ouvrir le fichier source: %w", err)
	}
	defer sourceFile.Close()

	// Obtenir les informations du fichier
	info, err := sourceFile.Stat()
	if err != nil {
		return "", fmt.Errorf("impossible d'obtenir les informations du fichier: %w", err)
	}

	// Créer l'en-tête ZIP
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return "", fmt.Errorf("impossible de créer l'en-tête ZIP: %w", err)
	}

	// Utiliser le nom de base du fichier dans le ZIP
	header.Name = filepath[strings.LastIndex(filepath, string(os.PathSeparator))+1:]
	header.Method = zip.Deflate

	// Créer le writer pour ce fichier dans le ZIP
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return "", fmt.Errorf("impossible de créer l'entrée ZIP: %w", err)
	}

	// Copier le contenu
	_, err = io.Copy(writer, sourceFile)
	if err != nil {
		return "", fmt.Errorf("impossible de copier le contenu: %w", err)
	}

	// Supprimer le fichier original
	if err := os.Remove(filepath); err != nil {
		fmt.Printf("Avertissement: impossible de supprimer le fichier original: %v\n", err)
	}

	return zipPath, nil
}