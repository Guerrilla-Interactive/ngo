package files

import "text/template"

// page.tsx for dynamic page
const slugPage = `import { draftMode } from 'next/headers'
import { notFound } from 'next/navigation'

import { runDraftQuery, runQuery } from '@/sanity/groqd-client'

import { {{.CamelCaseComponentName}}SlugQuery } from "../({{.KebabCaseComponentName}}-slug-server)/{{.KebabCaseComponentName}}.slug-query"
import { Preview{{.PascalCaseComponentName}}Slug } from "./{{.KebabCaseComponentName}}.slug-preview"
import {{.PascalCaseComponentName}}SlugBody  from "./{{.KebabCaseComponentName}}.slug-component"
import { generatePageMeta } from 'src/shame-utils/generate-page-meta-util'


interface PageParams extends Record<string, string> {
	slug: string
}

interface PageProps {
	params: PageParams
}

export const generateMetadata = async ({ params }: PageProps) => {
    const data = await runQuery(
		{{.CamelCaseComponentName}}SlugQuery,
		{ slug: params.slug }
		)
	return generatePageMeta(data?.metadata)
}

const {{.PascalCaseComponentName}}SlugRoute = async ({ params }: PageProps) => {

  const { isEnabled: draftModeEnabled } = draftMode()
  const fetchClient = draftModeEnabled ? runDraftQuery : runQuery

  const data = await fetchClient(
		{{.CamelCaseComponentName}}SlugQuery,
		{ slug: params.slug }
  )

  if (!data) {
	return notFound()
  }

  if (draftModeEnabled) {
	  return <Preview{{.PascalCaseComponentName}}Slug initial={data} queryParams={params} slug={params.slug} />
  }

  return <{{.PascalCaseComponentName}}SlugBody data={data} />
}

export default {{.PascalCaseComponentName}}SlugRoute
`

var SlugPage = template.Must(template.New("slugPage").Parse(slugPage))
