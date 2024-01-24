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

var PageIndexBody = template.Must(template.New("pageIndexBody").Parse(pageIndexBody))
