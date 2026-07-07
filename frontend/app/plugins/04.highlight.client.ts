import hljs from 'highlight.js/lib/core'
import javascript from 'highlight.js/lib/languages/javascript'
import yaml from 'highlight.js/lib/languages/yaml'
import typescript from 'highlight.js/lib/languages/typescript'
import highlightJS from '@highlightjs/vue-plugin'
import 'highlight.js/styles/atom-one-dark.css'

export default defineNuxtPlugin((nuxtApp) => {
  hljs.registerLanguage('javascript', javascript)
  hljs.registerLanguage('typescript', typescript)
  hljs.registerLanguage('yaml', yaml)
  nuxtApp.vueApp.use(highlightJS)
})
