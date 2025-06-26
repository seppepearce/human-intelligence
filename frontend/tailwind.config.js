/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{html,js,svelte,ts}"],
  darkMode: "class",
  theme: {
    extend: {
      // Vaporwave color palette
      colors: {
        // Primary vaporwave colors
        neon: {
          pink: "#ff006e",
          cyan: "#00f5ff",
          purple: "#8338ec",
          blue: "#3a86ff",
          green: "#06ffa5",
          yellow: "#ffbe0b",
          orange: "#fb8500",
        },

        // Dark theme backgrounds
        dark: {
          900: "#0a0a0a",
          800: "#1a1a1a",
          700: "#2a2a2a",
          600: "#3a3a3a",
          500: "#4a4a4a",
          400: "#5a5a5a",
        },

        // Marble/Greek inspired neutrals
        marble: {
          50: "#fdfcfa",
          100: "#f8f6f2",
          200: "#f1ede5",
          300: "#e8e2d7",
          400: "#ddd5c7",
          500: "#d0c6b6",
          600: "#c1b5a3",
          700: "#afa08c",
          800: "#968570",
          900: "#6b5d4f",
        },

        // Hacker terminal green
        terminal: {
          50: "#f0fff4",
          100: "#dcfce7",
          200: "#bbf7d0",
          300: "#86efac",
          400: "#4ade80",
          500: "#22c55e",
          600: "#16a34a",
          700: "#15803d",
          800: "#166534",
          900: "#14532d",
        },

        // Custom semantic colors
        primary: {
          50: "#fdf2f8",
          100: "#fce7f3",
          200: "#fbcfe8",
          300: "#f9a8d4",
          400: "#f472b6",
          500: "#ec4899", // Hot pink
          600: "#db2777",
          700: "#be185d",
          800: "#9d174d",
          900: "#831843",
        },

        secondary: {
          50: "#ecfeff",
          100: "#cffafe",
          200: "#a5f3fc",
          300: "#67e8f9",
          400: "#22d3ee",
          500: "#06b6d4", // Cyan
          600: "#0891b2",
          700: "#0e7490",
          800: "#155e75",
          900: "#164e63",
        },

        // Custom accent colors
        accent: {
          retro: "#ff006e",
          cyber: "#00f5ff",
          synthwave: "#8338ec",
          vapor: "#ff9a00",
        },
      },

      // Custom fonts for the aesthetic
      fontFamily: {
        sans: ["Inter", "system-ui", "sans-serif"],
        mono: ["JetBrains Mono", "Fira Code", "Consolas", "monospace"],
        display: ["Orbitron", "Inter", "sans-serif"], // Futuristic display font
        retro: ["Audiowide", "Inter", "sans-serif"], // 80s style font
        greek: ["Cinzel", "serif"], // Classical Greek style
      },

      // Custom spacing for the grid aesthetic
      spacing: {
        18: "4.5rem",
        88: "22rem",
        128: "32rem",
        144: "36rem",
      },

      // Grid patterns and backgrounds
      backgroundImage: {
        "grid-pattern":
          "linear-gradient(rgba(255, 0, 110, 0.1) 1px, transparent 1px), linear-gradient(90deg, rgba(255, 0, 110, 0.1) 1px, transparent 1px)",
        "vaporwave-gradient":
          "linear-gradient(45deg, #ff006e, #8338ec, #3a86ff)",
        "neon-gradient": "linear-gradient(90deg, #ff006e, #00f5ff)",
        "cyber-grid":
          "radial-gradient(circle at 50% 50%, rgba(0, 245, 255, 0.1) 0%, transparent 50%)",
        "retro-lines":
          "repeating-linear-gradient(90deg, transparent, transparent 98px, rgba(255, 0, 110, 0.3) 100px)",
        "marble-texture":
          "radial-gradient(circle at 20% 80%, rgba(120, 119, 198, 0.3), transparent 50%), radial-gradient(circle at 80% 20%, rgba(255, 255, 255, 0.15), transparent 50%)",
      },

      // Custom background sizes for patterns
      backgroundSize: {
        grid: "20px 20px",
        dots: "10px 10px",
        lines: "100px 100px",
      },

      // Animations for the retro aesthetic
      animation: {
        glow: "glow 2s ease-in-out infinite alternate",
        "neon-flicker": "neon-flicker 1.5s infinite linear",
        "cyber-pulse": "cyber-pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite",
        glitch: "glitch 0.3s ease-in-out infinite",
        float: "float 3s ease-in-out infinite",
        "scan-line": "scan-line 2s linear infinite",
        "type-writer": "type-writer 3s steps(40) 1s infinite",
      },

      // Keyframes for custom animations
      keyframes: {
        glow: {
          "0%": { boxShadow: "0 0 5px rgba(255, 0, 110, 0.5)" },
          "100%": {
            boxShadow:
              "0 0 20px rgba(255, 0, 110, 0.8), 0 0 30px rgba(255, 0, 110, 0.6)",
          },
        },
        "neon-flicker": {
          "0%, 19%, 21%, 23%, 25%, 54%, 56%, 100%": {
            textShadow:
              "0 0 5px rgba(0, 245, 255, 0.8), 0 0 10px rgba(0, 245, 255, 0.6), 0 0 15px rgba(0, 245, 255, 0.4)",
          },
          "20%, 24%, 55%": { textShadow: "none" },
        },
        "cyber-pulse": {
          "0%, 100%": { opacity: "1" },
          "50%": { opacity: "0.5" },
        },
        glitch: {
          "0%": { transform: "translate(0)" },
          "20%": { transform: "translate(-2px, 2px)" },
          "40%": { transform: "translate(-2px, -2px)" },
          "60%": { transform: "translate(2px, 2px)" },
          "80%": { transform: "translate(2px, -2px)" },
          "100%": { transform: "translate(0)" },
        },
        float: {
          "0%, 100%": { transform: "translateY(0px)" },
          "50%": { transform: "translateY(-10px)" },
        },
        "scan-line": {
          "0%": { transform: "translateY(-100vh)" },
          "100%": { transform: "translateY(100vh)" },
        },
        "type-writer": {
          "0%": { width: "0" },
          "50%": { width: "100%" },
          "100%": { width: "0" },
        },
      },

      // Custom blur effects
      blur: {
        xs: "2px",
      },

      // Drop shadow effects for neon look
      dropShadow: {
        "neon-sm": "0 0 2px rgba(255, 0, 110, 0.8)",
        neon: "0 0 5px rgba(255, 0, 110, 0.8)",
        "neon-lg": "0 0 10px rgba(255, 0, 110, 0.8)",
        cyber: "0 0 5px rgba(0, 245, 255, 0.8)",
        "cyber-lg": "0 0 10px rgba(0, 245, 255, 0.8)",
      },

      // Box shadows for depth and glow
      boxShadow: {
        neon: "0 0 10px rgba(255, 0, 110, 0.5)",
        "neon-lg":
          "0 0 20px rgba(255, 0, 110, 0.6), 0 0 40px rgba(255, 0, 110, 0.4)",
        cyber: "0 0 10px rgba(0, 245, 255, 0.5)",
        "cyber-lg":
          "0 0 20px rgba(0, 245, 255, 0.6), 0 0 40px rgba(0, 245, 255, 0.4)",
        "inner-glow": "inset 0 0 10px rgba(255, 0, 110, 0.3)",
        terminal: "0 0 15px rgba(34, 197, 94, 0.4)",
      },

      // Border radius for retro shapes
      borderRadius: {
        "4xl": "2rem",
      },

      // Custom border widths
      borderWidth: {
        3: "3px",
      },

      // Typography scale adjustments
      fontSize: {
        "2xs": ["0.625rem", { lineHeight: "0.875rem" }],
        "3xl": ["1.875rem", { lineHeight: "2.25rem" }],
        "4xl": ["2.25rem", { lineHeight: "2.5rem" }],
        "5xl": ["3rem", { lineHeight: "1" }],
        "6xl": ["3.75rem", { lineHeight: "1" }],
        "7xl": ["4.5rem", { lineHeight: "1" }],
        "8xl": ["6rem", { lineHeight: "1" }],
        "9xl": ["8rem", { lineHeight: "1" }],
      },

      // Custom aspect ratios for content
      aspectRatio: {
        "4/3": "4 / 3",
        "3/2": "3 / 2",
        "2/3": "2 / 3",
        "9/16": "9 / 16",
      },

      // Z-index scale
      zIndex: {
        60: "60",
        70: "70",
        80: "80",
        90: "90",
        100: "100",
      },
    },
  },

  plugins: [
    require("@tailwindcss/forms"),
    require("@tailwindcss/typography"),

    // Custom plugin for vaporwave utilities
    function ({ addUtilities, addComponents, theme }) {
      // Utility classes for the vaporwave aesthetic
      addUtilities({
        ".text-glow": {
          textShadow: "0 0 10px currentColor",
        },
        ".text-neon": {
          textShadow: `
						0 0 5px ${theme("colors.neon.pink")},
						0 0 10px ${theme("colors.neon.pink")},
						0 0 15px ${theme("colors.neon.pink")}
					`,
        },
        ".text-cyber": {
          textShadow: `
						0 0 5px ${theme("colors.neon.cyan")},
						0 0 10px ${theme("colors.neon.cyan")},
						0 0 15px ${theme("colors.neon.cyan")}
					`,
        },
        ".bg-grid": {
          backgroundImage:
            "linear-gradient(rgba(255, 0, 110, 0.1) 1px, transparent 1px), linear-gradient(90deg, rgba(255, 0, 110, 0.1) 1px, transparent 1px)",
          backgroundSize: "20px 20px",
        },
        ".bg-cyber-grid": {
          backgroundImage:
            "radial-gradient(circle at 50% 50%, rgba(0, 245, 255, 0.1) 0%, transparent 50%)",
          backgroundSize: "30px 30px",
        },
        ".border-neon": {
          borderColor: theme("colors.neon.pink"),
          boxShadow: `0 0 10px ${theme("colors.neon.pink")}`,
        },
        ".border-cyber": {
          borderColor: theme("colors.neon.cyan"),
          boxShadow: `0 0 10px ${theme("colors.neon.cyan")}`,
        },
      });

      // Component classes for common UI patterns
      addComponents({
        ".btn-neon": {
          backgroundColor: "transparent",
          border: `2px solid ${theme("colors.neon.pink")}`,
          color: theme("colors.neon.pink"),
          padding: `${theme("spacing.2")} ${theme("spacing.6")}`,
          borderRadius: theme("borderRadius.md"),
          textTransform: "uppercase",
          fontWeight: theme("fontWeight.semibold"),
          letterSpacing: theme("letterSpacing.wide"),
          transition: "all 0.3s ease",
          "&:hover": {
            backgroundColor: theme("colors.neon.pink"),
            color: theme("colors.dark.900"),
            boxShadow: `0 0 20px ${theme("colors.neon.pink")}`,
          },
        },
        ".btn-cyber": {
          backgroundColor: "transparent",
          border: `2px solid ${theme("colors.neon.cyan")}`,
          color: theme("colors.neon.cyan"),
          padding: `${theme("spacing.2")} ${theme("spacing.6")}`,
          borderRadius: theme("borderRadius.md"),
          textTransform: "uppercase",
          fontWeight: theme("fontWeight.semibold"),
          letterSpacing: theme("letterSpacing.wide"),
          transition: "all 0.3s ease",
          "&:hover": {
            backgroundColor: theme("colors.neon.cyan"),
            color: theme("colors.dark.900"),
            boxShadow: `0 0 20px ${theme("colors.neon.cyan")}`,
          },
        },
        ".card-vapor": {
          backgroundColor: "rgba(26, 26, 26, 0.8)",
          border: `1px solid ${theme("colors.dark.600")}`,
          borderRadius: theme("borderRadius.lg"),
          backdropFilter: "blur(10px)",
          boxShadow: "0 8px 32px rgba(0, 0, 0, 0.3)",
        },
        ".input-neon": {
          backgroundColor: "transparent",
          border: `2px solid ${theme("colors.dark.600")}`,
          borderRadius: theme("borderRadius.md"),
          color: theme("colors.white"),
          padding: theme("spacing.3"),
          width: "100%",
          transition: "all 0.3s ease",
          "&:focus": {
            borderColor: theme("colors.neon.pink"),
            boxShadow: `0 0 10px ${theme("colors.neon.pink")}`,
            outline: "none",
          },
        },
      });
    },
  ],
};
