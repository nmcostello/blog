import { defineConfig } from 'astro/config';

export default defineConfig({
  site: 'https://nmcostello.github.io/blog',
  base: '/blog/',
  outDir: 'docs',
  output: 'static'
});
