import { Plugin } from '@nuxt/types'
import { AxiosResponse } from 'axios'
import qs from 'qs'

const production = process.env.NODE_ENV === 'production'

function trimSlashes (str: string) {
  return str.split('/').filter(p => !!p).join('/')
}

const plugin: Plugin = ({ $axios }) => {
  $axios.setBaseURL(production ? `//${window.location.host}/${trimSlashes(window.location.pathname)}` : `${process.env.NUXT_ENV_API_URL}`)
  $axios.interceptors.response.use((response: AxiosResponse) => {
    if (response.headers['content-type'] === 'text/html; charset=utf-8') {
      return Promise.reject(new Error('unexpected content-type'))
    }

    return Promise.resolve(response)
  })

  $axios.onRequest((config) => {
    config.paramsSerializer = params => qs.stringify(params, { arrayFormat: 'repeat' })
    return config
  })
}

export default plugin
