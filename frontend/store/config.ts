import { defineStore } from 'pinia'

import { type Cluster, type Config, DisplayMode, type SourceConfig } from '~/modules/core/types'

const DisplayModeKey = 'dm'

type State = {
    displayMode: DisplayMode;
    plugins: string[];
    clusters: Cluster[];
    oauth: boolean;
}

export const useConfigStore = defineStore('config', {
  state: (): State => ({
    displayMode: sessionStorage.getItem(DisplayModeKey) as DisplayMode,
    plugins: [],
    clusters: [],
    oauth: false,
  }),
  getters: {
    multiCluster: (state: State) => state.clusters.length > 0,
    theme: (state: State) => {
      let mode = state.displayMode

      if (![DisplayMode.UNSPECIFIED, DisplayMode.LIGHT, DisplayMode.DARK].includes(mode)) {
        return preferredDisplayMode()
      }

      return mode
    }
  },
  actions: {
    setConfig(config: Config) {
      this.clusters = config.clusters
      this.oauth = config.oauth
    },
    setDisplayMode(mode: DisplayMode) {
      if (!mode || mode === this.displayMode) return;

      if (![DisplayMode.DARK, DisplayMode.LIGHT].includes(mode)) {
        mode = preferredDisplayMode()
      }

      sessionStorage.setItem(DisplayModeKey, mode)
      this.displayMode = mode
    }
  }
});

const preferredDisplayMode = () => {
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
      return DisplayMode.DARK
    }

    return DisplayMode.LIGHT
}
