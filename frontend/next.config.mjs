import process from 'node:process'
import { PHASE_DEVELOPMENT_SERVER } from 'next/constants.js'

/** @type {import('next').NextConfig} */
const sharedConfig = {
  trailingSlash: true,
  images: {
    unoptimized: true,
  },
}

export default function nextConfig(phase) {
  if (phase === PHASE_DEVELOPMENT_SERVER) {
    const target = (process.env.API_PROXY_TARGET ?? 'http://localhost:8080').replace(/\/$/, '')
    return {
      ...sharedConfig,
      async rewrites() {
        return [{ source: '/api/:path*', destination: `${target}/api/:path*` }]
      },
    }
  }

  return {
    ...sharedConfig,
    output: 'export',
  }
}
