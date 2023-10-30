import type { InferType } from "groqd"
import { q } from "groqd"

import { basePageQuery } from "@/lib/queries/utils/base-page.query"
import { notDraft } from "@/src/lib/sanity/not-draft.query"
import { sectionBlockQuery } from "@/src/components/sections/sections/section-block.queries"

export const <%= camelCaseComponentName %> SlugQuery = q("*")
    .filterByType("<%= lowerCaseComponentName %>")
    .filter(`${notDraft} && $slug == slug.current`)
    .grab({
        title: q.string().optional(),
        ...sectionBlockQuery,
        ...basePageQuery,
    })
    .slice(0)
    .nullable()

export type <%= UpperCaseComponentName %> SlugQuery = NonNullable < InferType < typeof <%= camelCaseComponentName %> SlugQuery >>
