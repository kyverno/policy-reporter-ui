import { defineNuxtPlugin } from 'nuxt/app'
import { create } from '~/modules/core/api'
import { useConfigStore } from "~/store/config";
import { clusterKinds, kinds } from "~/modules/core/store/filter";

export default defineNuxtPlugin(async () => {
  const config = useRuntimeConfig()
  const api = create({ baseURL: config.public.coreApi as string, prefix: '' })

  const apiConfig = await api.config()

  api.setPrefix(apiConfig.default)

  kinds.value = apiConfig.defaultFilter.resources || []
  clusterKinds.value = apiConfig.defaultFilter.clusterResources || []

  const store = useConfigStore()

  store.setConfig(apiConfig)

  return {
    provide: {
      coreAPI: api
    }
  }
})
