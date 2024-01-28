/**
 * @type {import('next').NextConfig}
 */
const nextConfig = {
  output: 'export',

  // Optional: Change links `/me` -> `/me/` and emit `/me.html` -> `/me/index.html`
  // Uncomment the line below if you want to enable trailing slashes
  // trailingSlash: true,

  // Optional: Prevent automatic `/me` -> `/me/`, instead preserve `href`
  // Uncomment the line below if you want to skip automatic trailing slash redirects
  // skipTrailingSlashRedirect: true,

  // Optional: Change the output directory `out` -> `dist`
  // Uncomment the line below if you want to change the output directory to `dist`
  // distDir: 'dist',

  // Add domains for Next.js Image Optimization
  images: {
      domains: ['gophersignal-cloudflare-assets.s3.us-west-1.amazonaws.com'],
  },
}

module.exports = nextConfig;
