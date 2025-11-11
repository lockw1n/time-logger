/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{js,jsx,ts,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                bug: "#ef4444",
                feature: "#3b82f6",
                meeting: "#10b981",
                research: "#f59e0b",
            },
        },
    },
    plugins: [],
};
