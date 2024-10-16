/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../templ/**/*.{templ,go}"],
  theme: {
    extend: {
        colors:{
            background: {
                DEFAULT: "var(--background)",
                muted: "var(--background-muted)"
            },
            foreground: "var(--foreground)",
            border: "var(--border)",
            button: "var(--button)"
       }
    },
  },
  plugins: [],
}

