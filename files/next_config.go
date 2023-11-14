package files

import "text/template"

const nextConfig = `// @ts-check
/** @type {import("next").NextConfig} */
const config = {
  reactStrictMode: true,
  modularizeImports: {
    "@phospohor-icons/react": {
      transform: "@phosphor-icons/react/dist/icons/` + "{{`{{member}}`}}" + `",
	},
  },
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "cdn.sanity.io",
        port: "",
        pathname: "/**",
      },
    ],
  },
}
export default config`

var NextConfig = template.Must(template.New("nextConfig").Parse(nextConfig))
