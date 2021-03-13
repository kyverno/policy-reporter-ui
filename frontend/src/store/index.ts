import Vue from 'vue';
import Vuex from 'vuex';

import { Target, ClusterPolicyReport, PolicyReport } from '@/models';
import api from '@/api';

Vue.use(Vuex);

export const SET_TARGETS = 'SET_TARGETS';
export const SET_CLUSTER_REPORTS = 'SET_CLUSTER_REPORTS';
export const SET_REPORTS = 'SET_REPORTS';

export const FETCH_TARGETS = 'FETCH_TARGETS';
export const FETCH_CLUSTER_REPORTS = 'FETCH_CLUSTER_REPORTS';
export const FETCH_REPORTS = 'FETCH_REPORTS';

export type State = {
  targets: Target[];
  reports: PolicyReport[];
  clusterReports: ClusterPolicyReport[];
}

export default new Vuex.Store<State>({
  state: {
    targets: [],
    reports: [],
    clusterReports: [],
  },
  mutations: {
    [SET_TARGETS]: (state, targets: Target[]) => {
      state.targets = targets;
    },
    [SET_REPORTS]: (state, reports: PolicyReport[]) => {
      state.reports = reports;
    },
    [SET_CLUSTER_REPORTS]: (state, clusterReports: ClusterPolicyReport[]) => {
      state.clusterReports = clusterReports;
    },
  },
  actions: {
    [FETCH_TARGETS]: ({ commit }) => {
      api.targets().then((targets) => commit(SET_TARGETS, targets));
    },
    [FETCH_REPORTS]: ({ commit }) => {
      api.policyReports().then((reports) => commit(SET_REPORTS, reports));
    },
    [FETCH_CLUSTER_REPORTS]: ({ commit }) => {
      api.clusterPolicyReports().then((reports) => commit(SET_CLUSTER_REPORTS, reports));
    },
  },
  modules: {
  },
});
