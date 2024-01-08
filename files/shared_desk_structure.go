package files

import "text/template"

const desk_structure = ``

var SharedDeskStructure = template.Must(template.New("desk_structure").Parse(desk_structure))
