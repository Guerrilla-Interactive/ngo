package files

import "text/template"

const eslintrc = `{
  "plugins": ["@typescript-eslint", "tailwindcss", "simple-import-sort"],
  "extends": [
    "next/core-web-vitals",
    "plugin:@typescript-eslint/recommended",
    "plugin:tailwindcss/recommended",
    "prettier"
  ],
  "rules": {
    "@typescript-eslint/consistent-type-imports": "warn",
    "simple-import-sort/imports": "warn",
    "simple-import-sort/exports": "warn",
    "react-hooks/exhaustive-deps": "error"
    
  },
  "settings": {
    "tailwindcss": {
      "callees": ["cn"],
      "config": "./tailwind.config.cjs"
    }
  }
}`

var Eslintrc = template.Must(template.New("eslintrc").Parse(eslintrc))
