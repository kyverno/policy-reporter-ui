import { DisplayMode, ViewsCofig } from '~/policy-reporter-plugins/core/types'

type State = {
    refreshInterval: number;
    displayMode: DisplayMode;
    viewsConfig: ViewsCofig;
}

export const state = (): State => ({
  refreshInterval: 10000,
  displayMode: sessionStorage.getItem('displayMode') as DisplayMode || DisplayMode.LIGHT,
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
  dashboardConfig: (state: State) => state.viewsConfig.dashboard
}

export const mutations = {
  setRefreshInterval (state: State, interval: number) {
    state.refreshInterval = interval
  },
  setViewsConfig (state: State, config: ViewsCofig) {
    state.viewsConfig = config
  },
  setDisplayMode (state: State, mode: DisplayMode) {
    sessionStorage.setItem('displayMode', mode)
    state.displayMode = mode
  }
}
