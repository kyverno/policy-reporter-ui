import { Status, type Dictionary } from './types'


export const sortByKeys = (dic: Dictionary): Dictionary => Object.keys(dic).sort().reduce<Dictionary>((obj, key) => {
  obj[key] = dic[key]
  return obj
}, {})

export const capilize = (source: string) => source.charAt(0).toUpperCase() + source.slice(1)
