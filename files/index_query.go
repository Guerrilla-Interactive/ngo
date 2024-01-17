package files

import "text/template"

const indexQuery = `import type { InferType } from "groqd"
import { q } from "groqd"

import { basePageQuery } from "@/sanity/shame-queries/base-page.query"

export const {{.CamelCaseComponentName}}IndexQuery = q("*")
    .filterByType("{{.CamelCaseComponentName}}")
    .grab({
        title: q.string().optional(),
        ...basePageQuery,
    })
    .slice(0)

export type {{.PascalCaseComponentName}}IndexQuery = NonNullable<InferType<typeof {{.CamelCaseComponentName}}IndexQuery>>
`

var IndexQuery = template.Must(template.New("indexQuery").Parse(indexQuery))
