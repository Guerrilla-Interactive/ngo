package files

import "text/template"

const pageIndexBody = `export function {{.PascalCaseComponentName}}IndexBody() {
	return (
		<>Modify this component by the name {{.PascalCaseComponentName}}</>
	)}`

const pageSlugBody = `export function {{.PascalCaseComponentName}}SlugBody() {
	return (
		<>Modify this component by the name {{.PascalCaseComponentName}}</>
	)}`

var (
	PageIndexBody = template.Must(template.New("pageIndexBody").Parse(pageIndexBody))
	PageSlugBody  = template.Must(template.New("pageSlugBody").Parse(pageSlugBody))
)
