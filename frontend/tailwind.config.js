/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			fontFamily: {
				sans: ['Inter', 'system-ui', 'sans-serif'],
				kannada: ['Anek Kannada', 'system-ui', 'sans-serif']
			},
			colors: {
				primary: {
					50: '#eff6ff',
					100: '#dbeafe',
					200: '#bfdbfe',
					300: '#93c5fd',
					400: '#60a5fa',
					500: '#3b82f6',
					600: '#2563eb',
					700: '#1d4ed8',
					800: '#1e40af',
					900: '#1e3a8a'
				},
				success: {
					500: '#22c55e',
					600: '#16a34a'
				},
				warning: {
					500: '#f59e0b',
					600: '#d97706'
				},
				danger: {
					50: '#fef2f2',
					500: '#ef4444',
					600: '#dc2626'
				}
			},
			borderRadius: {
				card: '0.75rem',
				input: '0.5rem',
				button: '0.5rem'
			}
		}
	},
	plugins: []
};
