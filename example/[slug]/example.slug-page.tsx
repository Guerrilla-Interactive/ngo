"use client"

import { Container, Section } from "@/src/components/nextgen-core-ui"

import type { <%= UpperCaseComponentName %> SlugQuery } from "./(<%= camelCaseComponentName %>-server)/<%= camelCaseComponentName %>.slug-query"

export const <%= UpperCaseComponentName %> SlugPage = (page: <%= UpperCaseComponentName %> SlugQuery) => {
    {/* Add your component logic here */ }
    return (
        <>
            <Section className="min-h-screen bg-background">
                <Container className="py-14">
                    {/* Fill with your own content */}
                </Container>
            </Section>
        </>
    )
}
