package files

import "text/template"

const pageIndexBody = `import type { {{.PascalCaseComponentName}}IndexQuery } from "../({{.KebabCaseComponentName}}-index-server)/{{.KebabCaseComponentName}}.index-query"

interface PageProps {
  data: {{.PascalCaseComponentName}}IndexQuery
}
export default function {{.PascalCaseComponentName}}IndexBody(props: PageProps) {
  return (
    <div>{{.PascalCaseComponentName}}: {props.data.title}</div>
  )
}`

// Component for slug page.tsx
const (
	pageSlugBody                 = ``
	pageSlugBodyCatchAll         = ``
	pageSlugBodyCatchAllOptional = ``
)

var (
	PageIndexBody                = template.Must(template.New("pageIndexBody").Parse(pageIndexBody))
	PageSlugBody                 = template.Must(template.New("pageSlugBody").Parse(pageSlugBody))
	PageSlugBodyCatchAll         = template.Must(template.New("pageSlugBodyCatchAll").Parse(pageSlugBodyCatchAll))
	PageSlugBodyCatchAllOptional = template.Must(template.New("pageSlugBodyCatchAllOptional").Parse(pageSlugBodyCatchAllOptional))
)
