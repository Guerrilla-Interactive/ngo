package files

import "text/template"

// Note that slug route schema and catch all slug route schema are the same
var SlugSchemaCatchAll = template.Must(template.New("slugSchema").Parse(slugSchema))
