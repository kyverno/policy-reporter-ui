import { NuxtAxiosInstance } from '@nuxtjs/axios'
import { ListResult, NamespacedStatusCount, StatusCount, Target, Filter, CoreAPI, Result, Config } from './types'

export const create = (axios: NuxtAxiosInstance): CoreAPI => ({
  config: (): Promise<Config> => {
    return axios.$get<Config>('/api/config').catch(() => ({
      plugins: [],
      displayMode: '',
      views: {
        logs: true,
        policyReports: true,
        clusterPolicyReports: true,
        kyvernoPolicies: true,
        kyvernoVerifyImages: true,
        dashboard: {
          policyReports: true,
          clusterPolicyReports: true
        }
      }
    }))
  },
  logs: (): Promise<Result[]> => {
    return axios.$get<Result[]>('/api/result-log')
  },
  targets: (): Promise<Target[]> => {
    return axios.$get<Target[]>('/api/v1/targets')
  },
  categories: (source?: string): Promise<string[]> => {
    return axios.$get<string[]>('/api/v1/categories', { params: { sources: [source] } })
  },
  namespaces: (source?: string): Promise<string[]> => {
    return axios.$get<string[]>('/api/v1/namespaces', { params: { sources: [source] } })
  },
  ruleStatusCount: (policy: string, rule: string): Promise<StatusCount[]> => {
    return axios.$get<StatusCount[]>('/api/v1/rule-status-count', { params: { policy, rule } })
  },
  namespacedKinds: (source?: string): Promise<string[]> => {
    return axios.$get<string[]>('/api/v1/namespaced-resources/kinds', { params: { sources: [source] } })
  },
  namespacedPolicies: (source?: string): Promise<string[]> => {
    return axios.$get<string[]>('/api/v1/namespaced-resources/policies', { params: { sources: [source] } })
  },
  namespacedSources: (): Promise<string[]> => {
    return axios.$get<string[]>('/api/v1/namespaced-resources/sources')
  },
  namespacedStatusCount: (filter?: Filter): Promise<NamespacedStatusCount[]> => {
    return axios.$get<NamespacedStatusCount[]>('/api/v1/namespaced-resources/status-counts', { params: filter })
  },
  namespacedResults: (filter?: Filter): Promise<ListResult[]> => {
    return axios.$get<ListResult[]>('/api/v1/namespaced-resources/results', { params: filter })
  },
  clusterKinds: (source?: string): Promise<string[]> => {
    return axios.$get<string[]>('/api/v1/cluster-resources/kinds', { params: { sources: [source] } })
  },
  clusterPolicies: (source?: string): Promise<string[]> => {
    return axios.$get<string[]>('/api/v1/cluster-resources/policies', { params: { sources: [source] } })
  },
  clusterSources: (): Promise<string[]> => {
    return axios.$get<string[]>('/api/v1/cluster-resources/sources')
  },
  statusCount: (filter?: Filter): Promise<StatusCount[]> => {
    return axios.$get<StatusCount[]>('/api/v1/cluster-resources/status-counts', { params: filter })
  },
  results: (filter?: Filter): Promise<ListResult[]> => {
    return axios.$get<ListResult[]>('/api/v1/cluster-resources/results', { params: filter })
  }
})
