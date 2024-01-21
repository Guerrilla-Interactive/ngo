package files

import "text/template"

// page.tsx for index page
const page = `import { draftMode } from 'next/headers'
import { notFound } from "next/navigation"

import { serverEnv } from '@/lib/env/server'
import { runQuery } from '@/sanity/groqd-client'

import { {{.CamelCaseComponentName}}IndexQuery } from '../({{.KebabCaseComponentName}}-index-server)/{{.KebabCaseComponentName}}.index-query'
import {{.PascalCaseComponentName}}IndexBody from './{{.KebabCaseComponentName}}.index-component'
import { {{.PascalCaseComponentName}}Preview} from './{{.KebabCaseComponentName}}.index-preview'

const {{.PascalCaseComponentName}}IndexRoute = async () => {
  const { isEnabled: draftModeEnabled } = draftMode()
  const token = serverEnv.SANITY_API_READ_TOKEN
  const data = await runQuery(
    {{.CamelCaseComponentName}}IndexQuery,
    {},
    draftModeEnabled ? token : undefined
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
