package files

import "text/template"

const page = ``

var Page = template.Must(template.New("page").Parse(page))
