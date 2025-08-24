import type { Meta, StoryObj } from '@storybook/react'
import DamageAnalytics from './DamageAnalytics'

const meta: Meta<typeof DamageAnalytics> = {
  title: 'Gaming Analytics/DamageAnalytics',
  component: DamageAnalytics,
  parameters: {
    layout: 'padded',
    backgrounds: { default: 'herald-dark' },
    docs: {
      description: {
        component: 'Damage analytics component showing comprehensive damage statistics and team contribution metrics.'
      }
    }
  },
  tags: ['autodocs'],
  argTypes: {
    matchId: {
      control: 'text',
      description: 'Match ID from Riot API',
    },
    playerId: {
      control: 'text',
      description: 'Player ID for damage analysis',
    },
    timeRange: {
      control: 'select',
      options: ['last7days', 'last30days', 'lastSeason', 'allTime'],
      description: 'Time range for damage statistics',
    }
  }
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    matchId: 'NA1_4567890123',
    playerId: 'damage-player-id',
    timeRange: 'last30days',
  }
}

export const HighDamageGame: Story = {
  name: 'High Damage Performance',
  args: {
    matchId: 'NA1_high-damage',
    playerId: 'carry-player-id',
    timeRange: 'last7days',
  },
  parameters: {
    docs: {
      description: {
        story: 'Analytics showing exceptional damage performance with high team contribution.'
      }
    }
  }
}

export const SupportDamage: Story = {
  name: 'Support Role Damage',
  args: {
    matchId: 'NA1_support-game',
    playerId: 'support-player-id',
    timeRange: 'last30days',
  },
  parameters: {
    docs: {
      description: {
        story: 'Damage analytics for support role showing utility and peel damage.'
      }
    }
  }
}

export const LowDamageAnalysis: Story = {
  name: 'Low Damage Game Analysis',
  args: {
    matchId: 'NA1_low-damage',
    playerId: 'struggling-player-id',
    timeRange: 'last7days',
  },
  parameters: {
    docs: {
      description: {
        story: 'Analytics highlighting areas for improvement in damage output.'
      }
    }
  }
}

export const ComparisonView: Story = {
  name: 'Team Damage Comparison',
  render: () => (
    <div style={{ display: 'grid', gridTemplateColumns: '1fr', gap: '24px' }}>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '16px' }}>
        <DamageAnalytics matchId="NA1_team-game" playerId="player1" timeRange="last7days" />
        <DamageAnalytics matchId="NA1_team-game" playerId="player2" timeRange="last7days" />
        <DamageAnalytics matchId="NA1_team-game" playerId="player3" timeRange="last7days" />
        <DamageAnalytics matchId="NA1_team-game" playerId="player4" timeRange="last7days" />
        <DamageAnalytics matchId="NA1_team-game" playerId="player5" timeRange="last7days" />
      </div>
    </div>
  ),
  parameters: {
    docs: {
      description: {
        story: 'Comparison view showing damage analytics for all team members in the same match.'
      }
    }
  }
}