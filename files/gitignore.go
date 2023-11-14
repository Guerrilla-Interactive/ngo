package files

import "text/template"

const gitignore = `.vercel
# env
.env
.env.local
`

var Gitignore = template.Must(template.New("gitignore").Parse(gitignore))
