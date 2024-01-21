package files

import "text/template"

const (
	slugSchema                 = ``
	slugSchemaCatchAll         = ``
	slugSchemaCatchAllOptional = ``
)

var (
	SlugSchema                 = template.Must(template.New("slugSchema").Parse(slugSchema))
	SlugSchemaCatchAll         = template.Must(template.New("slugSchemaCatchAll").Parse(slugSchemaCatchAll))
	SlugSchemaCatchAllOptional = template.Must(template.New("slugSchemaCatchAllOptional").Parse(slugSchemaCatchAllOptional))
)
