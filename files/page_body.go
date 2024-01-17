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

const pageSlugBody = `export function {{.PascalCaseComponentName}}SlugBody() {
	return (
		<>Modify this component by the name {{.PascalCaseComponentName}}</>
	)}`

var (
	PageIndexBody = template.Must(template.New("pageIndexBody").Parse(pageIndexBody))
	PageSlugBody  = template.Must(template.New("pageSlugBody").Parse(pageSlugBody))
)
