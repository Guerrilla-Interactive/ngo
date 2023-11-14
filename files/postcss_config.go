package files

import "text/template"

const postcss = `module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};`

var Postcss = template.Must(template.New("postcss").Parse(postcss))
