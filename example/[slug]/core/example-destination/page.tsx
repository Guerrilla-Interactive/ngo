import { draftMode } from "next/headers"
import { notFound } from "next/navigation"

import { PreviewWrapper } from "@/components/preview-wrapper.component"
import { env } from "@/env/server.mjs"
import { generatePageMeta } from "@/lib/generate-page-meta.util"
import { tClient } from "@/lib/sanity/groqd-client"
import { <%= camelCaseComponentName %> SlugQuery } from "../(<%= lowerCaseComponentName %>-server)/<%= lowerCaseComponentName %>.slug-query"
import { <%= UpperCaseComponentName %> SlugPreview } from "./<%= lowerCaseComponentName %>.slug-preview"
import { <%= UpperCaseComponentName %> SlugPage } from "../<%= lowerCaseComponentName %>.slug-page"

type Props = {
    params: {
        slug: string
    }
}

export const generateMetadata = async ({ params }: Props) => {
    const data = await tClient(<%= camelCaseComponentName %> SlugQuery, params)
    return generatePageMeta(data?.metadata)
}

const <%= UpperCaseComponentName %> SlugRoute = async ({ params }: Props) => {
    const { isEnabled } = draftMode()
    const data = await tClient(<%= camelCaseComponentName %> SlugQuery, params)

    if (!data) {
        return notFound()
    }

    if (isEnabled) {
        return (
            <PreviewWrapper token={env.SANITY_API_TOKEN}>
        <<%= UpperCaseComponentName %>SlugPreview initialData={data} queryParams={params} />
            </PreviewWrapper>
        )
    }

    return <<%= UpperCaseComponentName %> SlugPage {...data } />
}

export default <%= UpperCaseComponentName %> SlugRoute
