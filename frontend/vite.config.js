import { defineConfig } from 'vite'
import mkcert from 'vite-plugin-mkcert'
import { svelte } from '@sveltejs/vite-plugin-svelte'

export default defineConfig({
  server: { https: true },
  plugins: [svelte(), mkcert()]
})