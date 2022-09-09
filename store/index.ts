import { Cluster, DisplayMode, ViewsCofig } from '~/policy-reporter-plugins/core/types'

type State = {
    refreshInterval: number;
    displayMode: DisplayMode;
    viewsConfig: ViewsCofig;
    clusters: Cluster[]
    currentCluster?: Cluster
}

export const state = (): State => ({
  currentCluster: undefined,
  // 2 ** 31 - 1 is the limit for setInterval
  refreshInterval: 2 ** 31 - 1,
  displayMode: sessionStorage.getItem('displayMode') as DisplayMode || DisplayMode.LIGHT,
  clusters: [],
  viewsConfig: {
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
})

export const getters = {
  refreshInterval: (state: State) => state.refreshInterval,
  isDarkMode: (state: State) => state.displayMode === DisplayMode.DARK,
  displayMode: (state: State) => state.displayMode,
  viewsConfig: (state: State) => state.viewsConfig,
  dashboardConfig: (state: State) => state.viewsConfig.dashboard,
  clusters: (state: State) => state.clusters,
  multiCluster: (state: State) => state.clusters.length > 0,
  currentCluster: (state: State) => state.currentCluster
}

export const mutations = {
  setRefreshInterval (state: State, interval: number) {
    state.refreshInterval = interval
  },
  setCluster (state: State, cluster?: Cluster) {
    state.currentCluster = cluster
  },
  setClusters (state: State, clusters: Cluster[]) {
    state.clusters = clusters
    state.currentCluster = clusters.find(c => c.id === '')
  },
  setViewsConfig (state: State, config: ViewsCofig) {
    state.viewsConfig = config
  },
  setDisplayMode (state: State, mode: DisplayMode) {
    sessionStorage.setItem('displayMode', mode)
    state.displayMode = mode
  }
}
