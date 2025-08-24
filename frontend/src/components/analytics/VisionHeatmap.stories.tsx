import type { Meta, StoryObj } from '@storybook/react'
import VisionHeatmap from './VisionHeatmap'

const meta: Meta<typeof VisionHeatmap> = {
  title: 'Gaming Analytics/VisionHeatmap',
  component: VisionHeatmap,
  parameters: {
    layout: 'padded',
    backgrounds: { default: 'herald-dark' },
    docs: {
      description: {
        component: 'Vision heatmap component displaying ward placement patterns and vision control analytics for League of Legends matches.'
      }
    }
  },
  tags: ['autodocs'],
  argTypes: {
    matchId: {
      control: 'text',
      description: 'Match ID for vision data',
    },
    playerId: {
      control: 'text',
      description: 'Player ID for ward analytics',
    },
    mapSide: {
      control: 'select',
      options: ['blue', 'red', 'both'],
      description: 'Map side for vision analysis',
    },
    wardType: {
      control: 'select',
      options: ['all', 'yellow', 'control', 'pink'],
      description: 'Type of wards to display',
    }
  }
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    matchId: 'NA1_vision_game',
    playerId: 'vision-player',
    mapSide: 'both',
    wardType: 'all',
  }
}

export const SupportVision: Story = {
  name: 'Support Ward Placement',
  args: {
    matchId: 'NA1_support_vision',
    playerId: 'support-main-id',
    mapSide: 'both',
    wardType: 'all',
  },
  parameters: {
    docs: {
      description: {
        story: 'Vision heatmap showing optimal ward placement patterns from a support player perspective.'
      }
    }
  }
}

export const JungleVision: Story = {
  name: 'Jungle Control Vision',
  args: {
    matchId: 'NA1_jungle_vision',
    playerId: 'jungle-main-id',
    mapSide: 'both',
    wardType: 'control',
  },
  parameters: {
    docs: {
      description: {
        story: 'Heatmap focusing on jungle vision control and objective warding patterns.'
      }
    }
  }
}

export const BlueTeamVision: Story = {
  name: 'Blue Team Perspective',
  args: {
    matchId: 'NA1_blue_team',
    playerId: 'blue-team-player',
    mapSide: 'blue',
    wardType: 'all',
  },
  parameters: {
    docs: {
      description: {
        story: 'Vision analysis from blue team perspective showing optimal ward spots.'
      }
    }
  }
}

export const RedTeamVision: Story = {
  name: 'Red Team Perspective', 
  args: {
    matchId: 'NA1_red_team',
    playerId: 'red-team-player',
    mapSide: 'red',
    wardType: 'all',
  },
  parameters: {
    docs: {
      description: {
        story: 'Vision analysis from red team perspective with different optimal ward locations.'
      }
    }
  }
}

export const ControlWardsOnly: Story = {
  name: 'Control Ward Analysis',
  args: {
    matchId: 'NA1_control_wards',
    playerId: 'control-ward-expert',
    mapSide: 'both',
    wardType: 'pink',
  },
  parameters: {
    docs: {
      description: {
        story: 'Specialized view focusing only on control ward (pink ward) placement patterns.'
      }
    }
  }
}