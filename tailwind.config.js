module.exports = {
  purge: [
    './docs/**/*.html',
  ],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      typography: {
        DEFAULT: {
          css: {
            maxWidth: '72ch',
            p: {
              '@apply break-words': '',
            },
            figure: {
              img: { margin: '0 auto' },
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
    extend: {}
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
};
