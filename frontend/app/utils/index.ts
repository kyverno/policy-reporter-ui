import { Status, type Dictionary } from '~/types/core'


export const sortByKeys = (dic: Dictionary): Dictionary => Object.keys(dic).sort().reduce<Dictionary>((obj, key) => {
  obj[key] = dic[key] || ''
  return obj
}, {})

export const capilize = (source: string) => source.charAt(0).toUpperCase() + source.slice(1)

export const execOnChange = <T>(n: T, o: T, cb: () => any ) => {
  if (JSON.stringify(n) === JSON.stringify(o)) { return; }

  cb()
}

export const onChange = <T>(cb: () => any) => (n: T, o: T) =>  {
  if (JSON.stringify(n) === JSON.stringify(o)) { return; }

  cb()
}
