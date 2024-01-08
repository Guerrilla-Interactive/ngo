package files

import "text/template"

const shared_query = ``

var SharedQuery = template.Must(template.New("shared_query").Parse(shared_query))
