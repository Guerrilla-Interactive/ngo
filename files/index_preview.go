package files

import "text/template"

const indexPreview = `"use client"

import { PreviewLoadingErrorHOC } from "@/src/components/sanity/preview-loading-error-hoc"

import { {{.PascalCaseComponentName}}IndexPage } from "../../pieces.index-page"
import type { PiecesIndexQuery } from "../(pieces-index-server)/pieces.index-query"
import { piecesIndexQuery } from "../(pieces-index-server)/pieces.index-query"

export const PiecesIndexPreview = ({
  initial,
  queryParams
}: {
  initial: PiecesIndexQuery
  queryParams: Record<string, string>
}) => {
  return (
    <PreviewLoadingErrorHOC
      initial={initial}
      query={piecesIndexQuery.query}
      successFn={(data) =>
        <PiecesIndexPage {...data} />
      }
    />
  )
}`

var IndexPreview = template.Must(template.New("indexPreview").Parse(indexPreview))
