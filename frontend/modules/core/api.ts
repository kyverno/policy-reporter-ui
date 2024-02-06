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
  type ResourceDetails, type SourceDetails, type PolicyDetails
} from './types'

type APIConfig = { baseURL: string; prefix?: string; };

export const cluster = ref('default')

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
    return $fetch<Profile>('/profile', { baseURL: this.baseURL })
  }

  layout () {
    return $fetch<LayoutConfig>(`/api/config/${this.cluster}/layout`, { baseURL: this.baseURL })
  }

  dashboard <T extends Boolean>(filter?: Filter) {
    return $fetch<Dashboard<T>>(`/api/config/${this.cluster}/dashboard`, { baseURL: this.baseURL, params: applyExcludes(filter, [...this.nsExcludes, ...this.clusterExcludes]) })
  }

  customBoard <T extends Boolean>(id: string, filter?: Filter) {
    return $fetch<Dashboard<T>>(`/api/config/${this.cluster}/custom-board/${id}`, { baseURL: this.baseURL, params: applyExcludes(filter, [...this.nsExcludes, ...this.clusterExcludes]) })
  }

  resource (id: string, filter?: Filter) {
    return $fetch<ResourceDetails>(`/api/config/${this.cluster}/resource/${id}`, { baseURL: this.baseURL, params: filter })
  }

  policySources (filter?: Filter) {
    return $fetch<SourceDetails[]>(`/api/config/${this.cluster}/policy-sources`, { baseURL: this.baseURL, params: applyExcludes(filter, [...this.nsExcludes, ...this.clusterExcludes]) })
  }

  policyDetails (source: string, policy: string, namespace?: string) {
    return $fetch<PolicyDetails>(`/api/config/${this.cluster}/${source}/policy/details`, { baseURL: this.baseURL, params: applyExcludes({ policies: [policy], namespace }, [...this.nsExcludes, ...this.clusterExcludes]) })
  }

  policyHTMLReport (source: string, filter: { namespaces: string[]; categories: string[]; kinds: string[]; clusterScope: boolean; }) {
    return $fetch<BlobPart>(`/api/config/${this.cluster}/${source}/policy-report`, { baseURL: this.baseURL, params: filter, responseType: 'blob' })
  }

  namespaceHTMLReport (source: string, filter: { namespaces: string[]; categories: string[]; kinds: string[]; }) {
    return $fetch<BlobPart>(`/api/config/${this.cluster}/${source}/namespace-report`, { baseURL: this.baseURL, params: filter, responseType: 'blob' })
  }

  policies (source: string, filter?: Filter) {
    return $fetch<{ [category: string]: PolicyResult[] }>(`/api/config/${this.cluster}/${source}/policies`, { baseURL: this.baseURL, params: applyExcludes(filter, this.nsExcludes)})
  }

  config () {
    return $fetch<Config>('/api/config', { baseURL: this.baseURL })
  }

  targets () {
    return $fetch<{ [type: string]: Target[] }>(`/proxy/${this.cluster}/core/v2/targets`, { baseURL: this.baseURL })
  }

  namespaces (filter?: Filter) {
    return $fetch<string[]>(`/proxy/${this.cluster}/core/v1/namespaces`, { baseURL: this.baseURL, params: { ...filter } })
  }

  namespacedKinds (source?: string) {
    return $fetch<string[]>('/proxy/'+this.cluster+'/core/v1/namespaced-resources/kinds', { baseURL: this.baseURL, params: { sources: source ? [source] : undefined } })
  }

  namespacedStatusCount (source: string, filter?: Filter) {
    return $fetch<NamespaceStatusCount>(`/proxy/${this.cluster}/core/v2/namespace-scoped/${source}/status-counts`, { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes) } })
  }

  namespacedResults (filter?: Filter, pagination?: Pagination) {
    return $fetch<ResultList>('/proxy/'+this.cluster+'/core/v1/namespaced-resources/results', { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes), ...pagination } })
  }

  statusCount (source: string, filter?: Filter) {
    return $fetch<{ [status in Status]: number }>(`/proxy/${this.cluster}/core/v2/namespace-scoped/${source}/status-counts`, { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.clusterExcludes) } })
  }

  clusterKinds (source?: string) {
    return $fetch<string[]>('/proxy/'+this.cluster+'/core/v1/cluster-resources/kinds', { baseURL: this.baseURL, params: { sources: source ? [source] : undefined } })
  }

  clusterResults (filter?: Filter, pagination?: Pagination) {
    return $fetch<ResultList>('/proxy/'+this.cluster+'/core/v1/cluster-resources/results', { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes), ...pagination } })
  }

  resultsWithoutResources (filter?: Filter, pagination?: Pagination) {
    return $fetch<ResultList>('/proxy/'+this.cluster+'/core/v2/results-without-resources', { baseURL: this.baseURL, params: { ...filter, ...pagination } })
  }

  countFindings (filter?: Filter) {
    return $fetch<FindingCounts>('/proxy/'+this.cluster+'/core/v2/findings', { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes) } })
  }

  namespacedResourceResults (filter?: Filter, pagination?: Pagination) {
    return $fetch<ResourceResultList>('/proxy/'+this.cluster+'/core/v2/namespace-scoped/resource-results', { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.nsExcludes), ...pagination } })
  }

  clusterResourceResults (filter?: Filter, pagination?: Pagination) {
    return $fetch<ResourceResultList>('/proxy/'+this.cluster+'/core/v2/cluster-scoped/resource-results', { baseURL: this.baseURL, params: { ...applyExcludes(filter, this.clusterExcludes), ...pagination } })
  }

  resourceResults (id: string, filter?: Filter) {
    return $fetch<ResourceResult[]>(`/proxy/${this.cluster}/core/v2/resource/${id}/resource-results`, { baseURL: this.baseURL, params: filter })
  }

  resourceStatusCount (id: string, filter?: Filter) {
    return $fetch<ResourceStatusCount[]>(`/proxy/${this.cluster}/core/v2/resource/${id}/status-counts`, { baseURL: this.baseURL, params: filter })
  }

  results (id: string, pagination?: Pagination, filter?: Filter) {
    return $fetch<ResultList>(`/proxy/${this.cluster}/core/v2/resource/${id}/results`, { baseURL: this.baseURL, params: { ...pagination, ...filter } })
  }

  sources (filter?: Filter) {
    return $fetch<string[]>('/proxy/'+this.cluster+'/core/v2/sources', { baseURL: this.baseURL, params: applyExcludes(filter, [...this.nsExcludes, ...this.clusterExcludes]) })
  }

  categoryTree (id?: string, filter?: Filter) {
    return $fetch<Source[]>('/proxy/'+this.cluster+'/core/v2/sources/categories', { baseURL: this.baseURL, params: { id, ...applyExcludes(filter, [...this.nsExcludes, ...this.clusterExcludes]) } })
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

const applyExcludes = <T extends Filter>(filter: T | undefined, exclude: string[] | undefined) => {
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
