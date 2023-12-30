import { defineNuxtPlugin } from 'nuxt/app'
import { create } from '~/modules/core/api'
import { useConfigStore } from "~/store/config";
import { clusterKinds, kinds } from "~/modules/core/store/filter";

export default defineNuxtPlugin(async () => {
  const config = useRuntimeConfig()
  const api = create({ baseURL: config.public.coreApi as string, prefix: '' })

  const apiConfig = await api.config()

  if (apiConfig.oauth) {
    api.profile().then(console.log)
  }

  api.setPrefix(apiConfig.default)

  api.setExcludes((apiConfig.sources || []).reduce<string[]>((acc, config) => {
    return [...acc, ...(config.excludes.namespaceKinds || []).map(k => `${config.name}:${k}`)]
  }, []),
    (apiConfig.sources || []).reduce<string[]>((acc, config) => {
    return [...acc, ...(config.excludes.clusterKinds || []).map(k => `${config.name}:${k}`)]
  }, []))

  const store = useConfigStore()

  store.setConfig(apiConfig)
  store.setDisplayMode(apiConfig.displayMode)

  return {
    provide: {
      coreAPI: api
    }
  }
})
