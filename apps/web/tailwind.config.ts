import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: 'var(--color-primary-50)',
          100: 'var(--color-primary-100)',
          500: 'var(--color-primary-500)',
          600: 'var(--color-primary-600)',
          700: 'var(--color-primary-700)',
          900: 'var(--color-primary-900)',
        },
        success: 'var(--color-success-500)',
        warning: 'var(--color-warning-500)',
        danger: 'var(--color-danger-500)',
        info: 'var(--color-info-500)',
        'success-bg': 'var(--color-success-bg)',
        'warning-bg': 'var(--color-warning-bg)',
        'danger-bg': 'var(--color-danger-bg)',
        'info-bg': 'var(--color-info-bg)',
        text: {
          primary: 'var(--color-text-primary)',
          secondary: 'var(--color-text-secondary)',
          muted: 'var(--color-text-muted)',
          disabled: 'var(--color-text-disabled)',
        },
        border: {
          default: 'var(--color-border-default)',
          strong: 'var(--color-border-strong)',
        },
        surface: {
          DEFAULT: 'var(--color-surface)',
          subtle: 'var(--color-surface-subtle)',
          page: 'var(--color-surface-page)',
          hover: 'var(--color-surface-hover)',
        },
        'focus-ring': 'var(--color-focus-ring)',
      },
      fontFamily: {
        display: 'var(--font-display-lg)',
        title: ['var(--font-title-lg)', 'sans-serif'],
        body: ['var(--font-body-md)', 'sans-serif'],
        metric: ['var(--font-metric-lg)', 'monospace'],
        mono: ['var(--font-metric-md)', 'monospace'],
      },
      borderRadius: {
        sm: 'var(--radius-sm)',
        md: 'var(--radius-md)',
        lg: 'var(--radius-lg)',
        full: 'var(--radius-full)',
      },
      boxShadow: {
        xs: 'var(--shadow-xs)',
        sm: 'var(--shadow-sm)',
        md: 'var(--shadow-md)',
        lg: 'var(--shadow-lg)',
      },
      spacing: {
        '1': 'var(--space-1)',
        '2': 'var(--space-2)',
        '3': 'var(--space-3)',
        '4': 'var(--space-4)',
        '6': 'var(--space-6)',
        '8': 'var(--space-8)',
        '10': 'var(--space-10)',
        '12': 'var(--space-12)',
      },
    },
  },
  plugins: [],
}

export default config
