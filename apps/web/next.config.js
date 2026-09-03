/** @type {import('next').NextConfig} */
const BUILD_TIMESTAMP = Date.now()

const nextConfig = {
  reactStrictMode: true,
  output: 'standalone',
  generateBuildId: async () => `build-${BUILD_TIMESTAMP}`,
  async headers() {
    return [
      {
        source: '/:path*',
        headers: [
          {
            key: 'Cache-Control',
            value: 'no-cache, no-store, must-revalidate',
          },
        ],
      },
      {
        source: '/_next/static/:path*',
        headers: [
          {
            key: 'Cache-Control',
            value: 'public, max-age=31536000, immutable',
          },
        ],
      },
    ]
  },
  webpack: (config, { isServer, dev }) => {
    if (!isServer && !dev) {
      // Add unique suffix to all chunk filenames to bust browser cache
      config.output.chunkFilename = `static/chunks/[chunkhash]-${BUILD_TIMESTAMP}.js`
      config.output.filename = `static/chunks/[chunkhash]-${BUILD_TIMESTAMP}.js`
    }
    return config
  },
}

module.exports = nextConfig
