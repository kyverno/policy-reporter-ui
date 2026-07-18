import {
  type Target,
  type Filter,
  type Config,
  type ResultList,
  type Pagination,
  type FindingCounts,
  type ResourceResultList,
  type ResourceResult,
  type ResourceStatusCount,
  type Source,
  type PolicyResult,
  type NamespaceStatusCount,
  Status,
  type Profile,
  type LayoutConfig,
  type Dashboard,
  type ResourceDetails,
  type SourceDetails,
  type PolicyFilter,
  type PolicyDetails,
  type ExceptionResponse,
  type ExceptionPolicy,
} from '~/types/core'

import type { NitroFetchOptions, NitroFetchRequest } from "nitropack";

type APIConfig = { baseURL: string; prefix?: string; };

export const cluster = useSessionStorage('cluster', 'default')

export class CoreAPI {
  private readonly baseURL: string
  private cluster: string = ''

  private nsExcludes: string[] = []
  private clusterExcludes: string[] = []

  constructor (config: APIConfig) {
    this.baseURL = config.baseURL
    this.cluster = config.prefix || ''
  }

  profile () {
    return exec<Profile>('/profile', { baseURL: this.baseURL })
  }

  layout () {
    return exec<LayoutConfig>(`/api/config/${this.cluster}/layout`, { baseURL: this.baseURL })
  }

  config () {
    return exec<Config>('/api/config', { baseURL: this.baseURL })
  }

  dashboard (filter?: Filter) {
    return exec<Dashboard>(`/api/${this.cluster}/dashboard`, { baseURL: this.baseURL, params: applyExcludes(filter, [...this.nsExcludes, ...this.clusterExcludes]) })
  }

  customBoard (id: string, filter?: Filter) {
    return exec<Dashboard>(`/api/${this.cluster}/custom-board/${id}`, { baseURL: this.baseURL, params: applyExcludes(filter, [...this.nsExcludes, ...this.clusterExcludes]) })
  }

  resource (id: string, filter?: Filter) {
    return exec<ResourceDetails>(`/api/${this.cluster}/resource/${id}`, { baseURL: this.baseURL, params: filter })
  }

  createException (id: string, source: string, policies?: ExceptionPolicy[], category?: string) {
    return exec<ExceptionResponse>(`/api/${this.cluster}/resource/${id}/exception`, { baseURL: this.baseURL, method: "POST", body: { policies, source, category } })
  }

  policySources (filter?: Filter) {
    return exec<{ filter: PolicyFilter; sources: SourceDetails[] }>(`/api/${this.cluster}/policy-sources`, { baseURL: this.baseURL, params: applyExcludes(filter, [...this.nsExcludes, ...this.clusterExcludes]) })
  }

  policyDetails (source: string, policy: string, namespace?: string, status?: Status[], kinds?: string[] ) {
    return exec<PolicyDetails>(`/api/${this.cluster}/${source}/policy/details`, { baseURL: this.baseURL, params: applyExcludes({ policies: [policy], namespace, status, kinds }, [...this.nsExcludes, ...this.clusterExcludes], source) })
  }

  policyHTMLReport (source: string, filter: { namespaces: string[]; categories: string[]; kinds: string[]; clusterScope: boolean; status: Status[]; }) {
    return exec<BlobPart>(`/api/${this.cluster}/${source}/policy-report`, { baseURL: this.baseURL, params: filter, responseType: 'blob' })
  }

  namespaceHTMLReport (source: string, filter: { namespaces: string[]; categories: string[]; kinds: string[]; status: Status[]; }) {
    return exec<BlobPart>(`/api/${this.cluster}/${source}/namespace-report`, { baseURL: this.baseURL, params: filter, responseType: 'blob' })
  }

  policies (source: string, filter?: Filter) {
    return exec<{ [category: string]: PolicyResult[] }>(`/api/${this.cluster}/${source}/policies`, { baseURL: this.baseURL, params: applyExcludes(filter, this.nsExcludes)})
  }

  targets () {
    return exec<{ [type: string]: Target[] }>(`/api/${this.cluster}/targets`, { baseURL: this.baseURL })
  }

  namespaces (filter?: Filter) {
    return exec<string[]>(`/api/${this.cluster}/namespaces`, { baseURL: this.baseURL, params: { ...filter } })
  }

  namespace (filter?: Filter) {
    return exec<Dashboard>(`/api/${this.cluster}/namespace`, { baseURL: this.baseURL, params: { ...filter } })
  }

  namespacedResults (filter?: Filter, pagination?: Pagination) {
    return exec<ResultList>('/api/'+this.cluster+'/namespace-scoped/results', { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes), ...pagination } })
  }

  clusterResults (filter?: Filter, pagination?: Pagination) {
    return exec<ResultList>('/api/'+this.cluster+'/cluster-scoped/results', { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes), ...pagination } })
  }

  namespacedResourceResults (filter?: Filter, pagination?: Pagination) {
    return exec<ResourceResultList>('/api/'+this.cluster+'/namespace-scoped/resource-results', { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes), ...pagination } })
  }

  clusterResourceResults (filter?: Filter, pagination?: Pagination) {
    return exec<ResourceResultList>('/api/'+this.cluster+'/cluster-scoped/resource-results', { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.clusterExcludes), ...pagination } })
  }

  resourceResults (id: string, filter?: Filter) {
    return exec<ResourceResult[]>(`/api/${this.cluster}/resource/${id}/resource-results`, { baseURL: this.baseURL, params: filter })
  }

  results (id: string, pagination?: Pagination, filter?: Filter) {
    return exec<ResultList>(`/api/${this.cluster}/resource/${id}/results`, { baseURL: this.baseURL, params: { ...pagination, ...filter } })
  }

  resultsWithoutResources (filter?: Filter, pagination?: Pagination) {
    return exec<ResultList>('/api/'+this.cluster+'/results-without-resource', { baseURL: this.baseURL, params: { ...filter, ...pagination } })
  }

  clustersDashboard (filter?: Filter) {
    return exec<Dashboard>(`/api/clusters`, { baseURL: this.baseURL, params: applyExcludes(filter, [...this.nsExcludes, ...this.clusterExcludes]) })
  }

  totalResults (cluster?: string, pagination?: Pagination, filter?: Filter) {
    return exec<ResourceResultList>(`/api/${cluster ?? this.cluster}/total-results`, { baseURL: this.baseURL, params: { ...applyExcludes(filter, [...this.nsExcludes, ...this.clusterExcludes]), ...pagination } })
  }

  customBoardResourceResults (id: string, filter?: Filter, pagination?: Pagination) {
    return exec<ResourceResultList>(`/api/${this.cluster}/custom-board/${id}/resource-results`, { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes), ...pagination } })
  }

  customBoardResults (id: string, filter?: Filter, pagination?: Pagination) {
    return exec<ResultList>(`/api/${this.cluster}/custom-board/${id}/results`, { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes), ...pagination } })
  }

  customBoardClusterResourceResults (id: string, filter?: Filter, pagination?: Pagination) {
    return exec<ResourceResultList>(`/api/${this.cluster}/custom-board/${id}/cluster-resource-results`, { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes), ...pagination } })
  }

  customBoardClusterResults (id: string, filter?: Filter, pagination?: Pagination) {
    return exec<ResultList>(`/api/${this.cluster}/custom-board/${id}/cluster-results`, { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes), ...pagination } })
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


const exec = <T>(api: string, opts: NitroFetchOptions<NitroFetchRequest>): Promise<T> => {
  return $fetch<T>(api, opts)
}

const applyExcludes = <T extends Filter>(filter: T | undefined, exclude: string[] | undefined, source?: string) => {
  if (!filter) return ({ exclude })

  if (filter.kinds && filter.kinds.length > 0) return filter

  if (!exclude || exclude.length < 1) return filter

  if (filter.sources && filter.sources.length > 0) {
    exclude = exclude.filter((e) => filter.sources?.some(s => e.startsWith(s)))
  }

  if (source) {
    exclude = exclude.filter((e) => e.startsWith(source))
  }

  return {
    ...filter,
    exclude
  }
}
