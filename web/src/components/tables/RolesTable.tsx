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
} from '@mui/material';
import {
  GetApp,
  TableChart,
  BarChart,
  Visibility,
} from '@mui/icons-material';
import { Row } from '../../types';

interface RolesTableProps {
  data: Row[];
  onRoleSelect: (role: string) => void;
  onExportPNG?: () => void;
  onExportExcel?: () => void;
}

interface RoleStats {
  id: string;
  role: string;
  games: number;
  wins: number;
  winrate: number;
  avgKda: number;
  avgKp: number;
  avgCsPerMin: number;
  avgGpm: number;
  avgDpm: number;
  avgVision: number;
  avgEarlyCs: number;
  avgEarlyGold: number;
  mostPlayedChampion: string;
  championGames: number;
}

// Icônes de rôle
const RoleIcon: React.FC<{ role: string }> = ({ role }) => {
  const getIcon = () => {
    switch (role.toUpperCase()) {
      case 'TOP': return '⚔️';
      case 'JUNGLE': return '🌳';
      case 'MID': case 'MIDDLE': return '⭐';
      case 'ADC': case 'BOTTOM': return '🏹';
      case 'SUPPORT': return '🛡️';
      default: return '❓';
    }
  };

  return (
    <Avatar 
      sx={{ 
        width: 32, 
        height: 32, 
        backgroundColor: 'primary.main',
        fontSize: '1rem',
      }}
    >
      {getIcon()}
    </Avatar>
  );
};

export const RolesTable: React.FC<RolesTableProps> = ({
  data,
  onRoleSelect,
  onExportPNG,
  onExportExcel,
}) => {
  const [loading, setLoading] = useState(false);

  const roleStats = useMemo((): RoleStats[] => {
    const stats = data.reduce((acc, row) => {
      const role = row.lane || 'Unknown';
      if (!acc[role]) {
        acc[role] = {
          games: 0,
          wins: 0,
          kdaSum: 0,
          kpSum: 0,
          csPerMinSum: 0,
          gpmSum: 0,
          dpmSum: 0,
          visionSum: 0,
          earlyCs10Sum: 0,
          earlyGold10Sum: 0,
          champions: {} as Record<string, number>,
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
      if (typeof row.cs10 === 'number') stat.earlyCs10Sum += row.cs10;
      if (typeof row.gold10 === 'number') stat.earlyGold10Sum += row.gold10;

      // Champion le plus joué
      const champion = row.champion || 'Unknown';
      stat.champions[champion] = (stat.champions[champion] || 0) + 1;

      return acc;
    }, {} as Record<string, any>);

    return Object.entries(stats).map(([role, stat]) => {
      const mostPlayedChampion = Object.entries(stat.champions)
        .sort(([,a], [,b]) => (b as number) - (a as number))[0];

      return {
        id: role,
        role,
        games: stat.games,
        wins: stat.wins,
        winrate: stat.wins / stat.games,
        avgKda: stat.kdaSum / stat.games,
        avgKp: stat.kpSum / stat.games,
        avgCsPerMin: stat.csPerMinSum / stat.games,
        avgGpm: stat.gpmSum / stat.games,
        avgDpm: stat.dpmSum / stat.games,
        avgVision: stat.visionSum / stat.games,
        avgEarlyCs: stat.earlyCs10Sum / stat.games,
        avgEarlyGold: stat.earlyGold10Sum / stat.games,
        mostPlayedChampion: mostPlayedChampion?.[0] || 'N/A',
        championGames: mostPlayedChampion?.[1] as number || 0,
      };
    }).sort((a, b) => b.games - a.games);
  }, [data]);

  const columns: GridColDef[] = [
    {
      field: 'role',
      headerName: 'Rôle',
      width: 150,
      renderCell: (params: GridRenderCellParams) => (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <RoleIcon role={params.value} />
          <Box>
            <Box sx={{ fontWeight: 600 }}>{params.value}</Box>
            <Box sx={{ fontSize: '0.75rem', color: 'text.secondary' }}>
              {params.row.games} games
            </Box>
          </Box>
        </Box>
      ),
    },
    {
      field: 'winrate',
      headerName: 'Taux de victoire',
      width: 130,
      renderCell: (params: GridRenderCellParams) => (
        <Chip
          label={`${(params.value * 100).toFixed(1)}%`}
          color={params.value >= 0.5 ? 'success' : 'error'}
          variant="outlined"
          size="small"
        />
      ),
    },
    {
      field: 'avgKda',
      headerName: 'KDA moy.',
      width: 100,
      valueFormatter: (params) => typeof params.value === 'number' ? params.value.toFixed(2) : '0.00',
    },
    {
      field: 'avgKp',
      headerName: 'KP moy.',
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
      field: 'avgEarlyCs',
      headerName: 'CS@10',
      width: 100,
      valueFormatter: (params) => typeof params.value === 'number' ? params.value.toFixed(1) : '0.0',
    },
    {
      field: 'avgEarlyGold',
      headerName: 'Gold@10',
      width: 100,
      valueFormatter: (params) => typeof params.value === 'number' ? params.value.toFixed(0) : '0',
    },
    {
      field: 'mostPlayedChampion',
      headerName: 'Champion principal',
      width: 160,
      renderCell: (params: GridRenderCellParams) => (
        <Box>
          <Box sx={{ fontWeight: 500 }}>{params.value}</Box>
          <Box sx={{ fontSize: '0.75rem', color: 'text.secondary' }}>
            {params.row.championGames} games
          </Box>
        </Box>
      ),
    },
    {
      field: 'actions',
      headerName: 'Actions',
      width: 120,
      sortable: false,
      renderCell: (params: GridRenderCellParams) => (
        <Box sx={{ display: 'flex', gap: 0.5 }}>
          <Tooltip title="Voir les champions">
            <IconButton
              size="small"
              onClick={(e) => {
                e.stopPropagation();
                onRoleSelect(params.row.role);
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
        title="Performance par Rôle"
        subheader={`Analyse de ${data.length} matchs répartis sur ${roleStats.length} rôles`}
        action={
          <Box sx={{ display: 'flex', gap: 1 }}>
            <Tooltip title="Exporter la vue en PNG">
              <IconButton 
                onClick={onExportPNG} 
                size="small" 
                color="primary"
                aria-label="Exporter le tableau des rôles en image PNG"
              >
                <GetApp />
              </IconButton>
            </Tooltip>
            <Tooltip title="Exporter les données en Excel">
              <IconButton 
                onClick={onExportExcel} 
                size="small" 
                color="primary"
                aria-label="Exporter les données des rôles en fichier Excel"
              >
                <TableChart />
              </IconButton>
            </Tooltip>
          </Box>
        }
      />
      <CardContent sx={{ height: 400, p: 0 }}>
        <DataGrid
          rows={roleStats}
          columns={columns}
          loading={loading}
          disableRowSelectionOnClick={false}
          onRowClick={(params) => onRoleSelect(params.row.role)}
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
              paginationModel: { pageSize: 10 },
            },
            sorting: {
              sortModel: [{ field: 'games', sort: 'desc' }],
            },
          }}
          pageSizeOptions={[5, 10, 25]}
          density="comfortable"
        />
      </CardContent>
    </Card>
  );
};