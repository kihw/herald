import html2canvas from 'html2canvas';
import * as XLSX from 'xlsx';
import { Row } from '../types';

export interface ExportOptions {
  filename?: string;
  quality?: number;
  scale?: number;
}

export class ExportService {
  /**
   * Exporte un élément DOM en image PNG
   */
  static async exportToPNG(
    element: HTMLElement,
    options: ExportOptions = {}
  ): Promise<void> {
    const {
      filename = 'lol-analytics-export',
      quality = 0.95,
      scale = 2
    } = options;

    try {
      // Configuration pour un rendu de haute qualité
      const canvas = await html2canvas(element, {
        backgroundColor: '#0F2027', // Background sombre par défaut
        scale,
        useCORS: true,
        allowTaint: true,
        foreignObjectRendering: true,
        imageTimeout: 30000,
        removeContainer: true,
        scrollX: 0,
        scrollY: 0,
        windowWidth: element.scrollWidth,
        windowHeight: element.scrollHeight,
      });

      // Création du lien de téléchargement
      const link = document.createElement('a');
      link.download = `${filename}.png`;
      link.href = canvas.toDataURL('image/png', quality);
      
      // Déclenchement du téléchargement
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);

      // Nettoyage
      canvas.remove();
    } catch (error) {
      console.error('Erreur lors de l\'export PNG:', error);
      throw new Error('Impossible d\'exporter en PNG');
    }
  }

  /**
   * Exporte des données en fichier Excel
   */
  static async exportToExcel<T extends Record<string, any>>(
    data: T[],
    options: ExportOptions & {
      sheetName?: string;
      columns?: Array<{
        key: keyof T;
        label: string;
        formatter?: (value: any) => string | number;
      }>;
    } = {}
  ): Promise<void> {
    const {
      filename = 'lol-analytics-data',
      sheetName = 'Données',
      columns
    } = options;

    try {
      // Préparation des données
      let processedData: any[];

      if (columns) {
        // Utiliser les colonnes spécifiées avec formatage
        processedData = data.map(row => {
          const processedRow: Record<string, any> = {};
          columns.forEach(col => {
            const value = row[col.key];
            processedRow[col.label] = col.formatter ? col.formatter(value) : value;
          });
          return processedRow;
        });
      } else {
        // Utiliser toutes les données telles quelles
        processedData = data.map(row => ({ ...row }));
      }

      // Création du workbook
      const wb = XLSX.utils.book_new();
      const ws = XLSX.utils.json_to_sheet(processedData);

      // Configuration des colonnes (largeurs automatiques)
      const colWidths = Object.keys(processedData[0] || {}).map(key => ({
        wch: Math.max(
          key.length,
          ...processedData.map(row => String(row[key] || '').length)
        )
      }));
      ws['!cols'] = colWidths;

      // Ajout de la feuille au workbook
      XLSX.utils.book_append_sheet(wb, ws, sheetName);

      // Export du fichier
      XLSX.writeFile(wb, `${filename}.xlsx`);
    } catch (error) {
      console.error('Erreur lors de l\'export Excel:', error);
      throw new Error('Impossible d\'exporter en Excel');
    }
  }

  /**
   * Exporte les données des rôles en Excel
   */
  static async exportRolesToExcel(data: Row[]): Promise<void> {
    const roleStats = this.calculateRoleStats(data);
    
    await this.exportToExcel(roleStats, {
      filename: 'lol-analytics-roles',
      sheetName: 'Statistiques par Rôle',
      columns: [
        { key: 'role', label: 'Rôle' },
        { key: 'games', label: 'Matchs' },
        { key: 'wins', label: 'Victoires' },
        { key: 'winrate', label: 'Taux de victoire (%)', formatter: (v) => `${(v * 100).toFixed(1)}%` },
        { key: 'avgKda', label: 'KDA moyen', formatter: (v) => v.toFixed(2) },
        { key: 'avgKp', label: 'KP moyen (%)', formatter: (v) => `${(v * 100).toFixed(1)}%` },
        { key: 'avgCsPerMin', label: 'CS/min', formatter: (v) => v.toFixed(1) },
        { key: 'avgGpm', label: 'GPM', formatter: (v) => v.toFixed(0) },
        { key: 'avgDpm', label: 'DPM', formatter: (v) => v.toFixed(0) },
        { key: 'avgVision', label: 'Vision Score', formatter: (v) => v.toFixed(0) },
      ]
    });
  }

  /**
   * Exporte les données des champions en Excel
   */
  static async exportChampionsToExcel(data: Row[], selectedRole?: string): Promise<void> {
    const championStats = this.calculateChampionStats(data);
    
    const filename = selectedRole 
      ? `lol-analytics-champions-${selectedRole.toLowerCase()}`
      : 'lol-analytics-champions';

    await this.exportToExcel(championStats, {
      filename,
      sheetName: `Champions${selectedRole ? ` - ${selectedRole}` : ''}`,
      columns: [
        { key: 'champion', label: 'Champion' },
        { key: 'games', label: 'Matchs' },
        { key: 'wins', label: 'Victoires' },
        { key: 'winrate', label: 'Taux de victoire (%)', formatter: (v) => `${(v * 100).toFixed(1)}%` },
        { key: 'avgKda', label: 'KDA moyen', formatter: (v) => v.toFixed(2) },
        { key: 'avgKp', label: 'KP moyen (%)', formatter: (v) => `${(v * 100).toFixed(1)}%` },
        { key: 'avgCsPerMin', label: 'CS/min', formatter: (v) => v.toFixed(1) },
        { key: 'avgGpm', label: 'GPM', formatter: (v) => v.toFixed(0) },
        { key: 'avgDpm', label: 'DPM', formatter: (v) => v.toFixed(0) },
        { key: 'avgVision', label: 'Vision Score', formatter: (v) => v.toFixed(0) },
        { key: 'recentWinrate', label: 'Forme récente (%)', formatter: (v) => `${(v * 100).toFixed(0)}%` },
      ]
    });
  }

  /**
   * Exporte un export combiné (PNG + Excel) d'une vue
   */
  static async exportCombined(
    element: HTMLElement,
    data: Row[],
    type: 'roles' | 'champions',
    options: ExportOptions & { selectedRole?: string } = {}
  ): Promise<void> {
    const { filename = `lol-analytics-${type}`, selectedRole } = options;

    // Export PNG
    await this.exportToPNG(element, { 
      filename: `${filename}-graphiques`,
      ...options 
    });

    // Export Excel selon le type
    if (type === 'roles') {
      await this.exportRolesToExcel(data);
    } else {
      await this.exportChampionsToExcel(data, selectedRole);
    }
  }

  /**
   * Calcule les statistiques par rôle
   */
  private static calculateRoleStats(data: Row[]) {
    const stats = data.reduce((acc, row) => {
      const role = row.lane || 'Unknown';
      if (!acc[role]) {
        acc[role] = {
          role,
          games: 0,
          wins: 0,
          kdaSum: 0,
          kpSum: 0,
          csPerMinSum: 0,
          gpmSum: 0,
          dpmSum: 0,
          visionSum: 0,
        };
      }

      const stat = acc[role];
      stat.games++;
      if (row.win) stat.wins++;
      if (typeof row.kda === 'number') stat.kdaSum += row.kda;
      if (typeof row.kp === 'number') stat.kpSum += row.kp;
      if (typeof row.cs_per_min === 'number') stat.csPerMinSum += row.cs_per_min;
      if (typeof row.gpm === 'number') stat.gpmSum += row.gpm;
      if (typeof row.dpm === 'number') stat.dpmSum += row.dpm;
      if (typeof row.vision_score === 'number') stat.visionSum += row.vision_score;

      return acc;
    }, {} as Record<string, any>);

    return Object.values(stats).map((stat: any) => ({
      role: stat.role,
      games: stat.games,
      wins: stat.wins,
      winrate: stat.wins / stat.games,
      avgKda: stat.kdaSum / stat.games,
      avgKp: stat.kpSum / stat.games,
      avgCsPerMin: stat.csPerMinSum / stat.games,
      avgGpm: stat.gpmSum / stat.games,
      avgDpm: stat.dpmSum / stat.games,
      avgVision: stat.visionSum / stat.games,
    })).sort((a, b) => b.games - a.games);
  }

  /**
   * Calcule les statistiques par champion
   */
  private static calculateChampionStats(data: Row[]) {
    const stats = data.reduce((acc, row) => {
      const champion = row.champion || 'Unknown';
      if (!acc[champion]) {
        acc[champion] = {
          champion,
          games: 0,
          wins: 0,
          kdaSum: 0,
          kpSum: 0,
          csPerMinSum: 0,
          gpmSum: 0,
          dpmSum: 0,
          visionSum: 0,
          recentGames: [] as boolean[],
        };
      }

      const stat = acc[champion];
      stat.games++;
      if (row.win) stat.wins++;
      if (typeof row.kda === 'number') stat.kdaSum += row.kda;
      if (typeof row.kp === 'number') stat.kpSum += row.kp;
      if (typeof row.cs_per_min === 'number') stat.csPerMinSum += row.cs_per_min;
      if (typeof row.gpm === 'number') stat.gpmSum += row.gpm;
      if (typeof row.dpm === 'number') stat.dpmSum += row.dpm;
      if (typeof row.vision_score === 'number') stat.visionSum += row.vision_score;
      stat.recentGames.push(!!row.win);

      return acc;
    }, {} as Record<string, any>);

    return Object.values(stats).map((stat: any) => {
      const recentGames = stat.recentGames.slice(-5);
      const recentWinrate = recentGames.length > 0 ? 
        recentGames.filter(Boolean).length / recentGames.length : 0;

      return {
        champion: stat.champion,
        games: stat.games,
        wins: stat.wins,
        winrate: stat.wins / stat.games,
        avgKda: stat.kdaSum / stat.games,
        avgKp: stat.kpSum / stat.games,
        avgCsPerMin: stat.csPerMinSum / stat.games,
        avgGpm: stat.gpmSum / stat.games,
        avgDpm: stat.dpmSum / stat.games,
        avgVision: stat.visionSum / stat.games,
        recentWinrate,
      };
    }).sort((a, b) => b.games - a.games);
  }
}