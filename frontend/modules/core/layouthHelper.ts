import { Status, type Dictionary } from './types'

type StatusCounters = { [status in Status]: { namespaces: string[]; counts: number[] } }

export const shortGraph = (status: Status, counters: StatusCounters): boolean => {
  const pass = counters[Status.PASS].namespaces.length ? 1 : 0
  const fail = counters[Status.FAIL].namespaces.length ? 1 : 0
  const warn = counters[Status.WARN].namespaces.length ? 1 : 0
  const error = counters[Status.ERROR].namespaces.length ? 1 : 0
  const skip = counters[Status.SKIP].namespaces.length ? 1 : 0

  switch (status) {
    case Status.FAIL:
      return (pass + warn + error + skip) > 0
    case Status.PASS:
      return fail === 1 || (warn + error + skip) > 0
    case Status.ERROR:
      return (pass + fail) === 1 || (warn + skip) > 0
    case Status.WARN:
      return (pass + fail + error) % 2 === 1 || skip > 0
    case Status.SKIP:
      return (pass + fail + error + warn) % 2 === 1
  }
}

type ClusterStatusCounters = { [status in Status]: number }

export const boxSizes = (counters: ClusterStatusCounters): { sm: number, md: number, col: number} => {
  const pass = counters[Status.PASS] > 0 ? 1 : 0
  const fail = counters[Status.FAIL] > 0 ? 1 : 0
  const warn = counters[Status.WARN] > 0 ? 1 : 0
  const error = counters[Status.ERROR] > 0 ? 1 : 0
  const skip = counters[Status.SKIP] > 0 ? 1 : 0

  const amount = pass + fail + warn + error + skip

  return {
    md: Math.floor(12 / amount),
    sm: amount === 1 ? 12 : 6,
    col: 12
  }
}

export const sortByKeys = (dic: Dictionary): Dictionary => Object.keys(dic).sort().reduce<Dictionary>((obj, key) => {
  obj[key] = dic[key]
  return obj
}, {})

export const capilize = (source: string) => source.charAt(0).toUpperCase() + source.slice(1)
