import Vue from 'vue'
import VueHighlightJS from 'vue-highlight.js'
import yaml from 'highlight.js/lib/languages/yaml'

import VueClipboard from 'vue-clipboard2'

import 'highlight.js/styles/github.css'

VueClipboard.config.autoSetContainer = true
Vue.use(VueClipboard)

Vue.use(VueHighlightJS, {
  languages: { yaml }
})
