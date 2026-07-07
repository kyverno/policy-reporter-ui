import vuetify, { transformAssetUrls } from 'vite-plugin-vuetify'
import { join } from 'node:path'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr: false,
  router: {
    options: {
      hashMode: true
    }
  },
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  css: [
    'vuetify/styles',
    '@mdi/font/css/materialdesignicons.min.css',
  ],
  build: { transpile: ["vuetify"] },

  modules: [
    '@pinia/nuxt',
    '@vueuse/nuxt',
  ],

  runtimeConfig: {
    public: {
      coreApi: '',
    }
  },

  vite: {
    plugins: [
      vuetify({ autoImport: true }),
    ],
    vue: {
      template: {
        transformAssetUrls,
      },
    },
  },

  nitro: {
    output: {
      publicDir: join(__dirname, 'dist')
    }
  },

  app: {
    baseURL: process.env.NODE_ENV === 'production' ? '.' : '/',
    head: {
      title: 'Policy Reporter UI',
      htmlAttrs: {
        lang: 'en'
      },
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'format-detection', content: 'telephone=no' }
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: './favicon.ico' }
      ]
    }
  },

})
