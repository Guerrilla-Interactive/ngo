package files

import "text/template"

const slugQuery = `import type { InferType } from "groqd"
import { q } from "groqd"

import { basePageQuery } from "@/sanity/shame-queries/base-page.query"

export const {{.CamelCaseComponentName}}IndexQuery = q("*")
    .filterByType("{{.CamelCaseComponentName}}")
    .grab({
        title: q.string().optional(),
	slug: ["slug.current", q.string().optional()],
        ...basePageQuery,
    })
    .slice(0)

export type {{.PascalCaseComponentName}}IndexQuery = NonNullable<InferType<typeof {{.CamelCaseComponentName}}IndexQuery>>
`

// Note that the query template for normal dynamic route, catch all dynamic route and
// optional catch all dynamic route is all the same.
var (
	SlugQuery                 = template.Must(template.New("slugQuery").Parse(slugQuery))
	SlugQueryCatchAll         = template.Must(template.New("slugQuery").Parse(slugQuery))
	SlugQueryCatchAllOptional = template.Must(template.New("slugQuery").Parse(slugQuery))
)
