package files

import "text/template"

const npmrc = ``

var Npmrc = template.Must(template.New("npmrc").Parse(npmrc))
