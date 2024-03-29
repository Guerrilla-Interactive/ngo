package files

import "text/template"

const slugQuery = `import type { InferType } from "groqd"
import { q } from "groqd"

import { basePageQuery } from "@/sanity/shame-queries/base-page.query"

export const {{.CamelCaseComponentName}}SlugQuery = q("*")
    .filterByType("{{.CamelCaseComponentName}}")
    .filter("slug.current == $slug")
    .grab({
        title: q.string().optional(),
        slug: ["slug.current", q.string().optional()],
        ...basePageQuery,
    })
    .slice(0)
    .nullable()

export type {{.PascalCaseComponentName}}SlugQuery = NonNullable<InferType<typeof {{.CamelCaseComponentName}}SlugQuery>>
`

// Note that the query template for normal dynamic route, catch all dynamic route and
// optional catch all dynamic route is all the same.
var (
	SlugQuery                 = template.Must(template.New("slugQuery").Parse(slugQuery))
	SlugQueryCatchAll         = template.Must(template.New("slugQuery").Parse(slugQuery))
	SlugQueryCatchAllOptional = template.Must(template.New("slugQuery").Parse(slugQuery))
)
