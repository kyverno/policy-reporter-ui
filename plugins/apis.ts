import { Plugin } from '@nuxt/types'
import { create as createCoreAPI } from '~/policy-reporter-plugins/core/api'
import { create as createKyvernoAPI } from '~/policy-reporter-plugins/kyverno/api'
import { CoreAPI } from '~/policy-reporter-plugins/core/types'
import { KyvernoAPI } from '~/policy-reporter-plugins/kyverno/types'

declare module 'vue/types/vue' {
    interface Vue {
        $coreAPI: CoreAPI
        $kyvernoAPI: KyvernoAPI
    }
}

declare module '@nuxt/types' {
    interface NuxtAppOptions {
        $coreAPI: CoreAPI
        $kyvernoAPI: KyvernoAPI
    }
    interface Context {
        $coreAPI: CoreAPI
        $kyvernoAPI: KyvernoAPI
    }
}

const plugin: Plugin = ({ app }, inject) => {
  inject('coreAPI', createCoreAPI(app.$axios))
  inject('kyvernoAPI', createKyvernoAPI(app.$axios))
}

export default plugin
