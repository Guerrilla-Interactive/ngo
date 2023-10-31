package files

import "text/template"

const querySchema = `import { defineType } from "sanity";
import { metaFields } from "@/sanity/schemas/generator-field/meta-fields.field";
import { sectionsField } from "@/sanity/schemas/generator-field/sections.field";
import { slugField } from "@/sanity/schemas/generator-field/slug.field";
import { stringField } from "@/sanity/schemas/generator-field/string.field";
import { defaultGroups } from "@/sanity/schemas/utils/default-groups.util";

export const {{.CamelCaseComponentName}}SlugSchema = defineType({
  type: "document",
  name: "{{.CamelCaseComponentName}}",
  title: "{{.PascalCaseComponentName}}",
  groups: defaultGroups,
  options: {
    previewable: true,
    linkable: true,
  },
  fields: [
    stringField({
      name: "title",
      title: "{{.PascalCaseComponentName}} title",
      required: true,
    }),
    slugField({}),
    sectionsField({}),
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
});
`

var QuerySchema = template.Must(template.New("querySchema").Parse(querySchema))
