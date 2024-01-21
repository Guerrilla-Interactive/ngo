package files

import "text/template"

// Component for slug page.tsx
const pageSlugBody = `import { {{.PascalCaseComponentName}}SlugQuery } from "../({{.KebabCaseComponentName}}-slug-server)/{{.KebabCaseComponentName}}.slug-query"

interface PageProps {
	data: {{.PascalCaseComponentName}}SlugQuery
}
export default function {{.PascalCaseComponentName}}SlugBody(props: PageProps) {
	return (
		<div>Title: {props.data.title}</div>
	)
}`

const (
	pageSlugBodyCatchAll         = pageSlugBody
	pageSlugBodyCatchAllOptional = pageSlugBody
)

var (
	PageSlugBody                 = template.Must(template.New("pageSlugBody").Parse(pageSlugBody))
	PageSlugBodyCatchAll         = template.Must(template.New("pageSlugBodyCatchAll").Parse(pageSlugBodyCatchAll))
	PageSlugBodyCatchAllOptional = template.Must(template.New("pageSlugBodyCatchAllOptional").Parse(pageSlugBodyCatchAllOptional))
)
