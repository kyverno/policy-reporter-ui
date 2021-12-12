import { Plugin } from '@nuxt/types'
import qs from 'qs'

const production = process.env.NODE_ENV === 'production'

const plugin: Plugin = ({ $axios }) => {
  $axios.setBaseURL(production ? `//${window.location.host}` : `${process.env.NUXT_ENV_API_URL}`)

  $axios.onRequest((config) => {
    config.paramsSerializer = params => qs.stringify(params, { arrayFormat: 'repeat' })
    return config
  })
}

export default plugin
