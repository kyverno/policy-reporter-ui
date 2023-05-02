import { NuxtAxiosInstance } from '@nuxtjs/axios'
import { NamespacedStatusCount, StatusCount, Target, Filter, CoreAPI, Result, Config, ResultList, Pagination } from './types'

class API {
  private axios: NuxtAxiosInstance
  private prefix: string = ''

  constructor (axios: NuxtAxiosInstance) {
    this.axios = axios
  }

  public config (): Promise<Config> {
    return this.axios.$get<Config>('/config').catch(() => ({
      plugins: [],
      displayMode: '',
      refreshInterval: 10000,
      labelFilter: [],
      clusters: [],
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
  }

  logs (): Promise<Result[]> {
    return this.axios.$get<Result[]>('/result-log')
  }

  targets (): Promise<Target[]> {
    return this.axios.$get<Target[]>(this.prefix + '/v1/targets')
  }

  namespaces (source?: string): Promise<string[]> {
    return this.axios.$get<string[]>(this.prefix + '/v1/namespaces', { params: { sources: [source] } })
  }

  ruleStatusCount (policy: string, rule: string): Promise<StatusCount[]> {
    return this.axios.$get<StatusCount[]>(this.prefix + '/v1/rule-status-count', { params: { policy, rule } })
  }

  namespacedCategories (source?: string): Promise<string[]> {
    return this.axios.$get<string[]>(this.prefix + '/v1/namespaced-resources/categories', { params: { sources: [source] } })
  }

  namespacedKinds (source?: string): Promise<string[]> {
    return this.axios.$get<string[]>(this.prefix + '/v1/namespaced-resources/kinds', { params: { sources: [source] } })
  }

  namespacedPolicies (source?: string): Promise<string[]> {
    return this.axios.$get<string[]>(this.prefix + '/v1/namespaced-resources/policies', { params: { sources: [source] } })
  }

  namespacedRules (source?: string): Promise<string[]> {
    return this.axios.$get<string[]>(this.prefix + '/v1/namespaced-resources/rules', { params: { sources: [source] } })
  }

  namespacedSources (): Promise<string[]> {
    return this.axios.$get<string[]>(this.prefix + '/v1/namespaced-resources/sources')
  }

  namespacedStatusCount (filter?: Filter): Promise<NamespacedStatusCount[]> {
    return this.axios.$get<NamespacedStatusCount[]>(this.prefix + '/v1/namespaced-resources/status-counts', { params: filter })
  }

  namespacedResults (filter?: Filter, pagination?: Pagination): Promise<ResultList> {
    return this.axios.$get<ResultList>(this.prefix + '/v1/namespaced-resources/results', { params: { ...filter, ...pagination } })
  }

  namespacedReportLabels (source?: string): Promise<{[key: string]: string[]}> {
    return this.axios.$get<{[key: string]: string[]}>(this.prefix + '/v1/namespaced-resources/report-labels', { params: { sources: [source] } })
  }

  clusterCategories (source?: string): Promise<string[]> {
    return this.axios.$get<string[]>(this.prefix + '/v1/cluster-resources/categories', { params: { sources: [source] } })
  }

  clusterKinds (source?: string): Promise<string[]> {
    return this.axios.$get<string[]>(this.prefix + '/v1/cluster-resources/kinds', { params: { sources: [source] } })
  }

  clusterPolicies (source?: string): Promise<string[]> {
    return this.axios.$get<string[]>(this.prefix + '/v1/cluster-resources/policies', { params: { sources: [source] } })
  }

  clusterRules (source?: string): Promise<string[]> {
    return this.axios.$get<string[]>(this.prefix + '/v1/cluster-resources/rules', { params: { sources: [source] } })
  }

  clusterSources (): Promise<string[]> {
    return this.axios.$get<string[]>(this.prefix + '/v1/cluster-resources/sources')
  }

  clusterReportLabels (source?: string): Promise<{[key: string]: string[]}> {
    return this.axios.$get<{[key: string]: string[]}>(this.prefix + '/v1/cluster-resources/report-labels', { params: { sources: [source] } })
  }

  statusCount (filter?: Filter): Promise<StatusCount[]> {
    return this.axios.$get<StatusCount[]>(this.prefix + '/v1/cluster-resources/status-counts', { params: filter })
  }

  results (filter?: Filter, pagination?: Pagination): Promise<ResultList> {
    return this.axios.$get<ResultList>(this.prefix + '/v1/cluster-resources/results', { params: { ...filter, ...pagination } })
  }

  setPrefix (prefix: string): void {
    this.prefix = prefix
  }
}

export const create = (axios: NuxtAxiosInstance): CoreAPI => new API(axios)
