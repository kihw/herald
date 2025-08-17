import { useState, useCallback } from 'react';
import { ExportService, ExportOptions } from '../services/ExportService';
import { Row } from '../types';

export interface UseExportOptions {
  onSuccess?: (message: string) => void;
  onError?: (error: string) => void;
}

export const useExport = (options: UseExportOptions = {}) => {
  const [isExporting, setIsExporting] = useState(false);
  const [exportProgress, setExportProgress] = useState(0);

  const { onSuccess, onError } = options;

  const exportToPNG = useCallback(async (
    elementId: string,
    exportOptions: ExportOptions = {}
  ) => {
    setIsExporting(true);
    setExportProgress(0);

    try {
      const element = document.getElementById(elementId);
      if (!element) {
        throw new Error(`Élément avec l'ID "${elementId}" non trouvé`);
      }

      setExportProgress(25);

      // Petite pause pour que l'utilisateur voit le feedback
      await new Promise(resolve => setTimeout(resolve, 100));
      setExportProgress(50);

      await ExportService.exportToPNG(element, exportOptions);
      
      setExportProgress(100);
      onSuccess?.('Image PNG exportée avec succès');
    } catch (error) {
      console.error('Erreur export PNG:', error);
      onError?.(error instanceof Error ? error.message : 'Erreur lors de l\'export PNG');
    } finally {
      setIsExporting(false);
      setExportProgress(0);
    }
  }, [onSuccess, onError]);

  const exportToExcel = useCallback(async (
    data: any[],
    exportOptions: Parameters<typeof ExportService.exportToExcel>[1] = {}
  ) => {
    setIsExporting(true);
    setExportProgress(0);

    try {
      setExportProgress(25);
      await new Promise(resolve => setTimeout(resolve, 100));
      setExportProgress(50);

      await ExportService.exportToExcel(data, exportOptions);
      
      setExportProgress(100);
      onSuccess?.('Fichier Excel exporté avec succès');
    } catch (error) {
      console.error('Erreur export Excel:', error);
      onError?.(error instanceof Error ? error.message : 'Erreur lors de l\'export Excel');
    } finally {
      setIsExporting(false);
      setExportProgress(0);
    }
  }, [onSuccess, onError]);

  const exportRolesToExcel = useCallback(async (data: Row[]) => {
    setIsExporting(true);
    setExportProgress(0);

    try {
      setExportProgress(25);
      await new Promise(resolve => setTimeout(resolve, 100));
      setExportProgress(50);

      await ExportService.exportRolesToExcel(data);
      
      setExportProgress(100);
      onSuccess?.('Données des rôles exportées en Excel');
    } catch (error) {
      console.error('Erreur export rôles Excel:', error);
      onError?.(error instanceof Error ? error.message : 'Erreur lors de l\'export des rôles');
    } finally {
      setIsExporting(false);
      setExportProgress(0);
    }
  }, [onSuccess, onError]);

  const exportChampionsToExcel = useCallback(async (
    data: Row[], 
    selectedRole?: string
  ) => {
    setIsExporting(true);
    setExportProgress(0);

    try {
      setExportProgress(25);
      await new Promise(resolve => setTimeout(resolve, 100));
      setExportProgress(50);

      await ExportService.exportChampionsToExcel(data, selectedRole);
      
      setExportProgress(100);
      onSuccess?.('Données des champions exportées en Excel');
    } catch (error) {
      console.error('Erreur export champions Excel:', error);
      onError?.(error instanceof Error ? error.message : 'Erreur lors de l\'export des champions');
    } finally {
      setIsExporting(false);
      setExportProgress(0);
    }
  }, [onSuccess, onError]);

  const exportCombined = useCallback(async (
    elementId: string,
    data: Row[],
    type: 'roles' | 'champions',
    exportOptions: ExportOptions & { selectedRole?: string } = {}
  ) => {
    setIsExporting(true);
    setExportProgress(0);

    try {
      const element = document.getElementById(elementId);
      if (!element) {
        throw new Error(`Élément avec l'ID "${elementId}" non trouvé`);
      }

      setExportProgress(20);
      await new Promise(resolve => setTimeout(resolve, 100));

      await ExportService.exportCombined(element, data, type, exportOptions);
      
      setExportProgress(100);
      onSuccess?.('Export combiné (PNG + Excel) réalisé avec succès');
    } catch (error) {
      console.error('Erreur export combiné:', error);
      onError?.(error instanceof Error ? error.message : 'Erreur lors de l\'export combiné');
    } finally {
      setIsExporting(false);
      setExportProgress(0);
    }
  }, [onSuccess, onError]);

  return {
    isExporting,
    exportProgress,
    exportToPNG,
    exportToExcel,
    exportRolesToExcel,
    exportChampionsToExcel,
    exportCombined,
  };
};