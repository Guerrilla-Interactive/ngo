package files

import "text/template"

const query = `import type { InferType } from "groqd"
import { q } from "groqd"

import { basePageQuery } from "@/lib/queries/utils/base-page.query"
import { notDraft } from "@/src/lib/sanity/not-draft.query"
import { sectionBlockQuery } from "@/src/components/sections/sections/section-block.queries"

export const {{.CamelCaseComponentName}}SlugQuery = q("*")
    .filterByType("{{.CamelCaseComponentName}}")
    .filter(` + "`${notDraft} && $slug == slug.current`" + `)
    .grab({
        title: q.string().optional(),
        ...sectionBlockQuery,
        ...basePageQuery,
    })
    .slice(0)
    .nullable()

export type {{.PascalCaseComponentName}}SlugQuery = NonNullable<InferType<typeof {{.CamelCaseComponentName}}SlugQuery>>
`

var Query = template.Must(template.New("query").Parse(query))
