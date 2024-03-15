import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { md3 } from 'vuetify/blueprints'
import { useConfigStore } from "~/store/config";
import 'vuetify/styles'
import {DisplayMode, Severity, Status} from "~/modules/core/types";

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
        [DisplayMode.LIGHT]: {
          colors: {
            'header-item': '#E0E0E0',
            'sub-item': '#EEEEEE',
            header: '#3483c7',
            category: '#CFD8DC',
            primary: '#424242',
            secondary: '#ECEFF1',
            [`status-${Status.SKIP}`]: '#E0E0E0',
            [`status-${Status.PASS}`]: '#43A047',
            [`status-${Status.WARN}`]: '#FB8C00',
            [`status-${Status.FAIL}`]: '#EF5350',
            [`status-${Status.ERROR}`]: '#950011',
            [`btn-${Status.SKIP}`]: '#E0E0E0',
            [`btn-${Status.PASS}`]: '#43A047',
            [`btn-${Status.WARN}`]: '#FB8C00',
            [`btn-${Status.FAIL}`]: '#EF5350',
            [`btn-${Status.ERROR}`]: '#950011',
            [`severity-${Severity.INFO}`]: '#9E9E9E',
            [`severity-${Severity.LOW}`]: '#2194F3',
            [`severity-${Severity.MEDIUM}`]: '#FB8A00',
            [`severity-${Severity.HIGH}`]: '#D32F2F',
            [`severity-${Severity.CRITICAL}`]: '#FF1744',
          }
        },
        [DisplayMode.DARK]: {
          dark: true,
          colors: {
            'header-item': '#424242',
            'sub-item': '#111111',
            header: '#215580',
            category: '#111111',
            primary: '#E0E0E0',
            secondary: '#616161',
            info: '#01579B',
            [`status-${Status.SKIP}`]: '#424242',
            [`status-${Status.PASS}`]: '#1B5E20',
            [`status-${Status.WARN}`]: '#FF6F00',
            [`status-${Status.FAIL}`]: '#D32F2F',
            [`status-${Status.ERROR}`]: '#950011',
            [`btn-${Status.SKIP}`]: '#EEEEEE',
            [`btn-${Status.PASS}`]: '#00E676',
            [`btn-${Status.WARN}`]: '#FF6D00',
            [`btn-${Status.FAIL}`]: '#FF5252',
            [`btn-${Status.ERROR}`]: '#E53935',
            [`severity-${Severity.INFO}`]: '#80D8FF',
            [`severity-${Severity.LOW}`]: '#82B1FF',
            [`severity-${Severity.MEDIUM}`]: '#FFC400',
            [`severity-${Severity.HIGH}`]: '#FF8A80',
            [`severity-${Severity.CRITICAL}`]: '#FF5252',
          }
        },
        [DisplayMode.COLOR_BLIND]: {
          dark: false,
          colors: {
            'header-item': '#E0E0E0',
            'sub-item': '#EEEEEE',
            header: '#3483c7',
            category: '#CFD8DC',
            primary: '#424242',
            secondary: '#ECEFF1',
            [`status-${Status.SKIP}`]: '#9E9E9E',
            [`status-${Status.PASS}`]: '#039BE5',
            [`status-${Status.WARN}`]: '#FFD54F',
            [`status-${Status.FAIL}`]: '#E64A19',
            [`status-${Status.ERROR}`]: '#4E342E',
            [`btn-${Status.SKIP}`]: '#9E9E9E',
            [`btn-${Status.PASS}`]: '#039BE5',
            [`btn-${Status.WARN}`]: '#FFD54F',
            [`btn-${Status.FAIL}`]: '#E64A19',
            [`btn-${Status.ERROR}`]: '#4E342E',
            [`severity-${Severity.INFO}`]: '#9E9E9E',
            [`severity-${Severity.LOW}`]: '#2194F3',
            [`severity-${Severity.MEDIUM}`]: '#FB8A00',
            [`severity-${Severity.HIGH}`]: '#BF360C',
            [`severity-${Severity.CRITICAL}`]: '#6D4C41',
          }
        },
        [DisplayMode.COLOR_BLIND_DARK]: {
          dark: true,
          colors: {
            'header-item': '#424242',
            'sub-item': '#111111',
            header: '#215580',
            category: '#111111',
            primary: '#E0E0E0',
            secondary: '#616161',
            info: '#01579B',
            [`status-${Status.SKIP}`]: '#9E9E9E',
            [`status-${Status.PASS}`]: '#01579B',
            [`status-${Status.WARN}`]: '#FF6F00',
            [`status-${Status.FAIL}`]: '#BF360C',
            [`status-${Status.ERROR}`]: '#5D4037',
            [`btn-${Status.SKIP}`]: '#9E9E9E',
            [`btn-${Status.PASS}`]: '#039BE5',
            [`btn-${Status.WARN}`]: '#FF6F00',
            [`btn-${Status.FAIL}`]: '#BF360C',
            [`btn-${Status.ERROR}`]: '#BCAAA4',
            [`severity-${Severity.INFO}`]: '#9E9E9E',
            [`severity-${Severity.LOW}`]: '#2194F3',
            [`severity-${Severity.MEDIUM}`]: '#FB8A00',
            [`severity-${Severity.HIGH}`]: '#BF360C',
            [`severity-${Severity.CRITICAL}`]: '#BCAAA4',
          }
        },
      }
    }
  })
  config.$subscribe((mutation, state) => {
    if (vuetify.theme.global.name.value === state.displayMode) return;

    vuetify.theme.global.name.value = state.displayMode
  })

  nuxtApp.vueApp.use(vuetify)
})
