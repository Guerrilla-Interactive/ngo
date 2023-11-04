package files

import "text/template"

const documents = `
{{range .Schemas}}export * from ".@{{.}}" 
{{end}}
`

var DocumentsSchemaQuery = template.Must(template.New("documents").Parse(documents))
