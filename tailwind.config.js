module.exports = {
  purge: [
    './www/**/*.html',
  ],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      typography: {
        DEFAULT: {
          css: {
            maxWidth: '72ch',
            figure: {
              img: {
                margin: '0 auto',
              },
              figcaption: {
                textAlign: 'center',
              },
            },
          },
        },
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
};
