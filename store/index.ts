import { DisplayMode } from '~/policy-reporter-plugins/core/types'

type State = {
    refreshInterval: number;
    displayMode: DisplayMode;
}

export const state = (): State => ({
  refreshInterval: 10000,
  displayMode: sessionStorage.getItem('displayMode') as DisplayMode || DisplayMode.LIGHT
})

export const getters = {
  refreshInterval: (state: State) => state.refreshInterval,
  isDarkMode: (state: State) => state.displayMode === DisplayMode.DARK,
  displayMode: (state: State) => state.displayMode
}

export const mutations = {
  setRefreshInterval (state: State, interval: number) {
    state.refreshInterval = interval
  },
  setDisplayMode (state: State, mode: DisplayMode) {
    sessionStorage.setItem('displayMode', mode)
    state.displayMode = mode
  }
}
