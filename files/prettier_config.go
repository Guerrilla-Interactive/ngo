package files

import "text/template"

const prettierConfig = `/** @type {import("prettier").Config} */
module.exports = {
  plugins: [require.resolve("prettier-plugin-tailwindcss")],
  printWidth: 80,
  tabWidth: 2,
  semi: false,
}`

var PrettierConfig = template.Must(template.New("prettierConfig").Parse(prettierConfig))
