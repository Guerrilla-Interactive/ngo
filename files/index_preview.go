package files

import "text/template"

const indexPreview = `"use client"
import { PreviewLoadingErrorHOC } from '@/components/preview/preview-wrapper'

import type { {{.PascalCaseComponentName}}IndexQuery} from '../({{.KebabCaseComponentName}}-index-server)/{{.KebabCaseComponentName}}.index-query';
import { {{.CamelCaseComponentName}}IndexQuery } from '../({{.KebabCaseComponentName}}-index-server)/{{.KebabCaseComponentName}}.index-query'
import {{.PascalCaseComponentName}}IndexBody from './{{.KebabCaseComponentName}}.index-component'

interface PreviewProps {
  initial: {{.PascalCaseComponentName}}IndexQuery
  queryParams?: { slug: string }
}

export function {{.PascalCaseComponentName}}Preview(props: PreviewProps) {
  return (
    <PreviewLoadingErrorHOC 
	initial={props.initial}
	query={ {{.CamelCaseComponentName}}IndexQuery.query} 
	successFn={(data) =>
		<{{.PascalCaseComponentName}}IndexBody data={data} />}
	/>
  )
}
`

var IndexPreview = template.Must(template.New("indexPreview").Parse(indexPreview))
