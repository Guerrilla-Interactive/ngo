package files

import "text/template"

// page.tsx for dynamic optional catch all page
const slugPageCatchAllOptional = `import { draftMode } from 'next/headers'
import { notFound } from 'next/navigation'

import { runDraftQuery, runQuery } from '@/sanity/groqd-client'

import { {{.CamelCaseComponentName}}SlugQuery } from "../({{.KebabCaseComponentName}}-slug-server)/{{.KebabCaseComponentName}}.slug-query"
import { Preview{{.PascalCaseComponentName}}Slug } from "./{{.KebabCaseComponentName}}.slug-preview"
import {{.PascalCaseComponentName}}SlugBody from "./{{.KebabCaseComponentName}}.slug-component"
import { generatePageMeta } from 'src/shame-utils/generate-page-meta-util'


interface PageParams extends Record<string, string | Array<string>> {
	slug: Array<string>
}

interface PageProps {
	params: PageParams
}

export const generateMetadata = async ({ params }: PageProps) => {
	if (params?.slug && params?.slug.length > 0) {
		const data = await runQuery(
		{{.CamelCaseComponentName}}SlugQuery,
		{ slug: params.slug[0] }
		)
		return generatePageMeta(data?.metadata)
	}
}

const {{.PascalCaseComponentName}}SlugRoute = async ({ params }: PageProps) => {

  const { isEnabled: draftModeEnabled } = draftMode()
  const fetchClient = draftModeEnabled ? runDraftQuery : runQuery

  // Note this must be placed before accessing the first element of parms.slug in the catch all route
  // This is because this is an optional catchall route. See details here:
  // https://nextjs.org/docs/pages/building-your-application/routing/dynamic-routes
  if (!params?.slug || params?.slug?.length == 0) {
		return <div>THIS IS AN OPTIONAL CATCHALL ROUTE. WRITE YOUR LOGIC TO HANDLE NO SLUG PARARMS HERE</div>
  }

  const data = await fetchClient(
		{{.CamelCaseComponentName}}SlugQuery,
		{ slug: params.slug[0] }
  )

  if (!data) {
	return notFound()
  }
  

  if (draftModeEnabled) {
	  return <Preview{{.PascalCaseComponentName}}Slug initial={data} queryParams={params} slug={params.slug[0]!} />
  }

  return <{{.PascalCaseComponentName}}SlugBody data={data} />
}

export default {{.PascalCaseComponentName}}SlugRoute
`

var SlugPageCatchAllOptional = template.Must(template.New("slugPageCatchAllOptional").Parse(slugPageCatchAllOptional))
