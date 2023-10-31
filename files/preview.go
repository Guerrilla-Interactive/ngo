package files

import "text/template"

const preview = ``

var Preview = template.Must(template.New("preview").Parse(preview))
