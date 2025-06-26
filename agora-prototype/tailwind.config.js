/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			// Agora color palette
			colors: {
				// Ancient Greek inspired colors
				'agora-marble': '#f8f9fa',
				'agora-stone': '#e9ecef',
				'agora-shadow': '#495057',
				'agora-night': '#1a1a2e',
				'agora-deep': '#16213e',
				'agora-gold': '#d4af37',
				'agora-bronze': '#cd7f32',
				'agora-crimson': '#dc143c',
				'agora-olive': '#808000',

				// Modern accent colors
				'agora-neon-cyan': '#00f5ff',
				'agora-neon-purple': '#8b5cf6',
				'agora-neon-pink': '#ff006e',

				// Semantic color variations
				'agora': {
					50: '#faf9f7',
					100: '#f3f1ed',
					200: '#e6e1d9',
					300: '#d4af37',
					400: '#cd7f32',
					500: '#d4af37',
					600: '#b8941f',
					700: '#9a7a19',
					800: '#7d6315',
					900: '#664f12'
				}
			},

			// Typography
			fontFamily: {
				'sans': ['Inter', 'system-ui', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'sans-serif'],
				'serif': ['Crimson Text', 'Georgia', 'serif'],
				'display': ['Crimson Text', 'Georgia', 'serif']
			},

			// Custom spacing for agora layouts
			spacing: {
				'18': '4.5rem',
				'88': '22rem',
				'128': '32rem'
			},

			// Custom border radius
			borderRadius: {
				'4xl': '2rem'
			},

			// Custom shadows with agora gold glow
			boxShadow: {
				'agora-glow': '0 0 20px rgba(212, 175, 55, 0.3)',
				'agora-glow-lg': '0 0 30px rgba(212, 175, 55, 0.5)',
				'agora-glow-xl': '0 0 40px rgba(212, 175, 55, 0.7)',
				'agora-purple': '0 0 20px rgba(139, 92, 246, 0.3)',
				'agora-cyan': '0 0 20px rgba(0, 245, 255, 0.3)'
			},

			// Custom animations
			animation: {
				'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
				'bounce-slow': 'bounce 2s infinite',
				'glow': 'glow 2s ease-in-out infinite alternate',
				'float': 'float 3s ease-in-out infinite',
				'shimmer': 'shimmer 2s linear infinite'
			},

			keyframes: {
				glow: {
					'0%': { boxShadow: '0 0 20px rgba(212, 175, 55, 0.3)' },
					'100%': { boxShadow: '0 0 30px rgba(212, 175, 55, 0.6)' }
				},
				float: {
					'0%, 100%': { transform: 'translateY(0)' },
					'50%': { transform: 'translateY(-5px)' }
				},
				shimmer: {
					'0%': { backgroundPosition: '-200% 0' },
					'100%': { backgroundPosition: '200% 0' }
				}
			},

			// Custom gradients
			backgroundImage: {
				'agora-gradient': 'linear-gradient(135deg, var(--agora-night) 0%, var(--agora-deep) 100%)',
				'agora-gold-gradient': 'linear-gradient(135deg, #d4af37 0%, #cd7f32 100%)',
				'agora-shimmer': 'linear-gradient(90deg, transparent, rgba(212, 175, 55, 0.2), transparent)'
			},

			// Custom backdrop blur
			backdropBlur: {
				'xs': '2px'
			},

			// Custom line heights for discourse text
			lineHeight: {
				'discourse': '1.75'
			},

			// Custom z-index scale
			zIndex: {
				'modal': '100',
				'dropdown': '50',
				'header': '40',
				'overlay': '30'
			},

			// Custom container sizes
			maxWidth: {
				'agora': '1200px'
			},

			// Custom grid template columns for agora layout
			gridTemplateColumns: {
				'agora-layout': '300px 1fr 300px',
				'agora-spaces': 'repeat(auto-fit, minmax(320px, 1fr))'
			},

			// Custom aspect ratios
			aspectRatio: {
				'golden': '1.618 / 1'
			}
		}
	},
	plugins: [
		// Custom utility plugins
		function({ addUtilities, theme }) {
			const newUtilities = {
				// Agora-specific utilities
				'.text-agora-glow': {
					textShadow: '0 0 10px rgba(212, 175, 55, 0.3)'
				},
				'.border-agora-glow': {
					borderColor: theme('colors.agora-gold'),
					boxShadow: '0 0 0 1px rgba(212, 175, 55, 0.3)'
				},
				'.bg-agora-glass': {
					backgroundColor: 'rgba(26, 26, 46, 0.8)',
					backdropFilter: 'blur(10px)',
					border: '1px solid rgba(212, 175, 55, 0.2)'
				},
				'.gradient-text': {
					background: 'linear-gradient(135deg, #d4af37 0%, #cd7f32 100%)',
					WebkitBackgroundClip: 'text',
					WebkitTextFillColor: 'transparent',
					backgroundClip: 'text'
				},
				// Discourse-specific utilities
				'.discourse-depth-1': { paddingLeft: '1rem' },
				'.discourse-depth-2': { paddingLeft: '2rem' },
				'.discourse-depth-3': { paddingLeft: '3rem' },
				'.discourse-depth-4': { paddingLeft: '4rem' }
			}
			addUtilities(newUtilities)
		},

		// Typography plugin for better text handling
		function({ addComponents, theme }) {
			addComponents({
				'.agora-heading': {
					fontFamily: theme('fontFamily.serif'),
					fontWeight: '600',
					color: theme('colors.agora-gold'),
					textShadow: '0 0 10px rgba(212, 175, 55, 0.3)'
				},
				'.agora-body': {
					fontFamily: theme('fontFamily.sans'),
					color: theme('colors.agora-marble'),
					lineHeight: theme('lineHeight.discourse')
				},
				'.agora-card': {
					backgroundColor: 'rgba(22, 33, 62, 0.5)',
					borderRadius: theme('borderRadius.lg'),
					border: '1px solid rgba(212, 175, 55, 0.2)',
					backdropFilter: 'blur(10px)'
				},
				'.agora-button': {
					padding: '0.5rem 1rem',
					borderRadius: theme('borderRadius.md'),
					fontWeight: '500',
					transition: 'all 0.2s ease',
					border: '1px solid rgba(212, 175, 55, 0.3)',
					backgroundColor: 'rgba(212, 175, 55, 0.1)',
					color: theme('colors.agora-gold'),
					'&:hover': {
						backgroundColor: 'rgba(212, 175, 55, 0.2)',
						boxShadow: '0 0 15px rgba(212, 175, 55, 0.3)'
					}
				}
			})
		}
	]
}
