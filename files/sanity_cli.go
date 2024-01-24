package files

import "text/template"

// const sanityCli = `import { env } from "@/env/server.mjs"
// import { defineCliConfig } from "sanity/cli"
//
// export default defineCliConfig({
//   api: {
//     projectId: env.NEXT_PUBLIC_SANITY_PROJECT_ID,
//     dataset: env.NEXT_PUBLIC_SANITY_DATASET,
//   },
// })`

var SanityCli = template.Must(template.New("nextConfig").Parse(nextConfig))
