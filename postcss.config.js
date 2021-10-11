module.exports = {
  plugins: [
    require('autoprefixer'),
    require('postcss-nested'),
    require('tailwindcss')('./tailwind.config.js'),
  ]
};
