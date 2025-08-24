import type { Meta, StoryObj } from '@storybook/react'
import ChampionAnalytics from './ChampionAnalytics'

const meta: Meta<typeof ChampionAnalytics> = {
  title: 'Gaming Analytics/ChampionAnalytics',
  component: ChampionAnalytics,
  parameters: {
    layout: 'padded',
    backgrounds: { default: 'herald-dark' },
    docs: {
      description: {
        component: 'Champion analytics component displaying detailed champion performance metrics for League of Legends.'
      }
    }
  },
  tags: ['autodocs'],
  argTypes: {
    playerId: {
      control: 'text',
      description: 'Player ID for analytics',
    },
    champion: {
      control: 'text', 
      description: 'Champion name',
    },
    timeRange: {
      control: 'select',
      options: ['7d', '30d', '90d', 'season'],
      description: 'Time range for analysis',
    },
    position: {
      control: 'select',
      options: ['', 'top', 'jungle', 'mid', 'adc', 'support'],
      description: 'Lane position filter',
    }
  }
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    playerId: 'sample-player-id',
    champion: 'Yasuo',
    timeRange: '30d',
    position: 'mid',
  }
}

export const WithLoading: Story = {
  args: {
    playerId: 'loading-player-id',
    champion: 'Yasuo',
    timeRange: '30d',
  },
  parameters: {
    msw: {
      handlers: [
        // Mock loading state
      ]
    }
  }
}

export const WithError: Story = {
  args: {
    playerId: 'error-player-id',
    champion: 'InvalidChampion',
    timeRange: '30d',
  },
  parameters: {
    docs: {
      description: {
        story: 'Error state when champion data cannot be loaded.'
      }
    }
  }
}

export const PopularChampions: Story = {
  name: 'Popular Champions',
  render: () => (
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(400px, 1fr))', gap: '16px' }}>
      <ChampionAnalytics playerId="yasuo-main" champion="Yasuo" timeRange="30d" position="mid" />
      <ChampionAnalytics playerId="lee-main" champion="LeeSin" timeRange="30d" position="jungle" />
      <ChampionAnalytics playerId="zed-main" champion="Zed" timeRange="30d" position="mid" />
      <ChampionAnalytics playerId="ahri-main" champion="Ahri" timeRange="30d" position="mid" />
    </div>
  ),
  parameters: {
    docs: {
      description: {
        story: 'Multiple champion analytics cards showing popular champions (Yasuo, Lee Sin, Zed, Ahri).'
      }
    }
  }
}