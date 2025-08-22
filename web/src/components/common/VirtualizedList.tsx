import React, { useMemo, useCallback, useRef, useEffect } from 'react';
import {
  Box,
  List,
  ListItem,
  ListItemText,
  Typography,
  useTheme,
} from '@mui/material';
import { useVirtualScroll } from '../../hooks/usePerformance';
import useResponsive from '../../hooks/useResponsive';

interface VirtualizedListProps<T> {
  items: T[];
  itemHeight: number;
  containerHeight: number;
  renderItem: (item: T, index: number) => React.ReactNode;
  onItemClick?: (item: T, index: number) => void;
  loading?: boolean;
  emptyMessage?: string;
  overscan?: number;
}

const VirtualizedList = <T extends any>({
  items,
  itemHeight,
  containerHeight,
  renderItem,
  onItemClick,
  loading = false,
  emptyMessage = 'Aucun élément à afficher',
  overscan = 5,
}: VirtualizedListProps<T>) => {
  const theme = useTheme();
  const { isMobile } = useResponsive();
  const containerRef = useRef<HTMLDivElement>(null);
  const { visibleItems, handleScroll } = useVirtualScroll(items, itemHeight, containerHeight);

  // Adjust overscan based on device performance
  const adjustedOverscan = isMobile ? Math.max(2, overscan / 2) : overscan;

  // Calculate visible range with overscan
  const visibleRange = useMemo(() => {
    const startIndex = Math.max(0, visibleItems.startIndex - adjustedOverscan);
    const endIndex = Math.min(items.length, visibleItems.endIndex + adjustedOverscan);
    
    return {
      startIndex,
      endIndex,
      items: items.slice(startIndex, endIndex),
      offsetY: startIndex * itemHeight,
    };
  }, [visibleItems, items, itemHeight, adjustedOverscan]);

  // Keyboard navigation
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!containerRef.current) return;
      
      const container = containerRef.current;
      const currentScroll = container.scrollTop;
      const itemsPerScreen = Math.floor(containerHeight / itemHeight);
      
      switch (e.key) {
        case 'ArrowDown':
          e.preventDefault();
          container.scrollTop = Math.min(
            currentScroll + itemHeight,
            visibleItems.totalHeight - containerHeight
          );
          break;
        case 'ArrowUp':
          e.preventDefault();
          container.scrollTop = Math.max(currentScroll - itemHeight, 0);
          break;
        case 'PageDown':
          e.preventDefault();
          container.scrollTop = Math.min(
            currentScroll + itemHeight * itemsPerScreen,
            visibleItems.totalHeight - containerHeight
          );
          break;
        case 'PageUp':
          e.preventDefault();
          container.scrollTop = Math.max(
            currentScroll - itemHeight * itemsPerScreen,
            0
          );
          break;
        case 'Home':
          e.preventDefault();
          container.scrollTop = 0;
          break;
        case 'End':
          e.preventDefault();
          container.scrollTop = visibleItems.totalHeight - containerHeight;
          break;
      }
    };

    const container = containerRef.current;
    if (container) {
      container.addEventListener('keydown', handleKeyDown);
      return () => container.removeEventListener('keydown', handleKeyDown);
    }
  }, [containerHeight, itemHeight, visibleItems.totalHeight]);

  if (loading) {
    return (
      <Box
        sx={{
          height: containerHeight,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <Typography variant="body2" color="text.secondary">
          Chargement...
        </Typography>
      </Box>
    );
  }

  if (items.length === 0) {
    return (
      <Box
        sx={{
          height: containerHeight,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <Typography variant="body2" color="text.secondary">
          {emptyMessage}
        </Typography>
      </Box>
    );
  }

  return (
    <Box
      ref={containerRef}
      sx={{
        height: containerHeight,
        overflow: 'auto',
        position: 'relative',
        '&:focus': {
          outline: `2px solid ${theme.palette.primary.main}`,
          outlineOffset: -2,
        },
      }}
      onScroll={handleScroll}
      tabIndex={0}
      role="listbox"
      aria-label="Liste virtualisée"
    >
      {/* Total height spacer */}
      <Box sx={{ height: visibleItems.totalHeight, position: 'relative' }}>
        {/* Visible items container */}
        <Box
          sx={{
            position: 'absolute',
            top: visibleRange.offsetY,
            left: 0,
            right: 0,
          }}
        >
          <List disablePadding>
            {visibleRange.items.map((item, virtualIndex) => {
              const actualIndex = visibleRange.startIndex + virtualIndex;
              
              return (
                <ListItem
                  key={actualIndex}
                  button={!!onItemClick}
                  onClick={onItemClick ? () => onItemClick(item, actualIndex) : undefined}
                  sx={{
                    height: itemHeight,
                    minHeight: itemHeight,
                    maxHeight: itemHeight,
                    cursor: onItemClick ? 'pointer' : 'default',
                    '&:hover': onItemClick ? {
                      backgroundColor: theme.palette.action.hover,
                    } : {},
                    '&:focus': {
                      backgroundColor: theme.palette.action.focus,
                    },
                  }}
                  role="option"
                  aria-selected={false}
                  tabIndex={-1}
                >
                  {renderItem(item, actualIndex)}
                </ListItem>
              );
            })}
          </List>
        </Box>
      </Box>

      {/* Scroll indicator */}
      {items.length > Math.floor(containerHeight / itemHeight) && (
        <Box
          sx={{
            position: 'absolute',
            right: 4,
            top: 4,
            bottom: 4,
            width: 4,
            backgroundColor: 'rgba(0, 0, 0, 0.1)',
            borderRadius: 2,
          }}
        >
          <Box
            sx={{
              position: 'absolute',
              right: 0,
              width: 4,
              backgroundColor: theme.palette.primary.main,
              borderRadius: 2,
              top: `${(visibleItems.startIndex / items.length) * 100}%`,
              height: `${(Math.min(visibleItems.endIndex - visibleItems.startIndex, items.length) / items.length) * 100}%`,
            }}
          />
        </Box>
      )}
    </Box>
  );
};

// Specialized components for common use cases
interface VirtualizedGroupListProps {
  groups: Array<{ id: number; name: string; description: string; member_count: number }>;
  onGroupClick: (group: any) => void;
  containerHeight?: number;
}

export const VirtualizedGroupList: React.FC<VirtualizedGroupListProps> = ({
  groups,
  onGroupClick,
  containerHeight = 400,
}) => {
  const renderGroupItem = useCallback((group: any, index: number) => (
    <ListItemText
      primary={group.name}
      secondary={
        <React.Fragment>
          <Typography variant="body2" component="span">
            {group.description || 'Aucune description'}
          </Typography>
          <Typography variant="caption" component="span" sx={{ ml: 1 }}>
            • {group.member_count} membre{group.member_count > 1 ? 's' : ''}
          </Typography>
        </React.Fragment>
      }
    />
  ), []);

  return (
    <VirtualizedList
      items={groups}
      itemHeight={72}
      containerHeight={containerHeight}
      renderItem={renderGroupItem}
      onItemClick={onGroupClick}
      emptyMessage="Aucun groupe trouvé"
    />
  );
};

interface VirtualizedMemberListProps {
  members: Array<{ id: number; user: { riot_id: string; riot_tag: string }; role: string }>;
  onMemberClick?: (member: any) => void;
  containerHeight?: number;
}

export const VirtualizedMemberList: React.FC<VirtualizedMemberListProps> = ({
  members,
  onMemberClick,
  containerHeight = 300,
}) => {
  const renderMemberItem = useCallback((member: any, index: number) => (
    <ListItemText
      primary={`${member.user.riot_id}#${member.user.riot_tag}`}
      secondary={member.role}
    />
  ), []);

  return (
    <VirtualizedList
      items={members}
      itemHeight={56}
      containerHeight={containerHeight}
      renderItem={renderMemberItem}
      onItemClick={onMemberClick}
      emptyMessage="Aucun membre trouvé"
    />
  );
};

export default VirtualizedList;