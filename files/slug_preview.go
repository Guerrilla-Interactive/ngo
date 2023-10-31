package files

import "text/template"

const slugPreview = `
"use client"
import { useLiveQuery } from "next-sanity/preview"

import { Occupier } from "@/components/utils/occupier.component"

import { {{.PascalCaseComponentName}}Component } from "./{{.KebabCaseComponentName}}.component"
import type { {{.PascalCaseComponentName}}SlugQuery } from "../({{.KebabCaseComponentName}}-server)/{{.KebabCaseComponentName}}.slug-query"
import { {{.CamelCaseComponentName}}SlugQuery } from "../({{.KebabCaseComponentName}}-server)/{{.KebabCaseComponentName}}.slug-query"

export const {{.PascalCaseComponentName}}Preview = ({
  initialData,
  queryParams,
}: {
  initialData: {{.PascalCaseComponentName}}Query
  queryParams: { slug: string }
}) => {
  const [data] = useLiveQuery(initialData, {{.CamelCaseComponentName}}Query.query, queryParams)

  if (!data) {
    return <Occupier container>No preview data found.</Occupier>
  }

  return <{{.PascalCaseComponentName}}Component {...data} />
}
`

var SlugPreview = template.Must(template.New("slugPreview").Parse(slugPreview))
