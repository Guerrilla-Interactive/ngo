"use client"
import { useLiveQuery } from "next-sanity/preview"

import { Occupier } from "@/components/utils/occupier.component"
import { <%= UpperCaseComponentName %> SlugPage } from "../<%= lowerCaseComponentName %>.slug-page"
import { <%= UpperCaseComponentName %> SlugQuery, <%= camelCaseComponentName %> SlugQuery } from "../(<%= lowerCaseComponentName %>-server)/<%= lowerCaseComponentName %>.slug-query"



export const <%= UpperCaseComponentName %> SlugPreview = ({
    initialData,
    queryParams,
}: {
    initialData: <%= UpperCaseComponentName %> SlugQuery
    queryParams: { slug: string }
}) => {
    const [data] = useLiveQuery(initialData, <%= camelCaseComponentName %> SlugQuery.query, queryParams)

    if (!data) {
        return <Occupier>No preview data found.</Occupier>
    }

    return <<%= UpperCaseComponentName %> SlugPage {...data } />
}
