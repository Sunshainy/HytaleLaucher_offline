import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '@wailsjs': resolve(__dirname, 'wailsjs')
    }
  },
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
    // Generate source maps for debugging
    sourcemap: false,
    // Use content-based hashing for cache busting - matches original binary
    rollupOptions: {
      output: {
        // Entry files get content hash
        entryFileNames: 'assets/[name].[hash].js',
        // Chunk files get content hash with component names
        chunkFileNames: (chunkInfo) => {
          const name = chunkInfo.name || 'chunk'
          // Map component names to match original bundle structure
          return `assets/${name}.[hash].js`
        },
        // Asset files get content hash
        assetFileNames: (assetInfo) => {
          const name = assetInfo.name || 'asset'
          const ext = name.split('.').pop()
          // Extract base name without extension
          const baseName = name.replace(/\.[^.]+$/, '')
          return `assets/${baseName}.[hash].${ext}`
        },
        // Manual chunk configuration to match original binary structure
        manualChunks: (id) => {
          // Keep Vue runtime in main bundle
          if (id.includes('node_modules')) {
            if (id.includes('vue-router')) return undefined
            if (id.includes('pinia')) return undefined
            if (id.includes('vue-i18n')) return undefined
            if (id.includes('vue')) return undefined
          }
          // View components get their own chunks
          if (id.includes('/views/')) {
            const match = id.match(/\/views\/([^/]+)\.vue/)
            if (match) {
              return match[1]
            }
          }
          // Component chunks
          if (id.includes('/components/')) {
            const match = id.match(/\/components\/([^/]+)\.vue/)
            if (match) {
              return `${match[1]}Component`
            }
          }
          // Store chunks
          if (id.includes('/stores/')) {
            const match = id.match(/\/stores\/([^/]+)\.ts/)
            if (match) {
              return match[1]
            }
          }
          return undefined
        }
      }
    },
    // CSS code splitting to match original
    cssCodeSplit: true,
    // Minification settings
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: false,
        drop_debugger: true
      },
      format: {
        comments: false
      }
    }
  },
  // CSS configuration
  css: {
    devSourcemap: false
  }
})
