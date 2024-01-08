package files

import "text/template"

const slugPreview = `"use client"

import { PreviewLoadingErrorHOC } from "@/src/components/sanity/preview-loading-error-hoc"

import type { {{.PascalCaseComponentName}}SlugQueryType } from "../({{.KebabCaseComponentName}}-slug-server)/{{.KebabCaseComponentName}}.slug-query";
import { {{.CamelCaseComponentName}}SlugQuery } from "../({{.KebabCaseComponentName}}-slug-server)/{{.KebabCaseComponentName}}.slug-query"
import { {{.PascalCaseComponentName}}SlugBody } from "./{{.KebabCaseComponentName}}.body"

export const {{.PascalCaseComponentName}}SlugPreview = ({
	initial,
	queryParams,
}: {
	initial: {{.PascalCaseComponentName}}SlugQueryType
	queryParams: { slug: string }
}) => {

	return (
		<PreviewLoadingErrorHOC
			initial={initial}
			query={ {{.CamelCaseComponentName}}SlugQuery.query} 
			queryParams={queryParams}
			successFn={(data) => {
				// TODO
				return (
				// TODO
				<{{.PascalCaseComponentName}}PageBody PASS YOUR CUSTOM PROPS HERE />
				)
			}
			}
		/>
	)
}
`

var SlugPreview = template.Must(template.New("slugPreview").Parse(slugPreview))
