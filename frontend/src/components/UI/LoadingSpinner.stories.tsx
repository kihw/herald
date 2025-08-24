import type { Meta, StoryObj } from '@storybook/react'
import LoadingSpinner from './LoadingSpinner'

const meta: Meta<typeof LoadingSpinner> = {
  title: 'UI Components/LoadingSpinner',
  component: LoadingSpinner,
  parameters: {
    layout: 'centered',
    backgrounds: { default: 'herald-dark' },
    docs: {
      description: {
        component: 'Loading spinner component with Herald.lol gaming theme styling and various size options.'
      }
    }
  },
  tags: ['autodocs'],
  argTypes: {
    size: {
      control: 'select',
      options: ['small', 'medium', 'large'],
      description: 'Size of the loading spinner',
    },
    color: {
      control: 'color',
      description: 'Custom color for the spinner',
    },
    message: {
      control: 'text',
      description: 'Optional loading message',
    }
  }
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    size: 'medium',
  }
}

export const Small: Story = {
  args: {
    size: 'small',
  }
}

export const Large: Story = {
  args: {
    size: 'large',
  }
}

export const WithMessage: Story = {
  args: {
    size: 'medium',
    message: 'Loading match data...',
  }
}

export const GamingAnalytics: Story = {
  name: 'Gaming Analytics Loading',
  args: {
    size: 'large',
    message: 'Analyzing your gaming performance...',
  },
  parameters: {
    docs: {
      description: {
        story: 'Loading state for gaming analytics with descriptive message.'
      }
    }
  }
}

export const RiotAPI: Story = {
  name: 'Riot API Loading',
  args: {
    size: 'medium',
    message: 'Syncing with Riot Games API...',
  },
  parameters: {
    docs: {
      description: {
        story: 'Loading state when fetching data from Riot Games API.'
      }
    }
  }
}

export const LiveMatch: Story = {
  name: 'Live Match Loading',
  args: {
    size: 'large',
    message: 'Connecting to live match data...',
  },
  parameters: {
    docs: {
      description: {
        story: 'Loading state for real-time match tracking connection.'
      }
    }
  }
}

export const AllSizes: Story = {
  name: 'Size Comparison',
  render: () => (
    <div style={{ 
      display: 'flex', 
      alignItems: 'center', 
      gap: '32px',
      flexDirection: 'column' 
    }}>
      <div style={{ textAlign: 'center' }}>
        <h3>Small</h3>
        <LoadingSpinner size="small" />
      </div>
      <div style={{ textAlign: 'center' }}>
        <h3>Medium</h3>
        <LoadingSpinner size="medium" />
      </div>
      <div style={{ textAlign: 'center' }}>
        <h3>Large</h3>
        <LoadingSpinner size="large" />
      </div>
    </div>
  ),
  parameters: {
    docs: {
      description: {
        story: 'Comparison of all available spinner sizes.'
      }
    }
  }
}