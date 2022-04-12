module.exports = {
  content: ['layouts/**/*.html', 'content/**/*.html', 'assets/**/*.tsx', 'assets/**/*.ts'],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}
