import {defineNuxtPlugin} from 'nuxt/app'
import {cluster, create} from '~/core/api'
import {useConfigStore} from "~/store/config";
import {type Config, DisplayMode} from "~/types/core";

const trimSlashes = (str: string) => str.split('/').filter(p => !!p).join('/')

export default defineNuxtPlugin(async () => {
  const config = useRuntimeConfig()
  const api = create({ baseURL: config.public.coreApi || `//${window.location.host}/${trimSlashes(window.location.pathname)}`, prefix: cluster.value })

  const apiConfig = await api.config().catch((error: any): Config => {
    console.error(`failed to load config: ${error}`)

    return {
      error,
      displayMode: DisplayMode.UNSPECIFIED,
      default: 'default',
      sources: [],
      plugins: [],
      clusters: [{ name: 'Default', slug: 'default', plugins: [] }],
      oauth: false,
      banner: '',
      logo: {
        disabled: false,
      }
    }
  })

  if (!apiConfig.clusters.some(c => c.slug === cluster.value)) {
    api.setPrefix(apiConfig.default)
  }

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
