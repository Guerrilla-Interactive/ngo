package files

import "text/template"

const slugPage = ``

var SlugPage = template.Must(template.New("slugPage").Parse(slugPage))
