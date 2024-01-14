/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [ "./**/*.html", "./**/*.templ", "./**/*.go", ],
  theme: {
    extend: {
      width: {
        '100px': '100px',
        '40px': '40px',
      },
      height: {
        '100px': '100px',
        '40px': '40px',
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-conic':
        'conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))',
      },
      colors : {
        'link-blue': '#38bdf8',
      },
      lineHeight: {
        'header': '1.35',
      },
      typography: {
        DEFAULT: {
          css: {
            color: "#FFF",
            a: {
              color: "#38bdf8",
              "&:hover": {
                fontWeight: "bold",
              },
            },
            h1: {
              color: "#FFF",
            },
            h2: {
              color: "#FFF",
            },
            h3: {
              color: "#FFF",
            },
            h4: { color: "#FFF" },
            em: { color: "#FFF" },
            strong: { color: "#FFF" },
            blockquote: { color: "#FFF" },
            'code::before': {
              content: '""',
            },
            'code::after': {
              content: '""',
            },
            code: {
              color: "#FFF",
              backgroundColor: "#343941",
              overFlowX: "scroll",
              fontWeight: "heavy",
              borderRadius: "0.375rem",
              padding: "4px",
            },
            pre: { backgroundColor: "#343941", color: "#FFF" },
          },
        },
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}

