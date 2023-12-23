import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { md3 } from 'vuetify/blueprints'


export default defineNuxtPlugin((nuxtApp) => {
  const vuetify = createVuetify({
    ssr: true,
    components: {
      ...components
    },
    directives,
    blueprint: md3,
    theme: {
      themes: {
        light: {
          colors: {
            primary: '#01579B',
            secondary: '#FFE0B2',
          }
        }
      }
    }
  })

  nuxtApp.vueApp.use(vuetify)
})
