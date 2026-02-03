/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Sora", "ui-sans-serif", "system-ui"],
        mono: ["JetBrains Mono", "ui-monospace", "SFMono-Regular"],
      },
      colors: {
        "panel": "rgba(255, 255, 255, 0.06)",
        "panel-border": "rgba(255, 255, 255, 0.12)",
      },
      boxShadow: {
        "glass": "0 20px 60px rgba(0,0,0,0.35)",
        "glow": "0 0 40px rgba(59,130,246,0.12)",
      },
      backgroundImage: {
        "grid": "linear-gradient(rgba(255,255,255,0.04) 1px, transparent 1px), linear-gradient(to right, rgba(255,255,255,0.04) 1px, transparent 1px)",
      },
    },
  },
  plugins: [],
};
