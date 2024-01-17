package files

import "text/template"

const indexSchema = `import { defineType, defineField } from "sanity";

import type { CustomDocumentDefinition } from '@/sanity/api.desk-structure.ts'
import { defaultGroups } from '@/sanity/schema-utils/default-groups.util'
import { metaFields } from '@/sanity/schema-utils/generator-field/meta-fields.field'

export const {{.CamelCaseComponentName}}SlugSchema = defineType({
  type: "document",
  name: "{{.CamelCaseComponentName}}",
  title: "{{.PascalCaseComponentName}}",
  groups: defaultGroups,
  options: {
    previewable: true,
    linkable: true,
    isSingleton: true,
  },
  fields: [
    defineField({
      name: 'title',
      title: '{{.PascalCaseComponentName}} title',
      type: 'string',
      validation: (Rule) => Rule.required(),
      group: 'basic',
    }),
    ...metaFields({}),
  ],
  preview: {
    select: {
      title: "title",   
    },
    prepare({ title }) {
      return {
        title: title,
      };
    },
  },
}) as CustomDocumentDefinition
`

var IndexSchema = template.Must(template.New("indexSchema").Parse(indexSchema))
