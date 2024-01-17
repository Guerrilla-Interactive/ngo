package files

import "text/template"

const indexQuery = `import { basePageQuery } from "@/sanity/shame-queries/base-page.query"
import type { InferType } from "groqd"
import { q } from "groqd"

export const {{.CamelCaseComponentName}}IndexQuery = q("*")
    .filterByType("{{.CamelCaseComponentName}}")
    .grab({
        title: q.string().optional(),
        ...basePageQuery,
    })
    .slice(0)
    .nullable()

export type {{.PascalCaseComponentName}}IndexQuery = NonNullable<InferType<typeof {{.CamelCaseComponentName}}IndexQuery>>
`

var IndexQuery = template.Must(template.New("indexQuery").Parse(indexQuery))
