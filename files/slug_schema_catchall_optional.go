package files

import "text/template"

// Note that slug route schema and optional catch all slug route schema are the same
var SlugSchemaCatchAllOptional = template.Must(template.New("slugSchema").Parse(slugSchema))
