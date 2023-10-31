package files

import "text/template"

const page = `import { draftMode } from "next/headers"
import { notFound } from "next/navigation"

import { PreviewWrapper } from "@/components/preview-wrapper.component"
import { env } from "@/env/server.mjs"
import { generatePageMeta } from "@/lib/generate-page-meta.util"
import { tClient } from "@/lib/sanity/groqd-client"
import { {{.CamelCaseComponentName}}SlugQuery } from "../({{.KebabCaseComponentName}}-server)/{{.KebabCaseComponentName}}.slug-query"
import { {{.PascalCaseComponentName}}SlugPreview } from "./{{.KebabCaseComponentName}}.slug-preview"
import { {{.PascalCaseComponentName}}SlugPage } from "../{{.KebabCaseComponentName}}.slug-page"

type Props = {
    params: {
        slug: string
    }
}

export const generateMetadata = async ({ params }: Props) => {
    const data = await tClient({{.CamelCaseComponentName}} SlugQuery, params)
    return generatePageMeta(data?.metadata)
}

const {{.PascalCaseComponentName}} SlugRoute = async ({ params }: Props) => {
    const { isEnabled } = draftMode()
    const data = await tClient({{.CamelCaseComponentName}} SlugQuery, params)

    if (!data) {
        return notFound()
    }

    if (isEnabled) {
        return (
            <PreviewWrapper token={env.SANITY_API_TOKEN}>
		{{.PascalCaseComponentName}}SlugPreview initialData={data} queryParams={params} />
            </PreviewWrapper>
        )
    }

    return {{.PascalCaseComponentName}} SlugPage {...data } />
}

export default {{.PascalCaseComponentName}} SlugRoute`

var Page = template.Must(template.New("page").Parse(page))
