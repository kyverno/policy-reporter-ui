import {defineNuxtPlugin} from 'nuxt/app'
import {create} from '~/modules/core/api'
import {useConfigStore} from "~/store/config";
import {type Config, DisplayMode} from "~/modules/core/types";

export default defineNuxtPlugin(async () => {
  const config = useRuntimeConfig()
  const api = create({ baseURL: config.public.coreApi as string, prefix: '' })

  const apiConfig = await api.config().catch((error): Config => {
    console.error(`failed to load config: ${error}`)

    return {
      error,
      displayMode: DisplayMode.UNSPECIFIED,
      default: 'default',
      sources: [],
      plugins: [],
      defaultFilter: { resources: [], clusterResources: [] },
      clusters: [{ name: 'Default', slug: 'default', plugins: [] }],
      oauth: false,
    }
  })

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
