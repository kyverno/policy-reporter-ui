<template>
  <v-app v-if="!loading">
    <app-navigation v-model="drawer" :plugins="config.plugins" :views="config.views" />

    <v-app-bar app class="elevation-1" clipped-left>
      <v-app-bar-nav-icon @click="drawer = !drawer" />

      <v-toolbar-title>Policy Reporter</v-toolbar-title>

      <v-spacer />

      <template v-if="$route.name !== 'kyverno-plugin-uid'">
        <target-chip v-for="target in targets" :key="target.name" :target="target" />
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
import { DisplayMode, Target, Config } from '~/policy-reporter-plugins/core/types'

type Data = {
  loading: boolean;
  targets: Target[];
  drawer: boolean;
  config: Config;
  links: string[][];
}

type Methdos = {
  setDisplayMode(mode: DisplayMode): void;
}

export default Vue.extend<Data, Methdos, {}, {}>({
  components: { AppNavigation },
  data: () => ({
    loading: true,
    targets: [],
    config: {
      plugins: [],
      displayMode: '',
      views: {
        logs: true,
        policyReports: true,
        clusterPolicyReports: true,
        kyvernoPolicies: true,
        kyvernoVerifyImages: true
      }
    },
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
        this.config = config

        if (sessionStorage.getItem('displayMode')) {
          return
        }

        if (this.config.displayMode) {
          this.setDisplayMode(this.config.displayMode)
          return
        }

        if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
          this.setDisplayMode(DisplayMode.DARK)
        }
      }).then().finally(() => { this.loading = false })
    ])
  },
  computed: mapGetters(['isDarkMode']),
  watch: {
    isDarkMode: {
      immediate: true,
      handler (isDarkMode: boolean) {
        this.$vuetify.theme.dark = isDarkMode
      }
    }
  },
  beforeMount () {
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
      this.setDisplayMode(!e.matches ? DisplayMode.LIGHT : DisplayMode.DARK)
    })
  },
  methods: mapMutations(['setDisplayMode'])
})
</script>

<style>
.apexcharts-svg {
  background-color: transparent !important;
}
</style>
