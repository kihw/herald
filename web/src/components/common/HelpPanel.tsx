import React, { useState } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  IconButton,
  Collapse,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Chip,
} from '@mui/material';
import {
  Help,
  ExpandMore,
  ExpandLess,
  Dashboard,
  History,
  Settings,
  Sync,
  Timeline,
} from '@mui/icons-material';

const HelpPanel: React.FC = () => {
  const [expanded, setExpanded] = useState(false);

  const features = [
    {
      icon: <Dashboard color="primary" />,
      title: 'Dashboard Overview',
      description: 'View your global statistics, win rate, and favorite champion',
    },
    {
      icon: <History color="secondary" />,
      title: 'Match History',
      description: 'Browse your complete match history with detailed information',
    },
    {
      icon: <Sync color="success" />,
      title: 'Sync Matches',
      description: 'Synchronize your latest matches from Riot Games API',
    },
    {
      icon: <Settings color="warning" />,
      title: 'User Settings',
      description: 'Customize your experience and data collection preferences',
    },
    {
      icon: <Timeline color="info" />,
      title: 'Analytics',
      description: 'Advanced match analysis and performance tracking',
    },
  ];

  return (
    <Card sx={{ mb: 3, background: 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)' }}>
      <CardContent>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <Help color="primary" sx={{ mr: 1 }} />
            <Typography variant="h6" component="div">
              Welcome to LoL Match Exporter!
            </Typography>
            <Chip
              label="New"
              color="primary"
              size="small"
              sx={{ ml: 2 }}
            />
          </Box>
          <IconButton
            onClick={() => setExpanded(!expanded)}
            size="small"
          >
            {expanded ? <ExpandLess /> : <ExpandMore />}
          </IconButton>
        </Box>
        
        <Collapse in={expanded} timeout="auto" unmountOnExit>
          <Box sx={{ mt: 2 }}>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              Explore all the features available to analyze your League of Legends performance:
            </Typography>
            
            <List dense>
              {features.map((feature, index) => (
                <ListItem key={index} sx={{ py: 0.5 }}>
                  <ListItemIcon sx={{ minWidth: 40 }}>
                    {feature.icon}
                  </ListItemIcon>
                  <ListItemText
                    primary={feature.title}
                    secondary={feature.description}
                    primaryTypographyProps={{ variant: 'body2', fontWeight: 'medium' }}
                    secondaryTypographyProps={{ variant: 'caption' }}
                  />
                </ListItem>
              ))}
            </List>
            
            <Box sx={{ mt: 2, p: 2, bgcolor: 'rgba(25, 118, 210, 0.1)', borderRadius: 1 }}>
              <Typography variant="caption" color="primary" sx={{ fontWeight: 'medium' }}>
                💡 Pro Tip: Start by clicking on "Sync Matches" in the Match History tab to load your recent games!
              </Typography>
            </Box>
          </Box>
        </Collapse>
      </CardContent>
    </Card>
  );
};

export default HelpPanel;
