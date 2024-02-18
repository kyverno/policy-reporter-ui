import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { md3 } from 'vuetify/blueprints'
import { useConfigStore } from "~/store/config";
import 'vuetify/styles'

export default defineNuxtPlugin((nuxtApp) => {
  const config = useConfigStore()
  const vuetify = createVuetify({
    ssr: true,
    components: {
      ...components
    },
    directives,
    blueprint: md3,
    theme: {
      defaultTheme: config.theme,
      themes: {
        light: {
          colors: {
            'header-item': '#E0E0E0',
            'sub-item': '#EEEEEE',
            header: '#3483c7',
            category: '#CFD8DC',
            primary: '#424242',
            secondary: '#ECEFF1',
          }
        },
        dark: {
          colors: {
            'header-item': '#424242',
            'sub-item': '#111111',
            header: '#215580',
            category: '#111111',
            primary: '#E0E0E0',
            secondary: '#616161',
            info: '#01579B'
          }
        }
      }
    }
  })

  config.$subscribe((mutation, state) => {
    if (vuetify.theme.global.name.value === state.displayMode) return;

    vuetify.theme.global.name.value = state.displayMode
  })

  nuxtApp.vueApp.use(vuetify)
})
