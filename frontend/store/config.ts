import { defineStore } from 'pinia'

import { type Cluster, type Config, DisplayMode, type SourceConfig } from '~/modules/core/types'

const DisplayModeKey = 'dm'

type State = {
    displayMode: DisplayMode;
    plugins: string[];
    clusters: Cluster[];
    sources: SourceConfig[];
    clusterSources: string[];
    namespaceSources: string[];
}

export const useConfigStore = defineStore('config', {
  state: (): State => ({
    displayMode: sessionStorage.getItem(DisplayModeKey) as DisplayMode,
    plugins: [],
    clusters: [],
    sources: [],
    clusterSources: [],
    namespaceSources: [],
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
      this.sources = config.sources
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
