package files

import "text/template"

const tsconfigJSON = `{
  "compilerOptions": {
    "target": "es2017",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "forceConsistentCasingInFileNames": true,
    "noEmit": true,
    "noUncheckedIndexedAccess": true,
    "incremental": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "node",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "baseUrl": ".",
    "paths": {
      "@/sanity/*": ["sanity/*"],
      "@/src/*": ["./src/*"],
      "@/*": ["./src/*"],
      "@/styles/*": ["./src/styles/*"],
      "@/schemaGen/*": ["./sanity/schemas/generator-field/*"],
      "@/public/*": ["./src/public/*"],
    },
    "plugins": [
      {
        "name": "next"
      }
    ]
  },
  "include": [
    "next-env.d.ts",
    "**/*.ts",
    "**/*.tsx",
    "**/*.cjs",
    "**/*.mjs",
    ".next/types/**/*.ts"
, "src/app/(site)/pieces/[slug]/components/mobile-product-page/(parts)/material-filter-tabs.componentstsx"  ],
  "exclude": ["node_modules"]
}`

var TSConfigJSON = template.Must(template.New("tsconfigJSON").Parse(tsconfigJSON))
