import { defineStore } from 'pinia'

import type { Cluster, Config } from '~/modules/core/types'

type State = {
    plugins: string[];
    clusters: Cluster[];
    clusterSources: string[];
    namespaceSources: string[];
}

export const useConfigStore = defineStore('config', {
  state: (): State => ({
    plugins: [],
    clusters: [],
    clusterSources: [],
    namespaceSources: [],
  }),
  getters: {
    multiCluster: (state: State) => state.clusters.length > 0,
  },
  actions: {
    setConfig(config: Config) {
      this.clusters = config.clusters
    }
  }
});
