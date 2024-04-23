/** @type {import('tailwindcss').Config} */
// npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch
module.exports = {
  content: ['./templates/**/*.html'],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('tailwindcss'),
  ],
}

