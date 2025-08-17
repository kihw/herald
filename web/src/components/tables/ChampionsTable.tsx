import React, { useMemo, useState } from 'react';
import {
  DataGrid,
  GridColDef,
  GridRenderCellParams,
  GridToolbar,
} from '@mui/x-data-grid';
import {
  Box,
  Card,
  CardContent,
  CardHeader,
  Chip,
  IconButton,
  Tooltip,
  Avatar,
  Typography,
} from '@mui/material';
import {
  GetApp,
  TableChart,
  BarChart,
  Visibility,
  TrendingUp,
  TrendingDown,
} from '@mui/icons-material';
import { Row } from '../../types';

interface ChampionsTableProps {
  data: Row[];
  selectedRole?: string;
  onChampionSelect: (champion: string) => void;
  onExportPNG?: () => void;
  onExportExcel?: () => void;
}

interface ChampionStats {
  id: string;
  champion: string;
  games: number;
  wins: number;
  winrate: number;
  avgKda: number;
  avgKp: number;
  avgCsPerMin: number;
  avgGpm: number;
  avgDpm: number;
  avgVision: number;
  avgKills: number;
  avgDeaths: number;
  avgAssists: number;
  lastPlayed: string;
  trend: 'up' | 'down' | 'stable';
  recentWinrate: number;
}

// Fonction pour obtenir l'URL de l'icône champion Data Dragon
const getChampionIconUrl = (championName: string, patch = '14.23.1') => {
  const champId = championName.replace(/[^A-Za-z0-9]/g, '');
  return `https://ddragon.leagueoflegends.com/cdn/${patch}/img/champion/${champId}.png`;
};

// Composant icône champion
const ChampionIcon: React.FC<{ champion: string; size?: number }> = ({ 
  champion, 
  size = 40 
}) => (
  <Avatar
    src={getChampionIconUrl(champion)}
    alt={champion}
    sx={{ 
      width: size, 
      height: size, 
      border: '2px solid',
      borderColor: 'primary.main',
    }}
  >
    {champion.slice(0, 2).toUpperCase()}
  </Avatar>
);

export const ChampionsTable: React.FC<ChampionsTableProps> = ({
  data,
  selectedRole,
  onChampionSelect,
  onExportPNG,
  onExportExcel,
}) => {
  const [loading, setLoading] = useState(false);

  const championStats = useMemo((): ChampionStats[] => {
    const stats = data.reduce((acc, row) => {
      const champion = row.champion || 'Unknown';
      if (!acc[champion]) {
        acc[champion] = {
          games: 0,
          wins: 0,
          kdaSum: 0,
          kpSum: 0,
          csPerMinSum: 0,
          gpmSum: 0,
          dpmSum: 0,
          visionSum: 0,
          killsSum: 0,
          deathsSum: 0,
          assistsSum: 0,
          dates: [] as string[],
          recentGames: [] as boolean[], // Pour calculer la tendance
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
      if (typeof row.kills === 'number') stat.killsSum += row.kills;
      if (typeof row.deaths === 'number') stat.deathsSum += row.deaths;
      if (typeof row.assists === 'number') stat.assistsSum += row.assists;

      if (row.date) stat.dates.push(row.date);
      stat.recentGames.push(!!row.win);

      return acc;
    }, {} as Record<string, any>);

    return Object.entries(stats).map(([champion, stat]) => {
      // Calculer la tendance (winrate des 5 derniers matchs vs global)
      const recentGames = stat.recentGames.slice(-5);
      const recentWinrate = recentGames.length > 0 ? 
        recentGames.filter(Boolean).length / recentGames.length : 0;
      const globalWinrate = stat.wins / stat.games;
      
      let trend: 'up' | 'down' | 'stable' = 'stable';
      if (recentWinrate > globalWinrate + 0.1) trend = 'up';
      else if (recentWinrate < globalWinrate - 0.1) trend = 'down';

      // Date du dernier match
      const lastPlayed = stat.dates.length > 0 ? 
        new Date(Math.max(...stat.dates.map((d: string) => new Date(d).getTime()))).toLocaleDateString('fr-FR') : 
        'N/A';

      return {
        id: champion,
        champion,
        games: stat.games,
        wins: stat.wins,
        winrate: stat.wins / stat.games,
        avgKda: stat.kdaSum / stat.games,
        avgKp: stat.kpSum / stat.games,
        avgCsPerMin: stat.csPerMinSum / stat.games,
        avgGpm: stat.gpmSum / stat.games,
        avgDpm: stat.dpmSum / stat.games,
        avgVision: stat.visionSum / stat.games,
        avgKills: stat.killsSum / stat.games,
        avgDeaths: stat.deathsSum / stat.games,
        avgAssists: stat.assistsSum / stat.games,
        lastPlayed,
        trend,
        recentWinrate,
      };
    }).sort((a, b) => b.games - a.games); // Trier par nombre de games
  }, [data]);

  const columns: GridColDef[] = [
    {
      field: 'champion',
      headerName: 'Champion',
      width: 200,
      renderCell: (params: GridRenderCellParams) => (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <ChampionIcon champion={params.value} size={32} />
          <Box>
            <Typography variant="body2" sx={{ fontWeight: 600 }}>
              {params.value}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              {params.row.games} games • Dernier: {params.row.lastPlayed}
            </Typography>
          </Box>
        </Box>
      ),
    },
    {
      field: 'winrate',
      headerName: 'Taux de victoire',
      width: 140,
      renderCell: (params: GridRenderCellParams) => (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Chip
            label={`${(params.value * 100).toFixed(1)}%`}
            color={params.value >= 0.5 ? 'success' : 'error'}
            variant="outlined"
            size="small"
          />
          {params.row.trend === 'up' && <TrendingUp color="success" fontSize="small" />}
          {params.row.trend === 'down' && <TrendingDown color="error" fontSize="small" />}
        </Box>
      ),
    },
    {
      field: 'avgKda',
      headerName: 'KDA',
      width: 100,
      renderCell: (params: GridRenderCellParams) => (
        <Box>
          <Typography variant="body2" sx={{ fontWeight: 600 }}>
            {params.value.toFixed(2)}
          </Typography>
          <Typography variant="caption" color="text.secondary">
            {params.row.avgKills.toFixed(1)}/{params.row.avgDeaths.toFixed(1)}/{params.row.avgAssists.toFixed(1)}
          </Typography>
        </Box>
      ),
    },
    {
      field: 'avgKp',
      headerName: 'KP%',
      width: 100,
      renderCell: (params: GridRenderCellParams) => 
        `${(params.value * 100).toFixed(1)}%`,
    },
    {
      field: 'avgCsPerMin',
      headerName: 'CS/min',
      width: 100,
      valueFormatter: (params) => typeof params.value === 'number' ? params.value.toFixed(1) : '0.0',
    },
    {
      field: 'avgGpm',
      headerName: 'GPM',
      width: 100,
      valueFormatter: (params) => typeof params.value === 'number' ? params.value.toFixed(0) : '0',
    },
    {
      field: 'avgDpm',
      headerName: 'DPM',
      width: 100,
      valueFormatter: (params) => typeof params.value === 'number' ? params.value.toFixed(0) : '0',
    },
    {
      field: 'avgVision',
      headerName: 'Vision',
      width: 100,
      valueFormatter: (params) => typeof params.value === 'number' ? params.value.toFixed(0) : '0',
    },
    {
      field: 'recentWinrate',
      headerName: 'Form récente',
      width: 120,
      renderCell: (params: GridRenderCellParams) => (
        <Chip
          label={`${(params.value * 100).toFixed(0)}%`}
          color={params.value >= 0.6 ? 'success' : params.value >= 0.4 ? 'warning' : 'error'}
          variant="outlined"
          size="small"
        />
      ),
    },
    {
      field: 'actions',
      headerName: 'Actions',
      width: 120,
      sortable: false,
      renderCell: (params: GridRenderCellParams) => (
        <Box sx={{ display: 'flex', gap: 0.5 }}>
          <Tooltip title="Détails du champion">
            <IconButton
              size="small"
              onClick={(e) => {
                e.stopPropagation();
                onChampionSelect(params.row.champion);
              }}
            >
              <Visibility />
            </IconButton>
          </Tooltip>
          <Tooltip title="Graphiques">
            <IconButton size="small">
              <BarChart />
            </IconButton>
          </Tooltip>
        </Box>
      ),
    },
  ];

  return (
    <Card>
      <CardHeader
        title={`Champions ${selectedRole ? `- ${selectedRole}` : ''}`}
        subheader={`${championStats.length} champions joués sur ${data.length} matchs`}
        action={
          <Box sx={{ display: 'flex', gap: 1 }}>
            <Tooltip title="Exporter la vue en PNG">
              <IconButton 
                onClick={onExportPNG} 
                size="small" 
                color="primary"
                aria-label="Exporter le tableau des champions en image PNG"
              >
                <GetApp />
              </IconButton>
            </Tooltip>
            <Tooltip title="Exporter les données en Excel">
              <IconButton 
                onClick={onExportExcel} 
                size="small" 
                color="primary"
                aria-label="Exporter les données des champions en fichier Excel"
              >
                <TableChart />
              </IconButton>
            </Tooltip>
          </Box>
        }
      />
      <CardContent sx={{ height: 500, p: 0 }}>
        <DataGrid
          rows={championStats}
          columns={columns}
          loading={loading}
          disableRowSelectionOnClick={false}
          onRowClick={(params) => onChampionSelect(params.row.champion)}
          slots={{
            toolbar: GridToolbar,
          }}
          slotProps={{
            toolbar: {
              showQuickFilter: true,
              quickFilterProps: { debounceMs: 500 },
            },
          }}
          sx={{
            border: 'none',
            '& .MuiDataGrid-row': {
              cursor: 'pointer',
              '&:hover': {
                backgroundColor: 'action.hover',
              },
            },
            '& .MuiDataGrid-cell:focus': {
              outline: 'none',
            },
            '& .MuiDataGrid-columnHeaders': {
              backgroundColor: 'background.paper',
              borderBottom: '2px solid',
              borderColor: 'divider',
            },
          }}
          initialState={{
            pagination: {
              paginationModel: { pageSize: 15 },
            },
            sorting: {
              sortModel: [{ field: 'games', sort: 'desc' }],
            },
          }}
          pageSizeOptions={[10, 15, 25, 50]}
          density="comfortable"
        />
      </CardContent>
    </Card>
  );
};