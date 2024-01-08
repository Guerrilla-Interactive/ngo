package files

import "text/template"

const slug_schema = ``

var SlugSchema = template.Must(template.New("slug_schema").Parse(slug_schema))
