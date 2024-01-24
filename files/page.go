package files

import "text/template"

// page.tsx for index page
const page = `import { draftMode } from 'next/headers'
import { notFound } from "next/navigation"

import { runDraftQuery, runQuery } from '@/sanity/groqd-client'

import { {{.CamelCaseComponentName}}IndexQuery } from '../({{.KebabCaseComponentName}}-index-server)/{{.KebabCaseComponentName}}.index-query'
import {{.PascalCaseComponentName}}IndexBody from './{{.KebabCaseComponentName}}.index-component'
import { {{.PascalCaseComponentName}}Preview} from './{{.KebabCaseComponentName}}.index-preview'
import { generatePageMeta } from 'src/shame-utils/generate-page-meta-util'

export const generateMetadata = async () => {
  const data = await runQuery({{.CamelCaseComponentName}}IndexQuery, {})
  return generatePageMeta(data?.metadata)
}

const {{.PascalCaseComponentName}}IndexRoute = async () => {
  const { isEnabled: draftModeEnabled } = draftMode()
  const fetchClient = draftModeEnabled ? runDraftQuery : runQuery
  const data = await fetchClient(
    {{.CamelCaseComponentName}}IndexQuery,
    {},
  )

  if (!data) {
    return notFound()
  }

  if (draftModeEnabled) {
    return <{{.PascalCaseComponentName}}Preview initial={data} />
  }

  return <{{.PascalCaseComponentName}}IndexBody data={data} />
}

export default {{.PascalCaseComponentName}}IndexRoute`

var Page = template.Must(template.New("page").Parse(page))
