import { NuxtAxiosInstance } from '@nuxtjs/axios'
import { Policy, PolicyGroups } from './types'

export const create = ($axios: NuxtAxiosInstance) => ({
  async policies (): Promise<{ policies: Policy[]; groups: PolicyGroups }> {
    const policies = await $axios.$get<Policy[]>('/api/kyverno/policies')

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
})
