import Vue from 'vue';
import Vuex from 'vuex';

import {
  Target, ClusterPolicyReport, PolicyReport, NamespacePolicyMap, GlobalPolicyReportMap,
} from '@/models';
import api from '@/api';
import { convertPolicyReports, generateGlobalPolicyReports } from '@/mapper';

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
  namespacePolicyMap: NamespacePolicyMap;
  globalPolicyMap: GlobalPolicyReportMap;
  namespaces: string[];
}

export default new Vuex.Store<State>({
  state: {
    targets: [],
    reports: [],
    clusterReports: [],
    namespacePolicyMap: {},
    globalPolicyMap: {},
    namespaces: [],
  },
  mutations: {
    [SET_TARGETS]: (state, targets: Target[]) => {
      state.targets = targets;
    },
    [SET_REPORTS]: (state, reports: PolicyReport[]) => {
      const namespacePolicyMap = convertPolicyReports(reports);

      state.reports = reports;
      state.namespaces = Object.keys(namespacePolicyMap);
      state.namespacePolicyMap = namespacePolicyMap;
      state.globalPolicyMap = generateGlobalPolicyReports(reports);
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
