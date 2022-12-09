import { NuxtAxiosInstance } from '@nuxtjs/axios'
import { Policy, PolicyGroups, VerifyImageRule, KyvernoAPI, Filter } from './types'

class API {
  private axios: NuxtAxiosInstance
  private prefix: string = ''

  constructor (axios: NuxtAxiosInstance) {
    this.axios = axios
  }

  async policies (): Promise<{ policies: Policy[]; groups: PolicyGroups }> {
    const policies = await this.axios.$get<Policy[]>(this.prefix + '/kyverno/policies')

    const unsorted = policies.reduce<PolicyGroups>((groups, policy) => {
      if (!policy.category) {
        groups['No Category'].push(policy)

        return groups
      }

      if (!groups[policy.category]) {
        return { ...groups, [policy.category]: [policy] }
      }

      groups[policy.category].push(policy)

      return groups
    }, { 'No Category': [] })

    const groups = Object.keys(unsorted).sort().reduce<PolicyGroups>((acc, key) => {
      if (unsorted[key].length === 0) { return acc }

      acc[key] = unsorted[key]

      return acc
    }, {})

    return { policies, groups }
  }

  verifyImageRules (): Promise<VerifyImageRule[]> {
    return this.axios.$get<VerifyImageRule[]>(this.prefix + '/kyverno/verify-image-rules')
  }

  policyReport (filter: Filter): Promise<BlobPart> {
    return this.axios.$get<BlobPart>(this.prefix + '/kyverno/policy-details-reporting', { responseType: 'blob', params: filter })
  }

  namespaceReport (filter: Filter): Promise<BlobPart> {
    return this.axios.$get<BlobPart>(this.prefix + '/kyverno/namespace-details-reporting', { responseType: 'blob', params: filter })
  }

  setPrefix (prefix: string): void {
    this.prefix = prefix
  }
}

export const create = ($axios: NuxtAxiosInstance): KyvernoAPI => new API($axios)
