package files

import "text/template"

const sanityEnv = `import type { Theme } from "@sanity/ui"

declare module "styled-components" {
  // eslint-disable-next-line @typescript-eslint/no-empty-interface
  interface DefaultTheme extends Theme {}
}

interface FieldOptions {
  required?: boolean
}

declare module "sanity" {
  export interface DocumentOptions {
    previewable?: boolean
    linkable?: boolean
    singleton?: boolean
  }

  export interface StringOptions extends FieldOptions {}
  export interface NumberOptions extends FieldOptions {}
  export interface SlugOptions extends FieldOptions {}
  export interface TextOptions extends FieldOptions {}
  export interface ArrayOptions extends FieldOptions {}
  export interface ImageOptions extends FieldOptions {}
  export interface ObjectOptions extends FieldOptions {}
  export interface DatetimeOptions extends FieldOptions {}
  export interface ReferenceBaseOptions extends FieldOptions {}
}`

var SanityEnv = template.Must(template.New("sanityEnv").Parse(sanityEnv))
