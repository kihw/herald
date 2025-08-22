import React, { useState } from 'react';
import {
  Box,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Typography,
  Divider,
  Chip,
  Collapse,
  useTheme,
} from '@mui/material';
import {
  ExpandLess,
  ExpandMore,
  SportsEsports,
} from '@mui/icons-material';
import { leagueColors } from '../../theme/leagueTheme';

interface NavigationItem {
  id: string;
  label: string;
  icon: React.ReactNode;
  path: string;
  color: string;
  badge?: string;
  children?: NavigationItem[];
}

interface NavigationSidebarProps {
  items: NavigationItem[];
  onItemClick?: (item: NavigationItem) => void;
}

const NavigationSidebar: React.FC<NavigationSidebarProps> = ({ items, onItemClick }) => {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === 'dark';
  const [expandedItems, setExpandedItems] = useState<string[]>([]);
  const [selectedItem, setSelectedItem] = useState<string>('dashboard');

  const handleItemClick = (item: NavigationItem) => {
    if (item.children && item.children.length > 0) {
      // Toggle expansion for items with children
      setExpandedItems(prev => 
        prev.includes(item.id) 
          ? prev.filter(id => id !== item.id)
          : [...prev, item.id]
      );
    } else {
      // Select item and call callback
      setSelectedItem(item.id);
      onItemClick?.(item);
    }
  };

  const isExpanded = (itemId: string) => expandedItems.includes(itemId);

  const renderNavigationItem = (item: NavigationItem, depth = 0) => {
    const isSelected = selectedItem === item.id;
    const hasChildren = item.children && item.children.length > 0;
    const expanded = isExpanded(item.id);

    return (
      <React.Fragment key={item.id}>
        <ListItem disablePadding sx={{ pl: depth * 2 }}>
          <ListItemButton
            onClick={() => handleItemClick(item)}
            selected={isSelected}
            sx={{
              borderRadius: 2,
              mx: 1,
              mb: 0.5,
              py: 1.5,
              '&.Mui-selected': {
                background: `linear-gradient(135deg, ${item.color}20 0%, ${item.color}10 100%)`,
                borderLeft: `4px solid ${item.color}`,
                '&:hover': {
                  background: `linear-gradient(135deg, ${item.color}30 0%, ${item.color}20 100%)`,
                },
              },
              '&:hover': {
                background: `linear-gradient(135deg, ${item.color}10 0%, ${item.color}05 100%)`,
                borderRadius: 2,
              },
            }}
          >
            <ListItemIcon
              sx={{
                color: isSelected ? item.color : 'text.secondary',
                minWidth: 40,
                transition: 'color 0.2s ease',
              }}
            >
              {item.icon}
            </ListItemIcon>
            <ListItemText
              primary={
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <Typography
                    variant="body2"
                    sx={{
                      fontWeight: isSelected ? 600 : 500,
                      color: isSelected ? item.color : 'text.primary',
                      transition: 'all 0.2s ease',
                    }}
                  >
                    {item.label}
                  </Typography>
                  {item.badge && (
                    <Chip
                      label={item.badge}
                      size="small"
                      sx={{
                        height: 20,
                        fontSize: '0.7rem',
                        fontWeight: 600,
                        background: `linear-gradient(135deg, ${leagueColors.gold[500]} 0%, ${leagueColors.gold[600]} 100%)`,
                        color: '#000',
                        '& .MuiChip-label': {
                          px: 1,
                        },
                      }}
                    />
                  )}
                </Box>
              }
            />
            {hasChildren && (
              <Box sx={{ color: 'text.secondary' }}>
                {expanded ? <ExpandLess /> : <ExpandMore />}
              </Box>
            )}
          </ListItemButton>
        </ListItem>

        {hasChildren && (
          <Collapse in={expanded} timeout="auto" unmountOnExit>
            <List component="div" disablePadding>
              {item.children?.map(child => renderNavigationItem(child, depth + 1))}
            </List>
          </Collapse>
        )}
      </React.Fragment>
    );
  };

  return (
    <Box sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      {/* Header */}
      <Box
        sx={{
          p: 3,
          background: isDarkMode
            ? `linear-gradient(135deg, ${leagueColors.blue[900]} 0%, ${leagueColors.dark[100]} 100%)`
            : `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
          color: '#fff',
          textAlign: 'center',
        }}
      >
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 1, mb: 1 }}>
          <SportsEsports sx={{ fontSize: 32, color: leagueColors.gold[400] }} />
          <Typography variant="h5" sx={{ fontWeight: 700, letterSpacing: 1 }}>
            Herald.lol
          </Typography>
        </Box>
        <Typography variant="caption" sx={{ opacity: 0.9, fontSize: '0.75rem' }}>
          Analysez vos performances League of Legends
        </Typography>
      </Box>

      <Divider />

      {/* Navigation Items */}
      <Box sx={{ flex: 1, overflow: 'auto', py: 1 }}>
        <List>
          {items.map(item => renderNavigationItem(item))}
        </List>
      </Box>

      <Divider />

      {/* Footer */}
      <Box sx={{ p: 2, textAlign: 'center' }}>
        <Typography variant="caption" color="text.secondary" sx={{ fontSize: '0.7rem' }}>
          Version 2.0 - Redesign LoL
        </Typography>
      </Box>
    </Box>
  );
};

export default NavigationSidebar;