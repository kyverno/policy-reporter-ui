import { useFetch } from 'nuxt/app'
import {
  type NamespacedStatusCount,
  type StatusCount,
  type Target,
  type Filter,
  type Result,
  type Config,
  type ResultList,
  type Pagination,
  type FindingCounts,
  type ResourceResultList,
  type ResourceResult,
  type ResourceStatusCount,
  type Resource,
  type Source,
  type CustomBoard,
  type CustomBoardDetails,
  type PolicyResult, type NamespaceStatusCount
} from './types'

type APIConfig = { baseURL: string; prefix?: string; };

export const cluster = ref('default')

export class CoreAPI {
  private baseURL: string
  private cluster: string = ''

  private nsExcludes: string[] = []
  private clusterExcludes: string[] = []

  constructor (config: APIConfig) {
    this.baseURL = config.baseURL
    this.cluster = config.prefix || ''
  }

  config () {
    return $fetch<Config>('/api/config', { baseURL: this.baseURL })
  }

  logs () {
    return $fetch<Result[]>('/result-log', { baseURL: this.baseURL })
  }

  targets () {
    return $fetch<Target[]>('/proxy/'+this.cluster+'/core/v1/targets', { baseURL: this.baseURL })
  }

  categories (source?: string, kind?: string) {
    return $fetch<string[]>('/proxy/'+this.cluster+'/core/v1/categories', { baseURL: this.baseURL, params: { sources: [source], kinds: kind ? [kind] : undefined } })
  }

  namespaces (filter?: Filter) {
    return $fetch<string[]>('/proxy/'+this.cluster+'/core/v1/namespaces', { baseURL: this.baseURL, params: { ...filter } })
  }

  ruleStatusCount (policy: string, rule: string) {
    return useFetch<StatusCount[]>('/proxy/'+this.cluster+'/core/v1/rule-status-count', { baseURL: this.baseURL, params: { policy, rule } })
  }

  namespacedKinds (source?: string) {
    return $fetch<string[]>('/proxy/'+this.cluster+'/core/v1/namespaced-resources/kinds', { baseURL: this.baseURL, params: { sources: source ? [source] : undefined } })
  }

  namespacedPolicies (source?: string) {
    return $fetch<string[]>('/proxy/'+this.cluster+'/core/v1/namespaced-resources/policies', { baseURL: this.baseURL, params: { sources: source ? [source] : undefined } })
  }

  namespacedRules (source?: string) {
    return $fetch<string[]>('/proxy/'+this.cluster+'/core/v1/namespaced-resources/rules', { baseURL: this.baseURL, params: { sources: source ? [source] : undefined } })
  }

  namespacedSources () {
    return $fetch<string[]>('/proxy/'+this.cluster+'/core/v1/namespaced-resources/sources', { baseURL: this.baseURL })
  }

  namespacedStatusCount (filter?: Filter) {
    return $fetch<NamespaceStatusCount>('/proxy/'+this.cluster+'/core/v2/namespace-scoped/status-counts', { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes) } })
  }

  namespacedResults (filter?: Filter, pagination?: Pagination) {
    return $fetch<ResultList>('/proxy/'+this.cluster+'/core/v1/namespaced-resources/results', { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes), ...pagination } })
  }

  clusterKinds (source?: string) {
    return $fetch<string[]>('/proxy/'+this.cluster+'/core/v1/cluster-resources/kinds', { baseURL: this.baseURL, params: { sources: source ? [source] : undefined } })
  }

  clusterPolicies (source?: string) {
    return $fetch<string[]>('/proxy/'+this.cluster+'/core/v1/cluster-resources/policies', { baseURL: this.baseURL, params: { sources: source ? [source] : undefined } })
  }

  clusterRules (source?: string) {
    return $fetch<string[]>('/proxy/'+this.cluster+'/core/v1/cluster-resources/rules', { baseURL: this.baseURL, params: { sources: source ? [source] : undefined } })
  }

  clusterSources () {
    return $fetch<string[]>('/proxy/'+this.cluster+'/core/v1/cluster-resources/sources', { baseURL: this.baseURL })
  }

  statusCount (filter?: Filter) {
    return $fetch<StatusCount[]>('/proxy/'+this.cluster+'/core/v1/cluster-resources/status-counts', { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.clusterExcludes) } })
  }

  clusterResults (filter?: Filter, pagination?: Pagination) {
    return $fetch<ResultList>('/proxy/'+this.cluster+'/core/v1/cluster-resources/results', { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes), ...pagination } })
  }

  countFindings (filter?: Filter) {
    return $fetch<FindingCounts>('/proxy/'+this.cluster+'/core/v1/finding-counts', { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes) } })
  }

  namespacedResourceResults (filter?: Filter, pagination?: Pagination) {
    return $fetch<ResourceResultList>('/proxy/'+this.cluster+'/core/v1/namespaced-resources/resource-results', { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes), ...pagination } })
  }

  clusterResourceResults (filter?: Filter, pagination?: Pagination) {
    return $fetch<ResourceResultList>('/proxy/'+this.cluster+'/core/v1/cluster-resources/resource-results', { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.clusterExcludes), ...pagination } })
  }

  resourceResults (id: string, filter?: Filter) {
    return $fetch<ResourceResult[]>('/proxy/'+this.cluster+'/core/v1/resource-results', { baseURL: this.baseURL, params: { id, ...filter } })
  }

  resourceStatusCount (id: string, filter?: Filter) {
    return $fetch<ResourceStatusCount[]>('/proxy/'+this.cluster+'/core/v1/resource-status-counts', { baseURL: this.baseURL, params: { id, ...filter }})
  }

  resource (id: string, filter?: Filter) {
    return $fetch<Resource>('/proxy/'+this.cluster+'/core/v1/resource', { baseURL: this.baseURL, params: { id, ...filter }})
  }

  results (id: string, pagination?: Pagination, filter?: Filter) {
    return $fetch<ResultList>('/proxy/'+this.cluster+'/core/v1/results', { baseURL: this.baseURL, params: { id, ...pagination, ...filter } })
  }

  sources (id?: string, filter?: Filter) {
    return $fetch<Source[]>('/proxy/'+this.cluster+'/core/v1/sources', { baseURL: this.baseURL, params: { id, ...applyExcludes(filter, [...this.nsExcludes, ...this.clusterExcludes]) } })
  }

  policies (filter?: Filter) {
    return $fetch<PolicyResult[]>('/proxy/'+this.cluster+'/core/v1/policies', { baseURL: this.baseURL, params: applyExcludes(filter, this.nsExcludes)})
  }

  customBoards () {
    return $fetch<CustomBoard[]>('/api/custom-board/list', { baseURL: this.baseURL })
  }

  customBoard (id: string) {
    return $fetch<CustomBoardDetails>(`/api/custom-board/${this.cluster}/${id}`, { baseURL: this.baseURL })
  }

  setPrefix (prefix: string): void {
    this.cluster = prefix
    cluster.value = prefix
  }

  setExcludes (nsFilter: string[], clusterFilter: string[]): void {
    this.nsExcludes = nsFilter
    this.clusterExcludes = clusterFilter
  }
}

export const create = (config: APIConfig): CoreAPI => new CoreAPI(config)

const applyExcludes = (filter: Filter | undefined, exclude: string[] | undefined) => {
  if (!filter) return ({ exclude })

  if (filter.kinds && filter.kinds.length > 0) return filter

  if (!exclude || exclude.length < 1) return filter

  if (filter.sources && filter.sources.length > 0) {
    exclude = exclude.filter((e) => filter.sources?.some(s => e.startsWith(s)))
  }

  return {
    ...filter,
    exclude
  }
}
