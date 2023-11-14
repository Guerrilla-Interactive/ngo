package files

import "text/template"

const packageJSON = `{
  "name": "{{.PackageName }}",
  "version": "1.0.0",
  "description": "",
  "main": "index.js",
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "format": "pnpm exec prettier --write . --ignore-path .gitignore",
    "lint": "next lint"
  },
  "keywords": [],
  "author": "Guerrilla",
  "license": "UNLICENSED",
  "dependencies": {
  },
  "devDependencies": {
  },
  "engines": {
    "node": ">=18"
  }
}`

var PackageJSON = template.Must(template.New("packageJSON").Parse(packageJSON))
