package files

import "text/template"

const slugQuery = ``

const (
	slugQueryCatchAll         = ``
	slugQueryCatchAllOptional = ``
)

var (
	SlugQuery                 = template.Must(template.New("slugQuery").Parse(slugQuery))
	SlugQueryCatchAll         = template.Must(template.New("slugQueryCatchAll").Parse(slugQueryCatchAll))
	SlugQueryCatchAllOptional = template.Must(template.New("slugQueryCatchAllOptional").Parse(slugQueryCatchAllOptional))
)
