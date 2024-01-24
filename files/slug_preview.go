package files

import "text/template"

const slugPreview = `"use client"
import { PreviewLoadingErrorHOC } from '@/components/preview/preview-wrapper'
import {{.PascalCaseComponentName}}SlugBody from "./{{.KebabCaseComponentName}}.slug-component"
import { {{.PascalCaseComponentName}}SlugQuery,  {{.CamelCaseComponentName}}SlugQuery } from "../({{.KebabCaseComponentName}}-slug-server)/{{.KebabCaseComponentName}}.slug-query"


interface PreviewProps {
	initial: {{.PascalCaseComponentName}}SlugQuery
	queryParams?: Record<string, string|Array<string>>
	slug?: string
}

export function Preview{{.PascalCaseComponentName}}Slug(props: PreviewProps){
	return (
		<PreviewLoadingErrorHOC
			initial={props.initial}
			query={ {{.CamelCaseComponentName}}SlugQuery.query}
			slug={props.slug}
			successFn={(data) =>
				<{{.PascalCaseComponentName}}SlugBody data={data} />
			} />
	)
}`

const (
	// Dynamic routes catch all and catch all optional are all the same
	slugPreviewCatchAll         = slugPreview
	slugPreviewCatchAllOptional = slugPreview
)

var (
	SlugPreview                 = template.Must(template.New("slugPreview").Parse(slugPreview))
	SlugPreviewCatchAll         = template.Must(template.New("slugPreviewCatchAll").Parse(slugPreviewCatchAll))
	SlugPreviewCatchAllOptional = template.Must(template.New("slugPreviewCatchAllOptional").Parse(slugPreviewCatchAllOptional))
)
