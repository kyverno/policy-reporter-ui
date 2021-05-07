import Vue from 'vue';
import Vuex from 'vuex';
import {
  Target, ClusterPolicyReport, PolicyReport, NamespacePolicyMap, GlobalPolicyReportMap, Result,
} from '@/models';
import api from '@/api';
import { convertPolicyReports, generateGlobalPolicyReports } from '@/mapper';
import kyverno, { NAMESPACE } from '@/plugins/kyverno/store';

Vue.use(Vuex);

export const SET_TARGETS = 'SET_TARGETS';
export const SET_CLUSTER_REPORTS = 'SET_CLUSTER_REPORTS';
export const SET_REPORTS = 'SET_REPORTS';
export const SET_LOG = 'SET_RESULT_LOG';
export const SET_PLUGINS = 'SET_PLUGINS';

export const FETCH_TARGETS = 'FETCH_TARGETS';
export const FETCH_CLUSTER_REPORTS = 'FETCH_CLUSTER_REPORTS';
export const FETCH_REPORTS = 'FETCH_REPORTS';
export const FETCH_LOG = 'FETCH_RESULT_LOG';
export const FETCH_PLUGINS = 'FETCH_PLUGINS';

export type State = {
  targets: Target[];
  log: Result[];
  reports: PolicyReport[];
  clusterReports: ClusterPolicyReport[];
  namespacePolicyMap: NamespacePolicyMap;
  globalPolicyMap: GlobalPolicyReportMap;
  namespaces: string[];
  plugins: string[];
}

export default new Vuex.Store<State>({
  modules: { [NAMESPACE]: kyverno },
  state: {
    log: [],
    targets: [],
    reports: [],
    clusterReports: [],
    namespacePolicyMap: {},
    globalPolicyMap: {},
    namespaces: [],
    plugins: [],
  },
  mutations: {
    [SET_PLUGINS]: (state, plugins: string[]) => {
      state.plugins = plugins;
    },
    [SET_TARGETS]: (state, targets: Target[]) => {
      state.targets = targets;
    },
    [SET_LOG]: (state, log: Result[]) => {
      state.log = log;
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
    [FETCH_PLUGINS]: ({ commit }) => {
      api.plugins().then((plugins) => commit(SET_PLUGINS, plugins));
    },
    [FETCH_TARGETS]: ({ commit }) => {
      api.targets().then((targets) => commit(SET_TARGETS, targets));
    },
    [FETCH_LOG]: ({ commit }) => {
      api.log().then((log) => commit(SET_LOG, log));
    },
    [FETCH_REPORTS]: ({ commit }) => {
      api.policyReports().then((reports) => commit(SET_REPORTS, reports));
    },
    [FETCH_CLUSTER_REPORTS]: ({ commit }) => {
      api.clusterPolicyReports().then((reports) => commit(SET_CLUSTER_REPORTS, reports));
    },
  },
});
