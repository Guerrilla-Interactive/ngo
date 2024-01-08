package files

import "text/template"

const slugPage = `import { draftMode } from "next/headers"
import { notFound } from "next/navigation"

import { tClient, tClientDraft } from "@/sanity/groqd-client"
import { generatePageMeta } from "@/src/lib/generate-page-meta.util"

import { {{.CamelCaseComponentName}}SlugQuery } from "../({{.KebabCaseComponentName}}-slug-server)/{{.KebabCaseComponentName}}.slug-query"
import { {{.PascalCaseComponentName}}SlugPreview } from "./{{.KebabCaseComponentName}}.slug-preview"
import { {{.PascalCaseComponentName}}SlugBody } from "./{{.KebabCaseComponentName}}.body"


type Props = {
  params: {
    slug: string
  }
}

export const generateMetadata = async ({ params }: Props) => {
  const data = await tClient({{.CamelCaseComponentName}}SlugQuery, params)
  return generatePageMeta(data?.metadata)
}

const {{.PascalCaseComponentName}}Page = async ({ params }: Props) => {

  const data = await tClient({{.CamelCaseComponentName}}SlugQuery, params)
  const draftData = await tClientDraft({{.CamelCaseComponentName}}SlugQuery, params)

  // Return not found if the {{.KebabCaseComponentName}} slug is invalid
  if (!data)
    return notFound()

  // TODO
  // Insert any necessary custom logic here

  return draftMode().isEnabled ? (
    <{{.PascalCaseComponentName}}SlugPreview initial={draftData!} queryParams={params} />
  ) : (
		// TODO
		<{{.PascalCaseComponentName}}PageBody PASS YOUR CUSTOM PROPS HERE />
  )
}

export default {{.PascalCaseComponentName}}Page
`

const slugPageCatchAll = ``

var (
	SlugPage          = template.Must(template.New("slugPage").Parse(slugPage))
	SlugPageCatchAlll = template.Must(template.New("slugPageCatchAll").Parse(slugPageCatchAll))
)
