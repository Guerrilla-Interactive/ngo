package files

import "text/template"

const slug_query = ``

var SlugQuery = template.Must(template.New("slug_query").Parse(slug_query))
