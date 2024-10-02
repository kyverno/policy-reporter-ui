import {defineStore} from 'pinia'

import {type Cluster, type Config, DisplayMode, type Navigation} from '~/modules/core/types'

const DisplayModeKey = 'dm'

type State = {
    displayMode: DisplayMode;
    plugins: string[];
    clusters: Cluster[];
    oauth: boolean;
    banner: string;
    error?: Error;
}

export const useConfigStore = defineStore('config', {
    state: (): State => ({
        displayMode: sessionStorage.getItem(DisplayModeKey) as DisplayMode,
        plugins: [],
        clusters: [],
        oauth: false,
        banner: '',
    }),
    getters: {
        multiCluster: (state: State) => state.clusters.length > 0,
        theme: (state: State) => {
            let mode = state.displayMode

            if (!Object.values(DisplayMode).includes(mode)) {
                return preferredDisplayMode()
            }

            return mode
        }
    },
    actions: {
        setConfig(config: Config) {
            this.error = config.error
            this.clusters = config.clusters
            this.oauth = config.oauth
            this.banner = config.banner
        },
        setDisplayMode(mode: DisplayMode) {
            if (!mode || mode === this.displayMode) return;

            if (!Object.values(DisplayMode).includes(mode)) {
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
