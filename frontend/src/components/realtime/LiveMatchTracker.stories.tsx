import type { Meta, StoryObj } from '@storybook/react'
import LiveMatchTracker from './LiveMatchTracker'

const meta: Meta<typeof LiveMatchTracker> = {
  title: 'Gaming Realtime/LiveMatchTracker',
  component: LiveMatchTracker,
  parameters: {
    layout: 'fullscreen',
    backgrounds: { default: 'herald-dark' },
    docs: {
      description: {
        component: 'Real-time match tracker component providing live game statistics, player performance, and match predictions.'
      }
    }
  },
  tags: ['autodocs'],
  argTypes: {
    matchId: {
      control: 'text',
      description: 'Live match ID for tracking',
    },
    region: {
      control: 'select',
      options: ['na1', 'euw1', 'kr', 'br1', 'eun1', 'jp1'],
      description: 'Region for live match tracking',
    },
    updateInterval: {
      control: 'number',
      description: 'Update interval in seconds',
      defaultValue: 30,
    }
  }
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    matchId: 'NA1_live_match_123',
    region: 'na1',
    updateInterval: 30,
  }
}

export const EarlyGame: Story = {
  name: 'Early Game (5 minutes)',
  args: {
    matchId: 'NA1_early_game',
    region: 'na1',
    updateInterval: 15,
  },
  parameters: {
    docs: {
      description: {
        story: 'Live tracker showing early game state with farming and first blood potential.'
      }
    }
  }
}

export const MidGame: Story = {
  name: 'Mid Game (20 minutes)',
  args: {
    matchId: 'NA1_mid_game',
    region: 'na1',
    updateInterval: 30,
  },
  parameters: {
    docs: {
      description: {
        story: 'Mid game tracker showing team fights, objectives, and power spikes.'
      }
    }
  }
}

export const LateGame: Story = {
  name: 'Late Game (40+ minutes)',
  args: {
    matchId: 'NA1_late_game',
    region: 'na1',
    updateInterval: 20,
  },
  parameters: {
    docs: {
      description: {
        story: 'Late game tracker with high stakes team fights and game-ending potential.'
      }
    }
  }
}

export const HighEloMatch: Story = {
  name: 'High Elo Match (Diamond+)',
  args: {
    matchId: 'KR_high_elo_match',
    region: 'kr',
    updateInterval: 25,
  },
  parameters: {
    docs: {
      description: {
        story: 'Live tracking of high elo match with advanced macro play and mechanics.'
      }
    }
  }
}

export const ProMatch: Story = {
  name: 'Professional Match',
  args: {
    matchId: 'TOURNAMENT_pro_match',
    region: 'na1',
    updateInterval: 10,
  },
  parameters: {
    docs: {
      description: {
        story: 'Professional esports match tracker with enhanced analytics and predictions.'
      }
    }
  }
}

export const CloseGame: Story = {
  name: 'Close Match (50-50)',
  args: {
    matchId: 'NA1_close_game',
    region: 'na1',
    updateInterval: 20,
  },
  parameters: {
    docs: {
      description: {
        story: 'Tracker showing a very close match with balanced team compositions and equal chances.'
      }
    }
  }
}