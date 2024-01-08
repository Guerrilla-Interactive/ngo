package files

import "text/template"

const slug_query = `import type { InferType } from "groqd"
import { q } from "groqd"

// TODO ADD YOUR QUERY
export const {{.CamelCaseComponentName}}SlugQuery = q("*")

export type {{.PascalCaseComponentName}}SlugQueryType = NonNullable<InferType<typeof {{.CamelCaseComponentName}}SlugQuery>>
`

var SlugQuery = template.Must(template.New("slug_query").Parse(slug_query))
