/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
        mono: ['JetBrains Mono', 'Fira Code', 'monospace'],
      },
      fontSize: {
        '2xs': ['0.6875rem', { lineHeight: '1rem' }],
      },
      colors: {
        surface: {
          DEFAULT: '#0c0c0e',
          card: '#151518',
          hover: '#1c1c20',
          border: '#27272a',
          'border-light': '#1f1f23',
        },
        accent: {
          DEFAULT: '#6366f1',
          hover: '#818cf8',
          muted: '#4f46e5',
          subtle: '#6366f10d',
        },
      },
      spacing: {
        '18': '4.5rem',
        '88': '22rem',
      },
      ringColor: {
        accent: '#6366f1',
      },
    },
  },
  plugins: [require('@tailwindcss/forms')],
}
