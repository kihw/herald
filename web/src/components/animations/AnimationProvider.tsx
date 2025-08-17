import React, { createContext, useContext, useState, ReactNode } from 'react';
import {
  Fade,
  Slide,
  Zoom,
  Grow,
  Collapse,
  Box,
  useTheme,
} from '@mui/material';
import { TransitionProps } from '@mui/material/transitions';

export type AnimationType = 'fade' | 'slide' | 'zoom' | 'grow' | 'collapse' | 'slideUp' | 'slideDown' | 'slideLeft' | 'slideRight';
export type AnimationSpeed = 'slow' | 'normal' | 'fast';

interface AnimationConfig {
  type: AnimationType;
  speed: AnimationSpeed;
  duration: number;
  delay: number;
  easing: string;
  enabled: boolean;
}

interface AnimationContextType {
  config: AnimationConfig;
  updateConfig: (updates: Partial<AnimationConfig>) => void;
  animate: (children: ReactNode, type?: AnimationType, options?: Partial<TransitionProps>) => ReactNode;
  getStaggerDelay: (index: number, baseDelay?: number) => number;
}

const defaultConfig: AnimationConfig = {
  type: 'fade',
  speed: 'normal',
  duration: 300,
  delay: 0,
  easing: 'ease-in-out',
  enabled: true,
};

const speedDurations = {
  slow: 600,
  normal: 300,
  fast: 150,
};

const AnimationContext = createContext<AnimationContextType | undefined>(undefined);

export const useAnimation = () => {
  const context = useContext(AnimationContext);
  if (!context) {
    throw new Error('useAnimation must be used within AnimationProvider');
  }
  return context;
};

interface AnimationProviderProps {
  children: ReactNode;
  initialConfig?: Partial<AnimationConfig>;
}

export const AnimationProvider: React.FC<AnimationProviderProps> = ({
  children,
  initialConfig = {},
}) => {
  const theme = useTheme();
  const [config, setConfig] = useState<AnimationConfig>({
    ...defaultConfig,
    ...initialConfig,
  });

  const updateConfig = (updates: Partial<AnimationConfig>) => {
    setConfig(prev => ({
      ...prev,
      ...updates,
      duration: updates.speed ? speedDurations[updates.speed] : prev.duration,
    }));
  };

  const getStaggerDelay = (index: number, baseDelay: number = 50): number => {
    return config.enabled ? baseDelay * index : 0;
  };

  const animate = (
    children: ReactNode, 
    type: AnimationType = config.type, 
    options: Partial<TransitionProps> = {}
  ): ReactNode => {
    if (!config.enabled) {
      return children;
    }

    const duration = options.timeout || config.duration;
    const commonProps = {
      in: true,
      timeout: duration,
      style: {
        transitionTimingFunction: config.easing,
      },
      ...options,
    };

    switch (type) {
      case 'fade':
        return (
          <Fade {...commonProps}>
            <Box component="div">{children}</Box>
          </Fade>
        );

      case 'slide':
      case 'slideLeft':
        return (
          <Slide {...commonProps} direction="left">
            <Box component="div">{children}</Box>
          </Slide>
        );

      case 'slideRight':
        return (
          <Slide {...commonProps} direction="right">
            <Box component="div">{children}</Box>
          </Slide>
        );

      case 'slideUp':
        return (
          <Slide {...commonProps} direction="up">
            <Box component="div">{children}</Box>
          </Slide>
        );

      case 'slideDown':
        return (
          <Slide {...commonProps} direction="down">
            <Box component="div">{children}</Box>
          </Slide>
        );

      case 'zoom':
        return (
          <Zoom {...commonProps}>
            <Box component="div">{children}</Box>
          </Zoom>
        );

      case 'grow':
        return (
          <Grow {...commonProps}>
            <Box component="div">{children}</Box>
          </Grow>
        );

      case 'collapse':
        return (
          <Collapse {...commonProps}>
            <Box component="div">{children}</Box>
          </Collapse>
        );

      default:
        return children;
    }
  };

  const contextValue: AnimationContextType = {
    config,
    updateConfig,
    animate,
    getStaggerDelay,
  };

  return (
    <AnimationContext.Provider value={contextValue}>
      {children}
    </AnimationContext.Provider>
  );
};