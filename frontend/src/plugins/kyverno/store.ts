import { Module } from 'vuex';
import api from '@/plugins/kyverno/api';
import { Policy, PolicyGroups } from './models';

export const SET_POLICIES = 'SET_KYVERNO_POLICIES';

export const FETCH_POLICIES = 'FETCH_KYVERNO_POLICIES';

export const NAMESPACE = 'kyverno';

export type State = {
  policies: Policy[];
}

const module: Module<State, any> = {
  namespaced: true,
  state: {
    policies: [],
  },
  getters: {
    policies: (state) => state.policies,
    policyGroups: (state) => {
      const unsorted = state.policies.reduce<PolicyGroups>((groups, policy) => {
        if (!policy.category) {
          groups['No Category'].push(policy);

          return groups;
        }

        if (!groups.hasOwnProperty(policy.category)) {
          return { ...groups, [policy.category]: [policy] };
        }

        groups[policy.category].push(policy);

        return groups;
      }, { 'No Category': [] });

      return Object.keys(unsorted).sort().reduce<PolicyGroups>((acc, key) => {
        acc[key] = unsorted[key];

        return acc;
      }, {});
    },
  },
  mutations: {
    [SET_POLICIES]: (state, policies: Policy[]) => {
      state.policies = policies;
    },
  },
  actions: {
    [FETCH_POLICIES]: ({ commit }) => {
      api.policies().then((policies) => commit(SET_POLICIES, policies));
    },
  },
};

export default module;
