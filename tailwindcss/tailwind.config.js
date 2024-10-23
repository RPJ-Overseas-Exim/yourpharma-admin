/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["../templ/**/*.{templ,go}"],
    theme: {
        extend: {
            colors: {
                background: {
                    DEFAULT: "var(--background)",
                    muted: "var(--background-muted)"
                },
                foreground: {
                    DEFAULT: "var(--foreground)",
                    muted: "var(--foreground-muted)"
                },
                border: {
                    DEFAULT: "var(--border)",
                    muted: "var(--border-muted)"
                },
                button: {
                    DEFAULT: "var(--button)",
                    muted: "var(--button-muted)"
                },
            }
        },
    },
    plugins: [],
}

