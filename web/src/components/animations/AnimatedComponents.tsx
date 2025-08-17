import React, { useState, useEffect, ReactNode } from 'react';
import {
  Card,
  CardProps,
  Button,
  ButtonProps,
  Box,
  BoxProps,
  Chip,
  ChipProps,
  Typography,
  TypographyProps,
  IconButton,
  IconButtonProps,
  Fade,
  Slide,
  Zoom,
  Grow,
  keyframes,
  styled,
  alpha,
  useTheme,
} from '@mui/material';
import { useAnimation } from './AnimationProvider';

// Keyframes pour les animations personnalisées
const pulse = keyframes`
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
  100% {
    transform: scale(1);
  }
`;

const shake = keyframes`
  0%, 100% {
    transform: translateX(0);
  }
  10%, 30%, 50%, 70%, 90% {
    transform: translateX(-2px);
  }
  20%, 40%, 60%, 80% {
    transform: translateX(2px);
  }
`;

const bounce = keyframes`
  0%, 20%, 53%, 80%, 100% {
    transform: translateY(0);
  }
  40%, 43% {
    transform: translateY(-8px);
  }
  70% {
    transform: translateY(-4px);
  }
  90% {
    transform: translateY(-2px);
  }
`;

const slideInFromRight = keyframes`
  0% {
    transform: translateX(100%);
    opacity: 0;
  }
  100% {
    transform: translateX(0);
    opacity: 1;
  }
`;

const slideInFromLeft = keyframes`
  0% {
    transform: translateX(-100%);
    opacity: 0;
  }
  100% {
    transform: translateX(0);
    opacity: 1;
  }
`;

// Styled components avec animations
const AnimatedCardStyled = styled(Card, {
  shouldForwardProp: (prop) => prop !== 'animationDelay' && prop !== 'hoverEffect',
})<{ animationDelay?: number; hoverEffect?: boolean }>(({ theme, animationDelay = 0, hoverEffect = true }) => ({
  animation: `${slideInFromLeft} 0.6s ease-out ${animationDelay}ms both`,
  transition: 'all 0.3s ease-in-out',
  
  ...(hoverEffect && {
    '&:hover': {
      transform: 'translateY(-4px)',
      boxShadow: theme.shadows[8],
      '& .MuiCardContent-root': {
        transform: 'scale(1.02)',
      },
    },
  }),
  
  '& .MuiCardContent-root': {
    transition: 'transform 0.3s ease-in-out',
  },
}));

const AnimatedButtonStyled = styled(Button)<{ rippleEffect?: boolean }>(({ theme, rippleEffect = true }) => ({
  transition: 'all 0.2s ease-in-out',
  position: 'relative',
  overflow: 'hidden',
  
  '&:hover': {
    transform: 'translateY(-2px)',
    boxShadow: theme.shadows[4],
  },
  
  '&:active': {
    transform: 'translateY(0)',
  },
  
  ...(rippleEffect && {
    '&::before': {
      content: '""',
      position: 'absolute',
      top: '50%',
      left: '50%',
      width: 0,
      height: 0,
      borderRadius: '50%',
      background: alpha(theme.palette.common.white, 0.3),
      transform: 'translate(-50%, -50%)',
      transition: 'width 0.6s, height 0.6s',
    },
    
    '&:active::before': {
      width: '300px',
      height: '300px',
    },
  }),
}));

const PulsingBox = styled(Box)<{ pulseSpeed?: 'slow' | 'normal' | 'fast' }>(({ pulseSpeed = 'normal' }) => {
  const duration = pulseSpeed === 'slow' ? '2s' : pulseSpeed === 'fast' ? '0.5s' : '1s';
  
  return {
    animation: `${pulse} ${duration} ease-in-out infinite`,
  };
});

const ShakingBox = styled(Box)({
  '&.shake': {
    animation: `${shake} 0.5s ease-in-out`,
  },
});

const BouncingBox = styled(Box)({
  '&.bounce': {
    animation: `${bounce} 1s ease-in-out`,
  },
});

// Composants animés exportables
interface AnimatedCardProps extends CardProps {
  animationDelay?: number;
  hoverEffect?: boolean;
  children: ReactNode;
}

export const AnimatedCard: React.FC<AnimatedCardProps> = ({
  animationDelay = 0,
  hoverEffect = true,
  children,
  ...props
}) => {
  const { config } = useAnimation();
  
  if (!config.enabled) {
    return <Card {...props}>{children}</Card>;
  }
  
  return (
    <AnimatedCardStyled
      animationDelay={animationDelay}
      hoverEffect={hoverEffect}
      {...props}
    >
      {children}
    </AnimatedCardStyled>
  );
};

interface AnimatedButtonProps extends ButtonProps {
  rippleEffect?: boolean;
  children: ReactNode;
}

export const AnimatedButton: React.FC<AnimatedButtonProps> = ({
  rippleEffect = true,
  children,
  ...props
}) => {
  const { config } = useAnimation();
  
  if (!config.enabled) {
    return <Button {...props}>{children}</Button>;
  }
  
  return (
    <AnimatedButtonStyled rippleEffect={rippleEffect} {...props}>
      {children}
    </AnimatedButtonStyled>
  );
};

interface StaggeredListProps {
  children: ReactNode[];
  staggerDelay?: number;
  animationType?: 'fade' | 'slide' | 'zoom' | 'grow';
  direction?: 'up' | 'down' | 'left' | 'right';
}

export const StaggeredList: React.FC<StaggeredListProps> = ({
  children,
  staggerDelay = 100,
  animationType = 'fade',
  direction = 'up',
}) => {
  const { animate, getStaggerDelay, config } = useAnimation();
  const [visibleItems, setVisibleItems] = useState<boolean[]>([]);

  useEffect(() => {
    if (!config.enabled) {
      setVisibleItems(children.map(() => true));
      return;
    }

    const timers: NodeJS.Timeout[] = [];
    
    children.forEach((_, index) => {
      const timer = setTimeout(() => {
        setVisibleItems(prev => {
          const newVisible = [...prev];
          newVisible[index] = true;
          return newVisible;
        });
      }, getStaggerDelay(index, staggerDelay));
      
      timers.push(timer);
    });

    return () => {
      timers.forEach(timer => clearTimeout(timer));
    };
  }, [children, staggerDelay, getStaggerDelay, config.enabled]);

  if (!config.enabled) {
    return <>{children}</>;
  }

  return (
    <>
      {children.map((child, index) => {
        const animationTypeWithDirection = direction !== 'up' ? 
          `slide${direction.charAt(0).toUpperCase() + direction.slice(1)}` as any : 
          animationType;
          
        return (
          <Box key={index} sx={{ mb: 1 }}>
            {animate(
              child,
              animationTypeWithDirection,
              {
                in: visibleItems[index] || false,
                timeout: config.duration,
              }
            )}
          </Box>
        );
      })}
    </>
  );
};

interface CounterAnimationProps {
  value: number;
  duration?: number;
  suffix?: string;
  prefix?: string;
  decimalPlaces?: number;
  onComplete?: () => void;
}

export const CounterAnimation: React.FC<CounterAnimationProps> = ({
  value,
  duration = 1000,
  suffix = '',
  prefix = '',
  decimalPlaces = 0,
  onComplete,
}) => {
  const [currentValue, setCurrentValue] = useState(0);
  const { config } = useAnimation();

  useEffect(() => {
    if (!config.enabled) {
      setCurrentValue(value);
      onComplete?.();
      return;
    }

    let startTime: number;
    let animationFrame: number;

    const animate = (currentTime: number) => {
      if (!startTime) startTime = currentTime;
      const progress = Math.min((currentTime - startTime) / duration, 1);
      
      // Utilisation d'une fonction d'easing
      const easedProgress = 1 - Math.pow(1 - progress, 3); // easeOutCubic
      
      setCurrentValue(value * easedProgress);

      if (progress < 1) {
        animationFrame = requestAnimationFrame(animate);
      } else {
        onComplete?.();
      }
    };

    animationFrame = requestAnimationFrame(animate);

    return () => {
      if (animationFrame) {
        cancelAnimationFrame(animationFrame);
      }
    };
  }, [value, duration, config.enabled, onComplete]);

  return (
    <span>
      {prefix}{currentValue.toFixed(decimalPlaces)}{suffix}
    </span>
  );
};

interface PulsingChipProps extends ChipProps {
  pulseSpeed?: 'slow' | 'normal' | 'fast';
  children?: ReactNode;
}

export const PulsingChip: React.FC<PulsingChipProps> = ({
  pulseSpeed = 'normal',
  children,
  ...props
}) => {
  const { config } = useAnimation();
  
  if (!config.enabled) {
    return <Chip {...props}>{children}</Chip>;
  }
  
  return (
    <PulsingBox pulseSpeed={pulseSpeed}>
      <Chip {...props}>{children}</Chip>
    </PulsingBox>
  );
};

interface ShakeOnErrorProps {
  error: boolean;
  children: ReactNode;
  onAnimationEnd?: () => void;
}

export const ShakeOnError: React.FC<ShakeOnErrorProps> = ({
  error,
  children,
  onAnimationEnd,
}) => {
  const [shaking, setShaking] = useState(false);
  const { config } = useAnimation();

  useEffect(() => {
    if (error && config.enabled) {
      setShaking(true);
      const timer = setTimeout(() => {
        setShaking(false);
        onAnimationEnd?.();
      }, 500);
      
      return () => clearTimeout(timer);
    }
  }, [error, config.enabled, onAnimationEnd]);

  return (
    <ShakingBox className={shaking ? 'shake' : ''}>
      {children}
    </ShakingBox>
  );
};

interface BounceOnSuccessProps {
  success: boolean;
  children: ReactNode;
  onAnimationEnd?: () => void;
}

export const BounceOnSuccess: React.FC<BounceOnSuccessProps> = ({
  success,
  children,
  onAnimationEnd,
}) => {
  const [bouncing, setBouncing] = useState(false);
  const { config } = useAnimation();

  useEffect(() => {
    if (success && config.enabled) {
      setBouncing(true);
      const timer = setTimeout(() => {
        setBouncing(false);
        onAnimationEnd?.();
      }, 1000);
      
      return () => clearTimeout(timer);
    }
  }, [success, config.enabled, onAnimationEnd]);

  return (
    <BouncingBox className={bouncing ? 'bounce' : ''}>
      {children}
    </BouncingBox>
  );
};

interface PageTransitionProps {
  children: ReactNode;
  transitionKey: string;
  direction?: 'left' | 'right' | 'up' | 'down';
}

export const PageTransition: React.FC<PageTransitionProps> = ({
  children,
  transitionKey,
  direction = 'left',
}) => {
  const { animate, config } = useAnimation();

  if (!config.enabled) {
    return <>{children}</>;
  }

  return (
    <Box key={transitionKey}>
      {animate(
        children,
        `slide${direction.charAt(0).toUpperCase() + direction.slice(1)}` as any,
        {
          timeout: { enter: config.duration, exit: config.duration / 2 },
        }
      )}
    </Box>
  );
};

// Hook pour les animations personnalisées
export const useCustomAnimation = () => {
  const theme = useTheme();
  
  const createKeyframe = (name: string, styles: Record<string, any>) => {
    return keyframes(styles);
  };
  
  const createStyledComponent = (component: any, animations: any) => {
    return styled(component)(animations);
  };
  
  return {
    createKeyframe,
    createStyledComponent,
    theme,
    predefinedAnimations: {
      pulse,
      shake,
      bounce,
      slideInFromRight,
      slideInFromLeft,
    },
  };
};