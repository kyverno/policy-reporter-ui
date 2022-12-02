<template>
  <v-app v-if="!loading">
    <app-navigation v-model="drawer" :plugins="plugins" :views="viewsConfig" />

    <v-app-bar app class="elevation-1" clipped-left>
      <v-app-bar-nav-icon @click="drawer = !drawer" />

      <v-toolbar-title>Policy Reporter</v-toolbar-title>

      <v-spacer />

      <template v-if="$route.name !== 'kyverno-plugin-uid'">
        <target-chip v-for="target in targets" :key="target.name" :target="target" />
        <cluster-select v-if="multiCluster" class="ml-3" />
      </template>

      <refresh-interval-select class="ml-3" />

      <display-mode-select class="ml-3" />
    </v-app-bar>

    <v-main>
      <Nuxt :key="$route.path" />
    </v-main>
  </v-app>

  <v-app v-else />
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters, mapMutations } from 'vuex'
import AppNavigation from '~/components/AppNavigation.vue'
import { Cluster, DisplayMode, Target, ViewsCofig } from '~/policy-reporter-plugins/core/types'

type Data = {
  loading: boolean;
  targets: Target[];
  drawer: boolean;
  plugins: string[];
  links: string[][];
}

type Methdos = {
  setDisplayMode(mode: DisplayMode): void;
  setRefreshInterval(interval: number): void;
  setViewsConfig(config: ViewsCofig): void;
  setClusters(clusters: Cluster[]): void;
  setLabelFilter(filter: string[]): void;
}

type Computed = {
  viewsConfig: ViewsCofig;
  isDarkMode: boolean;
  multiCluster: boolean;
  currentCluster?: Cluster;
}

export default Vue.extend<Data, Methdos, Computed, {}>({
  components: { AppNavigation },
  data: () => ({
    loading: true,
    targets: [],
    plugins: [],
    drawer: true,
    links: [
      ['mdi-view-dashboard', 'Dashboard', '/'],
      ['mdi-file-chart', 'Policy Reports', '/policy-reports'],
      ['mdi-file-chart', 'ClusterPolicy Reports', '/cluster-policy-reports'],
      ['mdi-console', 'Logs', '/logs'],
      ['', 'Kyverno Policies', '/kyverno-plugin', 'kyverno', 'lazy-kyverno-icon']
    ]
  }),
  async fetch () {
    await Promise.all([
      this.$coreAPI.targets().then((targets) => {
        this.targets = targets
      }),
      this.$coreAPI.config().then((config) => {
        this.plugins = config.plugins
        this.setViewsConfig(config.views)
        this.setClusters(config.clusters)
        this.setRefreshInterval(config.refreshInterval)
        this.setLabelFilter(config.labelFilter)

        if (sessionStorage.getItem('displayMode')) {
          return
        }

        if (config.displayMode) {
          this.setDisplayMode(config.displayMode)
          return
        }

        if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
          this.setDisplayMode(DisplayMode.DARK)
        }
      }).then().finally(() => { this.loading = false })
    ])
  },
  computed: mapGetters(['isDarkMode', 'viewsConfig', 'multiCluster', 'currentCluster']),
  watch: {
    isDarkMode: {
      immediate: true,
      handler (isDarkMode: boolean) {
        this.$vuetify.theme.dark = isDarkMode
      }
    },
    currentCluster (cluster?: Cluster) {
      this.$coreAPI.targets().then((targets) => {
        this.targets = targets
      })

      if (!cluster) { return }

      if (!cluster.kyverno && this.$route.name?.startsWith('kyverno-plugin')) {
        this.$router.push('/')
      }
    }
  },
  beforeMount () {
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
      this.setDisplayMode(!e.matches ? DisplayMode.LIGHT : DisplayMode.DARK)
    })
  },
  methods: mapMutations(['setDisplayMode', 'setViewsConfig', 'setClusters', 'setRefreshInterval', 'setLabelFilter'])
})
</script>

<style>
.apexcharts-svg {
  background-color: transparent !important;
}
</style>
