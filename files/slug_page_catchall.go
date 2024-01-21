package files

import "text/template"

// page.tsx for dynamic catch all page
const slugPageCatchAll = `import { draftMode } from 'next/headers'
import { notFound } from 'next/navigation'
import type { ZodType } from 'zod'

import { serverEnv } from '@/lib/env/server'
import { runQuery } from '@/sanity/groqd-client'
import { {{.CamelCaseComponentName}}SlugQuery {{.CamelCaseComponentName}}SlugQuery } from "../({{.KebabCaseComponentName}}-slug-server)/{{.KebabCaseComponentName}}.slug-query"
import { Preview{{.PascalCaseComponentName}}Slug } from "./{{.KebabCaseComponentName}}.slug-preview"
import { {{.PascalCaseComponentName}}SlugBody } from "./{{.KebabCaseComponentName}}.slug-component"


interface PageParams extends Record<string, any> {
	slug: string
}

interface PageProps {
	params: PageParams
}

const {{.PascalCaseComponentName}}SlugRoute = async ({ params }: Props) => {

  const { isEnabled: draftModeEnabled } = draftMode()
  const token = serverEnv.SANITY_API_READ_TOKEN

  const data = await runQuery<ZodType<ProductsSlugQuery | null>>(
		{{.CamelCaseComponentName}}SlugQuery,
		{ slug: props.params.slug },
		draftModeEnabled ? token : undefined
  )

  if (!data) {
	return notFound()
  }

  if (draftModeEnabled) {
	  return <{ Preview{{.PascalCaseComponentName}}Slug initial={data} />
  }

  return <{{.PascalCaseComponentName}}SlugBody data={data} />
}

export default {{.PascalCaseComponentName}}SlugRoute
`

var SlugPageCatchAll = template.Must(template.New("slugPageCatchAll").Parse(slugPageCatchAll))
