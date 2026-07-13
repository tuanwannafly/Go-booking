/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,jsx}'],
  theme: {
    extend: {
      colors: {
        ink: '#000000',
        onPrimary: '#ffffff',
        body: '#5e5e5e',
        mute: '#afafaf',
        hairlineMid: '#4b4b4b',
        canvas: '#ffffff',
        canvasSoft: '#efefef',
        canvasSofter: '#f3f3f3',
        surfacePressed: '#e2e2e2',
        link: '#0000ee',
        blackElevated: '#282828',
      },
      fontFamily: {
        display: ['Inter', 'system-ui', 'Helvetica Neue', 'Arial', 'sans-serif'],
        sans: ['Inter', 'system-ui', 'Helvetica Neue', 'Arial', 'sans-serif'],
      },
      borderRadius: {
        none: '0px',
        md: '8px',
        lg: '12px',
        xl: '16px',
        pill: '999px',
        pillTab: '36px',
        full: '9999px',
      },
      boxShadow: {
        l1: '0 4px 16px 0 rgba(0,0,0,0.12)',
        l2: '0 4px 16px 0 rgba(0,0,0,0.16)',
        l3: '0 2px 8px 0 rgba(0,0,0,0.16)',
      },
      maxWidth: {
        container: '1200px',
      },
      letterSpacing: {
        tighter: '-0.02em',
      },
    },
  },
  plugins: [],
};
